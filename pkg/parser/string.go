package parser

import (
	"fmt"
	"toml-parser/pkg/utils"
)

// basic-string / literal-string
// argument is quoted
func parseNormalString(str string) (string, error) {
	if str[0] == QuotationMark {
		return parseBasicString(str)
	} else if str[0] == Apostrophe {
		return parseLiteralString(str)
	}

	return "", fmt.Errorf("invalid string: %s\n", str)
}

/*
basic-string = quotation-mark *basic-char quotation-mark
basic-char = basic-unescaped / escaped
*/
func parseBasicString(str string) (string, error) {
	if !utils.IsEnclosedBy(str, string(QuotationMark)) {
		return "", fmt.Errorf("missing quotation-mark: %s\n", str)
	}
	unquoted := utils.TrimChar(str)

	if unquoted[0] == Escape {
		// TODO
		/*
			escaped = escape escape-seq-char
			escape = %x5C                   ; \
			escape-seq-char =  %x22         ; "    quotation mark  U+0022
			escape-seq-char =/ %x5C         ; \    reverse solidus U+005C
			escape-seq-char =/ %x62         ; b    backspace       U+0008
			escape-seq-char =/ %x65         ; e    escape          U+001B
			escape-seq-char =/ %x66         ; f    form feed       U+000C
			escape-seq-char =/ %x6E         ; n    line feed       U+000A
			escape-seq-char =/ %x72         ; r    carriage return U+000D
			escape-seq-char =/ %x74         ; t    tab             U+0009
			escape-seq-char =/ %x78 2HEXDIG ; xHH                  U+00HH
			escape-seq-char =/ %x75 4HEXDIG ; uHHHH                U+HHHH
			escape-seq-char =/ %x55 8HEXDIG ; UHHHHHHHH            U+HHHHHHHH
		*/
	} else {
		/*
		   basic-unescaped = wschar / %x21 / %x23-5B / %x5D-7E / non-ascii
		*/
		pattern := `[ \t\x21\x23-\x5b\x5d-\x7e\x80-\x{d7ff}\x{e000}-\x{10ffff}]*`
		if !utils.MatchesPattern(unquoted, pattern) {
			return "", fmt.Errorf("invalid basic-string: %s\n", str)
		}
	}

	return unquoted, nil
}

/*
literal-string = apostrophe *literal-char apostrophe
literal-char = %x09 / %x20-26 / %x28-7E / non-ascii
non-ascii = %x80-D7FF / %xE000-10FFFF
*/
func parseLiteralString(str string) (string, error) {
	if !utils.IsEnclosedBy(str, string(Apostrophe)) {
		return "", fmt.Errorf("missing apostrophe: %s\n", str)
	}
	unquoted := utils.TrimChar(str)

	pattern := `[\x09\x20-\x26\x28-\x7e\x80-\x{d7ff}\x{e000}-\x{10ffff}]*`
	if !utils.MatchesPattern(unquoted, pattern) {
		return "", fmt.Errorf("invalid literal-string: %s\n", str)
	}

	return unquoted, nil
}

// TODO: 実装
// ml-basic-string / ml-literal-string
func parseMultilineString(str string) (string, error) {
	return str, nil
}
