import { defineStore } from 'pinia';
import {
  ReportType,
  UgcReportSubType,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/UgcReport'

export interface UgcReportState {
  showDialog: boolean;
  reportParams?: UgcReportParams;
}

export interface UgcReportParams {
  reportType: ReportType
  ugcReportSubType?: UgcReportSubType
  contentId: string
  pageUrl: string
}

export const useUgcReportStore = defineStore('ugcReport', {
  state: (): UgcReportState => ({
    showDialog: false,
  }),
  actions: {
    showUgcReportDialog(payload: UgcReportParams) {
      this.showDialog = true;
      this.reportParams = payload
    },
    hideUgcReportDialog() {
      this.showDialog = false;
    },
  },
});
