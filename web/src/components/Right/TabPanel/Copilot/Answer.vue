<template>
  <div v-if="answers.length > 1" class="switch-answer" @click="switchAnswer">
    <span class="page">{{ curIdx + 1 }}/{{ answers.length }}</span
    ><swap-outlined />
  </div>
  <div class="text relative">
    <PendingAnswer
      v-if="
        curAnswer?.answerStatus === AnswerStatus.SUCCESS ||
        curAnswer?.answerStatus === AnswerStatus.PENDING
      "
      ref="answerComp"
      :key="curAnswer.id"
      :answer="{
        answer: curAnswer.answer === 'null' ? '' : curAnswer.answer || '',
        anchors: curAnswer.hanleProcessedQuestion,
        status:
          curAnswer.answerStatus === AnswerStatus.SUCCESS
            ? 'success'
            : 'pending',
      }"
      :scroll-to-bottom="autoScrollBottom"
      @quick-ask-question="onQuickAskQuestion"
      @textShowComplete="handleTextShowComplete"
    />
    <!-- <MarkdownViewer
      v-if="curAnswer?.answerStatus === AnswerStatus.SUCCESS || curAnswer?.answerStatus === AnswerStatus.PENDING"
      :raw="(showAnswer || '') + (isPendingAnswer ? blinkingSpan : '') " /> -->
    <div v-if="curAnswer.answerStatus === AnswerStatus.ERROR" class="error">
      <a v-if="curAnswer!.canRetry" class="error" @click="onRetry"
        >{{
          curAnswer.errorMessage || $t('aiCopilot.errorTip')
        }}&nbsp;<reload-outlined
      /></a>
      <span v-else>{{
        curAnswer.errorMessage || $t('aiCopilot.errorTip')
      }}</span>
    </div>
    <div
      v-if="
        curAnswer!.answerStatus === AnswerStatus.SUCCESS &&
          answerOperateList[curIdx]?.isTextShowComplete
      "
      class="btns"
    >
      <!-- <div>
        <UpDown
          v-if="curAnswer!.id"
          :answer="curAnswer"
          :scroll-ele="scrollEle"
        />
      </div> -->
      <div v-if="multi" class="icons flex gap-2 justify-end min-w-[80px]">
        <a-tooltip
          v-if="answers.length < changeAnswerLimit && curIdx === 0"
          :title="$t('aiCopilot.changeAnswer')"
        >
          <a @click="onChangeAnswer"
            ><redo-outlined
          /></a>
        </a-tooltip>
        <a-tooltip :title="$t('common.text.copy')">
          <a @click="onCopy"><copy-outlined /></a>
        </a-tooltip>
      </div>
    </div>
    <!-- 满意度调查-TODO:暂时隐藏 -->
    <!-- <div
      v-if="
        curAnswer!.answerStatus === AnswerStatus.SUCCESS &&
          answerOperateList[curIdx]?.isTextShowComplete
      "
      class="satisfaction mt-2"
    >
      <div class="h-px bg-white opacity-10" />
      <template v-if="!answerOperateList[curIdx]?.isFeeback">
        <div class="my-2 text-xs leading-6 opacity-70">
          {{ $t('aiCopilot.satisfacteAnswer') }}
        </div>
        <div class="list-wrap flex flex-wrap gap-1 mb-1">
          <div
            v-for="item in satisfactionList"
            :key="item.key"
            class="border border-solid border-[#66ABFF] text-[#66ABFF] rounded-sm opacity-80 cursor-pointer w-14 text-center text-xs leading-6 font-semibold hover:opacity-100"
            @click="submitFeedback(item.value)"
          >
            {{ $t(`aiCopilot.${item.key}`) }}
          </div>
        </div>
      </template>
      <div v-else class="mt-2 -mb-2 text-xs leading-6 text-center opacity-70">
        {{ $t('aiCopilot.feedback') }}
      </div>
    </div> -->
  </div>
  <div
    v-if="
      curAnswer.answerStatus === AnswerStatus.SUCCESS && curAnswer.messageId
    "
    class="followup"
  >
    <a @click="onFollowup"
      ><i class="aiknowledge-icon icon-reply" />{{
        $t('aiCopilot.continueAsk')
      }}</a
    >
  </div>
</template>
<script setup lang="ts">
// import { MarkdownViewer } from '@idea/aiknowledge-markdown';
import { ref, computed, watch, nextTick } from 'vue'
// import UpDown from './UpDown.vue';
import Optimize from './Optimize.vue'
import {
  RedoOutlined,
  ReloadOutlined,
  SwapOutlined,
  CopyOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import PendingAnswer from './components/PendingAnswer.vue'
import {
  ElementClick,
  PageType,
  reportClick,
  reportModuleImpression,
} from '@/api/report'
import { satisfactionFeedback } from '@/api/copilot'
import { useStore } from '@/store'
import { Satisfaction } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiAssistReading'
import { useI18n } from 'vue-i18n'
import { copyToPaste } from '~/src/util/copy'
import {
  ChatMessage,
  AnswerStatus,
} from 'go-sea-proto/gen/ts/copilot/chat/ChatFlow'

const props = defineProps<{
  questionId: string
  answers: ChatMessage[]
  scrollEle: HTMLDivElement
  multi: boolean
  changeAnswerLimit: number
  scrollToBottom: (auto: boolean) => void
  onQuickAskQuestion: (question: string) => void
}>()

const i18n = useI18n()
const store = useStore()
const answerComp = ref()
const answerOperateList = ref<
  { isFeeback: boolean; isTextShowComplete: boolean }[]
>([])
const curIdx = ref(0)

const curAnswer = computed(() => props.answers[curIdx.value])

const satisfactionList = [
  {
    key: 'excellent',
    value: Satisfaction.EXCELLENT,
  },
  {
    key: 'good',
    value: Satisfaction.GOOD,
  },
  {
    key: 'average',
    value: Satisfaction.AVERAGE,
  },
  {
    key: 'inferior',
    value: Satisfaction.INFERIOR,
  },
  {
    key: 'terrible',
    value: Satisfaction.TERRIBLE,
  },
]

watch(
  () => curAnswer.value.answerStatus,
  () => {
    if (curAnswer.value.answerStatus === AnswerStatus.SUCCESS) {
      const noteInfo =
        store.state.base.noteInfoMap[store.state.base.currentGroupId]
      // 出现了满意度调查框 需要上报
      reportModuleImpression({
        page_type: PageType.note,
        type_parameter: noteInfo.pdfId,
        module_type: 'popup_user_feedback',
      })
    }
  },
  { immediate: true }
)

const handleTextShowComplete = (isNoScrollBottom?: boolean) => {
  answerOperateList.value[curIdx.value] = {
    isFeeback: false,
    isTextShowComplete: true,
  }
  if (!isNoScrollBottom) {
    nextTick(() => {
      autoScrollBottom()
    })
  }
}

const submitFeedback = (satisfaction: Satisfaction) => {
  satisfactionFeedback({
    answerId: curAnswer.value.id,
    satisfaction,
  }).then(() => {
    message.success(i18n.t('aiCopilot.feedbackSucc'))
    answerOperateList.value[curIdx.value].isFeeback = true
  })
}

const autoScrollBottom = () => {
  const isScrollToBottomAuto = curIdx.value > 0
  props.scrollToBottom(isScrollToBottomAuto)
}

// const blinkingSpan = '<span class="copilot-blinking"></span>'

// const isPendingAnswer = computed(() =>
//   curAnswer.value?.answerStatus === AnswerStatus.PENDING ||
//   curAnswer.value?.answerStatus === AnswerStatus.SUCCESS &&
//   curAnswer.value.answer !== showAnswer.value);

// 根据curAnswer的新旧值的变化，来将差异部分的内平滑的一个一个字符的显示出来，而不是整个一起显示出来
// let charsArray: string[] = []
// const showAnswer = ref('');
// let charTimer: number = 0;
// const startShowChar = () => {
//   if (charTimer) {
//     window.clearTimeout(charTimer);
//   }
//   const isScrollToBottomAuto = curIdx.value > 0;
//   if (charsArray.length) {
//     let inteval = 80;
//     if (charsArray.length > 30) {
//       inteval = 30
//     } else if (charsArray.length > 10) {
//       inteval = 50
//     }
//     showAnswer.value += charsArray.shift() || '';
//     charTimer = window.setTimeout(startShowChar, inteval);

//     if (charsArray.length) {
//       emit('polling-render-answer', '', isScrollToBottomAuto);
//     } else {
//       emit('polling-render-answer', props.questionId + '-' + curIdx.value, isScrollToBottomAuto);
//     }
//   }
// }
// watch(() => curAnswer.value, (newValue, oldValue) => {
//   console.log('watchCurAnswer', newValue.id, oldValue?.id)
//   if (
//     newValue?.answerStatus === AnswerStatus.PENDING ||
//     newValue?.answerStatus === AnswerStatus.SUCCESS && props.currentPollingAnswers.includes(props.questionId + '-' + curIdx.value)) {
//     const isSameAnswer = oldValue?.id === newValue.id;
//     if (!isSameAnswer) {
//       showAnswer.value = '';
//       charsArray = [];
//       if (charTimer) {
//         window.clearTimeout(charTimer);
//       }
//     }
//     const oldAnswer = isSameAnswer ? oldValue?.answer || '' : '';
//     const newAnswer = newValue.answer || '';
//     charsArray = charsArray.concat(newAnswer.slice(oldAnswer.length).split(''));
//     console.log('charsArray', charsArray, newAnswer, oldAnswer)
//     startShowChar();
//   } else if (newValue?.answerStatus === AnswerStatus.SUCCESS) {
//     if (charTimer) {
//       window.clearTimeout(charTimer);
//     }
//     showAnswer.value = newValue.answer || '';
//   }
// }, {
//   immediate: true
// });

const emit = defineEmits<{
  (event: 'retry', questionId: string, curIdx: number): void
  (event: 'followup', curAnswer: ChatMessage): void
  (event: 'change-answer', questionId: string, cb: (idx: number) => void): void
  // (event: 'polling-render-answer', key: string, isScrollToBottomAuto: boolean): void
  (
    e: 'optimaze-answer',
    questionId: string,
    curAnswer: ChatMessage,
    optimizeId: string
  ): void
}>()

const onRetry = () => {
  emit('retry', props.questionId, curIdx.value)
}

const onFollowup = () => {
  emit('followup', curAnswer.value)
}

const onChangeAnswer = () => {
  if (props.answers.length >= props.changeAnswerLimit) {
    message.warn(`仅支持回答${props.changeAnswerLimit}个答案`)
    return
  }
  emit('change-answer', props.questionId, (idx) => {
    curIdx.value = idx
  })
}

const switchAnswer = () => {
  curIdx.value = (curIdx.value + 1) % props.answers.length
}

const onCopy = () => {
  copyToPaste(curAnswer.value.answer!)
  reportClick(ElementClick.ai_assist_copy)
}
</script>
<style lang="less" scoped>
.switch-answer {
  cursor: pointer;
  padding: 4px 0;

  .page {
    color: var(--site-theme-text-tertiary);
    margin-right: 6px;
  }
}

.text {
  padding: 16px 12px;
  background-color: var(--site-theme-bg-secondary);
  color: var(--site-theme-text-primary);
}

.followup {
  margin-top: 10px;

  a {
    color: var(--site-theme-brand);

    .icon-reply {
      margin-right: 8px;
    }
  }
}

.btns {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 16px;
}

.icons {
  & > * {
    padding: 4px 8px;
    color: var(--site-theme-text-primary);

    &:hover {
      background-color: var(--site-theme-background-hover);
    }
  }
}
</style>
<style lang="less">
@keyframes blink {
  to {
    visibility: hidden;
  }
}

.copilot-blinking {
  content: '';
  width: 8px;
  height: 16px;
  background-color: var(--site-theme-brand);
  display: inline-block;
  animation: blink 1s steps(5, start) infinite;
  vertical-align: sub;
  margin: 0 4px;
}

.idea-markdown-view-container.blinking {
  p:last-child::after {
    .copilot-blinking();
  }
}
.idea-markdown-view-container.blinking.cursor {
  .copilot-blinking();
}
</style>
