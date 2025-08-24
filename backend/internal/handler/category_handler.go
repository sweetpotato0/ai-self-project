package handler

import (
	"gin-web-framework/internal/models"
	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryHandler 分类处理器
type CategoryHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		categoryService: service.NewCategoryService(),
	}
}

// CreateCategory 创建分类
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 设置创建者ID
	userID := c.GetUint("user_id")
	category.CreatedBy = userID

	// 验证必填字段
	if category.Name == "" {
		response.BadRequest(c, "分类名称不能为空")
		return
	}

	if err := h.categoryService.CreateCategory(&category); err != nil {
		response.InternalServerError(c, "创建分类失败: "+err.Error())
		return
	}

	response.Success(c, category)
}

// GetCategoryTree 获取分类树
func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
	userID := c.GetUint("user_id")

	categories, err := h.categoryService.GetCategoryTree(userID)
	if err != nil {
		response.InternalServerError(c, "获取分类树失败: "+err.Error())
		return
	}

	response.Success(c, categories)
}

// GetAllCategories 获取所有分类（扁平结构）
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	userID := c.GetUint("user_id")

	categories, err := h.categoryService.GetAllCategories(userID)
	if err != nil {
		response.InternalServerError(c, "获取分类列表失败: "+err.Error())
		return
	}

	response.Success(c, categories)
}

// GetCategoryByID 根据ID获取分类
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的分类ID: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")

	category, err := h.categoryService.GetCategoryByID(uint(id), userID)
	if err != nil {
		response.NotFound(c, "分类不存在: "+err.Error())
		return
	}

	response.Success(c, category)
}

// UpdateCategory 更新分类
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的分类ID: "+err.Error())
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")

	if err := h.categoryService.UpdateCategory(uint(id), userID, updates); err != nil {
		response.InternalServerError(c, "更新分类失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteCategory 删除分类
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的分类ID: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")

	if err := h.categoryService.DeleteCategory(uint(id), userID); err != nil {
		response.InternalServerError(c, "删除分类失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// GetCategoryStats 获取分类统计信息
func (h *CategoryHandler) GetCategoryStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	stats, err := h.categoryService.GetCategoryStats(userID)
	if err != nil {
		response.InternalServerError(c, "获取分类统计失败: "+err.Error())
		return
	}

	response.Success(c, stats)
}
