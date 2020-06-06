package main

import (
	"fmt"
	"log"
)

type network struct {
	layers []layer
}

func (n *network) init(layerStructure []int, learningRate float64, activationFunction activationFunctionT,
	derivativeActivationFunction derivativeActivationFunctionT) {
	if len(layerStructure) < 3 {
		log.Fatal("Error: Should be at least 3 layers (1 input, 1 hidden, 1 output)")
	}
	n.layers = []layer{}
	inputLayer := layer{}
	inputLayer.init(nil, layerStructure[0], learningRate, activationFunction, derivativeActivationFunction)
	n.layers = append(n.layers, inputLayer)

	for layerNum, numNeurons := range layerStructure[1:] {
		nextLayer := layer{}
		nextLayer.init(&n.layers[layerNum], numNeurons, learningRate, activationFunction, derivativeActivationFunction)
		n.layers = append(n.layers, nextLayer)
	}
}

func (n *network) outputs(input []float64) []float64 {
	output := n.layers[0].outputs(input)
	for l := 1; l < len(n.layers); l++ {
		output = n.layers[l].outputs(output)
	}
	return output
}

func (n *network) backPropagate(expected []float64) {
	lastLayer := len(n.layers) - 1
	n.layers[lastLayer].calculateDeltasforOutputLayer(expected)
	for l := lastLayer - 1; l > 0; l-- {
		n.layers[l].calculateDeltasforHiddenLayer(n.layers[l+1])
	}
}

func (n *network) updateWeights() {
	for li := 1; li < len(n.layers); li++ {
		for ni := range n.layers[li].neurons {
			for w := 0; w < len(n.layers[li].neurons[ni].weights); w++ {
				n.layers[li].neurons[ni].weights[w] = n.layers[li].neurons[ni].weights[w] +
					(n.layers[li].neurons[ni].learningRate * (n.layers[li-1].outputCache[w]) * n.layers[li].neurons[ni].delta)
			}
		}
	}
}

func (n *network) train(inputs [][]float64, expecteds [][]float64) {
	for location, xs := range inputs {
		ys := expecteds[location]
		_ = n.outputs(xs)
		n.backPropagate(ys)
		n.updateWeights()
	}
}

type interpretOutputT func([]float64) int

func (n *network) validate(inputs [][]float64, expecteds []int, interpretOutput interpretOutputT) (int, int, float64) {
	correct := 0
	for i := 0; i < len(inputs); i++ {
		input := inputs[i]
		expected := expecteds[i]
		result := interpretOutput(n.outputs(input))
		if result == expected {
			correct++
		}
	}
	percentage := float64(correct) / float64(len(inputs))
	return correct, len(inputs), percentage
}

func (n *network) print() {
	fmt.Println("Camadas", len(n.layers))
	for l := 0; l < len(n.layers); l++ {
		fmt.Printf("Camada %d: %d neuronios\n", l, len(n.layers[l].neurons))
		fmt.Printf("outputCache da camada %f\n", n.layers[l].outputCache)
		for neu := 0; neu < len(n.layers[l].neurons); neu++ {
			fmt.Printf("Camada %d: Neuronio %d, %d pesos\n", l, neu, len(n.layers[l].neurons[neu].weights))
			fmt.Printf("outputCache %f, delta %f\n", n.layers[l].neurons[neu].outputCache, n.layers[l].neurons[neu].delta)
			for w := 0; w < len(n.layers[l].neurons[neu].weights); w++ {
				fmt.Printf("Camada %d: Neuronio %d: Pesos %d: %f\n", l, neu, w, n.layers[l].neurons[neu].weights[w])
			}
		}
	}
}
