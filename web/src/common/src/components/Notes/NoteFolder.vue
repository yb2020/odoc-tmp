<template>
  <div class="note-folder pl-1 pt-3 w-full h-full">
    <a-spin :spinning="noteFolderList.length <= 1 && noteFolderLoading">
      <a-tree
        v-model:expandedKeys="emptyFolderExpanded"
        :tree-data="noteFolderList"
        :selectedKeys="[preSelectedKey]"
        @dblclick="handleDblclick"
        @click="handleClick"
        class="list-tree"
      >
        <template #title="{ title, key, count, data }">
          <div
            class="dis-flex"
            :class="{
              active:
                key === preSelectedKey && (isInClient || userStore.isLogin()),
              folder: !data.isDoc && key !== noteAllFolder.key,
            }"
            :data-key="key"
          >
            <div class="drop-middle" :class="{'selected': key === preSelectedKey}">
              <img
                v-if="key === noteAllFolder.key || data.isDoc"
                class="aiknowledge-icon"
                style="width: 14px; height: 14px; margin-top: 8px;"
                :src="getIcon(key)"
                alt=""
              >
              <component
                :is="
                  noteFolderExpanded.includes(key) ||
                    emptyFolderExpanded.includes(key)
                    ? FolderOpenFilled
                    : FolderFilled
                "
                v-else
                class="aiknowledge-icon"
                :style="{color: key === preSelectedKey ? 'var(--site-theme-text-inverse, #ffffff)' : 'var(--site-theme-text-tertiary, #5c6b7a)'}"
                @click.stop="toggleEmptyFolder(key)"
              />
              <span class="folder-selected">
                <span class="title-selected">{{
                  key === noteAllFolder.key
                    ? `${$t(`common.text.all`)}${$t(
                      `common.notes.${NoteSubType2I18nKey[activeType]}`,
                      count
                    )}`
                    : title
                }}</span>
                <!-- <span class="title-selected">{{
                  key === noteAllFolder.key
                    ? `${$t(`common.text.all`)}${$t(
                      `common.notes.${NoteSubType2I18nKey[activeType]}`,
                      activeType === NoteSubTypes.Vocabulary ? data.noteWordCount : activeType === NoteSubTypes.Annotation ? data.noteAnnotateCount : count
                    )}`
                    : title
                }}</span> -->
                <span v-if="!data.isDoc" class="count-selected">{{
                  key === noteAllFolder.key
                    ? activeType === NoteSubTypes.Vocabulary
                      ? noteAllFolder.noteWordCount
                      : activeType === NoteSubTypes.Annotation
                        ? noteAllFolder.noteAnnotateCount
                        : noteAllFolder.count
                    : activeType === NoteSubTypes.Vocabulary
                      ? data.noteWordCount
                      : activeType === NoteSubTypes.Annotation
                        ? data.noteAnnotateCount
                        : count
                }}</span>
              </span>
            </div>
          </div>
        </template>
      </a-tree>
    </a-spin>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import { FolderFilled, FolderOpenFilled } from '@ant-design/icons-vue';

import {
  NoteSubTypes,
  NoteFolder,
  NoteSubType2Icons,
  NoteSubType2I18nKey,
} from './types';
import { useNote } from './useNote';
import { useUserStore } from '../../stores/user';

const props = defineProps<{
  isInClient?: boolean;
  activeType: NoteSubTypes;
  noteState: ReturnType<typeof useNote>;
}>();

const userStore = useUserStore();

const {
  noteAllFolder,
  noteFolderList,
  noteFolderLoading,
  emptyFolderExpanded,
  noteFolderExpanded,
  noteFolderSelected,
  clickNoteFolder,
} = props.noteState;

const preSelectedKey = ref<string>(noteAllFolder.key);

const getIcon = (key: string) => {
  return NoteSubType2Icons[props.activeType][
    preSelectedKey.value === key && userStore.isLogin() ? 1 : 0
  ];
};

const handleDblclick = (event: MouseEvent, node: { key: string }) => {
  // 双击选中并加载数据
  preSelectedKey.value = node.key;
  clickNoteFolder([node.key]);
};
const handleClick = (event: MouseEvent, node: { key: string }) => {
  // 单击选中不加载数据
  preSelectedKey.value = node.key;
};
const toggleEmptyFolder = (key: NoteFolder['key']) => {
  if (emptyFolderExpanded.value.includes(key)) {
    emptyFolderExpanded.value = emptyFolderExpanded.value.filter(
      (exp: string) => exp !== key
    );
  } else {
    emptyFolderExpanded.value.push(key);
  }
};

watch(
  noteFolderSelected,
  (value) => {
    preSelectedKey.value = value;
  },
  { immediate: true }
);
</script>

<style lang="less" scoped>
@import url('~common/assets/functions.less');

.note-folder {
  /*默认显示滚动条，设置和背景色同色看不出来*/
  overflow: scroll;
  border-right: 1px solid var(--site-theme-border-color, #e9ebf0);
  
  .psm();
  
  :deep(.ant-tree) {
    height: 100%;
    background: transparent;
    color: var(--site-theme-text-color, #213547);
    
    .ant-tree-treenode,
    .ant-tree-title {
      width: 100%;
    }
    
    .ant-tree-switcher {
      display: none; /* 隐藏树形目录前的小三角 */
    }
    
    .ant-tree-node-content-wrapper {
      flex: 1;
      width: 100%;
      height: 32px;
      line-height: 32px;
      padding: 0;
      
      &.ant-tree-node-selected {
        background-color: transparent;
      }
    }
  }
}

.list-tree {
  position: relative;
  
  .dis-flex {
    display: flex;
    position: relative;
    margin-top: -4px;
    margin-bottom: -4px;
    
    &.active {
      .drop-middle.selected {
        background-color: var(--site-theme-primary-color);
        border-radius: 2px;
        
        .title-selected,
        .count-selected {
          color: var(--site-theme-text-inverse);
        }
      }
    }
  }
  
  .drop-middle {
    display: flex;
    align-items: stretch;
    position: relative;
    margin-top: 4px;
    margin-bottom: 4px;
    user-select: none;
    padding: 0 5px;
    border-radius: 2px;
    width: 215px;
    
    &:hover {
      background: var(--site-theme-background-hover);
      border-radius: 2px;
    }
    
    &.selected {
      background-color: var(--site-theme-primary-color);
    }
  }
  
  .aiknowledge-icon {
    width: 28px;
    height: 32px;
    line-height: 32px;
    display: inline-block;
    text-align: center;
    vertical-align: top;
    color: var(--site-theme-text-tertiary);
    
    &:hover {
      color: var(--site-theme-primary-color);
    }
  }
  
  .folder-selected {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 32px;
    flex: 1;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    color: var(--site-theme-text-color);
  }
  
  .title-selected {
    flex: 1;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    padding-left: 5px;
    color: var(--site-theme-text-color);
  }
  
  .count-selected {
    margin-left: 8px;
    margin-right: 5px;
    color: var(--site-theme-text-tertiary);
  }
}
</style>
