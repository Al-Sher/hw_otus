package main

import (
	"bytes"
	"errors"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrNameContainsEqualsSign = errors.New("название файла содержит знак равенства")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entities, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make(Environment)

	for _, entity := range entities {
		if entity.IsDir() {
			continue
		}

		if strings.Contains(entity.Name(), "=") {
			return nil, ErrNameContainsEqualsSign
		}

		content, err := os.ReadFile(dir + "/" + entity.Name())
		if err != nil {
			return nil, err
		}

		if i := bytes.Index(content, []byte("\n")); i != -1 {
			content = content[:i]
		}

		content = bytes.ReplaceAll(content, []byte{0x00}, []byte("\n"))

		str := strings.TrimRight(string(content), " \t\n")

		needRemove := str == ""

		result[entity.Name()] = EnvValue{
			Value:      str,
			NeedRemove: needRemove,
		}
	}

	return result, nil
}
