package parser

import (
	"bufio"
	"fmt"
	"os"
)

func ParseFile(fileName string) (map[string]any, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]any)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// TODO: 行を解析
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
