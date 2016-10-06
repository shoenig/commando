// Author hoenig

// commando is a command line tool for executing expect scripts
// via ssh on remote servers. It is as terrible as it sounds.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

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

	return string(bs), err
}

func setupKeys(user, pass string, hosts []string) error {
	key, err := publickey()
	if err != nil {
		return err
	}
	fmt.Println("public key is:", key)

	for _, host := range hosts {
		if err := setupKey(user, pass, host, key); err != nil {
			return err
		}
	}

	return nil
}

func setupKey(user, pass, host, key string) error {
	fmt.Println("adding public ssh key to ~/.ssh/authorized_keys on:", host)

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

	fmt.Println("connected!")

	return client.Close()
}
