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
	const testWind = "15007KT"
	direction, wind := ParseWind(testWind)
	t.Logf("Received %v, %v", direction, wind)
	if direction != 150 {
		t.Error("Direction not correct")
	}
	if wind != 7.0 {
		t.Error("Wind not correct")
	}
	t.Log("OK")
}

func TestParseFullMetar(t *testing.T) {
	const testMetar = "KORD 210051Z 15007KT 10SM OVC060 05/01 A3010 RMK AO2 RAE02 SLP200 P0000 T00500011"
	metar := ParseMetar(testMetar)

	t.Logf("Evaluating %+v ", metar)
	if metar.station != "KORD" {
		t.Error("Station not correct")
	}
	if metar.day != 21 {
		t.Error("Day not correct")
	}
	t.Log("OK")

}
