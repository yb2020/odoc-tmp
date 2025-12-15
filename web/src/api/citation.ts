import {
  GetMyCslListResponse,
  GetCslListReq,
  GetCslListResponse,
  AddCslReq,
  DeleteCslReq,
  SortCslReq,
  CetDefaultCslListResponse,
  DocMetaInfoSimpleVo,
  GetDocMetaInfoReq,
  GetDocMetaInfoResponse,
  GetDocTypeListReq,
  GetDocTypeListResponse,
} from 'go-sea-proto/gen/ts/doc/CSL'
import {
  EnDocCiteSearchReq,
  EnDocCiteSearchResponse,
  ZhDocCiteSearchReq,
  ZhDocCiteSearchResponse,
} from 'go-sea-proto/gen/ts/doc/DocCiteSearch'
import {
  ManualUpdateDocCiteInfoReq,
  SearchUpdateDocCiteInfoReq,
  UpdateDocCiteInfoResponse,
} from 'go-sea-proto/gen/ts/doc/UserDocManage'
import { SuccessResponse } from './type';

import { api } from './axios';
import { reporter } from '@idea/aiknowledge-report';

// 获取我的格式列表
export const getMyCslList = async () => {
  const res = await api.post<SuccessResponse<GetMyCslListResponse>>(
    `/csl/myCslList`
  );

  return res.data.data;
};

// 获取全部格式列表带搜索
export const getCslList = async (params: GetCslListReq) => {
  const res = await api.post<SuccessResponse<GetCslListResponse>>(
    `/csl/list`,
    params
  );

  return res.data.data;
};

// 添加格式
export const addCsl = async (params: AddCslReq) => {
  const res = await api.post(
    `/csl/addCsl`,
    params
  );

  return res.data.data;
};

// 删除格式
export const deleteCsl = async (params: DeleteCslReq) => {
  const res = await api.post(
    `/csl/deleteCsl`,
    params
  );
  return res.data.data;
};

// 排序
export const sortCsl = async (params: SortCslReq) => {
  const res = await api.post(
    `/csl/sortCsl`,
    params
  );
  return res.data.data;
};

// 获取默认格式列表
export const getDefaultCslList = async () => {
  const res = await api.post<SuccessResponse<CetDefaultCslListResponse>>(
    `/docPublic/csl/getDefaultCslList`
  );

  return res.data.data;
};

export const getDocMetaInfo = async (
  params: GetDocMetaInfoReq
): Promise<DocMetaInfoWithVenue> => {
  const { data } = await api.post<GetDocMetaInfoResponse>(
    `/docPublic/docMetaInfo/getDocMetaInfo`,
    params
  )
  return withVenue(data.data as DocMetaInfoSimpleVo)
}


// 引用 -- 更新
export const enDocCiteSearch = async (params: EnDocCiteSearchReq) => {
  const { data: res } = await api.post<
    SuccessResponse<EnDocCiteSearchResponse>
  >(`/userDoc/citeSearch/en`, params, {
    timeout: 60000,
  })

  return res.data.result || []
}

export const zhDocCiteSearch = async (params: ZhDocCiteSearchReq) => {
  const { data: res } = await api.post<
    SuccessResponse<ZhDocCiteSearchResponse>
  >(`/userDoc/citeSearch/zh`, params, {
    timeout: 60000,
  })

  return res.data.result || []
}

export const searchUpdate = async (params: SearchUpdateDocCiteInfoReq) => {
  await api.post<UpdateDocCiteInfoResponse>(
    `/userDoc/searchUpdateDocCiteInfo`,
    params,
    {
      timeout: 60000,
    }
  )
}

export const createDefaultDocInfo = (): DocMetaInfoWithVenue => ({
  title: '',
  authorList: [],
  page: '',
  doi: '',
  volume: '',
  issue: '',
  partition: '',
  docTypeName: '',
  containerTitle: [],
  publishDateStr: '',
  publishTimestamp: '',
  docType: '',
  eventDate: [],
  venue: '',
})

export interface DocMetaInfoWithVenue extends DocMetaInfoSimpleVo {
  venue: string
}

export const withVenue = (data: DocMetaInfoSimpleVo): DocMetaInfoWithVenue => ({
  ...data,
  venue: data?.containerTitle?.[0] ?? '',
})

export const postManualUpdateDocCiteInfo = async (
  params: ManualUpdateDocCiteInfoReq
) => {
  const { data } = await api.post<SuccessResponse<void>>(
    `/userDoc/manualUpdateDocCiteInfo`,
    params
  )

  return data
}

export const getDocTypeList = async (params: GetDocTypeListReq) => {
  const { data } = await api.post<GetDocTypeListResponse>(
    `/docPublic/docMetaInfo/getDocTypeList`,
    params
  )

  return data.data.docTypeInfos
}

export enum MetaEventCode {
  readpaper_popup_paper_search_renew_click = 'readpaper_popup_paper_search_renew_click',
}

export const reportSearchRenewClick = (params: {
  page_type: string
  search_content: string
  language_type: 'chinese' | 'english'
  type: 'input' | 'title'
}) => {
  reporter.report(
    {
      event_code: MetaEventCode.readpaper_popup_paper_search_renew_click,
    },
    params
  )
}