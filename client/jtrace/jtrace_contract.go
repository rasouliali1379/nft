package jtrace

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

type ITracer interface {
	GetTracer() opentracing.Tracer
	FromContext(ctx context.Context, startName string) opentracing.Span
	StartSpan(str string) opentracing.Span
	ContextWithSpan(ctx context.Context, span opentracing.Span) context.Context
	SpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context)
}
