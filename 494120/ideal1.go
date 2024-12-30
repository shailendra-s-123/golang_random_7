package main

import (
	"fmt"
)

// Employee struct holds information about an employee
type Employee struct {
	ID           int
	Name         string
	Role         string
	Performance  float64 // Performance score (0 to 10)
	Workload     int     // Current workload (number of tasks)
	Availability bool    // Indicates if the employee is available for new tasks
}

// Workforce struct to manage the workforce and employee data
type Workforce struct {
	employees map[int]*Employee
}

// NewWorkforce initializes the Workforce with an empty map
func NewWorkforce() *Workforce {
	return &Workforce{
		employees: make(map[int]*Employee),
	}
}

// AddEmployee adds a new employee to the workforce
func (w *Workforce) AddEmployee(emp Employee) {
	w.employees[emp.ID] = &emp
}

// EvaluatePerformance evaluates employee performance based on workload
func (w *Workforce) EvaluatePerformance() {
	for _, emp := range w.employees {
		if emp.Workload > 8 {
			emp.Performance -= 1 // Penalize for high workload
		} else if emp.Workload < 3 {
			emp.Performance += 0.5 // Reward for low workload
		}
		if emp.Performance > 10 {
			emp.Performance = 10
		} else if emp.Performance < 0 {
			emp.Performance = 0
		}
	}
}

// OptimizeResources dynamically redistributes tasks to balance workload
func (w *Workforce) OptimizeResources() {
	for {
		reassigned := false
		for _, emp := range w.employees {
			if emp.Workload > 8 {
				// Find a less-burdened employee
				for _, target := range w.employees {
					if target.Workload < 5 {
						emp.Workload--
						target.Workload++
						fmt.Printf("Reassigned a task from %s to %s\n", emp.Name, target.Name)
						reassigned = true
						break
					}
				}
			}
		}
		if !reassigned {
			break // Stop if no more tasks can be reassigned
		}
	}
}

// PredictFutureNeeds estimates future workloads based on current data
func (w *Workforce) PredictFutureNeeds() {
	fmt.Println("\nFuture Task Predictions:")
	for _, emp := range w.employees {
		predictedTasks := emp.Workload + int(emp.Performance*0.5) // Simple prediction formula
		fmt.Printf("Employee %d (%s) is projected to handle %d tasks.\n", emp.ID, emp.Name, predictedTasks)
	}
}

// DisplayWorkforce prints the current state of the workforce
func (w *Workforce) DisplayWorkforce() {
	fmt.Println("\nWorkforce Summary:")
	for _, emp := range w.employees {
		fmt.Printf("ID: %d, Name: %s, Role: %s, Performance: %.2f, Workload: %d\n",
			emp.ID, emp.Name, emp.Role, emp.Performance, emp.Workload)
	}
}

func main() {
	// Initialize workforce system
	workforce := NewWorkforce()

	// Add sample employees
	workforce.AddEmployee(Employee{ID: 1, Name: "Alice", Role: "Engineer", Performance: 8, Workload: 5})
	workforce.AddEmployee(Employee{ID: 2, Name: "Bob", Role: "Manager", Performance: 6, Workload: 9})
	workforce.AddEmployee(Employee{ID: 3, Name: "Charlie", Role: "Technician", Performance: 7, Workload: 3})

	// Evaluate employee performance
	workforce.EvaluatePerformance()

	// Optimize resources by redistributing workload
	workforce.OptimizeResources()

	// Predict future workforce needs
	workforce.PredictFutureNeeds()

	// Display the final state of the workforce
	workforce.DisplayWorkforce()
}