package middleware

import (
	"strings"

	"github.com/thienel/go-backend-template/internal/usecase/service"
)

// Middleware holds all middleware dependencies
type Middleware struct {
	jwtService     service.JWTService
	authzService   service.AuthorizationService
	redisService   service.RedisService
	origins        string
	allowedOrigins []string
	allowAll       bool
}

// New creates a new Middleware instance
func New(jwtService service.JWTService, authzService service.AuthorizationService, redisService service.RedisService, origins string) *Middleware {
	allowed := strings.Split(origins, ",")
	allowAll := len(allowed) == 1 && strings.TrimSpace(allowed[0]) == "*"

	for i := range allowed {
		allowed[i] = strings.TrimSpace(allowed[i])
	}

	return &Middleware{
		jwtService:     jwtService,
		authzService:   authzService,
		redisService:   redisService,
		origins:        origins,
		allowedOrigins: allowed,
		allowAll:       allowAll,
	}
}
