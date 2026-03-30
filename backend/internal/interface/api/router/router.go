package router

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thienel/tlog"

	"github.com/thienel/go-backend-template/internal/interface/api/handler"
	"github.com/thienel/go-backend-template/internal/interface/api/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/thienel/go-backend-template/docs"
)

type routeRegister struct {
	auth         handler.AuthHandler
	user         handler.UserHandler
	policy       handler.PolicyHandler
	nanocmd      handler.NanoCMDHandler // Keep for webhook only
	mobileConfig handler.MobileConfigHandler
	ppd          handler.PayloadPropertyDefinitionHandler
	dashboard    handler.DashboardHandler
	device       handler.DeviceHandler
	deviceGroup  handler.DeviceGroupHandler
	profile      handler.ProfileHandler
	application  handler.ApplicationHandler
	alert        handler.AlertHandler
	report       handler.ReportHandler
	setting      handler.SettingHandler
	mw           *middleware.Middleware
}

// SetupRouter configures all routes following THD-Checkin-App pattern
func SetupRouter(
	authHandler handler.AuthHandler,
	userHandler handler.UserHandler,
	policyHandler handler.PolicyHandler,
	nanocmdHandler handler.NanoCMDHandler, // Keep for webhook only
	mobileConfigHandler handler.MobileConfigHandler,
	dashboardHandler handler.DashboardHandler,
	deviceHandler handler.DeviceHandler,
	deviceGroupHandler handler.DeviceGroupHandler,
	profileHandler handler.ProfileHandler,
	applicationHandler handler.ApplicationHandler,
	alertHandler handler.AlertHandler,
	reportHandler handler.ReportHandler,
	settingHandler handler.SettingHandler,
	payloadPropertyDefinitionHandler handler.PayloadPropertyDefinitionHandler,
	mw *middleware.Middleware,
) *gin.Engine {

	routes := routeRegister{
		auth:         authHandler,
		user:         userHandler,
		policy:       policyHandler,
		nanocmd:      nanocmdHandler,
		mobileConfig: mobileConfigHandler,
		ppd:          payloadPropertyDefinitionHandler,
		dashboard:    dashboardHandler,
		device:       deviceHandler,
		deviceGroup:  deviceGroupHandler,
		profile:      profileHandler,
		application:  applicationHandler,
		alert:        alertHandler,
		report:       reportHandler,
		setting:      settingHandler,
		mw:           mw,
	}

	router := gin.New()
	router.Use(gin.Recovery(), mw.CORS(), tlog.GinMiddleware(tlog.WithSkipPaths("/health")))

	// Health check - system endpoint, keep at root level
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public API
	api := router.Group("/api")

	// V1 API Group
	v1 := api.Group("/v1")
	{
		routes.registerAuthRoutes(v1)
		routes.registerMobileConfigRoutes(v1)
		routes.registerPayloadPropertyDefinitionRoutes(v1)

		// Protected V1 routes
		protected := v1.Group("", mw.Auth(), mw.Authorize())
		{
			routes.registerUserRoutes(protected)
			routes.registerPolicyRoutes(protected)
			routes.registerMobileConfigRoutes(protected)
			routes.registerDashboardRoutes(protected)
			routes.registerDeviceRoutes(protected)
			routes.registerDeviceGroupRoutes(protected)
			routes.registerProfileRoutes(protected)
			routes.registerApplicationRoutes(protected)
			routes.registerAlertRoutes(protected)
			routes.registerReportRoutes(protected)
			routes.registerSettingRoutes(protected)
		}

		// NanoCMD Webhook (Public) - now consistently in v1
		v1.POST("/nanocmd/webhook", routes.nanocmd.Webhook)
	}

	// Swagger documentation
	// Resolve host/scheme from the current request so Swagger works behind public tunnels.
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	router.GET("/swagger/*any", func(c *gin.Context) {
		resolvedHost, resolvedScheme := resolveSwaggerEndpoint(c)
		docs.SwaggerInfo.Host = resolvedHost
		docs.SwaggerInfo.Schemes = []string{resolvedScheme}

		swaggerHandler(c)
	})

	return router
}

func resolveSwaggerEndpoint(c *gin.Context) (string, string) {
	host := firstHeaderValue(c.GetHeader("X-Forwarded-Host"))
	scheme := firstHeaderValue(c.GetHeader("X-Forwarded-Proto"))

	if originHost, originScheme := parseHostScheme(c.GetHeader("Origin")); originHost != "" {
		if host == "" {
			host = originHost
		}
		if scheme == "" {
			scheme = originScheme
		}
	}

	if refHost, refScheme := parseHostScheme(c.GetHeader("Referer")); refHost != "" {
		if host == "" {
			host = refHost
		}
		if scheme == "" {
			scheme = refScheme
		}
	}

	if host == "" {
		host = c.Request.Host
	}

	if scheme == "" {
		scheme = "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
	}

	return host, strings.ToLower(scheme)
}

func firstHeaderValue(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.Split(value, ",")
	return strings.TrimSpace(parts[0])
}

func parseHostScheme(raw string) (string, string) {
	if raw == "" {
		return "", ""
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", ""
	}
	return parsed.Host, parsed.Scheme
}

func (r *routeRegister) registerAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", r.auth.Login)
		auth.POST("/refresh", r.auth.Refresh)
		auth.POST("/logout", r.mw.Auth(), r.auth.Logout)
		auth.GET("/me", r.mw.Auth(), r.auth.GetMe)
	}
}

func (r *routeRegister) registerUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.GET("", r.user.List)
		users.GET("/:id", r.user.GetByID)
		users.POST("", r.user.Create)
		users.PUT("/:id", r.user.Update)
		users.DELETE("/:id", r.user.Delete)
	}
}

func (r *routeRegister) registerPolicyRoutes(rg *gin.RouterGroup) {
	policies := rg.Group("/policies")
	{
		policies.GET("", r.policy.ListPolicies)
		policies.POST("", r.policy.AddPolicy)
		policies.DELETE("", r.policy.RemovePolicy)
		policies.GET("/role/:role", r.policy.GetPoliciesForRole)
	}

	roles := rg.Group("/roles")
	{
		roles.GET("", r.policy.ListRoles)
		roles.POST("", r.policy.AddRole)
		roles.DELETE("", r.policy.RemoveRole)
	}
}

func (r *routeRegister) registerMobileConfigRoutes(rg *gin.RouterGroup) {
	mobileConfig := rg.Group("/mobile-configs")
	{
		mobileConfig.GET("", r.mobileConfig.List)
		mobileConfig.GET("/:id", r.mobileConfig.GetByID)
		mobileConfig.POST("", r.mobileConfig.Create)
		mobileConfig.PUT("/:id", r.mobileConfig.Update)
		mobileConfig.DELETE("/:id", r.mobileConfig.Delete)
		mobileConfig.GET("/:id/xml", r.mobileConfig.GetXML)
	}
}

func (r *routeRegister) registerDashboardRoutes(rg *gin.RouterGroup) {
	dashboard := rg.Group("/dashboard")
	{
		dashboard.GET("/stats", r.dashboard.GetStats)
		dashboard.GET("/device-stats", r.dashboard.GetDeviceStats)
		dashboard.GET("/alerts-summary", r.dashboard.GetAlertsSummary)
		dashboard.GET("/charts/:type", r.dashboard.GetChartData)
	}
}

func (r *routeRegister) registerDeviceRoutes(rg *gin.RouterGroup) {
	devices := rg.Group("/devices")
	{
		// Read operations
		devices.GET("", r.device.List)
		devices.GET("/export", r.device.Export)
		devices.GET("/:id", r.device.GetByID)

		// Device actions
		devices.POST("/:id/lock", r.device.Lock)
		devices.POST("/:id/wipe", r.device.Wipe)
		devices.POST("/:id/restart", r.device.Restart)
		devices.POST("/:id/shutdown", r.device.Shutdown)
		devices.POST("/:id/install-profile", r.device.InstallProfile)
		devices.POST("/:id/remove-profile", r.device.RemoveProfile)
		devices.POST("/:id/request-info", r.device.RequestInfo)
	}
}

func (r *routeRegister) registerDeviceGroupRoutes(rg *gin.RouterGroup) {
	groups := rg.Group("/device-groups")
	{
		groups.GET("", r.deviceGroup.List)
		groups.POST("", r.deviceGroup.Create)
		groups.GET("/:id", r.deviceGroup.GetByID)
		groups.PUT("/:id", r.deviceGroup.Update)
		groups.DELETE("/:id", r.deviceGroup.Delete)
		groups.POST("/:id/devices", r.deviceGroup.AddDevices)
		groups.DELETE("/:id/devices/:deviceId", r.deviceGroup.RemoveDevice)
	}
}

func (r *routeRegister) registerProfileRoutes(rg *gin.RouterGroup) {
	profiles := rg.Group("/profiles")
	{
		profiles.GET("", r.profile.List)
		profiles.POST("", r.profile.Create)
		profiles.GET("/:id", r.profile.GetByID)
		profiles.PUT("/:id", r.profile.Update)
		profiles.DELETE("/:id", r.profile.Delete)
		profiles.PUT("/:id/status", r.profile.UpdateStatus)

		settings := profiles.Group("/:id/settings")
		{
			settings.PUT("/security", r.profile.UpdateSecuritySettings)
			settings.PUT("/network", r.profile.UpdateNetworkConfig)
			settings.PUT("/restrictions", r.profile.UpdateRestrictions)
			settings.PUT("/content-filter", r.profile.UpdateContentFilter)
			settings.PUT("/compliance", r.profile.UpdateComplianceRules)
		}

		assign := profiles.Group("/:id/assignments")
		{
			assign.GET("", r.profile.ListAssignments)
			assign.POST("", r.profile.Assign)
			assign.DELETE("/:assignmentId", r.profile.Unassign)
		}

		versions := profiles.Group("/:id/versions")
		{
			versions.GET("", r.profile.ListVersions)
			versions.POST("/:versionId/rollback", r.profile.Rollback)
		}

		profiles.GET("/:id/deployment-status", r.profile.GetDeploymentStatus)
		profiles.POST("/:id/repush", r.profile.Repush)
		profiles.POST("/:id/duplicate", r.profile.Duplicate)
	}
}

func (r *routeRegister) registerApplicationRoutes(rg *gin.RouterGroup) {
	apps := rg.Group("/applications")
	{
		apps.GET("", r.application.List)
		apps.POST("", r.application.Create)
		apps.GET("/:id", r.application.GetByID)
		apps.PUT("/:id", r.application.Update)
		apps.DELETE("/:id", r.application.Delete)

		versions := apps.Group("/:id/versions")
		{
			versions.GET("", r.application.ListVersions)
			versions.POST("", r.application.CreateVersion)
			versions.DELETE("/:versionId", r.application.DeleteVersion)
			versions.GET("/:versionId/deployments", r.application.ListDeployments)
		}

		apps.POST("/deployments", r.application.Deploy)
	}
}

func (r *routeRegister) registerAlertRoutes(rg *gin.RouterGroup) {
	alerts := rg.Group("/alerts")
	{
		alerts.GET("", r.alert.List)
		alerts.POST("", r.alert.Create)
		alerts.GET("/stats", r.alert.GetStats) // Needs to be above /:id
		alerts.GET("/:id", r.alert.GetByID)

		alerts.PUT("/:id/acknowledge", r.alert.Acknowledge)
		alerts.PUT("/:id/resolve", r.alert.Resolve)
		alerts.POST("/bulk-resolve", r.alert.BulkResolve)

		actions := alerts.Group("/:id/actions")
		{
			actions.POST("/lock", r.alert.LockDevice)
			actions.POST("/wipe", r.alert.WipeDevice)
			actions.POST("/push-policy", r.alert.PushPolicy)
			actions.POST("/message", r.alert.SendMessage)
		}

		rules := alerts.Group("/rules")
		{
			rules.GET("", r.alert.ListRules)
			rules.POST("", r.alert.CreateRule)
			rules.GET("/:id", r.alert.GetRuleByID)
			rules.PUT("/:id", r.alert.UpdateRule)
			rules.DELETE("/:id", r.alert.DeleteRule)
			rules.PUT("/:id/toggle", r.alert.ToggleRule)
		}
	}
}

func (r *routeRegister) registerReportRoutes(rg *gin.RouterGroup) {
	reports := rg.Group("/reports")
	{
		reports.GET("/devices/export", r.report.ExportDevices)
		reports.GET("/alerts/export", r.report.ExportAlerts)
		reports.GET("/applications/export", r.report.ExportApplications)
	}
}

func (r *routeRegister) registerSettingRoutes(rg *gin.RouterGroup) {
	settings := rg.Group("/settings")
	{
		settings.GET("", r.setting.List)
		settings.POST("", r.setting.Create)
		settings.GET("/:key", r.setting.GetByKey)
		settings.PUT("/:key", r.setting.Update)
		settings.DELETE("/:key", r.setting.Delete)
	}
}

func (r *routeRegister) registerPayloadPropertyDefinitionRoutes(rg *gin.RouterGroup) {
	definitions := rg.Group("/payload-property-definitions")
	{
		// definitions.GET("", r.ppd.List)
		definitions.GET("/payload-types", r.ppd.ListPayloadTypes)
		// definitions.GET("/:id", r.ppd.GetByID)
		// definitions.POST("", r.ppd.Create)
		// definitions.PUT("/:id", r.ppd.Update)
		// definitions.DELETE("/:id", r.ppd.Delete)
		definitions.POST("/import", r.ppd.Import)
		definitions.GET("/schema", r.ppd.GetNestedSchema)
	}
}
