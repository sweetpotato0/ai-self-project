package api

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// APIVersion API版本信息
type APIVersion struct {
	Version     string `json:"version"`
	Status      string `json:"status"` // stable, beta, deprecated
	Deprecated  bool   `json:"deprecated"`
	SunsetDate  string `json:"sunset_date,omitempty"`
	Description string `json:"description"`
}

// VersionManager API版本管理器
type VersionManager struct {
	versions map[string]APIVersion
	current  string
}

// NewVersionManager 创建版本管理器
func NewVersionManager() *VersionManager {
	vm := &VersionManager{
		versions: make(map[string]APIVersion),
		current:  "v1",
	}

	// 注册API版本
	vm.RegisterVersion("v1", APIVersion{
		Version:     "v1",
		Status:      "stable",
		Deprecated:  false,
		Description: "Initial stable API version",
	})

	return vm
}

// RegisterVersion 注册API版本
func (vm *VersionManager) RegisterVersion(version string, info APIVersion) {
	vm.versions[version] = info
}

// GetVersion 获取版本信息
func (vm *VersionManager) GetVersion(version string) (APIVersion, bool) {
	info, exists := vm.versions[version]
	return info, exists
}

// GetCurrentVersion 获取当前版本
func (vm *VersionManager) GetCurrentVersion() string {
	return vm.current
}

// SetCurrentVersion 设置当前版本
func (vm *VersionManager) SetCurrentVersion(version string) {
	vm.current = version
}

// GetAllVersions 获取所有版本
func (vm *VersionManager) GetAllVersions() map[string]APIVersion {
	return vm.versions
}

// IsVersionSupported 检查版本是否支持
func (vm *VersionManager) IsVersionSupported(version string) bool {
	_, exists := vm.versions[version]
	return exists
}

// IsVersionDeprecated 检查版本是否已弃用
func (vm *VersionManager) IsVersionDeprecated(version string) bool {
	if info, exists := vm.versions[version]; exists {
		return info.Deprecated
	}
	return false
}

// VersionMiddleware API版本中间件
func (vm *VersionManager) VersionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从URL路径中提取版本
		path := c.Request.URL.Path
		parts := strings.Split(path, "/")

		var version string
		for i, part := range parts {
			if strings.HasPrefix(part, "v") && len(part) > 1 {
				version = part
				// 将版本信息存储到上下文中
				c.Set("api_version", version)
				c.Set("version_index", i)
				break
			}
		}

		// 如果没有指定版本，使用当前版本
		if version == "" {
			version = vm.GetCurrentVersion()
			c.Set("api_version", version)
		}

		// 检查版本是否支持
		if !vm.IsVersionSupported(version) {
			c.JSON(400, gin.H{
				"error":              "Unsupported API version",
				"code":               "UNSUPPORTED_VERSION",
				"supported_versions": vm.GetSupportedVersions(),
				"current_version":    vm.GetCurrentVersion(),
			})
			c.Abort()
			return
		}

		// 检查版本是否已弃用
		if vm.IsVersionDeprecated(version) {
			info, _ := vm.GetVersion(version)
			c.Header("X-API-Version-Deprecated", "true")
			if info.SunsetDate != "" {
				c.Header("X-API-Version-Sunset", info.SunsetDate)
			}
		}

		c.Next()
	}
}

// GetSupportedVersions 获取支持的版本列表
func (vm *VersionManager) GetSupportedVersions() []string {
	var versions []string
	for version := range vm.versions {
		versions = append(versions, version)
	}
	return versions
}

// GetVersionFromContext 从上下文中获取版本
func GetVersionFromContext(c *gin.Context) string {
	if version, exists := c.Get("api_version"); exists {
		return version.(string)
	}
	return "v1" // 默认版本
}
