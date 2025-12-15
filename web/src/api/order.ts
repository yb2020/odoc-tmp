import api from './axios';
import { SuccessResponse } from './type';
import { CreateOrderRequest, CreateOrderResponse, GetOrderInfoRequest, GetOrderInfoResponse } from 'go-sea-proto/gen/ts/order/OrderApi';

export const createOrder = async (params: CreateOrderRequest) => {
    const res = await api.post<SuccessResponse<CreateOrderResponse>>(
      `/order/create-order`,
      params,
    );
    return res.data.data;
};

export const getOrderInfo = async (params: GetOrderInfoRequest) => {
    const res = await api.get<SuccessResponse<GetOrderInfoResponse>>(
      `/order/get-order-info`,
      params
    );
    return res.data.data;
};