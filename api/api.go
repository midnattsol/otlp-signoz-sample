package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"main/otlp" // Import the otlp package
)

// ProcessHandler simulates processing a request and records latency
func ProcessHandler(w http.ResponseWriter, r *http.Request, metrics *otlp.Metrics) {
	ctx := r.Context()
	startTime := time.Now()
	statusCode := http.StatusOK

	// Simulate a 30% chance of an error
	if rand.Float64() < 0.3 {
		statusCode = http.StatusInternalServerError
		metrics.RegisterError(ctx, statusCode)
	}

	processRequest()

	duration := time.Since(startTime).Seconds()
	metrics.LatencyHist.Record(ctx, duration)

	w.WriteHeader(statusCode)
	if statusCode == http.StatusOK {
		fmt.Fprintf(w, "Request processed successfully")
	} else {
		http.Error(w, "Error processing the request", statusCode)
	}
}

// AddToCartHandler increments the cart items counter
func AddToCartHandler(w http.ResponseWriter, r *http.Request, metrics *otlp.Metrics) {
	metrics.UpdateCartItems(1)
	fmt.Fprintf(w, "Item added to cart. Total items: %d", metrics.CartItems.Count)
}

// RemoveFromCartHandler decrements the cart items counter
func RemoveFromCartHandler(w http.ResponseWriter, r *http.Request, metrics *otlp.Metrics) {
	metrics.UpdateCartItems(-1)
	fmt.Fprintf(w, "Item removed from cart. Total items: %d", metrics.CartItems.Count)
}

// processRequest simulates a processing delay
func processRequest() {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}
