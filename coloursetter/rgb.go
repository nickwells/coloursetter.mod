package coloursetter

import (
	"image/color" //nolint:misspell

	"github.com/nickwells/colour.mod/v2/colour"
	"github.com/nickwells/param.mod/v7/psetter"
)

// RGB is used to set a colour value
//
//nolint:misspell
type RGB struct {
	psetter.ValueReqMandatory

	Value    *color.RGBA
	Families colour.Families
}

// SetWithVal (called with the value following the parameter) either parses
// the RGB value or else looks up the supplied colour name. The search is
// performed "case-blind" - all names are mapped to their lower-case
// equivalents.
func (s RGB) SetWithVal(_ string, paramVal string) error {
	nc, err := colour.ParseNamedColour(s.Families, paramVal)
	if err == nil {
		*s.Value = nc.Colour()
	}

	return err
}

// AllowedValues returns a string describing the allowed values
func (s RGB) AllowedValues() string {
	return colour.NamedColourAllowedValues(s.Families)
}

// ValDescribe returns a string describing the value that can follow the
// parameter
func (s RGB) ValDescribe() string {
	return "colour"
}

// CurrentValue returns the current setting of the parameter value
func (s RGB) CurrentValue() string {
	return colour.Describe(*s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or the Families value is incorrect. Possible problems with
// the Families member include duplicate Families in the set or an invalid
// Family constant being used.
func (s RGB) CheckSetter(name string) {
	intro := name + ": coloursetter.RGB Check failed:"

	if s.Value == nil {
		panic(intro + " RGB.Value: is nil")
	}

	if err := s.Families.Check(); err != nil {
		panic(intro + " RGB.Families: " + err.Error())
	}
}
