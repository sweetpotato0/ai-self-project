package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// 全局标志
	configFile string
	verbose    bool
)

// NewRootCmd 创建根命令
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gin-cli",
		Short: "Gin Web Framework CLI tool",
		Long: `A comprehensive CLI tool for managing the Gin Web Framework application.

This tool provides commands for:
- Server management (start, stop, status)
- Database operations (migrate, seed, backup)
- User management (create, list, update, delete)
- System health checks
- Configuration management

Examples:
  gin-cli serve                    # Start the server
  gin-cli migrate                  # Run database migrations
  gin-cli user create              # Create a new user
  gin-cli health check             # Check system health`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// 全局预处理逻辑
			if verbose {
				fmt.Println("Verbose mode enabled")
			}
		},
	}

	// 全局标志
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is .env)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// 添加子命令
	rootCmd.AddCommand(
		newServeCmd(),
		newMigrateCmd(),
		newUserCmd(),
		newHealthCmd(),
		newVersionCmd(),
	)

	return rootCmd
}

// newVersionCmd 版本命令
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Gin Web Framework CLI v1.0.0")
			fmt.Println("Go version: go1.23.0")
		},
	}
}
