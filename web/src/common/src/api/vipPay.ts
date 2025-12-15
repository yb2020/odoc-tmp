import {
  GetPayInfoReq,
  GetPayInfoResp,
  GetVipPrivilegeReq,
  GetVipPrivilegeResp,
  VipPayType,
  GetGroupBuyScanQRCodeReq,
  GetGroupBuyScanQRCodeResp,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { VipType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/user/VipUserInterface';
import {
  GetScanQRCodeReq,
  GetScanQRCodeResp,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { PayStatus } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { invert } from 'lodash';
// import { setupCache } from 'axios-cache-adapter'
import { HEADER_CANCLE_AUTO_ERROR, REQUEST_SERVICE_NAME_PAY } from './const';
import { SuccessResponse } from './type';
import api from './axios';
import { h } from 'vue';

export { PayStatus };

export type VipTypePayable = Exclude<
  VipType,
  VipType.FREE | VipType.UNRECOGNIZED | VipType.ENTERPRISE
>;

export enum ProductType {
  VIP = 'vip',
  STD = 'std',
  PRO = 'pro',
  ENT = 'ent',
  OSD = 'osd',
  AIBEAN = 'aibean',
  QS_REVIEW = 'qs_review',
}

export const VipPayType2ProductType = {
  [VipPayType.STANDARD]: ProductType.STD,
  [VipPayType.PROFESSIONAL]: ProductType.PRO,
  [VipPayType.OUTSTANDING]: ProductType.OSD,
  [VipPayType.AIBEAN]: ProductType.AIBEAN,
  [VipPayType.PAPER_AI_REVIEW]: ProductType.QS_REVIEW,
  [VipPayType.UNRECOGNIZED]: ProductType.VIP,
};

export const VipType2ProductType = {
  [VipType.STANDARD]: ProductType.STD,
  [VipType.PROFESSIONAL]: ProductType.PRO,
  [VipType.ENTERPRISE]: ProductType.ENT,
  [VipType.OUTSTANDING]: ProductType.OSD,
};

export const ProductWordingsMap = {
  [ProductType.VIP]: 'ReadPaper会员',
};

export const PayStatusWordingsMap: {
  [x in PayStatus]?: string;
} = {
  [PayStatus.PAY_PRE]: '拉取支付中',
  [PayStatus.PAY_WAITING]: '等待支付中',
  [PayStatus.PAY_TIMEOUT]: '支付超时',
  [PayStatus.PAY_SUCCESS]: '支付完成',
  // 后台说预留
  [PayStatus.PAY_COMPLETED]: '支付完成',
  [PayStatus.PAY_FAIL]: '支付失败',
};

export const VipType2PayType = {
  [VipType.STANDARD]: VipPayType.STANDARD,
  [VipType.PROFESSIONAL]: VipPayType.PROFESSIONAL,
  [VipType.OUTSTANDING]: VipPayType.OUTSTANDING,
};

export const PayType2VipType = invert(VipType2PayType) as unknown as {
  [k in VipPayType]: VipType;
};

export const getPayOrigin = () => {
  // 由于无法配置多个支付域名，改用URL参数控制环境
  // let prefix = ''
  // if (import.meta.env.VITE_API_ENV === 'dev') {
  //   prefix = 'dev.'
  // }

  // if (import.meta.env.VITE_API_ENV === 'uat') {
  //   prefix = 'uat.'
  // }

  // 支付宝生活号配置的回调链接强制要求https
  // 微信公众号配置的支付域名必须是https
  // return `https://rp.hanijiankang.com`
  return 'https://pay.readpaper.com';
};

export function getPayURL(
  data: Record<string, string>,
  env = '',
  origin = getPayOrigin()
) {
  const params = new URLSearchParams(data);
  if (env) {
    params.append('env', env);
  }
  return `${origin}/pay?${params.toString()}`;
}

// const cache = setupCache({
//   maxAge: 1 * 60 * 1000,
// })

declare module '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo' {
  interface GetVipPrivilegeResp {
    payOrigin: string;
  }
}

// 已弃用：/pay/public/getVipPrivilege接口已不再使用
// let p: undefined | Promise<GetVipPrivilegeResp>;
// export const getVipPayConfig = async (params?: GetVipPrivilegeReq) => {
//   if (!p) {
//     p = api
//       .get<SuccessResponse<GetVipPrivilegeResp>>(
//         `${REQUEST_SERVICE_NAME_PAY}/pay/public/getVipPrivilege`,
//         {
//           params,
//           // adapter: cache.adapter,
//           headers: {
//             [HEADER_CANCLE_AUTO_ERROR]: 'true',
//           },
//         }
//       )
//       .then((res) => res.data.data)
//       .finally(() => {
//         setTimeout(() => {
//           p = undefined;
//         }, 5000);
//       });
//   }

//   try {
//     return await p;
//   } catch (error) {
//     console.warn('Failed to fetch VIP pay config, returning mock data', error);
//     // 返回模拟数据以避免开发环境中的错误
//     return {
//       env: 'dev',
//       enabled: true,
//       payByDialog: true,
//       payOrigin: '',
//     } as GetVipPrivilegeResp;
//   }
// };

// 提供一个模拟的实现，以避免代码报错
export const getVipPayConfig = async (params?: GetVipPrivilegeReq) => {
  console.warn('调用已弃用的getVipPayConfig函数');
  return {
    env: 'dev',
    enabled: true,
    payByDialog: true,
    payOrigin: '',
  } as GetVipPrivilegeResp;
};

export const getQrCodeParams = async (params: GetScanQRCodeReq) => {
  const { vipPayType, ...rest } = params;
  const res = await api.get<SuccessResponse<GetScanQRCodeResp>>(
    `${REQUEST_SERVICE_NAME_PAY}/pay/getScanQRCode`,
    {
      params: {
        ...rest,
        vipPayType: VipPayType[vipPayType],
      },
    }
  );
  return res.data.data;
};

export const getPayInfo = async (params: GetPayInfoReq) => {
  const res = await api.get<SuccessResponse<GetPayInfoResp>>(
    `${REQUEST_SERVICE_NAME_PAY}/public/pay/getPayInfo`,
    {
      params,
      // headers: {
      //   "gray-route-rule": "{\"version\":\"yibing\"}"
      // },
    }
  );
  return res.data.data;
};

export const GetGroupBuyScanQRCode = async (
  params: GetGroupBuyScanQRCodeReq
) => {
  const { vipPayType, groupBuyItem, ...rest } = params;
  const searchParams = new URLSearchParams({
    ...rest,
    vipPayType: VipPayType[vipPayType],
  });

  // 处理 groupBuyItem 数组
  groupBuyItem.forEach((item, index) => {
    searchParams.append(
      `groupBuyItem[${index}].buyVipType`,
      VipType[item.buyVipType]
    );
    searchParams.append(
      `groupBuyItem[${index}].vipAmount`,
      String(item.vipAmount)
    );
  });

  const res = await api.get<SuccessResponse<GetGroupBuyScanQRCodeResp>>(
    `${REQUEST_SERVICE_NAME_PAY}/pay/getGroupBuyScanQRCode`,
    {
      params: searchParams,
      // headers: {
      //   "gray-route-rule": "{\"version\":\"yibing\"}"
      // },
      paramsSerializer: {
        encode: (param: string) => param, // 防止二次编码
      },
    }
  );
  return res.data.data;
};
