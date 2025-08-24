package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// StringUtils 字符串工具集合
type StringUtils struct{}

// NewStringUtils 创建字符串工具实例
func NewStringUtils() *StringUtils {
	return &StringUtils{}
}

// IsEmpty 检查字符串是否为空
func (s *StringUtils) IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// IsNotEmpty 检查字符串是否不为空
func (s *StringUtils) IsNotEmpty(str string) bool {
	return !s.IsEmpty(str)
}

// TrimToEmpty 修剪字符串，如果为nil则返回空字符串
func (s *StringUtils) TrimToEmpty(str *string) string {
	if str == nil {
		return ""
	}
	return strings.TrimSpace(*str)
}

// DefaultIfEmpty 如果字符串为空则返回默认值
func (s *StringUtils) DefaultIfEmpty(str, defaultStr string) string {
	if s.IsEmpty(str) {
		return defaultStr
	}
	return str
}

// CamelCase 转换为驼峰命名
func (s *StringUtils) CamelCase(str string) string {
	if str == "" {
		return ""
	}
	
	// 按分隔符分割
	words := regexp.MustCompile(`[^a-zA-Z0-9]+`).Split(str, -1)
	var result strings.Builder
	
	for i, word := range words {
		if word == "" {
			continue
		}
		
		if i == 0 {
			result.WriteString(strings.ToLower(word))
		} else {
			result.WriteString(s.Capitalize(word))
		}
	}
	
	return result.String()
}

// PascalCase 转换为帕斯卡命名（首字母大写的驼峰）
func (s *StringUtils) PascalCase(str string) string {
	camel := s.CamelCase(str)
	if camel == "" {
		return ""
	}
	return s.Capitalize(camel)
}

// SnakeCase 转换为蛇形命名
func (s *StringUtils) SnakeCase(str string) string {
	if str == "" {
		return ""
	}
	
	// 处理驼峰命名
	re1 := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	str = re1.ReplaceAllString(str, "${1}_${2}")
	
	// 处理连续大写字母
	re2 := regexp.MustCompile(`([A-Z])([A-Z][a-z])`)
	str = re2.ReplaceAllString(str, "${1}_${2}")
	
	// 替换非字母数字字符为下划线
	re3 := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	str = re3.ReplaceAllString(str, "_")
	
	// 移除首尾下划线
	str = strings.Trim(str, "_")
	
	return strings.ToLower(str)
}

// KebabCase 转换为短横线命名
func (s *StringUtils) KebabCase(str string) string {
	snake := s.SnakeCase(str)
	return strings.ReplaceAll(snake, "_", "-")
}

// Capitalize 首字母大写
func (s *StringUtils) Capitalize(str string) string {
	if str == "" {
		return ""
	}
	
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Uncapitalize 首字母小写
func (s *StringUtils) Uncapitalize(str string) string {
	if str == "" {
		return ""
	}
	
	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// Reverse 反转字符串
func (s *StringUtils) Reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Truncate 截断字符串
func (s *StringUtils) Truncate(str string, length int, suffix string) string {
	if len(str) <= length {
		return str
	}
	
	if suffix == "" {
		suffix = "..."
	}
	
	return str[:length-len(suffix)] + suffix
}

// Contains 检查字符串是否包含子字符串（忽略大小写）
func (s *StringUtils) ContainsIgnoreCase(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

// RemoveAll 移除所有指定的字符
func (s *StringUtils) RemoveAll(str, remove string) string {
	return strings.ReplaceAll(str, remove, "")
}

// RemovePrefix 移除前缀
func (s *StringUtils) RemovePrefix(str, prefix string) string {
	return strings.TrimPrefix(str, prefix)
}

// RemoveSuffix 移除后缀
func (s *StringUtils) RemoveSuffix(str, suffix string) string {
	return strings.TrimSuffix(str, suffix)
}

// PadLeft 左侧填充
func (s *StringUtils) PadLeft(str string, length int, pad string) string {
	if len(str) >= length {
		return str
	}
	
	if pad == "" {
		pad = " "
	}
	
	padding := strings.Repeat(pad, (length-len(str))/len(pad)+1)
	return padding[:length-len(str)] + str
}

// PadRight 右侧填充
func (s *StringUtils) PadRight(str string, length int, pad string) string {
	if len(str) >= length {
		return str
	}
	
	if pad == "" {
		pad = " "
	}
	
	padding := strings.Repeat(pad, (length-len(str))/len(pad)+1)
	return str + padding[:length-len(str)]
}

// CountWords 统计单词数量
func (s *StringUtils) CountWords(str string) int {
	words := regexp.MustCompile(`\S+`).FindAllString(str, -1)
	return len(words)
}

// RandomString 生成随机字符串
func (s *StringUtils) RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		// 如果加密随机失败，回退到简单方法
		return s.randomStringFallback(length)
	}
	
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	
	return string(b)
}

// randomStringFallback 随机字符串回退方法
func (s *StringUtils) randomStringFallback(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}

// RandomAlphabetic 生成随机字母字符串
func (s *StringUtils) RandomAlphabetic(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	rand.Read(b)
	
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	
	return string(b)
}

// RandomNumeric 生成随机数字字符串
func (s *StringUtils) RandomNumeric(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	rand.Read(b)
	
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	
	return string(b)
}

// Hash 计算字符串的SHA256哈希值
func (s *StringUtils) Hash(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}

// IsNumeric 检查字符串是否为数字
func (s *StringUtils) IsNumeric(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}

// IsAlpha 检查字符串是否只包含字母
func (s *StringUtils) IsAlpha(str string) bool {
	for _, r := range str {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return str != ""
}

// IsAlphanumeric 检查字符串是否只包含字母和数字
func (s *StringUtils) IsAlphanumeric(str string) bool {
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return str != ""
}

// SplitAndTrim 分割字符串并修剪每个部分
func (s *StringUtils) SplitAndTrim(str, sep string) []string {
	parts := strings.Split(str, sep)
	var result []string
	
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	
	return result
}

// Quote 为字符串添加引号
func (s *StringUtils) Quote(str string) string {
	return fmt.Sprintf(`"%s"`, str)
}

// Unquote 移除字符串的引号
func (s *StringUtils) Unquote(str string) string {
	if len(str) >= 2 {
		if (str[0] == '"' && str[len(str)-1] == '"') ||
			(str[0] == '\'' && str[len(str)-1] == '\'') {
			return str[1 : len(str)-1]
		}
	}
	return str
}

// 全局字符串工具实例
var Str = NewStringUtils()