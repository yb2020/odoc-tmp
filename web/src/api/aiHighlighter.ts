import api from './axios';
import {
  HEADER_CANCLE_AUTO_ERROR,
  REQUEST_SERVICE_NAME_AI_READING,
} from './const';
import { SuccessResponse } from './type';
import {
  GetSCIMConfigReq,
  GetSCIMConfigResp,
  GetSCIMResultReq,
  GetSCIMResultResp,
  StartSCIMReq,
  StartSCIMResp,
  SCIMItem,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/scim/AiSCIMInfo';

export enum Status {
  INITIAL,
  LOADING,
  SUCCESS,
  FAILURE,
}

export type { SCIMItem as Highlight };

/**
 * 已弃用：/scim/getConfig接口已不再使用
 */
export const getRightInfo = async (params?: GetSCIMConfigReq) => {
  // 注释掉原始实现
  // const res = await api.get<SuccessResponse<GetSCIMConfigResp>>(
  //   `${REQUEST_SERVICE_NAME_AI_READING}/scim/getConfig`,
  //   params
  // );
  // return res.data.data;
  
  console.log('aiHighlighter.ts: getRightInfo called, but API is deprecated');
  // 返回模拟数据
  return {
    enabled: false,
    scimConfig: {
      scimEnabled: false,
      scimItems: []
    }
  } as GetSCIMConfigResp;
};

export const getHighlightResult = async (params: GetSCIMResultReq) => {
  const res = await api.get<SuccessResponse<GetSCIMResultResp>>(
    `${REQUEST_SERVICE_NAME_AI_READING}/scim/getProgress`,
    {
      params,
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: true,
      },
    }
  );
  return res.data.data;
};

export const startHighlight = async (params: StartSCIMReq) => {
  const res = await api.post<SuccessResponse<StartSCIMResp>>(
    `${REQUEST_SERVICE_NAME_AI_READING}/scim/start`,
    params
  );
  return res.data.data;
};
