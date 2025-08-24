package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("ai-self-project-backend")

// Tracing 分布式追踪中间件
func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 从请求头中提取追踪信息
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request.Header))

		// 创建span
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
				attribute.String("http.user_agent", c.Request.UserAgent()),
				attribute.String("http.request_id", c.GetString("request_id")),
			),
		)
		defer span.End()

		// 将span context设置到gin context中
		c.Request = c.Request.WithContext(ctx)

		// 处理请求
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		// 记录响应信息
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
			attribute.Int("http.response_size", c.Writer.Size()),
			attribute.String("http.duration", duration.String()),
		)

		// 设置span状态
		if c.Writer.Status() >= 400 {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", c.Writer.Status()))
		} else {
			span.SetStatus(codes.Ok, "")
		}

		// 记录错误信息
		if len(c.Errors) > 0 {
			span.RecordError(c.Errors.Last().Err)
		}
	}
}

// DatabaseTracing 数据库操作追踪
func DatabaseTracing(ctx context.Context, operation string, query string) (context.Context, trace.Span) {
	return tracer.Start(ctx, fmt.Sprintf("db.%s", operation),
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.operation", operation),
			attribute.String("db.query", query),
		),
	)
}

// RedisTracing Redis操作追踪
func RedisTracing(ctx context.Context, operation string, key string) (context.Context, trace.Span) {
	return tracer.Start(ctx, fmt.Sprintf("redis.%s", operation),
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("redis.operation", operation),
			attribute.String("redis.key", key),
		),
	)
}

// BusinessTracing 业务操作追踪
func BusinessTracing(ctx context.Context, operation string, params map[string]interface{}) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		attribute.String("business.operation", operation),
	}

	for k, v := range params {
		attrs = append(attrs, attribute.String(fmt.Sprintf("business.%s", k), fmt.Sprintf("%v", v)))
	}

	return tracer.Start(ctx, fmt.Sprintf("business.%s", operation),
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(attrs...),
	)
}

// GetTraceID 获取追踪ID
func GetTraceID(ctx context.Context) string {
	if span := trace.SpanFromContext(ctx); span != nil {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

// GetSpanID 获取Span ID
func GetSpanID(ctx context.Context) string {
	if span := trace.SpanFromContext(ctx); span != nil {
		return span.SpanContext().SpanID().String()
	}
	return ""
}
