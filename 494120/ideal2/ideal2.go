package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Employee represents individual employee data
type Employee struct {
	ID            int
	Name          string
	Department    string
	Performance   float64 // A value between 0 and 1
	HoursWorked   int
	Skills        []string
}

// Workforce represents the workforce planning system
type Workforce struct {
	EmployeeData map[int]Employee // Using a map for quick lookup by employee ID
	DepartmentalData map[string][]int // Map of department names to employee IDs
}

// NewWorkforce initializes a new workforce system
func NewWorkforce() *Workforce {
	return &Workforce{
		EmployeeData:     make(map[int]Employee),
		DepartmentalData: make(map[string][]int),
	}
}

// AddEmployee adds a new employee to the workforce
func (w *Workforce) AddEmployee(e Employee) {
	w.EmployeeData[e.ID] = e
	w.DepartmentalData[e.Department] = append(w.DepartmentalData[e.Department], e.ID)
}

// PredictFutureNeeds uses a placeholder for ML-based prediction
func PredictFutureNeeds(currentData map[int]Employee) string {
	// Placeholder for ML model integration
	// Replace with actual ML predictions using libraries like gorgonia or external APIs
	return "Predicted need: 10% more resources in Engineering"
}

// EvaluatePerformance evaluates employees and returns low-performance employees
func (w *Workforce) EvaluatePerformance(threshold float64) []Employee {
	var underperformers []Employee
	for _, emp := range w.EmployeeData {
		if emp.Performance < threshold {
			underperformers = append(underperformers, emp)
		}
	}
	return underperformers
}

// Example of a simple tree structure (alternative data structure)
type TreeNode struct {
	ID       int
	Name     string
	Children []*TreeNode
}

// CreateSampleHierarchy generates a sample tree structure
func CreateSampleHierarchy() *TreeNode {
	root := &TreeNode{ID: 1, Name: "Company"}
	root.Children = append(root.Children, &TreeNode{ID: 2, Name: "Engineering"})
	root.Children = append(root.Children, &TreeNode{ID: 3, Name: "Marketing"})
	return root
}

func main() {
	rand.Seed(time.Now().UnixNano())
	workforce := NewWorkforce()

	// Adding employees
	for i := 1; i <= 10; i++ {
		workforce.AddEmployee(Employee{
			ID:            i,
			Name:          fmt.Sprintf("Employee %d", i),
			Department:    "Engineering",
			Performance:   rand.Float64(),
			HoursWorked:   rand.Intn(200),
			Skills:        []string{"Go", "Leadership"},
		})
	}

	// Performance evaluation
	lowPerformers := workforce.EvaluatePerformance(0.5)
	fmt.Println("Underperformers:", lowPerformers)

	// Predict future workforce needs
	fmt.Println(PredictFutureNeeds(workforce.EmployeeData))

	// Using an alternative data structure
	hierarchy := CreateSampleHierarchy()
	fmt.Printf("Hierarchy Root: %s\n", hierarchy.Name)
}
