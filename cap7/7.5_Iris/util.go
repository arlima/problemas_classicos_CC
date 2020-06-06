package main

import (
	"math"

	"gonum.org/v1/gonum/floats"
)

func normalizeByFeatureScaling(dataset [][]float64) {
	for colNum := 0; colNum < len(dataset[0]); colNum++ {
		column := []float64{}
		for _, row := range dataset {
			column = append(column, row[colNum])
		}
		maximum := floats.Max(column)
		minimum := floats.Min(column)
		for rowNum := 0; rowNum < len(dataset); rowNum++ {
			dataset[rowNum][colNum] = (dataset[rowNum][colNum] - minimum) / (maximum - minimum)
		}
	}
}

func dotProduct(xs []float64, ys []float64) float64 {
	sum := 0.0
	for i := 0; i < len(xs); i++ {
		sum += xs[i] * ys[i]
	}
	return sum
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func derivationSigmoid(x float64) float64 {
	sig := sigmoid(x)
	return sig * (1.0 - sig)
}
