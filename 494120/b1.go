
package main

import (
    "fmt"
    "math/rand"
    "time"
)

type employee struct {
    Id       int
    Name     string
    Department string
    PerformanceScore float32
}

// Map to store employee data
var employeeData map[int]*employee

func init() {
    employeeData = make(map[int]*employee)
}

// Function to add new employees
func addEmployee(id int, name string, department string) {
    e := &employee{
        Id:       id,
        Name:     name,
        Department: department,
        PerformanceScore: 0.0,
    }
    employeeData[id] = e
}

// Function to update employee performance score
func updateEmployeePerformance(id int, score float32) {
    if e, ok := employeeData[id]; ok {
        e.PerformanceScore = score
    }
}

// Function to calculate average performance score for a department
func calculateAveragePerformance(department string) float32 {
    var totalScore float32
    var count int
    for _, e := range employeeData {
        if e.Department == department {
            totalScore += e.PerformanceScore
            count++
        }
    }
    if count == 0 {
        return 0.0
    }
    return totalScore / float32(count)
}

// Function to dynamically allocate resources based on performance scores
func allocateResources() map[string]int {
    allocation := make(map[string]int)
    for _, e := range employeeData {
        allocation[e.Department]++
    }
    return allocation
}

func main() {
    // Add sample employees
    addEmployee(1, "Alice", "Engineering")
    addEmployee(2, "Bob", "Sales")
    addEmployee(3, "Charlie", "Engineering")
    addEmployee(4, "David", "Marketing")
    addEmployee(5, "Emma", "Sales")

    // Update employee performance scores (randomly for demonstration)
    rand.Seed(time.Now().UnixNano())
    for id := range employeeData {
        updateEmployeePerformance(id, rand.Float32()*5)
    }

    // Calculate average performance scores for each department
    avgEngPerformance := calculateAveragePerformance("Engineering")
    avgSalesPerformance := calculateAveragePerformance("Sales")
    avgMarketingPerformance := calculateAveragePerformance("Marketing")

    fmt.Println("Average Performance Scores:")
    fmt.Println("Engineering:", avgEngPerformance)
    fmt.Println("Sales:", avgSalesPerformance)
    fmt.Println("Marketing:", avgMarketingPerformance)

    // Dynamically allocate resources based on performance scores
    resourceAllocation := allocateResources()

    fmt.Println("\nResource Allocation:")
    for department, count := range resourceAllocation {
        fmt.Println(department, ":", count)
    }
}  
