package commands

import (
	"fmt"

	"gin-web-framework/config"
	"gin-web-framework/internal/database"
	"gin-web-framework/pkg/logger"

	"github.com/spf13/cobra"
)

// newUserCmd 用户管理命令
func newUserCmd() *cobra.Command {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
		Long: `Manage users in the system.

Available subcommands:
  create    Create a new user
  list      List all users`,
	}

	userCmd.AddCommand(
		newUserCreateCmd(),
		newUserListCmd(),
	)

	return userCmd
}

// newUserCreateCmd 创建用户命令
func newUserCreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Long: `Create a new user in the system.

Examples:
  gin-cli user create`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createUser()
		},
	}

	return createCmd
}

// newUserListCmd 列出用户命令
func newUserListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all users",
		Long: `List all users in the system.

Examples:
  gin-cli user list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listUsers()
		},
	}

	return listCmd
}

// createUser 创建用户
func createUser() error {
	// 初始化
	if err := initForUserCommands(); err != nil {
		return err
	}

	fmt.Println("✅ User creation functionality will be implemented here")
	return nil
}

// listUsers 列出用户
func listUsers() error {
	// 初始化
	if err := initForUserCommands(); err != nil {
		return err
	}

	fmt.Println("✅ User listing functionality will be implemented here")
	return nil
}

// initForUserCommands 初始化用户命令
func initForUserCommands() error {
	// 加载配置
	if err := config.Load(); err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// 初始化日志
	logger.Init()

	// 初始化数据库连接
	if err := database.Init(); err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	return nil
}
