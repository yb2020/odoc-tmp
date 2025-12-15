package util

import (
	"unicode"

	"github.com/abadojack/whatlanggo"
	"github.com/yb2020/odoc/services/parse/constant"
)

func DetectLanguage(title, content string) string {
	// 首先使用 whatlanggo 检测  使用whatlanggo检测标题，如果不可靠，再使用字符统计方法
	info := whatlanggo.Detect(title)
	// 如果置信度足够高，直接返回结果
	if info.IsReliable() {
		if info.Lang == whatlanggo.Cmn {
			return constant.LanguageZhCN
		} else if info.Lang == whatlanggo.Eng {
			return constant.LanguageEnUS
		}
	}
	// 如果不可靠，回退到字符统计方法
	return DetectLanguageFallback(content)
}

func DetectLanguageFallback(content string) string {
	chineseChars := 0
	for _, r := range content {
		if unicode.Is(unicode.Han, r) {
			chineseChars++
		}
	}

	if len(content) > 0 && float64(chineseChars)/float64(len(content)) > 0.1 {
		return constant.LanguageZhCN
	}
	return constant.LanguageEnUS
}
