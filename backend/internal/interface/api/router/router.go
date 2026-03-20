package router

import (
	"github.com/gin-gonic/gin"
	"github.com/thienel/tlog"

	"github.com/thienel/go-backend-template/internal/interface/api/handler"
	"github.com/thienel/go-backend-template/internal/interface/api/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/thienel/go-backend-template/docs"
)

type routeRegister struct {
	auth    handler.AuthHandler
	user    handler.UserHandler
	policy  handler.PolicyHandler
	mdm     handler.MDMHandler
	dep     handler.DEPHandler
	nanocmd handler.NanoCMDHandler
	mw      *middleware.Middleware
}

// SetupRouter configures all routes following THD-Checkin-App pattern
func SetupRouter(
	authHandler handler.AuthHandler,
	userHandler handler.UserHandler,
	policyHandler handler.PolicyHandler,
	mdmHandler handler.MDMHandler,
	depHandler handler.DEPHandler,
	nanocmdHandler handler.NanoCMDHandler,
	mw *middleware.Middleware,
) *gin.Engine {

	routes := routeRegister{
		auth:    authHandler,
		user:    userHandler,
		policy:  policyHandler,
		mdm:     mdmHandler,
		dep:     depHandler,
		nanocmd: nanocmdHandler,
		mw:      mw,
	}

	router := gin.New()
	router.Use(gin.Recovery(), mw.CORS(), tlog.GinMiddleware(tlog.WithSkipPaths("/health")))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public API
	api := router.Group("/api")

	// V1 API Group
	v1 := api.Group("/v1")
	{
		routes.registerAuthRoutes(v1)

		// Protected V1 routes
		protected := v1.Group("", mw.Auth(), mw.Authorize())
		{
			routes.registerUserRoutes(protected)
			routes.registerPolicyRoutes(protected)
			routes.registerMDMRoutes(protected)
			routes.registerDEPRoutes(protected)
			routes.registerNanoCMDRoutes(protected)
		}
	}

	// NanoCMD Webhook (Public) - keeping at root /api/v1 or root?
	// Spec shows it might be useful to keep it as /api/v1/nanocmd/webhook
	v1.POST("/nanocmd/webhook", routes.nanocmd.Webhook)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func (r *routeRegister) registerAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", r.auth.Login)
		auth.POST("/logout", r.auth.Logout)
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

func (r *routeRegister) registerMDMRoutes(rg *gin.RouterGroup) {
	mdm := rg.Group("/mdm")
	{
		mdm.POST("/pushcert", r.mdm.PushCert)
		mdm.GET("/pushcert", r.mdm.GetCert)
		mdm.PUT("/enqueue/:id", r.mdm.EnqueueCommand)
	}
}

func (r *routeRegister) registerDEPRoutes(rg *gin.RouterGroup) {
	dep := rg.Group("/dep")
	{
		dep.GET("/names", r.dep.ListNames)
		dep.PUT("/token/:name", r.dep.PutToken)
		dep.GET("/token/:name", r.dep.GetToken)
		dep.GET("/tokens/:name", r.dep.GetTokens)
		dep.GET("/config/:name", r.dep.GetConfig)
		dep.GET("/assigner/:name", r.dep.GetAssigner)
		dep.PUT("/assigner/:name", r.dep.SetAssigner)
		dep.POST("/sync", r.dep.SyncDevices)
		dep.POST("/profile", r.dep.DefineProfile)
		dep.GET("/profiles", r.dep.ListProfiles)
		dep.GET("/profile/:uuid", r.dep.GetProfile)
		dep.POST("/disown", r.dep.DisownDevice)

		proxy := dep.Group("/proxy/:name")
		{
			proxy.GET("/account", r.dep.GetAccount)
			proxy.GET("/profile", r.dep.GetProfile)
			proxy.POST("/profile", r.dep.DefineProfile)
			proxy.POST("/devices", r.dep.GetDevices)
			proxy.POST("/devices/sync", r.dep.SyncDevices)
			proxy.POST("/devices/disown", r.dep.DisownDevice)
		}
	}
}

func (r *routeRegister) registerNanoCMDRoutes(rg *gin.RouterGroup) {
	nanocmd := rg.Group("/nanocmd")
	{
		nanocmd.GET("/version", r.nanocmd.GetVersion)
		nanocmd.POST("/workflow/:name/start", r.nanocmd.StartWorkflow)
		nanocmd.GET("/event/:name", r.nanocmd.GetEvent)
		nanocmd.PUT("/event/:name", r.nanocmd.PutEvent)
		nanocmd.GET("/fvenable/profiletemplate", r.nanocmd.GetFVEnableProfileTemplate)
		nanocmd.GET("/profile/:name", r.nanocmd.GetProfile)
		nanocmd.PUT("/profile/:name", r.nanocmd.PutProfile)
		nanocmd.DELETE("/profile/:name", r.nanocmd.DeleteProfile)
		nanocmd.GET("/profiles", r.nanocmd.GetProfiles)
		nanocmd.GET("/cmdplan/:name", r.nanocmd.GetCMDPlan)
		nanocmd.PUT("/cmdplan/:name", r.nanocmd.PutCMDPlan)
		nanocmd.GET("/inventory", r.nanocmd.GetInventory)
	}
}
