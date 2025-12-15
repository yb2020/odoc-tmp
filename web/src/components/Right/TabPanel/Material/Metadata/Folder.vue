<template>
  <div class="metadata-folder-container">
    <div
      ref="triggerRef"
      class="metadata-folder-view"
    >
      {{ selectedFolderList[0]?.path ?? '' }}
    </div>
    <TippyVue
      v-if="triggerRef"
      ref="tippyRef"
      :trigger-ele="triggerRef"
      placement="left-start"
      trigger="click"
      :z-index="9"
      :delay="[0, 800]"
      :offset="[0, 76]"
    >
      <div class="metadata-folder-card">
        <CloseOutlined @click="tippyRef.hide()" />
        <div class="metadata-folder-count">
          {{
            selectedFolderList.length > 0
              ? $t('message.currentDocumentLocateFolders', {
                num: selectedFolderList.length,
              })
              : $t('message.currentDocumentNotClassified')
          }}
        </div>
        <div
          v-for="folder in selectedFolderList"
          class="metadata-folder-path"
        >
          {{ viewPath(folder.path) }}
          <Tooltip
            :title="$t('message.removeFromFolderTip')"
            placement="topRight"
          >
            <CloseCircleOutlined @click="remove(folder.id)" />
          </Tooltip>
        </div>
        <FolderMenu
          :folder-list="folderStore.folderList"
          :detect-disable="detectDisable"
          :selected-keys="selectedFolderList.map((item) => item.id)"
          :selected-folder-map="selectedFolderMap"
          :new-folder-parent="folderStore.newFolderParent"
          @clickItem="clickItem"
          @submitFolder="submitFolder"
          @cancelFolder="folderStore.cancelFolder()"
          @newFolder="folderStore.newFolder($event)"
        >
          <a-button class="metadata-folder-add">
            <PlusOutlined style="font-size: 10px" />
            {{ $t('viewer.copyToOtherFolder') }}
          </a-button>
        </FolderMenu>
      </div>
    </TippyVue>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { Tooltip } from 'ant-design-vue';
import {
  PlusOutlined,
  CloseCircleOutlined,
  CloseOutlined,
} from '@ant-design/icons-vue';
import { DocFolder } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/doc/folder';
import {
  copyDocOrFolderToAnotherFolder,
  getFolders,
  removeDocFromFolder,
} from '~/src/api/material';
import { useFolderStore } from '~/src/stores/folderStore';
import {
  SelectedFolderMap,
  PaperFolder,
  removePrefix,
  limit30,
} from './helper';
import TippyVue from '../../../../Tippy/index.vue';
import FolderMenu from './FolderMenu.vue';

const props = defineProps<{
  noteId?: string;
  docId?: string;
}>();

const folderStore = useFolderStore();

const selectedFolderList = ref<DocFolder[]>([]);
const selectedFolderMap = computed(() => {
  const map: SelectedFolderMap = {};
  walkArray(selectedFolderList.value, (folder) => {
    map[folder.id] = { ...folder };
  });
  return map;
});

function walkArray<Item>(
  array: Item[],
  func: (item: Item) => void,
  children: keyof Item = 'children' as keyof Item
) {
  array.forEach((item) => {
    func(item);

    if (Array.isArray(item[children])) {
      walkArray(item[children] as unknown as Item[], func);
    }
  });
}

const viewPath = (path: string) => {
  return path.split('/').map(limit30).join('/').replace(/^\//, '');
};

const triggerRef = ref();
const tippyRef = ref();

const fetchSelectedFolder = async () => {
  if (!props.noteId || !props.docId) {
    return;
  }

  const response = await getFolders({
    noteId: props.noteId,
    docId: props.docId,
  });

  selectedFolderList.value = response;
};

const detectDisable = () => false;

const clickItem = async (folder: PaperFolder) => {
  const key = folder.key;
  await copyDocOrFolderToAnotherFolder({
    docIds: [props.docId!],
    folderIds: [],
    targetFolderId: removePrefix(key),
  });
  return fetchSelectedFolder();
};
const remove = async (folderId: string) => {
  await removeDocFromFolder({
    removedDocItems: [
      {
        docId: props.docId!,
        folderId,
      },
    ],
    isHierarchicallyRemove: false,
  });
  selectedFolderList.value = selectedFolderList.value.filter(
    (item) => item.id !== folderId
  );
};

const submitFolder = (name: string) => {
  folderStore.newFolderName = name;
  folderStore.submitFolder()
};

onMounted(() => {
  folderStore.init();
  fetchSelectedFolder();
});

watch(() => [props.noteId, props.docId].join(':::'), fetchSelectedFolder);
</script>

<style lang="less">
.metadata-folder-container {
  flex: 1 1 100%;
  overflow: hidden;
  .metadata-folder-view {
    width: 100%;
    height: 32px;
    line-height: 32px;
    border-radius: 2px;
    overflow: hidden;
    padding-left: 8px;
    white-space: nowrap;
    text-overflow: ellipsis;
    color: var(--site-theme-brand);
    &:hover {
      color: var(--site-theme-brand-hover);
      background: var(--site-theme-bg-hover);
    }
  }
}

[data-theme='dark'] {
  .metadata-folder-card {
    background-color: var(--site-theme-pdf-panel-secondary);
  }
}

.metadata-folder-card {
  width: 332px;
  background-color: var(--site-theme-bg-light);
  padding-top: 17px;
  padding-bottom: 26px;
  padding-left: 24px;
  padding-right: 16px;
  .metadata-folder-count {
    color: var(--site-theme-text-secondary);
    margin-bottom: 16px;
  }
  .metadata-folder-path {
    color: var(--site-theme-text-primary);
    min-height: 32px;
    padding-left: 8px;
    padding-right: 10px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    word-break: break-all;
    > * {
      visibility: hidden;
    }
    &:hover {
      background: var(--site-theme-bg-hover);
      > * {
        visibility: visible;
      }
    }
  }
  > .anticon-close {
    position: absolute;
    right: 16px;
    color: var(--site-theme-text-secondary);
  }
}

.metadata-folder-add {
  margin-top: 16px;
  // height: 30px;
  color: var(--site-theme-brand);
  border: 1px solid var(--site-theme-brand);
  border-radius: 2px;
  display: flex;
  align-items: center;
  font-size: 12px;
  cursor: pointer;
  > * {
    // padding-left: 18px;
    padding-right: 10px;
  }
}
</style>
