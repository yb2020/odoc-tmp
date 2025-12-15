import { computed, getCurrentInstance, nextTick, reactive, ref } from 'vue';
import useFetch from './useFetch';
import PerfectScrollbarType from 'perfect-scrollbar';
import { PAGE_ROUTE_NAME } from '../routes/type';
import { emitter, COPILOT_ASK } from '../util/eventbus';
import pull from 'lodash-es/pull';

import {
  AnswerInfo,
  AskQuestionReq,
  BanStrategy,
  GetAiAnswerByIdReq,
  GetAiAnswerByIdResp,
  QuestionAnswerInfo,
  QuestionInfo,
  QuestionType,
  SelectTextQuestionReq,
  ImageQuestionRequest,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiAssistReading';
import { PageSelectText } from '~/src/pdf-annotate-viewer/type';
import {
  askQuestion,
  selectText,
  getChatMessages,
  getAiAnswerById,
  retryAiAnswerById,
  changeAnswerSSE,
  askImageQuestion,
  chatMessages,
  stopChatMessages,
  feedbacks,
} from '../api/copilot';
import { useStore } from '../store';
import { useUserStore } from '@common/stores/user'
import { polling } from '../util/polling';
import { message } from 'ant-design-vue';
import { SELF_NOTEINFO_GROUPID } from '../store/base/type';
import {
  ElementName,
  PageType,
  reportAiAssistReadingResponseTime,
  reportModuleImpression,
} from '~/src/api/report';
import { ImageStorageType, uploadImage } from '../api/upload';
import { ResponseError } from '../api/type';
import {
  BeanScenes,
  useAIBeans,
  useAIBeansBuy,
} from '@common/hooks/useAIBeans';
import {
  ERROR_CODE_BEANS_NOT_ENOUGH,
  ERROR_CODE_NEED_VIP,
} from '@common/api/const';
import { useVipStore } from '@common/stores/vip';
import {
  NeedAiBeanException,
  NeedVipException,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { useI18n } from 'vue-i18n';
import { ChatMessageObj, 
  ChatMessage, 
  ChatMessagesRequest, 
  SelectedText, 
  UploadFile, 
  AnswerStatus, 
  StopChatMessagesRequest,
  QuoteInfo,
} from 'go-sea-proto/gen/ts/copilot/chat/ChatFlow';
import { summarySinglePaper as apiSummarySinglePaper } from '~/src/api/copilot';
import { RectOptions } from 'go-sea-proto/gen/ts/common/RectOptions';
import QuestionItem from '../components/Right/TabPanel/QuestionAnswer/Answer/QuestionItem.vue';
import { useCopilotStore } from '../stores/copilotStore';
// import { abort } from 'process';

export interface QuestionItem {
  // 问题id，用于前端显示
  id: string;
  // 任务id，用于交互
  taskId: string;
  // noteid
  noteId: string;
  // pdfid
  pdfId: string;
  // 模型类型
  modelType: string;
  // 对话内容
  content: string;
  // 对话id 【复用dify的对话id】
  conversationId?: string;
  // 选区文本
  selectedText?: SelectedText;
  // 上传的文件组, 最多三个
  uploadFiles: UploadFile[];
  // 回答id,用于续写
  messageId?: string;
  // 回答列表，可能会有多个，比如：重试，换个答案什么的。
  answers?: ChatMessage[];
}

export interface ImageUpload {
  base64: string;
  pageNumber?: number;
}

const LIST_PAGESIZE = 10;

const DEFAULT_CHAMGE_ANSWER_LIMIT = 2;

const MAX_QUEUE = 1;
export interface CopilotState {
  list: QuestionItem[];
  hasmore: boolean;
  gptGrayTip: string;
  banInputMessage: string;
  limitTotal: number;
  currentPollingAnswers: string[];
  submitPending: boolean;
  activeTaskId: string; // 当前正在交互的任务id，如果为空则表示没有正在交互的任务或已经停止交互
  isActiveTask: boolean; // 当前task是否在运行
  abortController: (() => void) | null;
}

const userStore = useUserStore()

const genId = () => {
  return `fake-${Date.now()}`;
};

function isAskQuestionReq(s: unknown): s is AskQuestionReq {
  return !!s && !(s as SelectTextQuestionReq).boundingBox;
}

export const useCopilotState = () => {
  const copilotState = reactive<CopilotState>({
    list: [],
    hasmore: true,
    gptGrayTip: '',
    banInputMessage: '',
    limitTotal: DEFAULT_CHAMGE_ANSWER_LIMIT,
    currentPollingAnswers: [],
    submitPending: false,
    activeTaskId: '',
    isActiveTask: false,
    abortController: null,
  });
  return copilotState;
};

const useCopilot = (copilotState: ReturnType<typeof useCopilotState>) => {
  let queueLength = 0;
  const vipStore = useVipStore();
  const { consumeBeans, refundBeans } = useAIBeans();
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  const { showBuyDialog } = useAIBeansBuy(getCurrentInstance()?.appContext);

  const checkBeans = (err?: unknown) => {
    refundBeans(BeanScenes.COPILOT_ASK);
    const e = err as ResponseError;
    if (e?.code === ERROR_CODE_NEED_VIP) {
      vipStore.showVipLimitDialog(e?.message, {
        exception: e?.extra as NeedVipException,
        reportParams: {
          element_name: ElementName.upperAIAssistLimitPopup,
        },
      });
    } else if (e?.code === ERROR_CODE_BEANS_NOT_ENOUGH) {
      showBuyDialog();
      message.warn(e?.message);
      reportModuleImpression({
        page_type: PageType.note,
        type_parameter: 'none',
        element_name: (e.extra as NeedAiBeanException)?.needVipType
          ? ElementName.upper50PAIAssistToast
          : ElementName.upper100PAIAssistToast,
        element_parameter: 'none',
        module_type: '',
      });
    }
  };

  const ask = async (
    payload: {
      question: string;
      selectText?: PageSelectText[];
      quoteAnswer?: ChatMessage;
      firstAskQuestion?: boolean;
      imageBase64?: string;
      imagePageNumber?: number;
      isDeepThinking?: boolean;
    },
    streamPollingHandler?: (res: GetAiAnswerByIdResp) => void
  ) => {
    const params = formateQuestion(
      payload.question,
      payload.selectText,
      payload.quoteAnswer?.messageId,
      payload.firstAskQuestion,
      payload.imageBase64,
      payload.imagePageNumber,
      payload.isDeepThinking
    );

    if (!params) {
      throw Error('copilot params error');
    }

    const fakeId = genId();
    const getQuestionInfo = () => {
      const getQuestionType = () => {
        if (payload.selectText) {
          return QuestionType.SELECT_TEXT;
        }
        if (payload.imageBase64) {
          return QuestionType.IMAGE_QUESTION;
        }
        return QuestionType.ASK_QUESTION;
      };
      const question: QuestionInfo = {
        id: fakeId,
        questionType: getQuestionType(),
        question: payload.question,
        createDate: '',
        imageUrl: payload.imageBase64,
      };
      if (question.questionType === QuestionType.SELECT_TEXT) {
        question.quoteInfo = {
          quoteContent: (params as SelectTextQuestionReq).text,
        };
      } else if (payload.quoteAnswer) {
        question.quoteInfo = {
          quoteContent: payload.quoteAnswer.answer!,
        };
      }
      return question;
    };
    const fakeAnswer: AnswerInfo = {
      id: fakeId,
      createDate: '',
      isLike: false,
      isDisLike: false,
      answerStatus: AnswerStatus.PENDING,
      canRetry: false,
      createTime: '',
      errorMessage: '',
      hanleProcessedQuestion: [],
    };
    const fakeInfo: QuestionItem = {
      question: getQuestionInfo(),
      answer: fakeAnswer,
      answers: [fakeAnswer],
      req: params,
    };
    copilotState.list.push(fakeInfo);
    const curInfo = copilotState.list[copilotState.list.length - 1];
    try {
      if (
        curInfo.question!.questionType === QuestionType.IMAGE_QUESTION &&
        payload.imageBase64 &&
        !payload.imageBase64.startsWith('http')
      ) {
        // 先上传图片
        const uploadUrl = await uploadImage(
          payload.imageBase64,
          ImageStorageType.aiReadingImageQuestion
        );
        (params as ImageQuestionRequest).imageUrl = uploadUrl;
      }

      consumeBeans(BeanScenes.COPILOT_ASK);

      const info =
        curInfo.question!.questionType === QuestionType.ASK_QUESTION
          ? await askQuestion(params as AskQuestionReq)
          : curInfo.question!.questionType === QuestionType.IMAGE_QUESTION
            ? await askImageQuestion(params as ImageQuestionRequest)
            : await selectText(params as SelectTextQuestionReq);

      // curInfo.question = info.question;
      curInfo.question!.id = info.question!.id;
      curInfo.answer = info.answer;
      curInfo.answers = (info.answers || [info.answer]).map(a => a as ChatMessage);
      // copilotState.currentPollingAnswers.push(info.question!.id + '-0');
    } catch (err) {
      checkBeans(err);
      const curAnswer = curInfo.answers[0];
      curAnswer.answerStatus = AnswerStatus.ERROR;
      curAnswer.errorMessage = (err as Error)?.message;
      curAnswer.canRetry = false;
    }
    copilotState.list = [...copilotState.list];

    if (curInfo.answers[0].answerStatus === AnswerStatus.SUCCESS) {
      // ask之后就是success了，不需要polling，为了typing，做一次假polling
      curInfo.answers[0].answerStatus = AnswerStatus.PENDING;
      curInfo.answer = curInfo.answers[0];
      copilotState.list = [...copilotState.list];
      return {
        p: new Promise((resolve) => {
          setTimeout(() => {
            curInfo.answers[0].answerStatus = AnswerStatus.SUCCESS;
            curInfo.answer = curInfo.answers[0];
            resolve(void 0);
          }, 2000);
        }),
      };
    }

    return {
      p: polling(curInfo.question!.id, 0, streamPollingHandler),
    };
  };

  const store = useStore();

  const changeAnswer = async (questionId: string, cb: (curId: number) => void) => {
    check();
    const store = useStore();

    const curInfo = copilotState.list.find((item) => item.id === questionId);
    if (!curInfo) return;

    if (!curInfo.answers) {
      curInfo.answers = [];
    }

    const answerIdx = curInfo.answers.length;
    const fakeAnsId = BigInt(Date.now());
    const newAnswer: ChatMessage = {
      id: fakeAnsId,
      messageId: fakeAnsId as any,
      answer: '',
      answerStatus: AnswerStatus.PENDING,
      createAt: BigInt(Date.now()),
      canRetry: false,
      errorMessage: '',
      files: [],
      retrieverResources: [],
      rating: 0,
      conversationId: curInfo.conversationId || '',
      taskId: curInfo.taskId,
      modelType: curInfo.modelType,
      event: '',
      relatedQuestions: [],
    } as any;
    curInfo.answers.push(newAnswer);
    cb(answerIdx);

    await nextTick();

    try {
      const targetAnswer = curInfo.answers![answerIdx];
      const { connection, abort } = changeAnswerSSE(
        { questionId, lang: store.state.documents.userSettingInfo.copilotLanguage },
        (chunk) => {
          // 监听工作流事件
          if (chunk.event === 'workflow_finished' || chunk.event === 'finish') {
            userStore.refreshUserCredits(); // 刷新积分
          }
          
          // 根据消息状态分别处理，与 sendChatMessage 保持一致
          if (chunk.answerStatus === AnswerStatus.PENDING) {
            // PENDING 状态：累加答案内容
            targetAnswer.answer += chunk.answer || '';
            targetAnswer.answerStatus = AnswerStatus.PENDING;
            if (chunk.taskId) targetAnswer.taskId = chunk.taskId;
            if (chunk.conversationId) targetAnswer.conversationId = chunk.conversationId;
            if (chunk.messageId) targetAnswer.messageId = chunk.messageId;
            if (chunk.id) targetAnswer.id = chunk.id;  // 直接使用 chunk.id
            if (chunk.createAt) targetAnswer.createAt = chunk.createAt;
          } else if (chunk.answerStatus === AnswerStatus.SUCCESS) {
            // SUCCESS 状态：设置最终状态，SSE 会自动断开
            targetAnswer.answerStatus = AnswerStatus.SUCCESS;
            targetAnswer.canRetry = chunk.canRetry;
            targetAnswer.errorMessage = chunk.errorMessage;
            targetAnswer.answer += chunk.answer || '';  // 追加最后一段（如果有）
            targetAnswer.files = chunk.files;
            targetAnswer.retrieverResources = chunk.retrieverResources;
            targetAnswer.relatedQuestions = chunk.relatedQuestions;
            targetAnswer.rating = chunk.rating;
            if (chunk.conversationId) targetAnswer.conversationId = chunk.conversationId;
            if (chunk.createAt) targetAnswer.createAt = chunk.createAt;
            if (chunk.messageId) targetAnswer.messageId = chunk.messageId;
            if (chunk.id) targetAnswer.id = chunk.id;  // 直接使用 chunk.id
            if (chunk.taskId) targetAnswer.taskId = chunk.taskId;
          }
        }
      );

      copilotState.abortController = abort;
      copilotState.submitPending = true;

      await connection;
      // 状态已在 SUCCESS 回调中设置，不需要重复设置
    } catch (error) {
      const targetAnswer = curInfo.answers![answerIdx];
      if ((error as Error).name !== 'AbortError') {
        targetAnswer.answerStatus = AnswerStatus.ERROR;
        targetAnswer.errorMessage = 'Stream error';
      }
    } finally {
      copilotState.abortController = null;
      copilotState.submitPending = false;
    }
  };

  const askAgain = async (params: { questionId: string; answerId?: string; optimizeId?: string; }) => {
    console.warn('askAgain is deprecated');
  };

  const getHistoryList = async (pageSize: number) => {
    if (!copilotState.hasmore) {
      return;
    }
    const baseStore = useStore();
    const selfNoteInfo =
      baseStore.state.base.noteInfoMap[SELF_NOTEINFO_GROUPID];

    const res = await getChatMessages({
      limit: pageSize,
      firstId: copilotState.list[0]?.id,
      pdfId: selfNoteInfo.pdfId,
    });

    copilotState.limitTotal = DEFAULT_CHAMGE_ANSWER_LIMIT;

    const isEmpty =
      !copilotState.list.length && !res.data?.length && !res.hasMore;

    const list = res.data || [];

    list.forEach((item) => {
      buildMessageToQuestionItem(item);
    });
    copilotState.hasmore = res.hasMore;
    // copilotState.list = [...copilotState.list];

  };

  const buildMessageToQuestionItem = (message: ChatMessageObj) => {
    const messageRequest = message.chatMessageRequest;
    const chatMessages = message.chatMessages;
    const answers = chatMessages?.map((item) => {
      return ChatMessage.create({
        id: item.id,
        createAt: item.createAt,
        answerStatus: item.answerStatus,
        canRetry: item.canRetry,
        errorMessage: item.errorMessage,
        relatedQuestions: item.relatedQuestions,
        retrieverResources: item.retrieverResources,
        rating: item.rating,
        answer: item.answer,
        conversationId: item?.conversationId || '',
        taskId: item?.taskId || '',
        messageId: item?.messageId || '',
        modelType: item?.modelType,
      });
    });

    const questionItemTmp: QuestionItem = {
      id: messageRequest?.id,
      taskId: messageRequest?.taskId,
      noteId: messageRequest?.noteId,
      pdfId: messageRequest?.pdfId,
      modelType: messageRequest?.modelType,
      content: messageRequest?.content,
      conversationId: messageRequest?.conversationId,
      selectedText: messageRequest?.selectedText,
      uploadFiles: messageRequest?.uploadFiles,
      messageId: messageRequest?.messageId,
      quoteInfo: messageRequest?.quoteInfo,
      answers: answers || [],
    }
    copilotState.list.unshift(questionItemTmp);
  };
    
  

  const { t } = useI18n();

  const check = () => {
    if (queueLength >= MAX_QUEUE) {
      const msg = t('aiCopilot.waitTip');
      message.error(msg);
      throw Error(msg);
    }
  };

  const retry = async (questionId: string, answerIdx: number) => {
    check();
    const curInfo = copilotState.list.find((item) => {
      return item.question!.id === questionId;
    });
    const curAnswer = curInfo?.answers[answerIdx];
    if (curInfo && curAnswer) {
      if (/fake.*/.test(String(curAnswer.id)) || !curAnswer.canRetry) {
        curAnswer.answerStatus = AnswerStatus.ERROR;
        curAnswer.canRetry = false;
        curAnswer.errorMessage = '重试失败';
        throw Error('重试失败');
      }
      curAnswer.answerStatus = AnswerStatus.PENDING;
      try {
        consumeBeans(BeanScenes.COPILOT_ASK);
        await retryAiAnswerById({ answerId: String(curAnswer.id) });
      } catch (error) {
        checkBeans(error);
        curAnswer.answerStatus = AnswerStatus.ERROR;
        return;
      }
      polling(curInfo.question!.id, answerIdx);
    }
  };

      
  const sendChatMessage = async (
    payload: {
      question: string;
      selectedText?: PageSelectText[];
      continueMessage?: ChatMessage;
      images?: ImageUpload[];
      model?: string;
      conversationId?: string;
    },
  ) => {
    const fakeId = BigInt(Date.now());
    const fakeAnswer: ChatMessage = {
      id: fakeId,
      messageId: fakeId as any,
      answer: '', // 初始为空，将通过SSE流式更新
      createAt: 0,
      answerStatus: AnswerStatus.PENDING,
      canRetry: false,
      errorMessage: '',
      files: [],
      retrieverResources: [],
    } as any;

    // 从 store 获取笔记信息
    const store = useStore();
    const selfNoteInfo = store.state.base.noteInfoMap[SELF_NOTEINFO_GROUPID];
    const curNoteInfo = store.state.base.noteInfoMap[store.state.base.currentGroupId] || selfNoteInfo;


    const selectedText : SelectedText = {
      selectedText: payload.selectedText?.[0].text,
      selectedPageNum: payload.selectedText?.[0].pageNum,
      selectedBoundingBox: payload.selectedText?.[0].rects,
    }

    const uploadFiles : UploadFile[] = payload.images || [];
    const quoteInfo : QuoteInfo = {
      messageId: payload.continueMessage?.messageId,
      quoteContent: payload.continueMessage?.answer || '',
    }
    const fakeInfo: QuestionItem = {
      id: '', // 初始为空，SSE 回调中会设置为真实的 requestId
      taskId: '', // 初始为空，SSE 回调中会设置
      noteId: selfNoteInfo.noteId,
      pdfId: curNoteInfo.pdfId,
      modelType: payload.model!,
      content: payload.question.trim(),
      answers: [fakeAnswer],
      selectedText,
      uploadFiles,
      quoteInfo,
      messageId: payload.continueMessage?.messageId,
      conversationId: payload.conversationId,
    };

    copilotState.list.push(fakeInfo);
    const curInfo = copilotState.list[copilotState.list.length - 1];

    try {
      // consumeBeans(BeanScenes.COPILOT_ASK);
      const tmpInfo = { ...fakeInfo };
      delete tmpInfo.answers;
      const chatParams: ChatMessagesRequest = tmpInfo as any;

      const { connection, abort, messages } = await chatMessages(
        chatParams,
        (message) => {
          if (message.event === 'workflow_started') {
            copilotState.activeTaskId = message.taskId;
            copilotState.isActiveTask = true;
          } else if (message.event === 'workflow_finished' || message.event === 'finish') {
            copilotState.activeTaskId = '';
            copilotState.isActiveTask = false;
            copilotState.submitPending = false; // 在任务结束时重置提交状态
            userStore.refreshUserCredits(); //刷新积分
            
          }

          // 原来的消息处理方法
          if (message.answerStatus === AnswerStatus.PENDING) {
            // 确保 curInfo.answers 存在且 answer 是字符串
            if (curInfo.answers && curInfo.answers.length > 0) {
              curInfo.answers[0].answer += message.answer || '';
              curInfo.answers[0].taskId = message.taskId;
              curInfo.answers[0].conversationId = message.conversationId;
              curInfo.answers[0].messageId = message.messageId;
              if (message.id) curInfo.answers[0].id = message.id;  // 直接使用 message.id
              curInfo.answers[0].createAt = message.createAt;
              curInfo.answers[0].answerStatus = message.answerStatus;
              
              // 获取 requestId（后端已在 SSE 消息中返回）
              const requestId = (message as any).requestId;
              console.log('[DEBUG] requestId:', requestId, 'curInfo.id:', curInfo.id);
              if (requestId && (!curInfo.id || curInfo.id === '' || curInfo.id === '0' || (curInfo.id as any) === 0)) {
                curInfo.id = String(requestId);
                curInfo.taskId = message.taskId;
                console.log('[DEBUG] 已设置 curInfo.id =', curInfo.id);
              }
            }
          } else if (message.answerStatus === AnswerStatus.SUCCESS) {
            if (curInfo.answers && curInfo.answers.length > 0) {
              curInfo.answers[0].answerStatus = AnswerStatus.SUCCESS;
              curInfo.answers[0].canRetry = message.canRetry;
              curInfo.answers[0].errorMessage = message.errorMessage;
              curInfo.answers[0].answer += message.answer || '';
              if ((message as any).hanleProcessedQuestion) {
                (curInfo.answers[0] as any).hanleProcessedQuestion = (message as any).hanleProcessedQuestion;
              }
              curInfo.answers[0].files = message.files;
              curInfo.answers[0].retrieverResources = message.retrieverResources;
              curInfo.answers[0].relatedQuestions = message.relatedQuestions;
              curInfo.answers[0].rating = message.rating;
              curInfo.answers[0].conversationId = message.conversationId;
              curInfo.answers[0].createAt = message.createAt;
              curInfo.answers[0].messageId = message.messageId;
              if (message.id) curInfo.answers[0].id = message.id;  // 直接使用 message.id
              curInfo.answers[0].taskId = message.taskId;
              copilotState.submitPending = false;
            }
          }
        }
      );
      copilotState.abortController = abort;
      return {
        connection,
        abort,
        messages,
        // questionId: curInfo.question!.messageId,
      };
    } catch (err) {
      checkBeans(err);
      const curAnswer = curInfo.answers[0];
      curAnswer.answerStatus = AnswerStatus.ERROR;
      curAnswer.errorMessage = (err as Error)?.message;
      curAnswer.canRetry = false;

      copilotState.list = [...copilotState.list];
      copilotState.submitPending = false; // 在错误时重置
      copilotState.isActiveTask = false; // 在错误时重置
      throw err;
    }
  };

  const stopChat = async (taskId: string) => {
    
    const stopChatParams: StopChatMessagesRequest = {
      taskId: taskId,
    };
    const questionInfo = copilotState.list.find(
      (item) => item.answers?.find((answer) => answer.taskId === taskId && answer.answerStatus === AnswerStatus.PENDING)
    );

    const answerInfo = questionInfo?.answers?.find((answer) => answer.taskId === taskId && answer.answerStatus === AnswerStatus.PENDING);

    if (answerInfo) {
      answerInfo.answerStatus = AnswerStatus.SUCCESS;

      try {
        await stopChatMessages(stopChatParams);
        if (copilotState.abortController) {
          copilotState.abortController();
          copilotState.abortController = null;
        }
      } catch (err) {
        console.error('Failed to stop chat messages:', err);
      }

      if (answerInfo) {
        answerInfo.answerStatus = AnswerStatus.SUCCESS;
 
        copilotState.list = [...copilotState.list];
        copilotState.submitPending = false;
      }
    }
  };

  return {
    getHistoryList,
    ask,
    retry,
    check,
    copilotState,
    askAgain,
    changeAnswer,
    sendChatMessage,
    stopChat,
  };
};

export interface AskImagePayload {
  base64: string;
  pageNumber: number;
}


export const useAsk = (
  props: Readonly<{
    page?: PAGE_ROUTE_NAME;
    textValue?: string;
  }>,
  copilotState: ReturnType<typeof useCopilotState>,
  scrollToBottom: ReturnType<typeof useList>['scrollToBottom']
) => {
  const inputValue = ref('');
  const isDeepThinkingSwitch = ref(false);
  const followupAnswer = ref<null | ChatMessage>(null);
  const pageSelectText = ref<null | PageSelectText[]>(null);
  
  // 图片暂存相关状态
  const uploadedImages = ref<ImageUpload[]>([]);

  const selectTextValue = computed(() => {
    if (!pageSelectText.value) {
      return '';
    }
    return pageSelectText.value[0].text.substring(0, 8000);
  });

  const copilotHook = useCopilot(copilotState);

  const submitButtonDisabled = computed(() => !inputValue.value.trim());

  const inputRef = ref();

  const onSubmit = async (e: KeyboardEvent, model: string) => {
    if (!(e.ctrlKey || e.altKey || e.metaKey || e.shiftKey)) {
      e.preventDefault();

      if (copilotState.submitPending) {
        return;
      }
      copilotState.submitPending = true;
      
      try {
        // console.log("selectTextValue.value", pageSelectText.value);
        // console.log("followupAnswer.value", followupAnswer.value);
        // console.log("images", uploadedImages.value)
        // console.log("model", model)
        // console.log("question", inputValue.value)

        // return;
        copilotHook.sendChatMessage({
          question: inputValue.value,
          selectedText: pageSelectText.value || undefined,
          continueMessage: followupAnswer.value || undefined,
          images: uploadedImages.value,
          model,
          conversationId: copilotState.conversationId || undefined,
        });
        followupAnswer.value = null;
        inputValue.value = '';
        pageSelectText.value = null;
        clearUploadedImages()
        await nextTick();
        scrollToBottom();
      } catch (err) {
        console.error(err);
      }
    }
  };

  const onRetry = (questionId: string, answerIdx: number) => {
    copilotHook.retry(questionId, answerIdx);
  };

  const followup = (answerInfo: ChatMessage) => {
    followupAnswer.value = answerInfo;
    pageSelectText.value = null;
    inputRef.value.focus();
  };

  const clearQuote = () => {
    followupAnswer.value = null;
    pageSelectText.value = null;
  };

  emitter.on(COPILOT_ASK, (e) => {
    if (copilotState.banInputMessage) {
      // message.warning(copilotState.banInputMessage);
      return;
    }
    followupAnswer.value = null;
    pageSelectText.value = e as PageSelectText[];
  });

  const request = async (q: string) => {
    const value = q.trim();
    const isDeepThinking = isDeepThinkingSwitch.value;

    if (value) {
      copilotHook.check();
      if (props.page === PAGE_ROUTE_NAME.NOTE) {
        const quoteAnswer = followupAnswer.value;
        const selectText = pageSelectText.value;
        const req = copilotHook.ask(
          {
            question: value,
            selectText: selectText ?? undefined,
            quoteAnswer: quoteAnswer ?? undefined,
            isDeepThinking: isDeepThinking ?? false,
          },
          async () => {
            await nextTick();
            scrollToBottom();
          }
        );
        followupAnswer.value = null;
        inputValue.value = '';
        pageSelectText.value = null;
        await nextTick();
        scrollToBottom();
        const res = await req;
        submitPending.value = false;
        await nextTick();
        scrollToBottom();
        await res.p;
        await nextTick();
        scrollToBottom();
      }
    }
  };

  const submitPending = ref(false);

  const onChangeAnswer = (questionId: string, cb: (curId: number) => void) => {
    copilotHook.changeAnswer(questionId, cb);
  };

  const onOptimazeAnswer = async (
    questionId: string,
    answer: ChatMessage,
    optimizeId: string
  ) => {
    console.warn('onOptimazeAnswer needs to be refactored to use the new SSE logic.');
  };

  const removeCurrentPollingAnswers = (key: string) => {
    console.log(
      'removeCurrentPollingAnswers',
      key,
      copilotState.currentPollingAnswers
    );
    pull(copilotState.currentPollingAnswers, key);
  };

  const askImage = async (payload: { base64: string; pageNumber?: number }) => {
    copilotHook.check();
    const req = copilotHook.ask(
      {
        question: '',
        imageBase64: payload.base64,
        imagePageNumber: payload.pageNumber,
      },
      async () => {
        await nextTick();
        scrollToBottom();
      }
    );
    const res = await req;
    await nextTick();
    scrollToBottom();
    await res;
    await nextTick();
    scrollToBottom();
  };

  const addImageToUpload = (imageData: { base64: string; pageNumber?: number }): boolean => {
    if (uploadedImages.value.length >= 3) {
      return false;
    }
    uploadedImages.value.push(imageData);
    return true;
  };

  const removeUploadedImage = (index: number) => {
    if (index >= 0 && index < uploadedImages.value.length) {
      uploadedImages.value.splice(index, 1);
    }
  };

  const clearUploadedImages = () => {
    uploadedImages.value = [];
  };

  const summarySinglePaper = (onMessage?: (message: ChatMessage) => void) => {
    const store = useStore();
    const selfNoteInfo = store.state.base.noteInfoMap[SELF_NOTEINFO_GROUPID];
    const curNoteInfo = store.state.base.noteInfoMap[store.state.base.currentGroupId] || selfNoteInfo;

    const params: SummarySinglePaperRequest = {
      noteId: selfNoteInfo.noteId,
      pdfId: curNoteInfo.pdfId,
    };
    return apiSummarySinglePaper(params, onMessage);
  }

  return {
    inputValue,
    isDeepThinkingSwitch,
    followupAnswer,
    selectTextValue,
    pageSelectText,
    submitButtonDisabled,
    inputRef,
    onSubmit,
    onRetry,
    onRequest: request,
    request,
    followup,
    clearQuote,
    onChangeAnswer,
    onOptimazeAnswer,
    askImage,
    summarySinglePaper,
    addImageToUpload,
    removeUploadedImage,
    clearUploadedImages,
    uploadedImages,
    stopChat: copilotHook.stopChat
  };
};

export const useList = ({
  copilotState,
}: {
  copilotState: ReturnType<typeof useCopilotState>;
}) => {
  const copilotHook = useCopilot(copilotState);
  let flag = true;
  const { fetchState, fetch } = useFetch(async () => {
    await copilotHook.getHistoryList(LIST_PAGESIZE);
    if (flag) {
      await nextTick();
      // 因为有图片
      setTimeout(() => {
        scrollToBottom();
        setTimeout(() => {
          flag = false;
        }, 200);
      }, 500);
    }
  });

  const scrollbar = ref<{ ps: PerfectScrollbarType }>();

  const scrollToBottom = (auto?: boolean) => {
    const ps = scrollbar.value!.ps;
    ps?.update();
    setTimeout(() => {
      // console.log(
      //   ps.element.scrollTop < ps.contentHeight - ps.containerHeight - 52,
      //   ps.element.scrollTop,
      //   ps.contentHeight,
      //   ps.containerHeight,
      //   ps.element
      // );
      if (
        auto &&
        ps.element.scrollTop > ps.contentHeight - ps.containerHeight - 52
      ) {
        return;
      }
      ps.element.scrollTop = ps.contentHeight + ps.containerHeight;
    }, 100);
  };

  const onScrollTop = async () => {
    if (fetchState.pending || flag) {
      return;
    }
    const ps = scrollbar.value!.ps;
    const lastContentHeight = ps.contentHeight;
    await fetch();
    await nextTick();
    ps.element.scrollTop =
      ps.contentHeight - lastContentHeight - ps.containerHeight + 20;
  };

  const loading = computed(() => {
    return copilotHook.copilotState.hasmore;
  });

  return {
    fetchState,
    fetch,
    scrollbar,
    scrollToBottom,
    onScrollTop,
    loading,
  };
};
