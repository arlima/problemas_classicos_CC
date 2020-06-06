package main

import (
	"fmt"
	"math"
)

func permutations(arr []string) [][]string {
	var helper func([]string, int)
	res := [][]string{}

	helper = func(arr []string, n int) {
		if n == 1 {
			tmp := make([]string, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func main() {
	vtDistances := make(map[string]map[string]int)
	vtDistances["Rutland"] = make(map[string]int)
	vtDistances["Rutland"]["Burlington"] = 67
	vtDistances["Rutland"]["White River Junction"] = 46
	vtDistances["Rutland"]["Bennington"] = 55
	vtDistances["Rutland"]["Battleboro"] = 75

	vtDistances["Burlington"] = make(map[string]int)
	vtDistances["Burlington"]["Rutland"] = 67
	vtDistances["Burlington"]["White River Junction"] = 91
	vtDistances["Burlington"]["Bennington"] = 122
	vtDistances["Burlington"]["Battleboro"] = 153

	vtDistances["White River Junction"] = make(map[string]int)
	vtDistances["White River Junction"]["Rutland"] = 46
	vtDistances["White River Junction"]["Burlington"] = 91
	vtDistances["White River Junction"]["Bennington"] = 98
	vtDistances["White River Junction"]["Battleboro"] = 65

	vtDistances["Bennington"] = make(map[string]int)
	vtDistances["Bennington"]["Rutland"] = 55
	vtDistances["Bennington"]["Burlington"] = 122
	vtDistances["Bennington"]["White River Junction"] = 98
	vtDistances["Bennington"]["Battleboro"] = 40

	vtDistances["Battleboro"] = make(map[string]int)
	vtDistances["Battleboro"]["Rutland"] = 75
	vtDistances["Battleboro"]["Burlington"] = 153
	vtDistances["Battleboro"]["White River Junction"] = 65
	vtDistances["Battleboro"]["Bennington"] = 40

	vtCities := make([]string, 0, len(vtDistances))
	for k := range vtDistances {
		vtCities = append(vtCities, k)
	}

	cityPermutations := permutations(vtCities)

	tspPaths := [][]string{}

	for _, c := range cityPermutations {
		tspPaths = append(tspPaths, append(c, c[0]))
	}

	bestPath := []string{}
	minDistance := math.MaxInt16

	for _, path := range tspPaths {
		distance := 0
		last := path[0]
		for _, next := range path[1:] {
			distance += vtDistances[last][next]
			last = next
		}
		if distance < minDistance {
			minDistance = distance
			bestPath = path
		}
	}

	fmt.Printf("The shortest path is : %v\n", bestPath)
	fmt.Printf("The distance is : %d miles\n", minDistance)
}
