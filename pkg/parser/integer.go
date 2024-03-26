package parser

import (
	"fmt"
	"strconv"
	"strings"
	"toml-parser/pkg/utils"
)

/*
integer = dec-int / hex-int / oct-int / bin-int

dec-int = [ minus / plus ] unsigned-dec-int
unsigned-dec-int = DIGIT / digit1-9 1*( DIGIT / underscore DIGIT )

hex-int = hex-prefix HEXDIG *( HEXDIG / underscore HEXDIG )
oct-int = oct-prefix digit0-7 *( digit0-7 / underscore digit0-7 )
bin-int = bin-prefix digit0-1 *( digit0-1 / underscore digit0-1 )
*/
func parseInteger(numStr string) (int64, error) {
	var base int
	// TODO: 無駄にマッチさせない
	if utils.MatchesPattern(numStr, generateIntPattern("0x", "[0-9A-Z]")) {
		base = 16
	} else if utils.MatchesPattern(numStr, generateIntPattern("0o", "[0-7]")) {
		base = 8
	} else if utils.MatchesPattern(numStr, generateIntPattern("0b", "[01]")) {
		base = 2
	} else if utils.MatchesPattern(numStr, generateDecimalIntPattern()) {
		base = 10
	} else {
		return 0, fmt.Errorf("invalid integer: %s\n", numStr)
	}

	if base != 10 {
		numStr = numStr[2:] // remove prefix
	}
	numStrWithoutUnderscore := strings.ReplaceAll(numStr, "_", "")
	parsedInt, err := strconv.ParseInt(numStrWithoutUnderscore, base, 64)
	if err != nil {
		return 0, err
	}

	return parsedInt, nil
}

func generateDecimalIntPattern() string {
	return `(\+|\-)?(\d|(\d(\d|_\d)*))`
}

func generateIntPattern(prefixPattern, digitPattern string) string {
	return fmt.Sprintf(`%s%s(%s|_%s)*`, prefixPattern, digitPattern, digitPattern, digitPattern)
}
