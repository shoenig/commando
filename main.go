// Author hoenig

// commando is a command line tool for executing expect scripts
// via ssh on remote servers. It is as terrible as it sounds.
package main

import "fmt"

func main() {
	hosts, scripts := arguments()

	fmt.Println("hosts:", hosts)
	fmt.Println("scripts:", scripts)
}
