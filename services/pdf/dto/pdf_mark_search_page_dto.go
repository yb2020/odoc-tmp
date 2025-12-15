package dto

// PdfMarkSearchPageDto PDF标记搜索分页DTO
type PdfMarkSearchPageDto struct {
	FolderId      string   `json:"folderId"`
	SearchContent string   `json:"searchContent"`
	SortType      int32    `json:"sortType"`
	TagIdList     []string `json:"tagIdList"`
	StyleIdList   []int64  `json:"styleIdList"`
	DocId         string   `json:"docId"`

	CurrentPage int32 `json:"currentPage"`
	PageSize    int32 `json:"pageSize"`
}
