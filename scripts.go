// Author hoenig

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/pkg/errors"
)

const filenameFmt = `[0-9]+-`

var filenameRe = regexp.MustCompile(filenameFmt)

// script is an expect script with an ordered name.
// script names must start with a number and a dash.
type script struct {
	name    string
	command string
	stdin   []string
}

func (s script) String() string {
	return s.name
}

func load(cfg args) ([]script, error) {
	scripts := []script{}

	err := filepath.Walk(cfg.scriptdir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "failed to read scripts")
		}

		// skip directories
		if info.IsDir() {
			return nil
		}

		if !filenameRe.MatchString(info.Name()) {
			return errors.Errorf("script name must start with ([0-9]+)-")
		}

		script, err := read(info.Name(), path)
		if err != nil {
			return errors.Wrapf(err, "failed to read script file %s", info.Name())
		}

		scripts = append(scripts, script)
		return nil
	})

	if len(scripts) == 0 {
		return nil, errors.Errorf("no scripts found")
	}

	return scripts, err
}

func read(name, path string) (script, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return script{}, errors.Wrap(err, "failed to read script")
	}
	s := strings.TrimSpace(string(bs))
	return parse(name, s)
}

func parse(name, content string) (script, error) {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return script{}, errors.Errorf("no command in script %s", name)
	}
	s := script{
		name:    name,
		command: lines[0],
		stdin:   lines[1:],
	}
	return s, nil
}

func run(user, pass string, hosts []string, scripts []script) error {
	for _, host := range hosts {

		client, err := makeClient(user, pass, host)
		if err != nil {
			return errors.Wrap(err, "failed to dial host")
		}

		for _, script := range scripts {
			if err := execute(client, user, pass, host, script); err != nil {
				return errors.Wrapf(err, "failed to run %s on %s", script, host)
			}
			fmt.Println("")
		}
	}
	return nil
}

func substitute(stdin []string, substitutions map[string]string) []string {
	replaced := []string{}
	for _, line := range stdin {
		for old, new := range substitutions {
			line = strings.Replace(line, old, new, -1)
		}
		replaced = append(replaced, line)
	}
	return replaced
}

func combine(stdin []string) string {
	var b bytes.Buffer
	for _, line := range stdin {
		line = strings.TrimSpace(line)
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}

func execute(client *ssh.Client, user, pass, host string, script script) error {

	fmt.Printf("###### running script %s on %s at %s\n", script, host, time.Now())
	fmt.Printf("###### command is `%s`\n", script.command)

	session, err := client.NewSession()
	if err != nil {
		return errors.Wrap(err, "asdf")
	}

	stdin := combine(substitute(script.stdin, map[string]string{
		"PASSWORD": pass,
	}))

	session.Stdin = strings.NewReader(stdin)

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		return errors.Wrap(err, "request pty failed")
	}

	bs, err := session.CombinedOutput(script.command)

	// print the output regardless of err
	output := strings.TrimSpace(string(bs))
	if len(output) == 0 {
		fmt.Println("<no output>")
	} else {
		fmt.Println(output)
	}

	return err
}

func makeClient(user, pass, host string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
	}

	address := fmt.Sprintf("%s:22", host)
	return ssh.Dial("tcp", address, config)
}
