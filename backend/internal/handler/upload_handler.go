package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gin-web-framework/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// UploadImage 上传图片
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "获取文件失败")
		return
	}
	defer file.Close()

	// 检查文件类型
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		response.Error(c, http.StatusBadRequest, "只能上传图片文件")
		return
	}

	// 检查文件大小 (5MB)
	if header.Size > 5*1024*1024 {
		response.Error(c, http.StatusBadRequest, "文件大小不能超过5MB")
		return
	}

	// 创建上传目录
	uploadDir := "uploads/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建目录失败")
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// 创建文件
	dst, err := os.Create(filepath)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建文件失败")
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		response.Error(c, http.StatusInternalServerError, "保存文件失败")
		return
	}

	// 返回文件URL
	fileURL := fmt.Sprintf("/uploads/images/%s", filename)

	response.Success(c, gin.H{
		"url":      fileURL,
		"filename": filename,
		"size":     header.Size,
		"type":     contentType,
	})
}

// ServeImage 提供图片服务
func (h *UploadHandler) ServeImage(c *gin.Context) {
	filename := c.Param("filename")
	filepath := filepath.Join("uploads/images", filename)

	// 检查文件是否存在
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		c.Status(http.StatusNotFound)
		return
	}

	c.File(filepath)
}
