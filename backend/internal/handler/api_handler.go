package handler

import (
	"strconv"

	"gin-web-framework/internal/api"
	"gin-web-framework/internal/middleware"
	"gin-web-framework/internal/service"

	"github.com/gin-gonic/gin"
)

// APIHealthCheck 使用新API响应格式的健康检查
func APIHealthCheck(c *gin.Context) {
	api.SuccessResponse(c, gin.H{
		"status":  "ok",
		"message": "Service is running",
		"version": api.GetVersionFromContext(c),
	})
}

// APIGetMetrics 使用新API响应格式的性能指标
func APIGetMetrics(c *gin.Context) {
	metrics := middleware.GetPerformanceMetrics()

	api.SuccessResponse(c, gin.H{
		"performance": map[string]interface{}{
			"request_count":   metrics.RequestCount,
			"avg_duration_ms": float64(metrics.AvgDuration.Nanoseconds()) / 1e6,
			"max_duration_ms": float64(metrics.MaxDuration.Nanoseconds()) / 1e6,
			"min_duration_ms": float64(metrics.MinDuration.Nanoseconds()) / 1e6,
			"error_count":     metrics.ErrorCount,
			"active_requests": metrics.ActiveRequests,
			"memory_usage_mb": float64(metrics.MemoryUsage) / 1024 / 1024,
			"goroutine_count": metrics.GoroutineCount,
		},
		"api_version": api.GetVersionFromContext(c),
	})
}

// APIRegister 使用新API响应格式的用户注册
func APIRegister(c *gin.Context) {
	var req service.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		api.ValidationErrorResponse(c, "Invalid request data", err.Error())
		return
	}

	userService := service.NewUserService()
	user, err := userService.Register(req)
	if err != nil {
		api.BadRequestResponse(c, err.Error())
		return
	}

	api.SuccessResponse(c, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// APILogin 使用新API响应格式的用户登录
func APILogin(c *gin.Context) {
	var req service.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		api.ValidationErrorResponse(c, "Invalid request data", err.Error())
		return
	}

	userService := service.NewUserService()
	response, err := userService.Login(req)
	if err != nil {
		api.BadRequestResponse(c, err.Error())
		return
	}

	api.SuccessResponse(c, gin.H{
		"message": "Login successful",
		"token":   response.Token,
		"user":    response.User,
	})
}

// APIGetProfile 使用新API响应格式的获取用户资料
func APIGetProfile(c *gin.Context) {
	userID, err := middleware.RequireAuth(c)
	if err != nil {
		api.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userService := service.NewUserService()
	user, err := userService.GetUserByID(userID)
	if err != nil {
		api.NotFoundResponse(c, err.Error())
		return
	}

	api.SuccessResponse(c, gin.H{
		"user": user,
	})
}

// APIGetTodoList 使用新API响应格式的获取TODO列表
func APIGetTodoList(c *gin.Context) {
	userID, err := middleware.RequireAuth(c)
	if err != nil {
		api.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// 解析分页参数
	var pagination api.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		api.ValidationErrorResponse(c, "Invalid pagination parameters", err.Error())
		return
	}

	// 验证分页参数
	if err := api.ValidatePagination(&pagination); err != nil {
		api.ValidationErrorResponse(c, "Invalid pagination parameters", err.Error())
		return
	}

	todoService := service.NewTodoService()
	todos, err := todoService.GetTodoList(userID)
	if err != nil {
		api.InternalServerErrorResponse(c, err.Error())
		return
	}

	// 简单的分页处理（实际项目中应该使用数据库分页）
	total := len(todos)
	start := pagination.GetOffset()
	end := start + pagination.GetLimit()
	if end > total {
		end = total
	}
	if start > total {
		start = total
	}

	pagedTodos := todos[start:end]
	totalPages := (total + pagination.GetLimit() - 1) / pagination.GetLimit()

	response := gin.H{
		"todos": pagedTodos,
		"meta": gin.H{
			"pagination": gin.H{
				"page":        pagination.Page,
				"page_size":   pagination.PageSize,
				"total":       total,
				"total_pages": totalPages,
				"has_next":    pagination.Page < totalPages,
				"has_prev":    pagination.Page > 1,
			},
		},
	}

	api.SuccessResponse(c, response)
}

// APICreateTodo 使用新API响应格式的创建TODO
func APICreateTodo(c *gin.Context) {
	userID, err := middleware.RequireAuth(c)
	if err != nil {
		api.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req service.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.ValidationErrorResponse(c, "Invalid request data", err.Error())
		return
	}

	todoService := service.NewTodoService()
	todo, err := todoService.CreateTodo(req, userID)
	if err != nil {
		api.BadRequestResponse(c, err.Error())
		return
	}

	api.SuccessResponse(c, gin.H{
		"message": "TODO created successfully",
		"todo":    todo,
	})
}

// APIUpdateTodo 使用新API响应格式的更新TODO
func APIUpdateTodo(c *gin.Context) {
	userID, err := middleware.RequireAuth(c)
	if err != nil {
		api.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		api.BadRequestResponse(c, "Invalid TODO ID")
		return
	}

	var req service.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.ValidationErrorResponse(c, "Invalid request data", err.Error())
		return
	}

	todoService := service.NewTodoService()
	todo, err := todoService.UpdateTodo(uint(todoID), userID, req)
	if err != nil {
		api.BadRequestResponse(c, err.Error())
		return
	}

	api.SuccessResponse(c, gin.H{
		"message": "TODO updated successfully",
		"todo":    todo,
	})
}

// APIDeleteTodo 使用新API响应格式的删除TODO
func APIDeleteTodo(c *gin.Context) {
	userID, err := middleware.RequireAuth(c)
	if err != nil {
		api.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		api.BadRequestResponse(c, "Invalid TODO ID")
		return
	}

	todoService := service.NewTodoService()
	err = todoService.DeleteTodo(uint(todoID), userID)
	if err != nil {
		api.BadRequestResponse(c, err.Error())
		return
	}

	api.SuccessResponse(c, gin.H{
		"message": "TODO deleted successfully",
	})
}
