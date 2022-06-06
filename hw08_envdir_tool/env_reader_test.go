package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	test := []struct {
		name     string
		dir      string
		expected Environment
		wantErr  bool
		err      string
	}{
		{
			name: "testdata",
			dir:  "./testdata/env",
			expected: map[string]EnvValue{
				"BAR":   {Value: "bar"},
				"EMPTY": {Value: ""},
				"FOO":   {Value: "   foo\nwith new line"},
				"HELLO": {Value: `"hello"`},
				"UNSET": {Value: "", NeedRemove: true},
			},
		},
		{
			name:    "non-existent directory",
			dir:     "./testdata/non-existent",
			wantErr: true,
			err:     "open ./testdata/non-existent: no such file or directory",
		},
	}

	for _, tc := range test {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			env, err := ReadDir(tc.dir)
			if tc.wantErr {
				require.ErrorContains(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, env)
			}
		})
	}
}
