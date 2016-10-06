// Author hoenig

package main

import (
	"flag"
	"fmt"
	"os"
)

// returns arguments from
// --hosts list of hosts to execute scripts against
// --scripts list of scripts (in order) to execute
func arguments() (string, bool, []string, []script) {
	var user string
	var hostexp string
	var scriptdir string
	var setupkey bool

	flag.StringVar(&user, "user", os.Getenv("USER"), "ssh username")
	flag.StringVar(&hostexp, "hosts", "", "the list of hosts")
	flag.BoolVar(&setupkey, "setupkey", false, "set ssh public key in authorized_keys on remote")
	flag.StringVar(&scriptdir, "scripts", "", "the directory full of scripts")

	flag.Parse()

	if scriptdir != "" && setupkey {
		fmt.Fprintln(os.Stderr, "only one of --scripts and --setupkey at a time")
		os.Exit(1)
	}

	if scriptdir == "" && !setupkey {
		fmt.Fprintln(os.Stderr, "one of --scripts or --setupkey is required")
		os.Exit(1)
	}

	if hostexp == "" {
		fmt.Fprintln(os.Stderr, "--hosts is required")
		os.Exit(1)
	}

	hosts := hosts(hostexp)

	if setupkey {
		return user, true, hosts, nil
	}

	scripts, err := loadScripts(scriptdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bad argument: %v", err)
		os.Exit(1)
	}

	return user, false, hosts, scripts
}
