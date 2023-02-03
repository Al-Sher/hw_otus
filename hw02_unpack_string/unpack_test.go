package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
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
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestUnpackInvalidStringWithAsterisk(t *testing.T) {
	invalidStrings := []string{`s2\qe`, `s2q\@`, `\_`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			tc1, err := Unpack(tc)
			fmt.Println(tc1)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestUnpackWithSymbols(t *testing.T) {
	tests := []struct {
		input    string
		excepted string
		hasError bool
	}{
		{"t2_4\n6", "tt____\n\n\n\n\n\n", false},
		{"≈5", "≈≈≈≈≈", false},
		{"❄3", "❄❄❄", false},
		{`_5\\55_53`, "", true},
		{"=^_^=", "=^_^=", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := Unpack(test.input)
			if test.hasError {
				require.Error(t, err)
			} else {
				require.Equal(t, test.excepted, result)
			}
		})
	}
}

func TestNumberError(t *testing.T) {
	tests := []struct {
		prev     rune
		cur      rune
		hasError bool
	}{
		{'0', '1', true},
		{'w', '1', false},
	}

	for _, test := range tests {
		testName := []rune{test.prev, test.cur}
		t.Run(string(testName), func(t *testing.T) {
			err := numberError(test.prev, test.cur)
			if test.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRepeatIfNeed(t *testing.T) {
	tests := []struct {
		prev     rune
		cur      rune
		excepted string
	}{
		{'a', '2', "aa"},
		{'a', 'b', ""},
		{'a', '5', "aaaaa"},
		{'b', '2', "bb"},
	}

	for _, test := range tests {
		var buf strings.Builder
		testName := []rune{test.prev, test.cur}
		t.Run(string(testName), func(t *testing.T) {
			_, err := repeatIfNeed(&buf, test.prev, test.cur)
			if err != nil {
				require.Fail(t, err.Error())
			}
			require.Equal(t, test.excepted, buf.String())
		})
	}
}
