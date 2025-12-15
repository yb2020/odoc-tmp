// TODO
// import * as SentryTypes from '@sentry/core'
import {
  AutoCpRequest,
  AutoCpResponse,
  CpPaperInfo,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/search/AutoCP';
import { api } from './axios';
import {
  UserCenterSearchRequest,
  UserCenterSearchResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/search/SearchNote';
import { SearchResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/doc/UserDoc';

import { DocDetailInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/client/ClientDoc';
import { FOSAdministrator } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/fieldOfStudy/FieldOfStudyDetail';
import { SuccessResponse } from './type';
import { $PaperDetail } from './paper';
import { normalizeDocDetailInfo } from './document';
import { HEADER_CANCLE_AUTO_ERROR, REQUEST_SERVICE_NAME_APP } from './const';

interface $PaperSearchReq {
  keywords?: string;
  page: number;
  pageSize: number;
  fosId?: string;
  searchType?: number;
  startYear?: number;
  endYear?: number;
  venues?: string[];
  authors?: string[];
  sortType?: string;
  searchId?: string;
  searchHasPublicPdf?: boolean;
}

export type $FosInfo = {
  fosId?: string;
  fosName?: string;
  relatedPaperCount?: number;
  discussionCount?: number;
  subscribeCount?: number;
  administrators?: FOSAdministrator[];
};

type $PaperSearchRsp = SuccessResponse<{
  list: $PaperDetail[];
  total: number;
  conditions?: { authors?: string[]; venues?: string[] };
  fosInfo?: $FosInfo;
}>;

export const fetchSearchPaperResult = async (
  param: $PaperSearchReq,
  $sentry?: any // TODO typeof SentryTypes
) => {
  const startTime = Date.now();
  const res = await api.post<$PaperSearchRsp>(
    `${REQUEST_SERVICE_NAME_APP}/aiKnowledge/paper/search`,
    param,
    {
      timeout: 60000,
    }
  );
  const { data } = res.data;
  const endTime = Date.now();
  if ($sentry) {
    $sentry.captureMessage('api-cost-time', {
      tags: {
        api_url: `${REQUEST_SERVICE_NAME_APP}/aiKnowledge/paper/search`,
        cost_time: endTime - startTime,
      },
    });
  }
  return data;
};

export const fetchRecommended = async (
  params: AutoCpRequest
): Promise<CpPaperInfo[]> => {
  try {
    const res = await api.post<SuccessResponse<AutoCpResponse>>(
      `${REQUEST_SERVICE_NAME_APP}/paper/search/auto_cp`,
      {
        ...params,
        cpType: 0,
      },
      {
        headers: {
          [HEADER_CANCLE_AUTO_ERROR]: true,
        },
      }
    );
    return res.data.data.data || [];
  } catch (error) {
    return [];
  }
};

export const getSearchMyNoteList = async (params: UserCenterSearchRequest) => {
  const res = await api.post<SuccessResponse<UserCenterSearchResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/search/userCenterSearch`,
    params
  );
  return res.data.data;
};

export const getSearchMyNoteListV2 = async (
  params: UserCenterSearchRequest
) => {
  const response = await api.get<SuccessResponse<SearchResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/userDoc/search`,
    {
      params,
      paramsSerializer(p) {
        const search = new URLSearchParams();
        Object.keys(p).forEach((key) => {
          if (p[key] instanceof Array) {
            p[key].forEach((item: any) => {
              search.append(key, item);
            });
          } else {
            search.append(key, p[key]);
          }
        });
        return String(search);
      },
    }
  );
  const { data } = response.data;
  data.items.forEach((item) => {
    if (!item.docInfo) {
      item.docInfo = {} as DocDetailInfo;
    }

    normalizeDocDetailInfo(item.docInfo);
  });
  return data;
};
