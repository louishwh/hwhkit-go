package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hwh/hwhkit-go/pkg/config"
	"github.com/hwh/hwhkit-go/pkg/logger"
	"github.com/hwh/hwhkit-go/pkg/utils"
)

// 演示如何为应用编写测试

func TestApplicationComponents(t *testing.T) {
	// 设置测试环境
	os.Setenv("SERVER_MODE", "test")
	os.Setenv("LOG_LEVEL", "error")
	
	t.Run("Config Manager", func(t *testing.T) {
		configManager := config.New()
		cfg := configManager.Get()
		
		if cfg == nil {
			t.Fatal("Config should not be nil")
		}
		
		if cfg.Server.Mode != "test" {
			t.Errorf("Expected test mode, got %s", cfg.Server.Mode)
		}
		
		fmt.Printf("✓ Config loaded successfully\n")
	})
	
	t.Run("Logger Manager", func(t *testing.T) {
		logConfig := &config.LogConfig{
			Level:  "error",
			Format: "json",
			Output: "console",
		}
		
		logManager, err := logger.New(logConfig)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}
		
		if logManager.GetLevel() != "error" {
			t.Errorf("Expected error level, got %s", logManager.GetLevel())
		}
		
		fmt.Printf("✓ Logger created successfully\n")
	})
	
	t.Run("Utils Functions", func(t *testing.T) {
		// 测试字符串工具
		camelCase := utils.Str.CamelCase("hello_world")
		if camelCase != "helloWorld" {
			t.Errorf("Expected helloWorld, got %s", camelCase)
		}
		
		// 测试JSON工具
		data := map[string]string{"key": "value"}
		jsonStr, err := utils.JSON.ToJSON(data)
		if err != nil {
			t.Fatalf("Failed to convert to JSON: %v", err)
		}
		
		var result map[string]string
		err = utils.JSON.FromJSON(jsonStr, &result)
		if err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}
		
		if result["key"] != "value" {
			t.Errorf("Expected value, got %s", result["key"])
		}
		
		// 测试时间工具
		now := utils.Time.Now()
		if now.IsZero() {
			t.Error("Time should not be zero")
		}
		
		fmt.Printf("✓ Utils functions working correctly\n")
	})
}

func TestIntegration(t *testing.T) {
	// 集成测试示例
	t.Run("Full Application Stack", func(t *testing.T) {
		// 这里可以添加完整的应用栈测试
		// 包括数据库连接、缓存连接、HTTP服务器等
		
		fmt.Printf("✓ Integration test placeholder\n")
	})
}

func BenchmarkStringUtils(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.Str.CamelCase("hello_world_test_string")
	}
}

func BenchmarkJSONUtils(b *testing.B) {
	data := map[string]interface{}{
		"name":  "John",
		"age":   30,
		"email": "john@example.com",
	}
	
	for i := 0; i < b.N; i++ {
		jsonStr, _ := utils.JSON.ToJSON(data)
		var result map[string]interface{}
		utils.JSON.FromJSON(jsonStr, &result)
	}
}

// 运行测试的函数
func main() {
	fmt.Println("Running HWHKit-Go Tests...")
	
	// 运行单元测试
	fmt.Println("\n=== Running Unit Tests ===")
	testing.Main(func(pat, str string) (bool, error) {
		return true, nil
	}, []testing.InternalTest{
		{"TestApplicationComponents", TestApplicationComponents},
		{"TestIntegration", TestIntegration},
	}, []testing.InternalBenchmark{
		{"BenchmarkStringUtils", BenchmarkStringUtils},
		{"BenchmarkJSONUtils", BenchmarkJSONUtils},
	}, []testing.InternalExample{})
}