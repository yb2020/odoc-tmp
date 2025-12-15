import { MessageTypeEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/operator/admin/msg/SystemMsgInfo';
import api from './axios';
import {
  getHeadersWithCancelAutoError,
  HEADER_CANCLE_AUTO_ERROR,
  REQUEST_SERVICE_NAME_APP,
} from './const';
import { SuccessResponse } from './type';

import { IsDisplayVipAdjustContentResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/MessageCenter';
import {
  GetUserRedDotListResponse,
  ClearRedDotReq,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/RedDot';

import { User } from 'go-sea-proto/gen/ts/user/User'
import { GetMyStatisticalReps } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/homePage/StatisticalInformation'

// import { GetTrafficReq, GetTrafficResp } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/ab/ABTestInfo';

export const getUserInfo = async () => {
  try {
    const res = await api.get<SuccessResponse<User>>(
      `/user/profile`,
      getHeadersWithCancelAutoError()
    );
    return res.data.data;
  } catch (err) {
    return null;
  }
};

export const doLogout = async () => {
  try {
    await api.post<SuccessResponse<null>>(
      `/oauth2/sign_out`
    );
    return true;
  } catch (error) {
    return false;
  }
};

export interface $RedDotInfo {
  msgCount: number;
  detailMessageCounts: {
    messageTypeEnum: MessageTypeEnum;
    count: number;
  }[];
}

export const getRedDotInfo = async () => {
  try {
    const res = await api.post<SuccessResponse<$RedDotInfo>>(
      `${REQUEST_SERVICE_NAME_APP}/messageCenter/getRedDotInfo`,
      {},
      getHeadersWithCancelAutoError()
    );
    return res.data.data;
  } catch (error) {
    return null;
  }
};

/**
 * 已弃用：/redDot/getUserRedDotList接口已不再使用
 */
export const getRedDotConfig = async () => {
  // 注释掉原始实现
  // try {
  //   const res = await api.post<SuccessResponse<GetUserRedDotListResponse>>(
  //     `${REQUEST_SERVICE_NAME_APP}/redDot/getUserRedDotList`,
  //     {},
  //     getHeadersWithCancelAutoError()
  //   );
  //   return res.data.data;
  // } catch (error) {
  //   return null;
  // }
  
  console.log('user.ts: getRedDotConfig called, but API is deprecated');
  // 返回模拟数据
  return {
    redDotList: []
  } as GetUserRedDotListResponse;
};

export const clearRedDot = async (params: ClearRedDotReq) => {
  try {
    const res = await api.post(
      `${REQUEST_SERVICE_NAME_APP}/redDot/clearRedDot`,
      params,
      getHeadersWithCancelAutoError()
    );

    return res;
  } catch (error) {
    return null;
  }
};

export const postVipAnnouncementVisible = async () => {
  const res = await api.post<
    SuccessResponse<IsDisplayVipAdjustContentResponse>
  >(
    `${REQUEST_SERVICE_NAME_APP}/isDisplayVipAdjustContent`,
    {},
    getHeadersWithCancelAutoError()
  );

  return res.data.data.isDisplay;
};

export const fetchPersonalLogout = async () => {
  try {
    await api.get<SuccessResponse<null>>(
      `/oauth2/sign_out`,
      getHeadersWithCancelAutoError()
    )
    return true
  } catch (error) {
    return false
  }
}

/**
 * 注释掉AB测试相关接口，避免404错误
 */
export const getTraffic = async (
  params: any // GetTrafficReq
) => {
  // try {
  //   const res = await api.get<SuccessResponse<GetTrafficResp>>(
  //     `${REQUEST_SERVICE_NAME_APP}/ab/get/traffic`,
  //     {
  //       params,
  //       ...getHeadersWithCancelAutoError()
  //     }
  //   )
  //   return res.data.data
  // } catch (error) {
  //   return null
  // }
  return null; // 直接返回null，避免调用实际接口
}

export const getMyStatistical = async (
)=> {
  const res = await api.get<SuccessResponse<GetMyStatisticalReps>>(
    `${REQUEST_SERVICE_NAME_APP}/personal/userCenter/getMyStatistical`
  )
  const { data } = res.data

  return data
}
