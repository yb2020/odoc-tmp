import reporter, { Clock } from '@idea/aiknowledge-report';
import { RightSideBarType } from '../components/Right/TabPanel/type';
import { onMounted, onUnmounted } from 'vue';
import { getHostname, isInElectron } from '../util/env';
import { NoteSubTypes } from '../store/note/types';
import { MemoryInfo } from '@vueuse/core';

const clock = new Clock();

export enum EventCode {
  readpaperPageVisitStart = 'readpaper_page_visit_start',
  readpaperPdfUploadSuccess = 'readpaper_pdf_upload_success',
  readpaperAddNote = 'readpaper_add_note',
  readpaperModuleImpression = 'readpaper_module_impression',
  readpaperPageDurationTime30S = 'readpaper_page_duration_time_30S',
  readpaperPopupTranslateFeedbackClick = 'readpaper_popup_translate_feedback_click',
  readpaperTenQListClick = 'readpaper_ten_q_list_click',
  readpaperElementClick = 'readpaper_element_click',
  readpaperTranslateResult = 'readpaper_translate_result',
  readpaperNoteCrashPopupImpression = 'readpaper_note_crash_popup_impression',
  readpaperNoteCrashRetrySuccess = 'readpaper_note_crash_retry_success',
  readpaperPaperItemClick = 'readpaper_paper_item_click',
  readpaperPopupNoteCiteFigureImpression = 'readpaper_popup_note_cite_figure_impression',
  readpaperPopupPaperReferenceClick = 'readpaper_popup_paper_reference_click',
  reportAiAssistReadingResponseTime = 'readpaper_ai_assist_reading_response_time',
  reportGoogleTranslateError = 'readpaper_google_translate_error',
  reportGoogleAPITranslateError = 'readpaper_google_api_translate_error',
  readpaperElementImpression = 'readpaper_element_impression',
}

export enum PageType {
  note = 'note',
  aibeans_pay = 'aidou_pay',
}

export enum ElementName {
  upperCollectionPopup = 'upper_collection_popup',
  upperPaperAttachPopup = 'upper_paper_attach_popup',
  upperPaperTranslatePopup = 'upper_paper_translate_popup',
  upperPaperTranslateLimitPopup = 'upper_paper_translate_limit_popup',
  upperWordTranslatePopup = 'upper_word_translate_popup',
  upperNoteDownloadPopup = 'upper_note_download_popup',
  upperNotePDFDownloadPopup = 'upper_note_pdf_download_popup',
  upperAIAssistPopup = 'upper_ai_assist_popup',
  upperAIAssistLimitPopup = 'upper_ai_assist_limit_popup',
  upper50PAIAssistToast = 'upper_50p_ai_assist_toast',
  upper100PAIAssistToast = 'upper_100p_ai_assist_toast',
  upperTeamNumNotePopup = 'upper_team_num_note_popup',
  upperOcrTranslatePopup = 'upper_ocr_translate_popup',
  upperOcrTranslateLimitPopup = 'upper_ocr_translate_limit_popup',
  updateVersion = 'update_version',
}

export enum ElementClick {
  note_word_recite = 'note_word_recite',
  note_word_recite_card = 'note_word_recite_card',
  version_switch = 'version_switch',

  similar_picture = 'similar_picture',

  copy_note_link = 'copy_note_link',

  collect = 'collect',
  head_portrait = 'head_portrait',
  paper_link = 'paper_link',
  search_paper = 'search_paper',
  my_note = 'my_note',
  attachments = 'attachments',

  go_home = 'go_home',
  shortcut_key = 'shortcut_key',
  suggest = 'suggest',
  screenshots = 'screenshots',
  find = 'find',
  full_screen = 'full_screen',
  size_adjust = 'size_adjust',
  page_adjust = 'page_adjust',
  size_auto_narrow = 'size_auto_narrow',
  size_auto_large = 'size_auto_large',
  directory = 'directory',
  note_bar = 'note_bar',
  bottom_tool_bar = 'bottom_tool_bar',
  renew = 'renew',

  note_independent_window = 'note_independent_window',
  note_word_global_present = 'note_word_global_present',
  note_word_single_present = 'note_word_single_present',
  note_word_not_present = 'note_word_not_present',

  full_translationg_reading = 'full_translationg_reading',
  original_text = 'original_text',

  scim_generate = 'scim_generate',
  scim_keypoint = 'scim_keypoint',
  scim_reading = 'scim_reading',
  save_reference_paper = 'save_reference_paper',

  ai_assist_copy = 'ai_assist_copy',

  trans_box_lock = 'trans_box_lock',
  trans_box_release = 'trans_box_release',
  trans_lock_right = 'trans_lock_right',
}

if (isInElectron()) {
  reporter.config.api = `https://${getHostname()}` + reporter.config.api;
  if (window.location.protocol.startsWith('file')) {
    reporter.config.getUin = () => {
      const uidString = localStorage.getItem('rp-user/uid');
      if (uidString) {
        const uid = JSON.parse(uidString);
        if (uid) {
          return uid;
        }
      }

      return '';
    };
  }
}

export interface ReportVisitParams {
  page_type: PageType;
  type_parameter: string;
  element_parameter?: string;
  element_name?: string;
  status?: 'on' | 'off' | 'none';
}

export const reportPdfUploadSuccess = (pdfId: string) => {
  return reporter.report(
    {
      event_code: EventCode.readpaperPdfUploadSuccess,
    },
    {
      pdf_id: pdfId,
    }
  );
};

export const reportAddNote = () => {
  return reporter.report(
    {
      event_code: EventCode.readpaperAddNote,
    },
    {}
  );
};

export const useReportVisitDuration = (
  getTouin: () => string,
  getParams: () => ReportVisitParams,
  checkInactive = () => false
) => {
  const reportPageVisitStart = () => {
    return reporter.report(
      {
        event_code: EventCode.readpaperPageVisitStart,
        touin: getTouin(),
      },
      getParams()
    );
  };

  const reportDuration = () => {
    if (checkInactive && checkInactive()) {
      return;
    }

    reporter.report(
      {
        event_code: EventCode.readpaperPageDurationTime30S,
        touin: getTouin(),
      },
      {
        ...getParams(),
        memory: getMemory(),
      }
    );
  };

  const getMemory = () => {
    if (typeof performance === 'undefined' || !(performance as any).memory) {
      return '';
    }

    const toKB = (num: number) => Math.round(num / 1024);
    const { usedJSHeapSize, totalJSHeapSize, jsHeapSizeLimit } = (
      performance as any
    ).memory as MemoryInfo;
    return [usedJSHeapSize, totalJSHeapSize, jsHeapSizeLimit]
      .map(toKB)
      .join(',');
  };

  onMounted(() => {
    reportPageVisitStart();
    clock.add(reportDuration);
  });

  onUnmounted(() => {
    clock.remove(reportDuration);
  });
};

export const reportElementClick = (
  params: { element_name: ElementClick } & ReportVisitParams
) => {
  return reporter.report(
    {
      event_code: EventCode.readpaperElementClick,
    },
    params
  );
};

const reportedRightTab: Record<string, boolean> = {};
export const reportRightTab = (
  tab: RightSideBarType,
  subTab?: number | string
) => {
  console.log(tab, subTab);

  const values = {
    [RightSideBarType.Group]: 'team_discussion',
    [RightSideBarType.Matirial]: 'data',
    [RightSideBarType.Note]: {
      [NoteSubTypes.Summary]: 'note_summary',
      [NoteSubTypes.Vocabulary]: 'note_word',
      [NoteSubTypes.Annotation]: 'note_excerpt',
    },
    [RightSideBarType.Question]: 'qa',
    [RightSideBarType.Learn]: 'learning_task',
    [RightSideBarType.Copilot]: 'ai_auxiliary',
    [RightSideBarType.Translate]: 'trans_box',
  }[tab];
  const mod: undefined | string =
    typeof values === 'string'
      ? values
      : // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        values?.[subTab];

  if (!mod || (reportedRightTab[mod] && tab !== RightSideBarType.Translate)) {
    return;
  }

  reportedRightTab[mod] = true;
  return reporter.report(
    {
      event_code: EventCode.readpaperModuleImpression,
    },
    {
      page_type: 'note',
      module_type: mod,
    }
  );
};

export const reportTenQListClick = (
  params: {
    question_id: string;
  } & ReportVisitParams
) => {
  return reporter.report(
    {
      event_code: EventCode.readpaperTenQListClick,
    },
    params
  );
};

export const reportTranslateClick = (
  params: {
    tran_content: string;
    sources: string;
  } & ReportVisitParams
) => {
  return reporter.report(
    {
      event_code: EventCode.readpaperTranslateResult,
    },
    params
  );
};

export const reportModuleImpression = (
  params: { module_type: string } & ReportVisitParams
) => {
  return reporter.report(
    {
      event_code: EventCode.readpaperModuleImpression,
    },
    params
  );
};

export const getPdfIdFromUrl = () => {
  let pdfId = '';
  const PDF_ID = 'pdfId';

  const url = new URL(location.href);
  if (url.protocol === 'file:') {
    const parts = url.hash.split('?');
    parts.shift();
    const search = new URLSearchParams(parts.pop() || '');
    pdfId = search.get(PDF_ID) || '';
  } else {
    pdfId = url.searchParams.get(PDF_ID) || '';
  }

  return pdfId;
};

export const reportClick = (
  elementName: ElementClick,
  status?: 'on' | 'off'
) => {
  return reportElementClick({
    element_name: elementName,
    page_type: PageType.note,
    type_parameter: getPdfIdFromUrl(),
    status,
  });
};

export type ReportPaperItemParams = ReportVisitParams & {
  module_type: string;
  subject_id: string;
  paper_id: string;
  scene_id: string;
  order_num?: number;
  ref_content?: string;
};

export const reportPaperItemClick = (params: ReportPaperItemParams) => {
  return reporter.report(
    {
      event_code: EventCode.readpaperPaperItemClick,
    },
    params
  );
};

export type ReportPopupNoteCiteFigureImpressionParams = ReportVisitParams & {
  popup_type: string;
  popup_id: string;
  order_num: string;
};

export const reportPopupNoteCiteFigureImpression = (
  params: ReportPopupNoteCiteFigureImpressionParams
) => {
  return reporter.report(
    {
      event_code: EventCode.readpaperPopupNoteCiteFigureImpression,
    },
    params
  );
};

export const reportReadpaperPopupPaperReferenceClick = (
  params: ReportVisitParams
) => {
  return reporter.report(
    {
      event_code: EventCode.readpaperPopupPaperReferenceClick,
    },
    params
  );
};

export type ReportAiAssistReadingResponseTimeParams = {
  answer_id: string;
  pdf_title: string;
  first_word_response_time: number;
  last_word_response_time: number;
  requestion_statr_time: number;
  language_type: string;
};

export const reportAiAssistReadingResponseTime = (
  params: ReportAiAssistReadingResponseTimeParams
) => {
  return reporter.report(
    {
      event_code: EventCode.reportAiAssistReadingResponseTime,
    },
    {
      page_type: PageType.note,
      type_parameter: getPdfIdFromUrl(),
      ...params,
    }
  );
};

export type ReportGoogleTranslateErrorParams = {
  source_content: string;
  error_message: string;
  request_time: number;
};

export type ReportGoogleAPITranslateErrorParams = {
  source_content: string;
  error_message: string;
  request_time: number;
  api_ke: string;
};
export const reportGoogleAPITranslateError = (
  params: ReportGoogleAPITranslateErrorParams
) => {
  return reporter.report(
    {
      event_code: EventCode.reportGoogleAPITranslateError,
    },
    {
      page_type: PageType.note,
      type_parameter: getPdfIdFromUrl(),
      ...params,
    }
  );
};

export const reportGoogleTranslateError = (
  params: ReportGoogleTranslateErrorParams
) => {
  return reporter.report(
    {
      event_code: EventCode.reportGoogleTranslateError,
    },
    {
      page_type: PageType.note,
      type_parameter: getPdfIdFromUrl(),
      ...params,
    }
  );
};

export const reportElementImpression = (
  params: { element_name: string } & ReportVisitParams
) => {
  return reporter.report(
    {
      event_code: EventCode.readpaperElementImpression,
    },
    params
  );
};
