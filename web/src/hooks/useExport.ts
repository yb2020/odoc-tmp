import { exportMarkdownOrPDF, getExportPDFUrl } from '@/api/export';
import { polling } from '~/src/util/polling';
import {
  ExportPdfRequest,
  ExportPdfResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/Export';
import { ComputedRef, ref } from 'vue';
import fileDownload from 'js-file-download';
import { message } from 'ant-design-vue';

export const fetchExportPDFWithNoteStatus = (params: ExportPdfRequest) => {
  return polling<ExportPdfRequest, ExportPdfResponse>({
    fn: getExportPDFUrl,
    maxAttempts: 1,
    interval: 5000,
    params,
    validate: (res) => {
      if (res.needFetch) {
        return false;
      }
      return true;
    },
  });
};

export const useExportFile = (
  noteId: string,
  fileName?: ComputedRef<string>
) => {
  const loading = ref({
    md: false,
    pdf: false,
  });

  const exportFile = async (file: 'md' | 'pdf') => {
    if ((loading.value as any)[file]) {
      return;
    }
    (loading.value as any)[file as 'md'] = true;
    try {
      const res = exportMarkdownOrPDF(noteId, file);
      const data = await res;
      fileDownload(
        data,
        file === 'md'
          ? `[批注笔记]${fileName?.value || 'export'}.zip`
          : `[批注笔记]${fileName?.value || 'export'}.pdf`
      );
    } catch (error) {
      message.error((error as Error).message || '导出失败');
    }
    (loading.value as any)[file] = false;
  };
  return {
    exportFile,
    loading,
  };
};
