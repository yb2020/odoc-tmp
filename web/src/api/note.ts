import api from './axios';
import { REQUEST_SERVICE_NAME_APP } from './const';
import { SuccessResponse } from './type';
import { GetUserFolderTreeReq } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/protocol/doc/request/_GetUserFolderTreeReq';
import { GetUserFolderTreeResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/protocol/doc/response/_GetUserFolderTreeResponse';

// import {
  // ChangeWordConfigReq,
  // DeleteWordReq,
  // GetExtractListReq,
  // GetExtractListResponse,
  // GetSummaryListByFolderIdReq,
  // GetSummaryListByFolderIdResponse,
  // GetSummaryListReq,
  // GetSummaryListResponse,
  // GetWordListByFolderIdReq,
  // GetWordListByFolderIdResponse,
  // GetWordListReq,
  // GetWordListResponse,
  // GetWordsReq,
  // GetWordsResponse,
  // SaveOrUpdateSummaryReq,
  // SaveWordReq,
  // SaveWordResponse,
  // UpdateWordReq,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/NoteModule';
import { GetNoteSummaryByNoteIdRequest, GetNoteSummaryByNoteIdResponse, SaveOrUpdateSummaryReq } from 'go-sea-proto/gen/ts/note/NoteSummary'
import { 
  GetExtractListReq,
  GetExtractListResponse,
  GetSummaryListByFolderIdReq,
  GetSummaryListByFolderIdResponse,
  GetSummaryListReq,
  GetSummaryListResponse,
  GetWordListByFolderIdReq,
  GetWordListByFolderIdResponse,
  GetWordListReq, 
  GetWordListResponse 
} from 'go-sea-proto/gen/ts/note/NoteManage'
import { GetNoteWordsByNoteIdRequest, GetNoteWordsByNoteIdResponse, ChangeWordConfigRequest, DeleteNoteWordRequest, UpdateNoteWordRequest, SaveNoteWordRequest, SaveNoteWordResponse } from 'go-sea-proto/gen/ts/note/NoteWord'
// import {
//   // AddTagToAnnotateReq,
//   // CreateAnnotateTagReq,
//   // CreateAnnotateTagResp,
//   // DeleteAnnotateTagReq,
//   // DeleteAnnotateTagResp,
//   // DeleteTagToAnnotateReq,
//   // GetAnnotateTagsReq,
//   // GetAnnotateTagsResp,
//   // RenameAnnotateTagReq,
//   // RenameAnnotateTagResp,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Tag';
import { GetMarkTagListByFolderIdReq } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/protocol/doc/request/_GetMarkTagListByFolderIdReq';
import { GetMarkTagListByFolderIdResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/protocol/doc/response/_GetMarkTagListByFolderIdResponse';
import { GetMyNoteMarkListReq } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/protocol/note/request/_GetMyNoteMarkListReq';
import { GetMyNoteMarkListResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/protocol/note/response/_GetMyNoteMarkListResponse';
// import {
//   // UpdateNoteResp,
//   // WebNoteAnnotationModel,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Web';
import { WebNoteAnnotationModel } from 'go-sea-proto/gen/ts/note/Web'

import {AddTagToAnnotateRequest, CreateAnnotateTagRequest, CreateAnnotateTagResponse, DeleteAnnotateTagRequest, DeleteAnnotateTagResponse, DeleteTagToAnnotateRequest, GetAnnotateTagsRequest, GetAnnotateTagsResponse, RenameAnnotateTagRequest, RenameAnnotateTagResponse} from 'go-sea-proto/gen/ts/note/PaperNoteAnnotateTag';

// 从 profile.ts 导入，如果需要的话
export type $AuthorInfo = {
  id: string;
  name: string;
  avatar: string;
};

export type $PaperNoteItem = {
  noteId: string;
  likeCount: number;
  paperTitle: string;
  noteSummary: string;
  modifyDate: number;
  userInfo?: $AuthorInfo;
  paperNoteCount: number;
};

export interface NoteAnnotation extends WebNoteAnnotationModel {
  uuid: string;
}

export interface MyNoteMarkListResponse {
  total: number;
  annotationModelList: NoteAnnotation[];
}

export const getSummaryList = async (data: GetSummaryListReq = {}) => {
  const response = await api.post<SuccessResponse<GetSummaryListResponse>>(
    `/note/noteManage/summary/getList`,
    data
  );
  return response?.data?.data;
};

export const getSummaryListByFolderId = async (
  data: GetSummaryListByFolderIdReq
) => {
  const response = await api.post<
    SuccessResponse<GetSummaryListByFolderIdResponse>
  >(`/note/noteManage/summary/getListByFolderId`, data);
  return response?.data?.data;
};

export const getSummaryNote = async (params: GetNoteSummaryByNoteIdRequest) => {
  const { data: res } = await api.post<SuccessResponse<GetNoteSummaryByNoteIdResponse>>(
    `/note/paperNote/summary/getByNoteId`,
    params
  );

  return res.data;
};

export const updateSummaryNote = async (params: SaveOrUpdateSummaryReq) => {
  const { data: res } = await api.post<SuccessResponse<object>>(
    `/note/paperNote/summary/saveOrUpdate`,
    params
  );

  return res.data;
};

export const getWordList = async (data: GetWordListReq = {}) => {
  const response = await api.post<SuccessResponse<GetWordListResponse>>(
    `/note/noteManage/word/getList`,
    data
  );
  return response?.data?.data;
};

export const getWordListByFolderId = async (data: GetWordListByFolderIdReq) => {
  const response = await api.post<
    SuccessResponse<GetWordListByFolderIdResponse>
  >(`/note/noteManage/word/getListByFolderId`, data);
  return response?.data?.data;
};

export const getWordNotes = async (params: GetNoteWordsByNoteIdRequest) => {
  const { data: res } = await api.post<SuccessResponse<GetNoteWordsByNoteIdResponse>>(
    `/note/paperNote/word/getByNoteId`,
    params
  );

  return res.data;
};

export const addWordNote = async (params: SaveNoteWordRequest) => {
  const { data: res } = await api.post<SuccessResponse<SaveNoteWordResponse>>(
    `/note/paperNote/word/save`,
    params
  );

  return res.data;
};

export const updateWordNote = async (params: UpdateNoteWordRequest) => {
  const { data: res } = await api.post<SuccessResponse<object>>(
    `/note/paperNote/word/update`,
    params
  );

  return res.data;
};

export const delWordNote = async (params: DeleteNoteWordRequest) => {
  const { data: res } = await api.post<SuccessResponse<object>>(
    `/note/paperNote/word/delete`,
    params
  );

  return res.data;
};

export const setNoteWordsConfig = async (params: ChangeWordConfigRequest) => {
  const { data: res } = await api.post<SuccessResponse<object>>(
    `/note/paperNote/word/config`,
    params
  );

  return res.data;
};

export const getExtractList = async (data: GetExtractListReq = {}) => {
  const response = await api.post<SuccessResponse<GetExtractListResponse>>(
    `/note/noteManage/extract/getList`,
    data
  );
  return response?.data?.data;
};

export const getAnnotateTags = async (params: GetAnnotateTagsRequest) => {
  const response = await api.get<SuccessResponse<GetAnnotateTagsResponse>>(
    `/pdf/marktag/tags`,
    { params }
  );

  return response.data?.data?.tags ?? [];
};

export const getMarkTagListByFolderId = async (
  data: GetMarkTagListByFolderIdReq
) => {
  // const response = await api.post<
  //   SuccessResponse<GetMarkTagListByFolderIdResponse>
  // >(`${REQUEST_SERVICE_NAME_APP}/folder/getMarkTagListByFolderId`, data);
  // return response?.data?.data?.markTagList ?? [];
  const response = await api.post<
    SuccessResponse<GetMarkTagListByFolderIdResponse>
  >(`/note/noteManage/extract/getMarkTagListByFolderId`, data);
  return response?.data?.data?.markTagList ?? [];
};

export const getMyNoteMarkList = async (data: GetMyNoteMarkListReq) => {
  const response = await api.post<SuccessResponse<GetMyNoteMarkListResponse>>(
    `/pdf/pdfMark/v2/web/getMyNoteMarkList`,
    data
  );
  const { total = 0, annotationModelList = [] } = response?.data?.data ?? {};
  const withId: NoteAnnotation[] = annotationModelList.map((item) => ({
    uuid: item.select?.uuid ?? item.rect?.uuid ?? '',
    ...item,
    tags: item.tags || [],
  }));

  return {
    total,
    annotationModelList: withId,
  };
};

export const editAnnotation = (params: WebNoteAnnotationModel) => {
  // return api.put<SuccessResponse<UpdateNoteResp>>(
  //   `${REQUEST_SERVICE_NAME_APP}/pdfMark/v2/web/update`,
  //   params
  // );
  // return api.post<SuccessResponse<UpdateNoteResp>>(
  //   `/pdf/pdfMark/v2/web/update`,
  //   params
  // );
  return api.post(`/pdf/pdfMark/v2/web/update`, params);
};

export const deleteAnnotation = (annotationId: string) => {
  return api.post(`/pdf/pdfMark/v2/web/delete`, {
    id: annotationId
  });
};

export const  createAnnotationTag = async (params: CreateAnnotateTagRequest) => {
  const { data: res } = await api.post<SuccessResponse<CreateAnnotateTagResponse>>(
    `/pdf/marktag/save`,
    params
  );

  return res.data;
};

export const addTagToAnnotation = async (params: AddTagToAnnotateRequest) => {
  const { data: res } = await api.post(
    `/pdf/marktag/relation/mark/save`,
    params
  );

  return res;
};

export const deleteTagFromAnnotation = async (
  params: DeleteTagToAnnotateRequest
) => {
  const { data: res } = await api.post(
    `/pdf/marktag/relation/mark/delete`,
    params
  );

  return res;
};

export const paramsSerializer = (
  p: Record<string, number | string | (number | string)[]>
) => {
  const search = new URLSearchParams();

  Object.keys(p).forEach((key) => {
    const value = p[key];
    if (value instanceof Array) {
      value.forEach((item) => {
        search.append(key, String(item));
      });
    } else {
      search.append(key, String(value));
    }
  });

  return String(search);
};

export const renameAnnotateTag = (data: RenameAnnotateTagRequest) => {
  return api.post<SuccessResponse<RenameAnnotateTagResponse>>(
    `/pdf/marktag/update`,
    data
  );
};

export const deleteAnnotateTag = (data: DeleteAnnotateTagRequest) => {
  return api.delete<SuccessResponse<DeleteAnnotateTagResponse>>(
    `/pdf/marktag/delete`,
    { data }
  );
};

export const getUserFolderTree = async (data: GetUserFolderTreeReq = {}) => {
  const response = await api.post<SuccessResponse<GetUserFolderTreeResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/folder/getUserFolderTree`,
    data
  );
  return response?.data?.data;
};

export const getSummaryListByNoteId = async (data: GetSummaryListReq) => {
  const response = await api.post<SuccessResponse<GetSummaryListResponse>>(
    `/note/noteManage/summary/getList`,
    data
  );
  return response?.data?.data;
};
