// Author hoenig

package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh/terminal"
)

func prompt(args args) (string, error) {
	if args.nopassword {
		tracef(args.verbose, "skipping password prompt")
		return "", nil
	}

	fmt.Printf("  password for '%s' --> ", args.user)
	bs, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", errors.Wrap(err, "failed to read password")
	}
	fmt.Println("\nok")
	return string(bs), nil
}
