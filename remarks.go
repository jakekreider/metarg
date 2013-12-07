package main

import (
    "fmt"
	"regexp"
    "strconv"
    "strings"
)


func parseRemark(remark string) (translation string) {

    remarkMap := map[string]func(flatValue string) string{
            `^AO[1,2]$`:parseStationType,
            `^SLP\d\d\d$`:parseSeaLevelPressure,
            `^WEA\:something$`:parseWeatherAddl,
            `^PRES[FR]R$`:parsePressureChange,
            `^1\d{4}$`:parseMax6HrTemp,
            `^2\d{4}$`:parseMin6HrTemp,
            `^4\/\d{3}$`:parseSnowCoverage,
            `^5[01]\d{3}$`:parsePressureTendency,
            `^6\d{4}$`:parse6HourPrecipitation,
            `^7\d{4}$`:parse24HourPrecipitation,
            `^8/[lmh]$`:parseCloudType,
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
    translation = fmt.Sprintf("Max temp in 6 hrs:  %4.1f °C", floatValue)
    return 
}

func parseMin6HrTemp(remark string) (translation string){
    var floatValue float64
    expression := regexp.MustCompile(`2(\d{4})`)
    matches := expression.FindStringSubmatch(remark)
    floatValue = parseRemarkSignedValue(matches[1])
    translation = fmt.Sprintf("Min temp in 6 hrs:  %4.1f °C", floatValue)
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

func parse6HourPrecipitation(remark string) (translation string){
    var floatValue float64
    floatValue, _ = strconv.ParseFloat(remark[1:], 32)
    translation = fmt.Sprintf("6-hour precipitation: %4.1f\"", floatValue/100)
    return 
}

func parse24HourPrecipitation(remark string) (translation string){
    var floatValue float64
    floatValue, _ = strconv.ParseFloat(remark[1:], 32)
    translation = fmt.Sprintf("24-hour precipitation: %4.1f\"", floatValue/100)
    return 
}

func parseCloudType(remark string) (translation string){
    var cloudType string
    var code = remark[len(remark)-1:]
    switch code {
        case "l" : cloudType = "Low"
        case "m" : cloudType = "Medium"
        case "h" : cloudType = "High"
        default: return
    }

    translation = fmt.Sprintf("Clouds:  %s", cloudType)

    return 
}

func parsePressureTendency(remark string) (translation string){
    pressure := parseOneSignedFloat(remark[1:]) * .1

    translation = fmt.Sprintf("Pressure tendency:  %4.1f mb", pressure)

    return 
}
    
func parseOneSignedFloat(signedInteger string) (value float64){
    value, _ = strconv.ParseFloat(signedInteger[1:], 64)
    if signedInteger[0:1] == "1" {
        return -value
    }
    return value
}



