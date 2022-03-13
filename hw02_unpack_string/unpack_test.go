package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `qw\ne`, `qqq\a`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestIsEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{name: "letter", input: 'a', expected: false},
		{name: "digit", input: '5', expected: false},
		{name: "control", input: '\n', expected: false},
		{name: "escape", input: '\\', expected: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := isEscape(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestIsNotValidForEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{name: "digit", input: '5', expected: false},
		{name: "letter", input: 'a', expected: true},
		{name: "control", input: '\n', expected: false},
		{name: "escape", input: '\\', expected: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := isNotValidForEscaping(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestRuneIsNotRepeatable(t *testing.T) {
	tests := []struct {
		name     string
		input    *complexRune
		expected bool
	}{
		{name: "repeatable", input: &complexRune{repeatable: true}, expected: false},
		{name: "not_repeatable", input: &complexRune{repeatable: false}, expected: true},
		{name: "nil", input: nil, expected: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := runeIsNotRepeatable(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestRuneIsNotEscaped(t *testing.T) {
	tests := []struct {
		name     string
		input    *complexRune
		expected bool
	}{
		{name: "escaped", input: &complexRune{escaped: true}, expected: false},
		{name: "not_escaped", input: &complexRune{escaped: false}, expected: true},
		{name: "nil", input: nil, expected: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := runeIsNotEscaped(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}
