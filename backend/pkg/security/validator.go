package security

import (
	"context"
	"regexp"
	"strings"
	"unicode"

	"gin-web-framework/pkg/errors"
)

// InputValidator 输入验证器
type InputValidator struct {
	rules map[string][]ValidationRule
}

// ValidationRule 验证规则接口
type ValidationRule interface {
	Validate(ctx context.Context, value interface{}) error
	GetErrorMessage() string
}

// NewInputValidator 创建输入验证器
func NewInputValidator() *InputValidator {
	return &InputValidator{
		rules: make(map[string][]ValidationRule),
	}
}

// AddRule 添加验证规则
func (v *InputValidator) AddRule(field string, rule ValidationRule) {
	if v.rules[field] == nil {
		v.rules[field] = make([]ValidationRule, 0)
	}
	v.rules[field] = append(v.rules[field], rule)
}

// Validate 验证输入
func (v *InputValidator) Validate(ctx context.Context, data map[string]interface{}) error {
	var validationErrors []ValidationError

	for field, rules := range v.rules {
		value, exists := data[field]
		
		// 检查必填字段
		if !exists {
			for _, rule := range rules {
				if _, ok := rule.(*RequiredRule); ok {
					validationErrors = append(validationErrors, ValidationError{
						Field:   field,
						Message: rule.GetErrorMessage(),
					})
					break
				}
			}
			continue
		}

		// 应用所有规则
		for _, rule := range rules {
			if err := rule.Validate(ctx, value); err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field:   field,
					Message: rule.GetErrorMessage(),
				})
			}
		}
	}

	if len(validationErrors) > 0 {
		return errors.NewValidationError("输入验证失败", validationErrors)
	}

	return nil
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// 常用验证规则实现

// RequiredRule 必填规则
type RequiredRule struct {
	message string
}

func NewRequiredRule(message string) *RequiredRule {
	if message == "" {
		message = "此字段为必填项"
	}
	return &RequiredRule{message: message}
}

func (r *RequiredRule) Validate(ctx context.Context, value interface{}) error {
	if value == nil {
		return errors.NewValidationError(r.message, nil)
	}
	
	if str, ok := value.(string); ok && strings.TrimSpace(str) == "" {
		return errors.NewValidationError(r.message, nil)
	}
	
	return nil
}

func (r *RequiredRule) GetErrorMessage() string {
	return r.message
}

// StringLengthRule 字符串长度规则
type StringLengthRule struct {
	minLength int
	maxLength int
	message   string
}

func NewStringLengthRule(minLength, maxLength int, message string) *StringLengthRule {
	if message == "" {
		message = "字符串长度不符合要求"
	}
	return &StringLengthRule{
		minLength: minLength,
		maxLength: maxLength,
		message:   message,
	}
}

func (r *StringLengthRule) Validate(ctx context.Context, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return nil // 不是字符串，跳过验证
	}

	length := len([]rune(str)) // 正确处理Unicode字符
	if length < r.minLength || length > r.maxLength {
		return errors.NewValidationError(r.message, nil)
	}
	
	return nil
}

func (r *StringLengthRule) GetErrorMessage() string {
	return r.message
}

// EmailRule 邮箱验证规则
type EmailRule struct {
	message string
	regex   *regexp.Regexp
}

func NewEmailRule(message string) *EmailRule {
	if message == "" {
		message = "邮箱格式不正确"
	}
	return &EmailRule{
		message: message,
		regex:   regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (r *EmailRule) Validate(ctx context.Context, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return nil
	}

	if !r.regex.MatchString(str) {
		return errors.NewValidationError(r.message, nil)
	}
	
	return nil
}

func (r *EmailRule) GetErrorMessage() string {
	return r.message
}

// PasswordStrengthRule 密码强度规则
type PasswordStrengthRule struct {
	minLength      int
	requireUpper   bool
	requireLower   bool
	requireNumber  bool
	requireSpecial bool
	message        string
}

func NewPasswordStrengthRule(minLength int, requireUpper, requireLower, requireNumber, requireSpecial bool, message string) *PasswordStrengthRule {
	if message == "" {
		message = "密码强度不符合要求"
	}
	return &PasswordStrengthRule{
		minLength:      minLength,
		requireUpper:   requireUpper,
		requireLower:   requireLower,
		requireNumber:  requireNumber,
		requireSpecial: requireSpecial,
		message:        message,
	}
}

func (r *PasswordStrengthRule) Validate(ctx context.Context, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return nil
	}

	if len(str) < r.minLength {
		return errors.NewValidationError(r.message, nil)
	}

	hasUpper, hasLower, hasNumber, hasSpecial := false, false, false, false
	
	for _, char := range str {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsNumber(char) {
			hasNumber = true
		} else if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
	}

	if r.requireUpper && !hasUpper {
		return errors.NewValidationError(r.message, nil)
	}
	if r.requireLower && !hasLower {
		return errors.NewValidationError(r.message, nil)
	}
	if r.requireNumber && !hasNumber {
		return errors.NewValidationError(r.message, nil)
	}
	if r.requireSpecial && !hasSpecial {
		return errors.NewValidationError(r.message, nil)
	}

	return nil
}

func (r *PasswordStrengthRule) GetErrorMessage() string {
	return r.message
}

// SQLInjectionRule SQL注入防护规则
type SQLInjectionRule struct {
	message    string
	patterns   []*regexp.Regexp
	dangerous  []string
}

func NewSQLInjectionRule(message string) *SQLInjectionRule {
	if message == "" {
		message = "输入包含危险内容"
	}

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)`),
		regexp.MustCompile(`(?i)(script|javascript|vbscript|onload|onerror|onclick)`),
		regexp.MustCompile(`(?i)(or\s+1\s*=\s*1|and\s+1\s*=\s*1)`),
		regexp.MustCompile(`(?i)([\'\"];?\s*(or|and)\s+[\'\"])`),
	}

	dangerous := []string{
		"--", "/*", "*/", "@@", "@",
		"char", "nchar", "varchar", "nvarchar",
		"declare", "cast", "convert", "cursor",
	}

	return &SQLInjectionRule{
		message:   message,
		patterns:  patterns,
		dangerous: dangerous,
	}
}

func (r *SQLInjectionRule) Validate(ctx context.Context, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return nil
	}

	// 检查正则模式
	for _, pattern := range r.patterns {
		if pattern.MatchString(str) {
			return errors.NewValidationError(r.message, nil)
		}
	}

	// 检查危险关键词
	lowerStr := strings.ToLower(str)
	for _, dangerous := range r.dangerous {
		if strings.Contains(lowerStr, dangerous) {
			return errors.NewValidationError(r.message, nil)
		}
	}

	return nil
}

func (r *SQLInjectionRule) GetErrorMessage() string {
	return r.message
}

// XSSRule XSS防护规则
type XSSRule struct {
	message  string
	patterns []*regexp.Regexp
}

func NewXSSRule(message string) *XSSRule {
	if message == "" {
		message = "输入包含恶意脚本"
	}

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`),
		regexp.MustCompile(`(?i)javascript:`),
		regexp.MustCompile(`(?i)on\w+\s*=`),
		regexp.MustCompile(`(?i)<iframe[^>]*>`),
		regexp.MustCompile(`(?i)<object[^>]*>`),
		regexp.MustCompile(`(?i)<embed[^>]*>`),
		regexp.MustCompile(`(?i)<link[^>]*>`),
		regexp.MustCompile(`(?i)<meta[^>]*>`),
	}

	return &XSSRule{
		message:  message,
		patterns: patterns,
	}
}

func (r *XSSRule) Validate(ctx context.Context, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return nil
	}

	for _, pattern := range r.patterns {
		if pattern.MatchString(str) {
			return errors.NewValidationError(r.message, nil)
		}
	}

	return nil
}

func (r *XSSRule) GetErrorMessage() string {
	return r.message
}

// RegexRule 正则表达式规则
type RegexRule struct {
	pattern *regexp.Regexp
	message string
}

func NewRegexRule(pattern, message string) (*RegexRule, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	if message == "" {
		message = "格式不正确"
	}

	return &RegexRule{
		pattern: regex,
		message: message,
	}, nil
}

func (r *RegexRule) Validate(ctx context.Context, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return nil
	}

	if !r.pattern.MatchString(str) {
		return errors.NewValidationError(r.message, nil)
	}

	return nil
}

func (r *RegexRule) GetErrorMessage() string {
	return r.message
}

// RangeRule 数值范围规则
type RangeRule struct {
	min     float64
	max     float64
	message string
}

func NewRangeRule(min, max float64, message string) *RangeRule {
	if message == "" {
		message = "数值超出范围"
	}
	return &RangeRule{
		min:     min,
		max:     max,
		message: message,
	}
}

func (r *RangeRule) Validate(ctx context.Context, value interface{}) error {
	var num float64
	var ok bool

	switch v := value.(type) {
	case int:
		num, ok = float64(v), true
	case int64:
		num, ok = float64(v), true
	case float64:
		num, ok = v, true
	case float32:
		num, ok = float64(v), true
	}

	if !ok {
		return nil
	}

	if num < r.min || num > r.max {
		return errors.NewValidationError(r.message, nil)
	}

	return nil
}

func (r *RangeRule) GetErrorMessage() string {
	return r.message
}