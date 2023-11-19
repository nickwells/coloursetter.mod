package coloursetter

import (
	"image/color"
	"testing"

	"github.com/nickwells/colour.mod/colourtesthelper"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestAlphaCheck(t *testing.T) {
	c := color.RGBA{}

	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		v Alpha
	}{
		{
			ID: testhelper.MkID("No panic expected"),
			v:  Alpha{Value: &c},
		},
		{
			ID: testhelper.MkID("Panic expected, nil Value"),
			ExpPanic: testhelper.MkExpPanic(
				"test-param: coloursetter.Alpha Check failed:" +
					" the Value to be set is nil"),
			v: Alpha{},
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(func() {
			tc.v.CheckSetter("test-param")
		})
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}
}

func TestAlphaSetWithVal(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		v      string
		expVal color.RGBA
	}{
		{
			ID:     testhelper.MkID("0xff - no error expected"),
			v:      "0xff",
			expVal: color.RGBA{R: 0, G: 0, B: 0, A: 0xff},
		},
		{
			ID:     testhelper.MkID("0xaa - no error expected"),
			v:      "0xaa",
			expVal: color.RGBA{R: 0, G: 0, B: 0, A: 0xaa},
		},
		{
			ID: testhelper.MkID("too big - error expected"),
			ExpErr: testhelper.MkExpErr(
				`cannot convert the alpha value ("0xfff") to a valid number: ` +
					`value out of range`),
			v:      "0xfff",
			expVal: color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
		{
			ID: testhelper.MkID("bad number - error expected"),
			ExpErr: testhelper.MkExpErr(
				`cannot convert the alpha value ("blah") to a valid number: ` +
					`invalid syntax`),
			v:      "blah",
			expVal: color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
	}

	for _, tc := range testCases {
		c := color.RGBA{}
		s := Alpha{
			Value: &c,
		}
		err := s.SetWithVal("", tc.v)
		testhelper.CheckExpErr(t, err, tc)
		colourtesthelper.DiffRGBA(t, tc.IDStr(), "alpha", c, tc.expVal)
	}
}

func TestAlphaCurrentValue(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		v      color.RGBA
		expVal string
	}{
		{
			ID:     testhelper.MkID("zero"),
			v:      color.RGBA{R: 0, G: 0, B: 0, A: 0},
			expVal: "0x00",
		},
		{
			ID:     testhelper.MkID("255"),
			v:      color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
			expVal: "0xff",
		},
	}

	for _, tc := range testCases {
		s := Alpha{Value: &tc.v}
		actVal := s.CurrentValue()
		testhelper.DiffString[string](t, tc.IDStr(), "CurrentValue",
			actVal, tc.expVal)
	}
}

func TestAlphaAllowedValue(t *testing.T) {
	s := Alpha{}
	actVal := s.AllowedValues()
	testhelper.DiffString[string](t, "sole test", "AllowedValue",
		actVal, "some value in the range 0-255")
}
