package constant

// Redis key 前缀常量
const (
	// ParsePDFStatusKeyPrefix PDF解析状态的Redis key前缀
	ParsePDFStatusKeyPrefix = "parse:pdf_status:header:"
	// 解析全文的redis key前缀
	ParsePDFFulltextKeyPrefix = "parse:pdf_status:fulltext:"
	// pdf解析过程中的oss信息
	ParsePDFOssInfoKeyPrefix = "parse:pdf_oss_info:"

	//pdf文件解析状态token
	ParsePDFStatusTokenKeyPrefix = "parse:pdf_status_token:"

	//下载文件相关的key
	DownloadFileKeyPrefix = "download:token:"
	//doi搜索结果
	DocSearchResultByTitleMd5KeyPrefix = "doc:doi_search_result_by_title_md5:"
)
