<template>
  <a-drawer
    v-model:visible="visible"
    placement="right"
    :mask="false"
    width="100%"
    :closable="false"
    :style="{
      position: 'absolute',
      height: 'calc(100% - 40px)',
      top: '40px',
    }"
    :get-container="false"
    :push="false"
  >
    <div class="answer-container">
      <div
        class="back"
        @click="handleBack"
      >
        <left-outlined class="iconxiangzuo" />
        返回
      </div>
      <PerfectScrollbar>
        <div class="detail-container">
          <QuestionItem
            v-if="question"
            :question="question"
            :is-allowed-delete="isAllowedDelete"
            @delete="handleBack"
            @questionClick="handleQuestionClick"
          />
          <List
            ref="listRef"
            :question="question"
            :is-allowed-delete="isAllowedDelete"
            @answer-click="handleAnswerClick"
          />
        </div>
      </PerfectScrollbar>
      <div class="create-answer-wrap">
        <div
          class="create-answer-btn"
          @click="handleCreate"
        >
          输入问题回复
        </div>
      </div>
      <TextareaBox
        v-model:visible="createAnswerVisible"
        :placeholder="placeholder"
        :publish-fn="handlePublishAnswer"
        type="answer"
        @success="listRef?.refresh()"
      />
    </div>
  </a-drawer>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import QuestionItem from './QuestionItem.vue';
import {
  PaperAnswer,
  PaperQuestion,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/paper/QuestionAnswer';
import { Nullable } from '~/src/typings/global';
import TextareaBox from '../TextareaBox.vue';
import List from './List.vue';
import useAllowDelete from '~/src/hooks/QuestionAndAnswer/useAllowDelete';
import { addAnswer } from '~/src/api/question';
import { LeftOutlined } from '@ant-design/icons-vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  question: Nullable<PaperQuestion>;
}>();

const { t } = useI18n();

const visible = defineModel('visible', { default: false });

const { isAllowedDelete } = useAllowDelete();

const handleBack = () => {
  visible.value = false;
};

const listRef = ref();

const replyItem = ref<PaperAnswer | null>(null);

const placeholder = ref('输入问题回复');

const handleAnswerClick = ({
  item,
  answerName,
}: {
  item: PaperAnswer;
  answerName: string;
}) => {
  replyItem.value = item;
  placeholder.value = `${t('teams.reply')}  ${answerName}`;
  createAnswerVisible.value = true;
};

const handleQuestionClick = () => {
  handleCreate();
};

const createAnswerVisible = ref(false);

const handleCreate = () => {
  placeholder.value = '输入问题回复';
  replyItem.value = null;
  createAnswerVisible.value = true;
};

const handlePublishAnswer = async (content: string) => {
  if (!props.question) {
    return;
  }
  await addAnswer({
    questionId: props.question.questionId,
    content,
    replyAnswerId: replyItem.value?.answerId,
  });
};
</script>

<style lang="less" scoped>
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

.back {
  display: flex;
  align-items: center;
  flex: 0 0 auto;
  padding: 10px;
  cursor: pointer;
  color: #a8adb3;
  .iconxiangzuo {
    font-size: 12px;
    margin-right: 4px;
  }

  font-size: 14px;
  font-weight: 400;
  position: absolute;
  top: 0;
  width: 100%;
  z-index: 100;
  background: #383a3d;
}
.answer-container {
  height: calc(100% - 44px);
  margin-top: 44px;
  .ps {
    height: 100%;
  }

  .create-answer-wrap {
    background-color: #383a3d;
    padding: 16px 10px 0px;
    height: 58px;
    position: absolute;
    bottom: 0;
    width: 100%;
    .create-answer-btn {
      background: #ffffff;
      border-radius: 2px;
      border: 1px solid rgba(0, 0, 0, 0.15);
      padding: 5px 12px;
      font-size: 14px;
      font-weight: 400;
      color: rgba(0, 0, 0, 0.25);
      cursor: pointer;
      width: 100%;
    }
  }
}
.detail-container {
  // padding: 10px;
  min-height: 100%;
  background: #383a3d;
  display: flex;
  flex-direction: column;
  position: absolute;
  left: 0;
  left: 0;
  width: 100%;

  .list {
    margin-top: 10px;
    flex: 1 1 auto;
    overflow: hidden;

    .ps {
      height: 100%;
    }
  }
}

.initial-container {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  height: 100%;
  position: relative;

  .icon {
    width: 64px;
    height: 64px;
  }

  .detail {
    font-size: 13px;
    font-weight: 400;
    color: #929599;
    margin-top: 10px;
  }
}

.textbox {
  position: absolute;
  bottom: 0;
}
</style>
