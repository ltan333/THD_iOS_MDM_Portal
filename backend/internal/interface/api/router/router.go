package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/thienel/tlog"

	"github.com/thienel/go-backend-template/internal/interface/api/handler"
	"github.com/thienel/go-backend-template/internal/interface/api/middleware"
)

type routeRegister struct {
	auth          handler.AuthHandler
	user          handler.UserHandler
	policy        handler.PolicyHandler
	mobile_config handler.MobileConfigHandler
	mw            *middleware.Middleware
}

// SetupRouter configures all routes following THD-Checkin-App pattern
func SetupRouter(
	authHandler handler.AuthHandler,
	userHandler handler.UserHandler,
	policyHandler handler.PolicyHandler,
	mobileConfigHandler handler.MobileConfigHandler,
	mw *middleware.Middleware,
) *gin.Engine {

	routes := routeRegister{
		auth:          authHandler,
		user:          userHandler,
		policy:        policyHandler,
		mobile_config: mobileConfigHandler,
		mw:            mw,
	}

	router := gin.New()
	router.Use(gin.Recovery(), mw.CORS(), tlog.GinMiddleware(tlog.WithSkipPaths("/health")))

	// Health check
	router.GET("/health", handler.Health)

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public API
	api := router.Group("/api")
	{
		routes.registerAuthRoutes(api)
		routes.registerMobileConfigRoutes(api)
	}

	// Protected API: Authentication (JWT) + Authorization (Casbin)
	protected := api.Group("", mw.Auth(), mw.Authorize())
	{
		routes.registerUserRoutes(protected)
		routes.registerPolicyRoutes(protected)
	}

	return router
}

func (r *routeRegister) registerAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", r.auth.Login)
		auth.POST("/logout", r.auth.Logout)
	}

	// Protected auth routes
	authProtected := auth.Group("", r.mw.Auth())
	{
		authProtected.GET("/me", r.auth.GetMe)
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
	mobileConfigs := rg.Group("/mobile-configs")
	{
		mobileConfigs.GET("/:id/xml", r.mobile_config.GetXML)
		mobileConfigs.POST("", r.mobile_config.Create)
		mobileConfigs.PUT("/:id", r.mobile_config.Update)
		mobileConfigs.DELETE("/:id", r.mobile_config.Delete)
	}
}
