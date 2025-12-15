import {
  // ChangeWordConfigReq,
  // DeleteWordReq,
  GetExtractListReq,
  GetExtractListResponse,
  GetSummaryListByFolderIdReq,
  GetSummaryListByFolderIdResponse,
  GetSummaryListReq,
  GetSummaryListResponse,
  GetWordListByFolderIdReq,
  GetWordListByFolderIdResponse,
  GetWordListReq,
  GetWordListResponse,
  // GetWordsReq,
  // GetWordsResponse,
  // SaveOrUpdateSummaryReq,
  // SaveWordReq,
  // SaveWordResponse,
  // UpdateWordReq,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/NoteModule';



import api from './axios';
import { REQUEST_SERVICE_NAME_APP } from './const';
import { SuccessResponse } from './type';
import { $AuthorInfo } from './profile';
import {
  AddTagToAnnotateReq,
  // CreateAnnotateTagReq,
  // CreateAnnotateTagResp,
  DeleteAnnotateTagReq,
  DeleteAnnotateTagResp,
  DeleteTagToAnnotateReq,
  // GetAnnotateTagsReq,
  // GetAnnotateTagsResp,
  RenameAnnotateTagReq,
  RenameAnnotateTagResp,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Tag';
import { GetMarkTagListByFolderIdReq } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/protocol/doc/request/_GetMarkTagListByFolderIdReq';
import { GetMarkTagListByFolderIdResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/protocol/doc/response/_GetMarkTagListByFolderIdResponse';
import { GetMyNoteMarkListReq } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/protocol/note/request/_GetMyNoteMarkListReq';
import { GetMyNoteMarkListResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/protocol/note/response/_GetMyNoteMarkListResponse';
// import {
//   // UpdateNoteResp,
// WebNoteAnnotationModel,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Web';
import { WebNoteAnnotationModel } from 'go-sea-proto/gen/ts/note/Web'
// import { SaveOrUpdateSummaryReq } from 'go-sea-proto/gen/ts/note/NoteSummary';

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

// export const getSummaryList = async (data: GetSummaryListReq = {}) => {
//   const response = await api.post<SuccessResponse<GetSummaryListResponse>>(
//     `${REQUEST_SERVICE_NAME_APP}/noteManage/summary/getList`,
//     data
//   );
//   return response?.data?.data;
// };

// export const getSummaryListByFolderId = async (
//   data: GetSummaryListByFolderIdReq
// ) => {
//   const response = await api.post<
//     SuccessResponse<GetSummaryListByFolderIdResponse>
//   >(`${REQUEST_SERVICE_NAME_APP}/noteManage/summary/getListByFolderId`, data);
//   return response?.data?.data;
// };

// export const updateSummaryNote = async (params: SaveOrUpdateSummaryReq) => {
//   const { data: res } = await api.post<SuccessResponse<object>>(
//     `${REQUEST_SERVICE_NAME_APP}/paperNote/summary/saveOrUpdate`,
//     params
//   );

//   return res.data;
// };

// export const getWordList = async (data: GetWordListReq = {}) => {
//   const response = await api.post<SuccessResponse<GetWordListResponse>>(
//     `${REQUEST_SERVICE_NAME_APP}/noteManage/word/getList`,
//     data
//   );
//   return response?.data?.data;
// };

// export const getWordListByFolderId = async (data: GetWordListByFolderIdReq) => {
//   const response = await api.post<
//     SuccessResponse<GetWordListByFolderIdResponse>
//   >(`${REQUEST_SERVICE_NAME_APP}/noteManage/word/getListByFolderId`, data);
//   return response?.data?.data;
// };

// export const getWordNotes = async (params: GetWordsReq) => {
//   const { data: res } = await api.post<SuccessResponse<GetWordsResponse>>(
//     `${REQUEST_SERVICE_NAME_APP}/paperNote/word/getByNoteId`,
//     params
//   );

//   return res.data;
// };

// export const addWordNote = async (params: SaveWordReq) => {
//   const { data: res } = await api.post<SuccessResponse<SaveWordResponse>>(
//     `${REQUEST_SERVICE_NAME_APP}/paperNote/word/save`,
//     params
//   );

//   return res.data;
// };

// export const updateWordNote = async (params: UpdateWordReq) => {
//   const { data: res } = await api.post<SuccessResponse<object>>(
//     `${REQUEST_SERVICE_NAME_APP}/paperNote/word/update`,
//     params
//   );

//   return res.data;
// };

// export const delWordNote = async (params: DeleteWordReq) => {
//   const { data: res } = await api.post<SuccessResponse<object>>(
//     `${REQUEST_SERVICE_NAME_APP}/paperNote/word/delete`,
//     params
//   );

//   return res.data;
// };

// export const setNoteWordsConfig = async (params: ChangeWordConfigReq) => {
//   const { data: res } = await api.post<SuccessResponse<object>>(
//     `${REQUEST_SERVICE_NAME_APP}/paperNote/word/config`,
//     params
//   );

//   return res.data;
// };

// export const getExtractList = async (data: GetExtractListReq = {}) => {
//   const response = await api.post<SuccessResponse<GetExtractListResponse>>(
//     `${REQUEST_SERVICE_NAME_APP}/noteManage/extract/getList`,
//     data
//   );
//   return response?.data?.data;
// };

// export const getAnnotateTags = async (params: GetAnnotateTagsReq) => {
//   const response = await api.get<SuccessResponse<GetAnnotateTagsResp>>(
//     `${REQUEST_SERVICE_NAME_APP}/paperNote/tags`,
//     { params }
//   );

//   return response.data?.data?.tags ?? [];
// };

// export const getMarkTagListByFolderId = async (
//   data: GetMarkTagListByFolderIdReq
// ) => {
//   const response = await api.post<
//     SuccessResponse<GetMarkTagListByFolderIdResponse>
//   >(`${REQUEST_SERVICE_NAME_APP}/folder/getMarkTagListByFolderId`, data);
//   return response?.data?.data?.markTagList ?? [];
// };

// export const getMyNoteMarkList = async (data: GetMyNoteMarkListReq) => {
//   const response = await api.post<SuccessResponse<GetMyNoteMarkListResponse>>(
//     `${REQUEST_SERVICE_NAME_APP}/pdfMark/v2/web/getMyNoteMarkList`,
//     data
//   );
//   const { total = 0, annotationModelList = [] } = response?.data?.data ?? {};
//   const withId: NoteAnnotation[] = annotationModelList.map((item) => ({
//     uuid: item.select?.uuid ?? item.rect?.uuid ?? '',
//     ...item,
//     tags: item.tags || [],
//   }));

//   return {
//     total,
//     annotationModelList: withId,
//   };
// };

// export const editAnnotation = (params: WebNoteAnnotationModel) => {
//   return api.put<SuccessResponse<UpdateNoteResp>>(
//     `${REQUEST_SERVICE_NAME_APP}/pdfMark/v2/web/update`,
//     params
//   );
// };

// export const deleteAnnotation = (annotationId: string) => {
//   return api.delete(`${REQUEST_SERVICE_NAME_APP}/pdfMark/v2/web/delete`, {
//     data: { id: annotationId },
//   });
// };

// export const createAnnotationTag = async (params: CreateAnnotateTagReq) => {
//   const { data: res } = await api.post<SuccessResponse<CreateAnnotateTagResp>>(
//     `${REQUEST_SERVICE_NAME_APP}/paperNote/tag`,
//     params
//   );

//   return res.data;
// };

// export const addTagToAnnotation = async (params: AddTagToAnnotateReq) => {
//   const { data: res } = await api.post(
//     `${REQUEST_SERVICE_NAME_APP}/paperNote/tag/relation/mark`,
//     params
//   );

//   return res;
// };

// export const deleteTagFromAnnotation = async (
//   params: DeleteTagToAnnotateReq
// ) => {
//   const { data: res } = await api.delete(
//     `${REQUEST_SERVICE_NAME_APP}/paperNote/tag/relation/mark`,
//     { params, paramsSerializer }
//   );

//   return res;
// };

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

// export const renameAnnotateTag = (data: RenameAnnotateTagReq) => {
//   return api.put<SuccessResponse<RenameAnnotateTagResp>>(
//     `${REQUEST_SERVICE_NAME_APP}/paperNote/tag`,
//     data
//   );
// };

// export const deleteAnnotateTag = (data: DeleteAnnotateTagReq) => {
//   return api.delete<SuccessResponse<DeleteAnnotateTagResp>>(
//     `${REQUEST_SERVICE_NAME_APP}/paperNote/tag`,
//     { data }
//   );
// };
