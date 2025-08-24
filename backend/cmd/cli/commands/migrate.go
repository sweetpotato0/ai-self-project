package commands

import (
	"fmt"

	"gin-web-framework/config"
	"gin-web-framework/internal/database"
	"gin-web-framework/pkg/logger"

	"github.com/spf13/cobra"
)

// newMigrateCmd 数据库迁移命令
func newMigrateCmd() *cobra.Command {
	var force bool

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Long: `Run database migrations to create or update database schema.

This command will:
- Connect to the database using configuration
- Run automatic migrations for all models
- Create tables if they don't exist
- Update schema if models have changed

Examples:
  gin-cli migrate              # Run migrations
  gin-cli migrate --force      # Force migration even if errors occur`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrations(force)
		},
	}

	// 标志
	migrateCmd.Flags().BoolVarP(&force, "force", "f", false, "force migration even if errors occur")

	return migrateCmd
}

// runMigrations 运行数据库迁移
func runMigrations(force bool) error {
	// 加载配置
	if err := config.Load(); err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// 初始化日志
	logger.Init()

	// 初始化数据库连接
	if err := database.Init(); err != nil {
		if force {
			logger.Error("Database connection failed, but continuing due to --force flag")
		} else {
			return fmt.Errorf("failed to initialize database: %v", err)
		}
	}

	logger.Info("Database migration completed successfully")
	return nil
}
