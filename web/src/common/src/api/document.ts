import { normalizeAssign } from '@idea/aiknowledge-special-util';
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
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/client/ClientDoc';
import {
  GetAuthorsRequest,
  GetAuthorsResponse,
  GetUserDocFoldersRequest,
  GetUserDocFoldersResponse,
  UpdateAuthorsRequest,
  UpdateAuthorsResponse,
  UpdatePublishDateRequest,
  UpdatePublishDateResponse,
  UpdateVenueRequest,
  UpdateVenueResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/doc/UserDoc';
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
  // UpdateReadStatusRequest,
  UserDocInfo,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/UserDocManage';
import {
  REQUEST_SERVICE_NAME_APP,
  // REQUEST_SERVICE_READPAPER_DOC,
} from './const';
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
import {
  ExportByFolderIdReq,
  ExportByIdsReq,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/bibtex/BibtexInfo';

export const removeDocFromFolder = async (params: RemoveDocFromFolderReq) => {
  try {
    await api.post<SuccessResponse<null>>(
      `${REQUEST_SERVICE_NAME_APP}/client/doc/removeDocFromFolder`,
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
  >(`${REQUEST_SERVICE_NAME_APP}/client/doc/getDocListByFolderId`, params);
  const list = res.data.data || [];
  list.forEach(normalizeDocDetailInfo);
  return list;
};

export const getDocListNew = async (params: GetDocListReq) => {
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
      `${REQUEST_SERVICE_NAME_APP}/client/doc/copyDocOrFolderToAnotherFolder`,
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
      `${REQUEST_SERVICE_NAME_APP}/client/doc/moveDocOrFolderToAnotherFolder`,
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
      `${REQUEST_SERVICE_NAME_APP}/client/doc/addFolder`,
      params
    );
  } catch (error) {
    return null;
  }
};

export const moveFolderOrDoc = async (params: MoveFolderOrDocReq) => {
  try {
    await api.post<SuccessResponse<null>>(
      `${REQUEST_SERVICE_NAME_APP}/client/doc/moveFolderOrDoc`,
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
      `${REQUEST_SERVICE_NAME_APP}/userDoc/renameUserDoc`,
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
      `${REQUEST_SERVICE_NAME_APP}/client/doc/updateFolder`,
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
    `${REQUEST_SERVICE_NAME_APP}/client/doc/getFolderTreeAndSimpleDocInfoList`,
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
      `${REQUEST_SERVICE_NAME_APP}/client/doc/deleteFolder`,
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
      `${REQUEST_SERVICE_NAME_APP}/client/doc/deleteDoc`,
      params
    );
    return true;
  } catch (error) {
    return false;
  }
};

export const updateVenue = async (params: UpdateVenueRequest) => {
  return await api.put<SuccessResponse<UpdateVenueResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/userDoc/venue`,
    params
  );
};

export const updateJcr = async (params: JcrPartionUpdateReq) => {
  return await api.put<SuccessResponse<JcrPartionUpdateResponse>>(
    `/userDoc/jcr/partion`,
    params
  );
};

export const updateImpactFactor = async (params: ImpactFactorReq) => {
  return await api.put<SuccessResponse<ImpactFactorResponse>>(
    `/userDoc/impact/factor`,
    params
  );
};

export const updateScore = async (params: ImportanceScoreReq) => {
  return await api.post<SuccessResponse<boolean>>(
    `/userDoc/importance/score`,
    params
  );
};

export const updatePublishDate = async (params: UpdatePublishDateRequest) => {
  return await api.put<SuccessResponse<UpdatePublishDateResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/userDoc/publishDate`,
    params
  );
};

export const updateAuthors = async (params: Partial<UpdateAuthorsRequest>) => {
  const response = await api.put<SuccessResponse<UpdateAuthorsResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/userDoc/authors`,
    params
  );

  return response.data.data.newAuthor;
};

export const getAuthors = async (params: GetAuthorsRequest) => {
  const response = await api.get<SuccessResponse<GetAuthorsResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/userDoc/authors`,
    { params }
  );

  const { data } = response.data;
  if (!data.displayAuthor) {
    data.displayAuthor = {
      authors: [],
      rollbackEnable: false,
      originAuthors: [],
    };
  }
  data.displayAuthor.authors.forEach((item) => {
    if (typeof item.isAuthentication === 'undefined') {
      item.isAuthentication = false;
    }
    if (typeof item.id === 'undefined') {
      item.id = '';
    }
  });

  return data;
};

export const getFolders = async (params: GetUserDocFoldersRequest) => {
  const response = await api.get<SuccessResponse<GetUserDocFoldersResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/userDoc/folders`,
    { params }
  );

  return response.data.data.folders;
};

export const getUserAllClassifyList = async (
  params: GetUserAllClassifyListRequest
) => {
  try {
    const res = await api.post<GetUserAllClassifyListResponse>(
      `${REQUEST_SERVICE_NAME_APP}/userDoc/getUserAllClassifyList`,
      params
    );
    return res.data.data;
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
    `${REQUEST_SERVICE_NAME_APP}/userDoc/addClassify`,
    params
  );
  return res.data.data;
};

// 删除分类
export const deleteClassify = async (params: DeleteClassifyRequest) => {
  const res = await api.post<SuccessResponse<null>>(
    `${REQUEST_SERVICE_NAME_APP}/userDoc/deleteClassify`,
    params
  );
  return res.data.data;
};

export const updateDocRemark = async (params: UpdateDocRemarkRequest) => {
  const res = await api.post<GetMyCollectedDocFolderResponse>(
    `${REQUEST_SERVICE_NAME_APP}/userDoc/updateDocRemark`,
    params
  );
  return res.data.data || [];
};

export const attachDocToClassify = async (
  params: AttachDocToClassifyRequest
) => {
  try {
    await api.post<AttachDocToClassifyResponse>(
      `${REQUEST_SERVICE_NAME_APP}/userDoc/attachDocToClassify`,
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
    `${REQUEST_SERVICE_NAME_APP}/userDoc/removeDocFromClassify`,
    params
  );
  return res.data.data;
};

// export const updateDocReadStatus = async (params: UpdateReadStatusRequest) => {
//   const res = await api.post<object>(
//     `/userDoc/updateReadStatus`,
//     params
//   );

//   return res.data;
// };

// 获取用户的一级文件夹列表
export const getTopLevelFolderList = async (
  param: GetTopLevelFolderListReq
) => {
  const res = await api.post<GetTopLevelFolderListResponse>(
    `${REQUEST_SERVICE_NAME_APP}/client/doc/getTopLevelFolderList`,
    param
  );
  return res.data || {};
};

// 将论文添加到文件夹
export const attachPaperToFolder = async (param: AttachPaperToFolderReq) => {
  const res = await api.post<SuccessResponse<AttachPaperToFolderResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/client/doc/attachPaperToFolder`,
    param
  );
  return res.data.data;
};
