package validator

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"gin-web-framework/config"
)

// ValidateConfig 验证配置
func ValidateConfig(cfg *config.Config) error {
	var errs []string

	// 验证服务器配置
	if err := validateServerConfig(&cfg.Server); err != nil {
		errs = append(errs, fmt.Sprintf("server config: %v", err))
	}

	// 验证数据库配置
	if err := validateDatabaseConfig(&cfg.Database); err != nil {
		errs = append(errs, fmt.Sprintf("database config: %v", err))
	}

	// 验证Redis配置
	if err := validateRedisConfig(&cfg.Redis); err != nil {
		errs = append(errs, fmt.Sprintf("redis config: %v", err))
	}

	// 验证JWT配置
	if err := validateJWTConfig(&cfg.JWT); err != nil {
		errs = append(errs, fmt.Sprintf("jwt config: %v", err))
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

func validateServerConfig(cfg *config.ServerConfig) error {
	// 验证端口
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return fmt.Errorf("invalid port: %v", err)
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	// 验证模式
	validModes := []string{"debug", "release", "test"}
	found := false
	for _, mode := range validModes {
		if cfg.Mode == mode {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("invalid mode: %s, must be one of %v", cfg.Mode, validModes)
	}

	return nil
}

func validateDatabaseConfig(cfg *config.DatabaseConfig) error {
	// 验证驱动
	validDrivers := []string{"postgres", "mysql"}
	found := false
	for _, driver := range validDrivers {
		if cfg.Driver == driver {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("invalid driver: %s, must be one of %v", cfg.Driver, validDrivers)
	}

	// 验证主机
	if cfg.Host == "" {
		return errors.New("host cannot be empty")
	}

	// 验证端口
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return fmt.Errorf("invalid port: %v", err)
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	// 验证用户名和数据库名
	if cfg.Username == "" {
		return errors.New("username cannot be empty")
	}
	if cfg.DBName == "" {
		return errors.New("database name cannot be empty")
	}

	// 验证PostgreSQL特有配置
	if cfg.Driver == "postgres" {
		validSSLModes := []string{"disable", "require", "verify-ca", "verify-full"}
		found := false
		for _, mode := range validSSLModes {
			if cfg.SSLMode == mode {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid ssl mode: %s, must be one of %v", cfg.SSLMode, validSSLModes)
		}
	}

	return nil
}

func validateRedisConfig(cfg *config.RedisConfig) error {
	// 验证主机
	if cfg.Host == "" {
		return errors.New("host cannot be empty")
	}

	// 验证端口
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return fmt.Errorf("invalid port: %v", err)
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	// 验证数据库编号
	if cfg.DB < 0 || cfg.DB > 15 {
		return fmt.Errorf("db must be between 0 and 15")
	}

	// 验证主机连接
	if cfg.Host != "localhost" && cfg.Host != "127.0.0.1" {
		if net.ParseIP(cfg.Host) == nil {
			// 尝试解析主机名
			_, err := net.LookupHost(cfg.Host)
			if err != nil {
				return fmt.Errorf("invalid host: %s", cfg.Host)
			}
		}
	}

	return nil
}

func validateJWTConfig(cfg *config.JWTConfig) error {
	// 验证密钥
	if cfg.Secret == "" {
		return errors.New("secret cannot be empty")
	}
	if len(cfg.Secret) < 32 {
		return errors.New("secret should be at least 32 characters long")
	}

	// 验证过期时间
	if cfg.ExpireHours <= 0 {
		return errors.New("expire hours must be positive")
	}
	if cfg.ExpireHours > 24*30 { // 最长30天
		return errors.New("expire hours cannot exceed 30 days")
	}

	return nil
}
