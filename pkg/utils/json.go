package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// JSONUtils JSON工具集合
type JSONUtils struct{}

// NewJSONUtils 创建JSON工具实例
func NewJSONUtils() *JSONUtils {
	return &JSONUtils{}
}

// ToJSON 将对象转换为JSON字符串
func (j *JSONUtils) ToJSON(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(bytes), nil
}

// ToJSONBytes 将对象转换为JSON字节数组
func (j *JSONUtils) ToJSONBytes(v interface{}) ([]byte, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON bytes: %w", err)
	}
	return bytes, nil
}

// ToPrettyJSON 将对象转换为格式化的JSON字符串
func (j *JSONUtils) ToPrettyJSON(v interface{}) (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to pretty JSON: %w", err)
	}
	return string(bytes), nil
}

// FromJSON 从JSON字符串解析对象
func (j *JSONUtils) FromJSON(jsonStr string, v interface{}) error {
	if err := json.Unmarshal([]byte(jsonStr), v); err != nil {
		return fmt.Errorf("failed to unmarshal from JSON: %w", err)
	}
	return nil
}

// FromJSONBytes 从JSON字节数组解析对象
func (j *JSONUtils) FromJSONBytes(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal from JSON bytes: %w", err)
	}
	return nil
}

// IsValidJSON 检查字符串是否为有效的JSON
func (j *JSONUtils) IsValidJSON(jsonStr string) bool {
	var js interface{}
	return json.Unmarshal([]byte(jsonStr), &js) == nil
}

// IsValidJSONBytes 检查字节数组是否为有效的JSON
func (j *JSONUtils) IsValidJSONBytes(data []byte) bool {
	var js interface{}
	return json.Unmarshal(data, &js) == nil
}

// Clone 通过JSON序列化/反序列化深度克隆对象
func (j *JSONUtils) Clone(src, dst interface{}) error {
	data, err := j.ToJSONBytes(src)
	if err != nil {
		return fmt.Errorf("failed to clone object: %w", err)
	}
	
	return j.FromJSONBytes(data, dst)
}

// DeepClone 深度克隆对象并返回新对象
func (j *JSONUtils) DeepClone(src interface{}) (interface{}, error) {
	// 获取源对象的类型
	srcType := reflect.TypeOf(src)
	if srcType.Kind() == reflect.Ptr {
		srcType = srcType.Elem()
	}
	
	// 创建新的实例
	dst := reflect.New(srcType).Interface()
	
	// 执行克隆
	if err := j.Clone(src, dst); err != nil {
		return nil, err
	}
	
	return dst, nil
}

// CompareJSON 比较两个JSON字符串是否相等（忽略格式）
func (j *JSONUtils) CompareJSON(json1, json2 string) (bool, error) {
	var obj1, obj2 interface{}
	
	if err := j.FromJSON(json1, &obj1); err != nil {
		return false, fmt.Errorf("failed to parse first JSON: %w", err)
	}
	
	if err := j.FromJSON(json2, &obj2); err != nil {
		return false, fmt.Errorf("failed to parse second JSON: %w", err)
	}
	
	return reflect.DeepEqual(obj1, obj2), nil
}

// ExtractField 从JSON字符串中提取指定字段
func (j *JSONUtils) ExtractField(jsonStr, fieldPath string) (interface{}, error) {
	var data map[string]interface{}
	if err := j.FromJSON(jsonStr, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	// 简单实现，只支持一级字段
	// 可以扩展支持嵌套字段，如 "user.profile.name"
	if value, exists := data[fieldPath]; exists {
		return value, nil
	}
	
	return nil, fmt.Errorf("field '%s' not found", fieldPath)
}

// SetField 在JSON字符串中设置指定字段的值
func (j *JSONUtils) SetField(jsonStr, fieldPath string, value interface{}) (string, error) {
	var data map[string]interface{}
	if err := j.FromJSON(jsonStr, &data); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	// 简单实现，只支持一级字段
	data[fieldPath] = value
	
	return j.ToJSON(data)
}

// RemoveField 从JSON字符串中移除指定字段
func (j *JSONUtils) RemoveField(jsonStr, fieldPath string) (string, error) {
	var data map[string]interface{}
	if err := j.FromJSON(jsonStr, &data); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	delete(data, fieldPath)
	
	return j.ToJSON(data)
}

// MergeJSON 合并两个JSON对象
func (j *JSONUtils) MergeJSON(json1, json2 string) (string, error) {
	var obj1, obj2 map[string]interface{}
	
	if err := j.FromJSON(json1, &obj1); err != nil {
		return "", fmt.Errorf("failed to parse first JSON: %w", err)
	}
	
	if err := j.FromJSON(json2, &obj2); err != nil {
		return "", fmt.Errorf("failed to parse second JSON: %w", err)
	}
	
	// 合并对象（obj2 的字段会覆盖 obj1 中的同名字段）
	for key, value := range obj2 {
		obj1[key] = value
	}
	
	return j.ToJSON(obj1)
}

// ConvertToMap 将任意对象转换为map[string]interface{}
func (j *JSONUtils) ConvertToMap(v interface{}) (map[string]interface{}, error) {
	// 先转换为JSON，再解析为map
	jsonStr, err := j.ToJSON(v)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	if err := j.FromJSON(jsonStr, &result); err != nil {
		return nil, err
	}
	
	return result, nil
}

// ConvertFromMap 将map[string]interface{}转换为指定类型的对象
func (j *JSONUtils) ConvertFromMap(data map[string]interface{}, v interface{}) error {
	jsonStr, err := j.ToJSON(data)
	if err != nil {
		return err
	}
	
	return j.FromJSON(jsonStr, v)
}

// GetKeys 获取JSON对象的所有键
func (j *JSONUtils) GetKeys(jsonStr string) ([]string, error) {
	var data map[string]interface{}
	if err := j.FromJSON(jsonStr, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	
	return keys, nil
}

// FilterFields 过滤JSON对象，只保留指定的字段
func (j *JSONUtils) FilterFields(jsonStr string, fields []string) (string, error) {
	var data map[string]interface{}
	if err := j.FromJSON(jsonStr, &data); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	filtered := make(map[string]interface{})
	for _, field := range fields {
		if value, exists := data[field]; exists {
			filtered[field] = value
		}
	}
	
	return j.ToJSON(filtered)
}

// ExcludeFields 排除JSON对象中的指定字段
func (j *JSONUtils) ExcludeFields(jsonStr string, fields []string) (string, error) {
	var data map[string]interface{}
	if err := j.FromJSON(jsonStr, &data); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	// 创建排除字段的集合
	excludeSet := make(map[string]bool)
	for _, field := range fields {
		excludeSet[field] = true
	}
	
	// 过滤数据
	filtered := make(map[string]interface{})
	for key, value := range data {
		if !excludeSet[key] {
			filtered[key] = value
		}
	}
	
	return j.ToJSON(filtered)
}

// 全局JSON工具实例
var JSON = NewJSONUtils()