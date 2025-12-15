package constant

import pkgi18n "github.com/yb2020/odoc/pkg/i18n"

//这个里面的配置是固定的，后面看情况是否需要写入到配置文件中

// DefaultCslTitleOrder 默认的引用样式标题顺序
var DefaultCslTitleOrder = "BibTeX generic citation style;China National Standard GB/T 7714-2015 (numeric, 中文);Modern Language Association 9th edition;American Psychological Association 6th edition;Chicago Manual of Style 16th edition (author-date)"
var SPLIT_SEMICOLON = ";"
var CslDefaultLang = pkgi18n.LanguageEnUS
var CslDefaultLangMapRel = "{\"1676703778815535360\":\"" + pkgi18n.LanguageEnUS + "\"}"
var I18nFilterCslTitle = "China National Standard GB/T 7714-2015 (numeric, 中文)"

// DocTypeInfoListJsonDesc 文档类型信息列表的JSON描述
var DocTypeInfoListJsonDesc = `[
	{"code":"article-journal","name":"期刊论文"},
	{"code":"article-magazine","name":"杂志论文"},
	{"code":"article-newspaper","name":"报纸论文"},
	{"code":"thesis","name":"学位论文"},
	{"code":"paper-conference","name":"会议论文"},
	{"code":"book","name":"图书"},
	{"code":"chapter","name":"图书章节"},
	{"code":"patent","name":"专利"},
	{"code":"report","name":"报告"}
]`

// EnDocCiteSearchExtractPathSwitch 是否启用从URL中提取路径的开关，默认为true
var EnDocCiteSearchExtractPathSwitch = true
