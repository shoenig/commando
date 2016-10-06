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
func arguments() ([]string, []script) {
	var hostexp string
	var scriptdir string

	flag.StringVar(&hostexp, "hosts", "", "the list of hosts")
	flag.StringVar(&scriptdir, "scripts", "", "the directory full of scripts")

	flag.Parse()

	if scriptdir == "" {
		fmt.Fprintln(os.Stderr, "--scripts is required")
		os.Exit(1)
	}

	if hostexp == "" {
		fmt.Fprintln(os.Stderr, "--hosts is required")
		os.Exit(1)
	}

	scripts, err := loadScripts(scriptdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bad argument: %v", err)
		os.Exit(1)
	}

	hosts := hosts(hostexp)
	return hosts, scripts
}
