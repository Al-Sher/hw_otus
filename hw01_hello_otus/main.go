package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

// stringForReverse строка, которую необходимо перевернуть.
const stringForReverse = "Hello, OTUS!"

func main() {
	fmt.Print(stringutil.Reverse(stringForReverse))
}
