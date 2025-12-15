<template>
  <div class="question-container">
    <PerfectScrollbar>
      <div class="question-inner">
        <List
          ref="listRef"
          v-model:create-textarea-visible="createQuestionVisible"
          :paper-id="paperId"
          @questionClick="handleQuestionClick"
        />
      </div>
    </PerfectScrollbar>
    <TextareaBox
      v-model:visible="createQuestionVisible"
      placeholder="创建问题"
      :publish-fn="handlePublishQuestion"
      type="question"
      @success="listRef?.refresh()"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import List from './List.vue';
import TextareaBox from '../TextareaBox.vue';
import { PaperQuestion } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/paper/QuestionAnswer';
import { addQuestion } from '~/src/api/question';
import { currentNoteInfo } from '@/store';

const props = defineProps<{
  answerVisible: boolean;
  paperId: string;
}>();
const emit = defineEmits<{
  (event: 'answerQuestion', question: PaperQuestion): void;
}>();

const createQuestionVisible = ref(false);

const listRef = ref();

const handlePublishQuestion = async (content: string) => {
  await addQuestion({
    content,
    paperId: props.paperId,
    pdfId: currentNoteInfo.value?.pdfId || '',
  });
};

const handleQuestionClick = (question: PaperQuestion) => {
  emit('answerQuestion', question);
};

watch(
  () => props.answerVisible,
  (newVal) => {
    if (!newVal) {
      listRef.value.refresh();
    }
  }
);
</script>

<style lang="less" scoped>
.question-container {
  height: 100%;
  .ps {
    height: 100%;
  }

  .question-inner {
    position: relative;
    min-height: 100%;
    display: flex;
    flex-direction: column;
  }
}
</style>
