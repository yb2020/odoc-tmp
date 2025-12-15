// pkg/http_client/js_engine.go
package http_client

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dop251/goja"
	"github.com/yb2020/odoc/pkg/logging"
)

// JSEngine 封装JavaScript引擎
type JSEngine struct {
	vm     *goja.Runtime
	logger logging.Logger
}

// NewJSEngine 创建新的JavaScript引擎
func NewJSEngine(logger logging.Logger) *JSEngine {
	return &JSEngine{
		vm:     goja.New(),
		logger: logger,
	}
}

// LoadScriptFromFile 从文件加载JavaScript脚本
func (e *JSEngine) LoadScriptFromFile(scriptPath string) error {
	// 读取JavaScript文件
	scriptBytes, err := os.ReadFile(scriptPath)
	if err != nil {
		e.logger.Error("read JavaScript file failed", "error", err, "path", scriptPath)
		return fmt.Errorf("read JavaScript file failed: %w", err)
	}

	// 执行JavaScript代码
	_, err = e.vm.RunString(string(scriptBytes))
	if err != nil {
		e.logger.Error("execute JavaScript code failed", "error", err)
		return fmt.Errorf("execute JavaScript code failed: %w", err)
	}

	e.logger.Info("success load JavaScript script", "path", scriptPath)
	return nil
}

// LoadScriptFromResource 从资源目录加载JavaScript脚本
func (e *JSEngine) LoadScriptFromResource(scriptName string) error {
	// 确定资源目录路径
	resourcePath := filepath.Join("resources", "js", scriptName)

	// 检查是否存在于当前目录
	if _, err := os.Stat(resourcePath); os.IsNotExist(err) {
		// 尝试从项目根目录查找
		resourcePath = filepath.Join(".", "resources", "js", scriptName)
		if _, err := os.Stat(resourcePath); os.IsNotExist(err) {
			return fmt.Errorf("JavaScript resource file not found: %s", scriptName)
		}
	}

	return e.LoadScriptFromFile(resourcePath)
}

// LoadScriptContent 直接加载JavaScript代码内容
func (e *JSEngine) LoadScriptContent(scriptContent string) error {
	_, err := e.vm.RunString(scriptContent)
	if err != nil {
		e.logger.Error("execute JavaScript code content failed", "error", err)
		return fmt.Errorf("execute JavaScript code content failed: %w", err)
	}
	return nil
}

// CallFunction 调用JavaScript函数
func (e *JSEngine) CallFunction(functionName string, args ...interface{}) (interface{}, error) {
	// 获取JavaScript函数
	jsFunc := e.vm.Get(functionName)
	if jsFunc == nil {
		return nil, fmt.Errorf("function %s not found", functionName)
	}

	// 检查是否是函数
	fn, ok := goja.AssertFunction(jsFunc)
	if !ok {
		return nil, fmt.Errorf("%s is not a function", functionName)
	}

	// 将Go参数转换为JavaScript值
	jsArgs := make([]goja.Value, len(args))
	for i, arg := range args {
		jsArgs[i] = e.vm.ToValue(arg)
	}

	// 调用函数
	result, err := fn(goja.Undefined(), jsArgs...)
	if err != nil {
		e.logger.Error("call JavaScript function failed", "function", functionName, "error", err)
		return nil, fmt.Errorf("call JavaScript function failed: %w", err)
	}

	// 将JavaScript值转换为Go值
	goValue := result.Export()

	return goValue, nil
}

// SetGlobal 设置全局变量
func (e *JSEngine) SetGlobal(name string, value interface{}) {
	e.vm.Set(name, value)
}

// GetGlobal 获取全局变量
func (e *JSEngine) GetGlobal(name string) (interface{}, error) {
	value := e.vm.Get(name)
	if value == nil {
		return nil, fmt.Errorf("global variable %s not found", name)
	}

	return value.Export(), nil
}
