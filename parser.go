package main

import (
	"fmt"
	"github.com/mragh/metarg/compass"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Regexp wrapper
type MappableRegexp struct {
	regexp.Regexp
}

type Metar struct {
	Station, Phenomena, Visibility, WindDirection                             string
	Clouds                                                                    []string
	Time                                                                      time.Time
	WindSpeed, WindGust, Temperature, Dewpoint, Pressure, WindDirectionDegree float32
	Day                                                                       int32
}

// Returns a map of named groups to values from the given input string
func (this *MappableRegexp) GetMap(input string) (result map[string]string) {
	result = make(map[string]string)
	fields := this.SubexpNames()
	matches := this.FindStringSubmatch(input)
	for i, value := range fields[1:] {
		result[value] = matches[i+1]
	}
	return
}

func ParseMetar(flatMetar string) (metar Metar, success bool) {
	mappable := MappableRegexp{*(regexp.MustCompile(
			`^(?P<station>\w{4})\s(?P<time>\w{7})\s(AUTO\s)?(?P<wind>\w+)\s(?P<visibility>\w+)` +
			`\s+(?P<clouds>.*)\s(?P<tempdue>M?\d\d\/M?\d\d)\s(?P<pressure>A\d{4})\s.*`))}
	matches := mappable.GetMap(flatMetar)
	if len(matches) < 7 {
		return metar, false
	}
	metar.Station = matches["station"]
	metar.Day, metar.Time = ParseDayTime(matches["time"])

	metar.WindDirection, metar.WindSpeed,
		metar.WindDirectionDegree, metar.WindGust = ParseWind(matches["wind"])
	metar.Visibility = ParseVisibility(matches["visibility"])
	metar.Clouds = ParseClouds(matches["clouds"])
	metar.Temperature, metar.Dewpoint = ParseTempDew(matches["tempdue"])
	metar.Pressure = ParsePressure(matches["pressure"])
	return metar, true
}


func ParseWind(windFlat string) (direction string, speed float32,
	dirDegrees float32, gust float32) {
	mapRegex := MappableRegexp{*regexp.MustCompile(`(\d{3})(\d+)(G(\d+))?KT`)}
	match := mapRegex.FindStringSubmatch(windFlat)
	dirDegrees64, _ := strconv.ParseInt(match[1], 10, 32)
	speed64, _ := strconv.ParseFloat(match[2], 32)
	dirDegrees = float32(dirDegrees64)
	direction = compass.GetCompassAbbreviation(dirDegrees)
	speed = float32(speed64)
	if match[4] != "" {
		gust64, _ := strconv.ParseFloat(match[4], 32)
		gust = float32(gust64)
	} else {
		gust = speed
	}
	return
}

// returns a string (for now at least) since values can be 1/2, etc.
func ParseVisibility(visibilityFlat string) (metarVisibility string) {
	regex := regexp.MustCompile(`(.+)SM`)
	match := regex.FindStringSubmatch(visibilityFlat)
	metarVisibility = match[1]
	return
}

func ParseDayTime(timeFlat string) (day int32, metarTime time.Time) {
	var day64 int64
	var timeString string
	regex := regexp.MustCompile(`(\d{2})(\d{4})Z`)
	match := regex.FindStringSubmatch(timeFlat)
	day64, _ = strconv.ParseInt(match[1], 10, 32)
	timeString = match[2]
	day = int32(day64)
	metarTime, _ = time.Parse("1504", timeString)
	return
}

func ParseCloudDescription(cloudFlat string) (cloud string) {
	regex := regexp.MustCompile(`(?P<code>\D\D\D)(?P<altitude>\d\d\d)`)
	matches := regex.FindStringSubmatch(cloudFlat)[1:]
	alt64, _ := strconv.ParseInt(matches[1], 10, 64)
	cloud = fmt.Sprintf("%v at %v", matches[0], alt64*100)
	return

}

func ParseClouds(cloudFlat string) (clouds []string) {
	regex := regexp.MustCompile(`\D\D\D\d\d\d`)
	matches := regex.FindAllString(cloudFlat, -1)
	for _, match := range matches {
		clouds = append(clouds, ParseCloudDescription(match))
	}
	return
}

func ParseSignedFloat(mBasedValue string) (signedFloat float32) {
	var signedFloat64 float64
	if strings.Index(mBasedValue, "M") == 0 {
		floatVal, _ := strconv.ParseFloat(mBasedValue[1:], 64)
		signedFloat64 = (-floatVal)
	} else {
		signedFloat64, _ = strconv.ParseFloat(mBasedValue, 64)
	}
	return float32(signedFloat64)
}

func ParseTempDew(tempDueFlat string) (temperature float32, dewPoint float32) {
	regex := regexp.MustCompile(`(M?\d\d)\/(M?\d\d)`)
	matches := regex.FindStringSubmatch(tempDueFlat)[1:]
	temperature = ParseSignedFloat(matches[0])
	dewPoint = ParseSignedFloat(matches[1])
	return
}

func ParsePressure(pressureFlat string) (pressure float32) {
	regex := regexp.MustCompile(`A(\d{4})`)
	matches := regex.FindStringSubmatch(pressureFlat)[1:]
	pressure = ParseSignedFloat(matches[0]) / 100
	return
}
