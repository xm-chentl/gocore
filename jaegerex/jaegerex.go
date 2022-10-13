package jaegerex

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

type ILinkTracing interface {
	Close() error
	WithStartSpan(string) ITracingSpan
	WithSpanFromContext(context.Context, string) ITracingSpan
}

type ITracingSpan interface {
	Context() context.Context
	SetTag(key string, value interface{})
	Finish()
}

type linkTracingImpl struct {
	serviceName string
	reportAddr  string // host:port

	tracer opentracing.Tracer
	closer io.Closer
}

func (t linkTracingImpl) Close() error {
	if t.closer != nil {
		return t.closer.Close()
	}

	return nil
}

func (t linkTracingImpl) WithStartSpan(name string) ITracingSpan {
	span := t.tracer.StartSpan(name)
	return &tracingSpanImpl{
		ctx:  opentracing.ContextWithSpan(context.Background(), span),
		span: span,
	}
}

func (t linkTracingImpl) WithSpanFromContext(ctx context.Context, name string) ITracingSpan {
	span, ctx := opentracing.StartSpanFromContext(ctx, name)
	return &tracingSpanImpl{
		ctx:  ctx,
		span: span,
	}
}

type tracingSpanImpl struct {
	ctx  context.Context
	span opentracing.Span
}

func (s tracingSpanImpl) Context() context.Context {
	return s.ctx
}

func (s tracingSpanImpl) SetTag(key string, value interface{}) {
	if s.span != nil {
		s.span.SetTag(key, value)
	}
}

func (s tracingSpanImpl) Finish() {
	if s.span != nil {
		s.span.Finish()
	}
}

func New(serviceName, reportAddr string) ILinkTracing {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const", // 常量采样
			Param: 1,       // 开启所有轨迹采样
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: reportAddr,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic("new jaeger configuration failed: " + err.Error())
	}
	opentracing.SetGlobalTracer(tracer)

	return &linkTracingImpl{
		serviceName: cfg.ServiceName,
		reportAddr:  reportAddr,
		tracer:      tracer,
		closer:      closer,
	}
}
