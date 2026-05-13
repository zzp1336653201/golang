package main

import (
	"fmt"
	"log"
	"os"

	"cms/internal/config"
	"cms/internal/handler"
	"cms/internal/middleware"
	"cms/internal/service"
	"cms/pkg/cache"
	"cms/pkg/database"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("日志初始化失败: %v", err)
	}
	sugar := logger.Sugar()
	defer logger.Sync()

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		sugar.Fatalf("配置加载失败: %v", err)
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	if err := database.InitMySQL(&cfg.Database); err != nil {
		sugar.Fatalf("MySQL 初始化失败: %v", err)
	}
	defer database.CloseMySQL()
	sugar.Info("MySQL 连接成功")

	// 初始化 Redis
	if err := initRedis(cfg); err != nil {
		sugar.Warnf("Redis 连接失败: %v，使用本地缓存", err)
	}
	defer closeRedis()
	sugar.Info("Redis 连接成功")

	// 初始化服务
	userService := service.NewUserService(database.DB)
	productService := service.NewProductService(database.DB)
	orderService := service.NewOrderService(database.DB, productService)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)

	// 创建 Gin 引擎
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// 限流中间件
	r.Use(middleware.RateLimitMiddleware(cfg.RateLimit))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 路由
	v1 := r.Group("/api/v1")
	{
		// 认证相关（无需登录）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// 商品相关（无需登录）
		products := v1.Group("/products")
		{
			products.GET("", productHandler.List)
			products.GET("/:id", productHandler.Get)
			products.GET("/:id/stock", productHandler.GetStock)

			// 需要登录
			products.Use(middleware.JWTAuth(userService)).POST("", productHandler.Create)
			products.Use(middleware.JWTAuth(userService)).PUT("/:id", productHandler.Update)
			products.Use(middleware.JWTAuth(userService)).DELETE("/:id", productHandler.Delete)
		}

		// 订单相关（需要登录）
		orders := v1.Group("/orders")
		orders.Use(middleware.JWTAuth(userService))
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.GetOrders)
			orders.POST("/:order_no/pay", orderHandler.PayOrder)
			orders.DELETE("/:order_no", orderHandler.CancelOrder)
		}

		// 用户相关（需要登录）
		user := v1.Group("/user")
		user.Use(middleware.JWTAuth(userService))
		{
			user.GET("/profile", userHandler.GetProfile)
		}
	}

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	sugar.Infof("服务启动中，监听端口: %d", cfg.Server.Port)

	if err := r.Run(addr); err != nil {
		sugar.Fatalf("服务启动失败: %v", err)
	}

	_ = sugar
	_ = os.Args
}

func initRedis(cfg *config.Config) error {
	return cache.InitRedis(&cfg.Redis)
}

func closeRedis() {
	cache.CloseRedis()
}
