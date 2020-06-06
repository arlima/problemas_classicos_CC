package main

import (
	"math/rand"
)

type layer struct {
	previousLayer *layer
	neurons       []neuron
	outputCache   []float64
}

func (l *layer) init(previousLayer *layer, numNeurons int, learningRate float64, activationFunction activationFunctionT,
	derivativeActivationFunction derivativeActivationFunctionT) {
	l.previousLayer = previousLayer
	l.neurons = []neuron{}

	for i := 0; i < numNeurons; i++ {
		randomWeights := []float64{}
		if previousLayer != nil {
			for n := 0; n < len(previousLayer.neurons); n++ {
				randomWeights = append(randomWeights, rand.Float64())
			}
		}
		nNeuron := neuron{}
		nNeuron.init(randomWeights, activationFunction, derivativeActivationFunction, learningRate)
		l.neurons = append(l.neurons, nNeuron)
	}
	l.outputCache = make([]float64, numNeurons)
}

func (l *layer) outputs(inputs []float64) []float64 {
	if l.previousLayer == nil {
		l.outputCache = inputs
	} else {
		out := []float64{}
		for i := range l.neurons {
			out = append(out, l.neurons[i].output(inputs))
		}
		l.outputCache = out
	}
	return l.outputCache
}

func (l *layer) calculateDeltasforOutputLayer(expected []float64) {
	for n := 0; n < len(l.neurons); n++ {
		l.neurons[n].delta = l.neurons[n].derivativeActivationFunction(l.neurons[n].outputCache) * (expected[n] - l.outputCache[n])
	}
}

func (l *layer) calculateDeltasforHiddenLayer(nextLayer layer) {
	for index := range l.neurons {
		nextWeights := []float64{}
		nextDeltas := []float64{}
		for _, n := range nextLayer.neurons {
			nextWeights = append(nextWeights, n.weights[index])
			nextDeltas = append(nextDeltas, n.delta)
		}
		sumWeightsAndDeltas := dotProduct(nextWeights, nextDeltas)
		l.neurons[index].delta = l.neurons[index].derivativeActivationFunction(l.neurons[index].outputCache) * sumWeightsAndDeltas
	}
}
