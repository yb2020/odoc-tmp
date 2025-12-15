package constant

const (

	// dify API 路径常量
	FILES_UPLOAD_PATH = "/files/upload"
	INFO_PATH         = "/info"
	PARAMETERS_PATH   = "/parameters"

	// datasets API 路径常量
	DATASETS_PATH                = "/datasets"
	DOCUMENTS_PATH               = "/documents"
	SEGMENTS_PATH                = "/segments"
	CHILD_CHUNKS_PATH            = "/child_chunks"
	DOCUMENT_CREATE_BY_TEXT_PATH = "/document/create-by-text"
	DOCUMENT_CREATE_BY_FILE_PATH = "/document/create-by-file"
	UPDATE_BY_TEXT_PATH          = "/update-by-text"
	UPDATE_BY_FILE_PATH          = "/update-by-file"
	INDEXING_STATUS_PATH         = "/indexing-status"
	UPLOAD_FILE_PATH             = "/upload-file"
	RETRIEVE_PATH                = "/retrieve"
	METADATA_PATH                = "/metadata"

	//启用/禁用内置元数据路径
	METADATA_BUILT_IN_PATH = "/metadata/built-in"
	//更新文档元数据路径
	DOCUMENT_METADATA_PATH = "/documents/metadata"
	//嵌入模型列表路径
	EMBEDDING_MODEL_TYPES_PATH = "/workspaces/current/models/model-types/text-embedding"

	// 流式响应相关常量
	DATA_PREFIX = "data:"
	PING_EVENT  = "event: ping"

	// chat API 路径常量
	// 对话型应用相关路径
	CHAT_MESSAGES_PATH       = "/chat-messages"
	MESSAGES_PATH            = "/messages"
	CONVERSATIONS_PATH       = "/conversations"
	AUDIO_TO_TEXT_PATH       = "/audio-to-text"
	TEXT_TO_AUDIO_PATH       = "/text-to-audio"
	META_PATH                = "/meta"
	STOP_PATH                = "/stop"
	FEEDBACKS_PATH           = "/feedbacks"
	SUGGESTED_QUESTIONS_PATH = "/suggested"
	NAME_PATH                = "/name"
	// 文本生成型应用相关路径
	COMPLETION_MESSAGES_PATH = "/completion-messages"

	// 工作流应用相关路径
	WORKFLOWS_PATH       = "/workflows"
	WORKFLOWS_RUN_PATH   = "/workflows/run"
	WORKFLOWS_TASKS_PATH = "/workflows/tasks"
	WORKFLOWS_LOGS_PATH  = "/workflows/logs"

	//标注应用相关路径
	APPS_ANNOTATIONS_PATH       = "/apps/annotations"
	APPS_ANNOTATIONS_REPLY_PATH = "/apps/annotations-reply"
)
