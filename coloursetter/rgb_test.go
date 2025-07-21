package coloursetter

import (
	"image/color" //nolint:misspell
	"math"
	"testing"

	"github.com/nickwells/colour.mod/v2/colour"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramtest"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestRGBUseStandardColour(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		v   RGB
		exp bool
	}{
		{
			ID:  testhelper.MkID("nil entry: useStandardColours==true"),
			exp: true,
		},
		{
			ID: testhelper.MkID("empty entry: useStandardColours==true"),
			v: RGB{
				Families: colour.Families{},
			},
			exp: true,
		},
		{
			ID: testhelper.MkID("one entry (StandardColours): useStandardColours==true"),
			v: RGB{
				Families: colour.Families{colour.StandardColours},
			},
			exp: true,
		},
		{
			ID: testhelper.MkID("one entry (X11Colours): useStandardColours==false"),
			v: RGB{
				Families: colour.Families{colour.X11Colours},
			},
			exp: false,
		},
	}

	for _, tc := range testCases {
		act := tc.v.useStandardColours()
		testhelper.DiffBool(t, tc.IDStr(), "useStandardColours", act, tc.exp)
	}
}

func TestRGBCurrentValue(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		v      color.RGBA //nolint:misspell
		expVal string
	}{
		{
			ID:     testhelper.MkID("black"),
			v:      color.RGBA{R: 0, G: 0, B: 0, A: 0xff}, //nolint:misspell
			expVal: "black",
		},
		{
			ID:     testhelper.MkID("white"),
			v:      color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, //nolint:misspell
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

const (
	updFlagNameRGB     = "upd-gf-RGB"
	keepBadFlagNameRGB = "keep-bad-RGB"
)

var commonGFCRGB = testhelper.GoldenFileCfg{
	DirNames:               []string{"testdata", "RGB"},
	Pfx:                    "gf",
	Sfx:                    "txt",
	UpdFlagName:            updFlagNameRGB,
	KeepBadResultsFlagName: keepBadFlagNameRGB,
}

func init() {
	commonGFCRGB.AddUpdateFlag()
	commonGFCRGB.AddKeepBadResultsFlag()
}

func TestRGBSetter(t *testing.T) {
	const dfltParamName = "param-name"

	dfltVal := color.RGBA{R: math.MaxUint8, A: math.MaxUint8} //nolint:misspell
	val := dfltVal

	testCases := []paramtest.Setter{
		{
			ID: testhelper.MkID("value not set"),
			ExpPanic: testhelper.MkExpPanic(
				"coloursetter.RGB Check failed: RGB.Value: is nil"),
			PSetter: RGB{},
		},
		{
			ID: testhelper.MkID("bad families - nonesuch"),
			ExpPanic: testhelper.MkExpPanic(
				"coloursetter.RGB Check failed: RGB.Families: 1 problem found:",
				`"nonesuch" is not a valid Family (at position 0)`,
			),
			PSetter: RGB{
				Value: &val,
				Families: colour.Families{
					colour.Family("nonesuch"),
				},
			},
		},
		{
			ID: testhelper.MkID("bad families - duplicates"),
			ExpPanic: testhelper.MkExpPanic(
				"param-name: coloursetter.RGB Check failed: RGB.Families:" +
					" 1 problem found:" +
					` "CGA" appears 2 times in the Families list,` +
					" at positions: [1 3]"),
			PSetter: RGB{
				Value: &val,
				Families: colour.Families{
					colour.X11Colours,
					colour.CGAColours,
					colour.WebColours,
					colour.CGAColours,
				},
			},
		},
		{
			ID: testhelper.MkID("bad families and duplicates"),
			ExpPanic: testhelper.MkExpPanic(
				"param-name: coloursetter.RGB Check failed: RGB.Families:" +
					" 2 problems found:" +
					` "CGA" appears 2 times in the Families list,` +
					" at positions: [1 3]" +
					" and" +
					` "nonesuch" is not a valid Family (at position 0)`),
			PSetter: RGB{
				Value: &val,
				Families: colour.Families{
					colour.Family("nonesuch"),
					colour.CGAColours,
					colour.WebColours,
					colour.CGAColours,
				},
			},
		},
		{
			ID: testhelper.MkID("goodSetter.badval"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "blac",
			SetWithValErr: testhelper.MkExpErr(`bad colour name ("blac"),`,
				` did you mean "black" or "black ink"?`),
		},
		{
			ID: testhelper.MkID("goodSetter.badval.nonStd-CGA"),
			PSetter: RGB{
				Value: &val,
				Families: colour.Families{
					colour.CGAColours,
				},
			},
			ParamVal: "blac",
			SetWithValErr: testhelper.MkExpErr(`bad colour name ("blac"),`,
				` did you mean "black"?`),
		},
		{
			ID: testhelper.MkID("goodSetter.badval.nonStd-CGAAndWeb"),
			PSetter: RGB{
				Value: &val,
				Families: colour.Families{
					colour.CGAColours,
					colour.WebColours,
				},
			},
			ParamVal: "blac",
			SetWithValErr: testhelper.MkExpErr(`bad colour name ("blac"),`,
				` did you mean "black"?`),
		},
		{
			ID: testhelper.MkID("goodSetter.goodval.std.colourName"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "black",
		},
		{
			ID: testhelper.MkID("goodSetter.goodval.nonStd-CGAAndWeb"),
			PSetter: RGB{
				Value: &val,
				Families: colour.Families{
					colour.CGAColours,
					colour.WebColours,
				},
			},
			ParamVal: "black",
		},
		{
			ID: testhelper.MkID("goodSetter.goodval.familyColour"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "pantone:black olive",
		},
		{
			ID: testhelper.MkID("goodSetter.goodval.nonStdFamilyColour"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "crayola:banana mania",
		},
		{
			ID: testhelper.MkID("goodSetter.goodval.RGB.R"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "RGB{R: 0xf}",
		},
		{
			ID: testhelper.MkID("goodSetter.goodval.RGB.R.withSpace"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: " RGB { R : 0xf } ",
		},
		{
			ID: testhelper.MkID("goodSetter.goodval.rgb.R"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "rgb{R: 0xf}",
		},
		{
			ID: testhelper.MkID("goodSetter.goodval.rgb.R.withSpace"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: " rgb { R : 0xf } ",
		},
		{
			ID: testhelper.MkID("goodSetter.goodval.RGB.G"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "RGB{G: 0xf}",
		},
		{
			ID: testhelper.MkID("goodSetter.goodval.RGB.B"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "RGB{B: 0xf}",
		},
		{
			ID: testhelper.MkID("goodSetter.goodval.RGB.A"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "RGB{A: 0xf}",
		},
		{
			ID: testhelper.MkID("goodSetter.goodval.RGB.ARBG"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "RGB{a: 0xf, r: 0xf, b: 0xf, g: 0xf}",
		},
		{
			ID: testhelper.MkID("goodSetter.goodval.RGBA.ARBG"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "RGBA{A: 0xf, R: 0xf, B: 0xf, G: 0xf}",
		},
		{
			ID: testhelper.MkID("goodSetter.badval.RGB.no-colon"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "RGBA{A 0xf, R: 0xf, B: 0xf, G: 0xf}",
			SetWithValErr: testhelper.MkExpErr(
				`bad colour component: "A 0xf",`,
				`the name and value should be separated by a colon(:)`),
		},
		{
			ID: testhelper.MkID("goodSetter.badval.RGB.unknown-component"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "RGBA{X: 0xf, R: 0xf, B: 0xf, G: 0xf}",
			SetWithValErr: testhelper.MkExpErr(
				`unknown colour component: "X", allowed values: A, B, G or R`),
		},
		{
			ID: testhelper.MkID("goodSetter.badval.RGB.too-big"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "RGBA{A: 0x100, R: 0xf, B: 0xf, G: 0xf}",
			SetWithValErr: testhelper.MkExpErr(
				`cannot convert the A value ("0x100") to a valid number:`,
				` value out of range`),
		},
		{
			ID: testhelper.MkID("goodSetter.badval.RGB.negative"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "RGBA{A: -1, R: 0xf, B: 0xf, G: 0xf}",
			SetWithValErr: testhelper.MkExpErr(
				`cannot convert the A value ("-1") to a valid number:`,
				` invalid syntax`),
		},
		{
			ID: testhelper.MkID("goodSetter.badval.RGB.blah"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "RGBA{A: blah, R: 0xf, B: 0xf, G: 0xf}",
			SetWithValErr: testhelper.MkExpErr(
				`cannot convert the A value ("blah") to a valid number:`,
				` invalid syntax`),
		},
		{
			ID: testhelper.MkID("goodSetter.badval.RGB.missing.brace"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "RGBA{A: 0xf, R: 0xf, B: 0xf, G: 0xf",
			SetWithValErr: testhelper.MkExpErr(
				`the parameter value starts with "RGBA{"`,
				` but has no trailing '}'`),
		},
		{
			ID: testhelper.MkID("goodSetter.badval.famAndCol.badFam"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "pontone:black",
			SetWithValErr: testhelper.MkExpErr(
				`bad colour family name: "pontone" `,
				`did you mean "pantone"?`),
		},
		{
			ID: testhelper.MkID("goodSetter.badval.famAndCol.badCol"),
			PSetter: RGB{
				Value: &val,
			},
			ParamVal: "web:blac",
			SetWithValErr: testhelper.MkExpErr(
				`bad colour name: "blac" `,
				`did you mean "black"?`),
		},
	}

	for _, tc := range testCases {
		f := func(t *testing.T) {
			if tc.ParamName == "" {
				tc.ParamName = dfltParamName
			}

			tc.SetVR(param.Mandatory)
			tc.GFC = commonGFCRGB
			val = dfltVal // reset the value to its default value

			tc.Test(t)
		}
		t.Run(tc.IDStr(), f)
	}
}
