// Author hoenig

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/crypto/ssh"

	"github.com/pkg/errors"
)

const filenameFmt = `[0-9]+-`

var filenameRe = regexp.MustCompile(filenameFmt)

// script is an expect script with an ordered name.
// script names must start with a number and a dash.
type script struct {
	name    string
	content string
}

func (s script) String() string {
	return s.name
}

func loadScripts(dir string) ([]script, error) {
	scripts := []script{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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

		bs, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.Wrap(err, "failed to read script")
		}

		scripts = append(scripts, script{name: info.Name(), content: string(bs)})
		return nil
	})

	return scripts, err
}

func runScripts(user, pass string, hosts []string, scripts []script) error {
	for _, host := range hosts {

		client, err := makeClient(user, pass, host)
		if err != nil {
			return errors.Wrap(err, "failed to dial host")
		}

		for _, script := range scripts {
			if err := runScript(client, user, pass, host, script); err != nil {
				return errors.Wrapf(err, "failed to run %s on %s", script, host)
			}
			fmt.Println("")
		}
	}
	return nil
}

const scriptFmt = "#!/bin/bash\n%s"

func runScript(client *ssh.Client, user, pass, host string, script script) error {

	fmt.Printf("###### running script %s on %s\n", script, host)

	// 1. put the script at ~/.commando with permissions 300
	text := fmt.Sprintf(scriptFmt, script.content)
	mkCmd := fmt.Sprintf(`echo "%s" > ~/.commando; chmod 700 ~/.commando`, text)
	if err := run(client, host, mkCmd, false, false); err != nil {
		return errors.Wrap(err, "failed to create .commando file")
	}

	// 3. execute ~/.command
	excCmd := `exec ~/.commando`
	if err := run(client, host, excCmd, true, true); err != nil {
		return errors.Wrap(err, "failed running .commando file")
	}

	// 4. delete ~/.commando

	return nil
}
