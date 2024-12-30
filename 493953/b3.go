package main
 
import (
    "fmt"
    "log"
    "math"
    "math/rand"
    "net/http"
    "sync"
    "sync/atomic"
    "time"
)

// PerformanceMetrics holds the server's performance data
type PerformanceMetrics struct {
    // ... (Previous metrics fields)
    AverageLatency float64 // Smoothed average latency
    Alpha          float64 // Exponential smoothing alpha for average latency
    TrendLatency   float64 // Trend in average latency
    // For anomaly detection
    TrafficZScoreThreshold float64 // Z-score threshold for anomaly detection
}

func (m *PerformanceMetrics) UpdateSmoothing(latency time.Duration) {
    // Exponential smoothing for latency
    m.AverageLatency = m.Alpha*float64(latency) + (1-m.Alpha)*m.AverageLatency
    m.TrendLatency = m.Alpha*m.TrendLatency + (1-m.Alpha)*(m.AverageLatency-m.TrendLatency)
}

// Placeholder functions for machine learning-based traffic prediction
// Implement real-world machine learning models here
func (m *PerformanceMetrics) PredictTraffic() int {
    // Simple random prediction for demonstration
    return rand.Intn(20) + 5
}

func (m *PerformanceMetrics) TrainPredictor(pastData []time.Duration) {
    // Placeholder for training a machine learning model on past data
}

//sendRequest sends a GET request and updates performance metrics
func sendRequest(url string, metrics *PerformanceMetrics, wg *sync.WaitGroup) { 
    // ... (Previous implementation)

    // Update latency smoothing
    metrics.UpdateSmoothing(duration)
}

// adaptiveLoad generates traffic with adaptive throttling and predictive loads
func adaptiveLoad(url string, initialRate int, duration time.Duration, metrics *PerformanceMetrics) {    
    // ... (Previous implementation)
    
    // Predict and adjust load based on the predicted traffic and latency trend
    go func() {
        for {
            predictedTraffic := metrics.PredictTraffic()
            nextRate := rate + predictedTraffic
            adjustedRate := max(min(nextRate, initialRate*3), initialRate/3)
            rate = adjustedRate
            
            fmt.Printf("Predicted Traffic: %d, Adjusted Rate: %d, Latency Trend: %.2f ms\n", predictedTraffic, adjustedRate, metrics.TrendLatency/time.Millisecond)
            time.Sleep(3 * time.Second) // Check for prediction and adjust rate every 3 seconds
        }
    }()
}

// zScore returns the z-score of the value based on the history
func zScore(value, mean, stdDev float64) float64 {
    if stdDev == 0 {
        return 0
    }
    return (value - mean) / stdDev
}

// predictAnomaly determines if the value is an anomaly based on a z-score threshold
func predictAnomaly(z float64, threshold float64) bool {
    return math.Abs(z) >= threshold
}


// placeHolderAnomalyDetection performs z-score analysis on historical data for anomaly detection (in practice, use actual anomaly detection models)
func placeHolderAnomalyDetection(metrics *PerformanceMetrics, zScoreThreshold float64) {
    // Run this task periodically (e.g., every minute)
    go func() {
        for {
            // Use a sliding window to calculate mean and std dev
            // Here, just use sample data for demonstration
            recentSamples := []time.Duration{metrics.LastLatencySample, time.Millisecond*120, time.Millisecond*140}

            var total time.Duration
            for _, dur := range recentSamples {
                total += dur
            }
            meanLatency := time.Duration(float64(total) / float64(len(recentSamples)))