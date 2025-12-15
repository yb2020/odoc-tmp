<template>
  <a-tabs
    :activeKey="activeKey"
    class="summarize-wrap h-full"
    @change="changeTab"
  >
    <a-tab-pane
      v-for="item in tabList"
      :key="item.key"
    >
      <template #tab>
        <div>
          <UnorderedListOutlined
            v-if="item.key === 'summarizeList'"
            class="ml-2 !mr-2 text-base"
          />
          <span class="tab-title -mb-1">{{
            item.getTitle?.() || item.title
          }}</span>
          <CloseOutlined
            v-if="item.closable"
            class="ml-2 !mr-0 text-xs close-icon"
            @click="closeTab(item.key)"
          />
        </div>
      </template>
      <Summarize
        :noteState="noteState"
        @clickNote="openTab"
        @loadedNote="updateTt"
      />
    </a-tab-pane>
  </a-tabs>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import { UnorderedListOutlined, CloseOutlined } from '@ant-design/icons-vue';
import { useNote } from '../useNote';
import Summarize from './Summarize.vue';
import {
  GetSummaryResponse,
  NoteManageDocInfo,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/NoteModule';
import { useI18n } from 'vue-i18n';
import type { Optional } from 'utility-types';

const props = defineProps<{
  noteState: ReturnType<typeof useNote>;
}>();

const { t } = useI18n();

const { noteFolderMap, noteFolderSelected, clickNoteFolder } = props.noteState;

const initialTab = {
  key: 'summarizeList',
  title: 'Summaries',
  getTitle: () =>
    `${t('common.notes.summary')}${t('common.symbol.space')}${t(
      'common.text.list'
    )}`,
  isFolder: true,
  closable: false,
  selectedKey: '0',
};

const activeKey = ref();
const tabList = ref<Array<Optional<typeof initialTab, 'getTitle'>>>([
  initialTab,
]);

const openTab = (e: MouseEvent, record: NoteManageDocInfo) => {
  e.preventDefault();
  clickNoteFolder([record.noteId]);
};
const changeTab = (value: string) => {
  // 切换树选中节点
  const item = tabList.value.find((item) => item.key === value);
  if (activeKey.value === item?.selectedKey) {
    return;
  }

  if (item) {
    clickNoteFolder([item.selectedKey]);
  }
};
const closeTab = (key: string) => {
  const index = tabList.value.findIndex((item) => item.key === key);
  if (key === activeKey.value) {
    // 如果删除当前展示的tab 就切换到前一个tab页
    changeTab(tabList.value[index - 1].key);
  }
  tabList.value = tabList.value.filter((item) => item.key !== key);
};
const updateTt = (data: GetSummaryResponse) => {
  document.title = data.docName;
};

watch(
  () => noteFolderSelected.value,
  (value: string) => {
    if (
      !noteFolderMap.value?.[value ?? ''] ||
      value?.split('-')?.[1] === '0' ||
      !noteFolderMap.value?.[value].isDoc
    ) {
      // 全部总结 文件夹节点
      tabList.value[0].selectedKey = noteFolderSelected.value;
      activeKey.value = 'summarizeList';
      return;
    }
    const item = tabList.value.find((item) => item.key === value);
    activeKey.value = value;
    if (item) {
      // 标签页已经存在
      return;
    }
    tabList.value.push({
      ...noteFolderMap.value[value],
      key: value,
      isFolder: false,
      closable: true,
      selectedKey: noteFolderSelected.value,
    });
  },
  { immediate: true, deep: true }
);
</script>

<style lang="less" scoped>
:deep(.summarize-wrap) {
  background-color: var(--site-theme-bg-primary);
  
  .ant-tabs-nav {
    background-color: var(--site-theme-bg-secondary);
    margin-bottom: 0;
    border-bottom: 1px solid var(--site-theme-divider);
  }
  
  .ant-tabs-tab {
    color: var(--site-theme-text-secondary);
    background-color: var(--site-theme-bg-secondary);
    border: none;
    
    &:hover {
      color: var(--site-theme-primary-color);
    }
    
    &.ant-tabs-tab-active {
      .ant-tabs-tab-btn {
        color: var(--site-theme-primary-color);
      }
    }
  }
  
  .ant-tabs-ink-bar {
    background: var(--site-theme-primary-color);
  }
  
  .close-icon {
    color: var(--site-theme-text-tertiary) !important;
    
    &:hover {
      color: var(--site-theme-text-secondary) !important;
    }
  }
  
  .ant-tabs-tab-btn {
    .tab-title {
      max-width: 152px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      display: inline-block;
    }
  }
  .ant-tabs-content {
    height: 100%;
    background-color: var(--site-theme-bg-primary);
  }
  .ant-tabs-tabpane {
    display: flex;
    flex-direction: column;
  }
}
</style>
