<template>
  <div ref="containerRef" class="copilot-container">
    <!-- 内部内容  -->
    <PerfectScrollbar ref="scrollbar" class="copilot-ps" @ps-y-reach-start="onScrollTop">
      <div v-if="page === PAGE_ROUTE_NAME.NOTE && loading && copilotStore.isChat" class="copilot-loading">
        <a-spin />
      </div>
      <div class="copilot-list" v-if="copilotStore.isChat">
        <div v-for="item in copilotState.list" :key="item.id" class="copilot-item">
          <div class="copilot-item__simple">
            <!-- 显示问题部分 -->
            <div v-if="item.content">
              <!-- 如果有引用内容，先显示引用内容 -->
              <div v-if="item.selectedText" class="inbox question select">
                {{ item.selectedText.selectedText }}
              </div>
              <!-- 如果有引用内容，先显示引用内容 -->
              <div v-if="item.quoteInfo?.messageId" class="inbox question select">
                {{ item.quoteInfo?.quoteContent }}
              </div>
              <!-- 显示问题内容 -->
              <div class="inbox question">
                <div>{{ item.content }}</div>
                <div v-if="item.uploadFiles?.length" class="question-image-flex">
                  <span v-for="file in item.uploadFiles" class="question-image-box">
                    <img :src="file.base64" />
                  </span>
                </div>
              </div>
            </div>

            <div v-if="item.answers?.length" class="inbox answer">
              <Answer :change-answer-limit="copilotState.limitTotal" :question-id="item.id"
                :scroll-ele="scrollbar?.ps.element as HTMLDivElement" :multi="!!item.content" :answers="item.answers"
                :scroll-to-bottom="scrollToBottom" :on-quick-ask-question="onQuickAskQuestion" @retry="onRetry"
                @followup="followup" @change-answer="onChangeAnswer" @optimaze-answer="onOptimazeAnswer" />
            </div>
          </div>
        </div>
      </div>
      <div v-else class="copilot-container no-permission-container">
        <div class="copilot-description">
          <!-- <img class="ai-img" src="@/assets/images/right/ai.png" alt="" /> -->
          <h2 style="color: var(--site-theme-text-primary)">
            {{ $t('aiCopilot.chatFlow') }}
          </h2>
          <p style="color: var(--site-theme-text-primary)">
            {{ $t('aiCopilot.askSuggestions') }}
          </p>
          <div class="text">
            {{ $t('aiCopilot.paperSummary') }}
            <p v-if="isLoadingSummary" class="loading-text">
              <a-spin size="small" />
              <span class="ml-2">{{ $t('aiCopilot.generatingSummary') }}</span>
            </p>
            <div v-else>
              <div v-if="summaryMessages.length > 0" class="summary-list">
                <PendingAnswer v-if="
                    curAnswer?.answerStatus === AnswerStatus.SUCCESS ||
                    curAnswer?.answerStatus === AnswerStatus.PENDING
                  " ref="answerComp" :key="curAnswer?.conversationId" :answer="{
                    answer: formattedSummary || '',
                    anchors: curAnswer?.relatedQuestions,
                    status:
                      curAnswer?.answerStatus === AnswerStatus.SUCCESS
                        ? 'success'
                        : 'pending',
                  }" :scroll-to-bottom="autoScrollBottom" @quick-ask-question="onQuickAskQuestion"
                  @textShowComplete="handleTextShowComplete" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </PerfectScrollbar>
    
    <!-- 底部输入区-->
    <div class="copilot-bottom">
      <Suggestion v-if="!banInputMessage && !followupAnswer && !selectTextValue" @suggestion:fill="onFillInput" />
      <div class="copilot-input">
        <div class="copilot-tools-container">
          <div class="tools-left">
            <ScreenShot :clip-selecting="clipSelecting" :clip-action="clipAction"
              :add-image-to-upload="addImageToUpload" />
          </div>
          <div class="tools-right">
            <!-- <Language /> -->
            <!---Feekback 链接 --> 
            <!-- <a href="https://docs.qq.com/doc/DVUZmcnpuTFFlbmhW" target="_blank" class="join-test-group">
              {{ $t('aiCopilot.joinGroup') }}
              <ExclamationCircleOutlined style="margin-left: 4px" />
            </a> -->
            <!-- chat状态灯 -->
            <div class="status-light" :class="{ 'is-running': copilotState.submitPending }"></div>
          </div>
        </div>
        <a-alert v-if="banInputMessage" :message="banInputMessage" show-icon type="warning" class="copilot-ban" />
        <div class="input-wrapper">
          <div v-if="uploadedImages.length > 0" class="image-preview-container">
            <div v-for="(image, index) in uploadedImages" :key="index" class="image-preview-item">
              <div class="image-thumbnail">
                <img :src="image.base64" :alt="$t('aiCopilot.uploadedImage')" />
                <div class="image-delete-btn" @click="handleRemoveImage(index)">
                  <close-outlined />
                </div>
              </div>
            </div>
            <div v-if="uploadedImages.length < 3" class="image-count-hint">
              {{ uploadedImages.length }}/{{ $t('aiCopilot.maxImages') }}
            </div>
          </div>
          <div v-if="!banInputMessage && (followupAnswer || selectTextValue)" class="copilot-followup-box">
            <div v-if="followupAnswer" class="copilot-followup-box-answer">
              {{ followupAnswer.answer }}
            </div>
            <div v-else class="copilot-followup-box-answer">
              {{ selectTextValue }}
            </div>
            <div class="copilot-followup-box-close" @click="clearQuote">
              <div class="copilot-followup-box-inner">
                <close-outlined />
              </div>
            </div>
          </div>
          <a-textarea ref="inputRef" v-model:value="inputValue" show-count :placeholder="$t('aiCopilot.askTip')"
            :rows="inputRows" :disabled="!!banInputMessage || !isEmbeddingReady" :maxlength="200" @focus="isInputFocus = true"
            @blur="isInputFocus = false" @pressEnter="changeComponentAndSubmit" />
          <PaidPlanButton v-if="!copilotStore.accessAiCopilot" class="paid-button" />
          <!-- 当正在提交时显示停止按钮 -->
          <a-button v-else-if="copilotStore.accessAiCopilot && copilotState.submitPending" type="primary" danger
            class="ask-button" @click="stopCurrentChat">
            <div class="stop-icon">
              <div class="stop-icon-square"></div>
            </div>
          </a-button>
          <!-- 当未提交时显示提问按钮 -->
          <a-button v-else type="primary" class="ask-button" :disabled="submitButtonDisabled || !isEmbeddingReady"
            @click="changeComponentAndSubmit">
            {{ $t('aiCopilot.ask') }}
          </a-button>
        </div>
        <div class="btn">
          <div class="copilot-controls">
            <div class="left-controls">
              <CreditHref class="control-item" />
              <div
                class="thinking-switch control-item summary-button"
                @click="copilotStore.isChat ? gotoSummary() : gotoChat()"
              >
                <a-tooltip
                  :title="
                    copilotStore.isChat ? $t('aiCopilot.summaryChatButtonTextSummary') : $t('aiCopilot.summaryChatButtonTextChat')
                  "
                  placement="top"
                >
                  <file-text-outlined v-if="copilotStore.isChat" />
                  <message-outlined v-else />
                </a-tooltip>
              </div>
            </div>
            <div class="right-controls">
              <div class="thinking-switch control-item">
                <div class="model-dropdown">
                  <div class="selected-model" @click="toggleModelDropdown">
                    <span class="model-text">{{ selectedModel }}</span>
                    <up-outlined class="up-icon" />
                  </div>
                  <div v-if="showModelDropdown" class="model-dropdown-menu" @click.stop>
                    <div class="model-search">
                      <search-outlined />
                      <input type="text" :placeholder="$t('aiCopilot.searchModels')" v-model="modelSearchText" />
                    </div>
                    <div class="model-list">
                      <div v-for="model in filteredModels" :key="model.value" class="model-item"
                        :class="{ active: selectedModel === model.value }" @click="selectModel(model.value)">
                        <div class="model-info">
                          <div class="model-name">{{ model.label }}</div>
                          <div class="model-credit">{{ model.credit }}</div>
                        </div>
                        <check-outlined v-if="selectedModel === model.value" class="model-check" />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { computed, ref, onMounted, watch, onUnmounted, nextTick } from 'vue'
import {
  CloseOutlined,
  ExclamationCircleOutlined,
  FileTextOutlined,
  UpOutlined,
  SearchOutlined,
  CheckOutlined,
  MessageOutlined,
} from '@ant-design/icons-vue'
import { useAsk, useList, useCopilotState } from '~/src/hooks/useCopilot'
import {
  ElementName,
  PageType,
  reportElementImpression,
} from '~/src/api/report'
import { useResizeObserver } from '@vueuse/core'
import { PAGE_ROUTE_NAME } from '~/src/routes/type'
import { useUserStore } from '@/common/src/stores/user'
import Answer from './Answer.vue'
import { useCopilotStore } from '~/src/stores/copilotStore';
import { useDocStore } from '~/src/stores/docStore';
import { UserDocParsedStatusEnum } from 'go-sea-proto/gen/ts/doc/UserDocParsedStatus';
import { storeToRefs } from 'pinia';
// import Language from './components/Language.vue'
import Suggestion from './components/Suggestion.vue'
import ScreenShot from './components/ScreenShot.vue'
import PendingAnswer from './components/PendingAnswer.vue'
import { useClip } from '~/src/hooks/useHeaderScreenShot'
import { ChatMessage } from 'go-sea-proto/gen/ts/copilot/chat/ChatFlow'
import { AnswerStatus } from 'go-sea-proto/gen/ts/copilot/chat/ChatFlow'
import CreditHref from '@common/components/Pay/CreditHref.vue'
import PaidPlanButton from '@common/components/Pay/PaidPlanButton.vue'

const props = withDefaults(
  defineProps<{
    clipSelecting: boolean
    clipAction: ReturnType<typeof useClip>['clipAction']
    page?: PAGE_ROUTE_NAME.CHAT | PAGE_ROUTE_NAME.NOTE | PAGE_ROUTE_NAME.WRITE
    textValue?: string
  }>(),
  {
    page: PAGE_ROUTE_NAME.NOTE,
    textValue: '',
  }
)
const userStore = useUserStore()
console.log(userStore)

const copilotStore = useCopilotStore()

const copilotState = useCopilotState()

const banInputMessage = computed(() => copilotState.banInputMessage)

const { scrollbar, scrollToBottom, onScrollTop, loading } = useList({
  copilotState,
})

const {
  followupAnswer,
  selectTextValue,
  onRetry,
  onRequest,
  onSubmit,
  stopChat,
  addImageToUpload,
  removeUploadedImage,
  inputValue,
  followup,
  clearQuote,
  submitButtonDisabled,
  inputRef,
  onChangeAnswer,
  onOptimazeAnswer,
  summarySinglePaper,
  uploadedImages,
} = useAsk(props, copilotState, scrollToBottom)

const containerRef = ref()
useResizeObserver(containerRef, () => {
  scrollbar.value?.ps.update()
})

const isInputFocus = ref(false)

const inputRows = ref(2);

const docStore = useDocStore();
const { docInfo, isEmbeddingReady } = storeToRefs(docStore);

const summaryMessages = ref<string[]>([])
const isLoadingSummary = ref<boolean>(false)

const curAnswer = ref<ChatMessage | null>(null)

const formattedSummary = computed(() => {
  return summaryMessages.value.join('')
})

// 自动滚动到底部
const autoScrollBottom = () => {
  // 如果有滚动容器，滚动到底部
  if (scrollbar.value?.ps) {
    const container = scrollbar.value.ps.element
    container.scrollTop = container.scrollHeight
  }
}

// 处理快速提问
const onQuickAskQuestion = (question: string) => {
  // 将问题文本放入输入框，而不是直接提问
  if (question) {
    inputValue.value = question
    // 聚焦到输入框
    nextTick(() => {
      inputRef.value?.focus()
    })
  }
}

// 处理自定义问题填入输入框
const onFillInput = (question: string) => {
  // 将自定义问题文本放入输入框
  if (question) {
    inputValue.value = question
    // 聚焦到输入框
    nextTick(() => {
      inputRef.value?.focus()
    })
  }
}

// 处理文本显示完成
const handleTextShowComplete = (isNoScrollBottom?: boolean) => {
  if (!isNoScrollBottom) {
    autoScrollBottom()
  }
}

const fetchPaperSummary = () => {
  isLoadingSummary.value = true
  summaryMessages.value = [] // 清空之前的消息

  const sse = summarySinglePaper((message: ChatMessage) => {
    console.log('Received message:', message)
    // 只有当 answer 不为空时才添加到数组
    if (message.answer && message.answer.trim() !== '') {
      summaryMessages.value.push(message.answer)
    }

    // 收到消息后立即设置加载状态为 false，显示内容
    if (summaryMessages.value.length > 0) {
      isLoadingSummary.value = false
    }
    curAnswer.value = message
  })

  // 设置一个超时，确保即使没有收到有效消息也会结束加载状态
  setTimeout(() => {
    isLoadingSummary.value = false
  }, 20000) // 20秒后自动结束加载状态

  // 可以在组件卸载时中止连接
  onUnmounted(() => {
    sse.abort()
  })
}

const changeComponentAndSubmit = (e: KeyboardEvent) => {
  if (!copilotStore.isChat) {
    copilotStore.gotoChat()
  }
  onSubmit(e, selectedModel.value)
}

// 跳转到总结模式
const gotoSummary = () => {
  // if (!isEmpty) {
  //   useCopilotStore().gotoChat();
  // }
  copilotStore.gotoSummary()
  fetchPaperSummary()
}

// 跳转到对话模式
const gotoChat = () => {
  copilotStore.gotoChat()
}

const selectedModel = ref('gpt-4o-mini')
const showModelDropdown = ref(false)
const modelSearchText = ref('')

// 从用户权限中获取模型列表
const models = computed(() => {
  const modelList = userStore.getAiModelList() || [];
  
  // 如果没有获取到模型列表或列表为空，则使用默认模型列表
  if (!modelList || modelList.length === 0) {
    return [
      // { value: 'gpt-4o-mini', label: 'GPT-4o-mini', credit: 'free' },
      // { value: 'gpt-4o', label: 'GPT-4o', credit: '0.25x credit' },
    ];
  }
  
  // 将后端返回的模型列表转换为前端所需的格式，并过滤掉isEnable为false的模型
  return modelList
    .filter(model => model.isEnable !== false) // 只显示isEnable不为false的模型
    .map(model => ({
      value: model.key || model.name,
      label: model.name || model.key,
      credit: model.isFree ? 'free' : `${Number(model.creditCost) / 100}x credit`
    }));
})

const filteredModels = computed(() => {
  if (!modelSearchText.value) return models.value
  return models.value.filter((model) =>
    model.label.toLowerCase().includes(modelSearchText.value.toLowerCase())
  )
})

const toggleModelDropdown = (event: MouseEvent) => {
  event?.stopPropagation()
  showModelDropdown.value = !showModelDropdown.value
}

const selectModel = (value) => {
  selectedModel.value = value
  showModelDropdown.value = false
}

const isClickStop = ref(false)

// 停止当前聊天 (旧版逻辑，已注释)
// const stopCurrentChat = () => {
//   // 获取最新的消息ID
//   const latestQuestion = copilotState.list[copilotState.list.length - 1]
//   if (
//     latestQuestion &&
//     latestQuestion.answers &&
//     latestQuestion.answers.length > 0
//   ) {
//     let taskId = ''
//     latestQuestion.answers.forEach((answer) => {
//       if (answer.answerStatus === AnswerStatus.PENDING) {
//         taskId = answer.taskId
//         return // 注意：这里的 return 无法跳出 forEach 循环
//       }
//     })
//     console.log('停止当前聊天', taskId)
//     if (taskId && !isClickStop.value) {
//       console.log('停止当前聊天', taskId)
//       stopChat(taskId)
//       isClickStop.value = true
//     }
//     setTimeout(() => {
//       isClickStop.value = false
//     }, 3000)
//   }
// }

// 停止当前聊天 (新版重构逻辑)
const stopCurrentChat = () => {
  // if (isClickStop.value) return // 如果正在停止中，则直接返回

  const activeTaskId = copilotState.activeTaskId

  console.log('准备停止聊天，任务ID:', activeTaskId)
  
  if (activeTaskId) {
    // isClickStop.value = true // 锁定按钮
    stopChat(activeTaskId)
    console.log('已调用 stopChat, taskId:', activeTaskId)

    // 3秒后自动解锁，防止异常情况下永久锁定
    // setTimeout(() => {
    //   isClickStop.value = false
    // }, 200)
  } else {
    console.info('未能找到待处理的任务ID，无法停止。')
    // 如果没找到taskId，不应该锁定按钮，以便用户可以重试
    // isClickStop.value = false
  }
}


const handleRemoveImage = (index: number) => {
  removeUploadedImage(index)
}

onMounted(() => {
  setTimeout(() => {
    if (copilotState.list && copilotState.list.length > 0) {
      copilotStore.gotoChat()
    }
  }, 2000)
  // 自动获取摘要
  if (copilotStore.isSummary) {
    fetchPaperSummary()
  }
  // 添加全局点击事件监听器
  document.addEventListener('click', () => {
    if (showModelDropdown.value) {
      showModelDropdown.value = false
    }
  })
})

watch(
  () => copilotStore.accessAiCopilot,
  () => {
    if (!copilotStore.accessAiCopilot) {
      reportElementImpression({
        page_type: PageType.note,
        type_parameter: 'none',
        element_name: ElementName.upperAIAssistPopup,
        element_parameter: 'none',
        // module_type: 'none',
      })
    }
  }
)
</script>
<style lang="less" scoped>
.copilot {
  &-container {
    position: relative;
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  &-language {
    :deep(.ant-select-selector) {
      border: none !important;
    }
  }

  &-loading {
    text-align: center;
  }

  &-tip {
    color: var(--site-theme-text-tertiary);
    padding: 5px 16px;
    line-height: 22px;
    border-bottom: 1px solid var(--site-theme-border-color);
    margin-bottom: 16px;
    font-size: 13px;
    display: flex;
    justify-content: space-between;
    flex-wrap: wrap;

    .content {
      cursor: pointer;
    }

    .join-test-group {
      margin-left: 12px;
      color: var(--site-theme-text-tertiary);
    }
  }

  &-ps {
    flex: 1;
  }

  &-followup-box {
    display: flex;
    align-items: center;
    margin: 8px 8px;
    position: relative;

    &-answer {
      display: -webkit-box;
      overflow: hidden;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      border-left: 3px solid var(--site-theme-pdf-panel-blockquote);
      overflow: hidden;
      background-color: var(--site-theme-pdf-panel-secondary);
      padding: 0 4px;
      color: var(--site-theme-text-secondary);
      font-style: italic;
      width: 100%;
    }

    &:hover {
      .copilot-followup-box-close {
        opacity: 1;
      }
    }

    &-inner {
      width: 16px;
      height: 16px;
      border-radius: 50%;
      background: var(--site-theme-text-tertiary);
      display: flex;
      align-items: center;
      justify-content: center;
    }

    &-close {
      position: absolute;
      right: 0;
      top: 0;
      width: 20px;
      height: 20px;
      opacity: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      transform: translate(50%, -50%);
      z-index: 2;

      .anticon {
        font-size: 14px;
        transform: scale(0.5);
      }
    }
  }

  &-item {
    overflow: hidden;
    margin-bottom: 24px;

    .inbox {
      border-radius: 2px;
      max-width: 80%;
      padding: 10px 12px;
      clear: both;
      word-wrap: break-word;
    }

    &__simple {
      overflow: hidden;

      .question {
        float: right;
        background-color: var(--site-theme-brand);
        color: var(--site-theme-text-inverse);
        margin-bottom: 24px;
        margin-right: 12px;

        &.select {
          background-color: var(--site-theme-pdf-panel-secondary);
          border-left: 4px solid var(--site-theme-pdf-panel-blockquote);
          color: var(--site-theme-text-secondary);
          padding: 0 8px;
          font-style: italic;
          -webkit-line-clamp: 10;
          display: -webkit-box;
          overflow: hidden;
          -webkit-box-orient: vertical;
        }

        .quote {
          border-left: 3px solid var(--site-theme-pdf-panel-ref-content);
          color: var(--site-theme-text-tertiary);
          margin-bottom: 8px;
          padding-left: 8px;
          display: -webkit-box;
          overflow: hidden;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
        }

        .question-image-flex {
          display: flex;
          flex-direction: row;
          gap: 4px;
          margin-top: 4px;
          overflow-x: auto;

          .question-image-box {
            flex: 0 0 auto;
            width: 50px;
            height: 50px;
            margin: 0;
            padding: 3px;

            img {
              width: 100%;
              height: 100%;
              object-fit: cover;
              border-radius: 2px;
            }
          }
        }
      }

      .question+.question {
        margin-top: -20px;
      }

      .answer {
        float: left;

        .error {
          color: var(--site-theme-error);
        }
      }
    }
  }

  &-bottom {
    padding: 16px 8px;
  }

  &-input {
    border-top: 1px solid var(--site-theme-border-color);
    margin-top: 6px;

    :deep(textarea) {
      border: none;
      resize: none;
      /* 禁用 textarea 拉伸功能 */
      background-color: var(--site-theme-bg-primary) !important;
      color: var(--site-theme-text-primary) !important;

      &:hover,
      &:focus,
      &:active {
        border: none;
        box-shadow: none;
      }

      &::placeholder {
        color: var(--site-theme-text-tertiary) !important;
      }
    }

    :deep(.ant-input-textarea-show-count::after) {
      float: left !important;
      margin-top: 16px;
      margin-left: 12px;
      color: var(--site-theme-text-tertiary) !important;
    }

    .input-wrapper {
      position: relative;
      margin-bottom: 10px;

      .image-preview-container {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        padding-left: 14px;
        margin-bottom: 8px;

        .image-preview-item {
          position: relative;
          margin: 2px;

          .image-thumbnail {
            width: 50px;
            height: 50px;
            position: relative;

            img {
              width: 100%;
              height: 100%;
              border-radius: 4px;
              object-fit: cover;
            }
          }

          .image-delete-btn {
            position: absolute;
            top: -6px;
            right: -6px;
            width: 14px;
            height: 14px;
            opacity: 1;
            display: flex;
            align-items: center;
            justify-content: center;
            z-index: 2;
            background-color: var(--site-theme-bg-primary);
            border-radius: 50%;
            box-shadow: var(--site-theme-shadow-2);
            cursor: pointer;
            transition: all 0.2s;

            &:hover {
              background-color: var(--site-theme-bg-secondary);
              transform: scale(1.1);
            }

            .anticon {
              font-size: 10px;
            }
          }
        }

        .image-count-hint {
          color: var(--site-theme-text-tertiary);
          font-size: 12px;
          flex: 0 0 auto;
          align-self: flex-end;
        }
      }

      :deep(.ant-input) {
        padding-right: 75px;
        /* 为按钮预留空间 */
      }

      :deep(.ant-input-textarea-show-count) {

        /* 确保文本区域有足够的右侧内边距 */
        &::after {
          right: 75px;
          /* 调整字数统计的位置，避免与按钮重叠 */
        }

        /* 确保文本不会被按钮遮挡 */
        textarea {
          padding-right: 75px;
        }
      }

      .paid-button {
        position: absolute;
        right: 0;
        bottom: 24px;
        z-index: 10;
        margin-right: 4px;

        /* 覆盖基本样式以适应输入框 */
        padding: 4px 6px;
        height: 32px;
        line-height: 24px;
        font-size: 14px;
        border-radius: 4px;
        background-color: var(--site-theme-brand);
        color: var(--site-theme-text-inverse);
      }

      .ask-button {
        position: absolute;
        right: 0;
        bottom: 20px;
        z-index: 10;
        border-radius: 4px;
        margin-right: 12px;

        /* 悬停时增加不透明度 */
        &:hover {
          opacity: 1;
        }

        /* 方案三：添加交互动画 */
        transition: opacity 0.2s,
        background-color 0.2s;

        /* 当文本框获得焦点且有内容时，按钮更加突出 */
        .input-wrapper:has(textarea:focus) & {
          opacity: 1;
          background-color: var(--site-theme-brand);
        }

        .stop-icon {
          width: 24px;
          height: 24px;
          display: flex;
          justify-content: center;
          align-items: center;
        }

        .stop-icon-square {
          width: 12px;
          height: 12px;
          background-color: var(--site-theme-bg-primary);
        }
      }
    }

    .btn {
      margin-top: 12px;

      .copilot-controls {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding-left: 16px;
      }

      .left-controls {
        display: flex;
        align-items: center;
        color: var(--site-theme-text-tertiary);
      }

      .right-controls {
        display: flex;
        align-items: center;
      }

      .control-item {
        margin-right: 16px;
      }

      .thinking-switch {
        display: flex;
        align-items: center;
      }

      .switch-item {
        margin-right: 8px;
      }

      .switch-label {
        white-space: nowrap;
      }
    }

    .ant-btn-primary[disabled] {
      background-color: var(--site-theme-bg-mute);
      color: var(--site-theme-text-tertiary);
      border: none;
    }
  }

  &-ban.ant-alert-warning {
    background: var(--site-theme-bg-primary) !important;
    margin: 2px 0;
    box-shadow: var(--site-theme-shadow-3);
    border-radius: 2px;
    border: none;

    :deep(.ant-alert-message) {
      color: var(--site-theme-text-primary);
    }
  }
}

.no-permission-container {
  .copilot-description {
    padding: 20px;
  }

  .ai-img {
    width: 80px;
  }

  .text {
    font-family: Noto Sans SC;
    font-size: 15px;
    font-weight: 400;
    line-height: 24px;
    color: var(--site-theme-text-secondary);
    margin: 13px 0 16px;
    // text-align: center;
  }

  .ant-btn {
    min-width: 102px;
  }

  .fetch-summary-btn {
    background: var(--site-theme-brand);
    color: var(--site-theme-text-inverse);
    border: none;
    padding: 6px 12px;
    border-radius: 4px;
    cursor: pointer;

    &:hover {
      background: var(--site-theme-brand-hover);
    }
  }
}

.summary-button {
  cursor: pointer;
  color: var(--site-theme-text-secondary);
  font-size: 18px;

  &:hover {
    color: var(--site-theme-brand-hover);
  }
}

.model-dropdown {
  position: relative;
  min-width: 100px;
  user-select: none;
}

.selected-model {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 8px;
  cursor: pointer;
  font-size: 14px;
  color: var(--site-theme-text-primary);
  background-color: var(--site-theme-bg-primary);
  border-radius: 4px;
  border: 1px solid var(--site-theme-border-color);

  &:hover {
    border-color: var(--site-theme-brand);
  }

  .model-text {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    display: inline-block;
  }
}

.up-icon {
  margin-left: 6px;
  font-size: 12px;
  opacity: 0.7;
  color: var(--site-theme-text-tertiary);
}

.model-dropdown-menu {
  position: absolute;
  bottom: 100%;
  right: 0;
  margin-bottom: 4px;
  background-color: var(--site-theme-bg-primary);
  border-radius: 4px;
  box-shadow: var(--site-theme-shadow-2);
  z-index: 1000;
  width: 280px;
  font-size: 12px;
  border: 1px solid var(--site-theme-border-color);
}

.model-search {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid var(--site-theme-border-color);

  .anticon-search {
    color: var(--site-theme-text-tertiary);
  }
}

.model-search input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: var(--site-theme-text-primary);
  margin-left: 8px;
  font-size: 14px;
}

.model-list {
  max-height: 300px;
  overflow-y: auto;
}

.model-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 5px 8px;
  cursor: pointer;
}

.model-item:hover {
  background-color: var(--site-theme-bg-mute);
}

.model-item.active {
  background-color: var(--site-theme-brand-light);
}

.model-info {
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.model-name {
  color: var(--site-theme-text-primary);
}

.model-credit {
  color: var(--site-theme-text-secondary);
  font-size: 13px;
}

.model-check {
  color: var(--site-theme-text-primary);
  padding-left: 4px;
}

.copilot-tools-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.tools-left {
  display: flex;
  align-items: center;
  color: var(--site-theme-text-primary);
}

.tools-right {
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.status-light {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background-color: #2a2a2a; /* 外圈深灰色 */
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  box-shadow: 0 0 2px rgba(0, 0, 0, 0.6), inset 0 0 1px rgba(255, 255, 255, 0.1);
  position: relative;
  transition: all 0.3s ease;
}

.status-light::before {
  content: '';
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #8c8c8c; /* 内圈浅灰色 */
  transition: background-color 0.3s ease;
  box-shadow: 0 0 3px 1px rgba(0, 0, 0, 0.3);
}

.status-light.is-running::before {
  background-color: #52c41a; /* 运行时内圈为绿色 */
}

</style>
