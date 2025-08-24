package commands

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-web-framework/config"
	"gin-web-framework/internal/container"
	"gin-web-framework/internal/database"
	"gin-web-framework/internal/redis"
	"gin-web-framework/internal/router"
	"gin-web-framework/internal/telemetry"
	"gin-web-framework/pkg/logger"

	"github.com/spf13/cobra"
)

// newServeCmd 服务器命令
func newServeCmd() *cobra.Command {
	var port string
	var mode string

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the web server",
		Long: `Start the Gin web server with the specified configuration.

Examples:
  gin-cli serve                    # Start with default config
  gin-cli serve --port 3000        # Start on port 3000
  gin-cli serve --mode release     # Start in release mode`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return startServer(port, mode)
		},
	}

	// 标志
	serveCmd.Flags().StringVarP(&port, "port", "p", "", "server port (default from config)")
	serveCmd.Flags().StringVarP(&mode, "mode", "m", "", "server mode: debug, release, test (default from config)")

	return serveCmd
}

// startServer 启动服务器
func startServer(port, mode string) error {
	// 加载配置
	if err := config.Load(); err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// 覆盖配置
	cfg := config.Get()
	if port != "" {
		// 注意：这里需要更新配置结构
		// TODO: 实现配置的动态更新
		logger.Warnf("Port override not implemented in new config structure")
	}
	if mode != "" {
		// TODO: 实现配置的动态更新
		logger.Warnf("Mode override not implemented in new config structure")
	}

	// 初始化日志
	logger.Init()

	// 初始化数据库连接
	if err := database.Init(); err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	// 初始化Redis连接
	if err := redis.Init(); err != nil {
		return fmt.Errorf("failed to initialize Redis: %v", err)
	}

	// 初始化OpenTelemetry
	appCfg := cfg.GetApp()
	if err := telemetry.InitTelemetry(appCfg.Name, appCfg.Version); err != nil {
		logger.Warnf("Failed to initialize telemetry: %v", err)
	}

	// 创建依赖注入容器
	container := container.NewContainer(cfg, database.GetDB(), redis.GetRedisClient())

	// 设置路由
	r := router.Setup(container)

	// 启动服务器
	srv := &http.Server{
		Addr:    ":" + cfg.GetServer().Port,
		Handler: r,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}

		fmt.Println("ListenAndServe")
	}()

	logger.Info("Server started on port " + cfg.GetServer().Port)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭容器
	if err := container.Shutdown(ctx); err != nil {
		logger.Error("Failed to shutdown container: " + err.Error())
	}

	// 关闭OpenTelemetry
	if err := telemetry.ShutdownTelemetry(ctx); err != nil {
		logger.Error("Failed to shutdown telemetry: " + err.Error())
	}

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: " + err.Error())
	}

	// 关闭数据库连接
	database.Close()
	// 关闭Redis连接
	redis.Close()

	logger.Info("Server exited")
	return nil
}
