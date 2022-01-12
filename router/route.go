package router

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lilj_01/gin_gateway/controller"
	"github.com/lilj_01/gin_gateway/docs"
	"github.com/lilj_01/gin_gateway/middleware"
	"github.com/lilj_01/gin_gateway/public"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	docs.SwaggerInfo.Title = lib.GetStringConf("base.swagger.title")
	docs.SwaggerInfo.Description = lib.GetStringConf("base.swagger.desc")
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = lib.GetStringConf("base.swagger.host")
	docs.SwaggerInfo.BasePath = lib.GetStringConf("base.swagger.base_path")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	store, err := sessions.NewRedisStore(10,
		"tcp", "101.34.146.196:6379", "", []byte("secret"))
	if err != nil {
		log.Fatalf("sessions.NewRedisStore err: %v", err)
	}

	// 将login-controller 注册到router
	adminLoginRouter := router.Group("/admin_login")
	// 设置中间件
	adminLoginRouter.Use(
		sessions.Sessions(public.GoGatewaySessionName, store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.TranslationMiddleware())
	{
		controller.AdminLoginRegister(adminLoginRouter)
	}
	adminRouter := router.Group("/admin")
	// 设置中间件
	adminRouter.Use(
		sessions.Sessions(public.GoGatewaySessionName, store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware())
	{
		controller.AdminRegister(adminRouter)
	}
	return router
}
