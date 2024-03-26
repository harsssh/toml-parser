package main

import (
	"fmt"
	"os"
	"toml-parser/pkg/parser"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "invalid usage")
		return
	}
	fileName := args[1]

	_, err := parser.ParseFile(fileName)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error parsing file:", err)
		return
	}
}
