package parser

import (
	"strings"
	"toml-parser/pkg/utils"
)

/*
val = string / boolean / array / inline-table / date-time / float / integer
*/
func parseValue(value string) (any, error) {
	// boolean
	if value == "true" {
		return true, nil
	}
	if value == "false" {
		return false, nil
	}

	// string
	// TODO: multiline string 対応
	if value[0] == QuotationMark || value[0] == Apostrophe {
		return parseNormalString(value)
	}

	// TODO: array
	if value[0] == ArrayOpen {
		return value, nil
	}

	// TODO: inline-table
	if value[0] == InlineTableOpen {
		return value, nil
	}

	// TODO: date-time
	// full-date (2024-04-02) or partial-time (09:42) から始まる
	// float, integer 以外はパース済みなので, ":" が含まれたら date-time
	if utils.MatchesPattern(value, `\d{4}-\d{2}-\d{2}`) || strings.Contains(value, ":") {
		return value, nil
	}

	// float
	// integer 以外はパース済みなので, "e" または "." の有無を調べればよい
	// inf, nan の場合もある
	if strings.Contains(value, "e") || strings.Contains(value, ".") ||
		strings.HasSuffix(value, "inf") || strings.HasSuffix(value, "nan") {
		return parseFloat(value)
	}

	// integer
	parsedValue, err := parseInteger(value)
	if err != nil {
		return nil, err
	}
	return parsedValue, nil
}
