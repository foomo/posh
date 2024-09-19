package util

import (
	"bufio"
	"fmt"
	"os"
)

func Prompt(txt string) (string, error) {
	r := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fprint(os.Stderr, txt+": "); err != nil {
		return "", err
	}
	out, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return out, nil
}
