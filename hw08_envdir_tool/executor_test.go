package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	test := []struct {
		name         string
		cmd          []string
		env          Environment
		expectedOut  string
		expectedCode int
	}{
		{
			name: "success",
			cmd:  []string{"/bin/bash", "-c", "echo FOO=$FOO BAR=$BAR"},
			env: map[string]EnvValue{
				"FOO": {Value: "foo"},
				"BAR": {Value: "bar", NeedRemove: true},
			},
			expectedOut:  "FOO=foo BAR=\n",
			expectedCode: 0,
		},
		{
			name:         "fail",
			cmd:          []string{""},
			expectedCode: 1,
		},
	}

	for _, tc := range test {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// http://craigwickesser.com/2015/01/capture-stdout-in-go/
			old := os.Stdout
			r, w, err := os.Pipe()
			require.NoError(t, err)
			os.Stdout = w

			code := RunCmd(tc.cmd, tc.env)

			err = w.Close()
			require.NoError(t, err)
			os.Stdout = old

			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			require.NoError(t, err)

			require.Equal(t, tc.expectedOut, buf.String())
			require.Equal(t, tc.expectedCode, code)
		})
	}
}
