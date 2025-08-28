package handler

import (
	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/logger"
)

// QueryHandler 查询工具处理器
type QueryHandler struct {
	logger       logger.LoggerInterface
	toolsService service.ToolsServiceInterface
}

// NewQueryHandler 创建查询工具处理器
func NewQueryHandler(logger logger.LoggerInterface, toolsService service.ToolsServiceInterface) *QueryHandler {
	return &QueryHandler{
		logger:       logger,
		toolsService: toolsService,
	}
}

// TODO: 这里可以添加查询工具相关的方法
// 例如：IP查询、Whois查询、域名信息等