package logger

import (
	"context"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor logs one structured entry for each unary gRPC call.
// It records method, status code, latency, peer address, and optional request id.
func UnaryServerInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// Measure end-to-end request duration including handler execution.
		start := time.Now()

		resp, err := handler(ctx, req)

		// Convert handler error to canonical gRPC status code.
		code := status.Code(err)
		if err == nil {
			code = codes.OK
		}

		// Build common structured attributes emitted for both success and failure.
		attrs := []any{
			slog.String("component", "grpc"),
			slog.String("method", info.FullMethod),
			slog.String("code", code.String()),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
		}

		// Include remote peer address when available.
		if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
			attrs = append(attrs, slog.String("peer", p.Addr.String()))
		}

		// Propagate request correlation id from incoming metadata when provided.
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get("x-request-id"); len(values) > 0 {
				attrs = append(attrs, slog.String("request_id", values[0]))
			}
		}

		if err != nil {
			// Error path: log at error level and return original handler result.
			attrs = append(attrs, slog.String("error", err.Error()))
			log.Error("grpc request failed", attrs...)
			return resp, err
		}

		if count, responseErrors := extractResponseErrors(resp); count > 0 {
			attrs = append(attrs,
				slog.Int("response_error_count", count),
				slog.Any("response_errors", responseErrors),
				slog.Any("response_data", resp),
			)
			log.Error("grpc request completed with response errors", attrs...)
			return resp, nil
		}

		// Success path: log at info level.
		log.Info("grpc request completed: ", attrs...)
		return resp, nil
	}
}

// extractResponseErrors reads the protobuf "Errors" field (if present)
// and returns a fully structured payload for logging/tracing.
func extractResponseErrors(resp any) (int, []map[string]any) {
	if resp == nil {
		return 0, nil
	}

	v := reflect.ValueOf(resp)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return 0, nil
		}
		v = v.Elem()
	}

	if !v.IsValid() || v.Kind() != reflect.Struct {
		return 0, nil
	}

	errorsField := v.FieldByName("Errors")
	if !errorsField.IsValid() || errorsField.Kind() != reflect.Slice {
		return 0, nil
	}

	count := errorsField.Len()
	if count == 0 {
		return 0, nil
	}

	details := make([]map[string]any, 0, count)
	for i := 0; i < count; i++ {
		item := errorsField.Index(i)
		if item.Kind() == reflect.Ptr {
			if item.IsNil() {
				continue
			}
			item = item.Elem()
		}

		if item.Kind() != reflect.Struct {
			continue
		}

		detail := make(map[string]any)
		itemType := item.Type()
		for fieldIndex := 0; fieldIndex < item.NumField(); fieldIndex++ {
			fieldMeta := itemType.Field(fieldIndex)
			if fieldMeta.PkgPath != "" {
				continue
			}
			detail[strings.ToLower(fieldMeta.Name)] = item.Field(fieldIndex).Interface()
		}

		if len(detail) > 0 {
			details = append(details, detail)
		}
	}

	return count, details
}
