package coloursetter

import (
	"errors"
	"fmt"
	"strconv"
)

// parseColourPart takes the named part of a colour value (the Red, Green, Blue
// or Alpha component) as a string and converts it into an appropriate 8-bit
// value. It returns a non-nil error if the value cannot be converted.
func parseColourPart(val, partName string) (uint8, error) {
	rVal, err := strconv.ParseUint(val, 0, 8)
	if err != nil {
		errIntro := fmt.Sprintf(
			"cannot convert the %s value (%q) to a valid number", partName, val)
		if errors.Is(err, strconv.ErrRange) {
			return 0, fmt.Errorf("%s: %w", errIntro, strconv.ErrRange)
		}

		if errors.Is(err, strconv.ErrSyntax) {
			return 0, fmt.Errorf("%s: %w", errIntro, strconv.ErrSyntax)
		}

		return 0, fmt.Errorf("%s: %w", errIntro, err)
	}

	return uint8(rVal), nil
}
