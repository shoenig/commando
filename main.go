// Author hoenig

// commando is a command line tool for executing expect scripts
// via ssh on remote servers. It is as terrible as it sounds.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh"
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

func publickey() (string, error) {
	keypath := path.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub")
	bs, err := ioutil.ReadFile(keypath)
	if err != nil {
		return "", errors.Wrapf(err, "could not read %s", keypath)
	}

	return strings.TrimSpace(string(bs)), err
}

func setupKeys(user, pass string, hosts []string) error {
	key, err := publickey()
	if err != nil {
		return err
	}

	fmt.Println("setting public key for hosts:", hosts)

	for _, host := range hosts {
		if err := setupKey(user, pass, host, key); err != nil {
			return err
		}
	}

	return nil
}

func setupKey(user, pass, host, key string) error {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
	}

	address := fmt.Sprintf("%s:22", host)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return errors.Wrap(err, "failed to dial server")
	}

	// 1. ensure .ssh directory exists
	if err := run(client, host, "mkdir -p ~/.ssh", true); err != nil {
		return errors.Wrap(err, "mkdir .ssh failed")
	}

	// 2. append key to authroized_keys (if it is not already present)
	appendCmd := fmt.Sprintf(`if grep -q "%s" ~/.ssh/authorized_keys; then echo "key already exists"; else echo "%s" >> ~/.ssh/authorized_keys; fi`, key, key)
	if err := run(client, host, appendCmd, true); err != nil {
		return errors.Wrap(err, "echo key failed")
	}

	if err := client.Close(); err != nil {
		return errors.Wrap(err, "failed to close client")
	}

	fmt.Println("")

	return nil
}

func run(client *ssh.Client, host, cmd string, output bool) error {
	session, err := client.NewSession()
	if err != nil {
		return errors.Wrap(err, "failed to create ssh session")
	}

	fmt.Println(host, "<<<", cmd)
	bs, err := session.CombinedOutput(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to run command")
	}

	if output {
		fmt.Println(host, ">>>", string(bs))
	}
	fmt.Println("")

	// normal to get non-nil EOF when cmd is complete
	_ = session.Close()

	return nil
}
