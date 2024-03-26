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
	// TODO: inline-table
	// TODO: date-time
	// TODO: float
	// TODO: integer

	return value, nil
}
