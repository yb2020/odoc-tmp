<template>
  <div class="scroll-container">
    <ErrorVue
      v-if="!isOpenPaper"
      :style="{ height: '100%' }"
      :img-url="noPaperImgUrl"
      message="未入库论文不支持问答"
    />
    <Question
      v-else
      :paper-id="paperId"
      :answer-visible="answerVisible"
      @answer-question="handleQuestionClick"
    />
    <Answer
      v-if="isOpenPaper"
      v-model:visible="answerVisible"
      :question="answerQuestion"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import Answer from './Answer/index.vue';
import { PaperQuestion } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/paper/QuestionAnswer';
import { Nullable } from '~/src/typings/global';
import Question from './Question/index.vue';
import ErrorVue from './Error.vue';
import noPaperImgUrl from '@/assets/images/no_paper.png';
import { checkOpenPaper } from '~/src/api/helper';

const props = defineProps<{ paperId: string; isPrivatePaper: boolean }>();

const answerVisible = ref(false);
const answerQuestion = ref<Nullable<PaperQuestion>>(null);

const isOpenPaper = computed(() =>
  checkOpenPaper(props.paperId, props.isPrivatePaper)
);

const handleQuestionClick = (question: PaperQuestion) => {
  answerQuestion.value = question;
  answerVisible.value = true;
};
</script>

<style lang="less" scoped>
.scroll-container {
  height: 100%;
  overflow: hidden;
}
</style>
