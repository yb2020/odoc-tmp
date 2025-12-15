import { defineStore } from 'pinia';
import { GetAttachmentListResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/DocAttachment';

export interface AttachmentsState {
  attachmentsInfo?: GetAttachmentListResponse;
}

export const useAttachmentsStore = defineStore('attachments', {
  state: (): AttachmentsState => ({
    attachmentsInfo: undefined,
  }),
  actions: {
    setAttachmentData(data: GetAttachmentListResponse) {
      this.attachmentsInfo = data;
    },
  },
});
