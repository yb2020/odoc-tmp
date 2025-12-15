import api from './axios';
import { SuccessResponse } from './type';
import { GetInfoResponse, GetBaseInfoResponse, GetSubPlanInfosResponse, GetMembershipProfileResponse, GetLoginPageConfigResponse, GetUserCreditResponse } from 'go-sea-proto/gen/ts/membership/MembershipApi';
import { getHeadersWithCancelAutoError } from './const';

export const getMembershipUserInfo = async () => {
  try {
    const res = await api.get<SuccessResponse<GetMembershipProfileResponse>>(
      `/membership/user/profile`,
      getHeadersWithCancelAutoError()
    );
    return res.data.data;
  } catch (err) {
    return null;
  }
};

export const getMembershipInfo = async () => {
  try {
    const res = await api.get<SuccessResponse<GetInfoResponse>>(
      `/membership/get-info`,
      getHeadersWithCancelAutoError()
    );
    return res.data.data;
  } catch (err) {
    return null;
  }
};

export const getMembershipBaseInfo = async () => {
  try {
    const res = await api.get<SuccessResponse<GetBaseInfoResponse>>(
      `/membership/get-base-info`,
      getHeadersWithCancelAutoError()
    );
    return res.data.data;
  } catch (err) {
    return null;
  }
};

export const getMembershipSubPlanInfos = async () => {
  try {
    const res = await api.get<SuccessResponse<GetSubPlanInfosResponse>>(
      `/public/membership/get-subplan-infos`,
      getHeadersWithCancelAutoError()
    );
    return res.data.data;
  } catch (err) {
    return null;
  }
};

export const getLoginPageConfig = async () => {
  try {
    const res = await api.get<SuccessResponse<GetLoginPageConfigResponse>>(
      `/public/membership/get-login-page-config`,
      getHeadersWithCancelAutoError()
    );
    return res.data.data;
  } catch (err) {
    return null;
  }
};

export const getUserCredit = async () => {
  try {
    const res = await api.get<SuccessResponse<GetUserCreditResponse>>(
      `/membership/get-user-credit`,
      getHeadersWithCancelAutoError()
    );
    return res.data.data;
  } catch (err) {
    return null;
  }
};