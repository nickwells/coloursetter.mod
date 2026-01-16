package coloursetter

import (
	"errors"
	"fmt"
	"image/color" //nolint:misspell
	"maps"
	"math"
	"regexp"
	"slices"
	"strings"

	"github.com/nickwells/colour.mod/v2/colour"
	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/param.mod/v6/psetter"
)

var rgbIntroRE = regexp.MustCompile(
	`^[[:space:]]*[rR][gG][bB]([aA])?[[:space:]]*\{`)
var rgbOutroRE = regexp.MustCompile(`[[:space:]]*\}[[:space:]]*$`)

// RGB is used to set a colour value
//
//nolint:misspell
type RGB struct {
	psetter.ValueReqMandatory

	Value    *color.RGBA
	Families colour.Families
}

// useStandardColours returns true if the StandardColours family is to be
// used. This will be the case if there are no families given or the only
// family given is the StandardColours family
func (s RGB) useStandardColours() bool {
	return len(s.Families) == 0 ||
		(len(s.Families) == 1 && s.Families[0] == colour.StandardColours)
}

// suggestAltVal will suggest a possible alternative value for the parameter
// value. It will find those strings in the set of possible values that are
// closest to the given value
func (s RGB) suggestAltVal(val string) string {
	var names []string

	var err error

	if s.useStandardColours() {
		names, err = colour.StandardColours.ColourNames()
		if err != nil {
			return ""
		}
	} else {
		nameDedup := map[string]bool{}

		for _, f := range s.Families {
			fNames, err := f.ColourNames()
			if err != nil {
				return ""
			}

			for _, n := range fNames {
				nameDedup[n] = true
			}
		}

		for n := range nameDedup {
			names = append(names, n)
		}
	}

	return psetter.SuggestionString(psetter.SuggestedVals(val, names))
}

// parseRGBString converts a string like: "RGB{R: val, G: val, B: val}" into
// the corresponding RGB value and sets the Setter's Value appropriately.
func (s RGB) parseRGBString(paramVal string) error {
	strippedVal := rgbIntroRE.ReplaceAllString(paramVal, "")
	strippedVal = rgbOutroRE.ReplaceAllString(strippedVal, "")

	components := map[string]uint8{
		"R": 0,
		"G": 0,
		"B": 0,
		"A": math.MaxUint8,
	}

	for part := range strings.SplitSeq(strippedVal, ",") {
		name, val, ok := strings.Cut(part, ":")
		if !ok {
			return fmt.Errorf(
				"bad colour component: %q,"+
					" the name and value should be separated by a colon(:)",
				part)
		}

		name = strings.TrimSpace(name)
		name = strings.ToUpper(name)

		if _, valid := components[name]; !valid {
			aval := slices.Collect(maps.Keys(components))
			slices.Sort(aval)

			return fmt.Errorf(
				"unknown colour component: %q, allowed values: %s",
				name, english.Join(aval, ", ", " or "))
		}

		val = strings.TrimSpace(val)

		parsedVal, err := parseColourPart(val, name)
		if err != nil {
			return err
		}

		components[name] = parsedVal
	}

	s.Value.R = components["R"]
	s.Value.G = components["G"]
	s.Value.B = components["B"]
	s.Value.A = components["A"]

	return nil
}

// setByFamilyAndColourName sets the RGB from the family and colour names. It
// returns a non-nil error if the value can not be set.
func (s RGB) setByFamilyAndColourName(fName, cName string) error {
	f := colour.Family(fName)
	if !f.IsValid() {
		return fmt.Errorf("bad colour family name: %q %s",
			fName,
			psetter.SuggestionString(
				psetter.SuggestedVals(fName,
					slices.Collect(maps.Keys(colour.AllowedFamilies())),
				)))
	}

	cVal, err := f.Colour(strings.ToLower(cName))
	if err != nil {
		altNames := ""
		if cNames, err := f.ColourNames(); err == nil {
			altNames = psetter.SuggestionString(
				psetter.SuggestedVals(cName, cNames))
		}

		return fmt.Errorf("bad colour name: %q, %s %s", cName, err, altNames)
	}

	*s.Value = cVal

	return nil
}

// setByColourName sets the RGB from the colour name. It
// returns a non-nil error if the value can not be set.
func (s RGB) setByColourName(cName string) error {
	var cVal color.RGBA //nolint:misspell

	var err error

	if s.useStandardColours() {
		cVal, err = colour.StandardColours.Colour(strings.ToLower(cName))
	} else {
		for _, f := range s.Families {
			cVal, err = f.Colour(strings.ToLower(cName))
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		if errors.Is(err, errors.New(colour.BadColourName)) {
			return fmt.Errorf("bad colour name (%q)%s",
				cName, s.suggestAltVal(strings.ToLower(cName)))
		}

		return err
	}

	*s.Value = cVal

	return nil
}

// SetWithVal (called with the value following the parameter) either parses
// the RGB value or else looks up the supplied colour name. The search is
// performed "case-blind" - all names are mapped to their lower-case
// equivalents.
func (s RGB) SetWithVal(_ string, paramVal string) error {
	if rgbIntroRE.MatchString(paramVal) {
		if !rgbOutroRE.MatchString(paramVal) {
			return fmt.Errorf(
				"the parameter value starts with %q but has no trailing '}'",
				rgbIntroRE.FindString(paramVal))
		}

		return s.parseRGBString(paramVal)
	}

	if familyName, colourName, found := strings.Cut(paramVal, ":"); found {
		return s.setByFamilyAndColourName(familyName, colourName)
	}

	return s.setByColourName(paramVal)
}

// AllowedValues returns a string describing the allowed values
func (s RGB) AllowedValues() string {
	fName := " in the standard colour-name families"

	if !s.useStandardColours() {
		if len(s.Families) == 1 {
			fName = " in the " + string(s.Families[0]) + " colour-name family"
		} else {
			fName = " in one of the colour-name families: " +
				s.Families.String()
		}
	}

	return "Either a colour name" + fName +
		" or a family name, a colon (:) and a colour name" +
		" or else a string giving the Red/Green/Blue/Alpha values as follows:" +
		" RGB{R: #, G: #, B: #, A: #} (Red, Green and Blue default to 0," +
		" Alpha defaults to 0xFF)"
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
