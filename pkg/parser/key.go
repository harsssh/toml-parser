package parser

import (
	"fmt"
	"strings"
	"toml-parser/pkg/utils"
)

/*
key = simple-key / dotted-key
*/
func parseKey(key string) ([]string, error) {
	// quoted-key でなく, dot があったら dotted-key
	if key[0] != QuotationMark && key[0] != Apostrophe && strings.Contains(key, string(DotSep)) {
		splitKeys, err := parseDottedKey(key)
		if err != nil {
			return nil, err
		}
		return splitKeys, nil
	}

	// simple-key
	parsedKey, err := parseSimpleKey(key)
	if err != nil {
		return nil, err
	}
	return []string{parsedKey}, nil
}

/*
simple-key = quoted-key / unquoted-key
*/
func parseSimpleKey(key string) (string, error) {
	if key[0] == QuotationMark || key[0] == Apostrophe {
		/*
		   quoted-key = basic-string / literal-string
		*/
		return parseNormalString(key)
	}

	return parseUnquotedKey(key)
}

/*
unquoted-key = 1*unquoted-key-char
unquoted-key-char = ALPHA / DIGIT / %x2D / %x5F         ; a-z A-Z 0-9 - _
unquoted-key-char =/ %xB2 / %xB3 / %xB9 / %xBC-BE       ; superscript digits, fractions
unquoted-key-char =/ %xC0-D6 / %xD8-F6 / %xF8-37D       ; non-symbol chars in Latin block
unquoted-key-char =/ %x37F-1FFF                         ; exclude GREEK QUESTION MARK, which is basically a semi-colon
unquoted-key-char =/ %x200C-200D / %x203F-2040          ; from General Punctuation Block, include the two tie symbols and ZWNJ, ZWJ
unquoted-key-char =/ %x2070-218F / %x2460-24FF          ; include super-/subscripts, letterlike/numberlike forms, enclosed alphanumerics
unquoted-key-char =/ %x2C00-2FEF / %x3001-D7FF          ; skip arrows, math, box drawing etc, skip 2FF0-3000 ideographic up/down markers and spaces
unquoted-key-char =/ %xF900-FDCF / %xFDF0-FFFD          ; skip D800-DFFF surrogate block, E000-F8FF Private Use area, FDD0-FDEF intended for process-internal use (unicode)
unquoted-key-char =/ %x10000-EFFFF                      ; all chars outside BMP range, excluding
*/
func parseUnquotedKey(key string) (string, error) {
	// TODO: 英数字, ハイフン, アンダースコア以外の文字種の検証
	pattern := `[a-zA-Z0-9\-_]+`
	if !utils.MatchesPattern(key, pattern) {
		return "", fmt.Errorf("invalid key character: %s\n", key)
	}

	return key, nil
}

/*
dotted-key = simple-key 1*( dot-sep simple-key )
*/
func parseDottedKey(key string) ([]string, error) {
	simpleKeys, err := splitByUnquotedDot(key)
	if err != nil {
		return nil, err
	}

	for i, k := range simpleKeys {
		parsedKey, err := parseSimpleKey(k)
		if err != nil {
			return nil, err
		}
		simpleKeys[i] = parsedKey
	}

	return simpleKeys, nil
}

func splitByUnquotedDot(key string) ([]string, error) {
	var result []string
	inDoubleQuote := false
	inSingleQuote := false
	segmentStart := 0 // 分割するセグメントの開始位置
	for i, c := range key {
		if c == QuotationMark {
			inDoubleQuote = !inDoubleQuote
			continue
		}
		if c == Apostrophe {
			inSingleQuote = !inSingleQuote
			continue
		}

		if c == DotSep && !inDoubleQuote {
			segment := strings.Trim(key[segmentStart:i], WhiteSpace)
			if len(segment) == 0 {
				// dot が連続する場合
				return nil, fmt.Errorf("invalid key: consecutive dots: %s\n", key)
			}
			segmentStart = i + 1
			result = append(result, segment)
		}
	}

	if inSingleQuote || inDoubleQuote {
		return nil, fmt.Errorf("invalid key: unclosed quote: %s\n", key)
	}

	// 最後のセグメントを追加
	segment := strings.Trim(key[segmentStart:], WhiteSpace)
	if len(segment) == 0 {
		return nil, fmt.Errorf("invalid key: cannot end with a dot: %s\n", key)
	}
	result = append(result, segment)

	return result, nil
}
