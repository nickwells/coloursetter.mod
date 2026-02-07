package coloursetter

import (
	"fmt"
	"strings"

	"github.com/nickwells/colour.mod/v2/colour"
	"github.com/nickwells/param.mod/v7/psetter"
)

var familyAllowedValues = psetter.AllowedVals[string](colour.AllowedFamilies())

const farrowAndBallAlias = "fnb"

var familyAliases = psetter.Aliases[string]{
	farrowAndBallAlias: []string{colour.FarrowAndBallColours.Name()},
}

// Families is a parameter setter for a
// [github.com/nickwells/colour.mod/v2/colour.Families]
type Families struct {
	psetter.ValueReqMandatory

	// Value must be set, the program will panic if not. This is the
	// colour.Families that this setter is setting.
	Value *colour.Families
	// The StrListSeparator allows you to override the default separator
	// between list elements.
	psetter.StrListSeparator
}

// SetWithVal (called when a value follows the parameter) checks the value
// for validity and only if the value is allowed (if it's one of the allowed
// colour family names or an alias) does it set the paramerer. It returns an
// error if the parameter is invalid.
func (s Families) SetWithVal(_ string, paramVal string) error {
	fl := colour.Families{}
	sep := s.GetSeparator()

	vals := strings.SplitSeq(paramVal, sep)
	for v := range vals {
		v = strings.ToLower(v)
		if !familyAllowedValues.ValueAllowed(v) {
			if !familyAliases.IsAnAlias(v) {
				return fmt.Errorf("bad family name %q", v)
			}

			for _, av := range familyAliases.AliasVal(v) {
				f, err := colour.GetFamily(av)
				if err != nil {
					return err
				}

				fl = append(fl, f)
			}

			continue
		}

		f, err := colour.GetFamily(v)
		if err != nil {
			return err
		}

		fl = append(fl, f)
	}

	if err := fl.Check(); err != nil {
		return err
	}

	*s.Value = fl

	return nil
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or the initial value is not allowed.
func (s Families) CheckSetter(name string) {
	const setterName = "coloursetter.Families"

	// Check the value is not nil
	if s.Value == nil {
		panic(psetter.NilValueMessage(name, setterName))
	}

	intro := name + ": " + setterName + " Check failed:"

	if err := s.Value.Check(); err != nil {
		panic(intro + " Value: " + err.Error())
	}

	err := familyAliases.Check(familyAllowedValues)
	if err != nil {
		panic(intro + " common Aliases: " + err.Error())
	}
}

// CurrentValue returns the current setting of the parameter value
func (s Families) CurrentValue() string {
	return s.Value.String()
}

// ValDescribe returns a string describing the value that can follow the
// parameter.
func (s Families) ValDescribe() string {
	return "colour-families"
}

// AllowedValues returns a string listing the allowed values
func (s Families) AllowedValues() string {
	return s.ListValDesc("colour-family names")
}

// AllowedValuesMap returns the map of allowed values for the colour family
// setter
func (s Families) AllowedValuesMap() psetter.AllowedVals[string] {
	return familyAllowedValues
}

// AllowedValuesAliasMap returns the map of allowed alias values for the
// colour family setter
func (s Families) AllowedValuesAliasMap() psetter.Aliases[string] {
	return familyAliases
}
