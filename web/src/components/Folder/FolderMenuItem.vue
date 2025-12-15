<script lang="ts" setup>
import { Menu } from 'ant-design-vue';
import { CheckOutlined, PlusOutlined } from '@ant-design/icons-vue';
import FolderMenuItem from './FolderMenuItem.vue';
import { PaperFolder, useFolderStore, limit30 } from '~/src/stores/folderStore';

const folderStore = useFolderStore();

const props = defineProps<{
  folderInfo: PaperFolder;
  selectedFolderMap: Record<string, boolean>;
  detectDisable(key: string): boolean;
}>();

const emit = defineEmits<{
  (event: 'clickItem', key: string): void;
}>();

const { Item, SubMenu } = Menu;

const clickItem = (key: string) => {
  emit('clickItem', key);
};

const keydown = (event: KeyboardEvent) => {
  if (event.key === 'Escape') {
    folderStore.cancelFolder();
  }
};
</script>

<template>
  <SubMenu
    :key="folderInfo.key"
    theme="light"
    popup-class-name="metadata-folder-submenu"
    :class="{
      'metadata-folder-selected':
        folderInfo.key in selectedFolderMap,
    }"
    v-bind="props"
    :mode="'vertical'"
    class="thin-scroll"
  >
    <template #title>
      <span
        class="metadata-folder-name"
        :class="{
          'metadata-folder-name-selected':
            folderInfo.key in selectedFolderMap,
        }"
        :disabled="detectDisable(folderInfo.key)"
        @click.stop="clickItem(folderInfo.key)"
      >{{ limit30(folderInfo.title) }}
        <CheckOutlined />
      </span>
    </template>
    <FolderMenuItem
      v-for="item in folderInfo.children"
      :key="item.key"
      :folder-info="item"
      :selected-folder-map="selectedFolderMap"
      :disabled="detectDisable(item.key)"
      :detect-disable="detectDisable"
      @click-item="clickItem($event)"
    />
    <Item v-if="folderStore.newFolderParent === folderInfo">
      <input
        v-model="folderStore.newFolderName"
        :data-folder-key="folderStore.newFolderParent.key"
        class="metadata-folder-input blue-6-b"
        @keypress.stop.enter="folderStore.submitFolder()"
        @keypress.stop.esc="folderStore.cancelFolder()"
        @keydown.stop="keydown"
        @keyup.stop
        @click.stop
      >
    </Item>
    <Item v-else>
      <div
        class="metadata-folder-new blue-6-c"
        @click.stop="folderStore.newFolder(folderInfo)"
      >
        <PlusOutlined class="blue-6-c" />{{ $t('folder.new_folder') }}
      </div>
    </Item>
  </SubMenu>
</template>

<style lang="less">
// @import '~/src/assets/less/style.less';
@import './folder-style.less';

.metadata-folder-name {
  min-width: 120px;
  width: 100%;
  display: inline-flex;
  align-items: center;
  justify-content: space-between;
  > span {
    margin-right: 16px;
    margin-left: 8px;
    visibility: hidden;
  }

  &.metadata-folder-name-selected {
    color: #1f71e0;
    > span {
      visibility: visible !important;
    }
  }
}

.metadata-folder-selected {
  background-color: #e8f5ff;
}
</style>
