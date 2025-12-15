import { DeleteAttachmentReq } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/DocAttachment';
import { filter } from 'lodash-es';
import { useRequest } from 'ahooks-vue';
import { Ref, ref } from 'vue';
import {
  getDocAttachments,
  getUploadToken,
  removeAttachment,
  saveAttachment,
  uploadAttachment,
} from '@/api/attachments';
import { ResponseError } from '@/api/type';
import { AxiosInstance } from 'axios';
import { Modal, message } from 'ant-design-vue';
import { useVipStore } from '@common/stores/vip';
import { useAttachmentsStore } from '../stores/attachmentsStore';
import { ERROR_CODE_NEED_VIP } from '@common/api/const';
import { ElementName } from '../api/report';
import { useI18n } from 'vue-i18n';

export function useDocAttachments(
  $axios: AxiosInstance,
  docId: Ref<string>,
  ready: Ref<boolean>,
  getFileMD5?: (x: File) => Promise<string>
) {
  const vipStore = useVipStore();
  const attachmentsStore = useAttachmentsStore();
  const uploadFile = ref<null | File>(null);

  const { data, refresh } = useRequest(
    async () => {
      const info = await getDocAttachments($axios, { docId: docId.value });
      attachmentsStore.setAttachmentData(info);
      return info;
    },
    {
      ready,
    }
  );

  const {
    run: doUpload,
    loading: uploading,
    cancel: cancelUpload,
  } = useRequest(
    async (file: File) => {
      uploadFile.value = file;

      let md5;
      try {
        md5 = await getFileMD5?.(file);
      } catch (error) {
        console.error(error);
      }

      const tokens = await getUploadToken($axios, {
        docId: docId.value,
        md5,
        size: file.size,
        fileName: file.name,
        contentType: file.type,
      }).catch((e) => {
        const err = e as ResponseError;
        if (err.code === ERROR_CODE_NEED_VIP) {
          vipStore.showVipLimitDialog(err.message, {
            exception: e.extra,
            reportParams: {
              element_name: ElementName.upperPaperAttachPopup,
              element_parameter: docId.value,
            },
          });
          return;
        }
        throw e;
      });
      if (!tokens) {
        return false;
      }
      if (!tokens.isNeedUpload) {
        refresh();
        return true;
      }
      const res = await uploadAttachment(file, tokens).catch((e: Error) => {
        message.error(e.message);
        throw e;
      });

      if (!uploadFile.value || !res.objectName) {
        // 被取消
        return false;
      }
      const arr = res.objectName.split('/');
      const fileName = arr.pop()!;
      const result = await saveAttachment($axios, {
        docId: docId.value,
        size: file.size,
        fileName,
        ossObjectName: res.objectName,
        contentType: res.mimeType!,
      });

      uploadFile.value = null;
      refresh();
      return result;
    },
    {
      manual: true,
    }
  );

  const removingIds = ref<string[]>([]);

  const { t } = useI18n();

  const { run: doRemove } = useRequest(
    async (params: DeleteAttachmentReq) => {
      return new Promise<boolean>((resolve) => {
        Modal.confirm({
          title: t('info.deleteAttachmentTip1'),
          content: t('info.deleteAttachmentTip2'),
          okButtonProps: {
            danger: true,
          },
          cancelButtonProps: {
            type: 'primary',
          },
          onOk: async () => {
            removingIds.value.push(params.attachmentId);
            await removeAttachment($axios, params);
            removingIds.value = filter(
              removingIds.value,
              (id: string) => id !== params.attachmentId
            );
            resolve(true);
            refresh();
          },
          onCancel: () => resolve(false),
        });
      });
    },
    {
      manual: true,
    }
  );

  const doCancelUpload = () => {
    uploadFile.value = null;
    cancelUpload();
  };

  return {
    data,
    refresh,
    uploading,
    uploadFile,
    removingIds,
    doUpload,
    doRemove,
    cancelUpload: doCancelUpload,
  };
}
