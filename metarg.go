package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const METAR_PATH = "http://weather.noaa.gov/pub/data/observations/metar/stations/"

var decode, verbose, help bool
var flagSet *flag.FlagSet

func init() {
	flagSet = new(flag.FlagSet)
	flagSet.BoolVar(&decode, "d", false, "Decode")
	flagSet.BoolVar(&verbose, "v", false, "Be verbose")
	flagSet.BoolVar(&help, "h", false, "Help")
}

type Metar struct {
	station, phenomena, visibility                       string
	clouds                                               []string
	time                                                 time.Time
	windSpeed, windGust, temperature, dewPoint, pressure float32
	day                                                  int32
}

func ParseWind(windFlat string) (direction int32, speed float32) {
	regex := regexp.MustCompile(`(\d{3})(\d+)KT`)
	match := regex.FindStringSubmatch(windFlat)
	dir64, _ := strconv.ParseInt(match[1], 10, 32)
	speed64, _ := strconv.ParseFloat(match[2], 32)
	direction = int32(dir64)
	speed = float32(speed64)
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
	cloud = fmt.Sprintf("%v at %v", matches[0], alt64 * 100)
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

func ParseMetar(flatMetar string) (metar Metar) {
	regex := regexp.MustCompile(
		`^(?P<station>\w{4})\s(?P<time>\w{7})\s(?P<wind>\w+)\s(?P<visibility>\w+)\s+(?P<clouds>.*)\s(?P<tempdue>M?\d\d\/M?\d\d)\s.*`)
	match := regex.FindStringSubmatch(flatMetar)
	metar.station = match[1]
	metar.day, metar.time = ParseDayTime(match[2])
	//TODO wind direction
	_, metar.windSpeed = ParseWind(match[3])
	metar.visibility = ParseVisibility(match[4])
	metar.clouds = ParseClouds(match[5])
	return
}

//Parse command-line args
func ParseArgs(arguments []string) (flag.FlagSet, bool) {
	success := true
	err := flagSet.Parse(arguments)
	if err != nil || help {
		flag.PrintDefaults()
		success = false
	}
	if decode || verbose {
		fmt.Print("Shh, not implemented yet ...")
		success = false
	}

	return *flagSet, success
}

//Retrieve the METAR for the given station
//Returns the string and a status
func GetMetar(station string) (value string, ok bool) {
	station = strings.ToUpper(station)
	resp, _ := http.Get(METAR_PATH + station + ".TXT")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), false
}

func main() {
	args, valid := ParseArgs(os.Args[1:])
	if valid {
		if len(args.Args()) == 0 {
			fmt.Print("Usage: metarg [options] station\n")

		} else {
			metar, _ := GetMetar(args.Args()[0])
			fmt.Print(metar)
		}
	}
}
