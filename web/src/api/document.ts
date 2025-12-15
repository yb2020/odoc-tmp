import { normalizeAssign } from '@idea/aiknowledge-special-util';
// import {
//   GetFolderTreeAndSimpleDocInfoListReq,
//   GetFolderTreeAndSimpleDocInfoListResponse,
//   DeleteFolderReq,
//   DeleteDocReq,
//   UpdateFolderReq,
//   RenameUserDocReq,
//   MoveFolderOrDocReq,
//   AddFolderReq,
//   MoveDocOrFolderToAnotherFolderReq,
//   CopyDocOrFolderToAnotherFolderReq,
//   GetDocListByFolderIdReq,
//   GetDocListByFolderIdResponse,
//   RemoveDocFromFolderReq,
//   DocDetailInfo,
//   GetTopLevelFolderListReq,
//   GetTopLevelFolderListResponse,
//   AttachPaperToFolderReq,
//   AttachPaperToFolderResponse,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/client/ClientDoc';
import {
  GetFolderTreeAndSimpleDocInfoListReq,
  GetFolderTreeAndSimpleDocInfoListResponse,
  DeleteFolderReq,
  DeleteDocReq,
  UpdateFolderReq,
  RenameUserDocReq,
  MoveFolderOrDocReq,
  AddFolderReq,
  MoveDocOrFolderToAnotherFolderReq,
  CopyDocOrFolderToAnotherFolderReq,
  GetDocListByFolderIdReq,
  GetDocListByFolderIdResponse,
  RemoveDocFromFolderReq,
  DocDetailInfo,
  GetTopLevelFolderListReq,
  GetTopLevelFolderListResponse,
  AttachPaperToFolderReq,
  AttachPaperToFolderResponse,
} from 'go-sea-proto/gen/ts/doc/ClientDoc';
import {
  GetUserDocFoldersRequest,
  GetUserDocFoldersResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/doc/UserDoc';
import {
  GetUserDocCreateStatusRequest,
  GetUserDocCreateStatusResponse,
  GetUploadTokenReq,
  GetUploadTokenResp,
  HandleFileFastUploadReq,
  HandleFileFastUploadResp,
  GetLatestReadDocListReq,
  GetLatestReadDocListResp,
  GetParseTokenResponse,
  UploadPdfByUrlLinkRequest,
} from 'go-sea-proto/gen/ts/doc/UserDoc';

import { 
  GetMinioUploadTokenRequest,
  GetMinioUploadTokenResponse,
 } from 'go-sea-proto/gen/ts/oss/Oss';

import { 
  GetDocIndexReq,
  GetDocIndexResponse,
  GetDocListReq,
  GetDocListResponse,
  GetDocRelatedAuthorListReq,
  GetDocRelatedAuthorListResponse,
  GetDocRelatedClassifyListReq,
  GetDocRelatedClassifyListResponse,
  GetDocRelatedVenueListReq,
  GetDocRelatedVenueListResponse,
  ImpactFactorReq,
  ImpactFactorResponse,
  ImportanceScoreReq,
  JcrPartionUpdateReq,
  JcrPartionUpdateResponse,
  JcrPartionsReq,
  JcrPartionsResponse,
  UpdateReadStatusRequest,
  UserDocInfo,
} from 'go-sea-proto/gen/ts/doc/UserDocManage';

import { SuccessResponse } from './type';
import { api } from './axios';
import {
  AddClassifyRequest,
  AddClassifyResponse,
  AttachDocToClassifyRequest,
  AttachDocToClassifyResponse,
  DeleteClassifyRequest,
  GetMyCollectedDocFolderResponse,
  GetUserAllClassifyListRequest,
  GetUserAllClassifyListResponse,
  RemoveDocFromClassifyRequest,
  UpdateDocRemarkRequest,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/userCenter/UserDoc';
// import {
//   ExportByFolderIdReq,
//   ExportByIdsReq,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/bibtex/BibtexInfo';
import {
  ExportByFolderIdReq,
  ExportByIdsReq,
} from 'go-sea-proto/gen/ts/doc/CSL';

export const removeDocFromFolder = async (params: RemoveDocFromFolderReq) => {
  try {
    await api.post<SuccessResponse<null>>(
      `/client/doc/removeDocFromFolder`,
      params
    );
    return true;
  } catch (error) {
    return false;
  }
};

export const normalizeDocDetailInfo = (item: DocDetailInfo) => {
  normalizeAssign(
    item,
    {
      docName: '',
      classifyInfos: [],
      displayAuthor: {
        authors: item.authors || [],
        originAuthors: item.authors || [],
        rollbackEnable: false,
      },
      displayVenue: {
        venue: '',
        originVenue: '',
        rollbackEnable: false,
      },
      displayPublishDate: {
        publishDate: '',
        originPublishDate: '',
        rollbackEnable: false,
      },
      publishDate: '0',
      remark: '',
    },
    false
  );
};

export const getDocListByFolderId = async (params: GetDocListByFolderIdReq) => {
  const res = await api.post<
    SuccessResponse<GetDocListByFolderIdResponse['data']>
  >(`/client/doc/getDocListByFolderId`, params);
  const list = res.data.data || [];
  list.forEach(normalizeDocDetailInfo);
  return list;
};

export const getDocListNew = async (params: GetDocListReq) => {
  try {
    const response = await api.post<SuccessResponse<GetDocListResponse>>(
      `/userDoc/getDocList`,
      params
    );
    const { docList = [], total = 0 } = response.data.data;
    docList.forEach(normalizeDocDetailNew);
    return {
      docList: docList as Required<UserDocInfo>[],
      total,
    };
  } catch (error) {
    console.error('[ERROR] getDocListNew - 错误详情:', JSON.stringify(error));
    throw error;
  }
};

export const getLatestReadDocListNew = async (params: GetLatestReadDocListReq) => {
  try {
    const response = await api.post<SuccessResponse<GetLatestReadDocListResp>>(
      `/userDoc/getLatestReadDocList`,
      params
    );
    return response.data.data;
  } catch (error) {
    console.error('[ERROR] getLatestReadDocListNew - 错误详情:', JSON.stringify(error));
    throw error;
  }
};

export const normalizeDocDetailNew = (item: UserDocInfo) => {
  normalizeAssign(
    item,
    {
      docId: '',
      docName: '',
      paperId: '',
      pdfId: '',
      remark: '',
      isLatestRead: false,
      newPaper: false,
      classifyInfos: [],
      displayAuthor: {
        authorInfos: item.displayAuthor?.authorInfos || [],
        originAuthorInfos: item.displayAuthor?.originAuthorInfos || [],
        userEdited: false,
      },
      displayVenue: {
        venueInfos: [],
        originVenueInfos: [],
        userEdited: false,
      },
      displayPublishDate: {
        publishDate: '',
        originPublishDate: '',
        userEdited: false,
      },
      searchResult: {
        hitDocName: '',
        hitNote: '',
        hitRemark: '',
        hitPublishDate: '',
        hitAuthor: '',
        hitVenue: '',
        hitJcrVenuePartion: '',
      },
      jcrVenuePartion: {
        jcrVenuePartion: '',
        originJcrVenuePartion: '',
        userEdited: false,
      },
      impactOfFactor: {
        impactOfFactor: null as unknown as undefined,
        originImpactOfFactor: null as unknown as undefined,
        userEdited: false,
      },
      parsedStatus: 0,
    },
    true
  );
};

export const getClassifyOptionList = async (
  params: GetDocRelatedClassifyListReq
) => {
  const response = await api.post<
    SuccessResponse<GetDocRelatedClassifyListResponse>
  >(
    `/userDoc/getDocRelatedClassifyList`,
    params
  );

  return response.data.data.classifyInfos;
};

export const getAuthorOptionList = async (
  params: GetDocRelatedAuthorListReq
) => {
  const response = await api.post<
    SuccessResponse<GetDocRelatedAuthorListResponse>
  >(`/userDoc/getDocRelatedAuthorList`, params);

  return response.data.data.authorInfos;
};

export const getVenueOptionList = async (params: GetDocRelatedVenueListReq) => {
  const response = await api.post<
    SuccessResponse<GetDocRelatedVenueListResponse>
  >(`/userDoc/getDocRelatedVenueList`, params);

  return response.data.data.venueInfos;
};

export const getJcrOptionList = async (params: JcrPartionsReq) => {
  const response = await api.get<SuccessResponse<JcrPartionsResponse>>(
    `/userDoc/jcr/partions`,
    { params }
  );

  return response.data.data.jcrPartions;
};

export const copyDocOrFolderToAnotherFolder = async (
  params: CopyDocOrFolderToAnotherFolderReq
) => {
  try {
    await api.post<SuccessResponse<null>>(
      `/client/doc/copyDocOrFolderToAnotherFolder`,
      params
    );
    return true;
  } catch (error) {
    return false;
  }
};

export const moveDocOrFolderToAnotherFolder = async (
  params: MoveDocOrFolderToAnotherFolderReq
) => {
  try {
    await api.post<SuccessResponse<null>>(
      `/client/doc/moveDocOrFolderToAnotherFolder`,
      params
    );
    return true;
  } catch (error) {
    return false;
  }
};

export const addFolder = async (params: AddFolderReq) => {
  try {
    return await api.post<SuccessResponse<null>>(
      `/client/doc/addFolder`,
      params
    );
  } catch (error) {
    return null;
  }
};

export const moveFolderOrDoc = async (params: MoveFolderOrDocReq) => {
  try {
    await api.post<SuccessResponse<null>>(
      `/client/doc/moveFolderOrDoc`,
      params
    );
    return true;
  } catch (error) {
    return false;
  }
};

export const updateDoc = async (params: RenameUserDocReq) => {
  try {
    await api.post<SuccessResponse<null>>(
      `/userDoc/renameUserDoc`,
      params
    );
    return true;
  } catch (error) {
    return false;
  }
};

export const updateFolder = async (params: UpdateFolderReq) => {
  try {
    await api.post<SuccessResponse<null>>(
      `/client/doc/updateFolder`,
      params
    );
    return true;
  } catch (error) {
    return false;
  }
};

export const getFolderTreeAndSimpleDocInfoList = async (
  params: GetFolderTreeAndSimpleDocInfoListReq
) => {
  const response = await api.post<GetFolderTreeAndSimpleDocInfoListResponse>(
    `/client/doc/getFolderTreeAndSimpleDocInfoList`,
    params
  );
  const data =
    response.data.data ??
    ({} as NonNullable<GetFolderTreeAndSimpleDocInfoListResponse['data']>);

  normalizeAssign(
    data,
    {
      totalDocCount: 0,
      folderInfos: [],
      unclassifiedDocInfos: [],
    },
    false
  );

  return data;
};

export const getDocIndex = async (params: GetDocIndexReq) => {
  const response = await api.post<SuccessResponse<GetDocIndexResponse>>(
    `/userDoc/getDocIndex`,
    params
  );

  const { data } = response.data;
  normalizeAssign(
    data,
    {
      totalDocCount: 0,
      folderInfos: [],
      unclassifiedDocInfos: [],
    },
    true
  );

  return data;
};

export const deleteFolder = async (params: DeleteFolderReq) => {
  try {
    await api.post<SuccessResponse<null>>(
      `/client/doc/deleteFolder`,
      params
    );
    return true;
  } catch (error) {
    return false;
  }
};

export const deleteDoc = async (params: DeleteDocReq) => {
  try {
    await api.post<SuccessResponse<null>>(
      `/client/doc/deleteDoc`,
      params
    );
    return true;
  } catch (error) {
    return false;
  }
};

// export const updateVenue = async (params: UpdateVenueRequest) => {
//   return await api.put<SuccessResponse<UpdateVenueResponse>>(
//     `/userDoc/venue`,
//     params
//   );
// };

// export const updateJcr = async (params: JcrPartionUpdateReq) => {
//   return await api.put<SuccessResponse<JcrPartionUpdateResponse>>(
//     `/userDoc/jcr/partion`,
//     params
//   );
// };

// export const updateImpactFactor = async (params: ImpactFactorReq) => {
//   return await api.put<SuccessResponse<ImpactFactorResponse>>(
//     `/userDoc/impact/factor`,
//     params
//   );
// };

// export const updateScore = async (params: ImportanceScoreReq) => {
//   return await api.post<SuccessResponse<boolean>>(
//     `/userDoc/importance/score`,
//     params
//   );
// };

// export const updatePublishDate = async (params: UpdatePublishDateRequest) => {
//   return await api.put<SuccessResponse<UpdatePublishDateResponse>>(
//     `/userDoc/publishDate`,
//     params
//   );
// };

// export const updateAuthors = async (params: Partial<UpdateAuthorsRequest>) => {
//   const response = await api.put<SuccessResponse<UpdateAuthorsResponse>>(
//     `/userDoc/authors`,
//     params
//   );

//   return response.data.data.newAuthor;
// };

// export const getAuthors = async (params: GetAuthorsRequest) => {
//   const response = await api.get<SuccessResponse<GetAuthorsResponse>>(
//     `/userDoc/authors`,
//     { params }
//   );

//   const { data } = response.data;
//   if (!data.displayAuthor) {
//     data.displayAuthor = {
//       authors: [],
//       rollbackEnable: false,
//       originAuthors: [],
//     };
//   }
//   data.displayAuthor.authors.forEach((item) => {
//     if (typeof item.isAuthentication === 'undefined') {
//       item.isAuthentication = false;
//     }
//     if (typeof item.id === 'undefined') {
//       item.id = '';
//     }
//   });

//   return data;
// };

export const getFolders = async (params: GetUserDocFoldersRequest) => {
  const response = await api.get<SuccessResponse<GetUserDocFoldersResponse>>(
    `/userDoc/folders`,
    { params }
  );

  return response.data.data.folders;
};

export const getUserAllClassifyList = async (
  params: GetUserAllClassifyListRequest
) => {
  try {
    const res = await api.post<any>(
      `/userDoc/getUserAllClassifyList`,
      params
    );
    // 处理新的数据结构 {data: {results: []}, message: "获取用户所有文档分类列表成功", status: 1}
    return res.data.data.results || [];
  } catch (error) {
    return null;
  }
};

export const exportBibTexByIds = async (params: ExportByIdsReq) => {
  const res = await api.post(
    `/doc/bibtex/exportByIds`,
    params,
    {
      responseType: 'blob',
    }
  );

  return res.data;
};

export const exportBibTexByFolderId = async (params: ExportByFolderIdReq) => {
  const res = await api.post(
    `/doc/bibtex/exportByFolderId`,
    params,
    {
      responseType: 'blob',
    }
  );

  return res.data;
};

// 新增分类
export const addClassify = async (params: AddClassifyRequest) => {
  const res = await api.post<SuccessResponse<AddClassifyResponse>>(
    `/userDoc/addClassify`,
    params
  );
  return res.data.data;
};

// 删除分类
export const deleteClassify = async (params: DeleteClassifyRequest) => {
  const res = await api.post<SuccessResponse<null>>(
    `/userDoc/deleteClassify`,
    params
  );
  return res.data.data;
};

export const updateDocRemark = async (params: UpdateDocRemarkRequest) => {
  const res = await api.post<GetMyCollectedDocFolderResponse>(
    `/userDoc/updateDocRemark`,
    params
  );
  return res.data.data || [];
};

export const attachDocToClassify = async (
  params: AttachDocToClassifyRequest
) => {
  try {
    await api.post<AttachDocToClassifyResponse>(
      `/userDoc/attachDocToClassify`,
      params
    );
    return true;
  } catch (error) {
    return false;
  }
};

export const removeDocFromClassify = async (
  params: RemoveDocFromClassifyRequest
) => {
  const res = await api.post(
    `/userDoc/removeDocFromClassify`,
    params
  );
  return res.data.data;
};

export const updateDocReadStatus = async (params: UpdateReadStatusRequest) => {
  const res = await api.post<object>(
    `/doc/userDoc/updateReadStatus`,
    params
  );

  return res.data;
};

// 获取用户的一级文件夹列表
export const getTopLevelFolderList = async (
  param: GetTopLevelFolderListReq
) => {
  const res = await api.post<GetTopLevelFolderListResponse>(
    `/client/doc/getTopLevelFolderList`,
    param
  );
  return res.data || {};
};

// 将论文添加到文件夹
export const attachPaperToFolder = async (param: AttachPaperToFolderReq) => {
  const res = await api.post<SuccessResponse<AttachPaperToFolderResponse>>(
    `/client/doc/attachPaperToFolder`,
    param
  );
  return res.data.data;
};

/**
 * 获取文档创建状态
 * @param md5 文件的MD5值
 * @returns 文档创建状态信息
 */
export const getUserDocCreateStatus = async (param: GetUserDocCreateStatusRequest) => {
  try {
    const response = await api.post<SuccessResponse<GetUserDocCreateStatusResponse>>(
      `/userDoc/GetUserDocCreateStatus`,
      param
    );
    return response.data;
  } catch (error) {
    console.error('获取文档创建状态失败:', error);
    throw error;
  }
};

/**
 * 获取上传令牌
 * @param param 上传令牌请求参数
 * @returns 上传令牌信息
 */
export const getUploadToken = async (param: GetUploadTokenReq) => {
  try {
    const response = await api.post<SuccessResponse<GetUploadTokenResp>>(
      `/userDoc/GetUploadToken`,
      param
    );
    return response.data;
  } catch (error) {
    console.error('获取上传令牌失败:', error);
    throw error;
  }
};

/**
 * 获取解析令牌
 * @returns 解析令牌信息
 */
export const getParseToken = async () => {
  try {
    const response = await api.post<SuccessResponse<GetParseTokenResponse>>(
      `/userDoc/GetParseToken`
    );
    return response.data;
  } catch (error) {
    console.error('获取解析令牌失败:', error);
    throw error;
  }
};

/**
 * 获取上传令牌
 * @param param 上传令牌请求参数
 * @returns 上传令牌信息
 */
export const getS3UploadToken = async (param: GetMinioUploadTokenRequest) => {
  try {
    const response = await api.post<SuccessResponse<GetMinioUploadTokenResponse>>(
      `/oss/s3/getS3UploadToken`,
      param
    );
    return response.data;
  } catch (error) {
    console.error('获取上传令牌失败:', error);
    throw error;
  }
};

/**
 * 处理文件快速上传（秒传）
 * @param param 快速上传请求参数
 * @returns 快速上传响应
 */
export const handleFileFastUpload = async (param: HandleFileFastUploadReq) => {
  try {
    const response = await api.post<SuccessResponse<HandleFileFastUploadResp>>(
      `/userDoc/HandleFileFastUpload`,
      param
    );
    return response.data;
  } catch (error) {
    console.error('处理文件快速上传失败:', error);
    throw error;
  }
};


/**
 * 通过URL链接上传PDF文件
 * @param param 快速上传请求参数
 * @returns 快速上传响应
 */
export const uploadPdfByUrlLink = async (param: UploadPdfByUrlLinkRequest) => {
  try {
    const response = await api.post<SuccessResponse<null>>(
      `/userDoc/uploadPdfByUrlLink`,
      param
    );
    return response.data;
  } catch (error) {
    console.error('通过URL链接上传PDF文件失败:', error);
    throw error;
  }
};