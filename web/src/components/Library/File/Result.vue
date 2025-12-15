<template>
  <div
    v-show="
      uploadStore.uploadList.length &&
      uploaderResult.status !== FileUploaderResult.unstart
    "
    class="literature-file-uploader-result js-literature-drag-result"
  >
    <div class="header">
      <div
        v-if="uploaderResult.status === FileUploaderResult.doing"
        class="title"
      >
        <LoadingOutlined />{{ $t('home.upload.result.uploadingTitle') }}
        {{ uploaderResult.uploading }}/{{ uploadStore.uploadList.length }}
      </div>
      <div v-else class="title">
        <CheckCircleOutlined style="color: #1fe02f" />{{
          $t('home.upload.result.uploadedTitle')
        }}
      </div>
      <DownOutlined v-if="mode === 'down'" @click="mode = 'up'" />
      <UpOutlined v-else @click="mode = 'down'" />
      <CloseOutlined
        style="cursor: pointer"
        @click="handleClose"
      />
    </div>
    <div v-show="mode === 'down'" class="atable">
      <div class="tip">
        {{ $t('home.upload.result.tip') }}
      </div>
      <div class="header">
        <div class="filename">{{ $t('home.upload.result.fileName') }}</div>
        <div class="progress">{{ $t('home.upload.result.progressTitle') }}</div>
        <div class="action">{{ $t('home.upload.result.opt') }}</div>
      </div>
      <div
        v-for="item in uploadStore.uploadList"
        :key="item.extra.uid"
        class="item"
      >
        <div class="filename">{{ item.extra.name }}</div>
        <div class="progress">
          <!-- 使用新的状态系统 -->
          <div
            v-if="isNewErrorStatus(item.status)"
            class="error"
          >
            {{ getStatusDisplayText(item.status) }}
          </div>
          <div
            v-else-if="isSuccessStatus(item.status)"
            class="finish"
          >
            {{ getStatusDisplayText(item.status) }} ({{ item.progress }}%)
          </div>
          <div
            v-else-if="isProcessingStatus(item.status)"
            :class="getStatusClassName(item.status)"
          >
            {{ getStatusDisplayText(item.status) }} ({{ item.progress }}%)
          </div>
        </div>
        <div class="action">
          <Button
            v-if="isSuccessStatus(item.status)"
            type="link"
            @click="goReading(item.docInfo || item.groupDocInfo)"
          >
            {{
              fromCopilot
                ? $t('home.upload.result.optButton.startAiReading')
                : $t('home.upload.result.optButton.startReading')
            }}
          </Button>
          <Button
            v-else-if="isErrorStatus(item.status)"
            type="link"
            @click="showReason(item)"
          >
            {{ $t('home.upload.result.optButton.reason') }}
          </Button>
        </div>
      </div>
    </div>
    <Modal 
      v-model:visible="dialogVisible" 
      :footer="null" 
      width="408px"
      :zIndex="10001"
      :maskClosable="true"
      @cancel="closeErrorDialog"
    >
      <template #title>
        <div class="literature-uploader-fail-title">
          <ExclamationCircleFilled class="exclamation" />{{
            $t('home.upload.result.progressDetail.failed')
          }}
        </div>
      </template>
      <div v-if="currentFailFile" class="literature-uploader-fail">
        <div class="msg">
          {{ $t('home.upload.result.fileName') }}{{ $t('home.global.colon')
          }}{{ currentFailFile.extra?.name || '' }}
        </div>
        <div class="msg" v-if="currentFailFile.tokens?.queryStatusToken">
          {{ $t('home.upload.errorDetail.log') }}{{ $t('home.global.colon')
          }}{{ currentFailFile.tokens.queryStatusToken }}
        </div>
        <div class="msg">
          {{ $t('home.upload.errorDetail.reason') }}{{ $t('home.global.colon')
          }}{{ currentFailFile.message || 'unknown error' }}
        </div>
        <!-- <div class="qrcode">
          <span>{{ $t('home.upload.errorDetail.tip') }}</span>
          <img src="https://readpaper.com/doc/assets/img/Qrcode.png" alt="" />
        </div> -->
        <a :href="userGuideLink" target="_blank" class="user-guide-link">{{ $t('home.upload.errorDetail.guide') }}</a>
      </div>
    </Modal>
  </div>
</template>
<script lang="ts" setup>

import { computed, onMounted, ref, h } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  MyCollectedDocInfo,
} from 'go-sea-proto/gen/ts/doc/UserDoc'
import {
  UserDocParsedStatusEnum,
} from 'go-sea-proto/gen/ts/doc/UserDocParsedStatus'
import interact from 'interactjs'
import { message, Modal, Button } from 'ant-design-vue'
import { Item } from '@/utils/pdf-upload/index.js'
import { 
  getStatusText, 
  getStatusProgress, 
  isErrorStatus as isNewErrorStatus,
  isProcessingStatus,
  getStatusClassName,
  isSuccessStatus,
  isErrorStatus
} from '@/utils/pdf-upload/statusMapper.js'
import { goPdfPage } from '@/common/src/utils/url'
import { GroupDocInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/groups/GroupDoc'
import {
  CheckCircleOutlined,
  CloseOutlined,
  DownOutlined,
  ExclamationCircleFilled,
  LoadingOutlined,
  UpOutlined,
  InfoCircleOutlined,
} from '@ant-design/icons-vue'
import { PageType, reportElementClick } from '@/utils/report'
import { useUpload } from '@/stores/upload'
import { PAGE_ROUTE_NAME } from '@/routes/type'
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '@/hooks/useLanguage';

enum FileUploaderResult {
  unstart,
  doing,
  done,
}

const props = defineProps({
  fromCopilot: {
    type: Boolean,
    default: false,
  },
  isNeedReportReading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['conflictResolved', 'refreshList'])

const uploadStore = useUpload()
const { locale, t } = useI18n()

// 语言管理
const { isCurrentLanguage } = useLanguage();

// 创建用户指南链接的计算属性
const userGuideLink = computed(() => {
  if (isCurrentLanguage(Language.ZH_CN)) {
    return '/docs/zh/guide'
  }
  return '/docs/guide'
})

const uploaderResult = computed<{
  status: FileUploaderResult
  uploading: number
}>(() => {
  if (!uploadStore.uploadList.length) {
    return {
      status: FileUploaderResult.unstart,
      uploading: 0,
    }
  }
  
  // 检查是否有正在上传的文件
  const uploading = uploadStore.uploadList.filter(
    (item) =>
      isProcessingStatus(item.status)
  )
  
  if (uploading.length) {
    return {
      status: FileUploaderResult.doing,
      uploading: uploading.length,
    }
  }
  
  return {
    status: FileUploaderResult.done,
    uploading: 0,
  }
})

const mode = ref<'down' | 'up'>('down')

// 获取状态显示文本的辅助函数
const getStatusDisplayText = (status: number) => {
  return getStatusText(status, t);
};

const canCloseUploader = () => {
  // 如果所有上传项都已完成或失败，允许关闭
  const result = uploadStore.uploadList.every(item => {
    // 检查是否是完成或失败状态
    const isFinished = isSuccessStatus(item.status) || isNewErrorStatus(item.status);
    
    if (!isFinished) {
      return false;
    }
    
    return true;
  });
  
  return result;
}

const handleClose = () => {
  // 无条件关闭上传器
  uploadStore.clearUploadList();
  
  // 不触发刷新列表事件，因为已经在上传完成时触发过了
  // window.dispatchEvent(new CustomEvent('uploadFinished'));
}

const handleDragEvent = () => {}

onMounted(() => {
  // target elements with the "draggable" class
  interact('.js-literature-drag-result').draggable({
    // enable inertial throwing
    inertia: true,
    // keep the element within the area of it's parent
    modifiers: [
      interact.modifiers.restrictRect({
        restriction: 'body',
        endOnly: true,
      }),
    ],
    // enable autoScroll
    autoScroll: true,

    listeners: {
      // call this function on every dragmove event
      move(event) {
        const target = event.target
        // keep the dragged position in the data-x/data-y attributes
        const x = (parseFloat(target.getAttribute('data-x')) || 0) + event.dx
        const y = (parseFloat(target.getAttribute('data-y')) || 0) + event.dy

        // translate the element
        target.style.transform = 'translate(' + x + 'px, ' + y + 'px)'

        // update the posiion attributes
        target.setAttribute('data-x', x)
        target.setAttribute('data-y', y)
      },
    },
  })

  // 监听拖拽事件
  document.body.addEventListener('dragenter', handleDragEvent, false)
})

const dialogVisible = ref<boolean>(false)
const currentFailFile = ref<Item>()

const showReason = (fileInfo: Item) => {
  currentFailFile.value = fileInfo
  dialogVisible.value = true
}

const closeErrorDialog = () => {
  dialogVisible.value = false
  currentFailFile.value = undefined
}

const goReading = (docInfo: MyCollectedDocInfo | GroupDocInfo | undefined) => {
  if (docInfo && docInfo.pdfId) {
    if (props.isNeedReportReading) {
      reportElementClick({
        page_type: PageType[PAGE_ROUTE_NAME.WORKBENCH],
        type_parameter: 'none',
        element_name: 'ai_assist_popup_upload',
        element_parameter: docInfo.pdfId,
        status: 'none',
      })
    }
    if (props.fromCopilot) {
      goPdfPage({ pdfId: docInfo.pdfId, tab: 'copilot' })
    } else {
      goPdfPage({ pdfId: docInfo.pdfId })
    }
  } else {
    message.warn('文献信息为空')
  }
}

const conflictDialogVisible = ref<boolean>(false)

const showConflict = (item: Item) => {
  currentFailFile.value = item
  conflictDialogVisible.value = true
}

const handleConflict = async (useOrigin: boolean, $event?: any) => {
  if ($event && $event.target.localName !== 'button') {
    return
  }

  await uploadStore.resolveUploadConflict(currentFailFile.value!, useOrigin)
  conflictDialogVisible.value = false
  emit('conflictResolved')
  emit('refreshList')
}
const goCancel = (item: Item) => {
  Modal.confirm({
    title: '确认取消？',
    content: '',
    icon: h(InfoCircleOutlined),
    okText: '是',
    okType: 'danger',
    cancelText: '否',
    onOk: () => {
      uploadStore.cancelUpload(item)
    },
    onCancel() {},
  })
}
</script>
<style lang="less" scoped>
.literature-file-uploader-result {
  position: fixed;
  right: 16px;
  bottom: 32px;
  background: #fff;
  width: 480px;
  z-index: 10000;
  box-shadow: 0 4px 24px 0 rgba(0, 0, 0, 0.2);

  .tip {
    padding: 0 21px;
    color: #a6a4a1;
    font-size: 12px;
  }

  .header {
    background: #f5f7fa;
    padding: 8px 21px;
    display: flex;
    justify-content: space-between;
    .anticon {
      font-size: 16px;
      line-height: 32px;
      color: rgba(0, 0, 0, 0.85);
    }
    .title {
      font-size: 16px;
      font-weight: 600;
      color: #262625;
      line-height: 32px;
      flex: 1;
      .anticon {
        margin-right: 6px;
        cursor: pointer;
      }
    }

    .anticon + .anticon {
      margin-left: 40px;
    }
  }

  .atable {
    margin: 10px 0 30px 0;
    overflow-y: auto;
    max-height: 400px;
    .filename {
      width: 200px;
      word-break: break-all;
    }

    .header {
      display: flex;
      border-bottom: 1px solid #e8e8e8;
      font-weight: 500;
      color: #73716f;
      line-height: 22px;
      margin: 0 21px;
      padding: 9px 0;
      background: #fff;
    }

    .progress {
      flex: 1;
    }

    .action {
      width: 80px;
    }

    .item {
      display: flex;
      color: #262625;
      line-height: 18px;
      font-size: 13px;
      padding: 11px 21px;

      .progress {
        div {
          &::before {
            display: inline-block;
            width: 8px;
            height: 8px;
            border-radius: 4px;
            content: '';
            margin-right: 5px;
          }
          &.error {
            &::before {
              background: #e01f1f;
            }
          }
          &.finish {
            &::before {
              background: #1fe02f;
            }
          }
          &.repeat {
            &::before {
              background: #ad1fe0;
            }
          }
          &.waiting {
            &::before {
              background: #f48227;
            }
          }
          &.parse,
          &.match,
          &.upload {
            &::before {
              background: #1f71e0;
            }
          }
        }
      }

      .action {
        .ant-btn-link {
          color: #1f71e0;
          line-height: 18px;
          font-size: 13px;
          padding: 0;
          height: 18px;
        }
      }
    }
  }
}

.literature-uploader-fail-title {
  font-size: 16px;
  line-height: 32px;
  padding-left: 4px;
  .exclamation {
    color: #e01f1f;
    font-size: 18px;
    margin-right: 6px;
  }
}
.literature-uploader-fail {
  margin: -4px 4px;
  color: #262625;
  line-height: 22px;

  .msg {
    margin-bottom: 16px;
    word-break: break-all;
  }

  .qrcode {
    text-align: center;
    margin-top: 8px;
    margin-bottom: 14px;
    span {
      color: #73716f;
    }
    img {
      width: 184px;
      height: 184px;
      margin-top: 12px;
    }
  }
  .user-guide-link {
    display: block;
    text-align: center;
    margin-top: 8px;
    margin-bottom: 14px;
    color: #1f71e0;
    text-decoration: none;
    
    &:hover {
      color: #0d5cbf;
      text-decoration: underline;
    }
  }
}
</style>
