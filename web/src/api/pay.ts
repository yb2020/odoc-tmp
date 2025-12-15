import api from './axios';
import { SuccessResponse } from './type';
import { CreateCheckoutSessionReq, CreateCheckoutSessionResp, GetStripePublishableKeyReq, GetStripePublishableKeyResp } from 'go-sea-proto/gen/ts/pay/StripeCheckout';


export const getStripePublishableKey = async (params: GetStripePublishableKeyReq) => {
  const res = await api.get<SuccessResponse<GetStripePublishableKeyResp>>(
    `/pay/stripe/get-publishable-key`,
    params
  );
  return res.data.data;
};

export const createCheckoutSession = async (params: CreateCheckoutSessionReq) => {
    const res = await api.post<SuccessResponse<CreateCheckoutSessionResp>>(
      `/pay/stripe-checkout/create-session`,
      params,
    );
    return res.data.data;
};