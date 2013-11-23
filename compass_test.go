package main

import (
	"testing"
)

func TestGetCompassPoints(t *testing.T) {
	var values = map[float32]string{
		0:   "N",
		90:  "E",
		180: "S",
		270: "W",
		225: "SW",
		359: "N",
	}
	for value, expected := range values {
		result := GetCompassAbbreviation(value)
		if result != expected {
			t.Errorf("Incorrect value for %v, expected %s but got %s", value, expected, result)
		}
	}

}
