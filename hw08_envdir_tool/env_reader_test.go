package main

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("testdata", func(t *testing.T) {
		envs, err := ReadDir("./testdata/env")

		require.NoError(t, err)

		expected := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		for name, env := range envs {
			if v, ok := expected[name]; ok {
				if v.Value != env.Value || v.NeedRemove != env.NeedRemove {
					t.Errorf("Ожидалось %v, но получено %v", v, env)
				}
			} else if !ok {
				t.Errorf("Значение для ключа %s не найдено", name)
			}
		}
	})

	t.Run("no found dir", func(t *testing.T) {
		_, err := ReadDir("./testNotFoundDir")
		require.Truef(t, errors.Is(err, fs.ErrNotExist), "actual error %q", err)
	})

	t.Run("dir in dir", func(t *testing.T) {
		_, err := ReadDir("./")
		require.NoError(t, err)
	})
}
