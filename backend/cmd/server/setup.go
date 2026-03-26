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
	"github.com/thienel/go-backend-template/pkg/mdmcmd"
)

// setupDependencies wires up all layers and returns the configured router
func setupDependencies(cfg *config.Config) *gin.Engine {
	// Repositories
	client := database.GetClient()
	userRepo := persistence.NewUserRepository(client)
	mobileConfigRepo := persistence.NewMobileConfigRepository(client)

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
	redisService := serviceimpl.NewRedisService(cfg.Redis)
	authService := serviceimpl.NewAuthService(userRepo, jwtService, authzService, redisService)
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
	mobileConfigService := serviceimpl.NewMobileConfigService(mobileConfigRepo)
	dashboardService := serviceimpl.NewDashboardService(client)
	deviceGroupService := serviceimpl.NewDeviceGroupService(client)
	profileGenerator := serviceimpl.NewProfileGenerator("THD MDM", "com.thd.mdm")
	profileService := serviceimpl.NewProfileService(client, profileGenerator, nanomdmService)
	deviceService := serviceimpl.NewDeviceService(client, profileService)
	applicationService := serviceimpl.NewApplicationService(client)
	alertService := serviceimpl.NewAlertService(client)
	alertRuleService := serviceimpl.NewAlertRuleService(client)
	reportService := serviceimpl.NewReportService(client)
	settingService := serviceimpl.NewSettingService(client)

	// MDM Command Builder
	cmdBuilder := mdmcmd.NewBuilder("com.thd.mdm")

	// Middleware
	origins := strings.Join(cfg.CORSAllowedOrigins, ",")
	mw := middleware.New(jwtService, authzService, redisService, origins)

	// Handlers
	authHandler := handler.NewAuthHandler(authService, userService, authzService)
	userHandler := handler.NewUserHandler(userService, authzService)
	policyHandler := handler.NewPolicyHandler(authzService)
	nanocmdHandler := handler.NewNanoCMDHandler(nanocmdService, deviceService)
	mobileConfigHandler := handler.NewMobileConfigHandler(mobileConfigService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	deviceHandler := handler.NewDeviceHandler(deviceService, nanomdmService, profileService, cmdBuilder)
	deviceGroupHandler := handler.NewDeviceGroupHandler(deviceGroupService)
	profileHandler := handler.NewProfileHandler(profileService)
	applicationHandler := handler.NewApplicationHandler(applicationService)
	alertHandler := handler.NewAlertHandler(alertService, alertRuleService)
	reportHandler := handler.NewReportHandler(reportService)
	settingHandler := handler.NewSettingHandler(settingService)

	// Build router
	return router.SetupRouter(authHandler, userHandler, policyHandler, nanocmdHandler, mobileConfigHandler, dashboardHandler, deviceHandler, deviceGroupHandler, profileHandler, applicationHandler, alertHandler, reportHandler, settingHandler, mw)
}
