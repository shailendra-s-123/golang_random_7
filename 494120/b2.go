package main

import (
	"fmt"
	"gonum.org/v1/gonum/stat/regression"
)

// Employee struct remains the same

// Workforce struct remains the same

// New function to train a linear regression model for performance prediction
func (w *Workforce) TrainPerformancePredictionModel() *regression.Regression {
	// Prepare data for regression
	x := make([]float64, len(w.employees))
	y := make([]float64, len(w.employees))
	i := 0
	for _, emp := range w.employees {
		x[i] = float64(emp.Workload)
		y[i] = emp.Performance
		i++
	}

	// Train the regression model
	r := regression.New(regression.WithIntercept(true))
	r.Train(regression.Data(x, y))
	return r
}

// Use the trained model to predict performance for new employees
func (w *Workforce) PredictPerformance(workload int) float64 {
	r := w.TrainPerformancePredictionModel()
	predictedPerformance := r.Predict(float64(workload))
	return predictedPerformance
}

func main() {
	// Initialize workforce system
	workforce := NewWorkforce()
	// Add sample employees as before

	// Train the performance prediction model
	workforce.TrainPerformancePredictionModel()