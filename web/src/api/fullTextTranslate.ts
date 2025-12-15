import api from './axios';
import { SuccessResponse } from './type';
import {
  GetFullTextTranslateRightInfoResponse,
  GetTranslateStatusReq,
  GetTranslateStatusResponse,
  FullTextTranslateRequest,
  FullTextTranslateResponse,
} from 'go-sea-proto/gen/ts/translate/FullTextTranslate';
import { liteThrottle } from '@idea/aiknowledge-special-util/throttle';
import { ElementName, PageType, reportElementImpression } from './report';

export const getRightInfo = async () => {
  const res = await api.get<
    SuccessResponse<GetFullTextTranslateRightInfoResponse>
  >(`/fullTextTranslate/getRightInfo`, {});
  return res.data.data;
};

export const getTranslateStatus = async (params: GetTranslateStatusReq) => {
  const res = await api.get<SuccessResponse<GetTranslateStatusResponse>>(
    `/fullTextTranslate/getTranslateStatus`,
    {
      params,
    }
  );
  return res.data.data;
};

export const postTranslate = async (params: FullTextTranslateRequest) => {
  const res = await api.post<SuccessResponse<FullTextTranslateResponse>>(
    `/fullTextTranslate/translate`,
    params
  );
  return res.data.data;
};

export const reportTranslateExposure = liteThrottle(
  (reachLimit: boolean) => {
    reportElementImpression({
      page_type: PageType.note,
      type_parameter: 'none',
      element_name: reachLimit
        ? ElementName.upperPaperTranslatePopup
        : ElementName.upperPaperTranslateLimitPopup,
    });
  },
  1000,
  true
);
