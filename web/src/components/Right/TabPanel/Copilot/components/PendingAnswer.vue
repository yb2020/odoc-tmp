<script setup lang="ts">
import Typing from './Typing.vue';
import { computed, onMounted, ref, watch } from 'vue';
import MarkdownViewer from './MarkdownText.vue';
import { formatQuestion } from './tool';

interface IAnswer {
  answer: string;
  anchors?: string[];
  status: 'success' | 'pending';
}

const props = defineProps<{
  answer: IAnswer;
  scrollToBottom: () => void;
}>();

const curAnswer = computed<IAnswer>(() => {
  return props.answer;
});

const firstStatus = ref(curAnswer.value.status);

const answerWithAnchor = computed(() => {
  if (
    curAnswer.value.status !== 'success' ||
    !curAnswer.value.anchors?.length
  ) {
    return curAnswer.value.answer;
  }
  return formatQuestion(curAnswer.value.answer, curAnswer.value.anchors);
});

const emit = defineEmits<{
  (event: 'quickAskQuestion', questionContent: string): void;
  (event: 'textShowComplete', isNoScrollBottom: boolean): void;
}>();
const handleQuestion = (e: MouseEvent) => {
  const target = e.target as HTMLElement;
  if (target.classList.contains('js-copilot-question')) {
    const question = target.innerText;
    if (question) {
      // 将question里面的序号去掉
      emit('quickAskQuestion', question.replace(/^[0-9]+[\.\)）]\s?/, ''));
    }
  }
};

onMounted(() => {
  if (firstStatus.value === 'success') {
    emit('textShowComplete', true);
  }
});

defineExpose({
  reTyping: () => {
    firstStatus.value = 'pending';
  },
});
</script>
<template>
  <Typing
    v-if="firstStatus === 'pending'"
    class="rp-markdown-view-container"
    :anchors="curAnswer.anchors"
    :text="curAnswer.answer || ''"
    :is-pending="curAnswer.status === 'pending'"
    :scroll-to-bottom="scrollToBottom"
    @quick-ask-question="handleQuestion"
    @textShowComplete="emit('textShowComplete', false)"
  />
  <MarkdownViewer
    v-else
    class="rp-markdown-view-container"
    :text="answerWithAnchor || ''"
    @click="handleQuestion"
  />
</template>
<style lang="less">
.rp-markdown-view-container {
  .copilot-question-anchor {
    color: #1a9fff;
    cursor: pointer;
    &:hover {
      color: #1f71e0;
    }
  }
}
</style>
