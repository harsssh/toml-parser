package main

import (
	"fmt"
	"os"
	"toml-parser/pkg/parser"
)

func main() {
	parseResult, err := parser.ParseFile("./example.toml")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error parsing file:", err)
		return
	}
	_, _ = fmt.Fprintln(os.Stderr, parseResult)
}
