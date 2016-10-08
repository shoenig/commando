// Author hoenig

package main

import (
	"flag"
	"os"

	"github.com/pkg/errors"
)

type args struct {
	user       string
	hostexp    string
	scriptdir  string
	nopassword bool
	verbose    bool
}

func arguments() args {
	var args args

	flag.StringVar(&args.user, "user", os.Getenv("USER"), "ssh username")
	flag.StringVar(&args.hostexp, "hosts", "", "the list of hosts")
	flag.StringVar(&args.scriptdir, "scripts", "", "the directory full of scripts")
	flag.BoolVar(&args.nopassword, "no-password", false, "no-password skips password prompt")
	flag.BoolVar(&args.verbose, "verbose", false, "verbose mode")

	flag.Parse()

	return args
}

func validate(args args) error {
	if args.hostexp == "" {
		return errors.Errorf("--hosts is required")
	}

	if args.user == "" {
		return errors.Errorf("--user or $USER must be set")
	}

	if args.scriptdir == "" {
		return errors.Errorf("--scripts is required")
	}

	return nil
}

/*
	hosts := hosts(args.ostexp)

	if args.setupkey {
		return user, true, hosts, nil
	}

	scripts, err := loadScripts(args.scriptdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bad argument: %v", err)
		os.Exit(1)
	}
*/
