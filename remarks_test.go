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
        RemarkTestCase{"10270", "Max temp in 6 hrs:  27.0"},
        RemarkTestCase{"20221", "Min temp in 6 hrs:  22.1"},
        RemarkTestCase{"21221", "Min temp in 6 hrs:  -22.1"},
        RemarkTestCase{"4/012", "Snow coverage:  12.0\""},
    }
    for _, testCase := range testCases {
        result := parseRemark(testCase.RemarkValue)
        if result != testCase.ExpectedResult {
            t.Errorf("Invalid remark.  Expected %v, got %v", testCase.ExpectedResult, result)
        }
    }

}

