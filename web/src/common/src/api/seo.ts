import api from '@common/api/axios';
import { SuccessResponse } from '@common/api/type';
import {
  HEADER_CANCLE_AUTO_ERROR,
  REQUEST_SERVICE_NAME_APP,
} from '@common/api/const';
// import { setupCache } from 'axios-cache-adapter'
// import { isDev } from '@common/utils/env'
import {
  GetTDKReq,
  GetTDKResp,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/seo/SEOTdkInfo';
import { SeoPageType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/seo/SEOCommon';

export interface TDKOptions {
  tplTitle: string;
  tplKeywords: string;
  tplDescription: string;
}

// const cache = setupCache({
//   maxAge: 5 * 60 * 1000,
// })

const formatTDKOptions = (data: GetTDKResp): TDKOptions => {
  return {
    tplTitle: data.titleTemplate,
    tplKeywords: data.keywordsTemplate,
    tplDescription: data.descriptionTemplate,
  };
};

export const getTKDOptions = async ({ type, ...rest }: GetTDKReq) => {
  try {
    const res = await api.get<SuccessResponse<GetTDKResp>>(
      `${REQUEST_SERVICE_NAME_APP}/aiKnowledge/seo/getTDK`,
      {
        params: {
          ...rest,
          type: SeoPageType[type],
        },
        // adapter: isDev ? undefined : cache.adapter,
        headers: {
          'x-custom-handle-error': HEADER_CANCLE_AUTO_ERROR,
        },
      }
    );
    const { data } = res.data;
    return formatTDKOptions(data);
  } catch (error) {
    return null;
  }
};

export type $SiteOptionData = {
  title: string;
  content: string;
  keywords: string;
  domain: string;
  copyright: string;
  contact: string;
  beianhao: string;
  icon: string;
  description: string;
  author: string;
  logo: string;
  staticDomain: string;
  metaCopyright: string;
};

export const getWebSiteInfo = async () => {
  try {
    const res = await api.post<SuccessResponse<$SiteOptionData>>(
      `${REQUEST_SERVICE_NAME_APP}/aiKnowledge/site/getSite`
    );
    const { data } = res.data;
    return data;
  } catch (error) {
    return null;
  }
};
