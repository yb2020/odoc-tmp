<template>
  <div
    ref="domRef"
    class="h-[364px]"
  >
    <div class="px-4 pt-2">
      <a-dropdown :get-popup-container="getContainer">
        <span
          class="text-rp-neutral-8 cursor-pointer"
          @click.prevent
        >
          {{
            sortKey === SortKey.ADD_TIME
              ? $t('aiCopilot.addEarlist')
              : $t('aiCopilot.useMost')
          }}
          <DownOutlined />
        </span>
        <template #overlay>
          <a-menu>
            <a-menu-item @click="sortKey = SortKey.ADD_TIME">
              <a href="javascript:;">{{ $t('aiCopilot.addEarlist') }}</a>
            </a-menu-item>
            <a-menu-item @click="sortKey = SortKey.USE_COUNT">
              <a href="javascript:;">{{ $t('aiCopilot.useMost') }}</a>
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
    <a-spin :spinning="loading">
      <div class="h-[332px] py-2 px-4 overflow-auto">
        <div
          v-if="showEmpty"
          class="flex flex-col justify-center items-center min-h-[284px]"
        >
          <img
            src="@/assets/images/right/copilot/add.png"
            class="w-12 mb-3"
          >
          <a-button
            type="link"
            @click="showAdd"
          >
            <PlusOutlined />{{ $t('viewer.add') }}
          </a-button>
        </div>
        <div
          v-else
          class="reletive"
        >
          <CustomItem
            v-for="suggestion in data"
            :key="suggestion.id"
            :suggestion="suggestion"
            :on-confirm="onDeleteConfirm"
            @suggestion:fill="handleSuggestionClick"
          />
          <div
            v-if="canAdd && !isAdding"
            :style="{ marginLeft: '-15px' }"
          >
            <a-button
              type="link"
              @click="showAdd"
            >
              <PlusOutlined />{{ $t('viewer.add') }}
            </a-button>
          </div>
        </div>
        <div v-if="isAdding">
          <a-input
            ref="inputTarget"
            v-model:value="addingValue"
            :autofocus="true"
            class="!text-rp-neutral-8"
            @pressEnter="saveQuestion"
            @keydown.esc="cancelSave"
          />
        </div>
        <div
          v-if="error && !loading"
          class="flex flex-col justify-center items-center min-h-[284px]"
        >
          <a
            class="text-rp-neutral-8 space-x-2"
            @click="run"
          ><ReloadOutlined /><span>{{ error.message }}</span></a>
        </div>
      </div>
    </a-spin>
  </div>
</template>
<script setup lang="ts">
import {
  DownOutlined,
  PlusOutlined,
  ReloadOutlined,
  CloseOutlined,
} from '@ant-design/icons-vue';
import { useLocalStorage } from '@vueuse/core';
import { useRequest } from 'ahooks-vue';
import { computed, ref } from 'vue';
import { addCustomQuestion, getCustomQuestions } from '~/src/api/copilot';
import {
  CustomQuestionSortType,
  QuestionInfo,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiAssistantReading/CustomQuestion';
import CustomItem from './CustomItem.vue';
import { SelectProps } from 'ant-design-vue';

const domRef = ref<HTMLDivElement>();

const getContainer = () => {
  return domRef.value || document.body;
};

const options = ref<SelectProps['options']>([
  { value: 'jack', label: 'Jack' },
  { value: 'lucy', label: 'Lucy' },
  { value: 'tom', label: 'Tom' },
]);

enum SortKey {
  ADD_TIME = 'ADD_TIME',
  USE_COUNT = 'USE_COUNT',
}

const sortKey = useLocalStorage<SortKey>(
  'rp-annotate2.0/copilot-custom-sort',
  SortKey.ADD_TIME
);

const { loading, data, run, error } = useRequest(
  () => {
    return getCustomQuestions({
      sortType:
        sortKey.value === SortKey.ADD_TIME
          ? CustomQuestionSortType.EARLIEST
          : CustomQuestionSortType.MOST_USE,
    });
  },
  {
    refreshDeps: [sortKey],
  }
);

const showEmpty = computed(() => {
  return data.value && data.value.length === 0 && !isAdding.value;
});

const canAdd = computed(() => {
  return data.value && data.value.length < 10;
});

const isAdding = ref(false);
const addingValue = ref('');
const addingLoading = ref(false);
const saveQuestion = async () => {
  const val = addingValue.value.trim();
  if (!val) {
    return;
  }
  if (addingLoading.value) {
    return;
  }
  addingLoading.value = true;

  try {
    await addCustomQuestion({
      question: val,
    });
    isAdding.value = false;
    addingValue.value = '';
    await run();
  } catch (error) {
  } finally {
    addingLoading.value = false;
  }
};

const cancelSave = () => {
  isAdding.value = false;
  addingValue.value = '';
};

const inputTarget = ref();
const showAdd = () => {
  isAdding.value = true;
  setTimeout(() => {
    inputTarget.value.focus();
  }, 100);
};

const emit = defineEmits<{
  (event: 'suggestion:fill', question: string): void;
}>();
const handleSuggestionClick = (question: QuestionInfo) => {
  // 兼容不同的字段名，优先使用 question.question，如果不存在则使用 question.questionContent
  const questionText = question.question
  emit('suggestion:fill', questionText);
};

const onDeleteConfirm = async () => {
  await run();
};
</script>
