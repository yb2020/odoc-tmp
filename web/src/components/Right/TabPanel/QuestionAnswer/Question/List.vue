<template>
  <a-spin
    :spinning="fetchState.pending"
    :wrapperClassName="list?.length === 0 ? 'center' : ''"
  >
    <div class="question-list">
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
        v-else-if="list && list?.length > 0"
        class="list-container"
      >
        <div
          v-for="item in list"
          :key="item.questionId"
          class="item"
          @click="handleQuestionClick(item)"
        >
          <div class="title">
            {{ item.title }}
          </div>
          <div class="opreation">
            <div class="check">
              {{ item.viewCount || 0 }} 查看
            </div>
            <div class="check">
              {{ item.replyCount || 0 }} 回复
            </div>
          </div>
        </div>
        <div
          class="create-btn-default"
          @click="handleCreate"
        >
          <plus-outlined class="iconjia" />创建新问题
        </div>
      </div>
      <ErrorVue
        v-if="list?.length === 0"
        message="暂时没有人提出问题"
        :img-url="emptyImgUrl"
      >
        <a-button
          type="primary"
          class="create-btn"
          @click="handleCreate"
        >
          <plus-outlined />创建新问题
        </a-button>
      </ErrorVue>
    </div>
  </a-spin>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { getQuestionList } from '@/api/question';
import useFetch from '~/src/hooks/useFetch';
import { PaperQuestion } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/paper/QuestionAnswer';
import ErrorVue from '../Error.vue';
import { PlusOutlined } from '@ant-design/icons-vue';
import emptyImgUrl from '@/assets/images/initial.png';
import { useStore } from '@/store';

const props =
  defineProps<{ createTextareaVisible: boolean; paperId: string }>();

const emit = defineEmits<{
  (event: 'update:createTextareaVisible', val: boolean): void;
  (event: 'questionClick', question: PaperQuestion): void;
}>();

const list = ref<PaperQuestion[]>();

const store = useStore();

const { fetchState, fetch } = useFetch(async () => {
  const res = await getQuestionList({
    paperId: props.paperId,
  });
  list.value = res.questions || [];
});

const handleQuestionClick = (item: PaperQuestion) => {
  emit('questionClick', item);
};

const handleCreate = () => {
  emit('update:createTextareaVisible', true);
};

defineExpose({
  refresh: fetch,
});
</script>

<style lang="less" scoped>
.question-list {
  min-height: 160px;
  flex: 1;
  display: flex;
  width: 100%;
  flex-direction: column;
}

.center {
  flex: 1;
  align-items: center;
  display: flex;
  justify-content: center;
  flex-direction: column;

  :deep(.ant-spin-container) {
    width: 100%;
    flex: 1;
    display: flex;
    flex-direction: column;
  }
}

.list-container {
  padding: 10px;

  .item {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 4px;
    padding: 10px;
    cursor: pointer;
    margin-bottom: 10px;

    .title {
      font-size: 14px;
      font-weight: 400;
      color: #f0f2f5;

      display: -webkit-box;
      -webkit-line-clamp: 3;
      -webkit-box-orient: vertical;
      overflow: hidden;
      text-overflow: ellipsis;
      word-break: break-all;
    }

    .opreation {
      display: flex;
      align-items: center;
      justify-content: flex-end;
      margin-top: 10px;

      .check {
        font-size: 13px;
        font-weight: 400;
        color: #a8adb3;
        margin-right: 16px;

        &:last-child {
          margin-right: 0;
        }
      }
    }
  }
}

.create-btn {
  margin-top: 10px;
}

.create-btn-default {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 30px;
  cursor: pointer;
  color: #929599;
  .iconjia {
    font-size: 14px;

    margin-right: 5px;
  }
}
</style>
