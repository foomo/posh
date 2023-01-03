package onepassword

import (
	"github.com/pkg/errors"
)

var ErrNotSignedIn = errors.New("you're not signed into your 1password account")
