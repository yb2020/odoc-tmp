package interfaces

import (
	"context"

	pb "github.com/yb2020/odoc/proto/gen/go/doc"
	"github.com/yb2020/odoc/services/doc/model"
)

// IDocMetaInfoHandlerService 文档元信息处理器接口声明
type IDocMetaInfoHandlerService interface {
	// GetName 获取处理器的名称 对应移植的getType方法
	GetName() string

	// GetDocMetaInfo 获取文档元信息
	GetDocMetaInfo(ctx context.Context, req *pb.DocMetaInfoHandlerReq) (*pb.DocMetaInfoSimpleVo, error)

	// GetDocMetaInfosByTitle 获取文档元信息
	GetDocMetaInfosByTitle(ctx context.Context, title string) ([]*pb.DocMetaInfoSimpleVo, error)

	// InitMetaFieldFromContent 初始化元信息字段
	InitMetaFieldFromContent(fieldName string)

	// GetEmptyDocMetaInfo 获取空的文档元信息
	GetEmptyDocMetaInfo() *pb.DocMetaInfoSimpleVo

	// PersistentMetaInfo 持久化元信息
	PersistentMetaInfo(ctx context.Context, req *pb.DocMetaInfoHandlerReq) (*model.DoiMetaInfo, error)

	// ParseToMetaInfoVo 解析元信息
	ParseToMetaInfoVo(ctx context.Context, vo *pb.DocMetaInfoSimpleVo, metaInfoContent string)

	// UpdateDbMetaInfo 更新数据库中的元信息
	UpdateDbMetaInfo(ctx context.Context, vo *pb.DocMetaInfoSimpleVo, dbMetaInfo *model.DoiMetaInfo) error
}
