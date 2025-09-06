package service

import (
	"errors"
	"time"

	"gin-web-framework/internal/model"
	"gin-web-framework/pkg/logger"

	"gorm.io/gorm"
)

type AuditService struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewAuditService(db *gorm.DB, logger *logger.Logger) *AuditService {
	return &AuditService{
		db:     db,
		logger: logger,
	}
}

// GetAuditLogs 获取审计日志列表
func (s *AuditService) GetAuditLogs(query *model.AuditLogQuery) ([]*model.AuditLog, int64, error) {
	var logs []*model.AuditLog
	var total int64

	// 构建查询条件
	db := s.db.Model(&model.AuditLog{})

	// 应用筛选条件
	if query.UserID != 0 {
		db = db.Where("user_id = ?", query.UserID)
	}

	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}

	if query.Action != "" {
		db = db.Where("action = ?", query.Action)
	}

	if query.Resource != "" {
		db = db.Where("resource = ?", query.Resource)
	}

	if query.Method != "" {
		db = db.Where("method = ?", query.Method)
	}

	if query.IPAddress != "" {
		db = db.Where("ip_address = ?", query.IPAddress)
	}

	if query.StatusCode != 0 {
		db = db.Where("status_code = ?", query.StatusCode)
	}

	if !query.StartTime.IsZero() {
		db = db.Where("timestamp >= ?", query.StartTime)
	}

	if !query.EndTime.IsZero() {
		db = db.Where("timestamp <= ?", query.EndTime)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		s.logger.Errorf("Failed to count audit logs: %v", err)
		return nil, 0, err
	}

	// 应用排序
	orderBy := "timestamp"
	if query.OrderBy != "" {
		switch query.OrderBy {
		case "timestamp", "duration", "status_code":
			orderBy = query.OrderBy
		}
	}

	order := "DESC"
	if query.Order == "asc" {
		order = "ASC"
	}

	// 应用分页
	offset := (query.Page - 1) * query.Limit
	if err := db.Order(orderBy + " " + order).
		Offset(offset).
		Limit(query.Limit).
		Find(&logs).Error; err != nil {
		s.logger.Errorf("Failed to get audit logs: %v", err)
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAuditLogByID 根据ID获取审计日志
func (s *AuditService) GetAuditLogByID(id uint) (*model.AuditLog, error) {
	var auditLog model.AuditLog
	if err := s.db.First(&auditLog, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("audit log not found")
		}
		s.logger.Errorf("Failed to get audit log by ID: %v", err)
		return nil, err
	}

	return &auditLog, nil
}

// GetAuditLogStats 获取审计日志统计
func (s *AuditService) GetAuditLogStats(startTime, endTime time.Time) (*model.AuditLogStats, error) {
	stats := &model.AuditLogStats{}

	// 基础统计
	if err := s.getBasicStats(stats, startTime, endTime); err != nil {
		return nil, err
	}

	// 用户操作统计
	if err := s.getTopUsers(stats, startTime, endTime); err != nil {
		return nil, err
	}

	// 操作类型统计
	if err := s.getTopActions(stats, startTime, endTime); err != nil {
		return nil, err
	}

	// 资源类型统计
	if err := s.getTopResources(stats, startTime, endTime); err != nil {
		return nil, err
	}

	// 每小时分布统计
	if err := s.getHourlyDistribution(stats, startTime, endTime); err != nil {
		return nil, err
	}

	// 状态码分布统计
	if err := s.getStatusDistribution(stats, startTime, endTime); err != nil {
		return nil, err
	}

	return stats, nil
}

// getBasicStats 获取基础统计信息
func (s *AuditService) getBasicStats(stats *model.AuditLogStats, startTime, endTime time.Time) error {
	// 总数统计
	if err := s.db.Model(&model.AuditLog{}).
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Count(&stats.TotalCount).Error; err != nil {
		return err
	}

	// 今日统计
	today := time.Now().Format("2006-01-02")
	todayStart, _ := time.Parse("2006-01-02", today)
	todayEnd := todayStart.Add(24 * time.Hour)

	if err := s.db.Model(&model.AuditLog{}).
		Where("timestamp BETWEEN ? AND ?", todayStart, todayEnd).
		Count(&stats.TodayCount).Error; err != nil {
		return err
	}

	// 成功请求统计
	if err := s.db.Model(&model.AuditLog{}).
		Where("timestamp BETWEEN ? AND ? AND status_code >= 200 AND status_code < 400", startTime, endTime).
		Count(&stats.SuccessCount).Error; err != nil {
		return err
	}

	// 错误请求统计
	if err := s.db.Model(&model.AuditLog{}).
		Where("timestamp BETWEEN ? AND ? AND status_code >= 400", startTime, endTime).
		Count(&stats.ErrorCount).Error; err != nil {
		return err
	}

	return nil
}

// getTopUsers 获取活跃用户统计
func (s *AuditService) getTopUsers(stats *model.AuditLogStats, startTime, endTime time.Time) error {
	var results []model.UserActionCount

	if err := s.db.Model(&model.AuditLog{}).
		Select("user_id, username, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("user_id, username").
		Order("count DESC").
		Limit(10).
		Scan(&results).Error; err != nil {
		return err
	}

	stats.TopUsers = results
	return nil
}

// getTopActions 获取操作类型统计
func (s *AuditService) getTopActions(stats *model.AuditLogStats, startTime, endTime time.Time) error {
	var results []model.ActionCount

	if err := s.db.Model(&model.AuditLog{}).
		Select("action, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("action").
		Order("count DESC").
		Limit(10).
		Scan(&results).Error; err != nil {
		return err
	}

	stats.TopActions = results
	return nil
}

// getTopResources 获取资源类型统计
func (s *AuditService) getTopResources(stats *model.AuditLogStats, startTime, endTime time.Time) error {
	var results []model.ResourceCount

	if err := s.db.Model(&model.AuditLog{}).
		Select("resource, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("resource").
		Order("count DESC").
		Limit(10).
		Scan(&results).Error; err != nil {
		return err
	}

	stats.TopResources = results
	return nil
}

// getHourlyDistribution 获取每小时分布统计
func (s *AuditService) getHourlyDistribution(stats *model.AuditLogStats, startTime, endTime time.Time) error {
	var results []model.HourlyCount

	// 初始化24小时的数据
	hourlyMap := make(map[int]int64)
	for i := 0; i < 24; i++ {
		hourlyMap[i] = 0
	}

	// SQLite兼容的小时提取查询
	rows, err := s.db.Model(&model.AuditLog{}).
		Select("CAST(strftime('%H', timestamp) AS INTEGER) as hour, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("CAST(strftime('%H', timestamp) AS INTEGER)").
		Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var hour int
		var count int64
		if err := rows.Scan(&hour, &count); err != nil {
			continue
		}
		hourlyMap[hour] = count
	}

	// 转换为切片
	for hour := 0; hour < 24; hour++ {
		results = append(results, model.HourlyCount{
			Hour:  hour,
			Count: hourlyMap[hour],
		})
	}

	stats.HourlyDistribution = results
	return nil
}

// getStatusDistribution 获取状态码分布统计
func (s *AuditService) getStatusDistribution(stats *model.AuditLogStats, startTime, endTime time.Time) error {
	var results []model.StatusCount

	if err := s.db.Model(&model.AuditLog{}).
		Select("status_code, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("status_code").
		Order("count DESC").
		Scan(&results).Error; err != nil {
		return err
	}

	stats.StatusDistribution = results
	return nil
}

// DeleteAuditLogsBefore 删除指定时间之前的审计日志
func (s *AuditService) DeleteAuditLogsBefore(beforeTime time.Time) (int64, error) {
	result := s.db.Where("timestamp < ?", beforeTime).Delete(&model.AuditLog{})
	if result.Error != nil {
		s.logger.Errorf("Failed to delete audit logs: %v", result.Error)
		return 0, result.Error
	}

	s.logger.WithFields(map[string]interface{}{
		"deleted_count": result.RowsAffected,
		"before_time":   beforeTime,
	}).Info("Audit logs deleted successfully")

	return result.RowsAffected, nil
}

// CreateAuditLog 创建审计日志记录
func (s *AuditService) CreateAuditLog(auditLog *model.AuditLog) error {
	if err := s.db.Create(auditLog).Error; err != nil {
		s.logger.Errorf("Failed to create audit log: %v", err)
		return err
	}
	return nil
}
