package handler

import (
	"strconv"
	"time"

	"gin-web-framework/internal/model"
	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/logger"
	"gin-web-framework/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	auditService *service.AuditService
	logger       *logger.Logger
}

func NewAuditHandler(auditService *service.AuditService, logger *logger.Logger) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
		logger:       logger,
	}
}

// GetAuditLogs 获取审计日志列表
func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	var query model.AuditLogQuery
	
	// 设置默认值
	query.Page = 1
	query.Limit = 20
	query.OrderBy = "timestamp"
	query.Order = "desc"

	// 绑定查询参数
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "Invalid query parameters: "+err.Error())
		return
	}

	// 解析时间参数
	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", startTimeStr); err == nil {
			query.StartTime = t
		}
	}

	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", endTimeStr); err == nil {
			query.EndTime = t
		}
	}

	logs, total, err := h.auditService.GetAuditLogs(&query)
	if err != nil {
		h.logger.Errorf("Failed to get audit logs: %v", err)
		response.InternalServerError(c, "Failed to get audit logs")
		return
	}

	// 计算分页信息
	totalPages := int((total + int64(query.Limit) - 1) / int64(query.Limit))

	responseData := map[string]interface{}{
		"items":        logs,
		"total":        total,
		"page":         query.Page,
		"limit":        query.Limit,
		"total_pages":  totalPages,
		"has_next":     query.Page < totalPages,
		"has_previous": query.Page > 1,
	}

	response.Success(c, responseData)
	h.logger.Info("Audit logs retrieved successfully")
}

// GetAuditLogByID 获取单个审计日志详情
func (h *AuditHandler) GetAuditLogByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid audit log ID")
		return
	}

	auditLog, err := h.auditService.GetAuditLogByID(uint(id))
	if err != nil {
		if err.Error() == "audit log not found" {
			response.NotFound(c, "Audit log not found")
			return
		}

		h.logger.Errorf("Failed to get audit log: %v", err)
		response.InternalServerError(c, "Failed to get audit log")
		return
	}

	response.Success(c, auditLog)
	h.logger.Info("Audit log retrieved successfully")
}

// GetAuditLogStats 获取审计日志统计
func (h *AuditHandler) GetAuditLogStats(c *gin.Context) {
	var startTime, endTime time.Time
	
	// 解析时间参数，默认为最近7天
	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", startTimeStr); err == nil {
			startTime = t
		} else {
			startTime = time.Now().AddDate(0, 0, -7)
		}
	} else {
		startTime = time.Now().AddDate(0, 0, -7)
	}

	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", endTimeStr); err == nil {
			endTime = t
		} else {
			endTime = time.Now()
		}
	} else {
		endTime = time.Now()
	}

	stats, err := h.auditService.GetAuditLogStats(startTime, endTime)
	if err != nil {
		h.logger.Errorf("Failed to get audit log stats: %v", err)
		response.InternalServerError(c, "Failed to get audit log statistics")
		return
	}

	response.Success(c, stats)
	h.logger.Info("Audit log statistics retrieved successfully")
}

// DeleteAuditLogs 批量删除审计日志
func (h *AuditHandler) DeleteAuditLogs(c *gin.Context) {
	var req struct {
		BeforeTime string `json:"before_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters: "+err.Error())
		return
	}

	beforeTime, err := time.Parse("2006-01-02T15:04:05Z07:00", req.BeforeTime)
	if err != nil {
		response.BadRequest(c, "Invalid time format")
		return
	}

	deletedCount, err := h.auditService.DeleteAuditLogsBefore(beforeTime)
	if err != nil {
		h.logger.Errorf("Failed to delete audit logs: %v", err)
		response.InternalServerError(c, "Failed to delete audit logs")
		return
	}

	responseData := map[string]interface{}{
		"deleted_count": deletedCount,
	}

	response.Success(c, responseData)

	h.logger.WithFields(map[string]interface{}{
		"deleted_count": deletedCount,
		"before_time":   beforeTime,
	}).Info("Audit logs deleted successfully")
}