package utils

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/yb2020/odoc/pkg/logging"
)

// ToValidUTF8 将字符串转换为合法的UTF-8格式，会替换掉无效的字符。
func ToValidUTF8(s string, replacement string) string {
	return strings.ToValidUTF8(s, replacement)
}

// TextUtils 文本处理工具
type TextUtils struct {
	logger logging.Logger
}

// NewTextUtils 创建文本处理工具实例
func NewTextUtils(logger logging.Logger) *TextUtils {
	return &TextUtils{
		logger: logger,
	}
}

// RemoveSpecialWords 去除内容中的特殊字符
// content: 需要处理的内容
// list: 需要替换的特殊字符正则表达式列表
func (t *TextUtils) RemoveSpecialWords(content string, list []string) string {
	if len(list) == 0 {
		return content
	}

	result := content

	// 预处理：将 NULL 字符替换为空格
	result = strings.Map(func(r rune) rune {
		if r == 0 {
			return ' '
		}
		return r
	}, result)

	for _, specialWord := range list {
		// 尝试编译正则表达式
		reg, err := regexp.Compile(specialWord)
		if err != nil {
			t.logger.Error("msg", "编译正则表达式失败", "pattern", specialWord, "error", err.Error())
			continue
		}

		// 替换特殊字符为空格
		result = reg.ReplaceAllString(result, " ")
	}

	// 替换所有控制字符为空格
	result = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, result)

	// 清理连续的空格
	spaceReg := regexp.MustCompile(`\s+`)
	result = spaceReg.ReplaceAllString(result, " ")
	
	// 清理首尾空格
	result = strings.TrimSpace(result)

	return result
}
