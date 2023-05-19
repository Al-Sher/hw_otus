package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := Environment{}

	t.Run("success command", func(t *testing.T) {
		code := RunCmd([]string{"ls"}, env)

		require.Equal(t, 0, code)
	})

	t.Run("not found command", func(t *testing.T) {
		code := RunCmd([]string{"/bin/"}, env)

		require.Equal(t, 1, code)
	})

	t.Run("command with error", func(t *testing.T) {
		code := RunCmd([]string{"/bin/bash", "-c", "exit 2"}, env)

		require.Equal(t, 2, code)
	})
}
