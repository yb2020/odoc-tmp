import api from './axios';
import {
  GetMyLearningTaskInfoRequest,
  GetMyLearningTaskInfoResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/paper/PaperDetail';
import { SuccessResponse } from './type';
import { REQUEST_SERVICE_NAME_APP } from './const';
import {
  PaperCommentAggregationReq,
  PaperCommentAggregationRsp,
  PaperCommentUpdateReq,
  PaperCommentCreateReq,
  PaperCommentCreateRsp,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/paper/PaperCommentView';

export const getMyLearningTaskInfo = async (params: GetMyLearningTaskInfoRequest) => {
  const { data: res } = await api.post<SuccessResponse<GetMyLearningTaskInfoResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/aiKnowledge/paper/getMyLearningTaskInfo`,
    params
  )
  return res.data
}

export type $SavePaperTranslationRequest = {
  translation: string
  paperId: string
}

export const savePaperTranslation = async (
  params: $SavePaperTranslationRequest
) => {
  const { data: res } = await api.post(
    `${REQUEST_SERVICE_NAME_APP}/translate/savePaperTranslation`,
    params
  )

  return res.data
}

export type $UpdatePaperTranslationRequest = {
  id: string
  translation: string
  paperId: string
}

export const updatePaperTranslation = async (
  params: $UpdatePaperTranslationRequest
) => {
  const { data: res } = await api.post(
    `${REQUEST_SERVICE_NAME_APP}/translate/updatePaperTranslation`,
    params
  )
  return res.data
}


interface SaveAnswerParams {
  paperId: string
  id: string
  classicQuestionId: string
  answer: string
  htmlAnswer?: string
}

export const saveAnswer = async (
  params: SaveAnswerParams
) => {
  const { data: res } = await api.post(
    `${REQUEST_SERVICE_NAME_APP}/classicQuestionAnswer/saveAnswer`,
    params
  )

  return res.data
}

export const updateAnswer = async (
  params: SaveAnswerParams
) => {
  const { data: res } = await api.post(
    `${REQUEST_SERVICE_NAME_APP}/classicQuestionAnswer/updateAnswer`,
    params
  )

  return res.data
}

export const getCommentAggregation = async (
  params: PaperCommentAggregationReq
) => {
  const { data: res } = await api.get<
    SuccessResponse<PaperCommentAggregationRsp>
  >(`${REQUEST_SERVICE_NAME_APP}/aiKnowledge/paper/comment/aggregation`, {
    params,
  })

  return res.data
}

export const updatePaperComment = async (
  params: PaperCommentUpdateReq
) => {
  const { data: res } = await api.put(
    `${REQUEST_SERVICE_NAME_APP}/paper/comment`,
    params
  )

  return res.data
}

export const createPaperComment = async (
  params: PaperCommentCreateReq
) => {
  const { data: res } = await api.post<
    SuccessResponse<PaperCommentCreateRsp>
  >(`${REQUEST_SERVICE_NAME_APP}/paper/comment`, params)

  return res.data
}