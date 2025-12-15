import {
  GetMyCslListResponse,
  GetCslListReq,
  GetCslListResponse,
  AddCslReq,
  DeleteCslReq,
  SortCslReq,
  CetDefaultCslListResponse,
  GetDocMetaInfoReq,
  DocMetaInfoSimpleVo,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/CSL';
import { SuccessResponse } from './type';
import { REQUEST_SERVICE_READPAPER_DOC } from './const';
import { api } from './axios';

// 获取我的格式列表
export const getMyCslList = async () => {
  const res = await api.post<SuccessResponse<GetMyCslListResponse>>(
    `${REQUEST_SERVICE_READPAPER_DOC}/csl/myCslList`
  );

  return res.data.data;
};

// 获取全部格式列表带搜索
export const getCslList = async (params: GetCslListReq) => {
  const res = await api.post<SuccessResponse<GetCslListResponse>>(
    `${REQUEST_SERVICE_READPAPER_DOC}/csl/list`,
    params
  );

  return res.data.data;
};

// 添加格式
export const addCsl = async (params: AddCslReq) => {
  const res = await api.post(
    `${REQUEST_SERVICE_READPAPER_DOC}/csl/addCsl`,
    params
  );

  return res.data.data;
};

// 删除格式
export const deleteCsl = async (params: DeleteCslReq) => {
  const res = await api.post(
    `${REQUEST_SERVICE_READPAPER_DOC}/csl/deleteCsl`,
    params
  );
  return res.data.data;
};

// 排序
export const sortCsl = async (params: SortCslReq) => {
  const res = await api.post(
    `${REQUEST_SERVICE_READPAPER_DOC}/csl/sortCsl`,
    params
  );
  return res.data.data;
};

// 获取默认格式列表
export const getDefaultCslList = async () => {
  const res = await api.post<SuccessResponse<CetDefaultCslListResponse>>(
    `${REQUEST_SERVICE_READPAPER_DOC}/docPublic/csl/getDefaultCslList`
  );

  return res.data.data;
};

// 获取meta数据
export const getDocMetaInfo = async (params: GetDocMetaInfoReq) => {
  const res = await api.post<SuccessResponse<DocMetaInfoSimpleVo>>(
    `${REQUEST_SERVICE_READPAPER_DOC}/docPublic/docMetaInfo/getDocMetaInfo`,
    params
  );

  return res.data.data;
};
