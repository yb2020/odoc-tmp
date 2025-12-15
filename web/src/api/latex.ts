import {
  GetFileIsNeedUploadReq,
  GetFileIsNeedUploadResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishPdf';
import { api } from '@common/api/axios';
import { REQUEST_SERVICE_NAME_AI_POLISH } from './const';
import {
  ParseLatexReq,
  ParseLatexResp,
  GetParseResultReq,
  GetParseResultResp,
  ModifyParseResultReq,
  ModifyParseResultResp,
  ParseFileRequest,
  ParseFileResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText';
import { SuccessResponse } from '@common/api/type';
import { HEADER_CANCLE_AUTO_ERROR } from '@common/api/const';

export const getFileIsNeedUpload = async (p: GetFileIsNeedUploadReq) => {
  const res = await api.post<GetFileIsNeedUploadResponse>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/ai/polish/token`,
    p
  );
  if (res.data.data) {
    return res.data.data;
  }
  throw Error(
    'getFileIsNeedUpload error:  empty AcquirePolicyCallbackInfoResponse '
  );
};

export const parseLatex = async (p: ParseLatexReq) => {
  const res = await api.get<SuccessResponse<ParseLatexResp>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/parselatex`,
    {
      params: p,
      timeout: 60000,
    }
  );
  return res.data.data;
};

export const parseFile = async (p: ParseFileRequest) => {
  const res = await api.post<SuccessResponse<ParseFileResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/parseFile`,
    p,
    {
      timeout: 60000,
    }
  );
  return res.data.data;
};

export const getParseResult = async (p: GetParseResultReq) => {
  const res = await api.get<SuccessResponse<GetParseResultResp>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/getParseResult`,
    {
      params: p,
    }
  );
  return res.data.data;
};

export const modifyParseResult = async (p: ModifyParseResultReq) => {
  const res = await api.post<SuccessResponse<ModifyParseResultResp>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/modifyParseResult`,
    p,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );
  return res.data.data;
};
