import {
  AiAnswerCancelDisLikeReq,
  AiAnswerCancelLikeReq,
  AiAnswerDisLikeFeedbackReq,
  AiAnswerDisLikeReq,
  AiAnswerLikeReq,
  AnswerInfo,
  AskQuestionReq,
  GetAiAnswerByIdReq,
  GetListReq,
  GetListResponse,
  GptAskQuestionReq,
  GptAskQuestionResponse,
  QuestionAnswerInfo,
  SelectTextQuestionReq,
  RetryByAnswerIdReq,
  GetAiAnswerByIdResp,
  ChangeAnswerRequest,
  ChangeAnswerResponse,
  ImageQuestionRequest,
  SatisfactionFeedbackRequest,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiAssistReading';
import { OptimizeFeedbackReq } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiAssistantReading/ReadingFeedbackInfo';
import {
  AddQuestionReq,
  AddQuestionResponse,
  GetListReq as GetCustomListReq,
  GetListResp as GetCustomListResp,
  DeleteQuestionReq,
  QuestionUseCountReq,
} from 'go-sea-proto/gen/ts/copilot/CustomQuestion';

import {
  GetQuestionListReq,
  GetQuestionListResp,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiAssistantReading/SystemQuestion';
import api from './axios';
import {
  REQUEST_SERVICE_NAME_AI_READING,
  REQUEST_SERVICE_NAME_APP,
  getHeadersWithCancelAutoError,
} from './const';
import { SuccessResponse } from './type';

import {SummarySinglePaperRequest,ChatMessage, GetChatMessagesRequest, StopChatMessagesRequest, ChatFeedbacksRequest} from 'go-sea-proto/gen/ts/copilot/chat/ChatFlow'
import ssePostFetcher from '@common/hooks/aitools/sse';
import { APIResponse } from 'go-sea-proto/gen/ts/common/Common';
import { ChatMessagesRequest } from 'go-sea-proto/gen/ts/copilot/chat/ChatFlow';
import { GetChatMessagesResponse } from 'go-sea-proto/gen/ts/copilot/chat/ChatFlow';

export const getList = async (params: GetListReq) => {
  const res = await api.post<SuccessResponse<GetListResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/aiAssistReadingQa/getList`,
    params
  );
  return res.data.data;
};

export const askQuestion = async (params: AskQuestionReq) => {
  const res = await api.post<SuccessResponse<QuestionAnswerInfo>>(
    `${REQUEST_SERVICE_NAME_APP}/aiAssistReadingQa/askQuestion`,
    params,
    {
      timeout: 120000,
      ...getHeadersWithCancelAutoError(),
    }
  );
  return res.data.data;
};

export const selectText = async (params: SelectTextQuestionReq) => {
  const res = await api.post<SuccessResponse<QuestionAnswerInfo>>(
    `${REQUEST_SERVICE_NAME_APP}/aiAssistReadingQa/selectTextQuestion`,
    params,
    {
      timeout: 120000,
      ...getHeadersWithCancelAutoError(),
    }
  );
  return res.data.data;
};

export const getAiAnswerById = async (params: GetAiAnswerByIdReq) => {
  const res = await api.post<SuccessResponse<Required<GetAiAnswerByIdResp>>>(
    `${REQUEST_SERVICE_NAME_APP}/aiAssistReadingQa/getAiAnswerById/v2`,
    params
  );
  return res.data.data;
};

export const retryAiAnswerById = async (params: RetryByAnswerIdReq) => {
  const res = await api.post<SuccessResponse<Required<AnswerInfo>>>(
    `${REQUEST_SERVICE_NAME_APP}/aiAssistReadingQa/retryByAnswerId`,
    params
  );
  return res.data.data;
};


export const likeAnswer = async (params: AiAnswerLikeReq) => {
  const res = await api.post<SuccessResponse<null>>(
    `${REQUEST_SERVICE_NAME_APP}/aiAnswer/like`,
    params
  );
  return res.data.data;
};

export const cancelLikeAnswer = async (params: AiAnswerCancelLikeReq) => {
  const res = await api.post<SuccessResponse<null>>(
    `${REQUEST_SERVICE_NAME_APP}/aiAnswer/cancelLike`,
    params
  );
  return res.data.data;
};

export const dislikeAnswer = async (params: AiAnswerDisLikeReq) => {
  const res = await api.post<SuccessResponse<null>>(
    `${REQUEST_SERVICE_NAME_APP}/aiAnswer/dislike`,
    params
  );
  return res.data.data;
};

export const cancelDislikeAnswer = async (params: AiAnswerCancelDisLikeReq) => {
  const res = await api.post<SuccessResponse<null>>(
    `${REQUEST_SERVICE_NAME_APP}/aiAnswer/cancelDislike`,
    params
  );
  return res.data.data;
};

export const dislikeFeedback = async (params: AiAnswerDisLikeFeedbackReq) => {
  const res = await api.post<SuccessResponse<null>>(
    `${REQUEST_SERVICE_NAME_APP}/aiAnswer/dislikeFeedback`,
    params
  );
  return res.data.data;
};

// export const changeAnswer = async (params: ChangeAnswerRequest) => {
//   const res = await api.get<SuccessResponse<ChangeAnswerResponse>>(
//     `${REQUEST_SERVICE_NAME_APP}/aiAssistReadingQa/answer/change`,
//     { params }
//   );
//   return res.data.data;
// };
export const changeAnswerSSE = (
  params: ChangeAnswerRequest,
  onMessage: (message: ChatMessage) => void
) => {
  const abortController = new AbortController();

  const connection = ssePostFetcher<ChangeAnswerRequest, ChatMessage>(
    `/api/copilot/answer/change`,
    params,
    (data) => {
      const parsedData = typeof data === 'string' ? JSON.parse(data) : data;
      onMessage(parsedData.data);
    },
    abortController
  );

  return {
    connection,
    abort: () => abortController.abort(),
  };
};

export const getCustomQuestions = async (params: GetCustomListReq) => {
  const res = await api.post<SuccessResponse<Required<GetCustomListResp>>>(
    `/customQuestion/getList`,
    params
  );
  console.log(res.data.data);
  return res.data.data?.questionList || [];
};

export const addCustomQuestion = async (params: AddQuestionReq) => {
  const res = await api.post<SuccessResponse<AddQuestionResponse>>(
    `/customQuestion/add`,
    params
  );
  return {
    questionId: res.data.data.questionId,
    question: params.question,
  };
};

export const deleteCustomQuestion = async (params: DeleteQuestionReq) => {
  await api.post<SuccessResponse<null>>(
    `/customQuestion/delete`,
    params
  );
};

export const countCustomQuestion = async (params: QuestionUseCountReq) => {
  await api.post<SuccessResponse<null>>(
    `/customQuestion/useCount`,
    params
  );
};

export const getSystemQuestionList = async (params: GetQuestionListReq) => {
  const res = await api.get<SuccessResponse<GetQuestionListResp>>(
    `${REQUEST_SERVICE_NAME_AI_READING}/systemQuestion/list`,
    params
  );
  return res.data.data?.items || [];
};

export const askImageQuestion = async (params: ImageQuestionRequest) => {
  const res = await api.post<SuccessResponse<QuestionAnswerInfo>>(
    `${REQUEST_SERVICE_NAME_APP}/aiAssistReadingQa/imageQuestion`,
    params,
    {
      timeout: 120000,
      ...getHeadersWithCancelAutoError(),
    }
  );
  return res.data.data;
};

export const satisfactionFeedback = async (
  params: SatisfactionFeedbackRequest
) => {
  const res = await api.post<SuccessResponse<Required<null>>>(
    `${REQUEST_SERVICE_NAME_AI_READING}/aiAssistReadingQa/satisfaction/feedback`,
    params
  );
  return res.data.data;
};

export enum OptimizeFeedbackType {
  Better = 'better',
  Worse = 'worse',
  Almost = 'almost',
}

export const optimizeFeedback = async (params: OptimizeFeedbackReq) => {
  const res = await api.post<SuccessResponse<Required<null>>>(
    `${REQUEST_SERVICE_NAME_AI_READING}/optimize/feedback`,
    params
  );
  return res.data.data;
};

export const getChatMessages = async (params: GetChatMessagesRequest) => {
  const res = await api.get<SuccessResponse<GetChatMessagesResponse>>(
    `/copilot/chat/messages`,
    { params }
  );
  return res.data.data;
};

export const stopChatMessages = async (params: StopChatMessagesRequest) => {
  const res = await api.post<SuccessResponse<Required<null>>>(
    `/copilot/chat/stop-messages`,
    params
  );
  return res.data.data;
};

export const feedbacks = async (params: ChatFeedbacksRequest) => {
  const res = await api.post<SuccessResponse<Required<null>>>(
    `/copilot/chat/feedbacks`,
    params
  );
  return res.data.data;
};

export const chatList = async (params: GetChatMessagesRequest) => {
  const res = await api.get<SuccessResponse<GetChatMessagesResponse>>(
    `/api/copilot/messages`,
    {params}
  );
  return res.data.data;
};

export const chatMessages = async (
  params: ChatMessagesRequest,
  onMessage?: (message: ChatMessage) => void
) => {
  const abortController = new AbortController();
  const messages: ChatMessage[] = [];
  const connection = ssePostFetcher<ChatMessagesRequest, ChatMessage>(
    `/api/copilot/chat-messages`,
    params,
    (data) => {
      const parsedData = typeof data === 'string' ? JSON.parse(data) : data;
      
      // 如果提供了回调函数，则实时调用它
      if (onMessage) {
        onMessage(parsedData.data);
        messages.push(parsedData.data);
      }
    },
    abortController
  );
  
  // 返回一个对象，包含连接和中止控制器
  return {
    connection,
    abort: () => abortController.abort(),
    messages: () => messages
  };
};

export const summarySinglePaper = (
  params: SummarySinglePaperRequest,
  onMessage?: (message: ChatMessage) => void
) => {
  const abortController = new AbortController();
  const messages: ChatMessage[] = [];
  
  // 直接返回 SSE 连接，不等待完成
  const connection = ssePostFetcher<SummarySinglePaperRequest, ChatMessage>(
    `/api/copilot/summary-single-paper`,
    params,
    (data) => {
      const parsedData = typeof data === 'string' ? JSON.parse(data) : data;
      
      // 如果提供了回调函数，则实时调用它
      if (onMessage) {
        const res = parsedData as APIResponse;
        if (res.status === 0) {
          console.error('Error in SSE response:', res.message);
          return;
        }
        onMessage(res.data as ChatMessage);
        messages.push(res.data as ChatMessage);
      }
    },
    abortController
  );
  
  // 返回一个对象，包含连接和中止控制器
  return {
    connection,
    abort: () => abortController.abort(),
    messages: () => messages
  };
};