<template>
  <div class="library-wrap">
    <div :class="[LIBRARY_CONTAINER_CLASSNAME]">
      <LibraryMenu v-show="storeLibraryMenu.loadTree" />

      <div
        class="scroll"
        style="display: flex; flex-direction: column; height: calc(100vh - 64px);"
      >
        <div class="table-header">
          <div style="display: flex; margin-left: 20px">
            <div
              v-if="!storeLibraryMenu.loadTree"
              style="margin-right: 40px; cursor: pointer"
              @click="storeLibraryMenu.loadTree = true"
            >
              <strong>{{ $t('home.library.showFolder') }}</strong>
              <MenuUnfoldOutlined style="margin-left: 8px" />
            </div>
            <Breadcrumb
              v-if="storeLibraryIndex.selectedKey"
              style="cursor: pointer"
            >
              <BreadcrumbItem
                v-for="(item, index) in storeLibraryIndex.breadCrumbList"
                :key="index"
                @click.native="breadCrumbClick(item)"
                >{{
                  item.folderId === allKey
                    ? $t('home.library.folder.all_papers')
                    : item.docName
                }}</BreadcrumbItem
              >
            </Breadcrumb>
          </div>

          <div v-show="!storeLibraryList.paperListEmpty" class="banner">
            <div class="search">
              <Input
                v-model:value="storeLibraryList.searchInput"
                allow-clear
                :placeholder="$t('home.library.searchTips')"
              >
                <template #prefix>
                  <SearchOutlined />
                </template>
              </Input>
            </div>
            <Select
              class="select"
              :default-value="storeLibraryList.dropdownSortType"
              :get-popup-container="getPopupContainer"
              @change="changeSortType as any"
              :dropdown-style="{ backgroundColor: 'var(--site-theme-background-primary)' }"
              :dropdown-class-name="'theme-dropdown'"
            >
              <SelectOption :value="UserDocListSortType.LAST_ADD">{{
                $t('home.library.currentAdd')
              }}</SelectOption>
              <SelectOption :value="UserDocListSortType.LAST_READ">{{
                $t('home.library.currentRead')
              }}</SelectOption>
            </Select>
            <div class="divider" />
            <div class="btn">
              <Button @click="collating = !collating" class="ant-btn-default">{{
                collating
                  ? $t('home.library.endMultipleSelect')
                  : $t('home.library.multipleSelect')
              }}</Button>
            </div>
            <div class="btn">
              <Button
                type="primary"
                @click="storeLibraryIndex.uploaderVisible = true"
                >{{ $t('home.global.add') }}</Button
              >
            </div>
          </div>
        </div>
        <div v-if="collating" class="collating-button-group">
          <div>
            <Button
              :class="[
                'collating-delete',
                { 'disabled-delete-btn': itemList.length === 0 },
              ]"
              :danger="true"
              :disabled="itemList.length === 0"
              @click="removeModalVisible = true"
            >
              {{ $t('home.global.delete') }}
            </Button>
            <RemoveModal
              v-if="removeModalVisible"
              :item-list="itemList"
              :parent="
                storeLibraryIndex.currentFolder.folderId === allKey
                  ? null
                  : {
                      title: storeLibraryIndex.currentFolder.docName,
                      key: storeLibraryIndex.currentFolder.folderId,
                    }
              "
              :remove="true"
              :has-attachment="itemListHasAttachments"
              @close="removeModalVisible = false"
              @refresh="handleAddSuccess()"
            />
          </div>
          <div>
            <Button v-if="itemList.length === 0" type="primary" disabled>{{
              $t('home.library.moveTo')
            }}</Button>
            <SubMenu
              v-else
              :list="storeLibraryIndex.libraryIndexList"
              :type="1"
              :has-slot="true"
              :detect-disable="detectDisable"
              :move-or-copy="moveOrCopy"
            >
              <Button type="primary">{{ $t('home.library.moveTo') }}</Button>
            </SubMenu>
          </div>
          <div>
            <Button v-if="itemList.length === 0" type="primary" disabled>{{
              $t('home.library.copyTo')
            }}</Button>
            <SubMenu
              v-else
              :list="storeLibraryIndex.libraryIndexList"
              :type="2"
              :has-slot="true"
              :detect-disable="detectDisable"
              :move-or-copy="moveOrCopy"
            >
              <Button type="primary">{{ $t('home.library.copyTo') }}</Button>
            </SubMenu>
          </div>
          <Button
            type="primary"
            :loading="bibtexLoading"
            :disabled="itemList.length === 0"
            :style="{
              cursor: itemList.length ? 'pointer' : 'not-allowed',
            }"
            @click="createBibtex()"
            >{{ $t('home.library.bibtex') }}</Button
          >
        </div>
        <LoginFirst style="flex: 1; min-height: 0; display: flex; flex-direction: column;">
          <template #blank>
            <img
              class="w-20"
              src="/src/common/assets/images/notes/empty-notes.svg"
              alt="empty"
            />
            <p>
              {{ $t('common.tips.empty', [$t('home.library.literature')]) }}
            </p>
          </template>
          <LibraryList
            :style="{ marginLeft: storeLibraryMenu.loadTree ? '12px' : '0', flex: '1', minHeight: '0' }"
            :collating="collating"
            @buttonClick="storeLibraryIndex.uploaderVisible = true"
          />
        </LoginFirst>
        <FileUploader
          :visible="storeLibraryIndex.uploaderVisible"
          :folder-id="storeLibraryIndex.rawFolderId"
          :collect-limit-dialog-report-params="collectLimitDialogReportParams"
          :selected-key="storeLibraryIndex.selectedKey"
          :bread-crumb-list="storeLibraryIndex.breadCrumbList"
          @close="storeLibraryIndex.uploaderVisible = false"
          @addSuccess="handleAddSuccess"
          @refreshList="handleAddSuccess"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted, nextTick, computed } from 'vue'
import {
  CopyDocOrFolderToAnotherFolderReq,
  MoveDocOrFolderToAnotherFolderReq,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/client/ClientDoc'
import {
  message,
  Button,
  Breadcrumb,
  BreadcrumbItem,
  Select,
  SelectOption,
  Input,
} from 'ant-design-vue'
import { MenuUnfoldOutlined, SearchOutlined } from '@ant-design/icons-vue'
import fileDownload from 'js-file-download'
import {
  UserDocInfo,
  UserDocListSortType,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/UserDocManage'
import FileUploader from '@/components/Library/File/Uploader.vue'
import LibraryList from '@/components/Library/List/List.vue'
import LibraryMenu from '@/components/Library/Menu/Menu.vue'
import SubMenu from '@/components/Library/Menu/SubMenu.vue'
import RemoveModal from '@/components/Library/Menu/MenuRemoveModal.vue'
import {
  allKey,
  BreadCrumb,
  removePrefix,
  noteDirty,
  useLibraryIndex,
} from '@/stores/library'
import { useLibraryList } from '@/stores/library/list'

import {
  copyDocOrFolderToAnotherFolder,
  exportBibTexByIds,
  moveDocOrFolderToAnotherFolder,
} from '@/api/document'
import { LimitDialogReportParams } from '@/components/Library/helper'
import { ElementName, PageType, reportElementClick } from '@common/utils/report'
import { useI18n } from 'vue-i18n'
import LimitDialog from '@/common/src/components/Premium/LimitDialog.vue'
import { useLibraryMenu } from '@/stores/library/menu'
import LoginFirst from '@/common/src/components/LoginFirst/index.vue'
import { useUserStore } from '@/common/src/stores/user'
import { LIBRARY_CONTAINER_CLASSNAME } from './helper'

const { t } = useI18n()
const userStore = useUserStore()

const collating = ref(false)
const storeLibraryIndex = useLibraryIndex()
const storeLibraryList = useLibraryList()
const storeLibraryMenu = useLibraryMenu()

const loadNotice = ref(false as boolean)

const breadCrumbClick = (e: BreadCrumb) => {
  nextTick(() => {
    storeLibraryIndex.selectedKey = e.folderId
  })
}

const changeSortType = (
  sort: UserDocListSortType.LAST_ADD | UserDocListSortType.LAST_READ
) => {
  storeLibraryList.currentSortType = sort
  storeLibraryList.dropdownSortType = sort
  storeLibraryList.paperListLocalSortDirection = 1

  if (!userStore.userInfo) {
    return
  }

  storeLibraryList.getFilesByFolderId()
}

const initLoad = () => {
  if (localStorage.getItem('notice') !== '0526') {
    loadNotice.value = true
  } else {
    loadNotice.value = false
  }
}

const closeNotice = () => {
  loadNotice.value = false
}

const handleAddSuccess = () => {
  storeLibraryIndex.fetchLibraryIndex()
  storeLibraryList.refreshClassiyAuthorVenuePaperList()
}

const handleUploadFinished = () => {
  storeLibraryList.refreshClassiyAuthorVenuePaperList()
}

onMounted(() => {
  initLoad()
  
  // 监听上传完成事件，自动刷新列表
  window.addEventListener('uploadFinished', handleUploadFinished);
})

onUnmounted(() => {
  // 移除事件监听
  window.removeEventListener('uploadFinished', handleUploadFinished)
})

const docIdList = computed(() => {
  return Object.keys(storeLibraryList.paperListCheckedMap).filter(
    (id) => storeLibraryList.paperListCheckedMap[id]
  )
})
const docList = computed(() => {
  const value = storeLibraryList.paperListAll
  const map: Record<UserDocInfo['docId'], UserDocInfo> = {}
  value.forEach((item) => {
    map[item.docId] = item
  })
  return docIdList.value.map((id) => map[id]!)
})
const itemList = computed(() => {
  return docList.value.map((item) => ({
    title: item.docName,
    key: item.docId,
  }))
})
const itemListHasAttachments = computed(() => {
  return docList.value.some((item) => item.hasAttachment)
})
const detectDisable = () => false
const moveOrCopy = async (type: 1 | 2, { key }: { key: string }) => {
  const targetFolderId = removePrefix(key)
  noteDirty.is = true

  if (type === 1) {
    const req: MoveDocOrFolderToAnotherFolderReq = {
      movedDocItems: docList.value.map((item) => ({
        docId: item?.docId,
        sourceFolderId: storeLibraryIndex.rawFolderId,
      })),
      movedFolderIds: [],
      targetFolderId,
      /** true=将选中的文献从当前文件夹及其子文件夹下迁移出去,false=只将选中的文献从当前文件夹下迁移出去 */
      isHierarchicallyMoveDoc: true,
    }
    await moveDocOrFolderToAnotherFolder(req)
  } else {
    const req: CopyDocOrFolderToAnotherFolderReq = {
      docIds: docIdList.value,
      folderIds: [],
      targetFolderId,
    }
    await copyDocOrFolderToAnotherFolder(req)
  }

  storeLibraryIndex.fetchLibraryIndex()
  storeLibraryList.refreshClassiyAuthorVenuePaperList()
}
const removeModalVisible = ref(false)
const bibtexLoading = ref(false)
const createBibtex = async () => {
  reportElementClick({
    page_type: PageType.library,
    type_parameter: 'none',
    element_name: 'generate_bibtex',
    status: undefined,
  })

  bibtexLoading.value = true
  try {
    const data = await exportBibTexByIds({
      docIds: docIdList.value,
    })
    if (!data) {
      message.warn(t('home.library.generateFail') as string)
    } else {
      fileDownload(data, 'export.bib')
    }
  } catch (error) {
    //
  }
  bibtexLoading.value = false
}

const collectLimitDialogReportParams: LimitDialogReportParams = {
  page_type: PageType.library,
  element_name: ElementName.upperCollectionPopup,
}

const getPopupContainer = (triggerNode: HTMLElement) =>
  (triggerNode.parentNode || document.body) as HTMLElement
</script>

<style lang="less" scoped>
.library-wrap {
  overflow: hidden;
  width: 100%;
  display: flex;
  flex-direction: column;

  :deep {
    .ant-select-selection,
    .ant-select-selector {
      border: none !important;
      background-color: var(--site-theme-background-secondary) !important;
      color: var(--site-theme-text-color) !important;
    }
    
    // 修复选中和聚焦状态
    .ant-select-open,
    .ant-select-focused {
      
      .ant-select-selection-item {
        color: var(--site-theme-primary-text) !important;
      }
      
      .ant-select-arrow {
        color: var(--site-theme-primary-color) !important;
      }
    }
    
    // 修复搜索框输入文字透明的问题
    .ant-input {
      color: var(--site-theme-text-color) !important;
      opacity: 1 !important;
    }
    .anticon-search {
      color: var(--site-theme-text-color) !important;
    }
    
    .ant-input::placeholder {
      opacity: 0.6;
      color: var(--site-theme-text-color) !important;
    }
    
    // 覆盖禁用状态的样式，确保在浅色模式下可见
    a[disabled] {
      color: var(--site-theme-text-disabled) !important;
      cursor: not-allowed;
    }
    
    // 覆盖禁用状态的按钮样式
    button[disabled] {
      color: var(--site-theme-text-disabled) !important;
      background-color: var(--site-theme-background-disabled) !important;
      border-color: var(--site-theme-border-color) !important;
      cursor: not-allowed;
    }
  }
}
.readpaper-library-container {
  position: relative;
  flex: 1 1 100%;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  border-top: 1px solid var(--site-theme-border-color);
  overflow: hidden;
  .scroll {
    position: relative;
    flex: 1 1 100%;
    overflow: hidden;
    background-color: var(--site-theme-background-primary);

    .table-header {
      background: var(--site-theme-background-secondary);
      padding: 8px 16px 8px 1px;
      margin-bottom: 12px;
    }
  }

  .banner {
    display: flex;
    justify-content: space-between;
    margin-top: 10px;
    margin-left: 20px;
    margin-bottom: 12px;
    align-items: center;

    .total {
      margin-right: 12px;
      color: var(--site-theme-text-tertiary);
    }

    .search {
      flex: 1;
    }

    .select {
      margin-left: 32px;
    }

    .divider {
      margin: 0 24px 0 20px;
      height: 20px;
      width: 1px;
      background: var(--site-theme-border-color);
    }

    .btn + .btn {
      margin-left: 16px;
    }
  }

  .ant-breadcrumb-link {
    cursor: pointer;
  }
  :deep(.ant-breadcrumb-separator) {
    color: var(--site-theme-text-primary);
  }

  .ant-breadcrumb span {
    color: var(--site-theme-primary-color);
  }

  .ant-breadcrumb span:last-child {
    color: var(--site-theme-text-primary);
  }

  :deep(.ant-pagination-item-active) {
    background: var(--site-theme-primary-color);
  }
}

.collating-button-group {
  margin-left: 16px;
  margin-bottom: 12px;
  display: flex;
  flex-direction: row;

  .collating-delete {
    color: var(--site-theme-text-inverse) !important;
    background: var(--site-theme-danger-color) !important;
    border-color: var(--site-theme-danger-color) !important;
  }

  .disabled-delete-btn {
    background: var(--site-theme-danger-color) !important;
    border-color: var(--site-theme-danger-color);
  }

  > * {
    margin-right: 12px;
  }
}
</style>
