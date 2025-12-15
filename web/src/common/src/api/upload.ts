import { v4 as uuid } from 'uuid';
import { SuccessResponse } from './type';
import api from './axios';
import { REQUEST_APPID, REQUEST_SERVICE_NAME_APP } from './const';

export const upload = async (data: FormData) => {
  const res = await api.post<SuccessResponse<any>>(
    REQUEST_SERVICE_NAME_APP + '/file/uploadFile',
    data,
    {
      timeout: 60 * 1000,
    }
  );

  return res.data;
};

export enum ImageStorageType {
  screenshot = 'screenshot',
  markdown = 'note-markdown-photo',
  aiReadingImageQuestion = 'ai-reading-image-question',
}

export const uploadImage = async (
  content: File | string,
  type: ImageStorageType
) => {
  return uploadImageFile(
    typeof content === 'string' ? await urlToFile(content) : content,
    type
  );
};

export const uploadImageFile = async (file: File, type: ImageStorageType) => {
  const formdata = new FormData();
  const ext = file.type.replace(/^image\//, '');
  formdata.append('fileName', `${uuid()}.${ext}`);
  formdata.append('appId', REQUEST_APPID);
  formdata.append('type', type);
  formdata.append('file', file);
  const result = await upload(formdata);
  return `${result.data.staticDomain}/${result.data.fileFullName}`;
};

export const urlToFile = async (url: string) => {
  const blob = await fetch(url).then((res) => res.blob());
  const mime = getImageMime(url);
  const file = new File([blob], uuid(), { type: mime });
  return file;
};

const getImageMime = (url: string) => {
  if (url.startsWith('http') || url.startsWith('file')) {
    const extendName = url.split('.').pop() || 'png';
    if (extendName === 'jpg') {
      return 'image/jpeg';
    }

    return `image/${extendName}`;
  }

  const arr = url.split(',');
  const mime = arr?.[0]?.match(/:(.*?);/)?.[1];
  return mime || 'image/png';
};
