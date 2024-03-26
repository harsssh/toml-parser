package parser

import (
	"fmt"
	"strings"
	"toml-parser/pkg/utils"
)

/*
array = array-open [ array-values ] ws-comment-newline array-close
array-values =  ws-comment-newline val ws-comment-newline array-sep array-values
array-values =/ ws-comment-newline val ws-comment-newline [ array-sep ]

ws-comment-newline = *( wschar / [ comment ] newline )
*/
func parseArray(str string) ([]any, error) {
	if !utils.IsBracketedBy(str, string(ArrayOpen), string(ArrayClose)) {
		return nil, fmt.Errorf("invalid array: %s\n", str)
	}
	unwrapped := utils.TrimChar(str)

	// [ array-values ] ws-comment-newline
	// TODO: newline がある場合も対応
	// newline なしなら, array = [ array-values ] ws
	arrayValuesStr := strings.TrimRight(unwrapped, WhiteSpace)
	if len(arrayValuesStr) == 0 {
		return []any{}, nil
	}
	return parseArrayValues(arrayValuesStr)
}

func parseArrayValues(str string) ([]any, error) {
	// TODO: newline がある場合も対応
	// newline なしなら,
	// array-values = ws val ws array-sep array-values
	// array-values =/ ws val ws [ array-sep ]
	str = strings.TrimLeft(str, WhiteSpace)

	// val ws array-sep array-values
	// val ws [ array-sep ]
	sepIdx := findArraySep(str)
	if sepIdx == -1 {
		// val ws
		parsedValue, err := parseValue(strings.TrimRight(str, WhiteSpace))
		if err != nil {
			return nil, err
		}
		return []any{parsedValue}, nil
	}

	// val ws array-sep array-values
	// val ws array-sep
	left, right := utils.SplitAtIndex(str, sepIdx)

	parsedValue, err := parseValue(strings.TrimRight(left, WhiteSpace))
	if err != nil {
		return nil, err
	}
	if len(right) == 0 {
		return []any{parsedValue}, nil
	}

	values, err := parseArrayValues(right)
	if err != nil {
		return nil, err
	}
	return append([]any{parsedValue}, values...), nil
}

// TODO: もっとキレイに実装
// val ws array-sep array-values
// val ws [ array-sep ]
// この状態で次の array-sep を探す
// val に "," が含まれ得ることに注意
func findArraySep(arrayValuesStr string) int {
	waitDoubleQuote := false
	waitSingleQuote := false
	arrayCount := 0
	tableCount := 0
	for i, c := range arrayValuesStr {
		if arrayCount == 0 && tableCount == 0 &&
			!waitDoubleQuote && !waitSingleQuote && c == ArraySep {
			return i
		}

		switch c {
		case QuotationMark:
			if waitDoubleQuote {
				waitDoubleQuote = false
			} else {
				waitDoubleQuote = true
			}
		case Apostrophe:
			if waitSingleQuote {
				waitSingleQuote = false
			} else {
				waitSingleQuote = true
			}
		case ArrayOpen:
			arrayCount++
		case ArrayClose:
			arrayCount--
		case InlineTableOpen:
			tableCount++
		case InlineTableClose:
			tableCount--
		}
	}

	return -1
}
