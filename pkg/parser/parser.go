package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"toml-parser/pkg/utils"
)

const (
	KeyValSep          = '='
	WhiteSpace         = " \t"
	CommentStartSymbol = '#'
	StdTableOpen       = '['
	StdTableClose      = ']'
	ArrayTableOpen     = "[["
	ArrayTableClose    = "]]"
	DotSep             = '.'
	QuotationMark      = '"'
	Apostrophe         = '\''
	Escape             = '\\'
	ArrayOpen          = '['
	ArrayClose         = ']'
	InlineTableOpen    = '{'
	InlineTableClose   = '}'
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
			keys, err := parseTableExpression(line)
			if err != nil {
				return nil, err
			}
			// TODO: 抽出後の処理
			fmt.Printf("table: %v\n", keys)
		} else {
			// keyval
			keys, val, err := parseKeyValueExpression(line)
			if err != nil {
				return nil, err
			}
			// TODO: 抽出後の処理
			fmt.Printf("keyval: %v, %v\n", keys, val)
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
func parseKeyValueExpression(line string) ([]string, any, error) {
	sepIdx := strings.Index(line, string(KeyValSep))
	if sepIdx == -1 {
		return nil, nil, fmt.Errorf("missing keyval-sep: %s\n", line)
	}

	keyString := strings.TrimRight(line[:sepIdx], WhiteSpace)
	valueString := strings.TrimLeft(line[sepIdx+1:], WhiteSpace)
	if len(valueString) == 0 {
		return nil, nil, fmt.Errorf("missing val: %s\n", line)
	}

	// TODO: さらにパース
	parsedKey, err := parseKey(keyString)
	if err != nil {
		return nil, nil, err
	}
	parsedValue, err := parseValue(valueString)
	if err != nil {
		return nil, nil, err
	}
	return parsedKey, parsedValue, nil
}

func parseTableExpression(line string) ([]string, error) {
	if utils.IsBracketedBy(line, ArrayTableOpen, ArrayTableClose) {
		/*
			array-table = array-table-open key array-table-close

			array-table-open  = %x5B.5B ws  ; [[ Double left square bracket
			array-table-close = ws %x5D.5D  ; ]] Double right square bracket
		*/
		keyString := strings.Trim(utils.TrimNChar(line, 2), WhiteSpace)
		return parseKey(keyString)
	} else if utils.IsBracketedBy(line, string(StdTableOpen), string(StdTableClose)) {
		/*
			std-table = std-table-open key std-table-close
			std-table-open  = %x5B ws     ; [ Left square bracket
			std-table-close = ws %x5D     ; ] Right square bracket
		*/
		keyString := strings.Trim(utils.TrimChar(line), WhiteSpace)
		return parseKey(keyString)
	}

	// invalid table
	return nil, fmt.Errorf("invalid expression: %s\n", line)
}
