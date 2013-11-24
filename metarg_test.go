package main

import (
	"testing"
)

func TestParseArgsBasic(t *testing.T) {
	const metarCode = "KPKW"
	args := []string{"-d", metarCode}
	flagSet, success := ParseArgs(args)
	t.Logf("Got %v, %v", flagSet, success)
	if !success {
		t.Error("Failed to parse but should've succeeded")
	}
	if !flagSet.Parsed() {
		t.Error("Should be parsed")
	}
	if len(flagSet.Args()) != 1 {
		t.Error("Wrong args length")
	}
	if flagSet.Args()[0] != metarCode {
		t.Error("Wrong METAR code")
	}
	if !decode {
		t.Error("Decode should be true")
	}
	if help {
		t.Error("Help should be false")
	}
}

func TestParseArgsInvalid(t *testing.T) {
	args := []string{"-wrong", "KPKW"}
	_, success := ParseArgs(args)
	if success {
		t.Error("Should've failed")
	}
}

func TestParseArgsMissingMetar(t *testing.T) {
	args := []string{"-d"}
	_, success := ParseArgs(args)
	if success {
		t.Error("Should've failed")
	}
}
