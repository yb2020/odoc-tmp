// pkg/squid_proxy/example/simple.go
// 简单使用示例 - 演示核心功能
package main

import (
	"fmt"
	"time"

	"github.com/yb2020/odoc/pkg/squid_proxy"
)

func main() {
	fmt.Println("=== Squid代理工具类简单示例 ===")

	// 1. 创建代理配置
	config := &squid_proxy.SquidConfig{
		ProxyURL:  "http://192.168.218.19:31128",      // 使用指定的Squid代理服务器
		Username:  "",                                  // 如果需要认证，填写用户名
		Password:  "",                                  // 如果需要认证，填写密码
		Timeout:   60 * time.Second,                    // 增加超时时间，arXiv文件可能较大
		UserAgent: "Mozilla/5.0 (compatible; PDF-Downloader/1.0)", // 使用更通用的User-Agent
	}

	// 2. 创建代理客户端
	client, err := squid_proxy.NewSquidProxyClient(config)
	if err != nil {
		fmt.Printf("创建代理客户端失败: %v\n", err)
		return
	}

	// 3. 测试代理连接（可选）
	fmt.Println("测试代理连接...")
	if err := client.TestConnection(); err != nil {
		fmt.Printf("代理连接失败: %v\n", err)
		fmt.Println("注意: 请确保代理服务器地址正确且可访问")
		return
	}
	fmt.Println("代理连接成功!")

	// 4. 下载PDF文件到内存
	pdfURL := "https://arxiv.org/pdf/2508.18942"
	fmt.Printf("正在下载arXiv论文: %s\n", pdfURL)

	result, err := client.DownloadFile(pdfURL, squid_proxy.PDFDownloadOptions())
	if err != nil {
		fmt.Printf("下载失败: %v\n", err)
		return
	}

	fmt.Printf("下载成功!\n")
	fmt.Printf("文件大小: %s\n", squid_proxy.FormatFileSize(result.ContentLength))
	fmt.Printf("下载耗时: %s\n", squid_proxy.FormatDuration(result.Duration))

	// 5. 或者直接下载到文件
	fmt.Println("\n保存文件到本地...")
	filename := "./arxiv_2508.18942.pdf"
	err = client.DownloadFileToPath(pdfURL, filename, squid_proxy.PDFDownloadOptions())
	if err != nil {
		fmt.Printf("保存文件失败: %v\n", err)
		return
	}

	fmt.Printf("文件已保存到 %s\n", filename)
	fmt.Println("\n=== 示例完成 ===")
}
