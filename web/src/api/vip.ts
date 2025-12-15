import api from './axios';
import { SuccessResponse } from './type';
import {
  REQUEST_SERVICE_NAME_APP,
  REQUEST_SERVICE_NAME_PAY,
  REQUEST_SERVICE_NAME_TRANSLATE,
} from './const';
import { VipProfileResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/user/VipUserInterface';
import { GetActivitiesStatusResp } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipActivitiesInfo';
import {
  OcrRemainCountReq,
  OcrRemainCountResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/TranslateProto';

// 已弃用：/vip/profile接口已不再使用
export const getVipProfile = async () => {
  // 注释掉原始实现
  // try {
  //   const res = await api.get<SuccessResponse<VipProfileResponse>>(
  //     `${REQUEST_SERVICE_NAME_APP}/vip/profile`,
  //     {}
  //   );
  //   return res.data.data;
  // } catch (error) {
  //   console.warn('Failed to fetch VIP profile, returning mock data', error);
  //   // 返回模拟数据以避免开发环境中的错误
  //   return {
  //     vipRoles: [{
  //       vipPrivilege: {
  //         vipType: 1, // VipType.FREE
  //         documentCountLimit: 100,
  //         exportEnable: false,
  //         fullTextTranslationCountLimit: 0,
  //         wordSelectionTranslateCountLimit: 10,
  //         ocrTranslateCountLimit: 0,
  //         aiBeanCountLimit: 0,
  //         aiTranslateEnable: false,
  //         readingEnable: false,
  //         polishEnable: false
  //       }
  //     }]
  //   } as VipProfileResponse;
  // }
  
  console.log('vip.ts: Failed to fetch VIP profile, returning mock data');
  // 返回模拟数据
  return {
    vipType: 0, // FREE
    leftDays: 0,
    expireTime: '',
    vipPrivileges: [],
    vipRoles: [{
      vipPrivilege: {
        vipType: 0, // VipType.FREE
        documentCountLimit: 100,
        exportEnable: false,
        fullTextTranslationCountLimit: 0,
        wordSelectionTranslateCountLimit: 10,
        ocrTranslateCountLimit: 0,
        aiBeanCountLimit: 0,
        aiTranslateEnable: false,
        readingEnable: false,
        polishEnable: false
      }
    }],
    vipProducts: []
  } as VipProfileResponse;
}

export const getActivityCardsStatus = async () => {
  const res = await api.get<SuccessResponse<GetActivitiesStatusResp>>(
    `${REQUEST_SERVICE_NAME_PAY}/vip/activities/getActivitiesStatus`,
    {}
  );
  return res.data.data;
};

// 已弃用：/ts/ocr/remainCount接口已不再使用
// export const postOcrRemainCount = async (params: OcrRemainCountReq) => {
//   const { data: res } = await api.post<SuccessResponse<OcrRemainCountResponse>>(
//     REQUEST_SERVICE_NAME_TRANSLATE + '/ts/ocr/remainCount',
//     params
//   );
//
//   return res.data.remainCount || 0;
// };
