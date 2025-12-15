<template>
  <div>
    <Modal
      :visible="dialogVisible"
      :footer="null"
      width="696px"
      :destroy-on-close="true"
      :z-index="9"
      @cancel="closeModal"
    >
      <template #title>
        <div class="file-uploader-title">
          {{
            fromCopilot
              ? $t('home.upload.aiReadingTitle')
              : $t('home.upload.title')
          }}
        </div>
      </template>
      <!-- 注释搜索文件功能 -->
      <!-- <div class="file-uploader-search file-uploader-section">
        <div v-if="selectedKey" class="file-uploader-catalogue">
          {{ $t('home.library.addTo')
          }}<span style="color: #1f71e0">{{ catalogue }}</span>
        </div>
        <div class="title">
          {{
            fromCopilot
              ? $t('home.upload.aiSearchTitle')
              : $t('home.upload.searchTitle')
          }}
        </div>
        <Search
          v-if="visible"
          :from-copilot="fromCopilot"
          @onSearch="onSearch"
          @addSuccess="handleAddSuccess"
        >
          <template #addButton="{ paperId, isCollected }">
            <slot
              name="addButton"
              :paperId="paperId"
              :isCollected="isCollected"
              @addSuccess="handleAddSuccess"
            />
          </template>
        </Search>
      </div> -->
      <div v-show="uploaderVisible" class="flex space-x-2 flex-1">
        <div
          class="file-uploader-drag file-uploader-section flex-1 flex flex-col"
        >
          <!-- URL上传（Tauri 环境下隐藏） -->
          <template v-if="!isTauriEnv">
            <div class="title">{{ $t('home.upload.urlTitle') }}</div>
            <UrlUpload
              :upload-params="{ groupId, classifyId, folderId }"
              :disabled="dragDisabled"
              @addSuccess="handleAddSuccess"
              @closeUploadModal="closeModal"
            />
          </template>

          <!-- 文件拖拽上传区域 -->
          <div class="title">{{ $t('home.upload.fileTitle') }}</div>
          <Drag
            :upload-params="{ groupId, classifyId, folderId }"
            :need-protocol="needProtocol"
            :disabled="dragDisabled"
            :report-params="collectLimitDialogReportParams"
            class="flex-1"
            @beforeUpload="handleBeforeUpload"
            @addSuccess="handleAddSuccess"
            @closeUploadModal="closeModal"
          />
        </div>
        <slot name="latest" />
      </div>
    </Modal>
    <Result
      :from-copilot="fromCopilot"
      is-need-report-reading
      @conflictResolved="handleAddSuccess"
      @refreshList="handleRefreshList"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, PropType, computed } from 'vue'
import { Modal, message } from 'ant-design-vue'
import { BreadCrumb } from '../../../stores/library'
import Search from './Search.vue'
import Drag from './Drag.vue'
import Result from './Result.vue'
import UrlUpload from './UrlUpload.vue'
import { LimitDialogReportParams } from '../helper'
import { isInTauri } from '@/util/env'

const props = defineProps({
  fromCopilot: {
    type: Boolean,
    default: false,
  },
  visible: {
    type: Boolean,
    default: false,
  },
  groupId: {
    type: String,
    default: '',
  },
  classifyId: {
    type: String,
    default: '',
  },
  selectedKey: {
    type: String,
    default: '',
  },
  folderId: {
    type: String,
    default: '',
  },
  needProtocol: Boolean,
  dragDisabled: Boolean,
  breadCrumbList: {
    type: Array as PropType<BreadCrumb[]>,
    default: () => [],
  },
  collectLimitDialogReportParams: {
    type: Object as PropType<LimitDialogReportParams>,
    default: null,
  },
})

const emit = defineEmits(['close', 'addSuccess', 'beforeUpload', 'refreshList'])

const dialogVisible = ref<boolean>(false)
const uploaderVisible = ref<boolean>(true)
const catalogue = ref('' as string)

watch(
  () => props.visible,
  (newValue) => {
    dialogVisible.value = newValue
  }
)
watch(
  () => props.breadCrumbList,
  (newValue) => {
    catalogue.value = ''
    newValue.forEach((item, index) => {
      if (item.docName.length > 10) {
        item.docName = item.docName.slice(0, 9) + '...'
      }
      if (index < newValue.length - 1) {
        catalogue.value += item.docName + ' / '
      } else {
        catalogue.value += item.docName
      }
    })
  }
)
const closeModal = () => {
  dialogVisible.value = false
  uploaderVisible.value = true
  emit('close')
}
const onSearch = (isSearching: boolean) => {
  uploaderVisible.value = !isSearching
}
const handleAddSuccess = () => {
  emit('addSuccess')
}

const handleRefreshList = () => {
  emit('refreshList')
}

const handleBeforeUpload = () => {
  emit('beforeUpload')
}

// 检测是否为 Tauri 环境
const isTauriEnv = computed(() => isInTauri())

</script>
<style lang="less" scoped>
.file-uploader {
  &-section {
    .title {
      font-weight: 600;
      color: #262625;
      line-height: 22px;
      margin-top: 4px;
      margin-bottom: 12px;
    }
  }
  &-title {
    font-weight: 600;
    color: #262625;
    line-height: 32px;
  }
  &-catalogue {
    font-size: 14px;
    font-weight: 400;
    color: #73716f;
    line-height: 24px;
    padding-bottom: 5px;
    span {
      margin-left: 3px;
    }
  }
}
</style>
<style>
.z-index-1061 {
  z-index: 1061;
}
</style>
