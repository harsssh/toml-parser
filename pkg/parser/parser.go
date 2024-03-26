package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	KeyValSep          = '='
	WhiteSpace         = " \t"
	CommentStartSymbol = '#'
	StdTableOpen       = '['
	StdTableClose      = ']'
	ArrayTableOpen     = "[["
	ArrayTableClose    = "]]"
)

func ParseFile(fileName string) (map[string]any, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]any)
	scanner := bufio.NewScanner(file)
	/* toml = expression *( newline expression ) */
	for scanner.Scan() {
		// TODO: 複数行の array, string
		line := scanner.Text()

		/*
			expression = ws [ comment ]
			expression =/ ws keyval ws [ comment ]
			expression =/ ws table ws [ comment ]
		*/
		line = strings.Trim(removeComment(line), WhiteSpace)
		if len(line) == 0 {
			continue
		}

		// keyval or table
		if line[0] == StdTableOpen {
			/* table = std-table / array-table */
			// [key] or [[key]]
			key, err := parseTableExpression(line)
			if err != nil {
				return nil, err
			}
			// TODO: 抽出後の処理
			fmt.Printf("table: %v\n", key)
		} else {
			// keyval
			key, val, err := parseKeyValueExpression(line)
			if err != nil {
				return nil, err
			}
			// TODO: 抽出後の処理
			fmt.Printf("keyval: (%v, %v)\n", key, val)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func removeComment(line string) string {
	commentIdx := strings.LastIndex(line, string(CommentStartSymbol))
	if commentIdx != -1 {
		// comment exists
		line = line[:commentIdx]
	}
	return line
}

/*
keyval = key keyval-sep val
keyval-sep = ws %x3D ws ; =
*/
func parseKeyValueExpression(line string) (string, any, error) {
	sepIdx := strings.Index(line, string(KeyValSep))
	if sepIdx == -1 {
		return "", nil, fmt.Errorf("missing keyval-sep: %s\n", line)
	}

	keyString := strings.TrimRight(line[:sepIdx], WhiteSpace)
	valueString := strings.TrimLeft(line[sepIdx+1:], WhiteSpace)
	if len(valueString) == 0 {
		return "", nil, fmt.Errorf("missing val: %s\n", line)
	}

	// TODO: さらにパース
	return keyString, valueString, nil
}

// TODO: 何を返すか考える, key?
func parseTableExpression(line string) (string, error) {
	if strings.HasPrefix(line, ArrayTableOpen) {
		if !strings.HasSuffix(line, ArrayTableClose) {
			return "", fmt.Errorf("missing array-table-close: %s\n", line)
		}

		/*
			std-table = std-table-open key std-table-close
			std-table-open  = %x5B ws     ; [ Left square bracket
			std-table-close = ws %x5D     ; ] Right square bracket
		*/
		keyString := strings.Trim(line[2:len(line)-2], WhiteSpace)
		// TODO: さらにパース
		return keyString, nil
	} else if line[0] == StdTableOpen {
		if line[len(line)-1] != StdTableClose {
			return "", fmt.Errorf("missing std-table-close: %s\n", line)
		}

		/*
			array-table = array-table-open key array-table-close

			array-table-open  = %x5B.5B ws  ; [[ Double left square bracket
			array-table-close = ws %x5D.5D  ; ]] Double right square bracket
		*/
		keyString := strings.Trim(line[1:len(line)-1], WhiteSpace)
		// TODO: さらにパース
		return keyString, nil
	}

	// invalid table
	return "", fmt.Errorf("invalid expression: %s\n", line)
}
