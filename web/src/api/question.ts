/*
 * Created Date: August 9th 2021, 5:21:09 pm
 * Author: zhoupengcheng
 * -----
 * Last Modified: June 16th 2022, 10:48:00 am
 */
import api from './axios';
import {
  QueryPaperQuestionRequest,
  QueryPaperQuestionResponse,
  AddPaperQuestionRequest,
  AddPaperQuestionResponse,
  QueryPaperAnswerRequest,
  QueryPaperAnswerResponse,
  AddPaperAnswerRequest,
  AddPaperAnswerResponse,
  DeletePaperAnswerRequest,
  DeletePaperAnswerResponse,
  DeletePaperQuestionRequest,
  DeletePaperQuestionResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/paper/QuestionAnswer';
import { RequestParam, SuccessResponse } from './type';
import { REQUEST_SERVICE_NAME_APP } from './const';

export const getQuestionList = async (params: QueryPaperQuestionRequest) => {
  const { data: res } = await api.post<
    SuccessResponse<QueryPaperQuestionResponse>
  >(
    `${REQUEST_SERVICE_NAME_APP}/aiKnowledge/paper/qa/queryPaperQuestion`,
    params
  );

  return res.data;
};

export const addQuestion = async (params: AddPaperQuestionRequest) => {
  const res = await api.post<SuccessResponse<AddPaperQuestionResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/aiKnowledge/paper/qa/addQuestion`,
    params
  );

  return res;
};

export const deleteQuestion = async (params: DeletePaperQuestionRequest) => {
  const { data: res } = await api.post<
    SuccessResponse<DeletePaperQuestionResponse>
  >(`${REQUEST_SERVICE_NAME_APP}/aiKnowledge/paper/qa/deleteQuestion`, params);

  return res;
};

export const getQuestionDetail = async (params: QueryPaperAnswerRequest) => {
  const { data: res } = await api.post<
    SuccessResponse<QueryPaperAnswerResponse>
  >(
    `${REQUEST_SERVICE_NAME_APP}/aiKnowledge/paper/qa/queryPaperAnswer`,
    params
  );

  return res.data;
};

export const addAnswer = async (params: Partial<AddPaperAnswerRequest>) => {
  const { data: res } = await api.post<SuccessResponse<AddPaperAnswerResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/aiKnowledge/paper/qa/addAnswer`,
    params
  );

  return res.data;
};

export const deleteAnswer = async (params: DeletePaperAnswerRequest) => {
  const { data: res } = await api.post<
    SuccessResponse<DeletePaperAnswerResponse>
  >(`${REQUEST_SERVICE_NAME_APP}/aiKnowledge/paper/qa/deleteAnswer`, params);

  return res.data;
};

export const getAllowedDeleteQuestionAndAnswerUserIdList = async () => {
  const { data: res } = await api.post(
    `${REQUEST_SERVICE_NAME_APP}/aiKnowledge/paper/qa/getAllowedDeleteQuestionAndAnswerUserIdList`
  );

  return res.data;
};
