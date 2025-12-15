<template>
  <div class="copilot-updown">
    <i
      v-if="answer.isLike"
      class="aiknowledge-icon icon-accept-checked blue"
      @click="updown('cancelLike')"
    />
    <i
      v-else
      class="aiknowledge-icon icon-accept"
      @click="updown('like')"
    />
    <a-popover
      v-if="answer.isDisLike"
      :getPopupContainer="() => scrollEle"
      :visible="feedbackVisible"
      placement="bottomLeft"
      :title="null"
      trigger="click"
      @confirm="feedback"
    >
      <template #content>
        <div ref="popoverTarget">
          <a-input
            ref="inputTarget"
            v-model:value="feedbackValue"
            autofocus
            class="copilot-updown-input"
            size="small"
            :placeholder="$t('aiCopilot.reportTip')"
          />
          <div class="copilot-updown-btns">
            <a-button
              size="small"
              class="cancel-btn"
              @click="feedbackVisible = false"
            >
              {{ $t('viewer.cancel') }}
            </a-button>
            <a-button
              size="small"
              type="primary"
              @click="feedback"
            >
              {{
                $t('viewer.confirm')
              }}
            </a-button>
          </div>
        </div>
      </template>
      <i
        class="aiknowledge-icon icon-dislike-checked blue"
        @click="updown('cancelDislike')"
      />
    </a-popover>
    <i
      v-else
      class="aiknowledge-icon icon-dislike"
      @click="updown('dislike')"
    />
  </div>
</template>
<script lang="ts" setup>
import { AnswerInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiAssistReading';
import { onClickOutside } from '@vueuse/core';
import { ref, nextTick } from 'vue';
import {
  cancelDislikeAnswer,
  cancelLikeAnswer,
  dislikeAnswer,
  dislikeFeedback,
  likeAnswer,
} from '~/src/api/copilot';

const props = defineProps<{ answer: AnswerInfo; scrollEle: HTMLDivElement }>();

const inputTarget = ref();

const updown = async (
  type: 'like' | 'cancelLike' | 'dislike' | 'cancelDislike'
) => {
  if (type === 'like') {
    await likeAnswer({
      answerId: props.answer!.id,
    });
    props.answer.isLike = true;
    props.answer.isDisLike = false;
  } else if (type === 'cancelLike') {
    await cancelLikeAnswer({
      answerId: props.answer!.id,
    });
    props.answer.isLike = false;
  } else if (type === 'dislike') {
    await dislikeAnswer({
      answerId: props.answer!.id,
    });
    props.answer.isLike = false;
    props.answer.isDisLike = true;
    feedbackVisible.value = true;
    await nextTick();
    console.log(inputTarget.value, inputTarget.value.input);
    setTimeout(() => {
      inputTarget.value.focus();
    }, 500);
  } else {
    await cancelDislikeAnswer({
      answerId: props.answer.id,
    });
    props.answer.isDisLike = false;
  }
};

const popoverTarget = ref();

onClickOutside(popoverTarget, (e) => {
  feedbackVisible.value = false;
});

const feedbackVisible = ref(false);
const feedbackValue = ref('');
const feedback = async () => {
  const value = feedbackValue.value.trim();
  if (value) {
    await dislikeFeedback({
      answerId: props.answer.id,
      feedback: value,
    });
    feedbackValue.value = '';
    feedbackVisible.value = false;
  } else {
    feedbackVisible.value = false;
  }
};
</script>
<style lang="less" scoped>
.copilot-updown {
  .aiknowledge-icon {
    font-size: 16px;
    margin-right: 24px;
    cursor: pointer;
    &.blue {
      color: #4387d9;
    }
  }

  &-input {
    width: 250px;
  }

  &-btns {
    text-align: right;
    margin-top: 8px;

    .ant-btn + .ant-btn {
      margin-left: 8px;
    }

    .cancel-btn {
      background-color: #5b6167;
      border: none;
      color: #d5d8dd;
    }
  }
}
</style>
