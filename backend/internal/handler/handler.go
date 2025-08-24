package handler

import (
	"strconv"

	"gin-web-framework/internal/database"
	"gin-web-framework/internal/middleware"
	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/response"

	"github.com/gin-gonic/gin"
)

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "ok",
		"message": "Service is running",
	})
}

// GetMetrics 获取性能指标
func GetMetrics(c *gin.Context) {
	metrics := middleware.GetPerformanceMetrics()

	response.Success(c, gin.H{
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
		"timestamp": c.GetInt64("timestamp"),
	})
}

// Register 用户注册
func Register(c *gin.Context) {
	var req service.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data: "+err.Error())
		return
	}

	userService := service.NewUserService()
	user, err := userService.Register(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req service.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data: "+err.Error())
		return
	}

	userService := service.NewUserService()
	loginResponse, err := userService.Login(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Login successful",
		"token":   loginResponse.Token,
		"user":    loginResponse.User,
	})
}

// GetProfile 获取用户资料
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userService := service.NewUserService()
	user, err := userService.GetUserByID(userID.(uint))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"user": user,
	})
}

// UpdateProfile 更新用户资料
func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data: "+err.Error())
		return
	}

	userService := service.NewUserService()
	updateReq := service.UpdateProfileRequest{
		Email: req.Email,
	}
	user, err := userService.UpdateProfile(userID.(uint), updateReq)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

// GetProducts 获取产品列表
func GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 这里应该从数据库获取产品列表
	products := []gin.H{
		{
			"id":    1,
			"name":  "Product 1",
			"price": 99.99,
		},
		{
			"id":    2,
			"name":  "Product 2",
			"price": 149.99,
		},
	}

	response.Success(c, gin.H{
		"products": products,
		"page":     page,
		"limit":    limit,
		"total":    len(products),
	})
}

// GetProduct 获取单个产品
func GetProduct(c *gin.Context) {
	id := c.Param("id")

	// 这里应该从数据库获取产品信息
	response.Success(c, gin.H{
		"id":    id,
		"name":  "Product " + id,
		"price": 99.99,
	})
}

// CreateProduct 创建产品
func CreateProduct(c *gin.Context) {
	_ = c.GetString("user_id") // 获取用户ID但不使用，实际项目中应该用于权限验证

	var req struct {
		Name  string  `json:"name" binding:"required"`
		Price float64 `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	// 这里应该创建产品记录
	response.Success(c, gin.H{
		"message": "Product created successfully",
		"product": gin.H{
			"id":    1,
			"name":  req.Name,
			"price": req.Price,
		},
	})
}

// UpdateProduct 更新产品
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	_ = c.GetString("user_id") // 获取用户ID但不使用，实际项目中应该用于权限验证

	var req struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	// 这里应该更新产品记录
	response.Success(c, gin.H{
		"message": "Product updated successfully",
		"id":      id,
	})
}

// DeleteProduct 删除产品
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	_ = c.GetString("user_id") // 获取用户ID但不使用，实际项目中应该用于权限验证

	// 这里应该删除产品记录
	response.Success(c, gin.H{
		"message": "Product deleted successfully",
		"id":      id,
	})
}

// 通知相关函数
var notificationHandler = NewNotificationHandler()

// 文章相关函数
var articleHandler *ArticleHandler

func getArticleHandler() *ArticleHandler {
	if articleHandler == nil {
		articleHandler = NewArticleHandler(service.NewArticleService(database.GetDB()))
	}
	return articleHandler
}

// GetNotifications 获取用户通知列表
func GetNotifications(c *gin.Context) {
	notificationHandler.GetNotifications(c)
}

// GetUnreadCount 获取未读通知数量
func GetUnreadCount(c *gin.Context) {
	notificationHandler.GetUnreadCount(c)
}

// MarkAsRead 标记通知为已读
func MarkAsRead(c *gin.Context) {
	notificationHandler.MarkAsRead(c)
}

// MarkAllAsRead 标记所有通知为已读
func MarkAllAsRead(c *gin.Context) {
	notificationHandler.MarkAllAsRead(c)
}

// DeleteNotification 删除通知
func DeleteNotification(c *gin.Context) {
	notificationHandler.DeleteNotification(c)
}

// CheckNotifications 手动触发通知检查
func CheckNotifications(c *gin.Context) {
	notificationHandler.CheckNotifications(c)
}

// CreateArticle 创建文章
func CreateArticle(c *gin.Context) {
	getArticleHandler().CreateArticle(c)
}

// GetUserArticles 获取用户文章列表
func GetUserArticles(c *gin.Context) {
	getArticleHandler().GetUserArticles(c)
}

// GetArticleByID 根据ID获取文章
func GetArticleByID(c *gin.Context) {
	getArticleHandler().GetArticleByID(c)
}

// UpdateArticle 更新文章
func UpdateArticle(c *gin.Context) {
	getArticleHandler().UpdateArticle(c)
}

// DeleteArticle 删除文章
func DeleteArticle(c *gin.Context) {
	getArticleHandler().DeleteArticle(c)
}

// GetArticleStats 获取文章统计信息
func GetArticleStats(c *gin.Context) {
	getArticleHandler().GetArticleStats(c)
}

// 上传相关函数
var uploadHandler = NewUploadHandler()

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
	uploadHandler.UploadImage(c)
}

// 统计相关函数
var statisticsHandler = NewStatisticsHandler()

// GetStatistics 获取统计数据
func GetStatistics(c *gin.Context) {
	statisticsHandler.GetStatistics(c)
}

// GetTrends 获取趋势数据
func GetTrends(c *gin.Context) {
	statisticsHandler.GetTrends(c)
}

// 分类相关函数
var categoryHandler = NewCategoryHandler()

// CreateCategory 创建分类
func CreateCategory(c *gin.Context) {
	categoryHandler.CreateCategory(c)
}

// GetCategoryTree 获取分类树
func GetCategoryTree(c *gin.Context) {
	categoryHandler.GetCategoryTree(c)
}

// GetAllCategories 获取所有分类
func GetAllCategories(c *gin.Context) {
	categoryHandler.GetAllCategories(c)
}

// GetCategoryByID 根据ID获取分类
func GetCategoryByID(c *gin.Context) {
	categoryHandler.GetCategoryByID(c)
}

// UpdateCategory 更新分类
func UpdateCategory(c *gin.Context) {
	categoryHandler.UpdateCategory(c)
}

// DeleteCategory 删除分类
func DeleteCategory(c *gin.Context) {
	categoryHandler.DeleteCategory(c)
}

// GetCategoryStats 获取分类统计
func GetCategoryStats(c *gin.Context) {
	categoryHandler.GetCategoryStats(c)
}
