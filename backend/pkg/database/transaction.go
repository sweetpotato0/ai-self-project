package database

import (
	"context"
	"time"

	"gin-web-framework/pkg/errors"
	"gin-web-framework/pkg/logger"

	"gorm.io/gorm"
)

// TransactionService 事务服务
type TransactionService struct {
	db     *gorm.DB
	logger logger.LoggerInterface
}

// NewTransactionService 创建事务服务
func NewTransactionService(db *gorm.DB, logger logger.LoggerInterface) *TransactionService {
	return &TransactionService{
		db:     db,
		logger: logger,
	}
}

// TransactionFunc 事务函数类型
type TransactionFunc func(tx *gorm.DB) error

// ExecuteInTransaction 在事务中执行操作
func (ts *TransactionService) ExecuteInTransaction(ctx context.Context, fn TransactionFunc) error {
	return ts.executeWithRetry(ctx, fn, 3) // 默认重试3次
}

// ExecuteInTransactionWithRetry 在事务中执行操作，支持重试
func (ts *TransactionService) ExecuteInTransactionWithRetry(ctx context.Context, fn TransactionFunc, maxRetries int) error {
	return ts.executeWithRetry(ctx, fn, maxRetries)
}

// executeWithRetry 执行事务，支持重试
func (ts *TransactionService) executeWithRetry(ctx context.Context, fn TransactionFunc, maxRetries int) error {
	var lastErr error
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		startTime := time.Now()
		
		err := ts.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return fn(tx)
		})

		duration := time.Since(startTime)
		
		if err == nil {
			// 成功执行
			if attempt > 1 {
				ts.logger.WithFields(map[string]any{
					"attempt":     attempt,
					"duration_ms": duration.Milliseconds(),
				}).Info("事务重试成功")
			}
			return nil
		}

		lastErr = err
		
		// 判断是否应该重试
		if !ts.shouldRetry(err) || attempt >= maxRetries {
			break
		}

		// 计算重试延迟（指数退避）
		delay := time.Duration(attempt) * 100 * time.Millisecond
		if delay > 2*time.Second {
			delay = 2 * time.Second
		}

		ts.logger.WithFields(map[string]any{
			"attempt":     attempt,
			"max_retries": maxRetries,
			"error":       err.Error(),
			"delay_ms":    delay.Milliseconds(),
			"duration_ms": duration.Milliseconds(),
		}).Warn("事务执行失败，准备重试")

		// 等待重试
		select {
		case <-ctx.Done():
			return errors.NewTimeoutError("事务执行超时", ctx.Err())
		case <-time.After(delay):
			// 继续重试
		}
	}

	// 所有重试都失败了
	ts.logger.WithFields(map[string]any{
		"max_retries": maxRetries,
		"final_error": lastErr.Error(),
	}).Error("事务执行最终失败")

	return ts.wrapDatabaseError(lastErr)
}

// shouldRetry 判断错误是否应该重试
func (ts *TransactionService) shouldRetry(err error) bool {
	// 检查是否是可重试的数据库错误
	// 这里可以根据具体的数据库类型和错误类型来判断
	
	// 死锁错误通常可以重试
	if isDatabaseDeadlock(err) {
		return true
	}
	
	// 连接错误可能可以重试
	if isDatabaseConnectionError(err) {
		return true
	}
	
	// 事务冲突可能可以重试
	if isDatabaseConflictError(err) {
		return true
	}
	
	return false
}

// isDatabaseDeadlock 检查是否是死锁错误
func isDatabaseDeadlock(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := err.Error()
	// MySQL死锁错误码
	return contains(errStr, "Deadlock found") || 
		   contains(errStr, "deadlock detected") ||
		   contains(errStr, "Error 1213")
}

// isDatabaseConnectionError 检查是否是连接错误
func isDatabaseConnectionError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := err.Error()
	return contains(errStr, "connection refused") ||
		   contains(errStr, "connection reset") ||
		   contains(errStr, "connection lost") ||
		   contains(errStr, "server has gone away")
}

// isDatabaseConflictError 检查是否是冲突错误
func isDatabaseConflictError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := err.Error()
	return contains(errStr, "Lock wait timeout") ||
		   contains(errStr, "try restarting transaction")
}

// contains 检查字符串是否包含子串（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
			len(substr) == 0 || 
			indexIgnoreCase(s, substr) >= 0)
}

// indexIgnoreCase 不区分大小写的字符串查找
func indexIgnoreCase(s, substr string) int {
	s = toLowerCase(s)
	substr = toLowerCase(substr)
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// toLowerCase 转换为小写
func toLowerCase(s string) string {
	result := make([]byte, len(s))
	for i, b := range []byte(s) {
		if b >= 'A' && b <= 'Z' {
			result[i] = b + 32
		} else {
			result[i] = b
		}
	}
	return string(result)
}

// wrapDatabaseError 包装数据库错误
func (ts *TransactionService) wrapDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	if err == gorm.ErrRecordNotFound {
		return errors.NewNotFoundError("记录")
	}

	if isDatabaseDeadlock(err) {
		return errors.NewDatabaseError("数据库死锁", err)
	}

	if isDatabaseConnectionError(err) {
		return errors.NewDatabaseError("数据库连接失败", err)
	}

	return errors.NewDatabaseError("数据库操作失败", err)
}

// BatchOperation 批量操作结果
type BatchOperation struct {
	BatchSize    int
	TotalRecords int
	Processed    int
	Errors       []error
	Duration     time.Duration
}

// ExecuteBatch 批量执行操作（分批处理，避免长事务）
func (ts *TransactionService) ExecuteBatch(ctx context.Context, totalCount int, batchSize int, fn func(offset, limit int, tx *gorm.DB) error) (*BatchOperation, error) {
	if batchSize <= 0 {
		batchSize = 1000 // 默认批次大小
	}

	startTime := time.Now()
	result := &BatchOperation{
		BatchSize:    batchSize,
		TotalRecords: totalCount,
		Processed:    0,
		Errors:       make([]error, 0),
	}

	for offset := 0; offset < totalCount; offset += batchSize {
		limit := batchSize
		if offset+limit > totalCount {
			limit = totalCount - offset
		}

		// 每批操作都在单独的事务中执行
		err := ts.ExecuteInTransaction(ctx, func(tx *gorm.DB) error {
			return fn(offset, limit, tx)
		})

		if err != nil {
			result.Errors = append(result.Errors, err)
			ts.logger.WithFields(map[string]any{
				"offset": offset,
				"limit":  limit,
				"error":  err.Error(),
			}).Error("批量操作失败")
			
			// 根据错误类型决定是否继续
			if ts.shouldStopBatch(err) {
				break
			}
		} else {
			result.Processed += limit
		}

		// 检查是否需要取消
		select {
		case <-ctx.Done():
			result.Errors = append(result.Errors, ctx.Err())
			break
		default:
			// 继续处理
		}
	}

	result.Duration = time.Since(startTime)
	
	ts.logger.WithFields(map[string]any{
		"total_records":  result.TotalRecords,
		"processed":      result.Processed,
		"error_count":    len(result.Errors),
		"duration_ms":    result.Duration.Milliseconds(),
		"batch_size":     result.BatchSize,
	}).Info("批量操作完成")

	if len(result.Errors) > 0 {
		return result, errors.NewDatabaseError("批量操作部分失败", result.Errors[0])
	}

	return result, nil
}

// shouldStopBatch 判断是否应该停止批量操作
func (ts *TransactionService) shouldStopBatch(err error) bool {
	// 连接错误应该停止
	if isDatabaseConnectionError(err) {
		return true
	}
	
	// 其他错误可以继续处理下一批
	return false
}