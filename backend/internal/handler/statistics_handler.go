package handler

import (
	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/logger"
	"gin-web-framework/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StatisticsHandler 统计处理器
type StatisticsHandler struct {
	statisticsService service.StatisticsServiceInterface
	logger           logger.LoggerInterface
}

// NewStatisticsHandler 创建统计处理器
func NewStatisticsHandler(statisticsService service.StatisticsServiceInterface, logger logger.LoggerInterface) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
		logger:           logger,
	}
}

// GetStatistics 获取统一统计数据
func (h *StatisticsHandler) GetStatistics(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	statTypeStr := c.Query("type")
	if statTypeStr == "" {
		response.BadRequest(c, "统计类型不能为空")
		return
	}

	// 验证统计类型
	statType, err := service.ValidateStatisticsType(statTypeStr)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取统计数据
	stats, err := h.statisticsService.GetStatistics(statType, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取统计数据失败: "+err.Error())
		return
	}

	response.Success(c, stats)
}

// GetTrends 获取趋势数据
func (h *StatisticsHandler) GetTrends(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	statTypeStr := c.Query("type")
	if statTypeStr == "" {
		response.BadRequest(c, "统计类型不能为空")
		return
	}

	// 验证统计类型
	statType, err := service.ValidateStatisticsType(statTypeStr)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证天数参数
	daysStr := c.DefaultQuery("days", "7")
	days, err := service.ValidateDays(daysStr)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取趋势数据
	trends, err := h.statisticsService.GetTrends(statType, userID, days)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取趋势数据失败: "+err.Error())
		return
	}

	response.Success(c, trends)
}
