// Author hoenig

package main

import (
	"testing"

	"indeed/gophers/3rdparty/p/github.com/stretchr/testify/require"
)

const script1 = `
sudo whoami
PASSWORD
`

const script2 = `

herp derp
PASSWORD
foo
PASSWORD

`

func Test_parseScript(t *testing.T) {
	tests := []struct {
		content  string
		name     string
		expCmd   string
		expStdin []string
	}{
		{
			content:  script1,
			name:     "0-script1",
			expCmd:   "sudo whoami",
			expStdin: []string{"PASSWORD"},
		},
		{
			content:  script2,
			name:     "1-script2",
			expCmd:   "herp depr",
			expStdin: []string{"PASSWORD", "foo", "PASSWORD"},
		},
	}

	check := func(s script, expCmd string, expStdin []string) {
		require.Equal(t, expCmd, s.command)
	}

	for _, test := range tests {
		script, err := parseScript(test.name, test.content)
		require.NoError(t, err)
		check(script, test.expCmd, test.expStdin)
	}
}
