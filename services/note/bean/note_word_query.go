package bean

/**
 * 笔记单词查询参数
 */
type NoteWordQuery struct {
	NoteId      string
	CurrentPage int
	PageSize    int
	MinLoadedId string
}
