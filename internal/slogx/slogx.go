package slogx

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"

	"github.com/otakakot/sample-gorm-slog/internal/contextx"
)

var keys = []contextx.Key{
	contextx.UserIDKey,
}

type LogHandler struct {
	slog.Handler
}

func New(
	handler slog.Handler,
) *LogHandler {
	return &LogHandler{
		Handler: handler,
	}
}

func (hdl *LogHandler) Handle(
	ctx context.Context,
	record slog.Record,
) error {
	span := trace.SpanFromContext(ctx)

	record.AddAttrs(slog.Attr{Key: "trace", Value: slog.StringValue(span.SpanContext().TraceID().String())})

	record.AddAttrs(slog.Attr{Key: "span", Value: slog.StringValue(span.SpanContext().SpanID().String())})

	for _, key := range keys {
		if val := ctx.Value(key); val != nil {
			record.AddAttrs(slog.Attr{Key: key.String(), Value: slog.AnyValue(val)})
		}
	}

	return hdl.Handler.Handle(ctx, record)
}
