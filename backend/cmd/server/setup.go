package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thienel/tlog"
	"go.uber.org/zap"

	"github.com/thienel/go-backend-template/internal/infra/authorization"
	"github.com/thienel/go-backend-template/internal/infra/database"
	"github.com/thienel/go-backend-template/internal/infra/persistence"
	"github.com/thienel/go-backend-template/internal/interface/api/handler"
	"github.com/thienel/go-backend-template/internal/interface/api/middleware"
	"github.com/thienel/go-backend-template/internal/interface/api/router"
	"github.com/thienel/go-backend-template/internal/usecase/service/serviceimpl"
	"github.com/thienel/go-backend-template/pkg/config"
)

// setupDependencies wires up all layers and returns the configured router
func setupDependencies(cfg *config.Config) *gin.Engine {
	// Repositories
	client := database.GetClient()
	userRepo := persistence.NewUserRepository(client)
	depProfileRepo := persistence.NewDepProfileRepository(client)
	mobileConfigRepo := persistence.NewMobileConfigRepository(client)
	payloadPropertyDefinitionRepo := persistence.NewPayloadPropertyDefinitionRepository(client)

	// Casbin Enforcer
	enforcer, err := authorization.NewEnforcer(cfg.Casbin.ModelPath, database.GetDB())
	if err != nil {
		tlog.Fatal("Failed to initialize Casbin enforcer", zap.Error(err))
	}

	// Services
	jwtService := serviceimpl.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpiryMinutes,
		cfg.JWT.RefreshExpiryHours,
	)
	authzService := serviceimpl.NewAuthorizationService(enforcer)
	authService := serviceimpl.NewAuthService(userRepo, jwtService, authzService)
	userService := serviceimpl.NewUserService(userRepo, authzService)
	nanocmdService := serviceimpl.NewNanoCMDService(cfg.NanoCMD.BaseURL, cfg.NanoCMD.Username, cfg.NanoCMD.Password)
	nanomdmService := serviceimpl.NewNanoMDMService(
		cfg.NanoMDM.MDMBaseURL,
		cfg.NanoMDM.DEPBaseURL,
		cfg.NanoMDM.MDMUsername,
		cfg.NanoMDM.MDMPassword,
		cfg.NanoMDM.DEPUsername,
		cfg.NanoMDM.DEPPassword,
	)
	depProfileService := serviceimpl.NewDepProfileService(depProfileRepo, nanomdmService)
	mobileConfigService := serviceimpl.NewMobileConfigService(mobileConfigRepo)
	payloadPropertyDefinitionService := serviceimpl.NewPayloadPropertyDefinitionService(payloadPropertyDefinitionRepo)

	// Middleware
	origins := strings.Join(cfg.CORSAllowedOrigins, ",")
	mw := middleware.New(jwtService, authzService, origins)

	// Handlers
	authHandler := handler.NewAuthHandler(authService, userService, authzService)
	userHandler := handler.NewUserHandler(userService, authzService)
	policyHandler := handler.NewPolicyHandler(authzService)
	mdmHandler := handler.NewMDMHandler(client, nanomdmService)
	depHandler := handler.NewDEPHandler(client, authzService, nanomdmService, depProfileService)
	nanocmdHandler := handler.NewNanoCMDHandler(nanocmdService)
	mobileConfigHandler := handler.NewMobileConfigHandler(mobileConfigService)
	payloadPropertyDefinitionHandler := handler.NewPayloadPropertyDefinitionHandler(payloadPropertyDefinitionService)

	// Build router
	return router.SetupRouter(authHandler, userHandler, policyHandler, mdmHandler, depHandler, nanocmdHandler, mobileConfigHandler, payloadPropertyDefinitionHandler, mw)
}
