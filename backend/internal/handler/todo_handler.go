package handler

import (
	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TodoHandler 任务处理器
type TodoHandler struct {
	todoService service.TodoServiceInterface
}

// NewTodoHandler 创建任务处理器
func NewTodoHandler(todoService service.TodoServiceInterface) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// CreateTodo 创建任务
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req service.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data: "+err.Error())
		return
	}

	todo, err := h.todoService.CreateTodo(req, userID.(uint))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Todo created successfully",
		"todo":    todo,
	})
}

// GetTodo 获取单个任务
func (h *TodoHandler) GetTodo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid todo ID")
		return
	}

	todo, err := h.todoService.GetTodoByID(uint(id), userID.(uint))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"todo": todo,
	})
}

// GetTodos 获取任务列表
func (h *TodoHandler) GetTodos(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// 解析查询参数
	filter := service.TodoFilter{}
	if status := c.Query("status"); status != "" {
		filter.Status = status
	}
	if priority := c.Query("priority"); priority != "" {
		filter.Priority = priority
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := strconv.ParseUint(categoryID, 10, 32); err == nil {
			uid := uint(id)
			filter.CategoryID = &uid
		}
	}

	todos, err := h.todoService.GetTodos(userID.(uint), filter)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"todos": todos.Todos,
	})
}

// UpdateTodo 更新任务
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid todo ID")
		return
	}

	var req service.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data: "+err.Error())
		return
	}

	todo, err := h.todoService.UpdateTodo(uint(id), userID.(uint), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Todo updated successfully",
		"todo":    todo,
	})
}

// DeleteTodo 删除任务
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid todo ID")
		return
	}

	err = h.todoService.DeleteTodo(uint(id), userID.(uint))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Todo deleted successfully",
	})
}

// MarkCompleted 标记任务为完成
func (h *TodoHandler) MarkCompleted(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid todo ID")
		return
	}

	err = h.todoService.MarkCompleted(uint(id), userID.(uint))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Todo marked as completed",
	})
}

// GetTodoStats 获取任务统计
func (h *TodoHandler) GetTodoStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	stats, err := h.todoService.GetTodoStats(userID.(uint))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"stats": stats,
	})
}
