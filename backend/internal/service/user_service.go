package service

import (
	"errors"
	"fmt"

	"gin-web-framework/internal/models"
	"gin-web-framework/pkg/auth"
	pkgerrors "gin-web-framework/pkg/errors"
	"gin-web-framework/pkg/jwt"
	"gin-web-framework/pkg/logger"
	pkgdb "gin-web-framework/pkg/database"
	"gin-web-framework/pkg/utils"

	"gorm.io/gorm"
)

type UserService struct {
	db     *gorm.DB
	logger logger.LoggerInterface
}

func NewUserService(db *gorm.DB, logger logger.LoggerInterface) *UserService {
	return &UserService{
		db:     db,
		logger: logger,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type UpdateProfileRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=50"`
	Email    string `json:"email" validate:"omitempty,email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type PaginatedUsers struct {
	Users      []*models.User `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// Register 用户注册
func (s *UserService) Register(req RegisterRequest) (*models.User, error) {

	// 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		s.logger.WithFields(map[string]any{"username": req.Username}).Warn("Registration failed: username already exists")
		return nil, pkgerrors.NewConflictError("用户名已存在")
	}

	// 检查邮箱是否已存在
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		s.logger.WithFields(map[string]any{"email": req.Email}).Warn("Registration failed: email already exists")
		return nil, pkgerrors.NewConflictError("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, pkgerrors.NewInternalError("密码加密失败", err)
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.db.Create(&user).Error; err != nil {
		s.logger.WithFields(map[string]any{"username": req.Username, "error": err}).Error("Failed to create user")
		return nil, pkgerrors.NewDatabaseError("用户创建失败", err)
	}

	s.logger.WithFields(map[string]any{"user_id": user.ID, "username": user.Username}).Info("User registered successfully")
	return &user, nil
}

// Login 用户登录
func (s *UserService) Login(req LoginRequest) (*LoginResponse, error) {

	// 查找用户
	var user models.User
	if err := s.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgerrors.NewUnauthorizedError("用户名或密码错误")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	// 验证密码
	if !auth.CheckPassword(req.Password, user.Password) {
		s.logger.WithFields(map[string]any{"username": req.Username}).Warn("Login failed: invalid password")
		return nil, pkgerrors.NewUnauthorizedError("用户名或密码错误")
	}

	// 生成JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	s.logger.WithFields(map[string]any{"user_id": user.ID, "username": user.Username}).Info("User logged in successfully")
	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {

	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户（实现UserServiceInterface接口）
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {

	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &user, nil
}

// UpdateProfile 更新用户资料（实现UserServiceInterface接口）
func (s *UserService) UpdateProfile(userID uint, req UpdateProfileRequest) (*models.User, error) {

	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	// 更新字段
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return &user, nil
}

// ListUsers 获取用户列表（实现UserServiceInterface接口）
func (s *UserService) ListUsers(page, limit int) (*PaginatedUsers, error) {

	var users []*models.User
	var total int64

	// 获取总数
	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count users: %v", err)
	}

	// 获取分页数据
	pagination := utils.NewPaginationInfo(page, limit, total)
	if err := s.db.Offset(pagination.Offset).Limit(pagination.Limit).Find(&users).Error; err != nil {
		return nil, pkgdb.FindRecordError("users", err)
	}

	return &PaginatedUsers{
		Users:      users,
		Total:      pagination.Total,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalPages: pagination.TotalPages,
	}, nil
}

// UpdateUserStatus 更新用户状态（实现UserServiceInterface接口）
func (s *UserService) UpdateUserStatus(id uint, status string) error {

	result := s.db.Model(&models.User{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update user status: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// ValidateUser 验证用户（实现UserServiceInterface接口）
func (s *UserService) ValidateUser(user *models.User) error {
	if user.Username == "" {
		return errors.New("username is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// HashPassword 加密密码（实现UserServiceInterface接口）
func (s *UserService) HashPassword(password string) (string, error) {
	return auth.HashPassword(password)
}

// ChangePassword 修改密码（实现UserServiceInterface接口）
func (s *UserService) ChangePassword(userID uint, req ChangePasswordRequest) error {

	// 获取用户
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("database error: %v", err)
	}

	// 验证旧密码
	if !auth.CheckPassword(req.OldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	// 加密新密码
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// 更新密码
	user.Password = hashedPassword
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	return nil
}

// CheckPassword 检查密码（实现UserServiceInterface接口）
func (s *UserService) CheckPassword(hashedPassword, password string) bool {
	return auth.CheckPassword(password, hashedPassword)
}

// DeleteUser 删除用户（实现UserServiceInterface接口）
func (s *UserService) DeleteUser(id uint) error {

	result := s.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// GetUserByEmail 根据邮箱获取用户（实现UserServiceInterface接口）
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {

	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &user, nil
}

// RefreshToken 刷新token（实现UserServiceInterface接口）
func (s *UserService) RefreshToken(tokenString string) (*LoginResponse, error) {
	// 使用JWT包的RefreshToken方法
	newToken, err := jwt.RefreshToken(tokenString)
	if err != nil {
		s.logger.WithFields(map[string]any{"error": err}).Warn("Token refresh failed")
		return nil, fmt.Errorf("token refresh failed: %v", err)
	}

	// 解析新token获取用户信息
	claims, err := jwt.ParseToken(newToken)
	if err != nil {
		s.logger.WithFields(map[string]any{"error": err}).Error("Failed to parse new token")
		return nil, fmt.Errorf("failed to parse new token: %v", err)
	}

	// 获取用户信息
	user, err := s.GetUserByID(claims.UserID)
	if err != nil {
		s.logger.WithFields(map[string]any{"user_id": claims.UserID, "error": err}).Error("Failed to get user during token refresh")
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// 清除密码字段（安全考虑）
	user.Password = ""

	s.logger.WithFields(map[string]any{"user_id": claims.UserID, "username": claims.Username}).Info("Token refreshed successfully")

	return &LoginResponse{
		Token: newToken,
		User:  *user,
	}, nil
}
