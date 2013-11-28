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
    }
    for _, testCase := range testCases {
        result := parseRemark(testCase.RemarkValue)
        if result != testCase.ExpectedResult {
            t.Errorf("Invalid remark.  Expected %v, got %v", testCase.ExpectedResult, result)
        }
    }

}

