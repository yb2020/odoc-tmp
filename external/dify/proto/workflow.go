package proto

// WorkflowRunRequest 工作流运行请求
type WorkflowRunRequest struct {
	Inputs       map[string]interface{} `json:"inputs,omitempty"`
	Query        string                 `json:"query,omitempty"`
	User         string                 `json:"user"`
	ResponseMode string                 `json:"response_mode,omitempty"` // blocking 或 streaming
	Files        []string               `json:"files,omitempty"`
}

// NewWorkflowRunRequest 创建新的工作流运行请求
func NewWorkflowRunRequest(user string) *WorkflowRunRequest {
	return &WorkflowRunRequest{
		User: user,
	}
}

// WithQuery 设置查询
func (r *WorkflowRunRequest) WithQuery(query string) *WorkflowRunRequest {
	r.Query = query
	return r
}

// WithInputs 设置输入参数
func (r *WorkflowRunRequest) WithInputs(inputs map[string]interface{}) *WorkflowRunRequest {
	r.Inputs = inputs
	return r
}

// WithResponseMode 设置响应模式
func (r *WorkflowRunRequest) WithResponseMode(responseMode string) *WorkflowRunRequest {
	r.ResponseMode = responseMode
	return r
}

// WithFiles 设置文件列表
func (r *WorkflowRunRequest) WithFiles(files []string) *WorkflowRunRequest {
	r.Files = files
	return r
}

// WorkflowRunResponse 工作流运行响应
type WorkflowRunResponse struct {
	TaskID    string                 `json:"task_id"`
	Workflow  string                 `json:"workflow"`
	Outputs   map[string]interface{} `json:"outputs,omitempty"`
	Status    string                 `json:"status"`
	CreatedAt int64                  `json:"created_at"`
	EndedAt   int64                  `json:"ended_at,omitempty"`
}

// WorkflowRunLogResponse 工作流运行日志响应
type WorkflowRunLogResponse struct {
	Logs []WorkflowRunLog `json:"logs"`
}

// WorkflowRunLog 工作流运行日志
type WorkflowRunLog struct {
	NodeID    string                 `json:"node_id"`
	NodeType  string                 `json:"node_type"`
	NodeName  string                 `json:"node_name"`
	Status    string                 `json:"status"`
	Inputs    map[string]interface{} `json:"inputs,omitempty"`
	Outputs   map[string]interface{} `json:"outputs,omitempty"`
	Error     string                 `json:"error,omitempty"`
	StartedAt int64                  `json:"started_at"`
	EndedAt   int64                  `json:"ended_at,omitempty"`
}

// WorkflowStopResponse 停止工作流响应
type WorkflowStopResponse struct {
	Result string `json:"result"`
}

// WorkflowRunStatusResponse 工作流运行状态响应
type WorkflowRunStatusResponse struct {
	TaskID    string                 `json:"task_id"`
	Workflow  string                 `json:"workflow"`
	Outputs   map[string]interface{} `json:"outputs,omitempty"`
	Status    string                 `json:"status"`
	CreatedAt int64                  `json:"created_at"`
	EndedAt   int64                  `json:"ended_at,omitempty"`
}

// WorkflowLogsResponse 工作流日志响应
type WorkflowLogsResponse struct {
	Data  []WorkflowLogItem `json:"data"`
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}

// WorkflowLogItem 工作流日志项
type WorkflowLogItem struct {
	ID        string `json:"id"`
	TaskID    string `json:"task_id"`
	Workflow  string `json:"workflow"`
	User      string `json:"user"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
	EndedAt   int64  `json:"ended_at,omitempty"`
}
