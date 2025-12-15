package dto

// 定义扩展的GROBID解析模型，支持更丰富的元数据和位置信息

// BBox 表示文档中元素的边界框
type BBox struct {
	X0           float64 `json:"x0"`
	Y0           float64 `json:"y0"`
	X1           float64 `json:"x1"`
	Y1           float64 `json:"y1"`
	OriginHeight float64 `json:"originHeight"`
	OriginWidth  float64 `json:"originWidth"`
	Page         int     `json:"page,omitempty"` // 页码
}

// PageInfo 表示文档的页面信息
type PageInfo struct {
	PageNumber int     `json:"pageNumber"`
	Width      float64 `json:"width"`
	Height     float64 `json:"height"`
}

// Author 表示文档的作者信息
type Author struct {
	FullName  string `json:"full_name"`
	GivenName string `json:"given_name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	BBox      *BBox  `json:"bbox,omitempty"`
}

// Header 包含文档的头部信息
type Header struct {
	Authors []*Author `json:"authors"`
	Date    string    `json:"date"`
	Title   string    `json:"title"`
	ArxivID string    `json:"arxiv_id"`
}

// CatalogueItem 表示目录中的一个条目
type CatalogueItem struct {
	Title      string           `json:"title"`                // 章节标题
	PageNum    int              `json:"pageNum"`              // 章节所在页码
	Child      []*CatalogueItem `json:"child"`                // 子章节列表
	BBox       *BBox            `json:"bbox,omitempty"`       // 章节标题的边界框
	Level      string           `json:"level,omitempty"`      // 章节层级，如1、1.1、1.1.1等
	TitleOrder string           `json:"titleOrder,omitempty"` // 章节序号，如1、2、3等
	Order      int              `json:"order,omitempty"`      // 原始解析顺序
}

// Reference 表示文档中的引用
type Reference struct {
	Title       string    `json:"title"`       // 引用的标题
	PageNum     int       `json:"pageNum"`     // 引用所在页码
	BBox        *BBox     `json:"bbox"`        // 引用的边界框
	RefIdx      string    `json:"refIdx"`      // 引用的索引标识
	PublishDate string    `json:"publishDate"` // 发布日期
	SearchKey   string    `json:"searchKey"`   // 搜索关键词
	Authors     []*Author `json:"authors"`     // 作者列表
	ContentText string    `json:"contentText"` // 引用的原始文本内容
	ArxivID     string    `json:"arxiv_id"`    // arXiv ID
}

// RefMarker 表示文档中的引用标记
type RefMarker struct {
	RefIdx     string `json:"refIdx"`
	BBox       *BBox  `json:"bbox"`
	PaperID    int64  `json:"paperId"`
	RefContent string `json:"refContent"`
	PageNum    int    `json:"pageNum"`
}

// FigureTableMarker 表示文档中的图表标记
type FigureTableMarker struct {
	BBox       *BBox  `json:"bbox"`
	PageNum    int    `json:"pageNum"`
	RefContent string `json:"refContent"`
	RefIdx     string `json:"refIdx"`
}

// RefInfo 表示引用信息
type RefInfo struct {
	Text   string `json:"text"`           // 引用的显示文本，如 [13]
	Target string `json:"target"`         // 引用的目标ID，如 b12
	BBox   *BBox  `json:"bbox,omitempty"` // 引用的边界框
}

// Paragraph 表示文档中的一个段落
type Paragraph struct {
	Index        int       `json:"index"`                  // 段落在文档中的索引
	Text         string    `json:"text"`                   // 段落的文本内容
	Sentences    []string  `json:"sentences,omitempty"`    // 段落中的句子列表
	PageNum      int       `json:"pageNum"`                // 段落所在页码
	BBox         *BBox     `json:"bbox,omitempty"`         // 段落的边界框
	HasFormula   bool      `json:"hasFormula"`             // 是否包含公式
	HasTable     bool      `json:"hasTable"`               // 是否包含表格
	SectionTitle string    `json:"sectionTitle,omitempty"` // 所属章节标题
	SectionId    string    `json:"sectionId,omitempty"`    // 所属章节ID，用于关联到目录结构
	Type         string    `json:"type,omitempty"`         // 段落类型（正文、引用、标题等）
	References   []RefInfo `json:"references,omitempty"`   // 段落中引用的参考文献信息
	Formulas     []string  `json:"formulas,omitempty"`     // 段落中的公式文本
}

// ExtendedDocument 是扩展的文档结构，包含更多元数据和位置信息
type ExtendedDocument struct {
	GrobidVersion         string               `json:"grobid_version"`        // GROBID处理引擎的版本号
	GrobidTs              string               `json:"grobid_ts"`             // GROBID处理的时间戳
	Header                *Header              `json:"header"`                // 文档的头部信息，包含作者、标题等元数据
	PDFSHA256             string               `json:"pdfsha256"`             // 原始PDF文件的SHA256哈希值，用于唯一标识
	Lang                  string               `json:"lang"`                  // 文档的主要语言代码
	Abstract              string               `json:"abstract"`              // 文档的摘要内容
	Title                 string               `json:"title"`                 // 文档的标题
	Catalogue             []*CatalogueItem     `json:"catalogue"`             // 文档的目录结构信息，包含章节层次
	References            []*Reference         `json:"references"`            // 文档中的参考文献列表
	RefMarkers            []*RefMarker         `json:"refMarkers"`            // 文档中的引用标记，如[1]、[Smith et al.]等
	FigureAndTableMarkers []*FigureTableMarker `json:"figureAndTableMarkers"` // 文档中的图表标记，如Figure 1、Table 2等
	Pages                 []*PageInfo          `json:"pages,omitempty"`       // 文档的页面信息，包含每页的尺寸和编号
	Paragraphs            []*Paragraph         `json:"paragraphs,omitempty"`  // 文档的段落列表
}
