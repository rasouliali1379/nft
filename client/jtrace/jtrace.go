package jtrace

import (
	"context"
	"io"
	"log"
	"maskan/config"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type tracing string

var (
	tracer *jtracer
	Span   tracing = "span"
)

type jtracer struct{}

func InitGlobalTracer(lc fx.Lifecycle) {
	var err error
	var closer io.Closer
	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			closer, err = tracer.connect(*config.C())
			log.Printf("Jaeger loaded successfully \n")
			return err
		},
		OnStop: func(c context.Context) error {
			log.Printf("Jaeger shutdown \n")
			return closer.Close()
		},
	})
}

func T() ITracer {
	return tracer
}

// Connect method
func (j *jtracer) connect(confs config.Config) (io.Closer, error) {
	// Initialize tracer with a logger and a metrics factory

	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via  configured Logger.
	cfg := jaegercfg.Configuration{
		ServiceName: confs.App.Name,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           confs.Jaeger.LogSpans,
			LocalAgentHostPort: confs.Jaeger.HostPort,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.ZipkinSharedRPCSpan(true),
	)

	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, err
}

// GetTracer method
func (j *jtracer) GetTracer() opentracing.Tracer {
	return opentracing.GlobalTracer()
}

// FromContext method
func (j *jtracer) FromContext(ctx context.Context, startName string) opentracing.Span {

	// if context has a span for tracing then use spanFromContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		pctx := parent.Context()
		if trc := opentracing.GlobalTracer(); trc != nil {
			spn := trc.StartSpan(startName, opentracing.ChildOf(pctx))
			return spn
		}
	}

	// if we havent span in context, create new span
	return opentracing.GlobalTracer().StartSpan(startName)
}

// StartSpan method
func (j *jtracer) StartSpan(str string) opentracing.Span {
	return opentracing.GlobalTracer().StartSpan(str)
}

// ContextWithSpan methd
func (j *jtracer) ContextWithSpan(ctx context.Context, span opentracing.Span) context.Context {
	if qr := ctx.Value(Span); qr != nil {
		ctx := context.Background()
		return opentracing.ContextWithSpan(ctx, span)
	}
	return opentracing.ContextWithSpan(ctx, span)
}

// SpanFromContext method
func (j *jtracer) SpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, name, opts...)
}

// ChildOf method
func (j *jtracer) ChildOf(span opentracing.Span, name string) opentracing.Span {
	return opentracing.StartSpan(name, opentracing.ChildOf(span.Context()))
}
