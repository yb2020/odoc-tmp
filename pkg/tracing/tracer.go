package tracing

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/yb2020/odoc/pkg/logging"
)

// TracerProvider 定义了一个通用的跟踪提供者接口
// 这个接口允许我们在将来无缝切换到 OpenTelemetry
type TracerProvider interface {
	// GetTracer 返回一个跟踪器实例
	GetTracer() opentracing.Tracer
	// Close 关闭跟踪提供者及其资源
	Close() error
}

// JaegerProvider 实现了 TracerProvider 接口，使用 Jaeger 作为后端
type JaegerProvider struct {
	tracer opentracing.Tracer
	closer io.Closer
}

// GetTracer 返回 Jaeger 跟踪器
func (jp *JaegerProvider) GetTracer() opentracing.Tracer {
	return jp.tracer
}

// Close 关闭 Jaeger 跟踪器
func (jp *JaegerProvider) Close() error {
	if jp.closer != nil {
		return jp.closer.Close()
	}
	return nil
}

// OTelProvider 是 OpenTelemetry 提供者的占位符
// 这个结构体将在未来实现，当我们集成 OpenTelemetry 时
type OTelProvider struct {
	// 将在未来实现
}

// GetTracer 返回 OpenTelemetry 跟踪器的桥接
// 目前返回 NoopTracer，将在未来实现
func (op *OTelProvider) GetTracer() opentracing.Tracer {
	// 将在未来实现 OpenTelemetry 到 OpenTracing 的桥接
	return opentracing.NoopTracer{}
}

// Close 关闭 OpenTelemetry 提供者
func (op *OTelProvider) Close() error {
	// 将在未来实现
	return nil
}

// InitTracer initializes a new Jaeger tracer
// 这个函数保留用于向后兼容
func InitTracer(serviceName string, jaegerURL string, sampleRate float64, logger logging.Logger) (opentracing.Tracer, io.Closer, error) {
	// Configure Jaeger tracing
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: sampleRate,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: jaegerURL,
		},
	}

	// Create Jaeger tracer
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	// Set the global tracer
	opentracing.SetGlobalTracer(tracer)

	logger.Info("msg", "Jaeger tracer initialized", "service", serviceName, "jaegerURL", jaegerURL)

	return tracer, closer, nil
}

// NewJaegerProvider 创建一个新的 Jaeger 跟踪提供者
func NewJaegerProvider(serviceName string, jaegerURL string, sampleRate float64, logger logging.Logger) (TracerProvider, error) {
	tracer, closer, err := InitTracer(serviceName, jaegerURL, sampleRate, logger)
	if err != nil {
		return nil, err
	}

	return &JaegerProvider{
		tracer: tracer,
		closer: closer,
	}, nil
}

// NewTracerProvider 根据配置创建适当的跟踪提供者
// 目前只支持 Jaeger，但在未来会支持 OpenTelemetry
func NewTracerProvider(providerType string, serviceName string, endpoint string, sampleRate float64, logger logging.Logger) (TracerProvider, error) {
	switch providerType {
	case "jaeger":
		return NewJaegerProvider(serviceName, endpoint, sampleRate, logger)
	case "opentelemetry":
		// 将在未来实现
		logger.Warn("msg", "OpenTelemetry 尚未实现，使用 NoopTracer 代替")
		return &OTelProvider{}, nil
	default:
		logger.Warn("msg", "未知的跟踪提供者类型，使用 Jaeger", "type", providerType)
		return NewJaegerProvider(serviceName, endpoint, sampleRate, logger)
	}
}
