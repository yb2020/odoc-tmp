<template>
  <div class="word-wrap h-full">
    <NoteBreadcrumb
      :type="NoteSubTypes.Vocabulary"
      :noteState="noteState"
      class="px-5 py-1 w-full bg-[#faf9f8]"
    />
    <div class="h-11 leading-10 w-full pr-4 mb-2">
      <a-checkbox
        v-model:checked="isRecite"
        class="float-right"
        @change="changeCheckbox"
      >
        {{ $t('common.notes.wordings.vmode') }}
      </a-checkbox>
    </div>
    <div
      class="flex justify-center mx-4 word-content-wrap overflow-y-auto overflow-x-hidden py-2 pr-1"
      :style="{ maxHeight: 'calc(100% - 140px)' }"
    >
      <!-- 留出10px误差空间 -->
      <a-row
        :gutter="[{ xs: 16, sm: 16, md: 16, lg: 20 }, 20]"
        :style="{ maxWidth: column * 294 + 10 + 'px' }"
        class="flex-1"
      >
        <a-col
          v-for="data in wordList"
          :key="data.id"
          flex="0 0 294px"
        >
          <WordCard
            :data="data"
            isLimitedHeight
            :isHiddenContent="isRecite"

            :titleStyle="{
              color: 'var(--site-theme-text-inverse)',
              backgroundColor: 'var(--site-theme-primary-color)',
            }"
            :contentStyle="{
              color: 'var(--site-theme-text)',
              height: '100%',
            }"
            @deleted="refreshNoteExplorer"
            @mouseenter="handleHover(data)"
          />
        </a-col>
      </a-row>
    </div>
    <div
      v-if="wordTotal > 0"
      class="pagination absolute bottom-4 right-6"
    >
      <a-pagination
        :page-size-options="pageSizeOptions"
        :current="wordPageNumber"
        :page-size="wordPageSize"
        size="small"
        show-size-changer
        show-quick-jumper
        :total="wordTotal"
        @change="onNoteListChange"
        @showSizeChange="onNoteListChange"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted, ref } from 'vue';
import NoteBreadcrumb from './NoteBreadcrumb.vue';
import WordCard from './WordCard.vue';
import { useNote } from '../useNote';
import { reportElementClick } from '@common/utils/report';
import { WordInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/NoteModule';
import lodash from 'lodash';
import { NoteSubTypes } from '../types';

const props = defineProps<{
  noteState: ReturnType<typeof useNote>;
}>();

const pageSizeOptions = ref<string[]>(['20', '40', '60', '80']);
const isRecite = ref<boolean>(false);
const column = ref<number>(3);

const {
  fetchWordList,
  wordList,
  wordTotal,
  wordPageNumber,
  wordPageSize,
  refreshNoteExplorer,
} = props.noteState;

const onNoteListChange = (
  pageNumber: number,
  pageSize = wordPageSize.value
) => {
  wordPageSize.value = pageSize;
  fetchWordList(pageNumber);
};

const changeCheckbox = () => {
  reportElementClick({
    page_type: 'note_tab',
    element_name: 'note_word_recite',
    status: isRecite.value ? 'on' : 'off',
  });
};
const handleHover = (data: WordInfo) => {
  if (!isRecite.value) {
    // 背诵模式下上报
    return;
  }

  reportElementClick({
    page_type: 'note_tab',
    type_parameter: data.id || '',
    element_name: 'note_word_recite_card',
  });
};

const countColumnNum = lodash.debounce(() => {
  const ele = document.querySelectorAll('.word-content-wrap');
  column.value = Math.floor(ele?.[0]?.clientWidth / 294);
}, 300);

onMounted(() => {
  countColumnNum();
  window.addEventListener('resize', countColumnNum, false);
});
onUnmounted(() => {
  window.removeEventListener('resize', countColumnNum, false);
});
</script>

<style lang="less" scoped>
.word-wrap {
  background: #f0f2f5;
  :deep(.word-card) {
    height: 176px;
    display: flex;
    flex-direction: column;
    .btn-delete {
      top: 6px;
    }
    .word-ct {
      flex: 1;
      overflow: hidden;
      background: #fff;
      border-radius: 0 0 6px 6px;
      .word-note {
        color: #4e5969;
      }
    }
    .word-tt {
      line-height: 22px;
      border-radius: 6px 6px 0 0;
      font-size: 14px;
      font-weight: 600;
      width: 272px;
      word-break: break-word;
      padding: 9px 32px 9px 16px;
    }
    .phonetic {
      flex-wrap: wrap;
      gap: 2px 8px;
      margin-bottom: 0;
      min-height: 36px;
      align-items: flex-start;
    }
    .pronunciation-item {
      border: 1px solid #e5e6eb;
    }
    border-radius: 6px;
    border: 1px solid transparent;
    &:hover {
      border: 1px solid #1f71e0;
    }
  }
  ::-webkit-scrollbar {
    width: 6px; //修改垂直滚动条宽度
    height: 6px; //修改水平滚动条宽度
  }

  ::-webkit-scrollbar-thumb {
    border-radius: 10px;
    background: transparent;
  }

  &:hover {
    ::-webkit-scrollbar-thumb {
      background: rgba(201, 205, 212, 1);
      &:hover {
        background-color: rgba(170, 170, 170, 0.8);
      }
    }
  }
}
</style>

<style lang="less">
.word-wrap {
  .pagination {
    .ant-pagination-item-active {
      background: #fff !important;
      a {
        line-height: 24px;
        color: #1f71e0 !important;
        font-weight: bold;
      }
    }
    .ant-pagination-item-active:hover a {
      color: #000 !important;
    }
  }
}
</style>
