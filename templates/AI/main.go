package main

import (
	"fmt"
	"math"
)

func main() {
	var inputs []float64 = []float64{0.00, 1.00, 1.00, 0.00}
	var weights []float64 = []float64{0.00, 0.00, 0.00, 0.00}
	var desired float64 = 1.00
	var learningRate float64 = 1.20
	var trials int = 100
	train(trials, inputs, weights, desired, learningRate)
}

func train(trials int, inputs []float64, weights []float64, desired float64, learningRate float64) {
	for i := 1; i < trials; i++ {
		weights = learn(inputs, weights, learningRate)
		output := evaluate(inputs, weights)
		errorResult := evaluateError(desired, output)

		fmt.Print("Output: ")
		fmt.Print(math.Round(output*100) / 100)
		fmt.Print("\nError: ")
		fmt.Print(math.Round(errorResult*100) / 100)
		fmt.Print("\n\n")
	}
}

func learn(inputVector []float64, weightVector []float64, learningRate float64) []float64 {
	for index, inputValue := range inputVector {
		if inputValue > 0.00 {
			weightVector[index] = weightVector[index] + learningRate
		}
	}
	return weightVector
}

func evaluate(inputVector []float64, weightVector []float64) float64 {
	var result float64 = 0.00
	for index, inputValue := range inputVector {
		layerValue := inputValue * weightVector[index]
		result = result + layerValue
	}
	return result
}

func evaluateError(desired float64, actual float64) float64 {
	return desired - actual
}
