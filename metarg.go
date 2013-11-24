package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

//Parse command-line args
func ParseArgs(arguments []string) (flag.FlagSet, bool) {
	success := true
	err := flagSet.Parse(arguments)
	if err != nil || help {
		flag.PrintDefaults()
		success = false
	}
	if len(flagSet.Args()) == 0 {
		fmt.Println("Usage: metarg [options] station")
		success = false
	}
	if verbose {
		fmt.Print("Shh, not implemented yet ...")
		success = false
	}

	return *flagSet, success
}

type Metar struct {
	Station, Phenomena, Visibility, WindDirection                             string
	Clouds                                                                    []string
	Time                                                                      time.Time
	WindSpeed, WindGust, Temperature, Dewpoint, Pressure, WindDirectionDegree float32
	Day                                                                       int32
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
		metar, success := GetMetar(args.Args()[0])
		if success {
			fmt.Print(metar, "\n")

		} else {
			fmt.Print("Oh no, something went wrong!\n")
		}
	}
}
