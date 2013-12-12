package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var Output io.Writer

const METAR_PATH = "http://weather.noaa.gov/pub/data/observations/metar/stations/"
const METAR_LIST_REF = "http://www.cnrfc.noaa.gov/metar.php"

var decode, verbose, search, help bool
var flagSet *flag.FlagSet

func init() {
	flagSet = new(flag.FlagSet)
	flagSet.BoolVar(&decode, "d", false, "Decode")
	flagSet.BoolVar(&verbose, "v", false, "Be verbose")
	flagSet.BoolVar(&search, "s", false, "Search")
	flagSet.BoolVar(&help, "h", false, "Help")
	Output = os.Stdout
}

//Command-line entry point
func main() {
	args, valid := ParseArgs(os.Args[1:])
	if valid {
		var result string
		var success bool
		if search {
			var resultList []string
			resultList, success = SearchStations(args.Args()[0])
			result = strings.Join(resultList, "\n")
		}else {
			result, success = GetMetar(args.Args())	
		}
		if success {
			fmt.Fprint(Output, result, "\n")

		} else {
			fmt.Fprint(Output, "Oh no, something went wrong!\n")
		}
	}
}

//Retrieve the METAR for the given station
//Returns the string and a status
func GetMetar(stations []string) (value string, ok bool) {
	for _, station := range stations {
		var stationMetar string
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
			stationMetar = fmt.Sprintf("%s\n%s", metarLine, decodedValue)
		} else {
			stationMetar = metarLine
		}

		value += stationMetar + "\n"
	}

	return value, true
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
		fmt.Fprintln(Output, "Usage: metarg [options] station")
		success = false
	}
	if verbose {
		fmt.Fprintln(Output, "Shh, not implemented yet ...")
		success = false
	}

	return *flagSet, success
}

func SearchStations(search string) (results []string, success bool) {
	search = strings.ToUpper(search)
	resp, err := http.Get(METAR_LIST_REF)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprint(Output, "ERROR ", err)
		return
	}
	body := string(byteBody)
	body = strings.Replace(body, "\n", "", -1)
	stationSearch := regexp.MustCompile(`<tr><td.*?>(?P<state>\w{2})</td>` +
				`<td.*?>(?P<name>.*?)</td><td.*?><a.*?>(?P<code>\w+)`)
	matches := stationSearch.FindAllStringSubmatch(body, -1)
	success = true
	for _, match := range matches {
		_, state, name, _ := match[0], match[1], match[2], match[3]
		if strings.Contains(name, search) || strings.Contains(state, search) {
			results = append(results, strings.Join(match[1:], " - "))
		}

	}
	return results, success
}

func GetDetailMetar(metar Metar) (details string) {
	const stringTemplate = `Station       : {{.Station}}
Day           : {{.Day}}
Time          : {{.Time.Format "15:04"}} UTC
Wind direction: {{.WindDirectionDegree}} ({{.WindDirection}})
Wind speed    : {{.WindSpeed}} KT
Wind gust     : {{.WindGust}} KT
Visibility    : {{.Visibility}}
Temperature   : {{.Temperature}} C
Dewpoint      : {{.Dewpoint}} C
Pressure      : {{.Pressure}} "Hg
Clouds        : {{range .Clouds}}{{.}} ft {{end}}
Remarks       : 
{{range .Remarks}}{{.}}
{{end}}`
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
