package main

import (
	"fmt"
	"math"

	"github.com/sfomuseum/go-learn/base"
	"github.com/sfomuseum/go-learn/linear"
)

// Sample data for demonstration
var X = [][]float64{
	{1}, {2}, {3}, {4}, {5},
}
var y = []float64{4, 6, 9, 11, 13}

// Training a linear regression model
func trainModel() *linear.LinearRegression {
	model := linear.NewLinearRegression()
	dataset := base.NewDataset(X, y, base.Predictor(len(X[0])), base.Response(0))
	err := model.Fit(dataset)
	if err != nil {
		fmt.Printf("Failed to train model: %v\n", err)
		return nil
	}
	return model
}

// Predicting future workforce needs
func predict(model *linear.LinearRegression, months int) float64 {
	return model.Predict([]float64{float64(months)})[0]
}

func main() {
	model := trainModel()
	if model == nil {
		return
	}

	// Predict workforce needs in 6 months
	predictedNeeds := predict(model, 6)
	fmt.Printf("Predicted workforce needs in 6 months: %.2f\n", predictedNeeds)
}