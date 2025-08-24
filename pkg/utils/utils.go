package utils

// 这个文件导出所有工具类的实例，方便使用

// Utils 工具集合
type Utils struct {
	String *StringUtils
	JSON   *JSONUtils
	Time   *TimeUtils
	HTTP   *HTTPUtils
}

// New 创建工具集合实例
func New() *Utils {
	return &Utils{
		String: NewStringUtils(),
		JSON:   NewJSONUtils(),
		Time:   NewTimeUtils(),
		HTTP:   NewHTTPUtils(),
	}
}

// 全局工具实例，可以直接使用
var (
	// Global 全局工具实例
	Global = New()
)

// 为了方便使用，也可以直接访问各个工具实例
// 这些变量在各自的文件中已经定义：
// - Str (string.go)
// - JSON (json.go) 
// - Time (time.go)
// - HTTP (http.go)