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

func interpretOutput(output []float64) string {
	if floats.Max(output) == output[0] {
		return "Iris-setosa"
	} else if floats.Max(output) == output[1] {
		return "Iris-versicolor"
	} else {
		return "Iris-virginica"
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	irisParameters := [][]float64{}
	irisClassifications := [][]float64{}
	irisSpecies := []string{}

	file, err := os.Open("iris.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	k := 0
	irises := [][]string{}
	for scanner.Scan() {
		res := scanner.Text()
		str := strings.Split(res, ",")
		irises = append(irises, str)
		k++
	}

	rand.Shuffle(len(irises), func(i, j int) { irises[i], irises[j] = irises[j], irises[i] })

	for _, iris := range irises {
		parameters := []float64{}
		for _, n := range iris[0:4] {
			value, _ := strconv.ParseFloat(n, 1)
			parameters = append(parameters, float64(value))
		}
		irisParameters = append(irisParameters, parameters)

		species := iris[4]

		if species == "Iris-setosa" {
			irisClassifications = append(irisClassifications, []float64{1.0, 0.0, 0.0})
		} else if species == "Iris-versicolor" {
			irisClassifications = append(irisClassifications, []float64{0.0, 1.0, 0.0})
		} else {
			irisClassifications = append(irisClassifications, []float64{0.0, 0.0, 1.0})
		}
		irisSpecies = append(irisSpecies, species)
	}
	normalizeByFeatureScaling(irisParameters)

	net := network{}
	net.init([]int{4, 6, 3}, 0.3, sigmoid, derivationSigmoid)

	irisTrainers := irisParameters[0:140]
	irisTrainersCorrects := irisClassifications[0:140]

	for t := 0; t < 50; t++ {
		net.train(irisTrainers, irisTrainersCorrects)
	}

	irisTesters := irisParameters[140:150]
	irisTestersCorrects := irisSpecies[140:150]

	correct, inputs, percentage := net.validate(irisTesters, irisTestersCorrects, interpretOutput)

	fmt.Println(correct, inputs, percentage)

}
