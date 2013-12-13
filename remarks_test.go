package main

import (
	"testing"
)

type RemarkTestCase struct {
	RemarkValue    string
	ExpectedResult string
}

func TestParseRemarks(t *testing.T) {
    testCases := []RemarkTestCase{
        RemarkTestCase{"", ""}, //empty remarks should return empty string
        RemarkTestCase{"AO1", "AMOS station"},
        RemarkTestCase{"AO2", "ASOS station"},
        RemarkTestCase{"SLP123", "Sea level pressure 12.3 mb"},
        RemarkTestCase{"WEA:something", "something"},
        RemarkTestCase{"PRESFR", "Pressure falling rapidly"},
        RemarkTestCase{"PRESRR", "Pressure rising rapidly"},
        RemarkTestCase{"10270", "Max temp in 6 hrs:  27.0 °C"},
        RemarkTestCase{"20221", "Min temp in 6 hrs:  22.1 °C"},
        RemarkTestCase{"21221", "Min temp in 6 hrs:  -22.1 °C"},
        RemarkTestCase{"4/012", "Snow coverage:  12.0\""},
        RemarkTestCase{"51021", "Pressure tendency:  -2.1 mb"},
        RemarkTestCase{"60100", "6-hour precipitation:  1.0\""},
        RemarkTestCase{"70510", "24-hour precipitation:  5.1\""},
        RemarkTestCase{"8/l", "Clouds:  Low"},
        RemarkTestCase{"8/m", "Clouds:  Medium"},
        RemarkTestCase{"8/h", "Clouds:  High"},
        RemarkTestCase{"933012", "New snow coverage (water eq.):  12\""},
    }
    for _, testCase := range testCases {
        result := parseRemark(testCase.RemarkValue)
        if result != testCase.ExpectedResult {
            t.Errorf("Invalid remark.  Expected %v, got %v", testCase.ExpectedResult, result)
        }
    }

}

