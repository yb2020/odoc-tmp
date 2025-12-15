import api from './axios';
import {
  getHeadersWithCancelAutoError,
  REQUEST_SERVICE_NAME_APP,
} from './const';
import { ResponseError, SuccessResponse } from './type';
// import {
//   GetOwnerPaperNoteBaseInfoReq,
// } from '../mocks/types-readpaper-proto-all';
import { GetPaperNoteBaseInfoByIdReq, PaperNoteBaseInfoResponse, GetOwnerPaperNoteBaseInfoReq } from 'go-sea-proto/gen/ts/note/PaperNote'
import { GetPaperVersionsRequest, GetPaperVersionsResponse } from 'go-sea-proto/gen/ts/paper/Paper'
import {
  SelectPdfRequest,
  SelectPdfResponse,
} from '../mocks/types-readpaper-proto-all';
import {
  CloseAccessReq,
  CloseAccessResponse,
  OpenAccessReq,
  OpenAccessResponse,
} from '../mocks/types-readpaper-proto-all';
import { GetPdfStatusInfoRequest, GetPdfStatusInfoResponse } from 'go-sea-proto/gen/ts/pdf/PaperPDF'
import { normalizeAssign } from '@idea/aiknowledge-special-util/normalize-assign';

export const getPdfStatusInfo = async (params: GetPdfStatusInfoRequest) => {
  // const response = await api.post<SuccessResponse<GetPdfStatusInfoResponse>>(
  //   `${REQUEST_SERVICE_NAME_APP}/aiKnowledge/pdf/getPdfStatusInfo/v2`,
  //   params
  // );
  const response = await api.post<SuccessResponse<GetPdfStatusInfoResponse>>(
    `/pdf/getPdfStatusInfo/v2`,
    params
  );
  const { data } = response.data;
  // normalizeAssign(data, {
  //   authPdfId: '',
  //   paperId: '',
  // }, false);
  return data;
};

// export interface GetOwnerPaperNoteBaseInfoReq {
//   pdfId: string;
//   noteId?: string;
// }

type GetOwnerPaperNoteBaseInfoRsp = SuccessResponse<{
  noteId: string;
  licenceType: string;
}>;

export const getOwnerPaperNoteBaseInfo = async (
  params: GetOwnerPaperNoteBaseInfoReq
) => {
  const res = await api.post<GetOwnerPaperNoteBaseInfoRsp>(
    `/note/paperNote/getOwnerPaperNoteBaseInfo`,
    params,
    getHeadersWithCancelAutoError()
  );

  return res.data.data;
};

export type NoteBaseInfo = Omit<
PaperNoteBaseInfoResponse,
  'userInfo' | 'isPrivatePaper'
> &
  Required<Pick<PaperNoteBaseInfoResponse, 'userInfo' | 'isPrivatePaper'>>;

export const getPaperNoteBaseInfoById = async (
  params: GetPaperNoteBaseInfoByIdReq
) => {
  // const res = await api.post<SuccessResponse<NoteBaseInfo>>(
  //   `${REQUEST_SERVICE_NAME_APP}/paperNote/getPaperNoteBaseInfoById`,
  //   params
  // );
  const res = await api.post<SuccessResponse<NoteBaseInfo>>(
    `/note/paperNote/getPaperNoteBaseInfoById`,
    params
  );
  if (res.data.data) {
    // res.data.data.modifyDate = convertDate(res.data.data.modifyDate) as any;
    normalizeAssign(res.data.data, {
      isPrivatePaper: false,
    }, false);
    return res.data.data;
  }

  throw new ResponseError({
    code: 400,
    message: 'Invalid noteId',
  });
};

export const getUserTags = async () => {
  const { data: res } = await api.get(
    `/pdf/marktag/tags`,
    {
      params: {
        onlyUsed: false,
      },
    }
  );

  return res.data.tags;
};

export const getPaperVersions = async (params: GetPaperVersionsRequest) => {
  const res = await api.get<SuccessResponse<GetPaperVersionsResponse>>(
    `/paper/versions`,
    {
      params,
    }
  );
  return res.data.data;
};

export const changeCurrentPDF = async (params: SelectPdfRequest) => {
  await api.post<SuccessResponse<SelectPdfResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/pdf/select`,
    params
  );
};

export const openAccess = async (params: OpenAccessReq) => {
  await api.post<SuccessResponse<OpenAccessResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/paperNote/openAccess`,
    params,
    getHeadersWithCancelAutoError()
  );
};

export const closeAccess = async (params: CloseAccessReq) => {
  await api.post<SuccessResponse<CloseAccessResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/paperNote/closeAccess`,
    params,
    getHeadersWithCancelAutoError()
  );
};
