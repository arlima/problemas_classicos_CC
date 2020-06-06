package main

type activationFunctionT func(float64) float64
type derivativeActivationFunctionT func(float64) float64

type neuron struct {
	weights                      []float64
	activationFunction           activationFunctionT
	derivativeActivationFunction derivativeActivationFunctionT
	learningRate                 float64
	outputCache                  float64
	delta                        float64
}

func (n *neuron) init(weights []float64, activationFunction activationFunctionT,
	derivativeActivationFunction derivativeActivationFunctionT, learningRate float64) {
	n.weights = weights
	n.activationFunction = activationFunction
	n.derivativeActivationFunction = derivativeActivationFunction
	n.learningRate = learningRate
}

func (n *neuron) output(inputs []float64) float64 {
	n.outputCache = dotProduct(inputs, n.weights)
	return n.activationFunction(n.outputCache)
}
