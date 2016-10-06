// Author hoenig

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

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
