package onepassword

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/1Password/connect-sdk-go/connect"
	"github.com/1Password/connect-sdk-go/onepassword"
	"github.com/foomo/posh/pkg/cache"
	"github.com/foomo/posh/pkg/log"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type (
	OnePassword struct {
		l             log.Logger
		cache         cache.Namespace
		connect       connect.Client
		uuidRegex     *regexp.Regexp
		tokenFilename string
		watching      map[string]bool
		last          time.Time
	}
	Option func(*OnePassword) error
)

// ------------------------------------------------------------------------------------------------
// ~ Options
// ------------------------------------------------------------------------------------------------

func WithTokenFilename(v string) Option {
	return func(o *OnePassword) error {
		o.tokenFilename = v
		return nil
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func New(l log.Logger, cache cache.Cache, opts ...Option) (*OnePassword, error) {
	inst := &OnePassword{
		l:         l,
		cache:     cache.Get("onePasswordClient"),
		uuidRegex: regexp.MustCompile(`^[a-z0-9]{26}$`),
		watching:  map[string]bool{},
	}
	for _, opt := range opts {
		if opt != nil {
			if err := opt(inst); err != nil {
				return nil, err
			}
		}
	}
	if client, err := connect.NewClientFromEnvironment(); err != nil {
		l.Debug("not able to create connect client", err.Error())
	} else {
		inst.connect = client
	}
	return inst, nil
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (op *OnePassword) Session(account string) (bool, error) {
	var sessChanged bool
	sess := os.Getenv("OP_SESSION_" + account)

	if op.tokenFilename != "" {
		if err := godotenv.Overload(op.tokenFilename); err != nil {
			op.l.Debug("could not load session from env file")
			sessChanged = true
		} else if value := os.Getenv("OP_SESSION_" + account); sess != value {
			sessChanged = true
		} else {
			op.l.Debug("INFO: loaded op session from file " + value)
		}
	}

	if sessChanged || op.last.IsZero() || time.Since(op.last) > time.Minute*10 {
		out, err := exec.Command("op", "account", "--account", account, "get", "--format", "json").Output()
		if err != nil {
			return false, fmt.Errorf("%w: %s", err, string(out))
		}

		var data struct {
			Name string `json:"name"`
		}

		if err := json.Unmarshal(out, &data); err != nil {
			return false, err
		}

		if data.Name == account {
			op.last = time.Now()
			op.watch(account)
			return true, nil
		}
	}
	return true, nil
}

func (op *OnePassword) SignIn(ctx context.Context, account string) error {
	if ok, _ := op.Session(account); ok {
		return nil
	}

	// create command
	cmd := exec.CommandContext(ctx,
		"op", "signin",
		"--account", account,
		"--raw",
	)

	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stdin = os.Stdin

	// start the process and wait till it's finished
	if err := cmd.Start(); err != nil {
		return err
	} else if err := cmd.Wait(); err != nil {
		return err
	}

	token := strings.TrimSuffix(stdoutBuf.String(), "\n")
	if token == "" {
		return errors.New("failed to retrieve 1password token!")
	} else if err := os.Setenv(fmt.Sprintf("OP_SESSION_%s", account), token); err != nil {
		return err
	} else {
		op.l.Infof(`If you need op outside the shell, run:

$ export OP_SESSION_%s=%s
`, account, token)
	}

	if op.tokenFilename != "" {
		if err := os.MkdirAll(path.Dir(op.tokenFilename), os.ModePerm); err != nil {
			return err
		} else if err := os.WriteFile(op.tokenFilename, []byte(fmt.Sprintf("OP_SESSION_%s=%s\n", account, token)), 0600); err != nil {
			return err
		} else {
			op.l.Infof(`Session env has been stored for your convenience at:

%s
`, op.tokenFilename)
		}
	}
	op.watch(account)
	return nil
}

func (op *OnePassword) Get(ctx context.Context, account, vaultUUID, itemUUID, field string) (string, error) {
	if op.connect != nil {
		if fields := op.connectGet(vaultUUID, itemUUID); len(fields) == 0 {
			return "", fmt.Errorf("could not find secret '%s' '%s'", vaultUUID, itemUUID)
		} else if value, ok := fields[field]; !ok {
			return "", fmt.Errorf("could not find field %s", field)
		} else {
			return strings.ReplaceAll(strings.TrimSpace(value), "\\n", "\n"), nil
		}
	} else {
		if ok, _ := op.Session(account); !ok {
			return "", ErrNotSignedIn
		} else if fields := op.clientGet(ctx, vaultUUID, itemUUID); len(fields) == 0 {
			return "", fmt.Errorf("could not find secret '%s' '%s'", vaultUUID, itemUUID)
		} else if value, ok := fields[field]; !ok {
			return "", fmt.Errorf("could not find field %s", field)
		} else {
			return strings.ReplaceAll(strings.TrimSpace(value), "\\n", "\n"), nil
		}
	}
}

func (op *OnePassword) GetOnetimePassword(ctx context.Context, account, uuid string) (string, error) {
	if ok, _ := op.Session(account); !ok {
		return "", ErrNotSignedIn
	}

	out, err := exec.CommandContext(ctx,
		"op", "item", "get", "--otp", uuid,
	).Output()
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(strings.TrimSpace(string(out)), "\\n", "\n"), nil
}

func (op *OnePassword) Render(ctx context.Context, source string) ([]byte, error) {
	tpl, err := template.New("1password").
		Delims("<% ", " %>").
		Option("missingkey=error").
		Funcs(
			template.FuncMap{
				"env": func(name string) (string, error) {
					value := os.Getenv(name)
					if value == "" {
						return "", fmt.Errorf("env variable %q was empty", name)
					}
					return value, nil
				},
				"op": func(account, vaultUUID, itemUUID, field string) (string, error) {
					return op.Get(ctx, account, vaultUUID, itemUUID, field)
				},
				"indent": func(spaces int, v string) string {
					pad := strings.Repeat(" ", spaces)
					return strings.ReplaceAll(v, "\n", "\n"+pad)
				},
				"quote": func(v string) string {
					return "'" + v + "'"
				},
				"replace": func(o, n, v string) string {
					return strings.ReplaceAll(v, o, n)
				},
				"base64": func(v string) string {
					return base64.StdEncoding.EncodeToString([]byte(v))
				},
			},
		).
		Parse(source)
	if err != nil {
		return nil, err
	}

	out := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(out, nil); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func (op *OnePassword) RenderFile(ctx context.Context, source string) ([]byte, error) {
	in, err := os.ReadFile(source)
	if err != nil {
		return nil, err
	}
	out, err := op.Render(ctx, string(in))
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (op *OnePassword) RenderFileTo(ctx context.Context, source, target string) error {
	out, err := op.RenderFile(ctx, source)
	if err != nil {
		return err
	}
	value := fmt.Sprintf(
		"# Code generated by shell %s - DO NOT EDIT.\n%s",
		time.Now().Format("2006-01-02 15:04:05"),
		string(out),
	)
	return os.WriteFile(target, []byte(value), 0600)
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

//nolint:forcetypeassert
func (op *OnePassword) clientGet(ctx context.Context, vaultUUID string, itemUUID string) map[string]string {
	return op.cache.Get(strings.Join([]string{vaultUUID, itemUUID}, "#"), func() any {
		ret := map[string]string{}
		var v struct {
			Vault struct {
				ID string `json:"id"`
			} `json:"vault"`
			Fields []struct {
				ID    string      `json:"id"`
				Type  string      `json:"type"` // CONCEALED, STRING
				Label string      `json:"label"`
				Value interface{} `json:"value"`
			} `json:"fields"`
		}
		if res, err := exec.CommandContext(ctx,
			"op", "item", "get", itemUUID, "--format", "json",
		).CombinedOutput(); err != nil {
			op.l.Error("failed to retrieve item", err.Error())
			return ret
		} else if err := json.Unmarshal(res, &v); err != nil {
			op.l.Error("failed to retrieve item", err.Error())
			return ret
		} else if v.Vault.ID != vaultUUID {
			op.l.Error("failed to retrieve item", "wrong vault UUID %s for item %s", vaultUUID, itemUUID)
			return ret
		} else {
			ret := map[string]string{}
			aliases := map[string]string{
				"notesPlain": "notes",
			}
			for _, field := range v.Fields {
				if alias, ok := aliases[field.Label]; ok {
					ret[alias] = fmt.Sprintf("%v", field.Value)
				} else {
					ret[field.Label] = fmt.Sprintf("%v", field.Value)
				}
			}
			return ret
		}
	}).(map[string]string)
}

//nolint:forcetypeassert
func (op *OnePassword) connectGet(vaultUUID, itemUUID string) map[string]string {
	return op.cache.Get(strings.Join([]string{vaultUUID, itemUUID}, "#"), func() any {
		ret := map[string]string{}
		var item *onepassword.Item
		if op.uuidRegex.Match([]byte(itemUUID)) {
			if v, err := op.connect.GetItem(itemUUID, vaultUUID); err != nil {
				op.l.Error("failed to retrieve item:", err.Error())
				return ret
			} else {
				item = v
			}
		} else {
			if v, err := op.connect.GetItemByTitle(itemUUID, vaultUUID); err != nil {
				op.l.Error("failed to retrieve item by title:", err.Error())
				return ret
			} else {
				item = v
			}
		}
		for _, f := range item.Fields {
			ret[f.Label] = f.Value
		}
		return ret
	}).(map[string]string)
}

func (op *OnePassword) watch(account string) {
	if v, ok := op.watching[account]; !ok || !v {
		go func() {
			for {
				if ok, err := op.Session(account); err != nil {
					op.l.Warnf("\n1password session keep alive failed for '%s' (%s)", account, err.Error())
					op.watching[account] = false
					return
				} else if !ok {
					op.l.Warnf("\n1password session keep alive failed for '%s'", account)
					op.watching[account] = false
					return
				}
				time.Sleep(time.Minute * 15)
			}
		}()
		op.watching[account] = true
	}
}
