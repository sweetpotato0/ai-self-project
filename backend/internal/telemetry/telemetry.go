package telemetry

import (
	"context"
	"fmt"

	"gin-web-framework/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// InitTelemetry 初始化OpenTelemetry
func InitTelemetry(serviceName, serviceVersion string) error {
	// 获取配置
	cfg := config.Get()
	telemetryCfg := cfg.GetTelemetry()

	// 检查是否启用遥测
	if !telemetryCfg.Enabled {
		return nil
	}

	// 创建资源
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(telemetryCfg.ServiceName),
			semconv.ServiceVersion(telemetryCfg.ServiceVersion),
			semconv.DeploymentEnvironment(telemetryCfg.Environment),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %v", err)
	}

	// 初始化追踪器
	if telemetryCfg.EnableTracing {
		if err := initTracer(res, telemetryCfg); err != nil {
			return fmt.Errorf("failed to init tracer: %v", err)
		}
	}

	// 初始化指标收集器
	if telemetryCfg.EnableMetrics {
		if err := initMeter(res); err != nil {
			return fmt.Errorf("failed to init meter: %v", err)
		}
	}

	return nil
}

// initTracer 初始化追踪器
func initTracer(res *resource.Resource, telemetryCfg config.TelemetryConfig) error {
	// 创建OTLP HTTP导出器
	exp, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(telemetryCfg.OTLPEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("failed to create OTLP exporter: %v", err)
	}

	// 创建追踪器提供者
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	// 设置全局追踪器提供者
	otel.SetTracerProvider(tp)

	// 设置全局传播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return nil
}

// initMeter 初始化指标收集器
func initMeter(res *resource.Resource) error {
	// 创建Prometheus导出器
	exp, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("failed to create prometheus exporter: %v", err)
	}

	// 创建指标提供者
	mp := metric.NewMeterProvider(
		metric.WithReader(exp),
		metric.WithResource(res),
	)

	// 设置全局指标提供者
	otel.SetMeterProvider(mp)

	return nil
}

// ShutdownTelemetry 关闭遥测系统
func ShutdownTelemetry(ctx context.Context) error {
	// 关闭追踪器提供者
	if tp := otel.GetTracerProvider(); tp != nil {
		if err := tp.(*trace.TracerProvider).Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown tracer provider: %v", err)
		}
	}

	// 关闭指标提供者
	if mp := otel.GetMeterProvider(); mp != nil {
		if err := mp.(*metric.MeterProvider).Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown meter provider: %v", err)
		}
	}

	return nil
}
