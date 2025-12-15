<template>
  <a-spin
    :spinning="fetchState.pending"
    :wrapperClassName="answers?.length === 0 ? 'center' : ''"
  >
    <div class="answer-list">
      <ErrorVue
        v-if="fetchState.error"
        :message="fetchState.error.message"
      >
        <a-button
          type="text"
          @click="fetch"
        >
          重试
        </a-button>
      </ErrorVue>
      <div
        v-else-if="answers?.length"
        class="list"
      >
        <AnswerItem
          v-for="item in answers"
          :key="item.answerId"
          :item="item"
          :isAllowedDelete="isAllowedDelete"
          @itemClick="handleItemClick"
          @delete="fetch"
        />
      </div>
      <ErrorVue
        v-else-if="answers?.length === 0"
        :img-url="emptyImgUrl"
        message="暂时没有人回复问题"
      />
    </div>
  </a-spin>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import {
  getQuestionDetail,
} from '~/src/api/question';
import emptyImgUrl from '@/assets/images/initial.png';
import AnswerItem from './AnswerItem.vue';
import { PaperAnswer, PaperQuestion } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/paper/QuestionAnswer';
import { Nullable } from '~/src/typings/global';
import ErrorVue from '../Error.vue';
import useFetch from '~/src/hooks/useFetch';

const props = defineProps<{
  question: Nullable<PaperQuestion>,
  isAllowedDelete: boolean,
}>()

const emit = defineEmits<{
  (event: 'answerClick', info: { item: PaperAnswer, answerName: string }): void
}>()

const answers = ref<Nullable<PaperAnswer[]>>(null);

const { fetch, fetchState } = useFetch(async () => {
  if (!props.question?.questionId) {
    return;
  }
  answers.value = null
  const res = await getQuestionDetail({
    questionId: props.question.questionId,
  });
  answers.value = (res && res.answers) || [];
}, false)

watch(() => props.question, () => {
  fetch()
})


const handleItemClick = (info: { item: PaperAnswer, answerName: string }) => {
  emit('answerClick', info)
};


defineExpose({
  refresh: fetch,
})



</script>

<style lang="less" scoped>
.answer-list {
  min-height: 160px;
  flex: 1;
  display: flex;
  width: 100%;
  flex-direction: column;
  padding-bottom: 50px;
}
</style>
