// Author hoenig

// commando is a command line tool for executing expect scripts
// via ssh on remote servers. It is as terrible as it sounds.
package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// example seting up ssh keys
// ./commando --setupkey --hosts "prod-executor{8..38}"

func main() {
	user, setupkey, hosts, scripts := arguments()

	fmt.Printf("\tgimmie password:  ")
	bs, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read password: %v\n", err)
	}
	fmt.Println("\t...thanks!")
	password := string(bs)

	if setupkey {
		if err := setupKeys(user, password, hosts); err != nil {
			fmt.Fprintf(os.Stderr, "failed to setup ssh keys: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("executing", scripts, "on hosts", hosts)

		if err := runScripts(user, password, hosts, scripts); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run scripts: %v\n", err)
			os.Exit(1)
		}
	}
}
