import {
  AiBeanCountResponse,
  BuyAiBeanResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { Optional } from 'utility-types';
import { SuccessResponse } from './type';
import {
  HEADER_CANCLE_AUTO_ERROR,
  REQUEST_SERVICE_NAME_ACL,
  REQUEST_SERVICE_NAME_PAY,
  REQUEST_SERVICE_NAME_AI_READING,
  REQUEST_SERVICE_NAME_AI_POLISH,
} from './const';
import api from './axios';
import {
  RefundAiBeanResp,
  RefundAiBeanReq as RefundAiBeanCopilotReq,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiAssistantReading/ReadingFeedbackInfo';
import { RefundAiBeanReq as RefundAiBeanPolishReq } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/polish/PolishFeedbackInfo';

// 注释掉getBeansInfo函数，因为/pay/bean/count接口已经取消
// export const getBeansInfo = async () => {
//   const res = await api.get<SuccessResponse<AiBeanCountResponse>>(
//     `${REQUEST_SERVICE_NAME_PAY}/pay/bean/count`,
//     {
//       headers: {
//         [HEADER_CANCLE_AUTO_ERROR]: true,
//       },
//     }
//   );
//   return res.data.data;
// };

export const getBeansBuyWays = async () => {
  const res = await api.get<SuccessResponse<BuyAiBeanResponse>>(
    `${REQUEST_SERVICE_NAME_PAY}/pay/bean/buy`,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: true,
      },
    }
  );
  return res.data.data;
};

export enum RefundReasonScene {
  Copilot = 'aibean_reasons_copilot',
  Paragraph = 'aibean_reasons_paragraph',
  Translation = 'aibean_reasons_translation',
  Polish = 'aibean_reasons_polish',
  Review = 'aibeans_reasons_review',
}

export interface RefundReasonInfo {
  id: string;
  idName: string;
  name: string;
  internationalize: string;
}

export const getBeansRefundReasons = async (scene: RefundReasonScene) => {
  const res = await api.get<SuccessResponse<Array<RefundReasonInfo>>>(
    `${REQUEST_SERVICE_NAME_ACL}/dic/public/app/getListByUniqueId`,
    {
      params: {
        appId: 'aiKnowledge',
        uniqueId: {
          [RefundReasonScene.Copilot]: 'ReadingRefundAibean',
          [RefundReasonScene.Review]: 'PolishRefundAibean',
          [RefundReasonScene.Polish]: 'PolishRefundAibean1',
          [RefundReasonScene.Paragraph]: 'PolishRefundAibean2',
          [RefundReasonScene.Translation]: 'PolishRefundAibean3',
        }[scene],
      },
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );

  return res.data.data;
};

export async function refundAiBeans(
  scene: RefundReasonScene,
  params:
    | Optional<RefundAiBeanCopilotReq, 'answerId'>
    | Optional<RefundAiBeanPolishReq, 'taskId' | 'bizType'>
) {
  const res = await api.post<SuccessResponse<Partial<RefundAiBeanResp>>>(
    `${
      scene === RefundReasonScene.Copilot
        ? REQUEST_SERVICE_NAME_AI_READING
        : REQUEST_SERVICE_NAME_AI_POLISH
    }/refund/aibean`,
    params
    // {
    //   headers: {
    //     [HEADER_CANCLE_AUTO_ERROR]: true,
    //   },
    // }
  );
  return res.data.data;
}
