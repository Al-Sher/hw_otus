package main

import (
	"fmt"
	"os"
)

func main() {
	c := len(os.Args)
	if c < 3 {
		// Вычитаем 1, так как первый аргумент - сама программа
		fmt.Printf("количество переданных аргументов должно быть более 2. Получено %d\n", c-1)
		os.Exit(1)
	}

	path := os.Args[1]
	args := os.Args[2:]

	envs, err := ReadDir(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	errCode := RunCmd(args, envs)

	os.Exit(errCode)
}
