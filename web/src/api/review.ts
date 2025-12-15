import { api } from '@common/api/axios';
import { SuccessResponse } from '@common/api/type';
import {
  HEADER_CANCLE_AUTO_ERROR,
  REQUEST_SERVICE_NAME_AI_REVIEW,
} from '@common/api/const';
import {
  GetConfigInfoResponse,
  // UploadReq,
  UploadResponse,
  SaveTaskReq,
  SaveTaskResponse,
  CancelTaskReq,
  GetTaskInfoReq,
  GetTaskListResponse,
  TaskInfo,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/review/AiReviewPaper';

export const getQSReviewConfig = async () => {
  const res = await api.get<SuccessResponse<GetConfigInfoResponse>>(
    `${REQUEST_SERVICE_NAME_AI_REVIEW}/graduationPaperReview/getConfigInfo`
  );
  return res.data.data;
};

export const uploadQSReviewFile = async (params: FormData) => {
  const res = await api.post<SuccessResponse<UploadResponse>>(
    `${REQUEST_SERVICE_NAME_AI_REVIEW}/graduationPaperReview/upload`,
    params
  );

  return res.data.data;
};

export const saveQSReviewTask = async (params: SaveTaskReq) => {
  const res = await api.post<SuccessResponse<SaveTaskResponse>>(
    `${REQUEST_SERVICE_NAME_AI_REVIEW}/graduationPaperReview/saveTask`,
    params,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );

  return res.data.data;
};

export const cancelQSReviewTask = async (params: CancelTaskReq) => {
  const res = await api.post<SuccessResponse<object>>(
    `${REQUEST_SERVICE_NAME_AI_REVIEW}/graduationPaperReview/cancelTask`,
    params,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );

  return res.data.data;
};

export const getQSReviewTask = async (params: GetTaskInfoReq) => {
  const res = await api.get<SuccessResponse<TaskInfo>>(
    `${REQUEST_SERVICE_NAME_AI_REVIEW}/graduationPaperReview/getTaskInfo`,
    {
      params,
    }
  );

  return res.data.data;
};

export const getQSReviewTaskList = async () => {
  const res = await api.get<SuccessResponse<GetTaskListResponse>>(
    `${REQUEST_SERVICE_NAME_AI_REVIEW}/graduationPaperReview/getTaskList`
  );

  return res.data.data;
};
