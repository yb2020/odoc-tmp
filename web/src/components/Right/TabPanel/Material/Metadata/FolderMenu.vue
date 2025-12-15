<template>
  <div>
    <Dropdown overlay-class-name="metadata-dropdown">
      <div>
        <slot />
      </div>
      <template #overlay>
        <Menu
          :mode="'vertical'"
          :selected-keys="selectedKeys"
          theme="light"
          class="metadata-menu metadata-scroll"
        >
          <FolderMenuItem
            v-for="item in folderList"
            :key="item.key"
            :folder-info="item"
            :disabled="props.detectDisable(item.key)"
            :detect-disable="props.detectDisable"
            :selected-folder-map="props.selectedFolderMap"
            :new-folder-parent="newFolderParent"
            @clickItem="clickItem($event)"
            @submitFolder="submitFolder($event)"
            @cancelFolder="cancelFolder()"
            @newFolder="newFolder($event)"
          />
          <Item v-if="newFolderParent?.key === '0'">
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
              @click.stop="
                newFolder({
                  key: '0',
                  title: $t('folder.all_papers'),
                  path: [],
                  children: folderList,
                })
              "
            >
              <PlusOutlined class="blue-6-c" />{{ $t('folder.new_folder') }}
            </div>
          </Item>
        </Menu>
      </template>
    </Dropdown>
  </div>
</template>

<script lang="ts" setup>
import { Dropdown, Menu } from 'ant-design-vue';
import { PlusOutlined } from '@ant-design/icons-vue';
import FolderMenuItem from './FolderMenuItem.vue';
import { SelectedFolderMap } from './helper';
import { ref } from 'vue';
import { PaperFolder } from '~/src/stores/folderStore';

const { Item } = Menu;

const props = defineProps<{
  folderList: PaperFolder[];
  selectedFolderMap: SelectedFolderMap;
  selectedKeys: PaperFolder['key'][];
  detectDisable: (key: string) => boolean;
  newFolderParent: PaperFolder | null;
}>();

const emit = defineEmits<{
  (event: 'clickItem', folder: PaperFolder): void
  (event: 'submitFolder', name: string): void
  (event: 'cancelFolder'): void
  (event: 'newFolder', folder: PaperFolder): void
}>();

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
@import '~/src/assets/less/style.less';
@import './style.less';
</style>
