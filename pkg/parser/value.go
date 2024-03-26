package parser

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

	// TODO: float

	// integer
	parsedValue, err := parseInteger(value)
	if err != nil {
		return nil, err
	}
	return parsedValue, nil
}
