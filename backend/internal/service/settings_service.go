package service

import (
	"errors"

	"gin-web-framework/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SettingsService struct {
	db *gorm.DB
}

func NewSettingsService(db *gorm.DB) *SettingsService {
	return &SettingsService{db: db}
}

func (s *SettingsService) GetUserSettings(userID uint) (*models.UserSettings, error) {
	var settings models.UserSettings
	err := s.db.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建默认设置
			settings = models.UserSettings{
				UserID:                 userID,
				DueReminder:            true,
				CompletionNotification: true,
				NewTaskNotification:    true,
				EmailNotification:      false,
				Theme:                  "light",
				Language:               "zh-CN",
				Timezone:               "Asia/Shanghai",
			}
			if err := s.db.Create(&settings).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &settings, nil
}

func (s *SettingsService) UpdateProfile(userID uint, req *models.UpdateProfileRequest) error {
	// 检查用户名是否已存在
	var existingUser models.User
	err := s.db.Where("username = ? AND id != ?", req.Username, userID).First(&existingUser).Error
	if err == nil {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	err = s.db.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error
	if err == nil {
		return errors.New("邮箱已存在")
	}

	// 更新用户信息
	updates := models.User{
		Username: req.Username,
		Email:    req.Email,
		Nickname: req.Nickname,
	}
	
	return s.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (s *SettingsService) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证当前密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("当前密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	return s.db.Model(&models.User{}).Where("id = ?", userID).Update("password", string(hashedPassword)).Error
}

func (s *SettingsService) UpdateNotificationSettings(userID uint, req *models.UpdateNotificationSettingsRequest) error {
	settings, err := s.GetUserSettings(userID)
	if err != nil {
		return err
	}

	updates := models.UserSettings{
		DueReminder:            req.DueReminder,
		CompletionNotification: req.CompletionNotification,
		NewTaskNotification:    req.NewTaskNotification,
		EmailNotification:      req.EmailNotification,
	}

	return s.db.Model(settings).Updates(updates).Error
}

func (s *SettingsService) UpdateInterfaceSettings(userID uint, req *models.UpdateInterfaceSettingsRequest) error {
	settings, err := s.GetUserSettings(userID)
	if err != nil {
		return err
	}

	updates := models.UserSettings{
		Theme:    req.Theme,
		Language: req.Language,
		Timezone: req.Timezone,
	}

	return s.db.Model(settings).Updates(updates).Error
}

func (s *SettingsService) GetUserData(userID uint) (map[string]interface{}, error) {
	var todos []models.Todo
	if err := s.db.Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		return nil, err
	}

	var articles []models.Article
	if err := s.db.Where("user_id = ?", userID).Find(&articles).Error; err != nil {
		return nil, err
	}

	var categories []models.Category
	if err := s.db.Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		return nil, err
	}

	settings, err := s.GetUserSettings(userID)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"todos":      todos,
		"articles":   articles,
		"categories": categories,
		"settings":   settings,
		"exportTime": "",
	}

	return data, nil
}

func (s *SettingsService) ClearCompletedTasks(userID uint) error {
	return s.db.Where("user_id = ? AND completed = ?", userID, true).Delete(&models.Todo{}).Error
}