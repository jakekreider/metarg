package main

import (
	"testing"
)

func TestParseDayTime(t *testing.T) {
	const testDateTime = "210051Z"
	day, time := ParseDayTime(testDateTime)
	t.Logf("Received %s, %v", day, time)
	if day != 21 {
		t.Error("day not correct")
	}
	if time.Hour() != 00 {
		t.Error("Time hour not correct")
	}
	if time.Minute() != 51 {
		t.Error("Time minute not correct")
	}
	t.Log("OK")
}

func TestParseWind(t *testing.T) {
	const testWind = "18055KT"
	direction, wind := ParseWind(testWind)
	t.Logf("Received %v, %v", direction, wind)
	if direction != 180 {
		t.Error("Direction not correct")
	}
	if wind != 55 {
		t.Error("Wind not correct")
	}
	t.Log("OK")
}

func TestParseVisibilityFraction(t *testing.T) {
	const testVisibility = "1/2SM"
	distance := ParseVisibility(testVisibility)
	t.Logf("Received %v ", distance)
	if distance != "1/2" {
		t.Error("Visiblity not correct")
	}
	t.Log("OK")
}

func TestParseCloudsMultiple(t *testing.T) {
	const testClouds = "FEW200 SCT250"
	clouds := ParseClouds(testClouds)
	t.Logf("Received %v ", clouds)
	if len(clouds) != 2 {
		t.Error("Received wrong count of clouds")
	}
	t.Log("OK")
}

func TestParseCloudItem(t *testing.T) {
	const testCloud = "FEW200"
	cloud := ParseCloudDescription(testCloud)
	t.Logf("Received %v ", cloud)
	if cloud != "FEW at 20000" {
		t.Error("Received wrong cloud value")
	}
	t.Log("OK")
}

func TestParseTempDew(t *testing.T) {
	const testTemp = "05/M01"
	temperature, dewPoint := ParseTempDew(testTemp)
	t.Logf("Received %v, %v", temperature, dewPoint)
	if temperature != 5.0 {
		t.Error("Received wrong temperature")
	}

	if dewPoint != -1.0 {
		t.Error("Received wrong dew point")
	}
	t.Log("OK")
}

func TestParsePressure(t *testing.T) {
	const testPressure = "A3006"
	pressure := ParsePressure(testPressure)
	t.Logf("Received %v", pressure)
	if pressure != 30.06 {
		t.Error("Received wrong pressure")
	}
}

func TestParseFullMetar(t *testing.T) {
	const testMetar = "KORD 210051Z 15007KT 10SM OVC060 05/01 A3010 RMK AO2 RAE02 SLP200 P0000 T00500011"
	metar := ParseMetar(testMetar)

	t.Logf("Evaluating %+v ", metar)
	if metar.Station != "KORD" {
		t.Error("Station not correct")
	}
	if metar.Day != 21 {
		t.Error("Day not correct")
	}
	if metar.Visibility != "10" {
		t.Error("Visiblity not correct")
	}
	if metar.WindSpeed != 7 {
		t.Error("Wind speed not correct")
	}
	details := GetDetailMetar(metar)
	if details == "" {
		t.Error("Details failed to parse")
	}
	t.Log("OK")
}


