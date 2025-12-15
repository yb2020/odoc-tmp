import { SuccessResponse } from './type';
import { bridgeAdaptor } from '../adaptor/bridge';
import api from './axios';
import {
  HEADER_CANCLE_AUTO_ERROR,
  REQUEST_SERVICE_NAME_APP,
  REQUEST_SERVICE_NAME_TRANSLATE,
} from './const';
import {
  TencentTranslateReq,
  TencentTranslateResp,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/v1/TencentTranslateProto';
import {
  IdeaTranslateReq,
  IdeaTranslateResp,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/IdeaTranslateProto';
import { TranslateCorrectInfoReq } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/TranslateCorrectInfoProto';
import {
  GoogleClientEnableReq,
  GoogleClientEnableResp,
  GoogleTranslateReq,
  GoogleTranslateResp,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/GoogleTranslateProto';
import {
  GetTranslateTabListResp,
  OcrTranslateReq,
  TranslateReq,
  TranslateResp,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/TranslateProto';
import {
  AddInterfaceReq,
  GetInterfaceListResponse,
  DeleteInterfaceReq,
  TxConfig,
  AliConfig,
  GoogleConfig,
  DeepLConfig,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/CustomTranslateInterface';
import axios from 'axios';
import { Optional } from 'utility-types';
import { GetTranslateTabsResponse, TextTranslateRequest, TranslateResponse } from 'go-sea-proto/gen/ts/translate/TextTranslate';

import { TranslateChannel } from 'go-sea-proto/gen/ts/translate/TranslateEnum';

export enum TranslateTabKey {
  youdao = "YOUDAO",
  google = "GOOGLE",
  ai = "AI",
  customGoogle = "CUSTOM_GOOGLE",
  tencent = 'tencent',
  idea = 'idea',
  baidu = 'baidu',
  deepl = 'deepl',
  ali = 'ali',
  other = 'other',
}

export enum CustomTranslateChannel {
  /** TENCENT - 腾讯 */
  TENCENT = 0,
  /** ALI - 阿里 */
  ALI = 1,
  /** GOOGLE - Google */
  GOOGLE = 2,
  /** DEEPL - DeepL */
  DEEPL = 3,
  /** OPENAI - OpenAI */
  OPENAI = 4,
  UNRECOGNIZED = -1,
}

export type UniTranslateResp = Optional<
  Omit<TranslateResp, 'targetContent'>,
  'targetResp'
> & {
  targetContent: string;
};

interface FetchTranslate {
  type: TranslateTabKey;
  content: string;
  pdfId: string;
  param: unknown;
  useGlossary: boolean;
}

const translateApi = axios.create({
  timeout: 15000,
});

export const fetchTextTranslate = async (params: TextTranslateRequest) => {
  const { data: res } = await api.post<SuccessResponse<TranslateResponse>>(
    '/text/translate',
    params
  );

  return res.data;
};

export const fetchTenCentTranslate = async (params: TencentTranslateReq) => {
  const { data: res } = await api.post<SuccessResponse<TencentTranslateResp>>(
    REQUEST_SERVICE_NAME_APP + '/ts/v1/tencent/zh-cn',
    params
  );

  return res.data;
};

export const fetchIdeaTranslate = async (params: IdeaTranslateReq) => {
  const { data: res } = await api.post<SuccessResponse<IdeaTranslateResp>>(
    REQUEST_SERVICE_NAME_APP + '/ts/idea/zh-cn',
    params
  );

  return res.data;
};

export const fetchGoogleTranslateApi = async (params: GoogleTranslateReq) => {
  const res = await api.post<SuccessResponse<GoogleTranslateResp>>(
    REQUEST_SERVICE_NAME_APP + '/ts/google/zh-cn',
    params
  );
  return res.data.data;
};

export const fetchGoogleTranslate = async (params: string) => {
  const content = encodeURIComponent(params).replaceAll('%20', '+');

  try {
    const res = await translateApi.get(
      `https://translate.googleapis.com/translate_a/single?client=gtx&dt=t&sl=auto&tl=zh-CN&q=${content}`
    );
    const str = res.data?.[0]
      ?.map((item: string[]) => {
        return item?.[0] || '';
      })
      .join('');
    return str;
  } catch (error1) {
    try {
      const res = await translateApi.get(
        `https://translate.google.cn/m?q=${content}&tl=zh-CN&sl=auto`
      );
      const text = res.data;
      const m = text.match(/class="result-container">(.*?)</);
      return (m && m[1]) || '';
    } catch (error2) {
      throw error2;
    }
  }
};

export const fetchUniTranslate = async (params: TranslateReq) => {
  const { data: res } = await api.post<SuccessResponse<UniTranslateResp>>(
    REQUEST_SERVICE_NAME_TRANSLATE + '/ts/zh-cn',
    params
  );

  return res.data;
};

export const fetchTranslate = async (params: FetchTranslate) => {
  const { type, ...rest } = params;
  let res;
  switch (type) {
    case TranslateTabKey.youdao: {
      const textTranslateRequest = {
        text: rest.content,
        channel: TranslateChannel.YOUDAO,
        pdfId: rest.pdfId ? rest.pdfId : undefined,
        useGlossary: rest.useGlossary
      } as TextTranslateRequest;
      res = await fetchTextTranslate(textTranslateRequest);
      break;
    }
    case TranslateTabKey.google: {
      const textTranslateRequest = {
        text: rest.content,
        channel: TranslateChannel.GOOGLE,
        pdfId: rest.pdfId ? rest.pdfId : undefined,
        useGlossary: rest.useGlossary
      } as TextTranslateRequest;
      res = await fetchTextTranslate(textTranslateRequest);
      break;
    }
    case TranslateTabKey.customGoogle: {
      const req = JSON.parse(JSON.stringify(rest.param));
      const resp = await bridgeAdaptor.translateOnGoogle({
        text: rest.content,
        projectId: req.googleApiKey,
      });
      res = {
        targetContent: Array.isArray(resp) ? resp : [resp],
      };
      break;
    }
    case TranslateTabKey.tencent: {
      const req = rest.param as TxConfig;
      const resp = await bridgeAdaptor.translateOnTX({
        text: rest.content,
        secretId: req.txSecretId,
        secretKey: req.txSecretKey,
      });

      res = {
        targetContent: [resp],
      };
      break;
    }
    case TranslateTabKey.deepl: {
      const req = rest.param as DeepLConfig;
      const resp = await bridgeAdaptor.translateOnDeepl({
        text: rest.content,
        authKey: req.deepLKey as string,
        api: req.deepLApi as string,
      });

      res = {
        targetContent: [resp],
      };
      break;
    }

    case TranslateTabKey.ali: {
      const req = rest.param as AliConfig;
      const resp = await bridgeAdaptor.translateOnAli({
        text: rest.content,
        scene: req.aliInterfaceVersion as string,
        accessKeyId: req.aliAccessKeyId as string,
        accessKeySecret: req.aliAccessKeySecret as string,
      });

      res = {
        targetContent: [resp],
      };
      break;
    }

    case TranslateTabKey.idea:
    case TranslateTabKey.baidu:
      res = await fetchUniTranslate({
        ...rest,
        channel: type,
      });
      break;

    default:
  }

  if (!res) {
    throw new Error('翻译失败');
  }

  const { targetContent: s } = res;
  const targetContent = Array.isArray(s) ? s?.[0] : s;
  return {
    ...res,
    targetContent,
    requestId: res.requestId,
  } as UniTranslateResp;
};

export const correctTranslate = async (params: TranslateCorrectInfoReq) => {
  await api.post<SuccessResponse<null>>(
    REQUEST_SERVICE_NAME_APP + '/ts/correct',
    params
  );
};

export const enableGoogleTranslate = async (params: GoogleClientEnableReq) => {
  try {
    const res = await api.post<SuccessResponse<GoogleClientEnableResp>>(
      REQUEST_SERVICE_NAME_APP + '/ts/google/client/enable',
      params,
      {
        headers: {
          [HEADER_CANCLE_AUTO_ERROR]: true,
        },
      }
    );
    return !!res.data.data.isOpen;
  } catch (error) {
    return true;
  }
};

// 删除自定义翻译接口
export const DeleteCustomerTranslateTab = async (
  params: DeleteInterfaceReq
) => {
  const res = await api.post<SuccessResponse<null>>(
    `${REQUEST_SERVICE_NAME_TRANSLATE}/customTranslateInterface/delete`,
    params,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: true,
      },
      params: {},
    }
  );
  return res;
};

// 获取自定义翻译接口列表
export const getCustomerTranslateTabList = async () => {
  const res = await api.post<SuccessResponse<GetInterfaceListResponse>>(
    `${REQUEST_SERVICE_NAME_TRANSLATE}/customTranslateInterface/getList`,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: true,
      },
      params: {},
    }
  );
  return res.data.data;
};

export const AddAliTranslateType = async (params: AliConfig) => {
  const res = await addCustomerTranslateType({
    channel: CustomTranslateChannel.ALI,
    aliConfig: params,
  });
  return res;
};

export const AddTxTranslateType = async (params: TxConfig) => {
  const res = await addCustomerTranslateType({
    channel: CustomTranslateChannel.TENCENT,
    txConfig: params,
  });
  return res;
};

export const AddGoogleTranslateTpe = async (params: GoogleConfig) => {
  const res = await addCustomerTranslateType({
    channel: CustomTranslateChannel.GOOGLE,
    googleConfig: params,
  });
  return res;
};

export const AddDeeplTranslateType = async (params: DeepLConfig) => {
  const res = await addCustomerTranslateType({
    channel: CustomTranslateChannel.DEEPL,
    deepLConfig: params,
  });
  return res;
};

export const addCustomerTranslateType = async (params: AddInterfaceReq) => {
  const res = await api.post<SuccessResponse<null>>(
    `${REQUEST_SERVICE_NAME_TRANSLATE}/customTranslateInterface/add`,
    params,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: true,
      },
    }
  );
  return res.data.data;
};

export const getTranslateTabList = async () => {
    const res = await api.get<SuccessResponse<GetTranslateTabsResponse>>(
    `/text/getTranslateTabs`,
      {
        headers: {
          [HEADER_CANCLE_AUTO_ERROR]: true,
        },
        params: {},
      }
    );
    return res.data.data;
};

// type OcrExtractTextResponse = any;

// export const postOcrExtractText = async (params: OcrExtractTextReq) => {
//   const { data: res } = await api.post<SuccessResponse<OcrExtractTextResponse>>(
//     REQUEST_SERVICE_NAME_TRANSLATE + '/ts/ocr/extractText',
//     params
//   );

//   return res.data.text || '';
// };

export const postOcrZhCn = async (
  params: OcrTranslateReq
): Promise<UniTranslateResp> => {
  const { data: res } = await api.post<SuccessResponse<TranslateResp>>(
    '/text/ocr/translate',
    params,
    { timeout: 30000 }
  );

  const { targetContent, ...rest } = res.data;
  const content = Array.isArray(targetContent)
    ? targetContent?.[0]
    : targetContent;
  return {
    ...rest,
    targetContent: content,
  };
};
