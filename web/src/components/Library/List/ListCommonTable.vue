<template>
  <div class="file-table" :style="{ minWidth }">
    <Dropdown
      v-for="(item, index) in computedList"
      :key="`${index}-${item.pdfId}`"
      :trigger="['contextmenu']"
      @visibleChange="onVisibleChange($event, item)"
    >
      <!-- 不能销毁dropdown，有类似Cite组件需要挂载引文弹窗 -->
      <!-- destroy-popup-on-hide -->
      <Row
        :index="index"
        :item="item"
        :paper-id="item.paperId || item.pdfId"
        :collating="collating"
        @attach="toggleAttachment"
      />
      <template #overlay>
        <Menu v-model:open-keys="openKeys">
          <Item
            key="open"
            @click="openPage(item)"
            @mouseenter="handleOpen('', false)"
            >{{ $t('home.global.openPDF') }}</Item
          >
          <Item key="move" :title="$t('home.library.moveTo')">
            <SubMenu
              :list="storeLibraryIndex.libraryIndexList"
              :type="1"
              :active="openKeys.includes('move')"
              :detect-disable="detectDisable"
              :move-or-copy="moveOrCopy"
              @visible-change="handleOpen('move', $event)"
            >
            </SubMenu>
          </Item>
          <Item
            key="copy"
            v-if="(item as any).isCopy"
            :title="$t('home.library.copyTo')"
          >
            <SubMenu
              :list="storeLibraryIndex.libraryIndexList"
              :type="2"
              :active="openKeys.includes('copy')"
              :detect-disable="detectDisable"
              :move-or-copy="moveOrCopy"
              @visible-change="handleOpen('copy', $event)"
            >
            </SubMenu>
          </Item>
          <Item key="cite" @mouseenter="handleOpen('', false)">
            <Cite
              class="share"
              :paper-id="item.paperId"
              :page-type="PageType.library"
              :pdf-id="item.pdfId"
              @update:success="refresh"
              >{{ $t('home.paper.quote') }}</Cite
            >
          </Item>
          <Divider />
          <Item
            key="refresh"
            @click="refresh"
            @mouseenter="handleOpen('', false)"
            >{{ $t('home.global.refresh') }}</Item
          >
          <Item
            key="add"
            @click="addButtonClick"
            @mouseenter="handleOpen('', false)"
            >{{ $t('home.library.addLiterature') }}</Item
          >
          <Divider />
          <Item
            key="delete"
            class="del-doc"
            @click="removeModalVisible = true"
            @mouseenter="handleOpen('', false)"
            >{{ $t('home.library.delLiterature') }}</Item
          >
        </Menu>
      </template>
    </Dropdown>

    <AttachmentModal
      v-if="attachmentDocId"
      :doc-id="attachmentDocId"
      @close="toggleAttachment"
      @attached="mutateDocItem"
      @removed="mutateDocItem"
    />
    <RemoveModal
      v-if="removeModalVisible"
      :item-list="selectedNodeFamily ? [selectedNodeFamily.node] : []"
      :parent="selectedNodeFamily && selectedNodeFamily.parent"
      :has-attachment="selectedPaper && selectedPaper.hasAttachment"
      :is-folder="
        !(
          selectedNodeFamily &&
          selectedNodeFamily.node &&
          selectedNodeFamily.node.isLeaf
        )
      "
      @close="removeModalVisible = false"
      @refresh="refresh()"
    />
  </div>
</template>
<script lang="ts" setup>
import { ref, computed, provide, onUnmounted } from 'vue'
import { Dropdown, Menu } from 'ant-design-vue'
import { UserDocClassifyInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/userCenter/UserDoc'
import { UserDocInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/UserDocManage'
import Row from './Row.vue'

import {
  useLibraryIndex,
  seperatorKey,
  allKey,
  LiteratureNodeFamily,
} from '@/stores/library'
import { useLibraryMenu } from '@/stores/library/menu'

import SubMenu from '../Menu/SubMenu.vue'
import RemoveModal from '../Menu/MenuRemoveModal.vue'
import AttachmentModal from './Attachment/AttachmentModal.vue'
import { goPathPage, goPdfPage } from '@/common/src/utils/url'
import { PageType } from '@/utils/report'

import Cite from '@/components/Paper/Quote/cite.vue'
import { TableHeadDetail } from '@/stores/library/list'

const { Item, Divider } = Menu

// ========== 全局刷新状态管理 ==========
// 待刷新的文档ID数组（所有 Row 组件共享）
const pendingReparseDocIds = ref<string[]>([])

// 全局轮询定时器（所有 Row 组件共享）
const globalPollingTimer = ref<ReturnType<typeof setTimeout> | null>(null)

// 正在进行 API 请求的标志（所有 Row 组件共享）
const isFetching = ref(false)

// 提供给子组件使用
provide('pendingReparseDocIds', pendingReparseDocIds)
provide('globalPollingTimer', globalPollingTimer)
provide('isFetching', isFetching)

// 组件卸载时清空
onUnmounted(() => {
  pendingReparseDocIds.value = []
  if (globalPollingTimer.value) {
    clearTimeout(globalPollingTimer.value)
    globalPollingTimer.value = null
  }
  isFetching.value = false
})

export interface UserCenterSearchListData {
  paperId: string
  pdfId: string
  sort: number
  type: number
  remark: string
  docName: string
  classifyInfos: UserDocClassifyInfo[]
  id: string
  content: string
  searchRemark: string
  title: string
}

const props = defineProps({
  paperHeadVisibleList: {
    type: Array as () => Required<TableHeadDetail>[],
    required: true,
  },
  list: {
    type: Array as () => Required<UserDocInfo>[],
    default: () => [],
  },
  collating: Boolean,
})

const emit = defineEmits([
  'addButtonClick',
  'refreshClassiyAuthorVenuePaperList',
])

const minWidth = computed(() => {
  const checkboxWidth = props.collating ? 20 : 0
  const configWidth = 14
  return `${props.paperHeadVisibleList.reduce(
    (sum, item) => sum + item.width,
    checkboxWidth + configWidth
  )}px`
})

const selectedKey = ref('')
const selectedNodeFamily = computed(() => {
  if (
    !selectedKey.value ||
    selectedKey.value === allKey ||
    selectedKey.value === seperatorKey
  ) {
    return null
  }

  return getLiteratureNodeFamily(selectedKey.value)
})
const selectedPaper = computed(() => {
  return props.list.find(
    (x) => x.docId === selectedNodeFamily.value?.node.docId
  )
})
const openKeys = ref<string[]>([])
const handleOpen = (k: string, v: boolean) => (openKeys.value = v ? [k] : [])

const storeLibraryIndex = useLibraryIndex()

const attachmentDocId = ref('')
const toggleAttachment = (id?: string) => {
  attachmentDocId.value =
    id ?? (attachmentDocId.value ? '' : attachmentDocId.value)
}
const mutateDocItem = (id: string, l: number) => {
  const item = props.list.find((x) => x.docId === id)
  if (item) {
    item.hasAttachment = l > 0
  }
}

// 刷新当前列表
const refresh = () => {
  storeLibraryIndex.fetchLibraryIndex()
  emit('refreshClassiyAuthorVenuePaperList')
}

const openPage = (item: UserDocInfo) => {
  if (item.pdfId && item.pdfId !== '0') {
    goPdfPage({ pdfId: item.pdfId })
  } else if (item.paperId && item.paperId !== '0') {
    goPathPage(`/paper/${item.paperId}`)
  }
}

const removeModalVisible = ref(false)

const storeLibraryMenu = useLibraryMenu()

const { moveTo, copyTo, getLiteratureNodeFamily } = storeLibraryMenu

const computedList = computed(() => {
  const prefix =
    storeLibraryIndex.breadCrumbList[
      storeLibraryIndex.breadCrumbList.length - 1
    ]
  const keyList = Object.keys(storeLibraryIndex.libraryIndexExtra)

  const resList = props.list.map((data) => {
    let key = ''

    if (prefix.folderId === '0') {
      key = keyList.find((keyItem: string) =>
        keyItem.includes(data.docId)
      ) as string
    } else {
      key = keyList.find(
        (keyItem: string) =>
          keyItem.includes(data.docId) && keyItem.includes(prefix.folderId)
      ) as string
    }

    ;(data as any).isCopy =
      storeLibraryIndex.libraryIndexExtra[key] &&
      storeLibraryIndex.libraryIndexExtra[key].path.length >= 2
    ;(data as any).eventKey = key

    return data
  })

  return resList
})

const onVisibleChange = (visible: boolean, data: any) => {
  if (visible) {
    selectedKey.value = data.eventKey
  }
}

const addButtonClick = () => {
  emit('addButtonClick')
}

const detectDisable = (key: string) => {
  return key === selectedKey.value
}

const moveOrCopy = (type: 1 | 2, { key }: { key: string }) => {
  if (type === 1) {
    moveTo(key, selectedNodeFamily.value as LiteratureNodeFamily)
  } else {
    copyTo(key, selectedNodeFamily.value as LiteratureNodeFamily)
  }
}
</script>

<style lang="less" scoped>
:deep(.ant-dropdown-menu) {
  background-color: var(--site-theme-background-primary);
  color: var(--site-theme-text-primary);
}
</style>
