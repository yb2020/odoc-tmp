import {
  GetPicMd5Response,
  SimilarPicSearchReq,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/search/PicSearch';
import api from './axios';
import { REQUEST_SERVICE_NAME_APP } from './const';
import { SuccessResponse } from './type';

export const getImgMd5FromUrl = async (params: SimilarPicSearchReq) => {
  const { data: res } = await api.post<SuccessResponse<GetPicMd5Response>>(
    `${REQUEST_SERVICE_NAME_APP}/search/similarPicSearch`,
    params
  );

  return res.data;
};
