package coloursetter

import (
	"errors"
	"image/color" //nolint:misspell
	"strings"

	"github.com/nickwells/colour.mod/v2/colour"
	"github.com/nickwells/param.mod/v7/psetter"
)

// RGBPair is used to set a pair of colour value
//
//nolint:misspell
type RGBPair struct {
	psetter.ValueReqMandatory

	Value1   *color.RGBA
	Value2   *color.RGBA
	Families colour.Families
}

// SetWithVal (called with the value following the parameter) either parses
// the RGB value or else looks up the supplied colour name. The search is
// performed "case-blind" - all names are mapped to their lower-case
// equivalents.
func (s RGBPair) SetWithVal(_ string, paramVal string) error {
	colour1, colour2, ok := strings.Cut(paramVal, ";")
	if !ok {
		return errors.New("missing ';' - two colours separated by ; are needed")
	}

	{
		nc, err := colour.ParseNamedColour(s.Families, colour1)
		if err != nil {
			return err
		}

		*s.Value1 = nc.Colour()
	}

	{
		nc, err := colour.ParseNamedColour(s.Families, colour2)
		if err != nil {
			return err
		}

		*s.Value2 = nc.Colour()
	}

	return nil
}

// AllowedValues returns a string describing the allowed values
func (s RGBPair) AllowedValues() string {
	return "a pair of colours separated by ';' where:" +
		colour.NamedColourAllowedValues(s.Families)
}

// ValDescribe returns a string describing the value that can follow the
// parameter
func (s RGBPair) ValDescribe() string {
	return "colour;colour"
}

// CurrentValue returns the current setting of the parameter value
func (s RGBPair) CurrentValue() string {
	return colour.Describe(*s.Value1) + ";" + colour.Describe(*s.Value2)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or the Families value is incorrect. Possible problems with
// the Families member include duplicate Families in the set or an invalid
// Family constant being used.
func (s RGBPair) CheckSetter(name string) {
	intro := name + ": coloursetter.RGB Check failed:"

	if s.Value1 == nil {
		panic(intro + " RGB.Value1: is nil")
	}

	if s.Value2 == nil {
		panic(intro + " RGB.Value2: is nil")
	}

	if err := s.Families.Check(); err != nil {
		panic(intro + " RGB.Families: " + err.Error())
	}
}
