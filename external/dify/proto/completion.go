package proto

// CompletionRequest 完成请求
type CompletionRequest struct {
	Inputs       map[string]interface{} `json:"inputs,omitempty"`
	Query        string                 `json:"query,omitempty"`
	User         string                 `json:"user"`
	ResponseMode string                 `json:"response_mode,omitempty"` // blocking 或 streaming
	Files        []string               `json:"files,omitempty"`
}

// NewCompletionRequest 创建新的完成请求
func NewCompletionRequest(query, user string) *CompletionRequest {
	return &CompletionRequest{
		Query: query,
		User:  user,
	}
}

// WithInputs 设置输入参数
func (r *CompletionRequest) WithInputs(inputs map[string]interface{}) *CompletionRequest {
	r.Inputs = inputs
	return r
}

// WithResponseMode 设置响应模式
func (r *CompletionRequest) WithResponseMode(responseMode string) *CompletionRequest {
	r.ResponseMode = responseMode
	return r
}

// WithFiles 设置文件列表
func (r *CompletionRequest) WithFiles(files []string) *CompletionRequest {
	r.Files = files
	return r
}

// CompletionResponse 完成响应
type CompletionResponse struct {
	ID        string                 `json:"id"`
	Answer    string                 `json:"answer"`
	CreatedAt int64                  `json:"created_at"`
	TaskID    string                 `json:"task_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
