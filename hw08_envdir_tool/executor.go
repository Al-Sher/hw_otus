package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	c := exec.Command(cmd[0], cmd[1:]...)

	var newEnv []string
	var envs = c.Environ()

	for _, value := range envs {
		envKeyWithValue := strings.Split(value, "=")
		key := envKeyWithValue[0]
		if v, ok := env[key]; ok && v.NeedRemove {
			continue
		}
		newEnv = append(newEnv, value)
	}

	for k, value := range env {
		if !value.NeedRemove {
			newEnv = append(newEnv, fmt.Sprintf("%s=%s", k, value.Value))
		}
	}

	c.Env = newEnv
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
		// Произошла ошибка приложения при запуске команды
		return 1
	}

	return 0
}
