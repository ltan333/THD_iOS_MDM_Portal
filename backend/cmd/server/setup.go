package main

import (
	"context"
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
	"github.com/thienel/go-backend-template/internal/usecase/worker"
	"github.com/thienel/go-backend-template/pkg/config"
	"github.com/thienel/go-backend-template/pkg/event"
	"github.com/thienel/go-backend-template/pkg/mdmcmd"
)

// setupDependencies wires up all layers and returns the configured router.
//
// Dependency flow (after decoupling):
//
//	DeviceService ──publish──► EventBus ──subscribe──► ProfileDeployWorker
//	                                      └─subscribe──► InventorySyncWorker
//
// Neither ProfileService nor InventoryLogic is injected into DeviceService
// directly, eliminating the circular dependency.
func setupDependencies(cfg *config.Config) *gin.Engine {
	// Repositories
	client := database.GetClient()
	userRepo := persistence.NewUserRepository(client)
	mobileConfigRepo := persistence.NewMobileConfigRepository(client)
	nanoRepo := persistence.NewNanoRepository(database.GetDB()) // read-only: nano server tables
	_ = nanoRepo                                                  // injected into services as needed


	// Casbin Enforcer
	enforcer, err := authorization.NewEnforcer(cfg.Casbin.ModelPath, database.GetDB())
	if err != nil {
		tlog.Fatal("Failed to initialize Casbin enforcer", zap.Error(err))
	}

	// ---- Event Bus --------------------------------------------------------
	// A single in-process event bus shared across all services and workers.
	// Workers subscribe before services are created so no events are missed.
	eventBus := event.NewBus()

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

	// DeviceService now receives the event bus instead of profileService.
	// Profile deployment and inventory sync are handled by background workers
	// that subscribe to the bus, eliminating the circular dependency.
	deviceService := serviceimpl.NewDeviceService(client, eventBus)

	applicationService := serviceimpl.NewApplicationService(client)
	alertService := serviceimpl.NewAlertService(client)
	alertRuleService := serviceimpl.NewAlertRuleService(client)
	reportService := serviceimpl.NewReportService(client)
	settingService := serviceimpl.NewSettingService(client)

	// MDM Command Builder
	cmdBuilder := mdmcmd.NewBuilder("com.thd.mdm")

	// ---- Background Workers -----------------------------------------------
	// Workers subscribe to the event bus and handle async post-enrollment work.
	// They must be started after services are created.
	workerCtx := context.Background()

	profileDeployWorker := worker.NewProfileDeployWorker(profileService, eventBus)
	profileDeployWorker.Start(workerCtx)

	inventorySyncWorker := worker.NewInventorySyncWorker(nanomdmService, eventBus)
	inventorySyncWorker.Start(workerCtx)

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
