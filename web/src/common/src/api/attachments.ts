import {
  GetAttachmentListReq,
  GetAttachmentListResponse,
  GetUploadTokenReq,
  GetUploadTokenResponse,
  SaveAttachmentReq,
  DeleteAttachmentReq,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/DocAttachment';
import {
  AcquirePolicyCallbackInfoResponse,
  ObjectStoreInfo,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/oss/AliOSS';
import axios, { AxiosResponse } from 'axios';
import { SuccessResponse } from './type';
import { api } from './axios';

const REQUEST_SERVICE_NAME_DOC = '/microservice-readpaper-doc';

/**
 * 已弃用：/userDoc/attachment/getAttachmentList接口已不再使用
 */
export const getDocAttachments = async (params: GetAttachmentListReq) => {
  // 注释掉原始实现
  // const res = await api.post<SuccessResponse<GetAttachmentListResponse>>(
  //   REQUEST_SERVICE_NAME_DOC + '/userDoc/attachment/getAttachmentList',
  //   params
  // );
  // return res.data.data;
  
  console.log('attachments.ts: getDocAttachments called, but API is deprecated');
  // 返回模拟数据
  return {
    attachmentList: [],
    totalCount: 0
  } as GetAttachmentListResponse;
};

export const getUploadToken = async (params: GetUploadTokenReq) => {
  const res = await api.post<GetUploadTokenResponse>(
    REQUEST_SERVICE_NAME_DOC + '/userDoc/attachment/getUploadToken',
    params
  );

  return res.data.data;
};

export const uploadAttachment = async (
  file: File,
  params: AcquirePolicyCallbackInfoResponse
): Promise<ObjectStoreInfo & { mimeType?: string }> => {
  const { accessid, callback, dir, host, policy, signature, ...rest } = params;

  if (!rest.isNeedUpload && rest.objectStoreInfo) {
    return rest.objectStoreInfo;
  }

  const formData = new FormData();

  formData.append('name', file.name);
  formData.append('key', `${dir}${file.name}`);
  if (params.isNeedAttachmentHeader) {
    formData.append('Content-Disposition', 'attachment');
  }
  formData.append('policy', policy);
  formData.append('OSSAccessKeyId', accessid);
  formData.append('success_action_status', '200'); // 让服务端返回200,不然，默认会返回204
  formData.append('callback', callback);
  formData.append('signature', signature);
  formData.append('file', file);

  const res: AxiosResponse<{
    bucketName: string;
    fileName: string;
    isSuccess: boolean;
    size: string;
    mimeType: string;
  }> = await axios(host, {
    method: 'post',
    data: formData,
  });

  if (res.status === 200 && res.data?.isSuccess) {
    return {
      bucketName: res.data.bucketName,
      objectName: res.data.fileName,
      mimeType: res.data.mimeType,
    };
  }

  throw new Error(
    `Upload failed: ${res.status === 200 ? res.data?.isSuccess : res.status}`
  );
};

export const saveAttachment = async (params: SaveAttachmentReq) => {
  const res = await api.post<SuccessResponse<unknown>>(
    REQUEST_SERVICE_NAME_DOC + '/userDoc/attachment/save',
    params
  );

  return res.data.data;
};

export const removeAttachment = async (params: DeleteAttachmentReq) => {
  const res = await api.post<SuccessResponse<unknown>>(
    REQUEST_SERVICE_NAME_DOC + '/userDoc/attachment/delete',
    params
  );

  return res.data.data;
};
