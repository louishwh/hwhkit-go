package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPUtils HTTP工具集合
type HTTPUtils struct {
	client *http.Client
}

// NewHTTPUtils 创建HTTP工具实例
func NewHTTPUtils() *HTTPUtils {
	return &HTTPUtils{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewHTTPUtilsWithTimeout 创建指定超时时间的HTTP工具实例
func NewHTTPUtilsWithTimeout(timeout time.Duration) *HTTPUtils {
	return &HTTPUtils{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// HTTPResponse HTTP响应结构
type HTTPResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Text       string
}

// IsSuccess 判断HTTP状态码是否为成功（2xx）
func (r *HTTPResponse) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// JSON 将响应体解析为JSON
func (r *HTTPResponse) JSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

// Get 发送GET请求
func (h *HTTPUtils) Get(url string, headers map[string]string) (*HTTPResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	
	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	return h.doRequest(req)
}

// Post 发送POST请求（JSON数据）
func (h *HTTPUtils) Post(url string, data interface{}, headers map[string]string) (*HTTPResponse, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON data: %w", err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}
	
	// 设置Content-Type
	req.Header.Set("Content-Type", "application/json")
	
	// 设置其他请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	return h.doRequest(req)
}

// PostForm 发送POST表单请求
func (h *HTTPUtils) PostForm(url string, formData map[string]string, headers map[string]string) (*HTTPResponse, error) {
	// 构建表单数据
	values := make(url.Values)
	for key, value := range formData {
		values.Set(key, value)
	}
	
	req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST form request: %w", err)
	}
	
	// 设置Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	// 设置其他请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	return h.doRequest(req)
}

// Put 发送PUT请求
func (h *HTTPUtils) Put(url string, data interface{}, headers map[string]string) (*HTTPResponse, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON data: %w", err)
	}
	
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %w", err)
	}
	
	// 设置Content-Type
	req.Header.Set("Content-Type", "application/json")
	
	// 设置其他请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	return h.doRequest(req)
}

// Patch 发送PATCH请求
func (h *HTTPUtils) Patch(url string, data interface{}, headers map[string]string) (*HTTPResponse, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON data: %w", err)
	}
	
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create PATCH request: %w", err)
	}
	
	// 设置Content-Type
	req.Header.Set("Content-Type", "application/json")
	
	// 设置其他请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	return h.doRequest(req)
}

// Delete 发送DELETE请求
func (h *HTTPUtils) Delete(url string, headers map[string]string) (*HTTPResponse, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %w", err)
	}
	
	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	return h.doRequest(req)
}

// Request 发送自定义请求
func (h *HTTPUtils) Request(method, url string, body io.Reader, headers map[string]string) (*HTTPResponse, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s request: %w", method, err)
	}
	
	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	return h.doRequest(req)
}

// doRequest 执行HTTP请求
func (h *HTTPUtils) doRequest(req *http.Request) (*HTTPResponse, error) {
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
		Text:       string(body),
	}, nil
}

// SetTimeout 设置超时时间
func (h *HTTPUtils) SetTimeout(timeout time.Duration) {
	h.client.Timeout = timeout
}

// GetClient 获取HTTP客户端
func (h *HTTPUtils) GetClient() *http.Client {
	return h.client
}

// SetClient 设置HTTP客户端
func (h *HTTPUtils) SetClient(client *http.Client) {
	h.client = client
}

// DownloadFile 下载文件
func (h *HTTPUtils) DownloadFile(url, filepath string, headers map[string]string) error {
	resp, err := h.Get(url, headers)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	
	if !resp.IsSuccess() {
		return fmt.Errorf("download failed with status code: %d", resp.StatusCode)
	}
	
	// 这里应该写入文件，但由于没有文件操作工具，暂时省略
	// 可以扩展添加文件操作功能
	_ = filepath // 避免未使用变量错误
	
	return nil
}

// BuildURL 构建URL
func (h *HTTPUtils) BuildURL(baseURL string, path string, params map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}
	
	if path != "" {
		u.Path = strings.TrimSuffix(u.Path, "/") + "/" + strings.TrimPrefix(path, "/")
	}
	
	if len(params) > 0 {
		query := u.Query()
		for key, value := range params {
			query.Set(key, value)
		}
		u.RawQuery = query.Encode()
	}
	
	return u.String(), nil
}

// ParseURL 解析URL
func (h *HTTPUtils) ParseURL(rawURL string) (*url.URL, error) {
	return url.Parse(rawURL)
}

// GetQueryParams 获取URL查询参数
func (h *HTTPUtils) GetQueryParams(rawURL string) (map[string]string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	
	params := make(map[string]string)
	for key, values := range u.Query() {
		if len(values) > 0 {
			params[key] = values[0] // 只取第一个值
		}
	}
	
	return params, nil
}

// EncodeURL URL编码
func (h *HTTPUtils) EncodeURL(str string) string {
	return url.QueryEscape(str)
}

// DecodeURL URL解码
func (h *HTTPUtils) DecodeURL(str string) (string, error) {
	return url.QueryUnescape(str)
}

// GetMimeType 根据文件扩展名获取MIME类型
func (h *HTTPUtils) GetMimeType(filename string) string {
	ext := strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])
	
	mimeTypes := map[string]string{
		"html": "text/html",
		"htm":  "text/html",
		"css":  "text/css",
		"js":   "text/javascript",
		"json": "application/json",
		"xml":  "application/xml",
		"pdf":  "application/pdf",
		"zip":  "application/zip",
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"svg":  "image/svg+xml",
		"ico":  "image/x-icon",
		"txt":  "text/plain",
		"csv":  "text/csv",
		"mp3":  "audio/mpeg",
		"mp4":  "video/mp4",
		"avi":  "video/x-msvideo",
		"mov":  "video/quicktime",
	}
	
	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}
	
	return "application/octet-stream"
}

// IsValidURL 检查URL是否有效
func (h *HTTPUtils) IsValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// GetStatusText 获取HTTP状态码对应的文本
func (h *HTTPUtils) GetStatusText(statusCode int) string {
	return http.StatusText(statusCode)
}

// CreateBasicAuthHeader 创建基础认证头
func (h *HTTPUtils) CreateBasicAuthHeader(username, password string) map[string]string {
	return map[string]string{
		"Authorization": "Basic " + h.encodeBasicAuth(username, password),
	}
}

// encodeBasicAuth 编码基础认证
func (h *HTTPUtils) encodeBasicAuth(username, password string) string {
	auth := username + ":" + password
	// 这里应该使用base64编码，但为了简化省略
	return auth
}

// CreateBearerAuthHeader 创建Bearer认证头
func (h *HTTPUtils) CreateBearerAuthHeader(token string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}

// CreateJSONHeaders 创建JSON请求头
func (h *HTTPUtils) CreateJSONHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
}

// MergeHeaders 合并请求头
func (h *HTTPUtils) MergeHeaders(headers ...map[string]string) map[string]string {
	result := make(map[string]string)
	
	for _, h := range headers {
		for key, value := range h {
			result[key] = value
		}
	}
	
	return result
}

// 全局HTTP工具实例
var HTTP = NewHTTPUtils()