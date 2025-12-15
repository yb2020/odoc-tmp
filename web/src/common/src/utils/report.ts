import { reporter, Clock } from '@idea/aiknowledge-report';
import { onMounted, onUnmounted } from 'vue';

export const enum PageType {
  NOTE = 'note',
  NOTE_TAB = 'note_tab',
  REVIEWER = 'readpaper_ai_reviewer',
  REVISE = 'readpaper_ai_revise',
  UNKNOWN = 'readpaper_ai_unknown',
  POLISH = 'ai_polish',
  aibeans_pay = 'aidou_pay',
  qs_pay = 'ai_mentor_gra',
  premium = 'vip',
  premium_h5 = 'h5_vip',
  library = 'library',
}

export const enum ModuleType {
  POLISH_REWRITE_ZH = 'polish_rewrite_zh',
  POLISH_REWRITE = 'polish_rewrite',
  ZH_TO_EN = 'zh_to_en',
  AI_POLISH_MENTOR = 'ai_polish_mentor',
  AI_MENTOR_GRA = 'ai_mentor_gra',
  REVISE = 'ai_polish',
  UNKNOWN = 'ai_unknown',
}

export enum EventCode {
  readpaperPopupPaperReferenceClick = 'readpaper_popup_paper_reference_click',
  readpaper_page_visit_start = 'readpaper_page_visit_start',
  readpaper_page_duration_time_30s = 'readpaper_page_duration_time_30S',
  readpaper_element_click = 'readpaper_element_click',
  readpaper_module_impression = 'readpaper_module_impression',
  readpaper_ai_polish_result_feedback_click = 'readpaper_ai_polish_result_feedback_click',
  readpaper_ai_polish_response_time = 'readpaper_ai_polish_response_time',
  readpaperElementImpression = 'readpaper_element_impression',
  readpaperVipPayPopupImpression = 'readpaper_vip_pay_popup_impression',
}

export enum ElementName {
  upperCollectionPopup = 'upper_collection_popup',
  upperTeamPopup = 'upper_team_popup',
  upperTeamPaperSharePopup = 'upper_team_paper_share_popup',
  upperTeamPaperShareLimitToast = 'upper_team_paper_share_limit_toast',
  upperPaperAttachPopup = 'upper_paper_attach_popup',
  popupReference = 'popup_reference',
  popupRealNameAuthe = 'popup_real_name_authe',
  popupAutheFile = 'popup_authe_file',
  popupAutheSubmitResult = 'popup_authe_submit_result',
  attachments = 'attachments',
  copilotUploaderPopup = 'copilot_uploader_popup',
}

export interface ReportParams {
  [k: string]: undefined | number | string;
  page_type: PageType | string;
  type_parameter?: string;
  element_name?: string;
  element_parameter?: string;
}

export const reportElementClick = (
  params: ReportParams & {
    status?: 'on' | 'off';
  }
) => {
  return reporter.report(
    {
      event_code: EventCode.readpaper_element_click,
    },
    {
      status: 'none',
      ...params,
    }
  );
};

export const reportModuleImpression = (params: ReportParams) => {
  return reporter.report(
    {
      event_code: EventCode.readpaper_module_impression,
    },
    {
      type_parameter: 'none',
      element_parameter: 'none',
      ...params,
    }
  );
};

export const reportEvent = (eventCode: EventCode, params: ReportParams) => {
  return reporter.report(
    {
      event_code: eventCode,
    },
    {
      type_parameter: 'none',
      element_parameter: 'none',
      ...params,
    }
  );
};

export interface ReportVisitParams {
  [k: string]: undefined | number | string;
  page_type: PageType | string;
  type_parameter: string;
  element_parameter?: string;
  element_name?: string;
  status?: 'on' | 'off';
}

const clock = new Clock();

export const useReportVisitDuration = (
  getTouin: () => string,
  getParams: () => ReportVisitParams,
  checkInactive = () => false
) => {
  const reportPageVisitStart = () => {
    return reporter.report(
      {
        event_code: EventCode.readpaper_page_visit_start,
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
        event_code: EventCode.readpaper_page_duration_time_30s,
        touin: getTouin(),
      },
      getParams()
    );
  };

  onMounted(() => {
    reportPageVisitStart();
    clock.add(reportDuration);
  });

  onUnmounted(() => {
    clock.remove(reportDuration);
  });
};


export const reportElementImpression = (
  event_code: EventCode,
  params: {
    [k: string]: any
    page_type: string
    type_parameter?: string
    element_parameter?: string
    element_name?: string
  }
) => {
  return reporter.report(
    {
      event_code,
    },
    params
  )
}



export enum PopupPaperReferenceElementName {
  bibTex = 'bibtex',
  endNote = 'endnote',
  copy = 'copy',
  selectCopy = 'select_copy',
}


export interface ReportVisitNormalParams {
  page_type: string
  type_parameter?: string
  uid?: string
}

export type ReportPopupPaperReferenceClickParams = ReportVisitNormalParams & {
  literature_format: string
  element_name: string
}

export const reportPopupPaperReferenceClick = (
  params: ReportPopupPaperReferenceClickParams
) => {
  return reporter.report(
    {
      event_code: EventCode.readpaperPopupPaperReferenceClick,
    },
    params
  )
}