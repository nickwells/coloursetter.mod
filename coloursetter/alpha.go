package coloursetter

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"

	"github.com/nickwells/param.mod/v6/psetter"
)

// Alpha is used to set an RGBA colour's alpha value
type Alpha struct {
	psetter.ValueReqMandatory

	Value *color.RGBA
}

// SetWithVal (called when a value follows the parameter) sets the Value's
// alpha to the result of converting the passed string to a uint8.
func (s Alpha) SetWithVal(_ string, paramVal string) error {
	v64, err := strconv.ParseUint(paramVal, 0, 8)
	if err != nil {
		errIntro := fmt.Sprintf(
			"cannot convert the alpha value (%q) to a valid number", paramVal)
		if errors.Is(err, strconv.ErrRange) {
			return fmt.Errorf("%s: %w", errIntro, strconv.ErrRange)
		}
		if errors.Is(err, strconv.ErrSyntax) {
			return fmt.Errorf("%s: %w", errIntro, strconv.ErrSyntax)
		}
		return fmt.Errorf("%s: %w", errIntro, err)
	}

	s.Value.A = uint8(v64)
	return nil
}

// AllowedValues returns a string describing the allowed values
func (s Alpha) AllowedValues() string {
	return "some value in the range 0-255"
}

// CurrentValue returns the current setting of the parameter value
func (s Alpha) CurrentValue() string {
	return fmt.Sprintf("0x%02x", s.Value.A)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil, if the base unit is invalid or if one of the check functions
// is nil.
func (s Alpha) CheckSetter(name string) {
	intro := name + ": coloursetter.Alpha Check failed:"

	if s.Value == nil {
		panic(intro + " the Value to be set is nil")
	}
}
