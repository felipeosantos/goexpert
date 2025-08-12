package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/felipeosantos/curso-go/goexpert/observabilidade/internal/infra/webserver/clients"
	"github.com/felipeosantos/curso-go/goexpert/observabilidade/internal/infra/webserver/handlers"
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

	godotenv.Load(".env")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

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

	// name := "go.opentelemetry.io/contrib/examples/otel-collector"
	// tracer := otel.Tracer(name)

	// // Attributes represent additional key-value descriptors that can be bound
	// // to a metric observer or recorder.
	// commonAttrs := []attribute.KeyValue{
	// 	attribute.String("attrA", "chocolate"),
	// 	attribute.String("attrB", "raspberry"),
	// 	attribute.String("attrC", "vanilla"),
	// }

	// // Work begins
	// ctx, span := tracer.Start(
	// 	ctx,
	// 	"CollectorExporter-Example",
	// 	trace.WithAttributes(commonAttrs...))
	// defer span.End()
	// for i := range 10 {
	// 	_, iSpan := tracer.Start(ctx, fmt.Sprintf("Sample-%d", i))
	// 	log.Printf("Doing really hard work (%d / 10)\n", i+1)

	// 	<-time.After(time.Second)
	// 	iSpan.End()
	// }

	spanClientClimaCEPB := func(operation string, r *http.Request) string {
		return fmt.Sprintf("invoke-climacepb %s", r.URL.Path) // nome customizado
	}

	httpClientClimaCEPB := clients.NewCustomClient(&http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithSpanNameFormatter(spanClientClimaCEPB),
		),
	})

	// Melhorar a injeção de dependência usando um MAP, onde a chave são N parametros e o valor é o client
	climaCEPHandler := handlers.NewClimaCEPHandler(httpClientClimaCEPB)

	// Instrumenta o handler com OpenTelemetry
	mux := http.NewServeMux()
	buscaClimaCEPHandler := otelhttp.NewHandler(http.HandlerFunc(climaCEPHandler.BuscaClimaCEP), "buscaclimacep-handler")

	// http.HandleFunc("/", climaCEPHandler.BuscaClimaCEP)
	// http.Handle("/", buscaClimaCEPHandler)

	mux.Handle("/", buscaClimaCEPHandler)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	go func() {
		log.Println("Starting server on port 8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Failed to shutdown server: %s", err)
	}

	log.Println("Server gracefully stopped")
}
