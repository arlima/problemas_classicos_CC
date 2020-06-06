package main

import "math"

type datapoint struct {
	originals  []float64
	dimensions []float64
	label      string
}

func (d *datapoint) init(initial []float64, label string) {
	for _, v := range initial {
		d.originals = append(d.originals, v)
		d.dimensions = append(d.dimensions, v)
	}
	d.label = label
}

func (d *datapoint) numDimensions() int {
	return len(d.dimensions)
}

func (d *datapoint) distance(other datapoint) float64 {
	distance := 0.0
	for i := 0; i < len(d.dimensions); i++ {
		distance += math.Pow(d.dimensions[i]-other.dimensions[i], 2.0)
	}
	return math.Sqrt(distance)
}
