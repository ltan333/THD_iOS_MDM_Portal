package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/thienel/tlog"
	"go.uber.org/zap"

	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/response"
)

// Authorize returns Casbin authorization middleware
// Must be used AFTER Auth() middleware (requires JWT claims in context)
func (m *Middleware) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := GetRequestID(c)
		claims := GetUserClaims(c)
		if claims == nil {
			tlog.Warn("Authorization failed: missing user claims",
				zap.String("request_id", requestID),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
			)
			response.WriteErrorResponse(c, apperror.ErrUnauthorized)
			c.Abort()
			return
		}

		// Get request info
		role := claims.Role
		path := c.FullPath() // e.g. "/api/users/:id"
		method := c.Request.Method

		// If FullPath is empty (no matching route), use the raw URL path
		if path == "" {
			path = c.Request.URL.Path
		}

		// Enforce policy via Casbin
		allowed, err := m.authzService.Enforce(role, path, method)
		if err != nil {
			tlog.Error("Authorization failed: casbin enforce error",
				zap.String("request_id", requestID),
				zap.String("role", role),
				zap.String("method", method),
				zap.String("path", path),
				zap.Error(err),
			)
			response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithMessage("Lỗi kiểm tra quyền truy cập").WithError(err))
			c.Abort()
			return
		}

		if !allowed {
			tlog.Warn("Authorization denied",
				zap.String("request_id", requestID),
				zap.Uint("user_id", claims.UserID),
				zap.String("username", claims.Username),
				zap.String("role", role),
				zap.String("method", method),
				zap.String("path", path),
			)
			response.WriteErrorResponse(c, apperror.ErrForbidden.WithMessage("Bạn không có quyền thực hiện hành động này"))
			c.Abort()
			return
		}

		tlog.Debug("Authorization granted",
			zap.String("request_id", requestID),
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
			zap.String("role", role),
			zap.String("method", method),
			zap.String("path", path),
		)

		c.Next()
	}
}
