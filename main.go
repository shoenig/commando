// Author hoenig

package main

import (
	"fmt"
	"os"
)

// commando is a command line tool for executing expect scripts
// via ssh on remote servers.
// It is as terrible as it sounds.

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "useage: ./commando [scripts directory]")
		os.Exit(1)
	}

	scripts, err := loadScripts(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load scripts: %v", err)
		os.Exit(1)
	}

	fmt.Printf("loaded %d scripts\n", len(scripts))
}
