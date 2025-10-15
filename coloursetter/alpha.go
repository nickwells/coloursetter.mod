package coloursetter

import (
	"fmt"
	"image/color" //nolint:misspell

	"github.com/nickwells/param.mod/v6/psetter"
)

// Alpha is used to set an RGBA colour's alpha value
//
//nolint:misspell
type Alpha struct {
	psetter.ValueReqMandatory

	Value *color.RGBA
}

// SetWithVal (called when a value follows the parameter) sets the Value's
// alpha to the result of converting the passed string to a uint8.
func (s Alpha) SetWithVal(_ string, paramVal string) error {
	alpha, err := parseColourPart(paramVal, "alpha")
	if err != nil {
		return err
	}

	s.Value.A = uint8(alpha)

	return nil
}

// AllowedValues returns a string describing the allowed values
func (s Alpha) AllowedValues() string {
	return "some value in the range 0-255"
}

// CurrentValue returns the current setting of the parameter value
func (s Alpha) CurrentValue() string {
	return fmt.Sprintf("%#02x", s.Value.A)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s Alpha) CheckSetter(name string) {
	intro := name + ": coloursetter.Alpha Check failed:"

	if s.Value == nil {
		panic(intro + " the Value to be set is nil")
	}
}
