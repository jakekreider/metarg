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