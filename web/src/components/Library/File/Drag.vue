<template>
  <div 
    class="literature-uploader-wrap"
    @dragenter="onNativeDragEnter"
    @dragover="onNativeDragOver"
    @dragleave="onNativeDragLeave"
    @drop="onNativeDrop"
  >
    <!-- Tauri 环境：使用原生文件选择对话框 -->
    <template v-if="isTauriEnv">
      <div class="literature-uploader-drag tauri-uploader">
        <a-button type="primary" :disabled="disabled || !userStore.isLogin()" @click="handleTauriFileSelect">
          <UploadOutlined />{{ $t('home.upload.upload') }}
        </a-button>
        <p class="hint">
          {{ $t('home.upload.uploadTip1') }}<span style="text-decoration: underline" @click="handleTauriFolderSelect">{{ $t('home.upload.uploadTip2') }}</span>
        </p>
      </div>
    </template>
    
    <!-- 浏览器环境：使用 Ant Design Upload 组件 -->
    <template v-else>
      <UploadDragger
        name="file1"
        :multiple="true"
        :custom-request="customRequest"
        :show-upload-list="false"
        :before-upload="beforeUpload"
        :directory="false"
        :accept="'.pdf'"
        :disabled="disabled || !userStore.isLogin()"
        style="pointer-events: none;"
      >
        <div style="position: relative; pointer-events: auto;" @click="handleProtocol">
          <slot v-if="$slots.default" name="default" />
          <div v-else class="literature-uploader-drag">
            <a-button type="primary" :disabled="disabled || !userStore.isLogin()">
              <UploadOutlined />{{ $t('home.upload.upload') }}</a-button
            >
            <p class="hint">
              {{ $t('home.upload.uploadTip1')
              }}<span style="text-decoration: underline">{{
                $t('home.upload.uploadTip2')
              }}</span>
            </p>
          </div>
          <div v-show="isDragging" class="literature-uploader-dragging">
            <div class="bg">
              <ArrowUpOutlined />
              <p class="hint1">{{ $t('home.upload.dragTip1') }}</p>
              <p class="hint2">{{ $t('home.upload.dragTip2') }}</p>
            </div>
          </div>
        </div>
      </UploadDragger>
      <Upload
        class="folder"
        name="file2"
        :multiple="true"
        :custom-request="customRequest"
        :show-upload-list="false"
        :before-upload="beforeUpload"
        directory
        :accept="'.pdf'"
        :disabled="disabled || !userStore.isLogin()"
        style="pointer-events: none;"
      >
        <span class="mask" style="pointer-events: auto;" @click="handleProtocol">&nbsp;</span>
      </Upload>
    </template>
  </div>
</template>
<script lang="ts" setup>
import { onMounted, onUnmounted, PropType, ref, computed } from 'vue'
import { Upload, UploadDragger, message } from 'ant-design-vue'
import { isSuccessStatus } from '@/utils/pdf-upload/statusMapper.js'
import { EventCode } from '@/utils/report'
import { useI18n } from 'vue-i18n'
import { ArrowUpOutlined, UploadOutlined } from '@ant-design/icons-vue'
import {
  RcFile,
  UploadRequestOption,
} from 'ant-design-vue/lib/vc-upload/interface'
import { LimitDialogReportParams } from '../helper'
import reporter from '@idea/aiknowledge-report'
import { useUpload } from '@/stores/upload'
import { useUserStore } from '@/common/src/stores/user'
import { isInTauri } from '@/util/env'
import { selectPdfFiles, selectPdfFolder, type FileSelectResult } from '@/utils/fileSelector'

const userStore = useUserStore()

const props = defineProps({
  offset: {
    type: Object as PropType<{
      top: string
      bottom: string
      left: string
      right: string
    }>,
    default: () => ({ top: '0', bottom: '0', left: '0', right: '0' }),
  },
  uploadParams: {
    type: Object as PropType<{
      groupId: string
      classifyId: string
      folderId: string
    }>,
    default: () => ({ groupId: '', classifyId: '', folderId: '' }),
  },
  needProtocol: Boolean,
  disabled: Boolean,
  reportParams: {
    type: Object as PropType<LimitDialogReportParams>,
    default: null,
  },
})

const emit = defineEmits(['addSuccess', 'beforeUpload', 'isDragging', 'closeUploadModal'])

const { t } = useI18n()
const uploadStore = useUpload()

const uploadingFiles = new Map<string, boolean>();

const customRequest = (data: UploadRequestOption) => {
  emit('closeUploadModal');
  
  const fileName = (data.file as RcFile).name;
  const fileKey = `${fileName}-${(data.file as RcFile).size}-${(data.file as RcFile).lastModified}`;
  if (uploadingFiles.has(fileKey)) {
    return;
  }
  uploadingFiles.set(fileKey, true);
  
  uploadStore
    .addUpload(
      data.file as RcFile,
      {
        groupId: props.uploadParams.groupId,
        classifyId: props.uploadParams.classifyId,
        folderId: props.uploadParams.folderId,
        paperId: '',
      },
      {
        name: (data.file as RcFile).name,
        uid: (data.file as RcFile).uid,
        size: (data.file as RcFile).size,
        type: (data.file as RcFile).type,
      }
    )
    .then((fileInfo) => {      
      // 清理重复检测标记
      uploadingFiles.delete(fileKey);
      
      emit('addSuccess')
      if (!fileInfo) {
        return;
      }
      
      if (fileInfo.docInfo && isSuccessStatus(fileInfo.status)) {
        reporter.report(
          {
            event_code: EventCode.readpaperPdfUploadSuccess,
          },
          {
            pdf_id: fileInfo.docInfo.pdfId,
          }
        )
      }
      
      return fileInfo; 
    })
    .catch((error) => {
      
      // 清理重复检测标记
      uploadingFiles.delete(fileKey);
      
      console.error('Upload failed:', error)
      message.error(t('home.upload.uploadFailed') || '上传失败，请稍后再试')
    })
}

const beforeUpload = (file: File) => {
  if (props.needProtocol) {
    emit('beforeUpload')
    return false
  }
  if (file.type === 'application/pdf' || file.name.endsWith('.pdf')) {
    emit('beforeUpload')
    return true
  }

  message.warning(t('home.upload.drag.uploadWarning', { name: file.name }))
  return false
}

const handleProtocol = (e: MouseEvent) => {
  if (props.disabled) {
    message.info('功能升级中，将尽快上线，敬请期待')
    return
  }
  if (props.needProtocol) {
    emit('beforeUpload')
    e.stopPropagation()
  }
}

// 检测是否为 Tauri 环境
const isTauriEnv = computed(() => isInTauri())

/**
 * Tauri 环境下的文件选择处理
 */
const handleTauriFileSelect = async () => {
  if (props.disabled || !userStore.isLogin()) {
    return
  }
  
  if (props.needProtocol) {
    emit('beforeUpload')
    return
  }
  
  try {
    const results = await selectPdfFiles()
    await processFileSelectResults(results)
  } catch (error) {
    console.error('Tauri file selection failed:', error)
    message.error(t('home.upload.uploadFailed') || '文件选择失败')
  }
}

/**
 * Tauri 环境下的文件夹选择处理
 */
const handleTauriFolderSelect = async () => {
  if (props.disabled || !userStore.isLogin()) {
    return
  }
  
  if (props.needProtocol) {
    emit('beforeUpload')
    return
  }
  
  try {
    const results = await selectPdfFolder()
    await processFileSelectResults(results)
  } catch (error) {
    console.error('Tauri folder selection failed:', error)
    message.error(t('home.upload.uploadFailed') || '文件夹选择失败')
  }
}

/**
 * 处理文件选择结果
 */
const processFileSelectResults = async (results: FileSelectResult[]) => {
  if (results.length === 0) {
    return
  }
  
  emit('closeUploadModal')
  
  for (const result of results) {
    if (result.isLocalPathMode && result.localPath) {
      // Tauri 模式：使用本地路径上传
      await uploadWithLocalPath(result.localPath, result.fileName)
    } else if (result.file) {
      // 浏览器模式：使用 File 对象上传
      const uploadOption = {
        file: result.file,
        filename: result.fileName,
        onProgress: () => {},
        onSuccess: () => {},
        onError: () => {}
      }
      customRequest(uploadOption as any)
    }
  }
}

/**
 * 使用本地路径上传文件（Tauri 模式）
 * TODO: 后续需要根据后端接口返回的标识决定是传路径还是上传文件
 */
const uploadWithLocalPath = async (localPath: string, fileName: string) => {
  console.log('Tauri local path upload:', { localPath, fileName })
  
  // 目前先读取文件内容，转换为 File 对象进行上传
  // 后续可以根据后端接口返回的标识决定是直接传路径还是上传文件
  try {
    const { readFile } = await import('@tauri-apps/plugin-fs')
    const fileData = await readFile(localPath)
    const file = new File([fileData], fileName, { type: 'application/pdf' })
    
    const fileKey = `${fileName}-${file.size}-${file.lastModified}`;
    if (uploadingFiles.has(fileKey)) {
      return;
    }
    uploadingFiles.set(fileKey, true);
    
    // 调用 uploadStore.addUpload，传递 localPath 到 extra 参数
    uploadStore
      .addUpload(
        file,
        {
          groupId: props.uploadParams.groupId,
          classifyId: props.uploadParams.classifyId,
          folderId: props.uploadParams.folderId,
          paperId: '',
        },
        {
          name: fileName,
          uid: `tauri-${Date.now()}`,
          size: file.size,
          type: file.type,
          localPath: localPath, // 传递本地路径
        }
      )
      .then((fileInfo) => {
        uploadingFiles.delete(fileKey);
        emit('addSuccess')
        if (!fileInfo) {
          return;
        }
        if (fileInfo.docInfo && isSuccessStatus(fileInfo.status)) {
          reporter.report(
            {
              event_code: EventCode.readpaperPdfUploadSuccess,
            },
            {
              pdf_id: fileInfo.docInfo.pdfId,
            }
          )
        }
        return fileInfo;
      })
      .catch((error) => {
        uploadingFiles.delete(fileKey);
        console.error('Upload failed:', error)
        message.error(t('home.upload.uploadFailed') || '上传失败，请稍后再试')
      })
  } catch (error) {
    console.error('Failed to read local file:', error)
    message.error(t('home.upload.uploadFailed') || '读取本地文件失败')
  }
}

const isDragging = ref<boolean>(false)
let enterTarget: EventTarget | null = null

const changDragzoneStyle = (e: DragEvent) => {
  const citationStyleDragClassName = ['personal-style-item-wrap', 'style-name']
  const isPersonalCitationStyleDrag = citationStyleDragClassName.findIndex(
    (item) => {
      const className =
        (e as any)?.target?.className ||
        (e as any)?.fromElement?.className ||
        (e as any)?.toElement?.className
      return className.includes(item)
    }
  )
  if (isPersonalCitationStyleDrag !== -1) {
    return
  }
  if (e.type === 'dragenter') {
    enterTarget = e.target
  }
  if (
    e.type === 'dragenter' ||
    e.type === 'dragover' ||
    (e.type === 'dragleave' && e.target !== enterTarget)
  ) {
    isDragging.value = true
  } else {
    isDragging.value = false
  }
  emit('isDragging', isDragging.value)
}

const bindDragEvent = (add: boolean) => {
  if (props.disabled) {
    return
  }
  
  ;['dragenter', 'dragleave', 'drop', 'dragover'].forEach((event) => {
    // 绑定到window
    window[add ? 'addEventListener' : 'removeEventListener'](
      event,
      changDragzoneStyle as any,
      false
    );
  })
}

onMounted(() => {
  bindDragEvent(true)
})

onUnmounted(() => {
  bindDragEvent(false)
})

const onNativeDragEnter = (e: DragEvent) => {
  if (e.dataTransfer?.items) {
    for (let i = 0; i < e.dataTransfer.items.length; i++) {
      const item = e.dataTransfer.items[i];
      if (item.webkitGetAsEntry) {
      }
    }
  }
  
  e.preventDefault(); 
}

const onNativeDragOver = (e: DragEvent) => {
  e.preventDefault();
}

const onNativeDragLeave = (e: DragEvent) => {
}

const onNativeDrop = (e: DragEvent) => {
  
  // 只使用 DataTransferItems 处理，避免重复处理
  if (e.dataTransfer?.items) {
    const items = Array.from(e.dataTransfer.items);
    const promises: Promise<void>[] = [];
    
    items.forEach((item, index) => {
      if (item.kind === 'file') {
        const entry = item.webkitGetAsEntry();
        
        if (entry) {
          if (entry.isDirectory) {
            promises.push(processDirectory(entry as any));
          } else if (entry.isFile) {
            promises.push(processFile(entry as any));
          }
        }
      }
    });
    
    if (promises.length > 0) {
      Promise.all(promises).then(() => {
      }).catch((error) => {
        console.error('Error processing items:', error);
      });
      
      e.preventDefault();
      e.stopPropagation();
      return;
    }
  }
  
  // 如果没有 DataTransferItems 支持，让 Ant Design 的 Upload 组件处理
  e.preventDefault();
}

const processDirectory = async (dirEntry: any): Promise<void> => {
  return new Promise((resolve, reject) => {
    const dirReader = dirEntry.createReader();
    
    const readEntries = () => {
      dirReader.readEntries((entries: any[]) => {
        if (entries.length === 0) {
          resolve();
          return;
        }
        
        const promises = entries.map(entry => {
          if (entry.isDirectory) {
            return processDirectory(entry);
          } else if (entry.isFile) {
            return processFile(entry);
          }
          return Promise.resolve();
        });
        
        Promise.all(promises).then(() => {
          readEntries(); 
        }).catch(reject);
      }, reject);
    };
    
    readEntries();
  });
};

const processFile = async (fileEntry: any): Promise<void> => {
  return new Promise((resolve, reject) => {
    fileEntry.file((file: File) => {
      if (file.type === 'application/pdf' || file.name.toLowerCase().endsWith('.pdf')) {
        const uploadOption = {
          file: file,
          filename: file.name,
          onProgress: () => {},
          onSuccess: () => {},
          onError: () => {}
        };
        
        customRequest(uploadOption as any);
      } else {
        message.warning(t('home.upload.drag.uploadWarning', { name: file.name }));
      }
      
      resolve();
    }, reject);
  });
};
</script>
<style lang="less" scoped>
.literature-uploader-wrap {
  position: relative;
  .folder {
    position: absolute;
    bottom: 126px;
    left: calc(50% + 75px);
    .mask {
      display: inline-block;
      width: 70px;
      cursor: pointer;
    }
  }
}
.literature-uploader-drag {
  padding: 126px 0;
  text-align: center;
  .hint {
    color: #73716f;
    line-height: 22px;
    margin-top: 13px;
    span {
      cursor: pointer;
    }
  }
  
  // Tauri 环境下的上传区域样式
  &.tauri-uploader {
    border: 1px dashed #e4e7ed;
    background: #fff;
  }
}
// stylelint-disable-next-line
::v-deep .ant-upload-drag {
  background: #fff;
  border: 1px dashed #e4e7ed;
}
// stylelint-disable-next-line
::v-deep .ant-upload.ant-upload-drag .ant-upload {
  padding: 0;
}

.literature-uploader-dragging {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  pointer-events: none;
  .bg {
    background: rgba(0, 0, 0, 0.2);
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    color: rgba(0, 0, 0, 0.65);

    .anticon-arrow-up {
      color: #1f71e0;
      font-size: 50px;
      margin-bottom: 20px;
    }
    .hint1 {
      font-weight: 600;
      line-height: 25px;
      font-size: 18px;
      margin-bottom: 8px;
    }
    .hint2 {
      line-height: 20px;
    }
  }
}
</style>
