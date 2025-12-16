package interfaces

import (
	"context"

	"github.com/yb2020/odoc/pkg/eventbus"
	pb "github.com/yb2020/odoc/proto/gen/go/membership"
	"github.com/yb2020/odoc/services/membership/model"
)

// IMembershipService 会员服务接口
type IMembershipService interface {

	// GetBaseInfo 获取用户会员基础信息
	GetBaseInfo(ctx context.Context, userId string) (*pb.UserMembershipBaseInfo, error)

	// GetInfo 获取用户会员信息
	GetInfo(ctx context.Context, userId string) (*pb.UserMembershipInfo, error)

	// CheckAccountAndNew 检查用户会员账号并创建
	CheckAccountAndNew(ctx context.Context, userId string) (string, error)

	// NewMembershipAccount 创建用户会员账号并自动订阅Free会员
	NewMembershipAccount(ctx context.Context, userId string) (string, error)

	// DeleteMembershipAccount 删除用户会员账号
	DeleteMembershipAccount(ctx context.Context, userId string) error

	// AutoSubscribeMemberFree 自动订阅Free会员，无需支付自动完成订单状态
	AutoSubscribeMemberFree(ctx context.Context, userId string) error

	// CreditFunDocsUpload 调用积分接口上传接口功能
	// fileSize: 文件大小(byte)
	// filePageCount: 文件页数
	// storageCapacity: 已使用总存储容量(byte)
	// autoConfirm: 是否自动确认
	CreditFunDocsUpload(ctx context.Context, fileSize int64, filePageCount int32, storageCapacity int64, invokeFun func(xctx context.Context, sessionId string) error, autoConfirm bool) error

	// CreditFunAi 调用积分接口AI辅读功能
	// funType: AI辅读类型
	// modelKey: AI模型key
	// invokeFun: 调用AI辅读接口
	// autoConfirm: 是否自动确认
	CreditFunAi(ctx context.Context, funType pb.CreditServiceType, modelKey string, invokeFun func(xctx context.Context, sessionId string) error, autoConfirm bool) error

	// CreditFunTranslate 调用翻译接口功能
	// funType: 翻译类型
	// filePageCount: 文件页数 ，只有全文翻译需要，其他类型可为0
	// invokeFun: 调用翻译接口
	// autoConfirm: 是否自动确认
	CreditFunTranslate(ctx context.Context, funType pb.CreditServiceType, filePageCount int32, invokeFun func(xctx context.Context, sessionId string) error, autoConfirm bool) error

	// CreditFunNote 调用笔记接口功能
	// funType: 笔记类型
	// invokeFun: 调用笔记接口
	// autoConfirm: 是否自动确认
	CreditFunNote(ctx context.Context, funType pb.CreditServiceType, invokeFun func(xctx context.Context, sessionId string) error, autoConfirm bool) error

	// ConfirmCreditFun 二段提交成功确认接口
	// sessionId: 会话ID
	ConfirmCreditFun(ctx context.Context, sessionId string) error

	// RetrieveCreditFun 二段提交失败确认接口
	// sessionId: 会话ID
	RetrieveCreditFun(ctx context.Context, sessionId string) error

	// RefundCreditFun 发起退款接口，用于积分接口功能执行失败后，发起退款，暂不实现功能
	// sessionId: 会话ID
	// RefundCreditFun(ctx context.Context, sessionId string) error

	// 处理用户通知事件（由user模块发布消息--》会员模块订阅处理）
	HandleUserNotifyEventHandler(ctx context.Context, eventType eventbus.EventType, userId string) error

	// GetAccountExpiredList 获取过期的用户会员列表
	GetAccountExpiredList(ctx context.Context, size int) ([]model.UserMembership, error)

	// GetFreeAccountExpiredList 获取过期的Free用户会员列表
	GetFreeAccountExpiredList(ctx context.Context, size int) ([]model.UserMembership, error)

	// GetProAccountExpiredList 获取过期的Pro用户会员列表
	GetProAccountExpiredList(ctx context.Context, size int) ([]model.UserMembership, error)

	// HandleAccountExpired 处理过期的用户会员, id为会员id
	HandleAccountExpired(ctx context.Context, id string, userId string) error
}
