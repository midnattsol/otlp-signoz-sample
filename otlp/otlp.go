package otlp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc/credentials"
)

type Metrics struct {
	ErrorCounter metric.Int64Counter
	LatencyHist  metric.Float64Histogram
	CartItems    struct {
		Count int64
		Gauge metric.Int64ObservableGauge
	}
	mutex sync.Mutex
}

func NewMetrics(meter metric.Meter) *Metrics {
	metrics := &Metrics{}
	errorCounter, err := meter.Int64Counter(
		"poc_error_requests",
		metric.WithDescription("Counter for error requests"),
	)
	if err != nil {
		log.Fatalf("Error creating errorCounter: %v", err)
	}

	latencyHist, err := meter.Float64Histogram(
		"poc_request_latency",
		metric.WithDescription("Histogram for request latency"),
	)
	if err != nil {
		log.Fatalf("Error creating latencyHist: %v", err)
	}

	cartItemsGauge, err := meter.Int64ObservableGauge(
		"poc_cart_items",
		metric.WithDescription("Number of items in the cart"),
		metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
			observer.Observe(metrics.CartItems.Count)
			return nil
		}))
	if err != nil {
		log.Fatalf("Error creating cartItems gauge: %v", err)
	}

	metrics.ErrorCounter = errorCounter
	metrics.LatencyHist = latencyHist
	metrics.CartItems.Gauge = cartItemsGauge
	metrics.CartItems.Count = 0
	return metrics
}

func (m *Metrics) UpdateCartItems(delta int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.CartItems.Count += delta
	if m.CartItems.Count < 0 {
		m.CartItems.Count = 0
	}
}

func (m *Metrics) RegisterLatency(ctx context.Context, duration float64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.LatencyHist.Record(ctx, duration)
}

func (m *Metrics) RegisterError(ctx context.Context, statusCode int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.ErrorCounter.Add(ctx, 1, metric.WithAttributes(semconv.HTTPStatusCodeKey.Int(statusCode)))
}

// InitMeterProvider sets up the metrics provider to send metrics to an OTLP collector with TLS
func InitMeterProvider() (*sdkmetric.MeterProvider, error) {
	ctx := context.Background()
	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if collectorURL == "" {
		collectorURL = "localhost:4317"
	}

	// Load the CA certificate
	certPool := x509.NewCertPool()
	caCertPath := os.Getenv("CA_CERT_PATH")
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("could not read CA file: %w", err)
	}
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("could not add CA certificate")
	}

	// Configure TLS options
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: certPool,
	})

	// Create the OTLP metrics exporter with TLS
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(collectorURL),
		otlpmetricgrpc.WithTLSCredentials(creds),
		otlpmetricgrpc.WithHeaders(map[string]string{"signoz-access-token": os.Getenv("SIGNOZ_ACCESS_TOKEN")}))
	if err != nil {
		return nil, fmt.Errorf("create OTLP metrics exporter: %w", err)
	}

	// Configure metrics resources
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("DEMO"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	// Create the PeriodicReader to send metrics at regular intervals
	reader := sdkmetric.NewPeriodicReader(
		metricExporter,
		sdkmetric.WithInterval(5*time.Second),
	)

	// Configure the MeterProvider
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)
	log.Println("MeterProvider configured correctly with TLS")
	return mp, nil
}
