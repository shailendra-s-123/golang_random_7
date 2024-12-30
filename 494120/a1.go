package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Employee represents a single employee in the system.
type Employee struct {
	ID        int
	Name      string
	Role      string
	Performance int // Performance rating on a scale of 1-10
}

// EmployeeData represents a map of employees by ID.
type EmployeeData map[int]*Employee

// PerformanceEvaluation evaluates employee performance.
func PerformanceEvaluation(employees EmployeeData) {
	for _, emp := range employees {
		rand.Seed(time.Now().UnixNano())
		emp.Performance = rand.Intn(10) + 1
		log.Printf("Updated Performance of %s (%d): %d\n", emp.Name, emp.ID, emp.Performance)
	}
}

// ResourceAllocator allocates resources dynamically based on employee performance and roles.
func ResourceAllocator(employees EmployeeData, totalResources int) map[string]int {
	resourceAllocation := make(map[string]int)
	sumPerformance := 0

	// Calculate total performance
	for _, emp := range employees {
		sumPerformance += emp.Performance
	}

	// Allocate resources based on performance
	for _, emp := range employees {
		share := int(float64(emp.Performance) / float64(sumPerformance) * float64(totalResources))
		if share < 1 {
			share = 1
		}

		resourceAllocation[emp.Role] += share
		log.Printf("Allocated %d resources to %s role.\n", share, emp.Role)
	}

	return resourceAllocation
}

// PredictFutureNeeds predicts future workforce needs based on a growth factor.
func PredictFutureNeeds(employees EmployeeData, growthFactor float64) EmployeeData {
	futureEmployees := make(EmployeeData)
	for _, emp := range employees {
		futureEmployees[emp.ID] = emp
	}

	// For simplicity, we'll increase each role's count based on growth factor
	for role := range futureEmployees {
		var numEmployees int
		for id, emp := range futureEmployees {
			if emp.Role == role {
				numEmployees++
			}
		}
		futureEmployees[numEmployees+1] = &Employee{ID: numEmployees + 1, Name: fmt.Sprintf("New %s", role), Role: role, Performance: 5}
	}

	return futureEmployees
}

func main() {
	employees := EmployeeData{
		1: &Employee{ID: 1, Name: "Alice", Role: "Engineer", Performance: 7},
		2: &Employee{ID: 2, Name: "Bob", Role: "Designer", Performance: 8},
		3: &Employee{ID: 3, Name: "Charlie", Role: "Manager", Performance: 6},
		4: &Employee{ID: 4, Name: "David", Role: "Engineer", Performance: 9},
		5: &Employee{ID: 5, Name: "Eve", Role: "Designer", Performance: 8},
	}

	// Evaluate employee performance
	PerformanceEvaluation(employees)

	// Allocate resources
	totalResources := 100
	resourceAllocation := ResourceAllocator(employees, totalResources)

	// Predict future workforce needs
	growthFactor := 1.1 // 10% growth
	futureEmployees := PredictFutureNeeds(employees, growthFactor)

	// Output results
	log.Println("Current Resource Allocation:")
	for role, resources := range resourceAllocation {
		log.Printf("%s: %d resources\n", role, resources)
	}

	log.Println("\nPredicted Future Employees:")
	for id, emp := range futureEmployees {
		log.Printf("ID: %d, Name: %s, Role: %s, Performance: %d\n", emp.ID, emp.Name, emp.Role, emp.Performance)
	}
}