package database

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Common error types
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidID      = errors.New("invalid ID")
	ErrDuplicateKey   = errors.New("duplicate key constraint")
	ErrForeignKey     = errors.New("foreign key constraint")
)

// HandleDBError 处理数据库错误，将GORM错误转换为业务错误
func HandleDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrRecordNotFound
	}

	// 检查是否为重复键错误
	errMsg := err.Error()
	if containsAny(errMsg, []string{"duplicate key", "UNIQUE constraint", "Duplicate entry"}) {
		return ErrDuplicateKey
	}

	// 检查是否为外键约束错误
	if containsAny(errMsg, []string{"foreign key constraint", "FOREIGN KEY constraint"}) {
		return ErrForeignKey
	}

	// 其他数据库错误
	return fmt.Errorf("database error: %w", err)
}

// HandleDBErrorWithContext 处理数据库错误并添加上下文信息
func HandleDBErrorWithContext(err error, operation string) error {
	if err == nil {
		return nil
	}

	dbErr := HandleDBError(err)
	return fmt.Errorf("%s failed: %w", operation, dbErr)
}

// IsNotFoundError 检查是否为记录未找到错误
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrRecordNotFound) || errors.Is(err, gorm.ErrRecordNotFound)
}

// IsDuplicateKeyError 检查是否为重复键错误
func IsDuplicateKeyError(err error) bool {
	return errors.Is(err, ErrDuplicateKey)
}

// IsForeignKeyError 检查是否为外键约束错误
func IsForeignKeyError(err error) bool {
	return errors.Is(err, ErrForeignKey)
}

// containsAny 检查字符串是否包含任意一个子字符串
func containsAny(s string, substrings []string) bool {
	for _, substring := range substrings {
		if len(s) >= len(substring) {
			for i := 0; i <= len(s)-len(substring); i++ {
				if s[i:i+len(substring)] == substring {
					return true
				}
			}
		}
	}
	return false
}

// CreateRecordError 创建记录失败错误
func CreateRecordError(entityName string, err error) error {
	return HandleDBErrorWithContext(err, fmt.Sprintf("create %s", entityName))
}

// UpdateRecordError 更新记录失败错误
func UpdateRecordError(entityName string, err error) error {
	return HandleDBErrorWithContext(err, fmt.Sprintf("update %s", entityName))
}

// DeleteRecordError 删除记录失败错误
func DeleteRecordError(entityName string, err error) error {
	return HandleDBErrorWithContext(err, fmt.Sprintf("delete %s", entityName))
}

// FindRecordError 查找记录失败错误
func FindRecordError(entityName string, err error) error {
	return HandleDBErrorWithContext(err, fmt.Sprintf("find %s", entityName))
}