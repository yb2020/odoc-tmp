package interfaces

import (
	"context"

	pb "github.com/yb2020/odoc-proto/gen/go/order"
	"github.com/yb2020/odoc/pkg/eventbus"
	"github.com/yb2020/odoc/services/membership/dto"
	"github.com/yb2020/odoc/services/membership/model"
)

// IOrderService 会员订单服务接口
type IOrderService interface {
	// GetById 根据ID获取订阅订单
	GetById(ctx context.Context, id string) (*model.Order, error)

	// GetOrderInfoById 根据ID获取订阅订单信息
	GetOrderInfoById(ctx context.Context, id string) (*pb.OrderInfo, error)

	// Subscribe 订阅会员订单
	Subscribe(ctx context.Context, userId string, orderType pb.OrderType, numberCount int32) (string, error)

	// NewOrder 创建订阅订单
	NewOrder(ctx context.Context, userId string, msId string, subInfo *dto.MembershipSubBaseInfo, numberCount int32) (string, error)

	// CancelOrder 取消订阅订单
	CancelOrder(ctx context.Context, orderId string) error

	// DoOrderPaySuccessHandler 订阅订单支付成功处理
	DoOrderPaySuccessHandler(ctx context.Context, orderId string, payOrderId string, stripeSubscriptionId string) error

	// DoOrderPayFailedHandler 订阅订单支付失败处理
	DoOrderPayFailedHandler(ctx context.Context, orderId string) error

	// DoOrderPayExpireHandler 订阅订单支付过期处理
	DoOrderPayExpireHandler(ctx context.Context, orderId string) error

	// HandlePayNotifyEventHandler 处理支付子系统的异步通知消息机制的支付结果回调事件
	HandlePayNotifyEventHandler(ctx context.Context, event eventbus.Event) error
}
