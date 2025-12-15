<template>
  <div
    class="note-search-container"
    :style="{
      paddingTop: '8px',
      paddingBottom: '8px',
      fontSize: '12px',
      background: 'var(--site-theme-bg-secondary)',
      marginLeft: '-17px',
      paddingLeft: '17px',
    }"
  >
    <NoteBreadcrumb
      :type="NoteSubTypes.Annotation"
      :noteState="noteState"
      class="mx-5"
    />
    <div class="note-search-keyword note-search-keyword-electron">
      <a-input
        v-model:value="noteSearchKeyword"
        :placeholder="placeholder"
        @change="search()"
        class="theme-search-input"
      >
        <template #prefix>
          <search-outlined class="search-icon" />
        </template>
      </a-input>
    </div>
    <div class="note-search-rule">
      <a-dropdown
        v-model:visible="ruleVisible"
        overlay-class-name="note-search-rule-overlay"
      >
        <a
          class="ant-dropdown-link"
          @click.prevent
        >
          {{ $t('common.notes.wordings.filterRule') }}
          <span v-if="ruleCheckedLength < ruleCheckedFullLength">({{ ruleCheckedLength }}/{{ ruleCheckedFullLength }})</span>
          <down-outlined
            style="transition: 0.24s"
            :style="{
              transform: ruleVisible ? 'rotate(-180deg)' : '',
            }"
          />
        </a>
        <template #overlay>
          <a-menu class="note-search-rule-overlay">
            <a-menu-item
              v-for="style in noteStyleList"
              :key="style.type"
              @click="clickNoteStyle(style?.type)"
            >
              <em :style="{ backgroundColor: style?.color, display: 'inline-block', width: '10px', height: '10px', borderRadius: '50%', marginRight: '12px' }" />
              <span>{{ $t(style.i18n) }}</span>
              <div>
                <check-outlined v-show="noteStyle[style?.type]" />
              </div>
            </a-menu-item>
            <a-menu-item @click="toggleNoteQuote()">
              <img
                src="@common/../assets/images/notes/quote.svg"
                alt=""
              >
              <span>{{ $t('common.notes.wordings.quoteLabel') }}</span>
              <div>
                <check-outlined v-show="noteQuote" />
              </div>
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
    <div class="note-search-sort">
      <a-dropdown
        v-model:visible="sortVisible"
        :overlay-class-name="'note-search-rule-overlay-electron'"
      >
        <a
          class="ant-dropdown-link"
          @click.prevent
        >
          {{ $t(sortTypeMap[noteSort]) }}
          <down-outlined
            style="transition: 0.24s"
            :style="{
              transform: sortVisible ? 'rotate(-180deg)' : '',
            }"
          />
        </a>
        <template #overlay>
          <a-menu :class="['note-search-rule-overlay-electron']">
            <a-menu-item @click="selectNoteSort(SortTypeEnum.NEWEST)">
              {{ $t(sortTypeMap[SortTypeEnum.NEWEST]) }}
            </a-menu-item>
            <a-menu-item @click="selectNoteSort(SortTypeEnum.EARLIEST)">
              {{ $t(sortTypeMap[SortTypeEnum.EARLIEST]) }}
            </a-menu-item>
            <a-menu-item @click="selectNoteSort(SortTypeEnum.DEFAULT)">
              {{ $t(sortTypeMap[SortTypeEnum.DEFAULT]) }}
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { SortTypeEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/protocol/note/request/_GetMyNoteMarkListReq';
import lodash from 'lodash';
import { noteStyleList, useNote } from '../useNote';
import { NoteSubTypes, ColorKey } from '../types';
import NoteBreadcrumb from './NoteBreadcrumb.vue';
import {
  SearchOutlined,
  DownOutlined,
  CheckOutlined,
} from '@ant-design/icons-vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const props = defineProps<{
  noteState: ReturnType<typeof useNote>;
}>();
const sortTypeMap = {
  [SortTypeEnum.NEWEST]: 'common.notes.wordings.sortNewest',
  [SortTypeEnum.EARLIEST]: 'common.notes.wordings.sortEarliest',
  [SortTypeEnum.DEFAULT]: 'common.notes.wordings.sortDefault',
  [SortTypeEnum.UNRECOGNIZED]: '',
};

const {
  noteAllFolder,
  noteSearchKeyword,
  noteList,
  noteLoading,
  noteSort,
  noteStyle,
  noteQuote,
  noteBreadcrumbList,
  noteFolderSelected,
  fetchNoteList,
} = props.noteState;

const ruleVisible = ref(false);
const ruleCheckedFullLength = Object.keys(noteStyle.value).length + 1;
const ruleCheckedLength = computed(
  () =>
    Object.values(noteStyle.value).filter(Boolean).length +
    Number(noteQuote.value)
);
const sortVisible = ref(false);
const placeholder = computed(() => {
  if (noteFolderSelected.value === noteAllFolder.key) {
    return t('common.notes.wordings.searchAll');
  }

  const { title } =
    noteBreadcrumbList.value[noteBreadcrumbList.value.length - 1];
  return t('common.notes.wordings.searchFolder', { title });
});

const clickNoteStyle = (style: ColorKey) => {
  noteStyle.value[style] = !noteStyle.value[style];
  fetchNoteList(1);
};
const toggleNoteQuote = () => {
  noteQuote.value = !noteQuote.value;
  fetchNoteList(1);
};
const selectNoteSort = (sort: SortTypeEnum) => {
  sortVisible.value = false;
  noteSort.value = sort;
  fetchNoteList(1);
};
const search = lodash.debounce(() => {
  noteLoading.value = true;
  noteList.value = [];
  fetchNoteList(1);
}, 300);
</script>

<style lang="less" scoped>
.note-search-container {
  padding-right: 30px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  .note-search-keyword {
    flex: 1 0 200px;
    &.note-search-keyword-electron {
      > * {
        font-size: 12px !important;
        /deep/input {
          height: 28px;
          font-size: 12px !important;
        }
      }
    }
  }
}

.note-search-rule {
  flex: 0 0 120px;
  display: flex;
  justify-content: center;
}

.note-search-rule-overlay {
  .ant-dropdown-menu-item {
    font-size: 14px;
  }
}
.note-search-rule-overlay-electron {
  .ant-dropdown-menu-item {
    font-size: 12px;
  }
}

.ant-dropdown-link {
  color: #4e5969;
}

::v-deep.note-search-rule-overlay,
.note-search-rule-overlay-electron {
  .ant-dropdown-menu-title-content {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    width: 120px;
    padding-left: 15px;
    padding-right: 12px;
    color: #4e5969;

    > em:first-of-type {
      display: inline-block;
      height: 10px;
      flex: 0 0 10px;
      border-radius: 5px;
    }

    > span:first-of-type {
      flex: 1 1 100%;
      padding-left: 8px;
    }

    > div:last-of-type {
      flex: 0 0 16px;
      height: 22px;
      color: #1f71e0;
    }
  }
}

.note-search-sort-select {
  color: #4e5969;
}

.note-search-sort {
  flex: 0 0 110px;
  display: flex;
  justify-content: center;
}
</style>
