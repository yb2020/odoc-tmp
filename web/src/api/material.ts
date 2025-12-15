/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-non-null-assertion */
import api from './axios';

import {
  GetCitationRequest,
  GetCitationResponse,
  GetReferenceRequest,
  GetReferenceResponse,
} from 'go-sea-proto/gen/ts/pdf/PdfParse';
import {
  UpdateAuthorsRequest,
  UpdateAuthorsResponse,
  UpdatePublishDateRequest,
  UpdatePublishDateResponse,
  UpdateVenueRequest,
  UpdateVenueResponse,
  GetUserDocFoldersRequest,
  GetUserDocFoldersResponse,
  GetAuthorsRequest,
  GetAuthorsResponse,
  UserDocStatusByIdsRequest,
  UserDocStatusByIdsResponse,
  // GetUserDocRequest,
} from 'go-sea-proto/gen/ts/doc/UserDoc';
import {GetUserDocRequest} from 'go-sea-proto/gen/ts/doc/UserDoc'
import {
  GetDocIndexReq,
  GetDocIndexResponse,
} from 'go-sea-proto/gen/ts/doc/UserDocManage';
import {
  AddFolderReq,
  AddFolderResponse,
  CopyDocOrFolderToAnotherFolderReq,
  // DocDetailInfo,
  RemoveDocFromFolderReq,
  RenameUserDocReq,
} from 'go-sea-proto/gen/ts/doc/ClientDoc';
import { DocDetailInfo } from 'go-sea-proto/gen/ts/doc/ClientDoc'
import {
  AddClassifyRequest,
  AddClassifyResponse,
  AttachDocToClassifyRequest,
  AttachDocToClassifyResponse,
  CancelCollectPaperRequest,
  CollectPaperRequest,
  GetMyCollectedDocFolderResponse,
  GetUserAllClassifyListRequest,
  GetUserAllClassifyListResponse,
  MyCollectedDocInfo,
  RemoveDocFromClassifyRequest,
  UpdateDocRemarkRequest,
} from 'go-sea-proto/gen/ts/userCenter/UserDoc';
import { normalizeAssign } from '@idea/aiknowledge-special-util/normalize-assign';
import { ResponseError, SuccessResponse } from './type';
import {
  HEADER_CANCLE_AUTO_ERROR,
  REQUEST_SERVICE_NAME_APP,
  REQUEST_SERVICE_NAME_DOC,
} from './const';
import { AxiosRequestConfig } from 'axios';
import merge from 'lodash-es/merge';

export function isGetReferenceResponse(
  res: GetReferenceResponse | GetCitationResponse
): res is GetReferenceResponse {
  return (res as GetReferenceResponse).referenceInfoList !== undefined;
}

export const getReference = async (
  params: GetReferenceRequest
): Promise<GetReferenceResponse> => {
  const res = await api.post<SuccessResponse<GetReferenceResponse>>(
    `/pdf/parser/getReference`,
    params,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: true,
      },
    }
  );
  return res.data.data;
};

export const getReferenceFinal = (() => {
  let count = 100;
  let error: null | ResponseError = null;
  return (params: GetReferenceRequest): Promise<GetReferenceResponse> => {
    return new Promise(async (resolve, reject) => {
      const fetch = async () => {
        try {
          const data = await getReference(params);
          if (!data.needFetch) {
            return resolve(data);
          }
        } catch (err) {
          error = err as ResponseError;
        }
        if (count-- <= 0) {
          return reject(error);
        }
        setTimeout(() => {
          fetch();
        }, 1000);
      };

      fetch();
    });
  };
})();

export const getCitation = async (
  params: GetCitationRequest,
  options: AxiosRequestConfig = {}
): Promise<GetCitationResponse> => {
  const res = await api.post<SuccessResponse<GetCitationResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/pdfApi/parser/getCitation`,
    params,
    merge(
      {
        headers: {
          [HEADER_CANCLE_AUTO_ERROR]: true,
        },
      },
      options
    )
  );

  return res.data.data;
};

export const getCitationFinal = () => {
  let count = 100;
  let error: null | ResponseError = null;
  let timer = 0;
  return {
    cancle: () => {
      if (timer) {
        window.clearTimeout(timer);
      }
    },
    request: (
      params: GetCitationRequest,
      options: AxiosRequestConfig = {}
    ): Promise<GetCitationResponse> => {
      return new Promise(async (resolve, reject) => {
        const fetch = async () => {
          try {
            const data = await getCitation(params, options);
            if (!data.needFetch) {
              return resolve(data);
            }
          } catch (err) {
            error = err as ResponseError;
          }
          if (count-- <= 0) {
            return reject(error);
          }
          if (timer) {
            window.clearTimeout(timer);
          }
          timer = window.setTimeout(() => {
            fetch();
          }, 2000);
        };

        fetch();
      });
    },
  };
};

export const updateDoc = async (params: RenameUserDocReq) => {
  return await api.post<SuccessResponse<null>>(
    `/userDoc/renameUserDoc`,
    params
  );
};

export const updateVenue = async (params: Partial<UpdateVenueRequest>) => {
  return await api.post<SuccessResponse<UpdateVenueResponse>>(
    `/userDoc/updateVenue`,
    params
  );
};

export const updatePublishDate = async (
  params: Partial<UpdatePublishDateRequest>
) => {
  const response = await api.post<SuccessResponse<UpdatePublishDateResponse>>(
    `/userDoc/updatePublishDate`,
    params
  );

  return response.data.data;
};

export const updateAuthors = async (params: Partial<UpdateAuthorsRequest>) => {
  const response = await api.post<SuccessResponse<UpdateAuthorsResponse>>(
    `/userDoc/updateAuthors`,
    params
  );

  return response.data.data.newAuthor;
};

export const getAuthors = async (params: GetAuthorsRequest) => {
  const response = await api.post<SuccessResponse<GetAuthorsResponse>>(
    `/userDoc/getAuthors`,
    params
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
    `/userDoc/folders`,
    { params }
  );

  return response.data.data.folders;
};

export const getUserDoc = async (params: GetUserDocRequest) => {
  const response = await api.get<SuccessResponse<DocDetailInfo>>(
    `/doc/userDoc`,
    { params }
  );
  const item = response.data.data || {};
  normalizeAssign(item, {
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
  });

  return item;
};

export const getDocsIndex = async (params: GetDocIndexReq) => {
  const res = await api.post<SuccessResponse<GetDocIndexResponse>>(
    `/userDoc/getDocIndex`,
    params
  );
  return res.data.data;
};

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

export const updateDocRemark = async (params: UpdateDocRemarkRequest) => {
  const res = await api.post<GetMyCollectedDocFolderResponse>(
    `/userDoc/updateDocRemark`,
    params
  );
  return res.data.data || [];
};

export const addClassify = async (params: AddClassifyRequest) => {
  const res = await api.post<SuccessResponse<AddClassifyResponse>>(
    `/userDoc/addClassify`,
    params
  );
  return res.data.data;
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

export const getUserAllClassifyList = async (
  params: GetUserAllClassifyListRequest
) => {
  try {
    const res = await api.post<GetUserAllClassifyListResponse>(
      `/userDoc/getUserAllClassifyList`,
      params
    );
    return res.data.data.results || [];
  } catch (error) {
    return null;
  }
};

export const addFolder = async (param: AddFolderReq) => {
  const res = await api.post<SuccessResponse<AddFolderResponse>>(
    `/client/doc/addFolder`,
    param
  );
  return res.data.data;
};

// 收藏论文
export const postCollectPaper = async (param: Partial<CollectPaperRequest>) => {
  const res = await api.post<SuccessResponse<MyCollectedDocInfo>>(
    `/userDoc/collectPaper`,
    param
  );
  return res.data.data || {};
};
// 取消收藏论文
export const postCancelCollectPaper = async (
  param: CancelCollectPaperRequest
) => {
  const res = await api.post<SuccessResponse<unknown>>(
    `/userDoc/cancelCollectPaper`,
    param
  );
  return res.data.data || {};
};


export const updateImpactFactor = async (params: ImpactFactorReq) => {
  return await api.post<SuccessResponse<ImpactFactorResponse>>(
    `/userDoc/update/impact/factor`,
    params
  );
};

export const updateJcr = async (params: JcrPartionUpdateReq) => {
  return await api.post<SuccessResponse<JcrPartionUpdateResponse>>(
    `/userDoc/update/jcr/partion`,
    params
  );
};

export const updateScore = async (params: ImportanceScoreReq) => {
  return await api.post<SuccessResponse<boolean>>(
    `/userDoc/importance/score`,
    params
  );
};

/**
 * 批量获取文档状态
 * @param params - 文档ID列表
 */
export const getUserDocStatusByIds = async (params: UserDocStatusByIdsRequest): Promise<UserDocStatusByIdsResponse> => {
  const response = await api.post<SuccessResponse<UserDocStatusByIdsResponse>>(
    `/doc/userDocStatusByIds`,
    params
  );
  return response.data.data;
};