package tracing

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	trace_sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Tracing interface {
	GetTracerProvider() *trace_sdk.TracerProvider
	ShutDown()
	LogError(span trace.Span, err error)
	Tracer(name string) trace.Tracer
}

type TracingImpl struct {
	tp *trace_sdk.TracerProvider
}

func tracerProvider(url string) (*trace_sdk.TracerProvider, error) {
	appID, err := strconv.ParseInt(os.Getenv("APP_ID"), 10, 64)
	if err != nil {
		return nil, err
	}

	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(url),
		uptrace.WithServiceName(os.Getenv("APP_NAME")),
		uptrace.WithServiceVersion("v0.0.1"),
		uptrace.WithDeploymentEnvironment(os.Getenv("ENV")),
		uptrace.WithResourceAttributes(
			attribute.Int64("ID", appID),
		),
	)

	tp := uptrace.TracerProvider()

	return tp, nil
}

func New() (Tracing, error) {
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

func (t *TracingImpl) Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}
