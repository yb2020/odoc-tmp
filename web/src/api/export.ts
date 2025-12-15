import api from './axios';
import { SuccessResponse } from './type';
import {
  ExportPdfRequest,
  ExportPdfResponse,
  ExportCountIncrRequest,
  ExportCountIncrResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/Export';

export const getExportEnable = async () => {
  const res = await api.get<SuccessResponse<boolean>>(
    `/microservice-readpaper-pdf/note/isVip`,
    {}
  );
  return !!res.data.data;
};

export const getExportPDFUrl = async (params: ExportPdfRequest) => {
  const res = await api.get<SuccessResponse<ExportPdfResponse>>(
    `/micrgo-pdf/export`,
    {
      params,
    }
  );
  return res.data.data;
};

export const validExportCountIncr = async (params: ExportCountIncrRequest) => {
  const res = await api.post<SuccessResponse<ExportCountIncrResponse>>(
    `/micrgo-pdf/export/countIncr`,
    params
  );
  return res.data.data;
};

export const readJsonFile = async (data: Blob) => {
  const reader = new FileReader();
  reader.readAsText(data);
  return new Promise((resolve, reject) => {
    reader.onload = function () {
      resolve(reader.result);
    };
    reader.onerror = function () {
      reject(reader.error);
    };
  });
};

export const exportMarkdownOrPDF = async (
  noteId: string,
  type: 'md' | 'pdf'
) => {
  // const res = await api.get<Blob>(
  //   `/microservice-readpaper-pdf/note/${
  //     type === 'md' ? 'downloadMarkDown' : 'downloadPdf'
  //   }`,
  //   {
  //     params: {
  //       noteId,
  //     },
  //     responseType: 'blob',
  //     timeout: 1000 * 60 * 5,
  //   }
  // );

  const res = await api.get<Blob>(
    `/note/paperNote/${
      type === 'md' ? 'downloadNoteMarkdown' : 'downloadNotePdf'
    }`,
    {
      params: {
        noteId,
      },
      responseType: 'blob',
      timeout: 1000 * 60 * 5,
    }
  );

  console.debug('exportMarkdownOrPDF', res);

  if (res.data.type === 'application/json') {
    const error = await readJsonFile(res.data);
    const errorJson = JSON.parse(error as string);
    throw new Error((errorJson as any).message || '下载失败');
  }

  return res.data;
};
