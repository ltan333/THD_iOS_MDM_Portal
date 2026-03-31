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
	_ = nanoRepo                                                // injected into services as needed
	payloadPropertyDefinitionRepo := persistence.NewPayloadPropertyDefinitionRepository(client)

	dashboardRepo := persistence.NewDashboardRepository(client)
	deviceGroupRepo := persistence.NewDeviceGroupRepository(client)
	alertRepo := persistence.NewAlertRepository(client)
	alertRuleRepo := persistence.NewAlertRuleRepository(client)
	reportRepo := persistence.NewReportRepository(client)
	settingRepo := persistence.NewSettingRepository(client)
	appRepo := persistence.NewApplicationRepository(client)
	deviceRepo := persistence.NewDeviceRepository(client)
	profileRepo := persistence.NewProfileRepository(client)
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
		cfg.NanoMDM.DEPSyncerContainer,
	)
	mobileConfigService := serviceimpl.NewMobileConfigService(mobileConfigRepo)
	dashboardService := serviceimpl.NewDashboardService(dashboardRepo, alertRepo, appRepo)
	deviceGroupService := serviceimpl.NewDeviceGroupService(deviceGroupRepo)
	profileGenerator := serviceimpl.NewProfileGenerator("THD MDM", "com.thd.mdm")
	profileService := serviceimpl.NewProfileService(profileRepo, profileGenerator, nanomdmService)

	// DeviceService now receives the event bus instead of profileService.
	// Profile deployment and inventory sync are handled by background workers
	// that subscribe to the bus, eliminating the circular dependency.
	deviceService := serviceimpl.NewDeviceService(deviceRepo, eventBus)

	applicationService := serviceimpl.NewApplicationService(appRepo, nanomdmService)
	alertService := serviceimpl.NewAlertService(alertRepo, alertRuleRepo, nanomdmService)
	alertRuleService := serviceimpl.NewAlertRuleService(alertRuleRepo)
	reportService := serviceimpl.NewReportService(reportRepo)
	settingService := serviceimpl.NewSettingService(settingRepo)

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

	payloadPropertyDefinitionService := serviceimpl.NewPayloadPropertyDefinitionService(payloadPropertyDefinitionRepo)

	// Middleware
	origins := strings.Join(cfg.CORSAllowedOrigins, ",")
	mw := middleware.New(jwtService, authzService, redisService, origins)

	// Handlers
	authHandler := handler.NewAuthHandler(authService, userService, authzService)
	userHandler := handler.NewUserHandler(userService, authzService)
	policyHandler := handler.NewPolicyHandler(authzService)
	nanocmdHandler := handler.NewNanoCMDHandler(nanocmdService, deviceService, nanomdmService, cfg)
	mobileConfigHandler := handler.NewMobileConfigHandler(mobileConfigService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	deviceHandler := handler.NewDeviceHandler(deviceService, nanomdmService, profileService, cmdBuilder)
	deviceGroupHandler := handler.NewDeviceGroupHandler(deviceGroupService)
	profileHandler := handler.NewProfileHandler(profileService)
	applicationHandler := handler.NewApplicationHandler(applicationService)
	alertHandler := handler.NewAlertHandler(alertService, alertRuleService)
	reportHandler := handler.NewReportHandler(reportService)
	settingHandler := handler.NewSettingHandler(settingService, nanomdmService)
	payloadPropertyDefinitionHandler := handler.NewPayloadPropertyDefinitionHandler(payloadPropertyDefinitionService)

	// Build router
	return router.SetupRouter(authHandler, userHandler, policyHandler, nanocmdHandler, mobileConfigHandler, dashboardHandler, deviceHandler, deviceGroupHandler, profileHandler, applicationHandler, alertHandler, reportHandler, settingHandler, payloadPropertyDefinitionHandler, mw)
}
