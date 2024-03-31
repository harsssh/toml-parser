package main

import (
	"fmt"
	"github.com/k0kubun/pp/v3"
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

	toml, err := parser.ParseFile(fileName)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error parsing file:", err)
		return
	}
	pp.Print(toml)
}
