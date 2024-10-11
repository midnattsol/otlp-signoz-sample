package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"main/otlp" // Import the otlp package
)

// AddToCartHandler increments the cart items counter and records latency
func AddToCartHandler(w http.ResponseWriter, r *http.Request, metrics *otlp.Metrics) {
	ctx := r.Context()
	startTime := time.Now()
	statusCode := http.StatusOK

	// Simulate a 30% chance of an error when adding to cart
	if rand.Float64() < 0.3 {
		statusCode = http.StatusInternalServerError
		metrics.RegisterError(ctx, statusCode)
	} else {
		metrics.UpdateCartItems(1)
	}
	processRequest()

	duration := time.Since(startTime).Seconds()
	metrics.RegisterLatency(ctx, duration)

	w.WriteHeader(statusCode)
	if statusCode == http.StatusOK {
		fmt.Fprintf(w, "Item added to cart. Total items: %d", metrics.CartItems.Count)
	} else {
		http.Error(w, "Error adding item to cart", statusCode)
	}
}

// RemoveFromCartHandler decrements the cart items counter and records latency
func RemoveFromCartHandler(w http.ResponseWriter, r *http.Request, metrics *otlp.Metrics) {
	ctx := r.Context()
	startTime := time.Now()
	statusCode := http.StatusOK

	// Simulate a 10% chance of an error when removing from cart
	if rand.Float64() < 0.1 {
		statusCode = http.StatusInternalServerError
		metrics.RegisterError(ctx, statusCode)
	} else {
		metrics.UpdateCartItems(-1)
	}
	processRequest()
	duration := time.Since(startTime).Seconds()
	metrics.RegisterLatency(ctx, duration)

	w.WriteHeader(statusCode)
	if statusCode == http.StatusOK {
		fmt.Fprintf(w, "Item removed from cart. Total items: %d", metrics.CartItems.Count)
	} else {
		http.Error(w, "Error removing item from cart", statusCode)
	}
}

// processRequest simulates a processing delay
func processRequest() {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}
