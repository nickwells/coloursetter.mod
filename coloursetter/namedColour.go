package coloursetter

import (
	"fmt"

	"github.com/nickwells/colour.mod/v2/colour"
	"github.com/nickwells/param.mod/v7/psetter"
)

// NamedColour is used to set a colour value and to also record the name it
// was given to generate the colour. If the Families value is not set then
// the StandardColours families are used.
//
//nolint:misspell
type NamedColour struct {
	psetter.ValueReqMandatory

	Value    *colour.NamedColour
	Families colour.Families
}

// SetWithVal (called with the value following the parameter) either parses
// the NamedColour value or else looks up the supplied colour name. The
// search is performed "case-blind" - all names are mapped to their
// lower-case equivalents.
func (s NamedColour) SetWithVal(_ string, paramVal string) error {
	nc, err := colour.ParseNamedColour(s.Families, paramVal)
	if err == nil {
		*s.Value = nc
	}

	return err
}

// AllowedValues returns a string describing the allowed values
func (s NamedColour) AllowedValues() string {
	return colour.NamedColourAllowedValues(s.Families)
}

// ValDescribe returns a string describing the value that can follow the
// parameter
func (s NamedColour) ValDescribe() string {
	return "colour"
}

// CurrentValue returns the current setting of the parameter value
func (s NamedColour) CurrentValue() string {
	return s.Value.Name() + fmt.Sprintf("%#4.2v", s.Value.Colour())
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or the Families value is incorrect. Possible problems with
// the Families member include duplicate Families in the set or an invalid
// Family constant being used.
func (s NamedColour) CheckSetter(name string) {
	intro := name + ": coloursetter.NamedColour Check failed:"

	if s.Value == nil {
		panic(intro + " NamedColour.Value: is nil")
	}

	if err := s.Families.Check(); err != nil {
		panic(intro + " NamedColour.Families: " + err.Error())
	}
}
