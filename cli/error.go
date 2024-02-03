package cli

import (
	"errors"
)

var (
	ErrNoConfig       = errors.New("config.toml not found")
	ErrInvalidVersion = errors.New("requested version does not appear to be valid")
	ErrFailedUpgrade  = errors.New("failed to self-upgrade ovm")
)
