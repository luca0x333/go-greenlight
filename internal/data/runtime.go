package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrInvalidRuntimeFormat is a custom error returned by UnmarshalJSON() method.
// It is returned when the method is unable to parse or convert the JSON string successfully.
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

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

// UnmarshalJSON method attacched to the Runtime custom type so that it satisfies the
// json.Unmarshaler interface.
// This method must modify the receiver so we must use a pointer receiver. Otherwise we will only
// be modifying a copy which is discarded when this method returns.
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// We expect the incoming JSON value will be a string in the format "runtime mins" and the first
	// thing we need to do is remove the double quotes from this string.
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Split the string to get the number.
	parts := strings.Split(unquotedJSONValue, " ")

	// Sanity check the parts of the string to make sure it's in the correct format.
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	// Parse the string containing the number into a int32
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Convert the int32 to a Runtime type and assign this to the receiver.
	// Note that we use use the * operator to deference the receiver (which is a pointer to a Runtime
	// type) in order to set the underlying value of the pointer.
	*r = Runtime(i)

	return nil
}
