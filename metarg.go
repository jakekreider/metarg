package main
 
import (
    "io/ioutil"
    "flag"
    "fmt"
    "net/http"
    "os"
    "strings"
)
const METAR_PATH = "http://weather.noaa.gov/pub/data/observations/metar/stations/"

var decode, verbose, help bool
var flagSet *flag.FlagSet

func init() {
    flagSet = new(flag.FlagSet)
    flagSet.BoolVar(&decode, "d", false, "Decode" )
    flagSet.BoolVar(&verbose, "v", false, "Be verbose" )
    flagSet.BoolVar(&help, "h", false, "Help" )
}

//Parse command-line args
func ParseArgs(arguments []string)  (flag.FlagSet, bool) {
    success := true
    err := flagSet.Parse(arguments)
    if err != nil || help   {
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
    resp, _ := http.Get(METAR_PATH+station+".TXT")
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return string(body), false
}

func main() {
    args, valid := ParseArgs(os.Args[1:])
    if valid {
        if len(args.Args()) == 0 {
            fmt.Print("Usage: metarg [options] station\n")

        }else{
            metar, _ := GetMetar(args.Args()[0]) 
            fmt.Print(metar)
        }
    }
}