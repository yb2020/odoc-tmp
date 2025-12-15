package proto

import "time"

// Annotation 标注
type Annotation struct {
	ID        string    `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AnnotationListResponse 标注列表响应
type AnnotationListResponse struct {
	Data  []Annotation `json:"data"`
	Total int          `json:"total"`
	Page  int          `json:"page"`
	Limit int          `json:"limit"`
}

// AnnotationReply 标注回复
type AnnotationReply struct {
	Status string `json:"status"`
	JobID  string `json:"job_id,omitempty"`
}
