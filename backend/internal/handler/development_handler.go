package handler

import (
	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/logger"
)

// DevelopmentHandler 开发工具处理器
type DevelopmentHandler struct {
	logger       logger.LoggerInterface
	toolsService service.ToolsServiceInterface
}

// NewDevelopmentHandler 创建开发工具处理器
func NewDevelopmentHandler(logger logger.LoggerInterface, toolsService service.ToolsServiceInterface) *DevelopmentHandler {
	return &DevelopmentHandler{
		logger:       logger,
		toolsService: toolsService,
	}
}

// TODO: 这里可以添加开发工具相关的方法
// 例如：代码格式化、JSON处理、时间戳转换等