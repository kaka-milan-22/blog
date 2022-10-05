package routers

import (
	"blog/app/model"
	"blog/config"
	_ "blog/docs"
	"blog/mid"

	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var Enforcer *casbin.Enforcer

type Option func(*gin.Engine)

var options = []Option{}

// 注册app的路由配置
func Include(opts ...Option) {
	options = append(options, opts...)
}

// 初始化
func Init() *gin.Engine {
	// gin.DisableConsoleColor()
	gin.ForceConsoleColor()

	gin.SetMode(config.Cfg.Server.Mode)
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(mid.JWTAuth())

	adapter := gormadapter.NewAdapterByDB(model.DB)
	Enforcer = casbin.NewEnforcer("rbac_models.conf", adapter)
	Enforcer.EnableLog(true)
	r.Use(mid.Authorize(Enforcer))

	for _, opt := range options {
		opt(r)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
