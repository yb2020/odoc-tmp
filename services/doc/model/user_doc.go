package model

import (
	"time"

	"github.com/yb2020/odoc/pkg/model"
)

// UserDoc 用户文献结构体
type UserDoc struct {
	model.BaseModel             // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	AppId             string    `json:"appId" gorm:"column:app_id;size:36"`                            // 应用ID
	UserId            string    `json:"userId" gorm:"column:user_id;index;not null;size:36"`           // 用户ID
	PaperId           string    `json:"paperId" gorm:"column:paper_id;index;size:36"`                  // 论文ID
	PdfId             string    `json:"pdfId" gorm:"column:pdf_id;index;size:36"`                      // PDF文件ID
	Sort              int       `json:"sort" gorm:"column:sort;type:int"`                              // 排序
	Remark            string    `json:"remark" gorm:"column:remark;type:text"`                         // 备注
	PaperTitle        string    `json:"paperTitle" gorm:"column:paper_title;type:varchar(255)"`        // 论文标题
	DocName           string    `json:"docName" gorm:"column:doc_name;type:varchar(255)"`              // 文档名称
	LastReadTime      time.Time `json:"lastReadTime" gorm:"index"`                                     // 最后阅读时间
	Venue             string    `json:"venue" gorm:"column:venue;type:varchar(255)"`                   // 发表场所
	PublishDate       string    `json:"publishDate" gorm:"column:publish_date;type:varchar(255)"`      // 发布日期
	AuthorDesc        string    `json:"authorDesc" gorm:"column:author_desc;type:text"`                // 作者描述
	GraphPartition    string    `json:"graphPartition" gorm:"column:graph_partition;type:text"`        // 图表分区
	VenueEdited       bool      `json:"venueEdited" gorm:"column:venue_edited;type:bool"`              // 发表场所是否已编辑
	PublishDateEdited bool      `json:"publishDateEdited" gorm:"column:publish_date_edited;type:bool"` // 发布日期是否已编辑
	AuthorDescEdited  bool      `json:"authorDescEdited" gorm:"column:author_desc_edited;type:bool"`   // 作者描述是否已编辑
	FillExtMeta       bool      `json:"fillExtMeta" gorm:"column:fill_ext_meta;type:bool"`             // 是否填充扩展元数据

	GraphAuthors     string `json:"graphAuthors" gorm:"column:graph_authors;type:text"`                  // 图表作者
	GraphPublishDate string `json:"graphPublishDate" gorm:"column:graph_publish_date;type:varchar(255)"` // 图表发布日期
	GraphVenues      string `json:"graphVenues" gorm:"column:graph_venues;type:text"`                    // 图表发表场所

	MetaAuthors     string `json:"metaAuthors" gorm:"column:meta_authors;type:text"`                  // 元数据作者
	MetaPublishDate string `json:"metaPublishDate" gorm:"column:meta_publish_date;type:varchar(255)"` // 元数据发布日期
	MetaVenues      string `json:"metaVenues" gorm:"column:meta_venues;type:text"`                    // 元数据发表场所

	DisplayAuthors     string `json:"displayAuthors" gorm:"column:display_authors;type:text"`                  // 显示作者
	DisplayPublishDate string `json:"displayPublishDate" gorm:"column:display_publish_date;type:varchar(255)"` // 显示发布日期
	DisplayVenues      string `json:"displayVenues" gorm:"column:display_venues;type:text"`                    // 显示发表场所

	MetaDoi        string `json:"metaDoi" gorm:"column:meta_doi;type:varchar(255)"`                // 元数据DOI
	MetaDocName    string `json:"metaDocName" gorm:"column:meta_doc_name;type:varchar(255)"`       // 元数据文档名称
	MetaPartition  string `json:"metaPartition" gorm:"column:meta_partition;type:text"`            // 元数据分区
	MetaDocType    string `json:"metaDocType" gorm:"column:meta_doc_type;type:varchar(255)"`       // 元数据文档类型
	MetaVolume     string `json:"metaVolume" gorm:"column:meta_volume;type:varchar(255)"`          // 元数据卷号
	MetaIssue      string `json:"metaIssue" gorm:"column:meta_issue;type:varchar(255)"`            // 元数据期号
	MetaPage       string `json:"metaPage" gorm:"column:meta_page;type:varchar(255)"`              // 元数据页码
	MetaLanguage   string `json:"metaLanguage" gorm:"column:meta_language;type:varchar(255)"`      // 元数据语言
	MetaUrl        string `json:"metaUrl" gorm:"column:meta_url;type:varchar(255)"`                // 元数据URL
	MetaEventTitle string `json:"metaEventTitle" gorm:"column:meta_event_title;type:varchar(255)"` // 元数据事件标题
	MetaEventPlace string `json:"metaEventPlace" gorm:"column:meta_event_place;type:varchar(255)"` // 元数据事件地点
	MetaEventDate  string `json:"metaEventDate" gorm:"column:meta_event_date;type:varchar(255)"`   // 元数据事件日期

	UserEditedDoi            string  `json:"userEditedDoi" gorm:"column:user_edited_doi;type:varchar(255)"`                  // 用户编辑的DOI
	UserEditedDocName        string  `json:"userEditedDocName" gorm:"column:user_edited_doc_name;type:varchar(255)"`         // 用户编辑的文档名称
	UserEditedPartition      string  `json:"userEditedPartition" gorm:"column:user_edited_partition;type:text"`              // 用户编辑的分区
	UserEditedDocType        string  `json:"userEditedDocType" gorm:"column:user_edited_doc_type;type:varchar(255)"`         // 用户编辑的文档类型
	UserEditedVolume         string  `json:"userEditedVolume" gorm:"column:user_edited_volume;type:varchar(255)"`            // 用户编辑的卷号
	UserEditedIssue          string  `json:"userEditedIssue" gorm:"column:user_edited_issue;type:varchar(255)"`              // 用户编辑的期号
	UserEditedPage           string  `json:"userEditedPage" gorm:"column:user_edited_page;type:varchar(255)"`                // 用户编辑的页码
	UserEditedImpactOfFactor float32 `json:"userEditedImpactOfFactor" gorm:"column:user_edited_impact_of_factor;type:float"` // 用户编辑的影响因子
	UserEditedJcrPartion     string  `json:"userEditedJcrPartion" gorm:"column:user_edited_jcr_partion;type:varchar(255)"`   // 用户编辑的JCR分区
	ImportanceScore          int     `json:"importanceScore" gorm:"column:importance_score;type:int"`                        // 重要性评分

	DoiEdited       bool `json:"doiEdited" gorm:"column:doi_edited;type:bool"`             // DOI是否已编辑
	DocNameEdited   bool `json:"docNameEdited" gorm:"column:doc_name_edited;type:bool"`    // 文档名称是否已编辑
	PartitionEdited bool `json:"partitionEdited" gorm:"column:partition_edited;type:bool"` // 分区是否已编辑
	DocTypeEdited   bool `json:"docTypeEdited" gorm:"column:doc_type_edited;type:bool"`    // 文档类型是否已编辑
	VolumeEdited    bool `json:"volumeEdited" gorm:"column:volume_edited;type:bool"`       // 卷号是否已编辑
	IssueEdited     bool `json:"issueEdited" gorm:"column:issue_edited;type:bool"`         // 期号是否已编辑
	PageEdited      bool `json:"pageEdited" gorm:"column:page_edited;type:bool"`           // 页码是否已编辑

	PaperRepositoryStatus string `json:"paperRepositoryStatus" gorm:"column:paper_repository_status;type:varchar(255)"` // 论文存储库状态
	NoteId                string `json:"noteId" gorm:"column:note_id;index;size:36"`                                    // 笔记ID
	NewPaper              bool   `json:"newPaper" gorm:"column:new_paper;type:bool"`                                    // 是否新论文

	ReadingStatus string `json:"readingStatus" gorm:"column:reading_status;type:varchar(255)"` // 阅读状态
	Progress      int    `json:"progress" gorm:"column:progress;type:int"`                     // 进度
	ParseStatus   int    `json:"parseStatus" gorm:"column:parse_status;type:int"`              // 解析状态
	// 不持久化字段（在Java中使用transient修饰）
	// 在Go中，我们不会添加gorm标签，这样它们就不会被持久化到数据库
	// OriginalImpactOfFactor float32 `json:"originalImpactOfFactor" gorm:"-"` // 原始影响因子
	// OriginalJcrPortion     string  `json:"originalJcrPortion" gorm:"-"`     // 原始JCR分区
}

// TableName 返回表名
func (UserDoc) TableName() string {
	return "t_user_doc"
}
