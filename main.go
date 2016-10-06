// Author hoenig

// commando is a command line tool for executing expect scripts
// via ssh on remote servers. It is as terrible as it sounds.
package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	user, setupkey, hosts, scripts := arguments()

	fmt.Printf("gimmie password: ")
	bs, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read password: %v\n", err)
	}
	fmt.Println("\t...thanks!")

	if setupkey {
		if err := setupKeys(user, string(bs), hosts); err != nil {
			fmt.Fprintf(os.Stderr, "failed to setup ssh keys: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("executing", scripts, "on hosts", hosts)
	}
}
