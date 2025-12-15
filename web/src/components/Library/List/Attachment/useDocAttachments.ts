import { DeleteAttachmentReq } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/DocAttachment'
import { filter } from 'lodash'
import { useRequest } from 'ahooks-vue'
import { Ref, ref, h } from 'vue'
import {
  getDocAttachments,
  getUploadToken,
  removeAttachment,
  saveAttachment,
  uploadAttachment,
} from '@/common/src/api/attachments'
import { Modal, message } from 'ant-design-vue'
import { ExclamationCircleFilled } from '@ant-design/icons-vue'
import { ResponseError } from '@/common/src/api/type'
import { LimitDialogReportParams } from '../../helper'
import { useVipStore } from '@/common/src/stores/vip'
import { ERROR_CODE_NEED_VIP } from '@/common/src/api/const'
import { ElementName, PageType } from '@/common/src/utils/report'
import { NeedVipException } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo'

export function useDocAttachments(
  docId: Ref<string>,
  ready: Ref<boolean>,
  getFileMD5?: (x: File) => Promise<string>
) {
  const vipStore = useVipStore()

  const uploadFile = ref<null | File>(null)

  const { data, refresh } = useRequest(
    () => {
      return getDocAttachments({ docId: docId.value })
    },
    {
      ready,
    }
  )

  const {
    run: doUpload,
    loading: uploading,
    cancel: cancelUpload,
  } = useRequest(
    async (file: File) => {
      uploadFile.value = file

      let md5
      try {
        md5 = await getFileMD5?.(file)
      } catch (error) {
        console.error(error)
      }

      const tokens = await getUploadToken({
        docId: docId.value,
        md5,
        size: file.size,
        fileName: file.name,
        contentType: file.type,
      }).catch((e) => {
        const err = e as ResponseError
        if (err.code === ERROR_CODE_NEED_VIP) {
          vipStore.showVipLimitDialog(err.message, {
            exception: err.extra as NeedVipException,
            reportParams: {
              page_type: PageType.NOTE,
              element_name: ElementName.top_bar,
            },
          })
          return
        }
        throw e
      })
      if (!tokens) {
        return false
      }
      if (!tokens.isNeedUpload) {
        await refresh()
        return true
      }
      const res = await uploadAttachment(file, tokens).catch((e: Error) => {
        message.error(e.message)
        throw e
      })

      if (!uploadFile.value || !res.objectName) {
        // 被取消
        return false
      }
      const arr = res.objectName.split('/')
      const fileName = arr.pop() as string
      const result = await saveAttachment({
        docId: docId.value,
        size: file.size,
        fileName,
        ossObjectName: res.objectName,
        contentType: res.mimeType as string,
      })

      await refresh()
      uploadFile.value = null
      return result
    },
    {
      manual: true,
    }
  )

  const removingIds = ref<string[]>([])

  const { run: doRemove } = useRequest(
    (params: DeleteAttachmentReq) => {
      return new Promise<boolean>((resolve) => {
        Modal.confirm({
          class: 'attachments-rm-modal',
          title: '确认要删除这条附件吗？',
          icon: h(ExclamationCircleFilled, {
            style: {
              fontSize: '24px',
              color: '#FAAD14',
            },
          }),
          okText: '删除',
          okType: 'danger',
          cancelText: '取消',
          onOk: async () => {
            removingIds.value.push(params.attachmentId)
            await removeAttachment(params)
            removingIds.value = filter(
              removingIds.value,
              (id: string) => id !== params.attachmentId
            )
            await refresh()
            resolve(true)
          },
          onCancel: () => resolve(false),
        })
      })
    },
    {
      manual: true,
    }
  )

  const doCancelUpload = () => {
    uploadFile.value = null
    cancelUpload()
  }

  return {
    data,
    refresh,
    uploading,
    uploadFile,
    removingIds,
    doUpload,
    doRemove,
    cancelUpload: doCancelUpload,
  }
}
