<template>
  <div class="url-uploader-wrap">
    <div class="url-uploader-content">
      <!-- 输入框和按钮两个元素一行排列 -->
      <div class="url-input-section">
        <div class="input-row">
          <a-input
            v-model:value="urlValue"
            :placeholder="$t('home.upload.urlPlaceholder')"
            :status="urlError ? 'error' : ''"
            class="url-input"
            @blur="validateUrl"
            @input="clearError"
          />
          <a-button
            type="primary"
            :disabled="!urlValue.trim() || !!urlError"
            class="submit-button"
            @click="handleSubmit"
          >
            <LinkOutlined />
            {{ $t('home.upload.urlSubmit') }}
          </a-button>
        </div>
        <div v-if="urlError" class="error-message">
          {{ urlError }}
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import { Input, Button, Radio, RadioGroup, message } from 'ant-design-vue'
import { LinkOutlined } from '@ant-design/icons-vue'
import { useI18n } from 'vue-i18n'
import { getParseToken, uploadPdfByUrlLink, getUserDocCreateStatus } from '../../../api/document'
import { useUpload } from '../../../stores/upload'
import { createItem } from '../../../utils/pdf-upload/index.js'
import { UserDocParsedStatusEnum } from 'go-sea-proto/gen/ts/doc/UserDocParsedStatus'
import { getStatusProgress, isErrorStatus } from '../../../utils/pdf-upload/statusMapper.js'

// 使用 Ant Design Vue 组件别名
const AInput = Input
const AButton = Button
const ARadio = Radio
const ARadioGroup = RadioGroup

const props = defineProps({
  uploadParams: {
    type: Object,
    default: () => ({ groupId: '', classifyId: '', folderId: '' }),
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['addSuccess', 'closeUploadModal'])

const uploadStore = useUpload()
const { t } = useI18n()

// 响应式数据
const urlValue = ref('')
const urlError = ref('')

// URL验证函数
const isValidUrl = (url: string): boolean => {
  try {
    const urlObj = new URL(url)
    // 只允许 http 和 https 协议
    return urlObj.protocol === 'http:' || urlObj.protocol === 'https:'
  } catch {
    return false
  }
}

// 验证URL
const validateUrl = () => {
  const trimmedUrl = urlValue.value.trim()
  if (trimmedUrl && !isValidUrl(trimmedUrl)) {
    urlError.value = t('home.upload.urlInvalid')
  } else {
    urlError.value = ''
  }
}

// 清除错误
const clearError = () => {
  if (urlError.value) {
    urlError.value = ''
  }
}

// URL上传状态轮询函数 - 复用普通上传的轮询逻辑
const startUrlPollingStatus = async (item, token, onUpdate) => {
  try {
    const pollInterval = 2000 // 轮询间隔2秒
    
    const updateStatus = (status, progress) => {
      item.status = status
      item.progress = progress
      onUpdate()
    }
    
    const continuePolling = () => {
      setTimeout(pollStatus, pollInterval)
    }
    
    const pollStatus = async () => {
      try {
        const statusResponse = await getUserDocCreateStatus({ token })
        console.log('URL上传状态轮询:', statusResponse)
        
        const createStatus = statusResponse.data?.status
        
        // 如果状态为空则继续轮询
        if (createStatus === 0) {
          continuePolling()
          return
        }
        
        // 如果状态没有变化，继续轮询
        if (createStatus === item.status) {
          continuePolling()
          return
        }
        
        // 检查是否完成
        if (createStatus >= UserDocParsedStatusEnum.HEADER_DATA_PARSED) {
          updateStatus(createStatus, getStatusProgress(createStatus))
          // 设置文档信息
          if (statusResponse.data && statusResponse.data.docInfo) {
            item.docInfo = statusResponse.data.docInfo
          }
          
          // 添加延时让用户看到最终状态
          await new Promise((resolve) => setTimeout(resolve, 2000))
          
          // 触发上传完成事件
          if (typeof window !== 'undefined') {
            window.dispatchEvent(new CustomEvent('uploadFinished'))
          }
          return
        } else if (isErrorStatus(createStatus)) {
          // 处理错误状态
          updateStatus(createStatus, getStatusProgress(createStatus))
          return
        } else {
          // 更新到中间状态
          updateStatus(createStatus, getStatusProgress(createStatus))
          await new Promise((resolve) => setTimeout(resolve, 1000))
          continuePolling()
        }
      } catch (error) {
        console.error('状态轮询错误:', error)
        continuePolling()
      }
    }
    
    // 开始轮询
    pollStatus()
  } catch (error) {
    console.error('启动状态轮询失败:', error)
    item.status = UserDocParsedStatusEnum.UPLOAD_FAILED
    item.progress = getStatusProgress(UserDocParsedStatusEnum.UPLOAD_FAILED)
    onUpdate()
  }
}

// 事件处理 - 参考Drag.vue的实现模式
const handleSubmit = async () => {
  const trimmedUrl = urlValue.value.trim()
  if (!trimmedUrl) return
  
  // 验证URL格式
  if (!isValidUrl(trimmedUrl)) {
    urlError.value = t('home.upload.urlInvalid')
    return
  }
  
  try {
    console.log('URL上传开始:', urlValue.value)
    
    // 1. 调用getParseToken获取解析令牌
    const tokenResponse = await getParseToken()
    if (tokenResponse.status !== 1 || !tokenResponse.data) {
      message.error('获取解析令牌失败，请稍后再试')
      return
    }
    
    const parseToken = tokenResponse.data.token
    console.log('获取到解析令牌:', parseToken)
    
    // 2. 调用uploadPdfByUrlLink上传PDF文件
    const uploadRequest = {
      url: urlValue.value.trim(),
      uploadToken: parseToken,
      groupId: props.uploadParams.groupId,
      classifyId: props.uploadParams.classifyId,
      folderId: props.uploadParams.folderId,
    }
    
    const uploadResponse = await uploadPdfByUrlLink(uploadRequest)
    if (uploadResponse.status === 1) {
      // URL上传成功，直接创建上传项并使用parseToken进行状态轮询
      const mockFile = new File([''], `URL: ${urlValue.value}`, { type: 'application/pdf' })
      
      // 创建上传项
      const item = createItem(mockFile, {
        name: `URL: ${urlValue.value}`,
        uid: `url-upload-${Date.now()}`,
        size: 0,
        type: 'application/pdf',
      })
      
      // 设置初始状态为DOWNLOADING
      item.status = UserDocParsedStatusEnum.DOWNLOADING
      item.progress = getStatusProgress(UserDocParsedStatusEnum.DOWNLOADING)
      
      // 添加到上传列表
      uploadStore.uploadList.push(item)
      
      // 触发更新
      const triggerUpdate = () => {
        uploadStore.uploadList = [...uploadStore.uploadList]
      }
      
      // 通知父组件关闭上传对话框
      emit('closeUploadModal')
      
      // 使用parseToken开始状态轮询
      startUrlPollingStatus(item, parseToken, triggerUpdate)
      
      emit('addSuccess')
      //message.success('URL上传请求已提交，正在处理中...')
      // 清空输入框
      urlValue.value = ''
    } else {
      message.error('URL上传失败，请检查链接是否有效')
    }
    
  } catch (error) {
    console.error('URL上传失败:', error)
    message.error('URL上传失败，请稍后再试')
  }
}
</script>

<style lang="less" scoped>
.url-uploader-wrap {
  background: #fff;
  border: 1px dashed #e4e7ed;
  border-radius: 6px;
  padding: 24px;
  margin-bottom: 16px;
}

.url-uploader-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.url-input-section {
  .input-row {
    display: flex;
    gap: 12px;
    align-items: center;
  }
  
  .error-message {
    color: #ff4d4f;
    font-size: 12px;
    margin-top: 4px;
    line-height: 1.5;
  }
  
  .url-input {
    flex: 1;
    
    :deep(.ant-input) {
      border-color: #e4e7ed;
      
      &:hover {
        border-color: #1f71e0;
      }
      
      &:focus {
        border-color: #1f71e0;
        box-shadow: 0 0 0 2px rgba(31, 113, 224, 0.1);
      }
    }
    
    :deep(.ant-input-status-error) {
      border-color: #ff4d4f;
      
      &:hover {
        border-color: #ff7875;
      }
      
      &:focus {
        border-color: #ff4d4f;
        box-shadow: 0 0 0 2px rgba(255, 77, 79, 0.1);
      }
    }
  }
  
  .submit-button {
    background-color: #1f71e0;
    border-color: #1f71e0;
    flex-shrink: 0;
    
    &:hover:not(:disabled) {
      background-color: #1557c7;
      border-color: #1557c7;
    }
    
    &:disabled {
      background-color: #f5f5f5;
      border-color: #d9d9d9;
      color: rgba(0, 0, 0, 0.25);
    }
  }
}
</style>
