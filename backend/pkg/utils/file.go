package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"

	"github.com/google/uuid"
)

// ValidateImageFile 验证是否为图片文件
func ValidateImageFile(header *multipart.FileHeader) error {
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return fmt.Errorf("只能上传图片文件，当前文件类型：%s", contentType)
	}
	return nil
}

// ValidateFileSize 验证文件大小
func ValidateFileSize(header *multipart.FileHeader, maxSizeMB int) error {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	if header.Size > maxSize {
		return fmt.Errorf("文件大小不能超过%dMB，当前文件大小：%.2fMB", 
			maxSizeMB, float64(header.Size)/(1024*1024))
	}
	return nil
}

// GenerateUniqueFilename 生成唯一文件名
func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}

// GetFileExtension 获取文件扩展名（小写）
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// IsAllowedImageType 检查是否为允许的图片类型
func IsAllowedImageType(filename string) bool {
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"}
	ext := GetFileExtension(filename)
	
	return slices.Contains(allowedTypes, ext)
}

// IsAllowedDocumentType 检查是否为允许的文档类型
func IsAllowedDocumentType(filename string) bool {
	allowedTypes := []string{".pdf", ".doc", ".docx", ".txt", ".md", ".rtf"}
	ext := GetFileExtension(filename)
	
	return slices.Contains(allowedTypes, ext)
}

// GetContentTypeByExtension 根据文件扩展名获取Content-Type
func GetContentTypeByExtension(filename string) string {
	ext := GetFileExtension(filename)
	
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg", 
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".bmp":  "image/bmp",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".txt":  "text/plain",
		".md":   "text/markdown",
		".rtf":  "application/rtf",
	}
	
	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}
	return "application/octet-stream"
}

// FormatFileSize 格式化文件大小显示
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// SanitizeFilename 清理文件名，移除不安全字符
func SanitizeFilename(filename string) string {
	// 替换不安全字符为下划线
	unsafe := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", "\n", "\r", "\t"}
	sanitized := filename
	
	for _, char := range unsafe {
		sanitized = strings.ReplaceAll(sanitized, char, "_")
	}
	
	// 移除多余的空格和点
	sanitized = strings.TrimSpace(sanitized)
	sanitized = strings.Trim(sanitized, ".")
	
	// 如果文件名为空或太长，使用默认名称
	if sanitized == "" {
		sanitized = "unnamed_file"
	}
	if len(sanitized) > 255 {
		ext := filepath.Ext(sanitized)
		name := strings.TrimSuffix(sanitized, ext)
		sanitized = name[:255-len(ext)] + ext
	}
	
	return sanitized
}