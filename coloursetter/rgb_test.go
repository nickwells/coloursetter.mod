package coloursetter

import (
	"image/color"
	"testing"

	"github.com/nickwells/colour.mod/colour"
	"github.com/nickwells/colour.mod/colourtesthelper"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestRGBCheck(t *testing.T) {
	c := color.RGBA{}

	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		v RGB
	}{
		{
			ID: testhelper.MkID("No panic expected, no Families"),
			v:  RGB{Value: &c},
		},
		{
			ID: testhelper.MkID("No panic expected, Families={Any}"),
			v: RGB{
				Value: &c,
				Families: []colour.Family{
					colour.AnyColours,
				},
			},
		},
		{
			ID: testhelper.MkID("No panic expected, Families={X11,Web}"),
			v: RGB{
				Value: &c,
				Families: []colour.Family{
					colour.X11Colours,
					colour.WebColours,
				},
			},
		},
		{
			ID: testhelper.MkID("Panic expected, Families={X11,X11}"),
			ExpPanic: testhelper.MkExpPanic(
				"test-param: coloursetter.RGB Check failed: 1 problem found:",
				"X11Colours appears 2 times, at: Families[0] and Families[1]"),
			v: RGB{
				Value: &c,
				Families: []colour.Family{
					colour.X11Colours,
					colour.X11Colours,
				},
			},
		},
		{
			ID: testhelper.MkID("Panic expected, nil Value"),
			ExpPanic: testhelper.MkExpPanic(
				"test-param: coloursetter.RGB Check failed:" +
					" the Value to be set is nil"),
			v: RGB{},
		},
		{
			ID: testhelper.MkID("Panic expected, Families={Family(99)}"),
			ExpPanic: testhelper.MkExpPanic(
				"test-param: coloursetter.RGB Check failed: 1 problem found:",
				"bad Family: 99 (at Families[0])"),
			v: RGB{
				Value: &c,
				Families: []colour.Family{
					colour.Family(99),
				},
			},
		},
		{
			ID: testhelper.MkID("Panic expected, Families={99,99}"),
			ExpPanic: testhelper.MkExpPanic(
				"test-param: coloursetter.RGB Check failed: 3 problems found:",
				"bad Family: 99 (at Families[0])",
				"bad Family: 99 (at Families[1])",
				"and",
				"BadFamily:99 appears 2 times, at: Families[0] and Families[1]",
			),
			v: RGB{
				Value: &c,
				Families: []colour.Family{
					colour.Family(99),
					colour.Family(99),
				},
			},
		},
		{
			ID: testhelper.MkID("Panic expected, Families={Any,X11}"),
			ExpPanic: testhelper.MkExpPanic(
				"test-param: coloursetter.RGB Check failed: 1 problem found:",
				"AnyColour (at Families[0]) is not the only Family",
			),
			v: RGB{
				Value: &c,
				Families: []colour.Family{
					colour.AnyColours,
					colour.X11Colours,
				},
			},
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(func() {
			tc.v.CheckSetter("test-param")
		})
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}
}

func TestRGBUseAnyColour(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		v   RGB
		exp bool
	}{
		{
			ID:  testhelper.MkID("nil entry: useAnyColours==true"),
			exp: true,
		},
		{
			ID: testhelper.MkID("empty entry: useAnyColours==true"),
			v: RGB{
				Families: []colour.Family{},
			},
			exp: true,
		},
		{
			ID: testhelper.MkID("one entry (AnyColours): useAnyColours==true"),
			v: RGB{
				Families: []colour.Family{colour.AnyColours},
			},
			exp: true,
		},
		{
			ID: testhelper.MkID("one entry (X11Colours): useAnyColours==false"),
			v: RGB{
				Families: []colour.Family{colour.X11Colours},
			},
			exp: false,
		},
	}

	for _, tc := range testCases {
		act := tc.v.useAnyColours()
		testhelper.DiffBool(t, tc.IDStr(), "useAnyColours", act, tc.exp)
	}
}

func TestRGBSetWithVal(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		v        string
		families []colour.Family
		expVal   color.RGBA
	}{
		{
			ID:     testhelper.MkID("RGB{...} - no error expected"),
			v:      "RGB{R: 0xff, G: 0, B: 0}",
			expVal: color.RGBA{R: 0xff, G: 0, B: 0, A: 0xff},
		},
		{
			ID: testhelper.MkID("RGB{Red too big} - error expected"),
			ExpErr: testhelper.MkExpErr(
				`cannot convert the Red value ("0xfff") to a valid number: ` +
					`value out of range`),
			v:      "RGB{R: 0xfff, G: 0, B: 0}",
			expVal: color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
		{
			ID: testhelper.MkID("RGB{Green too big} - error expected"),
			ExpErr: testhelper.MkExpErr(
				`cannot convert the Green value ("0xfff") to a valid number: ` +
					`value out of range`),
			v:      "RGB{R: 0, G: 0xfff, B: 0}",
			expVal: color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
		{
			ID: testhelper.MkID("RGB{Blue too big} - error expected"),
			ExpErr: testhelper.MkExpErr(
				`cannot convert the Blue value ("0xfff") to a valid number: ` +
					`value out of range`),
			v:      "RGB{R: 0, G: 0, B: 0xfff}",
			expVal: color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
		{
			ID: testhelper.MkID("RGB{Red invalid} - error expected"),
			ExpErr: testhelper.MkExpErr(
				`cannot convert the Red value ("xxx") to a valid number: ` +
					`invalid syntax`),
			v:      "RGB{R: xxx, G: 0, B: 0}",
			expVal: color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
		{
			ID: testhelper.MkID("RGB{missing fields} - error expected"),
			ExpErr: testhelper.MkExpErr(
				`cannot get the RGB values from "RGB{R: 0, B: 0}"`),
			v:      "RGB{R: 0, B: 0}",
			expVal: color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
		{
			ID:     testhelper.MkID("BLUE - no error expected"),
			v:      "BLUE",
			expVal: color.RGBA{R: 0, G: 0, B: 0xff, A: 0xff},
		},
		{
			ID:     testhelper.MkID("blue - no error expected"),
			v:      "blue",
			expVal: color.RGBA{R: 0, G: 0, B: 0xff, A: 0xff},
		},
		{
			ID:     testhelper.MkID("green - no error expected"),
			v:      "green",
			expVal: color.RGBA{R: 0, G: 0x80, B: 0, A: 0xff},
		},
		{
			ID:       testhelper.MkID("green(X11) - no error expected"),
			v:        "green",
			families: []colour.Family{colour.X11Colours, colour.WebColours},
			expVal:   color.RGBA{R: 0, G: 0xff, B: 0, A: 0xff},
		},
		{
			ID:       testhelper.MkID("green(Web) - no error expected"),
			v:        "green",
			families: []colour.Family{colour.WebColours, colour.X11Colours},
			expVal:   color.RGBA{R: 0, G: 0x80, B: 0, A: 0xff},
		},
		{
			ID:     testhelper.MkID("no-such-colour - error expected"),
			ExpErr: testhelper.MkExpErr(`bad colour name ("no-such-colour")`),
			v:      "no-such-colour",
			expVal: color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
		{
			ID: testhelper.MkID("law green - error expected"),
			ExpErr: testhelper.MkExpErr(`bad colour name ("law green")` +
				`, did you mean "lawn green", "lawngreen" or "low green"?`),
			v:      "law green",
			expVal: color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
		{
			ID: testhelper.MkID("law green (x11,CGA) - error expected"),
			ExpErr: testhelper.MkExpErr(`bad colour name ("law green")` +
				`, did you mean "lawn green", "lawngreen" or "low green"?`),
			v:        "law green",
			families: []colour.Family{colour.X11Colours, colour.CGAColours},
			expVal:   color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
	}

	for _, tc := range testCases {
		c := color.RGBA{}
		s := RGB{
			Value:    &c,
			Families: tc.families,
		}
		err := s.SetWithVal("", tc.v)
		testhelper.CheckExpErr(t, err, tc)
		colourtesthelper.DiffRGBA(t, tc.IDStr(), "color", c, tc.expVal)
	}
}

func TestRGBCurrentValue(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		v      color.RGBA
		expVal string
	}{
		{
			ID:     testhelper.MkID("black"),
			v:      color.RGBA{R: 0, G: 0, B: 0, A: 0xff},
			expVal: "black",
		},
		{
			ID:     testhelper.MkID("white"),
			v:      color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
			expVal: "white",
		},
	}

	for _, tc := range testCases {
		s := RGB{Value: &tc.v}
		actVal := s.CurrentValue()
		testhelper.DiffString(t, tc.IDStr(), "CurrentValue",
			actVal, tc.expVal)
	}
}

func TestRGBAllowedValue(t *testing.T) {
	AVPrefix := "Either a colour name"
	AVSuffix := " or else a string giving the Red/Green/Blue values" +
		" as follows:" +
		" RGB{R: val, G: val, B: val}" +
		" (the Alpha value is forced to 0xFF)"
	testCases := []struct {
		testhelper.ID
		s      RGB
		expVal string
	}{
		{
			ID:     testhelper.MkID("AnyColours"),
			expVal: AVPrefix + AVSuffix,
		},
		{
			ID: testhelper.MkID("X11"),
			s:  RGB{Families: []colour.Family{colour.X11Colours}},
			expVal: AVPrefix +
				" in the X11 colour-name family" +
				AVSuffix,
		},
		{
			ID: testhelper.MkID("X11,Web"),
			s: RGB{
				Families: []colour.Family{
					colour.X11Colours,
					colour.WebColours,
				},
			},
			expVal: AVPrefix +
				" in one of the X11 or Web colour-name families" +
				AVSuffix,
		},
	}

	for _, tc := range testCases {
		actVal := tc.s.AllowedValues()
		testhelper.DiffString(t, tc.IDStr(), "AllowedValue",
			actVal, tc.expVal)
	}
}
