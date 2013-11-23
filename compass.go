package main

import "sort"

var namePoints = [...]string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", 
            "S", "SSW", "SW", "WSW", "W", "NWN", "NW", "NNW"}

var breakpoints = [...]float64{16.88, 39.38, 61.88, 84.38, 106.88, 129.38, 151.88, 174.38,
                196.88, 219.38, 241.88, 264.38, 286.88, 309.38, 331.88, 343.12}

//Get the compass abbreviation for cardinal or principal point (SE, NE, WSW, etc.) for the given degrees
func GetCompassAbbreviation(point float32) (compassPoint string) {
    var index = sort.SearchFloat64s(breakpoints[:], float64(point))
    if index == len(namePoints) {
        index = 0
    }
    compassPoint = namePoints[index]
    return
}