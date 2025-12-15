<script lang="ts" setup>
import { Dropdown, Menu } from 'ant-design-vue';
import { PlusOutlined } from '@ant-design/icons-vue';
import FolderMenuItem from './FolderMenuItem.vue';
import { useFolderStore } from '~/src/stores/folderStore';
import { computed, ref, watch } from 'vue';
import {
  copyDocOrFolderToAnotherFolder,
  getFolders,
  removeDocFromFolder,
} from '~/src/api/material';

const folderStore = useFolderStore();

const detectDisable = () => false;

const props = defineProps<{
  docId: string;
  trigger?: 'click' | 'hover' | 'contextmenu';
  getPopupContainer?(): HTMLElement;
}>();

const { Item } = Menu;

folderStore.init();

const folderIdList = ref<string[]>([]);
const selectedFolderMap = computed(() => {
  const map: Record<string, boolean> = {};
  folderIdList.value.forEach((id) => {
    map[id] = true;
  });
  return map;
});
const fetchSelectedFolder = async () => {
  if (!props.docId) {
    return;
  }

  const response = await getFolders({
    noteId: 0 as unknown as string,
    docId: props.docId,
  });
  folderIdList.value = response.map((item) => item.id);
};
fetchSelectedFolder();
watch(() => props.docId, fetchSelectedFolder);

const clickItem = async (folderId: string) => {
  const exist = folderIdList.value.includes(folderId);
  if (exist) {
    await removeDocFromFolder({
      removedDocItems: [
        {
          docId: props.docId,
          folderId,
        },
      ],
      isHierarchicallyRemove: false,
    });
  } else {
    await copyDocOrFolderToAnotherFolder({
      docIds: [props.docId],
      folderIds: [],
      targetFolderId: folderId,
    });
  }

  await fetchSelectedFolder();
};

const dropdownVisible = ref(false);

const visibleChange = async (visible: boolean) => {
  if (visible) {
    dropdownVisible.value = true;
    return;
  }

  if (folderStore.newFolderParent && folderStore.getFolderInput()?.offsetParent)
    return;

  folderStore.cancelFolder();
  dropdownVisible.value = false;
};

const defaultGetPopupContainer = () => document.body;

const keydown = (event: KeyboardEvent) => {
  if (event.key === 'Escape') {
    folderStore.cancelFolder();
  }
};
</script>

<template>
  <Dropdown
    overlay-class-name="metadata-dropdown"
    :destroy-popup-on-hide="true"
    :get-popup-container="getPopupContainer || defaultGetPopupContainer"
    :visible="dropdownVisible"
    :trigger="trigger || 'hover'"
    placement="bottomLeft"
    @visible-change="visibleChange"
  >
    <div>
      <slot />
    </div>
    <template #overlay>
      <Menu
        :mode="'vertical'"
        :selected-keys="folderIdList"
        theme="light"
        class="metadata-menu thin-scroll"
      >
        <Item class="metadata-folder-count">
          <div
            class="metadata-folder-new"
          >
            {{ $t('folder.saved_to')
            }}{{ folderIdList.length ? `（${folderIdList.length}）` : '' }}
          </div>
        </Item>
        <FolderMenuItem
          v-for="item in folderStore.folderList"
          :key="item.key"
          :folder-info="item"
          :selected-folder-map="selectedFolderMap"
          :disabled="detectDisable()"
          :detect-disable="detectDisable"
          @click-item="clickItem($event)"
        />
        <Item v-if="folderStore.newFolderParent?.key === '0'">
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
            @click.stop="
              folderStore.newFolder({
                key: '0',
                title: $t('folder.all_papers'),
                path: [],
                children: folderStore.folderList,
              })
            "
          >
            <PlusOutlined class="blue-6-c" />{{ $t('folder.new_folder') }}
          </div>
        </Item>
      </Menu>
    </template>
  </Dropdown>
</template>

<style lang="less" scoped>
.first-select {
  position: absolute;
  margin-top: -30px;
  margin-left: 95px;
  max-height: 300px;
  overflow-y: auto;
}
.ant-dropdown-link::after {
  content: '>';
  margin-left: 25px;
}
</style>

<style lang="less">
// @import '~/src/assets/less/style.less';
@import './folder-style.less';
</style>
