package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringUtils_IsEmpty(t *testing.T) {
	s := NewStringUtils()
	
	assert.True(t, s.IsEmpty(""))
	assert.True(t, s.IsEmpty("   "))
	assert.True(t, s.IsEmpty("\t\n"))
	assert.False(t, s.IsEmpty("hello"))
	assert.False(t, s.IsEmpty("  hello  "))
}

func TestStringUtils_DefaultIfEmpty(t *testing.T) {
	s := NewStringUtils()
	
	assert.Equal(t, "default", s.DefaultIfEmpty("", "default"))
	assert.Equal(t, "default", s.DefaultIfEmpty("   ", "default"))
	assert.Equal(t, "hello", s.DefaultIfEmpty("hello", "default"))
}

func TestStringUtils_CamelCase(t *testing.T) {
	s := NewStringUtils()
	
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "helloWorld"},
		{"hello-world", "helloWorld"},
		{"hello world", "helloWorld"},
		{"HelloWorld", "helloWorld"},
		{"HELLO_WORLD", "helloWorld"},
		{"", ""},
		{"a", "a"},
	}
	
	for _, test := range tests {
		result := s.CamelCase(test.input)
		assert.Equal(t, test.expected, result, "CamelCase(%q) = %q, want %q", test.input, result, test.expected)
	}
}

func TestStringUtils_PascalCase(t *testing.T) {
	s := NewStringUtils()
	
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"hello-world", "HelloWorld"},
		{"hello world", "HelloWorld"},
		{"helloWorld", "HelloWorld"},
		{"", ""},
		{"a", "A"},
	}
	
	for _, test := range tests {
		result := s.PascalCase(test.input)
		assert.Equal(t, test.expected, result, "PascalCase(%q) = %q, want %q", test.input, result, test.expected)
	}
}

func TestStringUtils_SnakeCase(t *testing.T) {
	s := NewStringUtils()
	
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello_world"},
		{"helloWorld", "hello_world"},
		{"hello-world", "hello_world"},
		{"hello world", "hello_world"},
		{"XMLHttpRequest", "xml_http_request"},
		{"", ""},
		{"a", "a"},
	}
	
	for _, test := range tests {
		result := s.SnakeCase(test.input)
		assert.Equal(t, test.expected, result, "SnakeCase(%q) = %q, want %q", test.input, result, test.expected)
	}
}

func TestStringUtils_KebabCase(t *testing.T) {
	s := NewStringUtils()
	
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello-world"},
		{"helloWorld", "hello-world"},
		{"hello_world", "hello-world"},
		{"hello world", "hello-world"},
		{"", ""},
		{"a", "a"},
	}
	
	for _, test := range tests {
		result := s.KebabCase(test.input)
		assert.Equal(t, test.expected, result, "KebabCase(%q) = %q, want %q", test.input, result, test.expected)
	}
}

func TestStringUtils_Capitalize(t *testing.T) {
	s := NewStringUtils()
	
	assert.Equal(t, "Hello", s.Capitalize("hello"))
	assert.Equal(t, "Hello", s.Capitalize("Hello"))
	assert.Equal(t, "H", s.Capitalize("h"))
	assert.Equal(t, "", s.Capitalize(""))
}

func TestStringUtils_Uncapitalize(t *testing.T) {
	s := NewStringUtils()
	
	assert.Equal(t, "hello", s.Uncapitalize("Hello"))
	assert.Equal(t, "hello", s.Uncapitalize("hello"))
	assert.Equal(t, "h", s.Uncapitalize("H"))
	assert.Equal(t, "", s.Uncapitalize(""))
}

func TestStringUtils_Reverse(t *testing.T) {
	s := NewStringUtils()
	
	assert.Equal(t, "olleh", s.Reverse("hello"))
	assert.Equal(t, "a", s.Reverse("a"))
	assert.Equal(t, "", s.Reverse(""))
	assert.Equal(t, "54321", s.Reverse("12345"))
}

func TestStringUtils_Truncate(t *testing.T) {
	s := NewStringUtils()
	
	assert.Equal(t, "hello", s.Truncate("hello", 10, "..."))
	assert.Equal(t, "hel...", s.Truncate("hello world", 6, "..."))
	assert.Equal(t, "hello", s.Truncate("hello", 5, "..."))
	assert.Equal(t, "he---", s.Truncate("hello", 5, "---"))
}

func TestStringUtils_ContainsIgnoreCase(t *testing.T) {
	s := NewStringUtils()
	
	assert.True(t, s.ContainsIgnoreCase("Hello World", "hello"))
	assert.True(t, s.ContainsIgnoreCase("Hello World", "WORLD"))
	assert.True(t, s.ContainsIgnoreCase("Hello World", "lo Wo"))
	assert.False(t, s.ContainsIgnoreCase("Hello World", "xyz"))
}

func TestStringUtils_RemoveAll(t *testing.T) {
	s := NewStringUtils()
	
	assert.Equal(t, "helloworld", s.RemoveAll("hello world", " "))
	assert.Equal(t, "hello world", s.RemoveAll("hello world", "x"))
	assert.Equal(t, "", s.RemoveAll("aaaa", "a"))
}

func TestStringUtils_RemovePrefix(t *testing.T) {
	s := NewStringUtils()
	
	assert.Equal(t, "world", s.RemovePrefix("hello world", "hello "))
	assert.Equal(t, "hello world", s.RemovePrefix("hello world", "hi "))
	assert.Equal(t, "", s.RemovePrefix("hello", "hello"))
}

func TestStringUtils_RemoveSuffix(t *testing.T) {
	s := NewStringUtils()
	
	assert.Equal(t, "hello", s.RemoveSuffix("hello world", " world"))
	assert.Equal(t, "hello world", s.RemoveSuffix("hello world", " there"))
	assert.Equal(t, "", s.RemoveSuffix("hello", "hello"))
}

func TestStringUtils_PadLeft(t *testing.T) {
	s := NewStringUtils()
	
	assert.Equal(t, "  hello", s.PadLeft("hello", 7, " "))
	assert.Equal(t, "00hello", s.PadLeft("hello", 7, "0"))
	assert.Equal(t, "hello", s.PadLeft("hello", 5, " "))
	assert.Equal(t, "hello", s.PadLeft("hello", 3, " "))
}

func TestStringUtils_PadRight(t *testing.T) {
	s := NewStringUtils()
	
	assert.Equal(t, "hello  ", s.PadRight("hello", 7, " "))
	assert.Equal(t, "hello00", s.PadRight("hello", 7, "0"))
	assert.Equal(t, "hello", s.PadRight("hello", 5, " "))
	assert.Equal(t, "hello", s.PadRight("hello", 3, " "))
}

func TestStringUtils_CountWords(t *testing.T) {
	s := NewStringUtils()
	
	assert.Equal(t, 2, s.CountWords("hello world"))
	assert.Equal(t, 3, s.CountWords("hello  world  test"))
	assert.Equal(t, 0, s.CountWords(""))
	assert.Equal(t, 0, s.CountWords("   "))
	assert.Equal(t, 1, s.CountWords("hello"))
}

func TestStringUtils_RandomString(t *testing.T) {
	s := NewStringUtils()
	
	// 测试长度
	result := s.RandomString(10)
	assert.Len(t, result, 10)
	
	// 测试不同调用产生不同结果
	result1 := s.RandomString(5)
	result2 := s.RandomString(5)
	assert.NotEqual(t, result1, result2) // 理论上应该不相等
	
	// 测试空字符串
	result = s.RandomString(0)
	assert.Equal(t, "", result)
}

func TestStringUtils_IsNumeric(t *testing.T) {
	s := NewStringUtils()
	
	assert.True(t, s.IsNumeric("123"))
	assert.True(t, s.IsNumeric("123.45"))
	assert.True(t, s.IsNumeric("-123"))
	assert.True(t, s.IsNumeric("0"))
	assert.False(t, s.IsNumeric("abc"))
	assert.False(t, s.IsNumeric("123abc"))
	assert.False(t, s.IsNumeric(""))
}

func TestStringUtils_IsAlpha(t *testing.T) {
	s := NewStringUtils()
	
	assert.True(t, s.IsAlpha("abc"))
	assert.True(t, s.IsAlpha("ABC"))
	assert.True(t, s.IsAlpha("aBc"))
	assert.False(t, s.IsAlpha("abc123"))
	assert.False(t, s.IsAlpha("123"))
	assert.False(t, s.IsAlpha(""))
	assert.False(t, s.IsAlpha("abc "))
}

func TestStringUtils_IsAlphanumeric(t *testing.T) {
	s := NewStringUtils()
	
	assert.True(t, s.IsAlphanumeric("abc123"))
	assert.True(t, s.IsAlphanumeric("ABC"))
	assert.True(t, s.IsAlphanumeric("123"))
	assert.False(t, s.IsAlphanumeric("abc 123"))
	assert.False(t, s.IsAlphanumeric("abc-123"))
	assert.False(t, s.IsAlphanumeric(""))
}

func TestStringUtils_SplitAndTrim(t *testing.T) {
	s := NewStringUtils()
	
	result := s.SplitAndTrim("a, b , c,  d  ", ",")
	expected := []string{"a", "b", "c", "d"}
	assert.Equal(t, expected, result)
	
	result = s.SplitAndTrim("a;b;c", ";")
	expected = []string{"a", "b", "c"}
	assert.Equal(t, expected, result)
	
	result = s.SplitAndTrim("", ",")
	assert.Empty(t, result)
}

func TestStringUtils_Quote(t *testing.T) {
	s := NewStringUtils()
	
	assert.Equal(t, `"hello"`, s.Quote("hello"))
	assert.Equal(t, `""`, s.Quote(""))
	assert.Equal(t, `"hello world"`, s.Quote("hello world"))
}

func TestStringUtils_Unquote(t *testing.T) {
	s := NewStringUtils()
	
	assert.Equal(t, "hello", s.Unquote(`"hello"`))
	assert.Equal(t, "hello", s.Unquote("'hello'"))
	assert.Equal(t, "", s.Unquote(`""`))
	assert.Equal(t, "hello", s.Unquote("hello"))
	assert.Equal(t, `"hello`, s.Unquote(`"hello`))
}

func TestGlobalStringInstance(t *testing.T) {
	// 测试全局实例
	assert.NotNil(t, Str)
	assert.True(t, Str.IsEmpty(""))
	assert.Equal(t, "Hello", Str.Capitalize("hello"))
}