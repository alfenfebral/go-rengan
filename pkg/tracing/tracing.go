package pkg_tracing

import (
	"context"
	"log"
	"os"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	trace_sdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Tracing interface {
	GetTracerProvider() *trace_sdk.TracerProvider
	ShutDown()
	LogError(span trace.Span, err error)
}

type TracingImpl struct {
	tp *trace_sdk.TracerProvider
}

func tracerProvider(url string) (*trace_sdk.TracerProvider, error) {
	appID, err := strconv.ParseInt(os.Getenv("APP_ID"), 10, 64)
	if err != nil {
		return nil, err
	}

	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	tp := trace_sdk.NewTracerProvider(
		// Always be sure to batch in production.
		trace_sdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		trace_sdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(os.Getenv("APP_NAME")),
			attribute.String("environment", os.Getenv("ENV")),
			attribute.Int64("ID", appID),
		)),
	)

	return tp, nil
}

func NewTracing() (Tracing, error) {
	tp, err := tracerProvider(os.Getenv("TRACER_PROVIDER_URL"))
	if err != nil {
		return nil, err
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return &TracingImpl{
		tp: tp,
	}, nil
}

func (t *TracingImpl) GetTracerProvider() *trace_sdk.TracerProvider {
	return t.tp
}

func (t *TracingImpl) ShutDown() {
	if err := t.tp.Shutdown(context.Background()); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
	}
}

func (t *TracingImpl) LogError(span trace.Span, err error) {
	span.SetAttributes(
		attribute.Key("error").Bool(true),
	)
	span.RecordError(err)
}
