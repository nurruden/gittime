package opentelemetry

import (
	"gitime/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const instrumentstionName = "gitime/web/middlewares/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func (m MiddlewareBuilder) Build() web.Middleware {
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentstionName)
	}
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			reqCtx := ctx.Req.Context()
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))
			reqCtx, span := m.Tracer.Start(reqCtx, "unknown")
			defer span.End()
			//defer func() {
			//	span.SetName(ctx.MatchedRoute)
			//	span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))
			//	span.End()
			//}()
			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.schema", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("http.host", ctx.Req.Host))
			ctx.Req = ctx.Req.WithContext(reqCtx)
			//ctx.Ctx = reqCtx
			next(ctx)
			span.SetName(ctx.MatchedRoute)
			span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))

		}
	}
}
