package util

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	pb "github.com/yb2020/odoc-proto/gen/go/parsed"
)

// 标题格式类型定义
type TitleFormat int

const (
	UnknownFormat TitleFormat = iota
	ArabicFormat              // 1. 1.1 2. 等
	ChineseFormat             // 一、 （一） （1）等
	RomanFormat               // I. II. A. B.等
)

// 配置常量
const (
	// FormatDetectionThreshold 格式检测的最小出现次数阈值
	FormatDetectionThreshold = 3

	// MixedFormatThreshold 混合格式检测的最小出现次数阈值
	MixedFormatThreshold = 2

	// ChineseFormatPriority 中文格式优先级的阈值
	ChineseFormatPriority = 2
)

// processedItem 处理过程中的中间结构体
type processedItem struct {
	original   *pb.CatalogueItem // 原始条目
	normalized string            // 规范化后的标题
	titleOrder string            // 用于排序的标题序号
	level      int               // 标题层级
	order      int               // 显示顺序
}

// 统一格式化标题
func NormalizeTitle(title string) string {
	title = strings.TrimSpace(title)

	// 匹配：数字部分 | 可能的多余点 | 空格部分 | 标题内容
	titleRe := regexp.MustCompile(`^(\d+(?:\.\d+)*)(\.?)(\s*)(.+)$`)

	if matches := titleRe.FindStringSubmatch(title); len(matches) >= 5 {
		numberPart := matches[1]                  // 数字部分 (如 "2", "2.1")
		extraDot := matches[2]                    // 可能多余的点 (可能为 "." 或 "")
		textPart := strings.TrimSpace(matches[4]) // 标题文本
		// 标准化规则：
		// 1. 如果数字部分已经包含点(如2.1)，则去掉多余的点
		// 2. 确保数字和文本之间只有一个空格
		if strings.Contains(numberPart, ".") && extraDot == "." {
			extraDot = ""
		}
		// 重新组合
		return fmt.Sprintf("%s%s %s", numberPart, extraDot, textPart)
	}

	return title
}

// 拆分格式化标题
func SplitNormalizedTitle(normalizedTitle string) (numberPart string, textPart string) {
	// 匹配序号和正文部分
	// 正则解释:
	// ^(\d+(?:\.\d+)*) 匹配数字部分（如 1, 1.2, 1.2.3 等）
	// \.? 匹配可选的点
	// \s+ 匹配至少一个空格
	// (.+) 匹配剩余的标题内容
	re := regexp.MustCompile(`^(\d+(?:\.\d+)*)\.?\s+(.+)$`)

	if matches := re.FindStringSubmatch(normalizedTitle); len(matches) >= 3 {
		numberPart = matches[1]
		textPart = matches[2]
	} else {
		// 如果不是标准标题格式，返回整个字符串作为正文
		textPart = normalizedTitle
	}

	return
}

func HandleCatalogueItems(catalogueItems []*pb.CatalogueItem) []*pb.CatalogueItem {
	if len(catalogueItems) == 0 {
		return make([]*pb.CatalogueItem, 0)
	}

	// 第一阶段：预处理并分析
	items := preprocessItems(catalogueItems)
	dominantFormat := analyzeTitleFormats(items)

	// 第二阶段：解析和过滤
	filtered := parseAndFilterItems(items, dominantFormat)

	// 第三阶段：构建目录树并恢复原始标题
	return buildCatalogueTree(filtered)
}

// preprocessItems 预处理所有条目
func preprocessItems(items []*pb.CatalogueItem) []*processedItem {
	var processed []*processedItem
	for _, item := range items {
		processed = append(processed, &processedItem{
			original:   item,
			normalized: item.Title,
		})
	}
	return processed
}

func analyzeTitleFormats(items []*processedItem) TitleFormat {
	formatCount := make(map[TitleFormat]int)

	// 统计各种格式的出现次数
	for _, item := range items {
		format := detectTitleFormat(item.normalized)
		formatCount[format]++
	}

	// 找出出现次数最多的格式
	var maxFormat TitleFormat
	maxCount := 0
	for format, count := range formatCount {
		// 只考虑达到阈值的格式
		if count >= FormatDetectionThreshold && count > maxCount {
			maxCount = count
			maxFormat = format
		}
	}

	// 如果没有格式达到阈值，检查是否存在混合格式
	if maxCount == 0 {
		arabicCount := formatCount[ArabicFormat]
		chineseCount := formatCount[ChineseFormat]

		// 混合格式处理：中文数字和阿拉伯数字混合
		if arabicCount >= MixedFormatThreshold && chineseCount >= MixedFormatThreshold {
			// 中文格式优先
			if chineseCount >= ChineseFormatPriority {
				return ChineseFormat
			}
			return ArabicFormat
		}

		// 检查罗马数字/字母格式
		if formatCount[RomanFormat] >= MixedFormatThreshold {
			return RomanFormat
		}

		return UnknownFormat
	}

	return maxFormat
}

// parseAndFilterItems 解析并过滤条目
func parseAndFilterItems(items []*processedItem, dominantFormat TitleFormat) []*processedItem {
	var filtered []*processedItem
	order := 1

	for _, item := range items {
		// 尝试解析标题
		titleOrder, level, ok := tryParseTitle(item.normalized, dominantFormat)
		if !ok {
			continue
		}

		item.titleOrder = titleOrder
		item.level = level
		item.order = order
		order++
		filtered = append(filtered, item)
	}

	return filtered
}

// buildCatalogueTree 构建目录树并恢复原始数据
func buildCatalogueTree(items []*processedItem) []*pb.CatalogueItem {
	// 按照标题序号排序
	sort.Slice(items, func(i, j int) bool {
		return items[i].titleOrder < items[j].titleOrder
	})

	// 创建映射表
	itemMap := make(map[string]*pb.CatalogueItem)
	var rootItems []*pb.CatalogueItem

	for _, item := range items {
		// 创建新条目，保留原始数据
		newItem := &pb.CatalogueItem{
			Title:      item.original.Title, // 使用原始标题
			Level:      strconv.Itoa(item.level),
			TitleOrder: item.titleOrder,
			Order:      int32(item.order),
			Child:      []*pb.CatalogueItem{},
			Bbox:       item.original.Bbox,
			// 复制其他需要的原始字段...
		}

		// 构建目录树逻辑保持不变...
		titleOrder := item.titleOrder
		lastDotIndex := strings.LastIndex(titleOrder, ".")
		if lastDotIndex == -1 {
			rootItems = append(rootItems, newItem)
			itemMap[titleOrder] = newItem
		} else {
			parentOrder := titleOrder[:lastDotIndex]
			if parent, exists := itemMap[parentOrder]; exists {
				parent.Child = append(parent.Child, newItem)
			} else {
				rootItems = append(rootItems, newItem)
			}
			itemMap[titleOrder] = newItem
		}
	}

	return rootItems
}

// tryParseTitle 尝试解析标题
func tryParseTitle(title string, preferredFormat TitleFormat) (string, int, bool) {
	// 尝试优先使用主导格式解析
	switch preferredFormat {
	case ArabicFormat:
		if result, level, ok := parseArabicFormat(title); ok {
			return result, level, true
		}
	case ChineseFormat:
		if result, level, ok := parseChineseFormat(title); ok {
			return result, level, true
		}
	case RomanFormat:
		if result, level, ok := parseRomanFormat(title); ok {
			return result, level, true
		}
	}

	// 主导格式解析失败，尝试所有格式
	if result, level, ok := parseArabicFormat(title); ok {
		return result, level, true
	}
	if result, level, ok := parseChineseFormat(title); ok {
		return result, level, true
	}
	if result, level, ok := parseRomanFormat(title); ok {
		return result, level, true
	}

	return "", 0, false
}

// parseArabicFormat 解析阿拉伯数字格式
func parseArabicFormat(title string) (string, int, bool) {
	matches := regexp.MustCompile(`^\s*(\d+(?:\.\d+)*)\.?\s*`).FindStringSubmatch(title)
	if len(matches) < 2 {
		return "", 0, false
	}

	titleOrder := matches[1]
	level := len(strings.Split(titleOrder, "."))
	return titleOrder, level, true
}

// parseChineseFormat 解析中文数字格式
func parseChineseFormat(title string) (string, int, bool) {
	// 情况1：一、 二、 三、
	if matches := regexp.MustCompile(`^\s*([一二三四五六七八九十]+)、\s*`).FindStringSubmatch(title); len(matches) > 1 {
		arabicNum := chineseToArabic(matches[1])
		return arabicNum, 1, true
	}

	// 情况2：（一） （二） （1） （2）
	if matches := regexp.MustCompile(`^\s*（([一二三四五六七八九十\d]+)）\s*`).FindStringSubmatch(title); len(matches) > 1 {
		numStr := matches[1]
		var arabicNum string
		if isChineseNumber(numStr) {
			arabicNum = chineseToArabic(numStr)
		} else {
			arabicNum = numStr
		}
		return arabicNum, 2, true // 括号格式默认为二级标题
	}

	return "", 0, false
}

// chineseToArabic 中文数字转阿拉伯数字
func chineseToArabic(chinese string) string {
	mapping := map[string]string{
		"一": "1", "二": "2", "三": "3", "四": "4", "五": "5",
		"六": "6", "七": "7", "八": "8", "九": "9", "十": "10",
		"十一": "11", "十二": "12", "十三": "13", "十四": "14", "十五": "15",
		"十六": "16", "十七": "17", "十八": "18", "十九": "19", "二十": "20",
	}
	return mapping[chinese]
}

// romanToArabic 罗马数字转阿拉伯数字
func romanToArabic(roman string) string {
	mapping := map[string]string{
		"I": "1", "II": "2", "III": "3", "IV": "4", "V": "5",
		"VI": "6", "VII": "7", "VIII": "8", "IX": "9", "X": "10",
		"XI": "11", "XII": "12", "XIII": "13", "XIV": "14", "XV": "15",
		"XVI": "16", "XVII": "17", "XVIII": "18", "XIX": "19", "XX": "20",
	}
	return mapping[roman]
}

// letterToArabic 字母转阿拉伯数字
func letterToArabic(letter string) string {
	return strconv.Itoa(int(letter[0] - 'A' + 1))
}

// parseRomanFormat 解析罗马数字/字母格式
func parseRomanFormat(title string) (string, int, bool) {
	// 情况1：I. II. III.
	if matches := regexp.MustCompile(`^\s*([IVXLCDM]+)\.\s*`).FindStringSubmatch(title); len(matches) > 1 {
		arabicNum := romanToArabic(matches[1])
		return arabicNum, 1, true
	}

	// 情况2：A. B. C.
	if matches := regexp.MustCompile(`^\s*([A-Z])\.\s*`).FindStringSubmatch(title); len(matches) > 1 {
		arabicNum := letterToArabic(matches[1])
		return arabicNum, 1, true
	}

	return "", 0, false
}

// 检测标题格式类型
func detectTitleFormat(title string) TitleFormat {
	// 检测阿拉伯数字格式
	if regexp.MustCompile(`^\s*\d+(?:\.\d+)*\.?\s*`).MatchString(title) {
		return ArabicFormat
	}

	// 检测中文数字格式
	if regexp.MustCompile(`^\s*([一二三四五六七八九十]+、|（[一二三四五六七八九十\d]+）)\s*`).MatchString(title) {
		return ChineseFormat
	}

	// 检测罗马数字/字母格式
	if regexp.MustCompile(`^\s*([IVXLCDM]+|[A-Z])\.\s*`).MatchString(title) {
		return RomanFormat
	}

	return UnknownFormat
}

// 判断是否是中文数字
func isChineseNumber(s string) bool {
	return regexp.MustCompile(`^[一二三四五六七八九十]+$`).MatchString(s)
}
