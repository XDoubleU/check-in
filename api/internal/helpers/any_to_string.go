package helpers

import "strconv"

func AnyToString(value any) string {
	var result string

	var strVal string
	var int64Val int64

	var ok bool

	if strVal, ok = value.(string); ok {
		result = strVal
	} else if int64Val, ok = value.(int64); ok {
		result = strconv.FormatInt(int64Val, 10)
	} else {
		panic("undefined type")
	}

	return result
}
