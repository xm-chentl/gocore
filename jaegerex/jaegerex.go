package jaegerex

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	openTracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

type ILinkTracing interface {
	Close()
	// WithHeaderSpan(http.Header) ITracingSpan
	WithStartSpan(string) ITracingSpan
	WithSpanFromContext(context.Context, string) ITracingSpan
}

type ITracingSpan interface {
	Context() context.Context
	Finish()
	SetTag(key string, value interface{})
	SetLogs(fields ...openTracingLog.Field)
	Span() opentracing.Span
}

type ITracingSpanExt interface {
}

type linkTracingImpl struct {
	serviceName string
	reportAddr  string // host:port

	tracer opentracing.Tracer
	closer io.Closer // todo: 链接使用完未关闭，最后一个次的链接会有丢失。 优化：在程序退出时调用close
}

func (t linkTracingImpl) Close() {
	if t.closer != nil {
		t.closer.Close()
	}
}

func (t linkTracingImpl) WithHeaderSpan(header http.Header, name string) (ts ITracingSpan, err error) {
	var span opentracing.Span
	spanCtx, err := t.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
	if err != nil {
		span = t.tracer.StartSpan(name)
	} else {
		span = t.tracer.StartSpan(
			name,
			opentracing.ChildOf(spanCtx),
			ext.SpanKindRPCClient,
		)
	}

	ts = &tracingSpanImpl{
		span: span,
	}

	return
}

func (t linkTracingImpl) Inject(header http.Header, tracingSpan ITracingSpan) {
	err := t.tracer.Inject(tracingSpan.Span().Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
	if err != nil {
		tracingSpan.SetLogs(openTracingLog.String("inject_error", err.Error()))
	}

	// todo: 预留添加traceID 放至ctx
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
	if s.ctx == nil {
		s.ctx = context.Background()
	}

	return s.ctx
}

func (s tracingSpanImpl) Span() opentracing.Span {
	return s.span
}

func (s tracingSpanImpl) SetTag(key string, value interface{}) {
	s.span.SetTag(key, value)
}

func (s tracingSpanImpl) SetLogs(fields ...openTracingLog.Field) {
	if len(fields) > 0 {
		s.span.LogFields(fields...)
	}
}

func (s tracingSpanImpl) Finish() {
	if s.span != nil {
		s.span.Finish()
	}
}

// 暂时方案不存在待久化
func New(serviceName, reportAddr string) ILinkTracing {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst, // 常量采样
			Param: 1,                       // 开启所有轨迹采样
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: reportAddr,
			LogSpans:           true,
		},
		// Collector 配置
		// Reporter: &config.ReporterConfig{
		// 	BufferFlushInterval: 100 * time.Millisecond,
		// 	CollectorEndpoint:   reportAddr,
		// 	LogSpans:            false,
		// },
	}
	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
		// config.ZipkinSharedRPCSpan(true), 需要单独开启
	)
	if err != nil {
		log.Println("new jaeger configuration failed: ", err)
		os.Exit(1)
		return nil
	}
	opentracing.SetGlobalTracer(tracer)

	return &linkTracingImpl{
		serviceName: cfg.ServiceName,
		reportAddr:  reportAddr,
		tracer:      tracer,
		closer:      closer,
	}
}
