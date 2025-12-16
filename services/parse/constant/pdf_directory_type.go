package constant

import (
	pkgi18n "github.com/yb2020/odoc/pkg/i18n"
	pb "github.com/yb2020/odoc/proto/gen/go/oss"
)

var (
	// ================= pdf 数据存储类型 =================
	// PdfOssTypeMetadata 元数据文件类型
	PdfOssTypeMetadata = pb.ParsedDataEnum_METADATA.String()

	// PdfOssTypeParagraphs 段落文件类型
	PdfOssTypeParagraphs = pb.ParsedDataEnum_PARAGRAPHS.String()

	// PdfOssTypeFullText 全文文件类型
	PdfOssTypeFullText = pb.ParsedDataEnum_FULLTEXT.String()

	// PdfOssTypeFigureAndTable 图片和表格文件类型
	PdfOssTypeFigureAndTable = pb.ParsedDataEnum_FIGURE_AND_TABLE.String()

	// PdfOssTypeMarkdown markdown文件类型
	PdfOssTypeMarkdown = pb.ParsedDataEnum_MARKDOWN.String()

	// PdfOssTypePageBlocks 页面块文件类型
	PdfOssTypePageBlocks = pb.ParsedDataEnum_PAGE_BLOCK.String()
)

// PDF 存储类型常量
const (
	// PdfOssTypeFigure 图片文件类型
	PdfOssTypeFigure = "figure"

	// PdfOssTypeTable 表格文件类型
	PdfOssTypeTable = "table"

	// PdfOssTypeReference 引用文件类型
	PdfOssTypeReference = "reference"

	// ================= pdf 目录分类 =================

	//上传的原始文件目录
	SourcePdfCatalog = "origin"

	//解析后的文件目录
	ParsedPdfCatalog = "parsed"

	//全文翻译后的结果
	FullTextCatalog = "full_text_translated"

	//用户
	UserCatalog = "user"

	// ================= json 文件后缀 =================

	//元数据文件后缀
	MetadataJsonSuffix = "_meta.json"

	//段落文件后缀
	ParagraphsJsonSuffix = "_paragraphs.json"

	//page blocks文件后缀
	PageBlocksJsonSuffix = "_page_blocks.json"

	//  ================= 解析版本 =================
	ParseVersion = "0.0.1"

	// ================= minerU文件后缀名 =================

	//minerU pdf文件后缀名
	MineruPdfSuffix = ".pdf"

	//minerU json文件后缀名
	MineruJsonSuffix = ".json"

	//minerU md文件后缀名
	MineruMdSuffix = ".md"

	// ================= minerU文件名称类型 =================
	//========== minerU json文件类型 ==========
	// json文件 content_list 类型
	MineruJsonContentList = "content_list"
	// json文件 middle类型
	MineruJsonMiddle = "middle"

	// json文件 model类型
	MineruJsonModel = "model"

	//========== minerU pdf文件类型 ==========
	// pdf文件 layout类型
	MineruPdfLayout = "layout"

	// pdf文件 origin类型
	MineruPdfOrigin = "origin"

	// pdf文件 spans类型
	MineruPdfSpans = "spans"
	//========== content_list 和 middle文件中的数据类型 ==========
	// content_list文件 text 类型   普通文本
	ContentListJsonText = "text"
	// content_list文件 image 类型   图片
	ContentListJsonImage = "image"
	// content_list文件 table 类型   表格
	ContentListJsonTable = "table"
	// content_list文件 equation 类型   公式
	ContentListJsonEquation = "equation"
	//middle文件 title 类型   标题
	MiddleBlockTypeTitle = "title"
	//middle文件 image and table 类型 的属性
	MiddleBlockImageAndTableBody = "_body"
	//middle文件 image and table 类型 的属性
	MiddleBlockCaption = "_caption"
	//middle文件 interline_equation 类型 的属性   这个属性对应的是content_list文件中的equation类型
	MiddleBlockInterlineEquation = "interline_equation"

	//========== 语言类型 ==========
	// 使用标准 RFC 5646 格式的语言常量
	// 中文（简体）
	LanguageZhCN = pkgi18n.LanguageZhCN // "zh-CN"
	// 英文（美国）
	LanguageEnUS = pkgi18n.LanguageEnUS // "en-US"

	// 向后兼容的简化格式（仅用于特定服务）
	LanguageZh = "zh"
	LanguageEn = "en"
)

// GetGrobidLanguageZhCN 获取 Grobid 服务的中文语言格式
func GetGrobidLanguageZhCN() string {
	return pkgi18n.GlobalConverter.GetGrobidLanguage(pkgi18n.LanguageZhCN)
}

// IsGrobidLanguageZhCN 判断是否为 Grobid 的中文语言格式
func IsGrobidLanguageZhCN(lang string) bool {
	return lang == GetGrobidLanguageZhCN()
}
