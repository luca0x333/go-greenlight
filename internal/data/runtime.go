package data

import (
	"fmt"
	"strconv"
)

// Runtime type with an underlying type int32.
type Runtime int32

// MarshalJSON method attached to the Runtime custom type so that it satisfies the
// json.Marshaler interface.
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	// strconv.Quote() wraps jsonValue in double quotes.
	// It needs to be wrapped in double quotes to be a valid JSON string.
	quotedJSONValue := strconv.Quote(jsonValue)

	// Convert the quoted string to []byte and return it.
	return []byte(quotedJSONValue), nil
}
