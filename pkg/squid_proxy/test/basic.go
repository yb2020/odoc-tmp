// pkg/squid_proxy/test/basic_test.go
package main

import (
	"fmt"
	"time"

	"github.com/yb2020/odoc/pkg/squid_proxy"
)

func main() {
	fmt.Println("=== Squid代理工具类基础测试 ===\n")

	// 测试1: 配置验证
	fmt.Println("1. 测试配置验证...")
	testConfigValidation()

	// 测试2: 错误处理
	fmt.Println("\n2. 测试错误处理...")
	testErrorHandling()

	// 测试3: 工具函数
	fmt.Println("\n3. 测试工具函数...")
	testUtilityFunctions()

	// 测试4: PDF验证功能
	fmt.Println("\n4. 测试PDF验证功能...")
	testPDFValidation()

	// 测试5: 无代理模式下的基础功能（用于测试代码逻辑）
	fmt.Println("\n5. 测试基础功能（无代理模式）...")
	testBasicFunctionality()

	fmt.Println("\n=== 测试完成 ===")
}

// testConfigValidation 测试配置验证
func testConfigValidation() {
	// 测试有效配置
	validConfig := &squid_proxy.SquidConfig{
		ProxyURL:  "http://proxy.example.com:3128",
		Username:  "testuser",
		Password:  "testpass",
		Timeout:   30 * time.Second,
		UserAgent: "TestAgent/1.0",
	}

	err := validConfig.Validate()
	if err != nil {
		fmt.Printf("  ❌ 有效配置验证失败: %v\n", err)
	} else {
		fmt.Printf("  ✅ 有效配置验证通过\n")
	}

	// 测试无效配置
	invalidConfig := &squid_proxy.SquidConfig{
		ProxyURL: "", // 空URL
	}

	err = invalidConfig.Validate()
	if err != nil {
		fmt.Printf("  ✅ 无效配置正确被拒绝: %v\n", err)
	} else {
		fmt.Printf("  ❌ 无效配置未被检测到\n")
	}

	// 测试默认配置
	defaultConfig := squid_proxy.DefaultConfig("http://localhost:3128")
	err = defaultConfig.Validate()
	if err != nil {
		fmt.Printf("  ❌ 默认配置验证失败: %v\n", err)
	} else {
		fmt.Printf("  ✅ 默认配置验证通过\n")
	}
}

// testErrorHandling 测试错误处理
func testErrorHandling() {
	// 创建各种类型的错误
	configErr := squid_proxy.NewSquidError(squid_proxy.ErrInvalidConfig, "测试配置错误")
	networkErr := squid_proxy.NewSquidError(squid_proxy.ErrProxyConnection, "测试网络错误")
	
	fmt.Printf("  配置错误: %v\n", configErr)
	fmt.Printf("  网络错误: %v\n", networkErr)

	// 测试错误类型判断
	if squid_proxy.IsConfigError(configErr) {
		fmt.Printf("  ✅ 配置错误类型判断正确\n")
	} else {
		fmt.Printf("  ❌ 配置错误类型判断错误\n")
	}

	if squid_proxy.IsNetworkError(networkErr) {
		fmt.Printf("  ✅ 网络错误类型判断正确\n")
	} else {
		fmt.Printf("  ❌ 网络错误类型判断错误\n")
	}

	// 测试重试判断
	if squid_proxy.IsRetryableError(networkErr) {
		fmt.Printf("  ✅ 网络错误可重试判断正确\n")
	} else {
		fmt.Printf("  ❌ 网络错误可重试判断错误\n")
	}

	if !squid_proxy.IsRetryableError(configErr) {
		fmt.Printf("  ✅ 配置错误不可重试判断正确\n")
	} else {
		fmt.Printf("  ❌ 配置错误不可重试判断错误\n")
	}
}

// testUtilityFunctions 测试工具函数
func testUtilityFunctions() {
	// 测试URL验证
	validURL := "https://example.com/document.pdf"
	invalidURL := "invalid-url"

	if squid_proxy.ValidateURL(validURL) == nil {
		fmt.Printf("  ✅ 有效URL验证通过: %s\n", validURL)
	} else {
		fmt.Printf("  ❌ 有效URL验证失败: %s\n", validURL)
	}

	if squid_proxy.ValidateURL(invalidURL) != nil {
		fmt.Printf("  ✅ 无效URL正确被拒绝: %s\n", invalidURL)
	} else {
		fmt.Printf("  ❌ 无效URL未被检测到: %s\n", invalidURL)
	}

	// 测试PDF URL判断
	pdfURL := "https://example.com/document.pdf"
	nonPdfURL := "https://example.com/image.jpg"

	if squid_proxy.IsPDFURL(pdfURL) {
		fmt.Printf("  ✅ PDF URL判断正确: %s\n", pdfURL)
	} else {
		fmt.Printf("  ❌ PDF URL判断错误: %s\n", pdfURL)
	}

	if !squid_proxy.IsPDFURL(nonPdfURL) {
		fmt.Printf("  ✅ 非PDF URL判断正确: %s\n", nonPdfURL)
	} else {
		fmt.Printf("  ❌ 非PDF URL判断错误: %s\n", nonPdfURL)
	}

	// 测试文件大小格式化
	sizes := []int64{1024, 1024 * 1024, 1024 * 1024 * 1024}
	expected := []string{"1.0 KB", "1.0 MB", "1.0 GB"}

	for i, size := range sizes {
		formatted := squid_proxy.FormatFileSize(size)
		if formatted == expected[i] {
			fmt.Printf("  ✅ 文件大小格式化正确: %d -> %s\n", size, formatted)
		} else {
			fmt.Printf("  ❌ 文件大小格式化错误: %d -> %s (期望: %s)\n", size, formatted, expected[i])
		}
	}

	// 测试时间格式化
	durations := []time.Duration{
		500 * time.Microsecond,
		500 * time.Millisecond,
		2 * time.Second,
	}

	for _, duration := range durations {
		formatted := squid_proxy.FormatDuration(duration)
		fmt.Printf("  ✅ 时间格式化: %v -> %s\n", duration, formatted)
	}

	// 测试文件名提取
	testURL := "https://example.com/path/document.pdf"
	fileName := squid_proxy.GetFileNameFromURL(testURL)
	if fileName == "document.pdf" {
		fmt.Printf("  ✅ 文件名提取正确: %s -> %s\n", testURL, fileName)
	} else {
		fmt.Printf("  ❌ 文件名提取错误: %s -> %s\n", testURL, fileName)
	}
}

// testBasicFunctionality 测试基础功能（无代理模式）
func testBasicFunctionality() {
	// 注意：这个测试不会真正使用代理，只是测试代码逻辑
	fmt.Println("  注意：此测试仅验证代码逻辑，不会实际连接代理服务器")

	// 创建一个无效的代理配置用于测试
	config := &squid_proxy.SquidConfig{
		ProxyURL:  "http://nonexistent-proxy.test:3128",
		Timeout:   5 * time.Second,
		UserAgent: "TestClient/1.0",
	}

	client, err := squid_proxy.NewSquidProxyClient(config)
	if err != nil {
		fmt.Printf("  ❌ 客户端创建失败: %v\n", err)
		return
	}

	fmt.Printf("  ✅ 客户端创建成功\n")

	// 测试选项创建
	pdfOptions := squid_proxy.PDFDownloadOptions()
	if pdfOptions != nil && pdfOptions.Headers["Accept"] == "application/pdf,*/*" {
		fmt.Printf("  ✅ PDF下载选项创建正确\n")
	} else {
		fmt.Printf("  ❌ PDF下载选项创建错误\n")
	}

	largeFileOptions := squid_proxy.LargeFileDownloadOptions()
	if largeFileOptions != nil && largeFileOptions.MaxSize == 500*1024*1024 {
		fmt.Printf("  ✅ 大文件下载选项创建正确\n")
	} else {
		fmt.Printf("  ❌ 大文件下载选项创建错误\n")
	}

	defaultOptions := squid_proxy.DefaultDownloadOptions()
	if defaultOptions != nil && defaultOptions.RetryCount == 3 {
		fmt.Printf("  ✅ 默认下载选项创建正确\n")
	} else {
		fmt.Printf("  ❌ 默认下载选项创建错误\n")
	}

	// 测试连接测试功能（预期会失败，因为代理不存在）
	err = client.TestConnection()
	if err != nil {
		fmt.Printf("  ✅ 连接测试正确检测到代理不可用: %v\n", err)
	} else {
		fmt.Printf("  ❌ 连接测试未检测到代理问题\n")
	}

	fmt.Println("  基础功能测试完成")
}

// testPDFValidation 测试PDF验证功能
func testPDFValidation() {
	// 创建一个最小的有效PDF文件
	validPDF := []byte(`%PDF-1.4
1 0 obj
<<
/Type /Catalog
/Pages 2 0 R
>>
endobj
2 0 obj
<<
/Type /Pages
/Kids [3 0 R]
/Count 1
>>
endobj
3 0 obj
<<
/Type /Page
/Parent 2 0 R
/MediaBox [0 0 612 792]
>>
endobj
xref
0 4
0000000000 65535 f
0000000009 00000 n
0000000058 00000 n
0000000115 00000 n
trailer
<<
/Size 4
/Root 1 0 R
>>
startxref
190
%%EOF`)

	// 测试PDF魔数验证
	fmt.Println("  测试PDF魔数验证...")
	if err := squid_proxy.ValidatePDFMagicNumber(validPDF); err == nil {
		fmt.Printf("    ✅ 有效PDF魔数验证通过\n")
	} else {
		fmt.Printf("    ❌ 有效PDF魔数验证失败: %v\n", err)
	}

	invalidMagic := []byte("<!DOCTYPE html>")
	if err := squid_proxy.ValidatePDFMagicNumber(invalidMagic); err != nil {
		fmt.Printf("    ✅ 无效魔数正确被拒绝: %v\n", err)
	} else {
		fmt.Printf("    ❌ 无效魔数未被检测到\n")
	}

	// 测试Content-Type验证
	fmt.Println("\n  测试Content-Type验证...")
	if err := squid_proxy.ValidatePDFContentType("application/pdf"); err == nil {
		fmt.Printf("    ✅ 有效Content-Type验证通过\n")
	} else {
		fmt.Printf("    ❌ 有效Content-Type验证失败: %v\n", err)
	}

	if err := squid_proxy.ValidatePDFContentType("application/pdf; charset=utf-8"); err == nil {
		fmt.Printf("    ✅ 带参数的Content-Type验证通过\n")
	} else {
		fmt.Printf("    ❌ 带参数的Content-Type验证失败: %v\n", err)
	}

	if err := squid_proxy.ValidatePDFContentType("text/html"); err != nil {
		fmt.Printf("    ✅ 无效Content-Type正确被拒绝: %v\n", err)
	} else {
		fmt.Printf("    ❌ 无效Content-Type未被检测到\n")
	}

	// 测试PDF结构完整性验证
	fmt.Println("\n  测试PDF结构完整性验证...")
	if err := squid_proxy.ValidatePDFStructure(validPDF); err == nil {
		fmt.Printf("    ✅ 有效PDF结构验证通过\n")
	} else {
		fmt.Printf("    ❌ 有效PDF结构验证失败: %v\n", err)
	}

	incompletePDF := []byte("%PDF-1.4\n1 0 obj\n<<\n/Type /Catalog\n>>\nendobj\n")
	if err := squid_proxy.ValidatePDFStructure(incompletePDF); err != nil {
		fmt.Printf("    ✅ 不完整PDF正确被拒绝: %v\n", err)
	} else {
		fmt.Printf("    ❌ 不完整PDF未被检测到\n")
	}

	// 测试综合PDF验证
	fmt.Println("\n  测试综合PDF验证...")
	if err := squid_proxy.ValidatePDFFile(validPDF, "application/pdf"); err == nil {
		fmt.Printf("    ✅ 完整PDF验证通过\n")
	} else {
		fmt.Printf("    ❌ 完整PDF验证失败: %v\n", err)
	}

	// 测试各种错误情况
	if err := squid_proxy.ValidatePDFFile(validPDF, "text/html"); err != nil {
		fmt.Printf("    ✅ 错误的Content-Type正确被拒绝\n")
	} else {
		fmt.Printf("    ❌ 错误的Content-Type未被检测到\n")
	}

	if err := squid_proxy.ValidatePDFFile(invalidMagic, "application/pdf"); err != nil {
		fmt.Printf("    ✅ 错误的魔数正确被拒绝\n")
	} else {
		fmt.Printf("    ❌ 错误的魔数未被检测到\n")
	}

	if err := squid_proxy.ValidatePDFFile(incompletePDF, "application/pdf"); err != nil {
		fmt.Printf("    ✅ 不完整的PDF正确被拒绝\n")
	} else {
		fmt.Printf("    ❌ 不完整的PDF未被检测到\n")
	}

	fmt.Println("\n  PDF验证功能测试完成")
}
