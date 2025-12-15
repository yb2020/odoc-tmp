package bean

import (
	docpb "github.com/yb2020/odoc-proto/gen/go/doc"
)

// UserDocUpdateType 更新类型枚举
type UserDocUpdateType int

const (
	UpdateTypeAuthors UserDocUpdateType = iota
	UpdateTypePublishDate
	UpdateTypeVenue
	UpdateTypeJcrPartion
	UpdateTypeRemark
	UpdateTypeImpactFactor
	UpdateTypeImportanceScore
)

// 更新请求结构体
type UserDocUpdateRequest struct {
	UpdateType      UserDocUpdateType
	Authors         []*docpb.BaseAuthorInfo
	PublishDate     string
	Venue           string
	JcrPartion      string
	Remark          string
	ImpactFactor    *float32
	ImportanceScore int32
}

type UserDocUpdateResponse struct {
	AuthorsResponse          *docpb.UpdateAuthorsResponse
	PublishDateResponse      *docpb.UpdatePublishDateResponse
	VenueResponse            *docpb.UpdateVenueResponse
	JcrPartionUpdateResponse *docpb.JcrPartionUpdateResponse
	ImpactFactorResponse     *docpb.ImpactFactorResponse
}
