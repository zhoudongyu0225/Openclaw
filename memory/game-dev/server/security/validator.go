// Package security provides security utilities including input validation
package security

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	// MaxTextLength 最大文本长度
	MaxTextLength = 1000
	// MaxNameLength 最大名称长度
	MaxNameLength = 50
	// MaxPasswordLength 最大密码长度
	MaxPasswordLength = 128
	// MinPasswordLength 最小密码长度
	MinPasswordLength = 6
)

// InputValidator 输入验证器
type InputValidator struct {
	// 用户名正则
	usernameRegex *regexp.Regexp
	// 房间名正则
	roomNameRegex *regexp.Regexp
	// 弹幕正则
	danmakuRegex *regexp.Regexp
	// 敏感词列表
	sensitiveWords []string
}

// NewInputValidator 创建输入验证器
func NewInputValidator() *InputValidator {
	return &InputValidator{
		usernameRegex:  regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`),
		roomNameRegex:  regexp.MustCompile(`^[\p{L}\p{N}\s_\-]{2,30}$`),
		danmakuRegex:   regexp.MustCompile(`^[\p{L}\p{N}\p{P}\p{Z}\s]{1,200}$`),
		sensitiveWords: getDefaultSensitiveWords(),
	}
}

// ValidateUsername 验证用户名
func (v *InputValidator) ValidateUsername(username string) ValidationResult {
	if username == "" {
		return ValidationResult{Valid: false, Error: "用户名不能为空"}
	}
	
	if len(username) < 3 || len(username) > 20 {
		return ValidationResult{Valid: false, Error: "用户名长度需在3-20个字符之间"}
	}
	
	if !v.usernameRegex.MatchString(username) {
		return ValidationResult{Valid: false, Error: "用户名只能包含字母、数字和下划线"}
	}
	
	// 检查敏感词
	if v.containsSensitiveWord(username) {
		return ValidationResult{Valid: false, Error: "用户名包含敏感词"}
	}
	
	return ValidationResult{Valid: true}
}

// ValidatePassword 验证密码
func (v *InputValidator) ValidatePassword(password string) ValidationResult {
	if password == "" {
		return ValidationResult{Valid: false, Error: "密码不能为空"}
	}
	
	if len(password) < MinPasswordLength || len(password) > MaxPasswordLength {
		return ValidationResult{Valid: false, Error: "密码长度需在6-128个字符之间"}
	}
	
	return ValidationResult{Valid: true}
}

// ValidateRoomName 验证房间名
func (v *InputValidator) ValidateRoomName(name string) ValidationResult {
	if name == "" {
		return ValidationResult{Valid: false, Error: "房间名不能为空"}
	}
	
	if len(name) < 2 || len(name) > 30 {
		return ValidationResult{Valid: false, Error: "房间名长度需在2-30个字符之间"}
	}
	
	if !v.roomNameRegex.MatchString(name) {
		return ValidationResult{Valid: false, Error: "房间名包含非法字符"}
	}
	
	if v.containsSensitiveWord(name) {
		return ValidationResult{Valid: false, Error: "房间名包含敏感词"}
	}
	
	return ValidationResult{Valid: true}
}

// ValidateDanmaku 验证弹幕内容
func (v *InputValidator) ValidateDanmaku(content string) ValidationResult {
	if content == "" {
		return ValidationResult{Valid: false, Error: "弹幕内容不能为空"}
	}
	
	// 检查长度（UTF-8编码）
	if utf8.RuneCountInString(content) > 200 {
		return ValidationResult{Valid: false, Error: "弹幕长度不能超过200个字符"}
	}
	
	if !v.danmakuRegex.MatchString(content) {
		return ValidationResult{Valid: false, Error: "弹幕包含非法字符"}
	}
	
	// 检查敏感词
	if v.containsSensitiveWord(content) {
		return ValidationResult{Valid: false, Error: "弹幕包含敏感词"}
	}
	
	return ValidationResult{Valid: true}
}

// ValidateGift 验证礼物数据
func (v *InputValidator) ValidateGift(giftID string, amount int) ValidationResult {
	if giftID == "" {
		return ValidationResult{Valid: false, Error: "礼物ID不能为空"}
	}
	
	if amount <= 0 || amount > 9999 {
		return ValidationResult{Valid: false, Error: "礼物数量需在1-9999之间"}
	}
	
	return ValidationResult{Valid: true}
}

// ValidatePlayerInput 验证玩家输入（通用）
func (v *InputValidator) ValidatePlayerInput(input string, maxLen int) ValidationResult {
	if input == "" {
		return ValidationResult{Valid: false, Error: "输入不能为空"}
	}
	
	if utf8.RuneCountInString(input) > maxLen {
		return ValidationResult{Valid: false, Error: "输入长度超出限制"}
	}
	
	// 检查空字符
	if strings.Contains(input, "\x00") {
		return ValidationResult{Valid: false, Error: "输入包含非法字符"}
	}
	
	return ValidationResult{Valid: true}
}

// containsSensitiveWord 检查是否包含敏感词
func (v *InputValidator) containsSensitiveWord(text string) string {
	lower := strings.ToLower(text)
	for _, word := range v.sensitiveWords {
		if strings.Contains(lower, word) {
			return word
		}
	}
	return ""
}

// FilterSensitiveWords 过滤敏感词（替换为*）
func (v *InputValidator) FilterSensitiveWords(text string) string {
	result := text
	for _, word := range v.sensitiveWords {
		if strings.Contains(strings.ToLower(result), word) {
			result = strings.ReplaceAll(result, word, strings.Repeat("*", len(word)))
		}
	}
	return result
}

// SanitizeInput 清理输入（移除危险字符）
func (v *InputValidator) SanitizeInput(input string) string {
	// 移除非打印字符
	var builder strings.Builder
	for _, r := range input {
		if r == '\t' || r == '\n' || r == '\r' {
			builder.WriteRune(' ')
		} else if r >= 32 && r < 127 {
			builder.WriteRune(r)
		} else if unicode.IsPrint(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

// ValidationResult 验证结果
type ValidationResult struct {
	Valid bool
	Error string
}

// getDefaultSensitiveWords 获取默认敏感词列表
func getDefaultSensitiveWords() []string {
	return []string{
		// 色情相关
		"av", "a片", "成人", "激情", "黄色",
		// 赌博相关
		"赌博", "赌场", "博彩", "老虎机",
		// 诈骗相关
		"诈骗", "钓鱼", "木马", "病毒",
		// 政治相关（简化示例）
		"敏感", "违规",
		// 其他
		"黑客", "外挂", "作弊",
	}
}

// --- XSS 防护 ---

// XSSSanitizer XSS防护器
type XSSSanitizer struct {
	// 危险标签
	dangerousTags *regexp.Regexp
	// 危险属性
	dangerousAttrs *regexp.Regexp
}

// NewXSSSanitizer 创建XSS防护器
func NewXSSSanitizer() *XSSSanitizer {
	return &XSSSanitizer{
		dangerousTags:  regexp.MustCompile(`(?i)<(script|iframe|object|embed|form|input|button|select|textarea|style)`),
		dangerousAttrs: regexp.MustCompile(`(?i)\s+(onclick|onload|onerror|onmouse|onkey|onfocus|onblur|onsubmit)`),
	}
}

// Sanitize 清理HTML/XSS
func (x *XSSSanitizer) Sanitize(html string) string {
	// 移除危险标签
	result := x.dangerousTags.ReplaceAllString(html, "")
	// 移除危险属性
	result = x.dangerousAttrs.ReplaceAllString(result, "")
	// 移除javascript: 协议
	result = strings.ReplaceAll(result, "javascript:", "")
	// 移除data: 协议（可能用于base64攻击）
	result = strings.ReplaceAll(result, "data:", "")
	return result
}

// --- SQL 注入防护（MongoDB） ---

// MongoDBSanitizer MongoDB查询清理
type MongoDBSanitizer struct {
	dangerousRegex *regexp.Regexp
}

// NewMongoDBSanitizer 创建MongoDB清理器
func NewMongoDBSanitizer() *MongoDBSanitizer {
	return &MongoDBSanitizer{
		// 防止 $ 开头的MongoDB操作符注入
		dangerousRegex: regexp.MustCompile(`[\$\{\}]`),
	}
}

// SanitizeFieldName 清理字段名
func (m *MongoDBSanitizer) SanitizeFieldName(field string) string {
	// 字段名不能包含 $ 和 .
	if m.dangerousRegex.MatchString(field) {
		return ""
	}
	return field
}

// SanitizeValue 清理值
func (m *MongoDBSanitizer) SanitizeValue(value string) string {
	// 转义 $ 开头的内容（可能被解释为操作符）
	if strings.HasPrefix(value, "$") {
		return "\\" + value
	}
	return value
}
