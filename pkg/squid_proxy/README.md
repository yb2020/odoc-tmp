# Squid代理工具类

一个专门为Go项目设计的独立Squid代理客户端工具类，专门优化了PDF等文件下载场景。

## 特性

✅ **完全独立** - 零依赖，只使用Go标准库  
✅ **专门优化** - 针对PDF等文件下载场景特别优化  
✅ **简单易用** - 清晰的API，最少的配置  
✅ **健壮可靠** - 完善的错误处理和重试机制  
✅ **PDF验证** - 多层次PDF文件合法性验证  

## 文件结构

```
pkg/squid_proxy/
├── squid_proxy.go      # 核心代理客户端实现
├── config.go          # 配置结构体和验证
├── errors.go          # 错误处理系统
├── utils.go           # 工具函数
├── example/           # 使用示例
│   └── simple.go
├── test/              # 测试文件
│   └── basic.go
└── README.md         # 使用文档
```

## 快速开始

### 1. 基础使用

```go
package main

import (
    "fmt"
    "time"
    "github.com/yb2020/go-sea/pkg/squid_proxy"
)

func main() {
    // 创建代理配置
    config := &squid_proxy.SquidConfig{
        ProxyURL:  "http://your-squid-proxy.com:3128",
        Username:  "proxy_user",     // 可选，如果代理需要认证
        Password:  "proxy_pass",     // 可选，如果代理需要认证
        Timeout:   30 * time.Second,
        UserAgent: "PDF-Downloader/1.0",
    }
    
    // 创建代理客户端
    client, err := squid_proxy.NewSquidProxyClient(config)
    if err != nil {
        panic(err)
    }
    
    // 下载PDF文件到内存
    result, err := client.DownloadFile(
        "https://example.com/document.pdf",
        squid_proxy.PDFDownloadOptions(),
    )
    if err != nil {
        fmt.Printf("下载失败: %v\n", err)
        return
    }
    
    fmt.Printf("下载成功: %s\n", squid_proxy.FormatFileSize(result.ContentLength))
}
```

### 2. 下载到文件

```go
// 直接下载到文件
err := client.DownloadFileToPath(
    "https://example.com/document.pdf", 
    "./document.pdf", 
    squid_proxy.PDFDownloadOptions(),
)
```

### 3. 流式下载大文件

```go
file, _ := os.Create("large_file.pdf")
defer file.Close()

err := client.GetWithStream(
    "https://example.com/large_document.pdf",
    file,
    squid_proxy.LargeFileDownloadOptions(),
)
```

### 4. 测试代理连接

```go
// 测试代理是否可用
err := client.TestConnection()
if err != nil {
    fmt.Printf("代理连接失败: %v\n", err)
}
```

## 配置说明

### SquidConfig 配置参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| ProxyURL | string | ✅ | Squid代理服务器URL | `"http://proxy.com:3128"` |
| Username | string | ❌ | 代理认证用户名 | `"user123"` |
| Password | string | ❌ | 代理认证密码 | `"password"` |
| Timeout | time.Duration | ❌ | 请求超时时间 | `30 * time.Second` |
| UserAgent | string | ❌ | 用户代理字符串 | `"MyApp/1.0"` |

### DownloadOptions 下载选项

| 参数 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| Headers | map[string]string | 自定义请求头 | `{}` |
| MaxSize | int64 | 最大文件大小限制（字节） | `100MB` |
| Timeout | time.Duration | 单次请求超时 | `60s` |
| RetryCount | int | 重试次数 | `3` |
| RetryDelay | time.Duration | 重试间隔 | `2s` |

## 预设选项

### PDF下载优化选项
```go
options := squid_proxy.PDFDownloadOptions()
// 自动设置PDF相关的Accept头和大小限制
```

### 大文件下载选项
```go
options := squid_proxy.LargeFileDownloadOptions()
// 适用于大文件的超时和重试设置
```

### 默认下载选项
```go
options := squid_proxy.DefaultDownloadOptions()
// 通用的默认设置
```

## API 方法

### 主要方法

#### NewSquidProxyClient
```go
func NewSquidProxyClient(config *SquidConfig) (*SquidProxyClient, error)
```
创建新的Squid代理客户端。

#### DownloadFile
```go
func (c *SquidProxyClient) DownloadFile(targetURL string, options *DownloadOptions) (*DownloadResult, error)
```
通过代理下载文件到内存。

#### DownloadFileToPath
```go
func (c *SquidProxyClient) DownloadFileToPath(targetURL, savePath string, options *DownloadOptions) error
```
通过代理下载文件并保存到指定路径。

#### GetWithStream
```go
func (c *SquidProxyClient) GetWithStream(targetURL string, writer io.Writer, options *DownloadOptions) error
```
通过代理获取流式响应，适用于大文件下载。

#### TestConnection
```go
func (c *SquidProxyClient) TestConnection() error
```
测试代理连接是否正常。

### 工具方法

#### 文件大小格式化
```go
squid_proxy.FormatFileSize(1024) // "1.0 KB"
```

#### 持续时间格式化
```go
squid_proxy.FormatDuration(time.Second) // "1.00s"
```

#### URL验证
```go
err := squid_proxy.ValidateURL("https://example.com/file.pdf")
```

#### PDF URL判断
```go
isPDF := squid_proxy.IsPDFURL("https://example.com/document.pdf") // true
```

## 错误处理

工具类提供了详细的错误分类：

```go
result, err := client.DownloadFile(url, options)
if err != nil {
    if squid_proxy.IsNetworkError(err) {
        // 网络相关错误，可以重试
        fmt.Println("网络错误，请检查网络连接")
    } else if squid_proxy.IsConfigError(err) {
        // 配置错误，需要修正配置
        fmt.Println("配置错误，请检查代理设置")
    } else {
        // 其他错误
        fmt.Printf("下载失败: %v\n", err)
    }
}
```

### 错误类型

- `ErrInvalidConfig` - 配置无效
- `ErrInvalidURL` - URL无效
- `ErrInvalidProxy` - 代理设置无效
- `ErrProxyConnection` - 代理连接失败
- `ErrRequestFailed` - 请求失败
- `ErrHTTPError` - HTTP错误
- `ErrFileTooLarge` - 文件过大
- `ErrReadFailed` - 读取失败
- `ErrWriteFailed` - 写入失败
- `ErrInvalidResponse` - 响应无效（包括PDF验证失败）

## PDF文件验证

工具类提供了多层次的PDF文件合法性验证，确保下载的文件是有效的PDF文档。

### 验证层次

#### 1. Content-Type 验证

验证HTTP响应头的Content-Type是否为PDF类型：

```go
err := squid_proxy.ValidatePDFContentType("application/pdf")
if err != nil {
    fmt.Printf("无效的Content-Type: %v\n", err)
}
```

**支持的Content-Type：**
- `application/pdf`
- `application/x-pdf`
- `application/x-bzpdf`
- `application/x-gzpdf`
- `application/octet-stream` （自动跳过验证，因为很多服务器如 GitHub 会返回通用二进制流）

**注意：** 当 Content-Type 为 `application/octet-stream` 时，验证会自动跳过此步骤，仅依赖魔数和结构验证，这是因为很多 CDN 和文件服务器会将 PDF 作为通用二进制流返回。

#### 2. PDF魔数验证

验证文件头是否以PDF魔数 `%PDF-` 开头：

```go
fileData := []byte("%PDF-1.4\n...")
err := squid_proxy.ValidatePDFMagicNumber(fileData)
if err != nil {
    fmt.Printf("无效的PDF魔数: %v\n", err)
}
```

**验证规则：**
- 文件必须以字节序列 `0x25 0x50 0x44 0x46 0x2D` 开头
- 即ASCII字符串 `%PDF-`
- 这是PDF规范要求的文件头标识

#### 3. PDF结构完整性验证

验证PDF文件是否包含 `%%EOF` 结尾标记：

```go
err := squid_proxy.ValidatePDFStructure(fileData)
if err != nil {
    fmt.Printf("PDF结构不完整: %v\n", err)
}
```

**验证规则：**
- 检查文件末尾（最后1KB）是否包含 `%%EOF` 标记
- PDF规范要求文件必须以此标记结束
- 缺少此标记可能表示文件下载不完整

#### 4. 综合验证

一次性执行所有验证：

```go
fileData := []byte("%PDF-1.4\n...") 
contentType := "application/pdf"

err := squid_proxy.ValidatePDFFile(fileData, contentType)
if err != nil {
    fmt.Printf("PDF验证失败: %v\n", err)
    return
}

fmt.Println("PDF文件验证通过！")
```

### 业务集成示例

在实际业务中使用PDF验证：

```go
func (s *PDFDownloadService) DownloadAndValidatePDF(url string) ([]byte, error) {
    // 1. 下载PDF文件
    result, err := s.proxyClient.DownloadFile(url, squid_proxy.PDFDownloadOptions())
    if err != nil {
        return nil, fmt.Errorf("下载失败: %w", err)
    }
    
    // 2. 验证PDF文件
    contentType := result.Headers.Get("Content-Type")
    if err := squid_proxy.ValidatePDFFile(result.Body, contentType); err != nil {
        return nil, fmt.Errorf("PDF验证失败: %w", err)
    }
    
    // 3. 返回验证通过的PDF数据
    return result.Body, nil
}
```

### 验证失败场景

PDF验证会在以下情况下失败：

1. **错误的Content-Type**
   - 服务器返回的不是PDF类型（如HTML页面、错误页面）
   - 示例：`Content-Type: text/html`

2. **错误的文件魔数**
   - 下载的文件不是PDF格式
   - 示例：下载到了HTML页面、图片或其他格式文件

3. **不完整的PDF文件**
   - 网络中断导致下载不完整
   - 文件缺少 `%%EOF` 结尾标记

4. **文件太小**
   - 文件小于100字节，不可能是有效的PDF

### 最佳实践

1. **总是验证下载的PDF文件**
   ```go
   // 推荐：下载后立即验证
   result, err := client.DownloadFile(url, options)
   if err == nil {
       err = squid_proxy.ValidatePDFFile(result.Body, result.Headers.Get("Content-Type"))
   }
   ```

2. **记录验证失败的详细信息**
   ```go
   if err := squid_proxy.ValidatePDFFile(data, contentType); err != nil {
       log.Printf("错误: PDF验证失败 - URL: %s, Content-Type: %s, 错误: %v", 
           url, contentType, err)
       return err
   }
   ```

3. **分层验证以获取更详细的错误信息**
   ```go
   // 分步验证，便于调试
   if err := squid_proxy.ValidatePDFContentType(contentType); err != nil {
       return fmt.Errorf("Content-Type验证失败: %w", err)
   }
   if err := squid_proxy.ValidatePDFMagicNumber(data); err != nil {
       return fmt.Errorf("魔数验证失败: %w", err)
   }
   if err := squid_proxy.ValidatePDFStructure(data); err != nil {
       return fmt.Errorf("结构验证失败: %w", err)
   }
   ```

## 环境变量配置

推荐使用环境变量配置代理信息：

```bash
export SQUID_PROXY_URL="http://proxy.company.com:3128"
export SQUID_USERNAME="your_username"
export SQUID_PASSWORD="your_password"
```

```go
config := &squid_proxy.SquidConfig{
    ProxyURL:  os.Getenv("SQUID_PROXY_URL"),
    Username:  os.Getenv("SQUID_USERNAME"),
    Password:  os.Getenv("SQUID_PASSWORD"),
    Timeout:   30 * time.Second,
    UserAgent: "MyApp/1.0",
}
```

## 运行示例

查看 `example/simple.go` 文件获取完整的使用示例。

```bash
# 运行基础测试
go run pkg/squid_proxy/test/basic.go

# 运行使用示例（需要配置真实的代理地址）
go run pkg/squid_proxy/example/simple.go

# 测试编译
go build ./pkg/squid_proxy
```

## 业务集成示例

在实际业务中，可以这样封装使用：

```go
type PDFDownloadService struct {
    proxyClient *squid_proxy.SquidProxyClient
}

func NewPDFDownloadService(proxyURL string) (*PDFDownloadService, error) {
    config := &squid_proxy.SquidConfig{
        ProxyURL:  proxyURL,
        Timeout:   30 * time.Second,
        UserAgent: "BusinessApp/1.0",
    }
    
    client, err := squid_proxy.NewSquidProxyClient(config)
    if err != nil {
        return nil, err
    }
    
    return &PDFDownloadService{proxyClient: client}, nil
}

func (s *PDFDownloadService) DownloadPDF(url string) ([]byte, error) {
    result, err := s.proxyClient.DownloadFile(url, squid_proxy.PDFDownloadOptions())
    if err != nil {
        return nil, err
    }
    return result.Body, nil
}
```

## 最佳实践

1. **配置管理**: 使用环境变量管理敏感的代理配置信息
2. **错误处理**: 根据错误类型进行适当的重试或用户提示
3. **文件大小限制**: 设置合理的MaxSize避免内存溢出
4. **超时设置**: 根据文件大小和网络状况调整超时时间
5. **连接测试**: 在批量下载前先测试代理连接

## 注意事项

- 确保Squid代理服务器配置正确且可访问
- 注意文件大小限制，避免下载过大文件导致内存问题
- 合理设置重试次数和间隔，避免对目标服务器造成压力
- 使用合适的User-Agent，避免被目标网站屏蔽

## 许可证

本工具类遵循项目的整体许可证。
