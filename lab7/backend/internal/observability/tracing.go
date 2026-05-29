package observability

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const defaultServiceName = "sport-platform-backend"

func ServiceName() string {
	serviceName := strings.TrimSpace(os.Getenv("OTEL_SERVICE_NAME"))
	if serviceName == "" {
		return defaultServiceName
	}
	return serviceName
}

func GinTracingMiddleware() gin.HandlerFunc {
	return otelgin.Middleware(ServiceName())
}

func ConfigureTracing(ctx context.Context) (func(context.Context) error, error) {
	endpoint := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if endpoint == "" {
		return func(context.Context) error { return nil }, nil
	}

	exporterOptions, err := otlpHTTPOptions(endpoint)
	if err != nil {
		return func(context.Context) error { return nil }, err
	}

	exporter, err := otlptracehttp.New(ctx, exporterOptions...)
	if err != nil {
		return func(context.Context) error { return nil }, err
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", ServiceName()),
			attribute.String("service.namespace", "sport-platform"),
			attribute.String("deployment.environment", deploymentEnvironment()),
		),
	)
	if err != nil {
		return func(context.Context) error { return nil }, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tracerProvider.Shutdown, nil
}

func otlpHTTPOptions(endpoint string) ([]otlptracehttp.Option, error) {
	options := []otlptracehttp.Option{
		otlptracehttp.WithTimeout(5 * time.Second),
	}

	parsedURL, err := url.Parse(endpoint)
	if err == nil && parsedURL.Scheme != "" {
		if parsedURL.Scheme == "http" {
			options = append(options, otlptracehttp.WithInsecure())
		}

		if parsedURL.Path != "" && parsedURL.Path != "/" {
			options = append(options, otlptracehttp.WithURLPath(parsedURL.Path))
		}

		return append(options, otlptracehttp.WithEndpoint(parsedURL.Host)), nil
	}

	if strings.Contains(endpoint, "://") {
		return nil, fmt.Errorf("invalid OTEL_EXPORTER_OTLP_ENDPOINT %q", endpoint)
	}

	return append(options, otlptracehttp.WithEndpoint(endpoint), otlptracehttp.WithInsecure()), nil
}

func deploymentEnvironment() string {
	value := strings.TrimSpace(os.Getenv("DEPLOYMENT_ENVIRONMENT"))
	if value == "" {
		return "local"
	}
	return value
}
