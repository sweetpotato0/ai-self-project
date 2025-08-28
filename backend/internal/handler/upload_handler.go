package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"gin-web-framework/pkg/response"
	"gin-web-framework/pkg/utils"

	"github.com/gin-gonic/gin"
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
	if err := utils.ValidateImageFile(header); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 检查文件大小 (5MB)
	if err := utils.ValidateFileSize(header, 5); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 创建上传目录
	uploadDir := "uploads/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建目录失败")
		return
	}

	// 生成唯一文件名
	filename := utils.GenerateUniqueFilename(header.Filename)
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
	contentType := header.Header.Get("Content-Type")

	response.Success(c, gin.H{
		"url":      fileURL,
		"filename": filename,
		"size":     utils.FormatFileSize(header.Size),
		"size_bytes": header.Size,
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
