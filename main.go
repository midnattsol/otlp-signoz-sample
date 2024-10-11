package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"main/api"
	"main/otlp"

	"go.opentelemetry.io/otel"
)

func main() {
	// Initialize the MeterProvider for metrics
	mp, err := otlp.InitMeterProvider()
	if err != nil {
		log.Fatalf("Error initializing MeterProvider: %v", err)
	}
	defer func() {
		if err := mp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down MeterProvider: %v", err)
		}
	}()

	meter := otel.Meter("POC")
	metrics := otlp.NewMetrics(meter)

	http.HandleFunc("/add-to-cart", func(w http.ResponseWriter, r *http.Request) {
		api.AddToCartHandler(w, r, metrics)
	})
	http.HandleFunc("/remove-from-cart", func(w http.ResponseWriter, r *http.Request) {
		api.RemoveFromCartHandler(w, r, metrics)
	})

	// Start the server
	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
