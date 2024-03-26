package parser

import (
	"fmt"
	"math"
	"strconv"
	"toml-parser/pkg/utils"
)

/*
float = float-int-part ( exp / frac [ exp ] )
float =/ special-float
*/
// 変換自体は ParseFloat でできる: https://go.dev/ref/spec#Floating-point_literals
func parseFloat(value string) (float64, error) {
	/*
	   special-float = [ minus / plus ] ( inf / nan )
	*/
	// +nan, -nan はエラーになるので別途処理
	if value == "nan" || value == "+nan" || value == "-nan" {
		return math.NaN(), nil
	}
	if value == "inf" || value == "+inf" {
		return math.Inf(0), nil
	}
	if value == "-inf" {
		return math.Inf(-1), nil
	}

	if !isValidFloat(value) {
		return 0, fmt.Errorf("invalid float: %s\n", value)
	}

	parsedFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse float: %s\n", value)
	}
	return parsedFloat, nil
}

/*
float = float-int-part ( exp / frac [ exp ] )

exp = "e" float-exp-part
frac = decimal-point zero-prefixable-int

float-exp-part = [ minus / plus ] zero-prefixable-int
*/
func isValidFloat(value string) bool {
	intPartPattern := generateDecimalIntPattern()
	zeroPrefixableIntPattern := `\d(\d|_\d)*`
	fracPattern := string(DecimalPoint) + zeroPrefixableIntPattern
	expPattern := `e(\+|\-)?` + zeroPrefixableIntPattern

	floatPattern := fmt.Sprintf("%s(%s|%s(%s)?)", intPartPattern, expPattern, fracPattern, expPattern)
	return utils.MatchesPattern(value, floatPattern)
}
