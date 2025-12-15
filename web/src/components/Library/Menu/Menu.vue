<template>
  <div class="literature-tree-container list-tree">
    <Tree :tree-data="storeLibraryIndex.libraryIndexList" v-model:expandedKeys.sync="storeLibraryMenu.expandedKeys"
      v-model:selectedKeys="treeSelectedKeys" :draggable="isSafari" @click="clickItem" @rightClick="openMenu"
      @drop="onDrop as any">
      <template #title="{
          title,
          key,
          isLeaf,
          pdfId,
          docCount,
          isEditing,
          docId,
          paperId,
        }">
        <Dropdown :trigger="['contextmenu']">
          <span v-if="key === allKey" class="root-folder" :data-key="key">
            <span class="folder-count">
              <CopyOutlined class="aiknowledge-icon iconfile_copy" />
              <span class="root-title">{{
                $t('home.library.folder.all_papers')
                }}</span>
              <span class="root-count">
                {{ docCount }}
              </span>
              <Tooltip placement="bottom" overlay-class-name="literature-tips" :overlay-style="{
                  whiteSpace: 'pre-line',
                }" :get-popup-container="getPopupContainer">
                <template #title>{{ $t('home.library.tips') }}</template>
                <QuestionCircleOutlined style="color: #86919c; margin-left: 10px" />
              </Tooltip>
            </span>

            <MenuFoldOutlined class="aiknowledge-icon iconoutdent"
              @click.stop.prevent.capture="storeLibraryMenu.loadTree = false" />
          </span>
          <span v-else-if="key === seperatorKey" class="root-folder" :data-key="key">
            <i class="unclassified-literature" aria-hidden="true" />
            <Tooltip>
              <template #title>
                {{ $t('home.library.folder.no_classified') }}
              </template>
              <span class="title-unclassified-selected">{{
                $t('home.library.folder.no_classified')
                }}</span>
            </Tooltip>
            <span v-if="!isLeaf" class="count-selected"> </span>
          </span>
          <Tooltip v-else placement="right" overlay-class-name="list-tree-tooltip"
            :get-popup-container="getPopupContainer">
            <template v-if="key !== newFolderKey" #title>{{ title }}</template>
            <span class="dis-flex" :class="{ isLeaf }" :data-key="key" @mousedown="onMouseDown(key)"
              @dragstart="onDragStart">
              <div v-if="storeLibraryMenu.dragKey" class="drop-before" @mouseup="onMouseUp(key, -1)"
                @dragstart="onDragStart">
              </div>
              <div class="drop-middle" :class="{
                  'drop-middle-over': storeLibraryMenu.dragKey,
                  'drop-middle-drag': storeLibraryMenu.dragKey === key,
                }" :style="{
                  outline:
                    storeLibraryMenu.dragKey === key ? 'dashed 1px gray' : '0',
                  opacity: storeLibraryMenu.dragKey === key ? 0.75 : 1,
                }" @mouseup="storeLibraryMenu.dragKey && onMouseUp(key, 0)" @dragstart="onDragStart">
                <template v-if="!isLeaf">
                  <FolderOpenFilled v-if="storeLibraryMenu.expandedKeys.includes(key)"
                    class="aiknowledge-icon iconfolder_open_fill" style="color: #919fb5" />
                  <FolderFilled v-else class="aiknowledge-icon iconfolder_fill" style="color: #919fb5" />
                </template>
                <i v-else-if="pdfId && pdfId !== '0'" class="aiknowledge-icon icon-file-pdf" style="color: #919fb5"
                  aria-hidden="true" />
                <LinkOutlined v-else class="aiknowledge-icon iconfile_link" style="color: #919fb5" />

                <span :style="{
                    width:
                      storeLibraryIndex.libraryIndexExtra[key].width + 'px',
                  }" class="folder-selected">
                  <input v-if="isEditing" :id="`list-tree-input:${key}`" :value="title" class="list-tree-input" :style="{
                      width:
                        storeLibraryIndex.libraryIndexExtra[key].width + 'px',
                    }" :placeholder="
                      $t(
                        `home.library.${
                          isLeaf ? 'enterLiteratureName' : 'enterFolderName'
                        }`
                      )
                    " @click.stop @contextmenu.stop.capture.prevent @blur="renameEnd($event, key)"
                    @keypress.enter="renameEnd($event, key)" @keypress.esc="
                      key === newFolderKey ? addFolderCancel() : renameCancel()
                    " />
                  <span v-if="!isEditing" class="title-selected">{{
                    title
                    }}</span>
                  <span v-if="!isEditing && !isLeaf" class="count-selected">
                    {{ docCount }}
                  </span>
                </span>
              </div>
              <div v-if="storeLibraryMenu.dragKey" class="drop-after" @mouseup="onMouseUp(key, 1)"
                @dragstart="onDragStart"></div>
            </span>
          </Tooltip>
          <template #overlay>
            <Menu v-if="key === allKey">
              <Item @click="showFileUploader">{{
                $t('home.library.addLiterature')
                }}</Item>
              <Item @click="allAddChildFolderStart">{{
                $t('home.library.addFolder')
                }}</Item>
              <Item @click="refresh">{{ $t('home.global.refresh') }}</Item>
            </Menu>

            <Menu v-else-if="isLeaf">
              <Item @click="renameStart">{{ $t('home.global.rename') }}</Item>
              <Item :title="$t('home.library.moveTo')">
                <SubMenu v-if="storeLibraryIndex.selectedKey === key" :list="storeLibraryIndex.libraryIndexList"
                  :type="1" :detect-disable="detectDisable" :move-or-copy="moveOrCopy" :has-slot="false">
                </SubMenu>
              </Item>
              <Item v-if="
                  storeLibraryIndex.libraryIndexExtra[key] &&
                  storeLibraryIndex.libraryIndexExtra[key].path.length >= 2
                " :title="$t('home.library.copyTo')">
                <SubMenu v-if="storeLibraryIndex.selectedKey === key" :list="storeLibraryIndex.libraryIndexList"
                  :type="2" :detect-disable="detectDisable" :move-or-copy="moveOrCopy" :has-slot="false">
                </SubMenu>
              </Item>
              <Item>
                <Cite :paper-id="paperId" :pdf-id="pdfId" :page-type="PageType.library">
                  {{ $t('home.paper.quote') }}
                </Cite>
              </Item>
              <Divider />
              <Item class="del-doc" @click="removeModalVisible = true">{{
                $t('home.library.delLiterature')
                }}</Item>
            </Menu>

            <Menu v-else>
              <Item @click="showFileUploader">{{
                $t('home.library.addLiterature')
                }}</Item>
              <Item @click="renameStart">{{ $t('home.global.rename') }}</Item>
              <Item @click="addSiblingFolderStart">{{
                $t('home.library.addFolder')
                }}</Item>
              <Item @click="addChildFolderStart">{{
                $t('home.library.addSubFolder')
                }}</Item>
              <Item key="test" :title="$t('home.library.moveTo')" style="position: relative">
                <SubMenu v-if="storeLibraryIndex.selectedKey === key" :list="storeLibraryIndex.libraryIndexList"
                  :type="1" :detect-disable="detectDisable" :move-or-copy="moveOrCopy" :has-slot="false">
                </SubMenu>
              </Item>
              <Item :title="$t('home.library.copyTo')">
                <SubMenu v-if="storeLibraryIndex.selectedKey === key" :list="storeLibraryIndex.libraryIndexList"
                  :type="2" :detect-disable="detectDisable" :move-or-copy="moveOrCopy" :has-slot="false">
                </SubMenu>
              </Item>
              <!--  暂时关闭分享文件夹功能  -->
              <!-- <Item :title="$t('home.library.shareFolder')" @click="
                  shareFolder = {
                    key: removePrefix(key),
                    title,
                  }
                ">{{ $t('home.library.shareFolder') }}</Item> -->
              <Item :disabled="docCount <= 0" @click="handleExportBibTex(removePrefix(key))">
                {{ $t('home.library.exportBibtex') }}
              </Item>
              <Divider />
              <Item class="del-folder" @click="removeModalVisible = true">{{
                $t('home.library.delFolder')
                }}</Item>
            </Menu>
          </template>
        </Dropdown>
      </template>
    </Tree>
    <RemoveModal v-if="removeModalVisible" :item-list="
        storeLibraryMenu.selectedNodeFamily
          ? [storeLibraryMenu.selectedNodeFamily.node]
          : []
      " :parent="
        storeLibraryMenu.selectedNodeFamily &&
        storeLibraryMenu.selectedNodeFamily.parent
      " :is-folder="
        !(
          storeLibraryMenu.selectedNodeFamily &&
          storeLibraryMenu.selectedNodeFamily.node &&
          storeLibraryMenu.selectedNodeFamily.node.isLeaf
        )
      " :has-attachment="
        storeLibraryMenu.selectedPaper &&
        storeLibraryMenu.selectedPaper.hasAttachment
      " :remove="false" @close="removeModalVisible = false" @dirty="dirty()" @refresh="refreshMenuListParams()" />
    <ShareFolder :visible="!!shareFolder.key" :folder-key="shareFolder.key" :folder-title="shareFolder.title"
      @cancel="shareFolder = emptyShareFolder" />
  </div>
</template>
<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { message, Dropdown, Menu, Tree, Tooltip } from 'ant-design-vue'
import {
  CopyOutlined,
  FolderOpenFilled,
  FolderFilled,
  LinkOutlined,
  MenuFoldOutlined,
  QuestionCircleOutlined,
} from '@ant-design/icons-vue'
import fileDownload from 'js-file-download'
import SubMenu from './SubMenu.vue'
import ShareFolder from './ShareFolder.vue'
import RemoveModal from './MenuRemoveModal.vue'
import {
  allKey,
  seperatorKey,
  newFolderKey,
  useLibraryIndex,
  removePrefix,
  noteDirty,
} from '@/stores/library'
import { useLibraryMenu } from '@/stores/library/menu'
import { useLibraryList } from '@/stores/library/list'
import { liteThrottle } from '@idea/aiknowledge-special-util'
import { exportBibTexByFolderId } from '@/api/document'
import { PageType, reportElementClick } from '@/utils/report'
import { useUserStore } from '@/common/src/stores/user'
import { useI18n } from 'vue-i18n'
import { LIBRARY_CONTAINER_CLASSNAME } from '../helper'
import Cite from '@/components/Paper/Quote/cite.vue'

const { Item, Divider } = Menu

const i18n = useI18n()

const userStore = useUserStore()

const storeLibraryIndex = useLibraryIndex()

const storeLibraryMenu = useLibraryMenu()

;(window as any).storeLibraryIndex = storeLibraryIndex
;(window as any).storeLibraryMenu = storeLibraryMenu

const treeSelectedKeys = computed({
  get: () => {
    return [storeLibraryIndex.selectedKey]
  },
  set: (keys) => {
    const [key] = keys
    if (key) {
      storeLibraryIndex.selectedKey = key
    }
  },
})

const {
  clickItem,
  openMenu,
  mouseDown,
  dropTo,
  renameStart,
  renameCancel,
  renameSubmit,
  addSiblingFolderStart,
  addChildFolderStart,
  addFolderCancel,
  addFolderSubmit,
  moveTo,
  copyTo,
} = storeLibraryMenu

const allAddChildFolderStart = () => {
  if (!userStore.userInfo) {
    message.info(i18n.t('common.tips.loginFirst'))
    return
  }

  addChildFolderStart()
}

const renameEnd = liteThrottle(
  ($event: Event, key: string) => {
    key === newFolderKey ? addFolderSubmit($event) : renameSubmit($event)
  },
  1000,
  false,
  true
)

const storeLibraryList = useLibraryList()
const paperMap = computed(() => {
  const res: Record<string, any> = {}

  storeLibraryList.paperListAll.forEach((paper) => {
    res[paper.docId] = paper
  })

  return paperMap
})

const isSafari =
  typeof navigator !== 'undefined' &&
  /^((?!chrome|android).)*safari/i.test(navigator.userAgent)
const onMouseDown = isSafari ? () => {} : mouseDown
const onMouseUp = isSafari ? () => {} : dropTo
const onDragStart = !isSafari
  ? () => {}
  : (event: DragEvent) => {
      event.preventDefault()
      event.stopPropagation()
    }

const onDrop = !isSafari
  ? () => {}
  : (info: {
      node: { eventKey: any; pos: string }
      dragNode: { eventKey: any }
      dropPosition: number
    }) => {
      const dragKey = info.dragNode.eventKey
      const dropKey = info.node.eventKey
      const dropPos = info.node.pos.split('-')
      const dropPosition =
        info.dropPosition - Number(dropPos[dropPos.length - 1])

      dropTo(dropKey, dropPosition, dragKey)
    }

onMounted(() => {
  storeLibraryIndex.fetchLibraryIndex()
})

const refresh = () => {
  storeLibraryIndex.fetchLibraryIndex()
  storeLibraryList.getFilesByFolderId()
}

const showFileUploader = () => {
  storeLibraryIndex.uploaderVisible = true
}
const emptyShareFolder = { key: '', title: '' }
const shareFolder = ref(emptyShareFolder)
const detectDisable = (key: string) => {
  return key === storeLibraryIndex.selectedKey
}
const moveOrCopy = (type: 1 | 2, { key }: { key: string }) => {
  console.log('[DEBUG][Menu] moveOrCopy 被调用:', type === 1 ? '移动到' : '复制到', '目标文件夹键值:', key);
  console.log('[DEBUG][Menu] 当前文件夹列表状态:', storeLibraryIndex.libraryIndexList.map(item => ({ key: item.key, title: item.title })));
  
  if (type === 1) {
    moveTo(key)
  } else {
    copyTo(key)
  }
  
  // 延迟打印操作后的状态
  setTimeout(() => {
    console.log('[DEBUG][Menu] moveOrCopy 操作后的文件夹列表状态:', storeLibraryIndex.libraryIndexList.map(item => ({ key: item.key, title: item.title })));
  }, 100);
}

const removeModalVisible = ref(false)

const handleExportBibTex = async (key: string) => {
  reportElementClick({
    page_type: PageType.library,
    type_parameter: 'none',
    element_name: 'generate_bibtex',
    status: 'none',
  })

  try {
    const data = await exportBibTexByFolderId({
      folderId: key,
    })

    if (!data) {
      message.warn('生成失败，文件为空')
    } else {
      fileDownload(data, 'export.bib')
    }
  } catch (error) {
    //
  }
}

const getPopupContainer = () => {
  return (document.getElementsByClassName(LIBRARY_CONTAINER_CLASSNAME)[0] ||
    document.body) as HTMLDivElement
}

const dirty = () => {
  noteDirty.is = true
}
const refreshMenuListParams = () => {
  storeLibraryIndex.fetchLibraryIndex()
  storeLibraryList.refreshClassiyAuthorVenuePaperList()
}
</script>
<style lang="less" scoped>
@import './Menu.less';
</style>
<style>
.ant-dropdown-placement-bottomLeft {
  position: absolute;
  margin-top: -67px;
  margin-left: 68px;
}

.ant-dropdown-menu .del-doc,
.ant-dropdown-menu .del-folder {
  color: #f00;
}

.ant-dropdown-menu .del-doc:hover,
.ant-dropdown-menu .del-folder:hover {
  color: #fff;
  background: #f00;
}

.list-tree-input {
  height: 32px;
  color: black;
  border: 1px solid #a19f9d;
  padding-left: 5px;
}

.literature-tips .ant-tooltip-inner,
.note-source-tips .ant-tooltip-inner {
  color: #4e5969;
  background: white;
  font-size: 12px;
  padding: 12px 16px;
  line-height: 18px;
}

.literature-tips .ant-tooltip-inner {
  width: 256px;
}

.literature-tips .ant-tooltip-arrow::before,
.note-source-tips .ant-tooltip-arrow::before {
  background: white;
}

.list-tree>.ant-tree {
  background: transparent !important;
}

/* 修改树形菜单中文档项的文字颜色 */
.ant-tree-node-content-wrapper .ant-tree-title {
  color: #262625 !important;
  font-weight: 500;
}

/* 确保文档图标可见 */
.ant-tree .iconfile_pdf,
.ant-tree .iconfile_link {
  color: #4e5969 !important;
}

/* 选中状态下的文字颜色 */
.ant-tree-node-selected .ant-tree-title {
  color: #ffffff !important;
}
</style>
