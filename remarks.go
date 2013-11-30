package main

import (
    "fmt"
	"regexp"
    "strconv"
    "strings"
)


func parseRemark(remark string) (translation string) {

    remarkMap := map[string]func(flatValue string) string{
            `AO[1,2]`:parseStationType,
            `SLP\d\d\d`:parseSeaLevelPressure,
            `WEA\:something`:parseWeatherAddl,
            `PRES[FR]R`:parsePressureChange,
            `1\d{4}`:parseMax6HrTemp,
            `2\d{4}`:parseMin6HrTemp,
            `4\/\d{3}`:parseSnowCoverage,
    }
    for rgx, evaluator := range remarkMap {
        expression := regexp.MustCompile(rgx)
        if expression.MatchString(remark) {
            return evaluator(remark)
        }
    }
    return ""
}

func parseStationType(remark string) (translation string){
    if strings.HasSuffix(remark, "1") {
        return "AMOS station"
    } else if strings.HasSuffix(remark, "2") {
        return "ASOS station"
    }
    return ""
}

func parseSeaLevelPressure(remark string) (translation string){
    pressure, _ := strconv.ParseFloat(remark[3:], 64)
    pressure = pressure/10
    return fmt.Sprintf("Sea level pressure %v mb", pressure)
}

func parseWeatherAddl(remark string) (translation string){
    translation = remark[4:]
    return
}

func parsePressureChange(remark string) (translation string){
    if strings.HasSuffix(remark, "RR") {
        return "Pressure rising rapidly"
    } else if strings.HasSuffix(remark, "FR") {
        return "Pressure falling rapidly"
    }
    return ""
}

func parseMax6HrTemp(remark string) (translation string){
    var floatValue float64
    expression := regexp.MustCompile(`1(\d{4})`)
    matches := expression.FindStringSubmatch(remark)
    floatValue = parseRemarkSignedValue(matches[1])
    translation = fmt.Sprintf("Max temp in 6 hrs:  %4.1f", floatValue)
    return 
}

func parseMin6HrTemp(remark string) (translation string){
    var floatValue float64
    expression := regexp.MustCompile(`2(\d{4})`)
    matches := expression.FindStringSubmatch(remark)
    floatValue = parseRemarkSignedValue(matches[1])
    translation = fmt.Sprintf("Min temp in 6 hrs:  %4.1f", floatValue)
    return 
}

func parseRemarkSignedValue(value string) (floatValue float64) {
    floatValue, _ = strconv.ParseFloat(value[1:], 32)
    if strings.HasPrefix(value, "1") {
        floatValue = -.1 * floatValue
    } else {
        floatValue = .1 * floatValue
    }
    return
}

func parseSnowCoverage(remark string) (translation string){
    var floatValue float64
    expression := regexp.MustCompile(`4\/(\d{3})`)
    matches := expression.FindStringSubmatch(remark)
    floatValue, _ = strconv.ParseFloat(matches[1], 32)
    translation = fmt.Sprintf("Snow coverage:  %4.1f\"", floatValue)
    return 
}