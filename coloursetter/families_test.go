package coloursetter

import (
	"testing"

	"github.com/nickwells/colour.mod/v2/colour"
	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/paramtest"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

const (
	updFlagNameFamilies     = "upd-gf-Families"
	keepBadFlagNameFamilies = "keep-bad-Families"
)

var commonGFCFamilies = testhelper.GoldenFileCfg{
	DirNames:               []string{"testdata", "Families"},
	Pfx:                    "gf",
	Sfx:                    "txt",
	UpdFlagName:            updFlagNameFamilies,
	KeepBadResultsFlagName: keepBadFlagNameFamilies,
}

func init() {
	commonGFCFamilies.AddUpdateFlag()
	commonGFCFamilies.AddKeepBadResultsFlag()
}

func TestFamiliesSetter(t *testing.T) {
	const dfltParamName = "param-name"

	dfltV1 := colour.Families{colour.CrayolaColours}
	badVal := colour.Families{
		colour.CGAColours,
		colour.CrayolaColors,
		colour.CGAColours,
		colour.CrayolaColors,
		colour.Family("nonesuch"),
	}

	var v1 colour.Families

	testCases := []paramtest.Setter{
		{
			ID: testhelper.MkID("badSetter.nilValue"),
			ExpPanic: testhelper.MkExpPanic(
				"param-name: coloursetter.Families Check failed:" +
					" the Value to be set is nil"),
			PSetter: Families{},
		},
		{
			ID: testhelper.MkID("badSetter.badInitialValue"),
			ExpPanic: testhelper.MkExpPanic(
				"param-name: coloursetter.Families Check failed: Value:" +
					" 3 problems found:" +
					` "CGA" appears 2 times` +
					" in the Families list, at positions: [0 2]," +
					` "Crayola" appears 2 times` +
					" in the Families list, at positions: [1 3]" +
					" and" +
					` "nonesuch" is not a valid Family`),
			PSetter: Families{
				Value: &badVal,
			},
		},
		{
			ID: testhelper.MkID("goodSetter.badParam"),
			PSetter: Families{
				Value: &v1,
			},
			ParamVal:      "nonesuch",
			SetWithValErr: testhelper.MkExpErr(`bad family name "nonesuch"`),
		},
		{
			ID: testhelper.MkID("goodSetter.badParam.dups"),
			PSetter: Families{
				Value: &v1,
			},
			ParamVal: string(colour.CGAColours) +
				"," +
				string(colour.CGAColours),
			SetWithValErr: testhelper.MkExpErr(
				"1 problem found:" +
					` "CGA" appears 2 times` +
					" in the Families list, at positions: [0 1]"),
		},
		{
			ID: testhelper.MkID("goodSetter.goodParam"),
			PSetter: Families{
				Value: &v1,
			},
			ParamVal: string(colour.WebColours),
		},
		{
			ID: testhelper.MkID("goodSetter.goodParam.alias"),
			PSetter: Families{
				Value: &v1,
			},
			ParamVal: farrowAndBallAlias,
		},
		{
			ID: testhelper.MkID("goodSetter.goodParam.multiple"),
			PSetter: Families{
				Value: &v1,
			},
			ParamVal: farrowAndBallAlias + "," + string(colour.WebColours),
		},
	}

	for _, tc := range testCases {
		f := func(t *testing.T) {
			if tc.ParamName == "" {
				tc.ParamName = dfltParamName
			}

			v1 = dfltV1
			tc.GFC = commonGFCFamilies
			tc.SetVR(param.Mandatory)

			tc.Test(t)
		}

		t.Run(tc.Name, f)
	}
}
