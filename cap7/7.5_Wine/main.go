package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gonum/floats"
)

func interpretOutput(output []float64) int {
	if floats.Max(output) == output[0] {
		return 1
	} else if floats.Max(output) == output[1] {
		return 2
	} else {
		return 3
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	wineParameters := [][]float64{}
	wineClassifications := [][]float64{}
	wineSpecies := []int{}

	file, err := os.Open("wine.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	k := 0
	wines := [][]string{}
	for scanner.Scan() {
		res := scanner.Text()
		str := strings.Split(res, ",")
		wines = append(wines, str)
		k++
	}

	rand.Shuffle(len(wines), func(i, j int) { wines[i], wines[j] = wines[j], wines[i] })

	for _, wine := range wines {
		parameters := []float64{}
		for _, n := range wine[1:14] {
			value, _ := strconv.ParseFloat(n, 1)
			parameters = append(parameters, float64(value))
		}
		wineParameters = append(wineParameters, parameters)

		species, _ := strconv.Atoi(wine[0])

		if species == 1 {
			wineClassifications = append(wineClassifications, []float64{1.0, 0.0, 0.0})
		} else if species == 2 {
			wineClassifications = append(wineClassifications, []float64{0.0, 1.0, 0.0})
		} else {
			wineClassifications = append(wineClassifications, []float64{0.0, 0.0, 1.0})
		}
		wineSpecies = append(wineSpecies, species)
	}
	normalizeByFeatureScaling(wineParameters)

	net := network{}
	net.init([]int{13, 7, 3}, 0.9, sigmoid, derivationSigmoid)

	wineTrainers := wineParameters[0:150]
	wineTrainersCorrects := wineClassifications[0:150]

	for t := 0; t < 10; t++ {
		net.train(wineTrainers, wineTrainersCorrects)
	}

	wineTesters := wineParameters[150:178]
	wineTestersCorrects := wineSpecies[150:178]

	correct, inputs, percentage := net.validate(wineTesters, wineTestersCorrects, interpretOutput)

	fmt.Println(correct, inputs, percentage)

}
