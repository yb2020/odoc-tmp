package proto

// DatasetRequest 创建数据集请求
type DatasetRequest struct {
	// 数据集名称（必填）
	Name string `json:"name"`
	// 数据集描述（选填）
	Description string `json:"description,omitempty"`
	// 索引模式（选填，建议填写）
	// high_quality 高质量
	// economy 经济
	IndexingTechnique string `json:"indexing_technique,omitempty"`
	// 权限（选填，默认 only_me）
	// only_me 仅自己
	// all_team_members 所有团队成员
	// partial_members 部分团队成员
	Permission string `json:"permission,omitempty"`
	// Provider（选填，默认 vendor）
	// vendor 上传文件
	// external 外部知识库
	Provider string `json:"provider,omitempty"`
	// 外部知识库 API_ID（选填）
	ExternalKnowledgeApiId string `json:"external_knowledge_api_id,omitempty"`
	// 外部知识库 ID（选填）
	ExternalKnowledgeId string `json:"external_knowledge_id,omitempty"`
	// 处理规则
	ProcessRule *ProcessRule `json:"process_rule,omitempty"`
	// 检索模式
	RetrievalModel *RetrievalModel `json:"retrieval_model,omitempty"`
	// Embedding 模型名称
	EmbeddingModel string `json:"embedding_model,omitempty"`
	// Embedding 模型供应商
	EmbeddingModelProvider string `json:"embedding_model_provider,omitempty"`
}

// 处理规则
type ProcessRule struct {
	// 清洗、分段模式
	// automatic 自动
	// custom 自定义
	Mode string `json:"mode"`
	// 自定义规则（自动模式下，该字段为空）
	Rules *Rules `json:"rules,omitempty"`
}

type Rules struct {
	// 预处理规则
	PreProcessingRules []*PreProcessingRule `json:"pre_processing_rules,omitempty"`
	// 分段规则
	Segmentation *Segmentation `json:"segmentation,omitempty"`
	// 父分段的召回模式
	// full-doc 全文召回
	// paragraph 段落召回
	ParentMode string `json:"parent_mode,omitempty"`
	// 子分段规则
	SubchunkSegmentation *SubchunkSegmentation `json:"subchunk_segmentation,omitempty"`
}

type PreProcessingRule struct {
	// 预处理规则的唯一标识符
	// remove_extra_spaces 替换连续空格、换行符、制表符
	// remove_urls_emails 删除 URL、电子邮件地址
	Id string `json:"id"`
	// 是否选中该规则
	Enabled bool `json:"enabled"`
}

type Segmentation struct {
	// 自定义分段标识符，目前仅允许设置一个分隔符。默认为 \n
	Separator string `json:"separator"`
	// 最大长度（token）默认为 1000
	MaxTokens int `json:"max_tokens"`
}

type SubchunkSegmentation struct {
	// 分段标识符，目前仅允许设置一个分隔符。默认为 ***
	Separator string `json:"separator"`
	// 最大长度 (token) 需要校验小于父级的长度
	MaxTokens int `json:"max_tokens"`
	// 分段重叠
	ChunkOverlap int `json:"chunk_overlap"`
}

// RetrievalModel 检索模型
type RetrievalModel struct {
	// 检索方法
	// hybrid_search 混合检索
	// semantic_search 语义检索
	// full_text_search 全文检索
	SearchMethod string `json:"search_method"`
	// 是否开启rerank
	RerankingEnable bool `json:"reranking_enable"`
	// rerank 模式 【返回结果时才有这个数据，提交时不要提交】
	// weighted_score 权重排名
	// reranking_mode rerank检索
	RerankingMode string `json:"reranking_mode,omitempty"`
	// rerank 模型配置
	RerankingModel *RerankingModel `json:"reranking_model,omitempty"`
	// 权重
	Weights *WeightedScore `json:"weights,omitempty"`
	// 召回条数
	TopK int `json:"top_k"`
	// 是否开启召回分数限制
	ScoreThresholdEnabled bool `json:"score_threshold_enabled"`
	// 召回分数限制
	ScoreThreshold float64 `json:"score_threshold,omitempty"`
}

// 权重
type WeightedScore struct {
	WeightType     string          `json:"weight_type"`
	KeywordSetting *KeywordSetting `json:"keyword_setting,omitempty"`
	VectorSetting  *VectorSetting  `json:"vector_setting,omitempty"`
}

// VectorSetting 向量权重设置
type VectorSetting struct {
	VectorWeight               float64 `json:"vector_weight"`
	EmbeddingModelName         string  `json:"embedding_model_name"`
	EmbeddingModelProviderName string  `json:"embedding_model_provider_name"`
}

// KeywordSetting 关键词权重设置
type KeywordSetting struct {
	KeywordWeight float64 `json:"keyword_weight"`
}

// RerankingModel rerank 模型配置
type RerankingModel struct {
	// rerank 模型的提供商
	RerankingProviderName string `json:"reranking_provider_name"`
	// rerank 模型的名称
	RerankingModelName string `json:"reranking_model_name"`
}

// DatasetResponse 数据集响应
type DatasetResponse struct {
	// 知识库ID
	ID string `json:"id"`
	// 知识库名称
	Name string `json:"name"`
	// 知识库描述
	Description string `json:"description"`
	// 提供者
	Provider string `json:"provider"` // vendor/external
	// 权限
	Permission string `json:"permission"`
	// 数据源类型
	DataSourceType string `json:"data_source_type"`
	// 索引技术
	IndexingTechnique string `json:"indexing_technique"`
	// 应用数量
	AppCount int `json:"app_count"`
	// 文档数量
	DocumentCount int `json:"document_count"`
	// 字数
	WordCount int `json:"word_count"`
	// 创建者
	CreatedBy string `json:"created_by"`
	// 创建时间
	CreatedAt int64 `json:"created_at"`
	// 更新者
	UpdatedBy string `json:"updated_by"`
	// 更新时间
	UpdatedAt int64 `json:"updated_at"`
	// Embedding模型
	EmbeddingModel string `json:"embedding_model,omitempty"`
	// Embedding模型提供商
	EmbeddingModelProvider string `json:"embedding_model_provider,omitempty"`
	// Embedding是否可用
	EmbeddingAvailable bool `json:"embedding_available"`
	// 检索模型
	RetrievalModelDict *RetrievalModel `json:"retrieval_model_dict,omitempty"`
	// 标签
	Tags []string `json:"tags,omitempty"`
	// 文档形式
	DocForm string `json:"doc_form,omitempty"`
	// 外部知识库信息
	ExternalKnowledgeInfo *ExternalKnowledgeInfo `json:"external_knowledge_info,omitempty"`
	// 外部检索模型
	ExternalRetrievalModel *ExternalRetrievalModel `json:"external_retrieval_model,omitempty"`
	// 部分成员列表
	PartialMemberList []string `json:"partial_member_list,omitempty"`
	// 元数据
	DocMetadata *[]DocumentMetadata `json:"doc_metadata,omitempty"`
}

// ExternalKnowledgeInfo 外部知识库信息
type ExternalKnowledgeInfo struct {
	ExternalKnowledgeId          string `json:"external_knowledge_id,omitempty"`
	ExternalKnowledgeApiId       string `json:"external_knowledge_api_id,omitempty"`
	ExternalKnowledgeApiName     string `json:"external_knowledge_api_name,omitempty"`
	ExternalKnowledgeApiEndpoint string `json:"external_knowledge_api_endpoint,omitempty"`
}

// ExternalRetrievalModel 外部检索模型
type ExternalRetrievalModel struct {
	TopK                  int     `json:"top_k,omitempty"`
	ScoreThreshold        float64 `json:"score_threshold,omitempty"`
	ScoreThresholdEnabled bool    `json:"score_threshold_enabled,omitempty"`
}

// DatasetsRequest 获取数据集列表请求
type DatasetsRequest struct {
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	Keyword    string   `json:"keyword,omitempty"`
	TagIds     []string `json:"tag_ids,omitempty"`
	IncludeAll bool     `json:"include_all,omitempty"`
}

// DatasetsResponse 数据集列表响应
type DatasetsResponse struct {
	Data    []DatasetResponse `json:"data"`
	HasMore bool              `json:"has_more"`
	Total   int               `json:"total"`
	Page    int               `json:"page"`
	Limit   int               `json:"limit"`
}

// DocumentCreateRequest 创建文档请求
type DocumentCreateRequest struct {
	DatasetID   string                 `json:"dataset_id"`
	Name        string                 `json:"name"`
	Content     string                 `json:"content,omitempty"`
	OriginalURL string                 `json:"original_url,omitempty"`
	MetaInfo    map[string]interface{} `json:"meta_info,omitempty"`
	Batch       string                 `json:"batch,omitempty"`
}

// DocumentCreateByTextRequest 通过文本创建文档请求
type DocumentCreateByTextRequest struct {
	// 文档名称【必填, 不会返回文件Id】
	Name string `json:"name"`

	// 文档内容【必填】
	Text string `json:"text"`

	// 索引方式【必填】
	// high_quality 高质量：使用 embedding 模型进行嵌入，构建为向量数据库索引
	// economy 经济：使用 keyword table index 的倒排索引进行构建
	IndexingTechnique string `json:"indexing_technique,omitempty"`

	// 索引内容的形式【必填】
	// text_model text 文档直接 embedding，经济模式默认为该模式
	// hierarchical_model parent-child 模式
	// qa_model Q&A 模式：为分片文档生成 Q&A 对，然后对问题进行 embedding
	DocForm string `json:"doc_form,omitempty"`

	// 在 Q&A 模式下，指定文档的语言，例如：English、Chinese
	DocLanguage string `json:"doc_language,omitempty"`

	// 处理规则【必填】
	ProcessRule *ProcessRule `json:"process_rule,omitempty"`

	// 检索模式【选填】
	RetrievalModel *RetrievalModel `json:"retrieval_model,omitempty"`

	// Embedding 模型名称【选填】
	EmbeddingModel string `json:"embedding_model,omitempty"`

	// Embedding 模型供应商【选填】
	EmbeddingModelProvider string `json:"embedding_model_provider,omitempty"`
}

// DocumentCreateByFileRequest 通过文件创建文档请求
type DocumentCreateByFileRequest struct {
	// 源文档ID【选填】
	OriginalDocumentId string `json:"original_document_id,omitempty"`
	// 索引方式【必填】
	// high_quality 高质量：使用 embedding 模型进行嵌入，构建为向量数据库索引
	// economy 经济：使用 keyword table index 的倒排索引进行构建
	IndexingTechnique string `json:"indexing_technique,omitempty"`

	// 索引内容的形式
	// text_model text 文档直接 embedding，经济模式默认为该模式
	// hierarchical_model parent-child 模式
	// qa_model Q&A 模式：为分片文档生成 Q&A 对，然后对问题进行 embedding
	DocForm string `json:"doc_form,omitempty"`

	// 文档元数据（如提供文档类型则必填）
	DocMetadata map[string]interface{} `json:"doc_metadata,omitempty"`

	// 在 Q&A 模式下，指定文档的语言，例如：English、Chinese
	DocLanguage string `json:"doc_language,omitempty"`

	// 处理规则
	ProcessRule *ProcessRule `json:"process_rule,omitempty"`

	// 检索模式
	RetrievalModel *RetrievalModel `json:"retrieval_model,omitempty"`

	// Embedding 模型名称
	EmbeddingModel string `json:"embedding_model,omitempty"`

	// Embedding 模型供应商
	EmbeddingModelProvider string `json:"embedding_model_provider,omitempty"`
}

// DocumentUpdateByTextRequest 通过文本更新文档请求
type DocumentUpdateByTextRequest struct {
	// 文档名称（选填）
	Name string `json:"name,omitempty"`

	// 文档内容
	Text string `json:"text,omitempty"`

	// 文档元数据（如提供文档类型则必填）
	DocMetadata map[string]interface{} `json:"doc_metadata,omitempty"`

	// 处理规则（选填）
	ProcessRule *ProcessRule `json:"process_rule,omitempty"`
}

// DocumentUpdateByFileRequest 通过文件更新文档请求
type DocumentUpdateByFileRequest struct {
	// 文档名称（选填）
	Name string `json:"name,omitempty"`

	// 文档内容（选填）
	Text string `json:"text,omitempty"`

	// 文档元数据（如提供文档类型则必填）
	DocMetadata map[string]interface{} `json:"doc_metadata,omitempty"`

	// 处理规则（选填）
	ProcessRule *ProcessRule `json:"process_rule,omitempty"`
}

// IndexingStatusResponse 文档嵌入状态响应
type IndexingStatusResponse struct {
	// 状态列表
	Data []IndexingStatus `json:"data"`
}

// 嵌入状态
type IndexingStatus struct {
	// 文档ID
	ID string `json:"id"`

	// 索引状态
	IndexingStatus string `json:"indexing_status"`

	// 处理开始时间
	ProcessingStartedAt float64 `json:"processing_started_at"`

	// 解析完成时间
	ParsingCompletedAt float64 `json:"parsing_completed_at"`

	// 清洗完成时间
	CleaningCompletedAt float64 `json:"cleaning_completed_at"`

	// 分段完成时间
	SplittingCompletedAt float64 `json:"splitting_completed_at"`

	// 完成时间
	CompletedAt float64 `json:"completed_at"`

	// 暂停时间
	PausedAt float64 `json:"paused_at"`

	// 错误信息
	Error string `json:"error"`

	// 停止时间
	StoppedAt float64 `json:"stopped_at"`

	// 已完成分段数
	CompletedSegments int `json:"completed_segments"`

	// 总分段数
	TotalSegments int `json:"total_segments"`
}

// Document 响应信息
type DocumentResponse struct {
	Document Document `json:"document"`
	Batch    string   `json:"batch"`
}

// 文档信息
type Document struct {
	// 文档ID
	ID string `json:"id"`

	// 位置
	Position int `json:"position"`

	// 数据源类型
	DataSourceType string `json:"data_source_type"`

	// 数据源信息
	DataSourceInfo DataSourceInfo `json:"data_source_info"`

	// 数据源详情
	DataSourceDetailDict DataSourceDetailInfo `json:"data_source_detail_dict"`

	// 知识库处理规则ID
	DatasetProcessRuleID string `json:"dataset_process_rule_id"`

	// 文档名称
	Name string `json:"name"`

	// 创建来源
	CreatedFrom string `json:"created_from"`

	// 创建者
	CreatedBy string `json:"created_by"`

	// 创建时间
	CreatedAt int64 `json:"created_at"`

	// 令牌数
	Tokens int `json:"tokens"`

	// 索引状态
	IndexingStatus string `json:"indexing_status"`

	// 错误信息
	Error string `json:"error"`

	// 是否启用
	Enabled bool `json:"enabled"`

	// 禁用时间
	DisabledAt float64 `json:"disabled_at"`

	// 禁用者
	DisabledBy string `json:"disabled_by"`

	// 是否归档
	Archived bool `json:"archived"`

	// 显示状态
	DisplayStatus string `json:"display_status"`

	// 字数
	WordCount int `json:"word_count"`

	// 命中次数
	HitCount int `json:"hit_count"`

	// 文档形式
	DocForm string `json:"doc_form"`

	// 文档元数据
	DocMetadata []DocumentMetadataDescription `json:"doc_metadata"`
}

// DataSourceInfo 数据源信息
type DataSourceInfo struct {
	// 上传文件ID
	UploadFileId string `json:"upload_file_id"`
}

type DataSourceDetailInfo struct {
	// 上传文件ID
	UploadFile DataSourceDetailDict `json:"upload_file"`
}

// DataSourceDetailDict 数据源详情
type DataSourceDetailDict struct {
	// ID
	ID string `json:"id"`
	// 名称
	Name string `json:"name"`
	// 大小
	Size int `json:"size"`
	// 扩展名
	Extension string `json:"extension"`
	// MIME类型
	MimeType string `json:"mime_type"`
	// 创建者
	CreatedBy string `json:"created_by"`
	// 创建时间
	CreatedAt float64 `json:"created_at"`
}

// DocumentMetadataDescription 文档元数据描述
type DocumentMetadataDescription struct {
	// ID
	ID string `json:"id"`
	// 名称
	Name string `json:"name"`
	// 类型
	Type string `json:"type"`
	// 值
	Value string `json:"value"`
}

// DocumentsResponse 文档列表响应
type DocumentsResponse struct {
	// 文档列表
	Data []Document `json:"data"`
	// 文档总数
	Total int `json:"total"`
	// 当前页码
	Page int `json:"page"`
	// 每页显示数量
	Limit int `json:"limit"`
	// 是否有更多
	HasMore bool `json:"has_more"`
}

// CreateSegmentsRequest 创建文档分段请求
type CreateSegmentsRequest struct {
	// 分段列表
	Segments []SegmentInfoRequest `json:"segments"`
}

// 分段信息
type SegmentInfoRequest struct {
	// 文本内容/问题内容
	Content string `json:"content"`
	// 答案内容
	Answer string `json:"answer"`
	// 关键字
	Keywords []string `json:"keywords"`
}

// SegmentResponse 文档段落响应
type SegmentResponse struct {
	// 分段列表
	Data SegmentInfo `json:"data"`
	// 文档形式
	DocForm string `json:"doc_form"`
}

// 分段信息
type SegmentInfo struct {
	// 分段ID
	ID string `json:"id"`
	// 位置
	Position int `json:"position"`
	// 文档ID
	DocumentID string `json:"document_id"`
	// 内容
	Content string `json:"content"`
	// 答案
	Answer string `json:"answer"`
	// 字数
	WordCount int `json:"word_count"`
	// 令牌数
	Tokens int `json:"tokens"`
	// 关键字
	Keywords []string `json:"keywords"`
	// 索引节点ID
	IndexNodeID string `json:"index_node_id"`
	// 索引节点哈希
	IndexNodeHash string `json:"index_node_hash"`
	// 命中次数
	HitCount int `json:"hit_count"`
	// 是否启用
	Enabled bool `json:"enabled"`
	// 禁用时间
	DisabledAt int64 `json:"disabled_at"`
	// 禁用者
	DisabledBy string `json:"disabled_by"`
	// 状态
	Status string `json:"status"`
	// 创建者
	CreatedBy string `json:"created_by"`
	// 创建时间
	CreatedAt int64 `json:"created_at"`
	// 索引时间
	IndexingAt int64 `json:"indexing_at"`
	// 完成时间
	CompletedAt int64 `json:"completed_at"`
	// 错误信息
	Error string `json:"error"`
	// 停止时间
	StoppedAt int64 `json:"stopped_at"`
}

// UpdateSegmentRequest 更新文档分段请求
type UpdateSegmentRequest struct {
	// 分段信息
	Segment UpdateSegmentInfo `json:"segment"`
}

// UpdateSegmentInfo 更新文档分段信息
type UpdateSegmentInfo struct {
	// 文本内容/问题内容
	Content string `json:"content"`

	// 答案内容
	Answer string `json:"answer"`

	// 关键字
	Keywords []string `json:"keywords"`

	// 是否启用
	Enabled bool `json:"enabled"`

	// 是否重新生成子分段
	RegenerateChildChunks bool `json:"regenerate_child_chunks"`
}

// SegmentListResponse 文档段落列表响应
type SegmentListResponse struct {
	// 分段列表
	Data []SegmentInfo `json:"data"`
	// 文档形式
	DocForm string `json:"doc_form"`
}

// SaveChildChunkRequest 保存子分段请求
type SaveChildChunkRequest struct {
	// 子分段信息
	Content string `json:"content"`
}

// ChildChunkResponse 子分段响应
type ChildChunkResponse struct {
	// 子分段信息
	Data ChildChunkResponseInfo `json:"data"`
}

// ChildChunkResponseInfo 子分段信息
type ChildChunkResponseInfo struct {
	// 子分段ID
	ID string `json:"id"`

	// 分段ID
	SegmentID string `json:"segment_id"`

	// 子分段内容
	Content string `json:"content"`

	// 字数
	WordCount int `json:"word_count"`

	// Token 数
	Tokens int `json:"tokens"`

	// 索引节点ID
	IndexNodeID string `json:"index_node_id"`

	// 索引节点哈希
	IndexNodeHash string `json:"index_node_hash"`

	// 状态
	Status string `json:"status"`

	// 创建者
	CreatedBy string `json:"created_by"`

	// 创建时间
	CreatedAt int64 `json:"created_at"`

	// 索引时间
	IndexingAt int64 `json:"indexing_at"`

	// 完成时间
	CompletedAt int64 `json:"completed_at"`

	// 错误
	Error string `json:"error"`

	// 停止时间
	StoppedAt int64 `json:"stopped_at"`
}

// ChildChunkListResponse 子分段列表响应
type ChildChunkListResponse struct {
	Data  []ChildChunkResponse `json:"data"`
	Total int                  `json:"total"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
}

// UpdateDocumentMetadataRequest 修改文档元数据请【赋值】求
type UpdateDocumentMetadataRequest struct {
	// 操作数据
	OperationData []OperationMetadataItem `json:"operation_data"`
}

// OperationMetadataItem 操作元数据项
type OperationMetadataItem struct {
	DocumentID   string     `json:"document_id"`
	MetadataList []Metadata `json:"metadata_list"`
}

// Metadata 元数据
type Metadata struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Name  string `json:"name"`
}

// UpdateDatasetMetadataRequest 更新元数据请求
type UpdateDatasetMetadataRequest struct {
	// 元数据名称，必填
	Name string `json:"name"`
}

// MetadataResponse 元数据响应
type MetadataResponse struct {
	// 元数据ID
	ID string `json:"id"`

	// 元数据类型
	Type string `json:"type"`

	// 元数据名称
	Name string `json:"name"`
}

// CreateMetadataRequest 创建元数据请求
type CreateMetadataRequest struct {
	// 元数据名称，必填，相当于字段名
	Name string `json:"name"`

	// 元数据类型，必填, 类型：string, number, time
	Type string `json:"type"`
}

type EmbeddingModelListResponse struct {
	// 数据
	Data []EmbeddingModel `json:"data"`
}

type EmbeddingModel struct {
	// 提供商
	Provider string `json:"provider"`

	// 标签
	Label Language `json:"label"`

	// 小图标
	IconSmall Language `json:"icon_small"`

	// 大图标
	IconLarge Language `json:"icon_large"`

	// 状态
	Status string `json:"status"`

	// 模型列表
	Models []Model `json:"models"`
}

type Language struct {
	// 简体中文
	ZhHans string `json:"zh_hans"`

	// 美国英语
	EnUS string `json:"en_us"`
}

type Model struct {
	// 模型名称
	Model string `json:"model"`

	// 标签
	Label Language `json:"label"`

	// 模型类型
	ModelType string `json:"model_type"`

	// 特征
	Features string `json:"features"`

	// 获取方式
	FetchFrom string `json:"fetch_from"`

	// 模型属性
	ModelProperties ModelProperties `json:"model_properties"`

	// 弃用
	Deprecated string `json:"deprecated"`

	// 状态
	Status string `json:"status"`

	// 负载均衡
	LoadBalancingEnabled string `json:"load_balancing_enabled"`
}

type ModelProperties struct {
	// 上下文大小
	ContextSize int64 `json:"context_size"`
}

type DocumentMetadataListResponse struct {
	// 文档元数据列表
	DocMetadata []DocumentMetadata `json:"doc_metadata"`
	// 是否启用内置字段
	BuiltInFieldEnabled bool `json:"built_in_field_enabled"`
}

// 文档元数据
type DocumentMetadata struct {
	// 元数据 ID
	ID string `json:"id"`

	// 元数据类型
	Type string `json:"type"`

	// 元数据名称
	Name string `json:"name"`

	// 元数据使用次数
	UseCount int `json:"use_count"`
}
