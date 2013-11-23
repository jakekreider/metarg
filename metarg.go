package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/mragh/metarg/compass"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
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
	Station, Phenomena, Visibility, WindDirection                             string
	Clouds                                                                    []string
	Time                                                                      time.Time
	WindSpeed, WindGust, Temperature, Dewpoint, Pressure, WindDirectionDegree float32
	Day                                                                       int32
}

func ParseWind(windFlat string) (direction string, speed float32, 
								dirDegrees float32, gust float32) {
	regex := regexp.MustCompile(`(\d{3})(\d+)(G(\d+))?KT`)
	match := regex.FindStringSubmatch(windFlat)
	dirDegrees64, _ := strconv.ParseInt(match[1], 10, 32)
	speed64, _ := strconv.ParseFloat(match[2], 32)
	dirDegrees = float32(dirDegrees64)
	direction = compass.GetCompassAbbreviation(dirDegrees)
	speed = float32(speed64)
	if match[4] != "" {
		gust64 , _ := strconv.ParseFloat(match[4], 32)
		gust = float32(gust64)
	}else{
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

func ParseMetar(flatMetar string) (metar Metar, success bool) {
	regex := regexp.MustCompile(
		`^(?P<station>\w{4})\s(?P<time>\w{7})\s(?P<wind>\w+)\s(?P<visibility>\w+)\s+(?P<clouds>.*)\s(?P<tempdue>M?\d\d\/M?\d\d)\s(?P<pressure>A\d{4})\s.*`)
	match := regex.FindStringSubmatch(flatMetar)
	if len(match) < 7 {
		return metar, false
	}
	metar.Station = match[1]
	metar.Day, metar.Time = ParseDayTime(match[2])
	//TODO wind direction
	metar.WindDirection, metar.WindSpeed, 
		metar.WindDirectionDegree, metar.WindGust = ParseWind(match[3])
	metar.Visibility = ParseVisibility(match[4])
	metar.Clouds = ParseClouds(match[5])
	metar.Temperature, metar.Dewpoint = ParseTempDew(match[6])
	metar.Pressure = ParsePressure(match[7])
	return metar, true
}

func GetDetailMetar(metar Metar) (details string) {
	const stringTemplate = `Station       : {{.Station}}
Day           : {{.Day}}
Time          : {{.Time.Format "15:04"}}
Wind direction: {{.WindDirectionDegree}} ({{.WindDirection}})
Wind speed    : {{.WindSpeed}} KT
Wind gust     : {{.WindGust}} KT
Visibility    : {{.Visibility}} SM
Temperature   : {{.Temperature}} C
Dewpoint      : {{.Dewpoint}} C
Pressure      : {{.Pressure}} "Hg
Clouds        : {{range .Clouds}}{{.}} ft {{end}}
Phenomena     :  //TODO`
	tmpl, err := template.New("metarDetail").Parse(stringTemplate)
	if err != nil {
		panic(err)
	}
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, metar)
	if err != nil {
		panic(err)
	}
	details = doc.String()
	return
}

func DecodeMetar(metarLine string) (details string, success bool) {
	metar, success := ParseMetar(metarLine)
	if !success {
		return details, success
	}
	details = GetDetailMetar(metar)
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
	if verbose {
		fmt.Print("Shh, not implemented yet ...")
		success = false
	}

	return *flagSet, success
}

//Retrieve the METAR for the given station
//Returns the string and a status
func GetMetar(station string) (value string, ok bool) {
	station = strings.ToUpper(station)
	resp, err := http.Get(METAR_PATH + station + ".TXT")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	metarLine := strings.Split(string(body), "\n")[1]
	if decode {
		decodedValue, ok := DecodeMetar(metarLine)
		if !ok {
			return value, ok
		}
		value = fmt.Sprintf("%s\n%s", metarLine, decodedValue)
	} else {
		value = metarLine
	}
	return value, true
}

func main() {
	args, valid := ParseArgs(os.Args[1:])
	if valid {
		if len(args.Args()) == 0 {
			fmt.Print("Usage: metarg [options] station\n")

		} else {
			metar, success := GetMetar(args.Args()[0])
			if success {
				fmt.Print(metar, "\n")

			} else {
				fmt.Print("Oh no, something went wrong!\n")
			}
		}
	}
}
