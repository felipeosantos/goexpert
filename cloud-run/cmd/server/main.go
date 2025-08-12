package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/felipeosantos/curso-go/goexpert/cloud-run/internal/infra/webserver/clients"
	"github.com/felipeosantos/curso-go/goexpert/cloud-run/internal/infra/webserver/handlers"
	"github.com/joho/godotenv"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serviceName = semconv.ServiceNameKey.String(os.Getenv("OTEL_SERVICE_NAME"))

// Initialize a gRPC connection to be used by both the tracer and meter
// providers.
func initConn() (*grpc.ClientConn, error) {
	// It connects the OpenTelemetry Collector through local gRPC connection.
	// You may replace `localhost:4317` with your endpoint.
	conn, err := grpc.NewClient(os.Getenv("OTEL_COLLECTOR_ENDPOINT"),
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, err
}

// Initializes an OTLP exporter, and configures the corresponding trace provider.
func initTracerProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (func(context.Context) error, error) {
	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	godotenv.Load(".env")

	conn, err := initConn()
	if err != nil {
		log.Fatal(err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			serviceName,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	shutdownTracerProvider, err := initTracerProvider(ctx, res, conn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {

		log.Println("Shutting down trace...")

		shutdownTraceCtx, shutdownTraceCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownTraceCancel()

		if err := shutdownTracerProvider(shutdownTraceCtx); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %s", err)
		}

		log.Println("Trace shutdown complete")
	}()

	spanClientViaCEP := func(operation string, r *http.Request) string {
		return fmt.Sprintf("invoke-viacep %s", r.URL.Path) // nome customizado
	}

	httpClientViaCEP := clients.NewCustomClient(&http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithSpanNameFormatter(spanClientViaCEP),
		),
	})

	spanClientWeatherAPI := func(operation string, r *http.Request) string {
		return fmt.Sprintf("invoke-weatherapi %s", r.URL.Path) // nome customizado
	}

	httpClientWeatherAPI := clients.NewCustomClient(&http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithSpanNameFormatter(spanClientWeatherAPI),
		),
	})

	climaCEPHandler := handlers.NewClimaCEPHandler(httpClientViaCEP, httpClientWeatherAPI)

	// Instrumenta o handler com OpenTelemetry
	mux := http.NewServeMux()
	buscaClimaCEPHandler := otelhttp.NewHandler(http.HandlerFunc(climaCEPHandler.BuscaClimaCEP), "buscaclimacep-handler")

	// http.HandleFunc("/", climaCEPHandler.BuscaClimaCEP)
	mux.Handle("/", buscaClimaCEPHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Server running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %s", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down server...")

	// Shutdown the server gracefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("failed to shutdown server: %s", err)
	}

	log.Println("Server stopped.")
}
