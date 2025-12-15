package api

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/order"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	"github.com/yb2020/odoc/services/membership/interfaces"
)

type OrderApi struct {
	logger            logging.Logger
	tracer            opentracing.Tracer
	orderService      interfaces.IOrderService
	membershipService interfaces.IMembershipService
}

func NewOrderApi(logger logging.Logger, tracer opentracing.Tracer, orderService interfaces.IOrderService, membershipService interfaces.IMembershipService) *OrderApi {
	return &OrderApi{
		logger:            logger,
		tracer:            tracer,
		orderService:      orderService,
		membershipService: membershipService,
	}
}

// @api /api/order/get-order-info
// @method GET
// @apiDescription 获取订单信息
func (api *OrderApi) GetOrderInfo(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OrderApi.GetOrderInfo")
	defer span.Finish()

	req := &pb.GetOrderInfoRequest{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "bad request params")
		return
	}
	// 1.获取订单信息
	orderInfo, err := api.orderService.GetOrderInfoById(ctx, req.OrderId)
	if err != nil {
		response.Error(c, "get order info failed", nil)
		return
	}

	if orderInfo == nil {
		response.Error(c, "order not exists", nil)
		return
	}

	result := &pb.GetOrderInfoResponse{
		OrderInfo: orderInfo,
	}

	response.Success(c, "Success", result)
}

// @api /api/order/create-order
// @method POST
// @apiDescription 创建订单
func (api *OrderApi) CreateOrder(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OrderApi.CreateOrder")
	defer span.Finish()

	req := &pb.CreateOrderRequest{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "bad request params")
		return
	}

	// validate order type
	if req.OrderType == pb.OrderType_ORDER_TYPE_SUB_FREE {
		response.ErrorNoData(c, "can not subscribe free order type")
		return
	} else if req.OrderType == pb.OrderType_ORDER_TYPE_SUB_PRO {
		req.NumberCount = 1
	} else if req.OrderType == pb.OrderType_ORDER_TYPE_SUB_PRO_ADD_ON_CREDIT {
		if req.NumberCount <= 0 {
			response.ErrorNoData(c, "numbercount cannot be less than or equal to 0")
			return
		}
	} else {
		response.ErrorNoData(c, "unsupported order type")
		return
	}
	api.logger.Info("msg", "create order", "req", req)

	userId, _ := userContext.GetUserID(ctx)
	orderId, err := api.orderService.Subscribe(ctx, userId, pb.OrderType(req.OrderType), int32(req.NumberCount))
	if err != nil {
		response.Error(c, err.Error(), nil)
		return
	}
	resp := &pb.CreateOrderResponse{
		OrderId: orderId,
	}

	// 自动完成订单状态
	// if req.OrderType == uint32(constant.OrderType_SubFree) {
	// 	err = api.orderService.DoOrderPaySuccessHandler(ctx, orderId)
	// 	if err != nil {
	// 		response.Error(c, "auto complete order status failed", nil)
	// 		return
	// 	}
	// }

	response.Success(c, "Success", resp)
}
