import { Ref, ref } from 'vue';
import axios, { AxiosResponse } from 'axios';
import type { UploadFile, UploadProps } from 'ant-design-vue';
import { message } from 'ant-design-vue';
import { calculateFileMD5 } from '@common/utils/md5';
import {
  AcquirePolicyCallbackInfoResponse,
  ObjectStoreInfo,
  UploadBizScene,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/oss/AliOSS';
import { useI18n } from 'vue-i18n';

type FileType = UploadFile;

const checkFile = async (
  file: File,
  fn: UsePDFUploadOptions['acquirePolicyCallbackInfo'],
  scene = UploadBizScene.PDF
) => {
  // 计算file文件的md5值
  let md5 = '';
  try {
    md5 = await calculateFileMD5(file);
  } catch (error) {
    // 计算失败，继续上传
    console.error(error);
  }

  // 检查文件是否已经上传过 并获取上传所需的policy和callback
  const res = await fn({ md5, scene });

  return res;
};

export interface UsePDFUploadOptions {
  scene?: UploadBizScene;
  type?: Ref<'Latex' | 'Word'>;
  acquirePolicyCallbackInfo: (p: {
    md5: string;
    scene: UploadBizScene;
  }) => Promise<AcquirePolicyCallbackInfoResponse>;
  i18n?: ReturnType<typeof useI18n>;
}

const uploadFileToOSS = async (
  file: FileType,
  emit: UploadEmitEvent,
  acquirePolicyCallbackInfoFn: UsePDFUploadOptions['acquirePolicyCallbackInfo'],
  scene = UploadBizScene.PDF
): Promise<ObjectStoreInfo> => {
  const { accessid, callback, dir, host, policy, signature, ...rest } =
    await checkFile(file.originFileObj!, acquirePolicyCallbackInfoFn, scene);

  if (!rest.isNeedUpload && rest.objectStoreInfo) {
    return rest.objectStoreInfo;
  }

  emit('import:bucket:progress', Math.floor(Math.random() * 30));

  const formData = new FormData();

  formData.append('name', file.name);
  formData.append('key', `${dir}${file.name}`);
  formData.append('policy', policy);
  formData.append('OSSAccessKeyId', accessid);
  formData.append('success_action_status', '200'); // 让服务端返回200,不然，默认会返回204
  formData.append('callback', callback);
  formData.append('signature', signature);
  formData.append('file', file.originFileObj!);
  const res: AxiosResponse<{
    bucketName: string;
    fileName: string;
    isSuccess: boolean;
    size: string;
    mimeType: string;
  }> = await axios(host, {
    method: 'post',
    data: formData,
  });

  if (res.status === 200 && res.data?.isSuccess) {
    return {
      bucketName: res.data.bucketName,
      objectName: res.data.fileName,
    };
  }

  throw new Error(
    `Request Failed With Status Code ${
      res.status === 200 ? res.data?.isSuccess : res.status
    }`
  );
};

export interface UploadEmitEvent {
  (
    event: 'import:bucket:progress',
    progress: number,
    bucket?: ObjectStoreInfo
  ): void;
  (event: 'import:bucket:error', error: Error): void;
}

export const usePDFUpload = (
  emit: UploadEmitEvent,
  options: UsePDFUploadOptions
) => {
  const uploadError = ref<Error | null>(null);
  const fileList = ref<UploadProps['fileList']>([]);

  const beforeUpload: UploadProps['beforeUpload'] = (file) => {
    console.log('hooks:usePDFUpload:beforeUpload', file);
    const { type = ref('Latex') } = options;
    const fileType =
      options.scene === UploadBizScene.AI_POLISH
        ? type.value === 'Latex'
          ? [
              'application/x-tar',
              'application/zip',
              'application/x-zip',
              'application/x-zip-compressed',
              'application/octet-stream',
              'application/zip-compressed',
              'multipart/x-zip',
            ]
          : [
              'application/msword',
              'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
            ]
        : ['application/pdf'];

    if (fileType.includes(file.type) === false) {
      // message.error(`${file.name} is not a ${fileType} file`)
      if (options.scene === UploadBizScene.AI_POLISH) {
        message.error(
          options.i18n?.t(`common.uploadFile.invalid${type.value}Filetype`) ||
            `${file.name} is not a ${fileType} file, please upload tar file`
        );
        return false;
      }
      message.error(
        options.i18n?.t('common.uploadFile.invalidFiletype', {
          filename: file.name,
          filetype: fileType,
        }) || `${file.name} is not a ${fileType} file`
      );
      return false;
    }
    if (file.size > 100 * 1024 * 1024) {
      if (options.scene === UploadBizScene.AI_POLISH) {
        message.error(
          options.i18n?.t(`common.uploadFile.invalid${type.value}Filesize`, {
            filesize: '100M',
          }) || `${file.name} is too large, please upload file less than 100Mb`
        );
        return false;
      }
      message.error(
        options.i18n?.t('common.uploadFile.invalidFilesize', {
          filename: file.name,
          filesize: '100Mb',
        }) || `${file.name} is too large, please upload file less than 100Mb`
      );
      // message.error(
      //   `${file.name} is too large, please upload file less than 100Mb`
      // )
      return false;
    }
    fileList.value = [...(fileList.value as any), file];
    return true;
  };

  const handleUpload: UploadProps['customRequest'] = async () => {
    console.debug('hooks:usePDFUpload:handleUpload', fileList.value);
    const cur = fileList.value?.[0];
    if (!cur) {
      return;
    }
    try {
      emit('import:bucket:progress', 0);
      const objectStoreInfo = await uploadFileToOSS(
        cur as FileType,
        emit,
        options.acquirePolicyCallbackInfo,
        options.scene
      );
      console.debug(
        'hooks:usePDFUpload:handleUpload:objectStoreInfo',
        objectStoreInfo
      );
      fileList.value = [];
      emit('import:bucket:progress', 100, objectStoreInfo);
    } catch (error) {
      uploadError.value = error as Error;
      emit('import:bucket:error', error as Error);
    }
  };

  return {
    fileList,
    beforeUpload,
    handleUpload,
  };
};
