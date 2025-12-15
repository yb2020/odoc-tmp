<template>
  <SubMenu
    :key="folderInfo.key"
    theme="light"
    popup-class-name="metadata-folder-submenu"
    :class="{
      'metadata-folder-selected': folderInfo.key in selectedFolderMap,
    }"
    v-bind="props"
    :mode="'vertical'"
    class="metadata-scroll"
  >
    <template #title>
      <span
        class="metadata-folder-name"
        :class="{
          'metadata-folder-name-selected': folderInfo.key in selectedFolderMap,
        }"
        :disabled="detectDisable(folderInfo.key)"
        @click.stop="clickItem(folderInfo)"
      >{{ limit30(folderInfo.title) }}
        <CheckOutlined />
      </span>
    </template>
    <FolderMenuItem
      v-for="item in folderInfo.children"
      :key="item.key"
      :folder-info="item"
      :disabled="detectDisable(item.key)"
      :detect-disable="detectDisable"
      :selected-folder-map="selectedFolderMap"
      :new-folder-parent="newFolderParent"
      @clickItem="clickItem($event)"
      @submitFolder="submitFolder($event)"
      @cancelFolder="cancelFolder()"
      @newFolder="newFolder($event)"
    />
    <Item v-if="newFolderParent === folderInfo">
      <input
        v-model="newFolderName"
        :data-folder-key="newFolderParent.key"
        class="metadata-folder-input blue-6-b"
        @keypress.stop.enter="submitFolder(newFolderName)"
        @keypress.stop.esc="cancelFolder()"
        @keydown.stop="keydown"
        @keyup.stop
        @click.stop
      >
    </Item>
    <Item v-else>
      <div
        class="metadata-folder-new blue-6-c"
        @click.stop="newFolder(folderInfo)"
      >
        <PlusOutlined class="blue-6-c" />{{ $t('folder.new_folder') }}
      </div>
    </Item>
  </SubMenu>
</template>
<script lang="ts" setup>
import { Menu } from 'ant-design-vue';
import { CheckOutlined, PlusOutlined } from '@ant-design/icons-vue';
import { SelectedFolderMap, limit30 } from './helper.js';
import FolderMenuItem from './FolderMenuItem.vue';
import { ref } from 'vue';
import { PaperFolder } from '~/src/stores/folderStore';

const { Item, SubMenu } = Menu;

const props = defineProps<{
  folderInfo: PaperFolder;
  selectedFolderMap: SelectedFolderMap;
  detectDisable(key: string): boolean;
  newFolderParent: PaperFolder | null;
}>();

const emit = defineEmits<{
  (event: 'clickItem', folder: PaperFolder): void
  (event: 'submitFolder', name: string): void
  (event: 'cancelFolder'): void
  (event: 'newFolder', folder: PaperFolder): void
}>()

const newFolderName = ref('')

const clickItem = (folder: PaperFolder) => {
  emit('clickItem', folder)
}

const submitFolder = async(name: string) => {
  emit('submitFolder', name)
  setTimeout(() => {
    newFolderName.value = ''
  }, 300)
}

const cancelFolder = () => {
  emit('cancelFolder')
  newFolderName.value = ''
}

const newFolder = (folder: PaperFolder) => {
  emit('newFolder', folder)
}

const keydown = (event: KeyboardEvent) => {
  if (event.key === 'Escape') {
    cancelFolder()
  }
}

</script>
<style lang="less">
@import '~/src/assets/less/style.less';
@import './style.less';

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
    > span {
      color: var(--site-theme-brand);
      visibility: visible !important;
    }
  }
}

.metadata-folder-selected {
  background-color: var(--site-theme-brand-bg);
}
</style>
