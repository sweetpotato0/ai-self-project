package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// newHealthCmd 健康检查命令
func newHealthCmd() *cobra.Command {
	var url string
	var timeout int

	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "Check system health",
		Long: `Check the health status of the application.

Examples:
  gin-cli health check
  gin-cli health check --url http://localhost:3000`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkHealth(url, timeout)
		},
	}

	// 标志
	healthCmd.Flags().StringVarP(&url, "url", "u", "http://localhost:8080", "service URL")
	healthCmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "timeout in seconds")

	return healthCmd
}

// checkHealth 检查健康状态
func checkHealth(url string, timeout int) error {
	healthURL := url + "/api/v1/health"

	fmt.Printf("🔍 Checking health at: %s\n", healthURL)

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Get(healthURL)
	if err != nil {
		return fmt.Errorf("❌ Health check failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("❌ Failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("❌ Health check failed with status %d: %s", resp.StatusCode, string(body))
	}

	var healthResp HealthResponse
	if err := json.Unmarshal(body, &healthResp); err != nil {
		return fmt.Errorf("❌ Failed to parse response: %v", err)
	}

	fmt.Printf("✅ Service is healthy: %s\n", healthResp.Message)
	return nil
}
