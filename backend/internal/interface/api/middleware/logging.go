package middleware

import (
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thienel/tlog"
	"go.uber.org/zap"

	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/response"
)

const requestIDHeader = "X-Request-ID"

// RequestContext ensures each request has a stable request_id for correlation.
func (m *Middleware) RequestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set(requestIDHeader, requestID)
		c.Writer.Header().Set(requestIDHeader, requestID)
		c.Next()
	}
}

// AccessLog writes a structured completion log for every request.
func (m *Middleware) AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		requestID := GetRequestID(c)
		latency := time.Since(start)

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("route", c.FullPath()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.Int("response_size", c.Writer.Size()),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		if claims := GetUserClaims(c); claims != nil {
			fields = append(fields,
				zap.Uint("user_id", claims.UserID),
				zap.String("username", claims.Username),
				zap.String("role", claims.Role),
			)
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		status := c.Writer.Status()
		switch {
		case status >= 500:
			tlog.Error("Request completed with server error", fields...)
		case status >= 400:
			tlog.Warn("Request completed with client error", fields...)
		default:
			tlog.Info("Request completed", fields...)
		}
	}
}

// Recovery catches panics and emits a structured log with request context.
func (m *Middleware) Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		requestID := GetRequestID(c)
		tlog.Error("Panic recovered",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("route", c.FullPath()),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Any("panic", recovered),
			zap.ByteString("stack", debug.Stack()),
		)

		response.WriteErrorResponse(c, apperror.ErrInternalServerError)
		c.Abort()
	})
}

// GetRequestID returns request correlation ID from context/header.
func GetRequestID(c *gin.Context) string {
	if v, ok := c.Get(requestIDHeader); ok {
		if requestID, ok := v.(string); ok {
			return requestID
		}
	}
	return c.GetHeader(requestIDHeader)
}
