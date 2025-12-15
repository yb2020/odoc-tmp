<template>
  <div
    class="cursor-pointer relative question-item"
    @click="handleSuggestionClick()"
  >
    <a-tooltip
      :title="suggestion.question"
      placement="left"
    >
      <div class="question-item-text">
        {{ suggestion.question }}
      </div>
    </a-tooltip>
    <TippyVue
      v-if="deleteTippyTriggerRef"
      ref="tippyVueRef"
      :append-to-parent="true"
      :trigger-ele="deleteTippyTriggerRef"
      :placement="'bottom-end'"
      :trigger="'click'"
      :disable-draggable="true"
      @on-hide="onVisibleChange(false)"
    >
      <div
        class="text-rp-neutral-10 bg-white px-4 py-4"
        @click.stop
      >
        <div class="flex items-center space-x-1">
          <ExclamationCircleFilled :style="{ color: '#faad14' }" /><span>{{
            $t('aiCopilot.deleteTip')
          }}</span>
        </div>
        <div class="flex justify-end space-x-2 mt-3">
          <a-button
            size="small"
            @click="close"
          >
            {{
              $t('viewer.cancel')
            }}
          </a-button>
          <a-button
            size="small"
            type="primary"
            @click="deleteQuestion(suggestion.id)"
          >
            {{ $t('viewer.delete') }}
          </a-button>
        </div>
      </div>
    </TippyVue>
    <div
      ref="deleteTippyTriggerRef"
      :class="[
        'absolute right-0 top-0 w-4 h-4 delete',
        { show: deleteQuestionId === suggestion.id },
      ]"
      :style="{ transform: 'translate(50%, -50%)' }"
      @click.stop="deleteQuestionId = suggestion.id"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="16"
        height="16"
        viewBox="0 0 16 16"
        fill="none"
      >
        <g opacity="0.85">
          <path
            d="M8 1C4.13438 1 1 4.13438 1 8C1 11.8656 4.13438 15 8 15C11.8656 15 15 11.8656 15 8C15 4.13438 11.8656 1 8 1ZM10.5844 10.6594L9.55313 10.6547L8 8.80313L6.44844 10.6531L5.41563 10.6578C5.34688 10.6578 5.29063 10.6031 5.29063 10.5328C5.29063 10.5031 5.30156 10.475 5.32031 10.4516L7.35313 8.02969L5.32031 5.60938C5.30143 5.58647 5.29096 5.5578 5.29063 5.52812C5.29063 5.45937 5.34688 5.40312 5.41563 5.40312L6.44844 5.40781L8 7.25938L9.55156 5.40938L10.5828 5.40469C10.6516 5.40469 10.7078 5.45937 10.7078 5.52969C10.7078 5.55937 10.6969 5.5875 10.6781 5.61094L8.64844 8.03125L10.6797 10.4531C10.6984 10.4766 10.7094 10.5047 10.7094 10.5344C10.7094 10.6031 10.6531 10.6594 10.5844 10.6594Z"
            fill="#86919C"
          />
        </g>
      </svg>
    </div>
  </div>
</template>
<script setup lang="ts">
import { QuestionInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiAssistantReading/CustomQuestion';
import { ref } from 'vue';
import { countCustomQuestion, deleteCustomQuestion } from '~/src/api/copilot';
import TippyVue from '@/components/Tippy/index.vue';
import { ExclamationCircleFilled } from '@ant-design/icons-vue';

const deleteTippyTriggerRef = ref<HTMLDivElement>();

const tippyVueRef = ref<InstanceType<typeof TippyVue>>();

const props = defineProps<{
  suggestion: QuestionInfo;
  onConfirm: () => Promise<void>;
}>();

const deleteQuestionId = ref('');
const deleteQuestion = async (id: string) => {
  await deleteCustomQuestion({
    questionId: id,
  });
  await props.onConfirm();
  deleteQuestionId.value = '';
  tippyVueRef.value?.hide();
};

const close = () => {
  tippyVueRef.value?.hide();
};

const emit = defineEmits<{
  (event: 'suggestion:fill', question: QuestionInfo): void;
}>();
const handleSuggestionClick = () => {
  emit('suggestion:fill', props.suggestion);
  // 恢复计数API调用，记录问题使用次数
  countCustomQuestion({
    questionId: props.suggestion.id,
  });
};

const onVisibleChange = (visible: boolean) => {
  if (!visible) {
    deleteQuestionId.value = '';
  }
};
</script>
<style scoped lang="less">
.question-item {
  &-text {
    color: #4e5969;
    max-width: 100%;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    padding: 5px 0;
  }
  .delete {
    display: none;
    &.show {
      display: block;
    }
  }
  &:hover {
    background-color: #f7f8fa;
    .delete {
      display: block;
    }
  }
}
</style>
