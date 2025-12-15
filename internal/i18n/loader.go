package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Loader i18n资源加载器
type Loader struct {
	bundle           *i18n.Bundle
	defaultLanguage  string
	supportedLangs   []string
	fallbackLanguage string
	resourceDirs     []string
}

// NewLoader 创建新的资源加载器
func NewLoader(defaultLang string, supportedLangs []string, fallbackLang string, resourceDirs []string) *Loader {
	// 创建语言包
	bundle := i18n.NewBundle(language.MustParse(defaultLang))
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	return &Loader{
		bundle:           bundle,
		defaultLanguage:  defaultLang,
		supportedLangs:   supportedLangs,
		fallbackLanguage: fallbackLang,
		resourceDirs:     resourceDirs,
	}
}

// LoadMessages 加载所有消息
func (l *Loader) LoadMessages() error {
	fmt.Printf("Starting to load i18n messages from %d directories\n", len(l.resourceDirs))

	// 用于跟踪所有已加载的键，检测跨文件的重复（按语言分组）
	langToKeysMap := make(map[string]map[string]string) // lang -> key -> 文件路径
	var duplicateKeysError error

	// 直接加载所有支持的语言文件
	for _, lang := range l.supportedLangs {
		// 初始化该语言的键映射
		if _, exists := langToKeysMap[lang]; !exists {
			langToKeysMap[lang] = make(map[string]string)
		}

		// 使用完整的语言代码作为文件名
		langFile := lang // 例如，使用 "zh-CN" 而不是 "zh"
		fmt.Printf("Loading messages for language: %s\n", lang)

		// 遍历所有目录
		for _, dir := range l.resourceDirs {
			// 查找该语言的文件
			jsonFile := filepath.Join(dir, langFile+".json")

			// 检查文件是否存在
			if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
				fmt.Printf("Language file does not exist: %s\n", jsonFile)
				continue
			}

			fmt.Printf("Loading language file: %s\n", jsonFile)

			// 读取文件内容
			content, err := os.ReadFile(jsonFile)
			if err != nil {
				fmt.Printf("Failed to read message file %s: %v\n", jsonFile, err)
				return fmt.Errorf("failed to read message file %s: %w", jsonFile, err)
			}

			// 解析为嵌套的 JSON 结构
			var nested map[string]interface{}
			if err := json.Unmarshal(content, &nested); err != nil {
				fmt.Printf("Failed to parse message file %s: %v\n", jsonFile, err)
				return fmt.Errorf("failed to parse message file %s: %w", jsonFile, err)
			}

			// 扁平化嵌套结构
			flat, duplicateKeys := flattenJSON(nested, "")

			// 检查文件内部的重复键
			if len(duplicateKeys) > 0 {
				errMsg := fmt.Sprintf("Warning: Found duplicate keys in language file %s: %v\n", jsonFile, duplicateKeys)
				fmt.Print(errMsg)

				// 如果是第一个重复键错误，记录下来
				if duplicateKeysError == nil {
					duplicateKeysError = fmt.Errorf(errMsg)
				}
			}

			// 检查同一语言内跨文件的重复键
			crossFileDuplicates := []string{}
			for key := range flat {
				if existingFile, exists := langToKeysMap[lang][key]; exists {
					crossFileDuplicates = append(crossFileDuplicates, fmt.Sprintf("%s (already in %s)", key, existingFile))
				} else {
					langToKeysMap[lang][key] = jsonFile
				}
			}

			if len(crossFileDuplicates) > 0 {
				errMsg := fmt.Sprintf("Warning: Found cross-file duplicate keys in %s: %v\n", jsonFile, crossFileDuplicates)
				fmt.Print(errMsg)

				// 如果是第一个重复键错误，记录下来
				if duplicateKeysError == nil {
					duplicateKeysError = fmt.Errorf(errMsg)
				}
			}

			// 将扁平化的结构添加到 bundle 中
			langTag := language.MustParse(lang)
			messageCount := 0

			for key, value := range flat {
				l.bundle.AddMessages(langTag, &i18n.Message{
					ID:    key,
					Other: value,
				})
				messageCount++
			}

			// 打印加载的消息数量和语言标签
			fmt.Printf("Successfully loaded message file: %s with %d messages for language %s\n",
				jsonFile, messageCount, langFile)

			// 确保语言标签与支持的语言匹配
			if langFile != lang {
				fmt.Printf("Remapping language tag from %s to %s\n", langFile, lang)
			}
		}
	}

	// 打印所有已加载的语言
	tags := l.bundle.LanguageTags()
	fmt.Printf("Completed loading i18n messages. Available languages: %v\n", tags)

	// 如果有重复键错误，返回错误
	if duplicateKeysError != nil {
		return fmt.Errorf("duplicate keys found during i18n loading: %w", duplicateKeysError)
	}

	return nil
}

// loadMessagesFromDir 从目录加载消息
func (l *Loader) loadMessagesFromDir(dir string) error {
	// 记录加载开始
	fmt.Printf("Loading i18n messages from directory: %s\n", dir)

	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %s\n", dir)
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	// 遍历目录
	fileCount := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err
		}

		// 只处理json文件
		if !info.IsDir() && strings.HasSuffix(path, ".json") {
			fmt.Printf("Found i18n file: %s\n", path)
			// 加载消息文件
			_, err := l.bundle.LoadMessageFile(path)
			if err != nil {
				fmt.Printf("Failed to load message file %s: %v\n", path, err)
				return fmt.Errorf("failed to load message file %s: %w", path, err)
			}
			fileCount++
			fmt.Printf("Successfully loaded message file: %s\n", path)
		}

		return nil
	})

	fmt.Printf("Loaded %d message files from directory: %s\n", fileCount, dir)
	if err != nil {
		fmt.Printf("Error walking directory %s: %v\n", dir, err)
	}

	return err
}

// GetBundle 获取语言包
func (l *Loader) GetBundle() *i18n.Bundle {
	return l.bundle
}

// GetDefaultLanguage 获取默认语言
func (l *Loader) GetDefaultLanguage() string {
	return l.defaultLanguage
}

// GetSupportedLanguages 获取支持的语言列表
func (l *Loader) GetSupportedLanguages() []string {
	return l.supportedLangs
}

// GetFallbackLanguage 获取回退语言
func (l *Loader) GetFallbackLanguage() string {
	return l.fallbackLanguage
}

// GetLanguage 获取当前上下文中的语言类型
func (l *Loader) GetLanguage(c *gin.Context) string {
	// 默认返回默认语言
	if c == nil {
		return l.defaultLanguage
	}

	// 尝试从请求头中获取语言
	acceptLanguage := c.GetHeader("Accept-Language")
	if acceptLanguage != "" {
		// 简单处理，取第一个语言标签
		lang := strings.Split(acceptLanguage, ",")[0]
		lang = strings.Split(lang, ";")[0]

		// 检查是否支持该语言
		for _, supportedLang := range l.supportedLangs {
			if strings.EqualFold(lang, supportedLang) {
				return supportedLang
			}
		}
	}

	// 如果没有找到匹配的语言，返回默认语言
	return l.defaultLanguage
}

// flattenJSON 将嵌套的 JSON 结构扁平化为键值对
// 例如，将 {"user":{"profile":{"name":"用户名"}}} 转换为 {"user.profile.name":"用户名"}
// 如果检测到重复的键，将返回错误
func flattenJSON(nested map[string]interface{}, prefix string) (map[string]string, []string) {
	flat := make(map[string]string)
	duplicateKeys := []string{}

	for key, value := range nested {
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			// 递归处理嵌套结构
			nestedFlat, nestedDuplicates := flattenJSON(v, newKey)
			for k, v := range nestedFlat {
				if _, exists := flat[k]; exists {
					duplicateKeys = append(duplicateKeys, k)
				}
				flat[k] = v
			}
			duplicateKeys = append(duplicateKeys, nestedDuplicates...)
		case string:
			// 直接添加字符串值
			if _, exists := flat[newKey]; exists {
				duplicateKeys = append(duplicateKeys, newKey)
			}
			flat[newKey] = v
		case float64:
			// 处理数字
			if _, exists := flat[newKey]; exists {
				duplicateKeys = append(duplicateKeys, newKey)
			}
			flat[newKey] = fmt.Sprintf("%v", v)
		case bool:
			// 处理布尔值
			if _, exists := flat[newKey]; exists {
				duplicateKeys = append(duplicateKeys, newKey)
			}
			flat[newKey] = fmt.Sprintf("%v", v)
		case nil:
			// 处理空值
			if _, exists := flat[newKey]; exists {
				duplicateKeys = append(duplicateKeys, newKey)
			}
			flat[newKey] = ""
		case []interface{}:
			// 处理数组（不支持复杂数组，只支持简单值）
			for i, item := range v {
				arrayKey := fmt.Sprintf("%s[%d]", newKey, i)
				if _, exists := flat[arrayKey]; exists {
					duplicateKeys = append(duplicateKeys, arrayKey)
				}

				switch itemValue := item.(type) {
				case string:
					flat[arrayKey] = itemValue
				case float64:
					flat[arrayKey] = fmt.Sprintf("%v", itemValue)
				case bool:
					flat[arrayKey] = fmt.Sprintf("%v", itemValue)
				default:
					// 对于复杂类型，转换为 JSON 字符串
					jsonBytes, _ := json.Marshal(item)
					flat[arrayKey] = string(jsonBytes)
				}
			}
		default:
			// 对于其他类型，转换为 JSON 字符串
			if _, exists := flat[newKey]; exists {
				duplicateKeys = append(duplicateKeys, newKey)
			}
			jsonBytes, _ := json.Marshal(v)
			flat[newKey] = string(jsonBytes)
		}
	}

	return flat, duplicateKeys
}
