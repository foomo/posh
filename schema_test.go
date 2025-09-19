package main_test

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	testingx "github.com/foomo/go/testing"
	tagx "github.com/foomo/go/testing/tag"
	"github.com/foomo/posh/pkg/config"
	"github.com/invopop/jsonschema"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Config struct {
	// Version of the config
	Version string `json:"version" jsonschema:"required,default=v1.0"`
	// Environment variables
	Env config.Env `json:"env"`
	// Prompt settings
	Prompt config.Prompt `json:"prompt"`
	// Require settings
	Require config.Require `json:"require"`
}

func TestSchema(t *testing.T) {
	t.Parallel()
	testingx.Tags(t, tagx.Short)

	cwd, err := os.Getwd()
	require.NoError(t, err)

	reflector := new(jsonschema.Reflector)
	reflector.AllowAdditionalProperties = true
	reflector.RequiredFromJSONSchemaTags = true
	require.NoError(t, reflector.AddGoComments("github.com/foomo/posh", "./"))
	schema := reflector.Reflect(&Config{})
	actual, err := json.MarshalIndent(schema, "", "  ")
	require.NoError(t, err)

	filename := path.Join(cwd, "posh.schema.json")
	expected, err := os.ReadFile(filename)
	if !errors.Is(err, os.ErrNotExist) {
		require.NoError(t, err)
	}

	if !assert.Equal(t, string(expected), string(actual)) {
		require.NoError(t, os.WriteFile(filename, actual, 0600))
	}
}
