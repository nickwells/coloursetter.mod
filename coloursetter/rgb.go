package coloursetter

import (
	"errors"
	"fmt"
	"image/color" //nolint:misspell
	"regexp"
	"strconv"
	"strings"

	"github.com/nickwells/colour.mod/colour"
	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/param.mod/v6/psetter"
)

var rgbRE = regexp.MustCompile(`RGB{R: (.*), G: (.*), B: (.*)}`)

// RGB is used to set a colour value
type RGB struct {
	psetter.ValueReqMandatory

	Value    *color.RGBA //nolint:misspell
	Families []colour.Family
}

// useAnyColours returns true if the AnyColours family is to be used. This
// will be the case if there are no families given or the only family given
// is the AnyColours family
func (s RGB) useAnyColours() bool {
	return len(s.Families) == 0 ||
		(len(s.Families) == 1 && s.Families[0] == colour.AnyColours)
}

// suggestAltVal will suggest a possible alternative value for the parameter
// value. It will find those strings in the set of possible values that are
// closest to the given value
func (s RGB) suggestAltVal(val string) string {
	var names []string

	if s.useAnyColours() {
		names = colour.AnyColours.ColourNames()
	} else {
		nameDedup := map[string]bool{}

		for _, f := range s.Families {
			for _, n := range f.ColourNames() {
				nameDedup[n] = true
			}
		}

		for n := range nameDedup {
			names = append(names, n)
		}
	}

	return psetter.SuggestionString(psetter.SuggestedVals(val, names))
}

// getColourVal converts the colour part value string into an appropriate
// value. It returns a non-nil error if the value cannot be converted.
func getColourVal(s, partName string) (uint8, error) {
	rVal, err := strconv.ParseUint(s, 0, 8)
	if err != nil {
		errIntro := fmt.Sprintf(
			"cannot convert the %s value (%q) to a valid number", partName, s)
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

// getRGBVals parses and returns the red, green and blue values from the
// slice. Any errors are returned as they are found.
func getRGBVals(rgb []string) (uint8, uint8, uint8, error) {
	var rVal, gVal, bVal uint8

	rVal, err := getColourVal(rgb[1], "Red")
	if err != nil {
		return 0, 0, 0, err
	}

	gVal, err = getColourVal(rgb[2], "Green")
	if err != nil {
		return 0, 0, 0, err
	}

	bVal, err = getColourVal(rgb[3], "Blue")
	if err != nil {
		return 0, 0, 0, err
	}

	return rVal, gVal, bVal, nil
}

// parseRGBString converts a string like: "RGB{R: val, G: val, B: val}" into
// the corresponding RGB value and sets the Setter's Value appropriately.
func (s RGB) parseRGBString(paramVal string) error {
	rgb := rgbRE.FindStringSubmatch(paramVal)
	if rgb == nil || len(rgb) != 4 {
		return fmt.Errorf("cannot get the RGB values from %q", paramVal)
	}

	rVal, gVal, bVal, err := getRGBVals(rgb)
	if err != nil {
		return err
	}

	s.Value.R = rVal
	s.Value.G = gVal
	s.Value.B = bVal
	s.Value.A = 0xff

	return nil
}

// SetWithVal (called with the value following the parameter) either parses
// the RGB value or else looks up the supplied colour name. The search is
// performed "case-blind" - all names are mapped to their lower-case
// equivalents.
func (s RGB) SetWithVal(_ string, paramVal string) error {
	if strings.HasPrefix(paramVal, "RGB{") &&
		strings.HasSuffix(paramVal, "}") {
		return s.parseRGBString(paramVal)
	}

	var cVal color.RGBA //nolint:misspell

	var err error

	if s.useAnyColours() {
		cVal, err = colour.AnyColours.Colour(strings.ToLower(paramVal))
	} else {
		for _, f := range s.Families {
			cVal, err = f.Colour(strings.ToLower(paramVal))
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		if errors.Is(err, colour.ErrBadColour) {
			return fmt.Errorf("bad colour name (%q)%s",
				paramVal, s.suggestAltVal(strings.ToLower(paramVal)))
		}

		return err
	}

	*s.Value = cVal

	return nil
}

// AllowedValues returns a string describing the allowed values
func (s RGB) AllowedValues() string {
	var fName string

	if !s.useAnyColours() {
		if len(s.Families) == 1 {
			fName = " in the " + s.Families[0].String() + " colour-name family"
		} else {
			var families []string
			for _, f := range s.Families {
				families = append(families, f.String())
			}

			fName = " in one of the " + english.Join(families, ", ", " or ") +
				" colour-name families"
		}
	}

	return "Either a colour name" + fName +
		" or else a string giving the Red/Green/Blue values as follows:" +
		" RGB{R: val, G: val, B: val} (the Alpha value is forced to 0xFF)"
}

// ValDescribe returns a string describing the value that can follow the
// parameter
func (s RGB) ValDescribe() string {
	return "colour-name"
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
		panic(intro + " the Value to be set is nil")
	}

	s.checkFamilies(intro)
}

// checkFamilies performs checks on the Colour.Families value. It
// panics, reporting all the problems found, if any problem is found.
func (s RGB) checkFamilies(intro string) {
	if s.useAnyColours() {
		return
	}

	problems, perFamilyIndices := s.findBadFamilies()

	for f, indices := range perFamilyIndices {
		if len(indices) > 1 {
			problems = append(problems, reportDuplicateFamily(f, indices))
		}
	}

	if len(problems) > 0 {
		panic(fmt.Sprintf("%s %d %s found:\n%s",
			intro,
			len(problems), english.Plural("problem", len(problems)),
			english.Join(problems, "\n", "\nand\n")))
	}
}

// reportDuplicateFamily generates a string describing the duplicate
// occurrence of the Family in the setter's Families list.
func reportDuplicateFamily(f colour.Family, indices []int) string {
	famName := ""
	if f.IsValid() {
		famName = f.Literal()
	} else {
		famName = f.String()
	}

	problem := fmt.Sprintf("%s appears %d times, at: ", famName, len(indices))

	idxStrs := []string{}
	for _, idx := range indices {
		idxStrs = append(idxStrs, fmt.Sprintf("Families[%d]", idx))
	}

	problem += english.Join(idxStrs, ", ", " and ")

	return problem
}

// findBadFamilies records problems with individual Family values and also
// records the indexes where each Family instance appears so we can notify of
// any duplicates.
func (s RGB) findBadFamilies() ([]string, map[colour.Family][]int) {
	indices := map[colour.Family][]int{}

	var problems []string

	for i, f := range s.Families {
		indices[f] = append(indices[f], i)

		if !f.IsValid() {
			problems = append(problems,
				fmt.Sprintf("bad Family: %d (at Families[%d])", f, i))
		}

		if f == colour.AnyColours {
			problems = append(problems,
				fmt.Sprintf(
					"AnyColour (at Families[%d]) is not the only Family", i))
		}
	}

	return problems, indices
}
