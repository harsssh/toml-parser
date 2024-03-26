package parser

import (
	"fmt"
	"strings"
)

/*
key = simple-key / dotted-key
*/
func parseKey(key string) ([]string, error) {
	// dot があったら dotted-key
	if strings.Contains(key, string(DotSep)) {
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

func parseUnquotedKey(key string) (string, error) {
	return key, nil
}

/*
dotted-key = simple-key 1*( dot-sep simple-key )
*/
func parseDottedKey(key string) ([]string, error) {
	// "k1..k2" -> ["k1", "", "k2"]
	// "k1 . k2" -> ["k1 ", " k2]
	splitKeys := strings.Split(key, string(DotSep))

	result := make([]string, len(splitKeys))
	// 空白要素がある場合, dot が連続しているためエラー
	// 各要素は trim してから, simple-key としてパース
	for i, k := range splitKeys {
		if len(k) == 0 {
			return nil, fmt.Errorf("invalid dotted-key: %s\n", key)
		}

		parsedKey, err := parseSimpleKey(strings.Trim(k, WhiteSpace))
		if err != nil {
			return nil, err
		}
		result[i] = parsedKey
	}

	return result, nil
}
