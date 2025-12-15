import { reporter } from '@idea/aiknowledge-report'
import { onMounted, onUnmounted } from 'vue'

import { PAGE_ROUTE_NAME } from '../routes/type'

export enum EventCode {
  readpaperPageVisitStart = 'readpaper_page_visit_start',
  readpaperPdfUploadSuccess = 'readpaper_pdf_upload_success',
  readpaperAddNote = 'readpaper_add_note',
  readpaperModuleImpression = 'readpaper_module_impression',
  readpaperPageDurationTime30S = 'readpaper_page_duration_time_30S',
  readpaperPaperItemClick = 'readpaper_paper_item_click',
  readpaperPaperItemImpression = 'readpaper_paper_item_impression',
  readpaperTenQListClick = 'readpaper_ten_q_list_click',
  readpaperElementClick = 'readpaper_element_click',
  readpaperSearchButtonClick = 'readpaper_search_button_click',
  readpaperInvitationToSignupResult = 'readpaper_invitation_to_signup_result',
  readpaperElementImpression = 'readpaper_element_impression',
  readpaperFeedsClick = 'readpaper_feeds_click',
  readpaperFeedsImpression = 'readpaper_feeds_impression',
  readpaperPublicAccountItemClick = 'readpaper_public_account_item_click',
  readpaperPopupPaperReferenceClick = 'readpaper_popup_paper_reference_click',
  readpaperPopupOfficialMsgsImpression = 'readpaper_popup_official_news_impression',
  readpaperNoticeClick = 'readpaper_notice_news_click',
  readpaperVipPayPopupImpression = 'readpaper_vip_pay_popup_impression',
}

export const PageType = {
  [PAGE_ROUTE_NAME.WORKBENCH]: 'home',
  [PAGE_ROUTE_NAME.LIBRARY]: 'library',
  [PAGE_ROUTE_NAME.NOTES]: 'note_tab',
  [PAGE_ROUTE_NAME.NOTE]: 'note',
  UNKNOWN: 'readpaper_ai_unknown',

  // vip弹窗上报
  search: 'search',

  // PaperItem上报
  subjectPage: 'subject_page',
  newPaper: 'paper_new',
  premium: 'vip',
  premium_h5: 'h5_vip',
  pay_premium: 'hani_pay',
  vipTaskPage: 'vip_task_page',
}

export interface ReportVisitNormalParams {
  page_type: string
  type_parameter?: string
  uid?: string
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

export enum PopupPaperReferenceElementName {
  bibTex = 'bibtex',
  endNote = 'endnote',
  copy = 'copy',
  selectCopy = 'select_copy',
}

export interface ReportElementClickParams extends ReportVisitNormalParams {
  element_name: string
  element_parameter?: string
  status?: string
  uid?: string
}
export const reportElementClick = (params: ReportElementClickParams) => {
  return reporter.report(
    {
      event_code: EventCode.readpaperElementClick,
    },
    params
  )
}

export interface ReportSearchButtonClickParams extends ReportVisitNormalParams {
  start_year?: string
  end_year?: string
  journals?: string
  authors?: string
  field?: string
  sorting?: string
  query_tag?: string
  query_type?: 'text' | 'picture'
  scene_id?: string
}

export type ReportFeedsClickParams = ReportVisitNormalParams & {
  module_type: string
  feed_id: string
  paper_id: string
  scene_id: string
  order_num: number
  feed_type: string
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

export interface ReportVisitPersonParams {
  page_type: string
  touin: string
}

export type ReportVisitParams =
  | ReportVisitNormalParams
  | ReportVisitPersonParams

export type ReportPaperItemParams = ReportVisitParams & {
  event_code:
    | EventCode.readpaperPaperItemClick
    | EventCode.readpaperPaperItemImpression
  module_type: string
  paper_id: string
  scene_id: string
  order_num?: number
  subject_id: string
}

export const reportPaperItem = ({
  event_code,
  ...params
}: ReportPaperItemParams) => {
  return reporter.report(
    {
      event_code,
    },
    params
  )
}

export const useReportVisitDuration = (
  getUid: () => string,
  getParams: () => ReportVisitNormalParams,
  isDestroyed?: () => boolean
) => {
  let timer: number | null = null
  let startTime = 0

  const startReport = () => {
    startTime = Date.now()
    if (timer) {
      clearInterval(timer)
    }
    timer = window.setInterval(() => {
      if (isDestroyed && isDestroyed()) {
        clearInterval(timer!)
        return
      }
      const params = getParams()
      const uid = getUid()
      reporter.report(
        {
          event_code: EventCode.readpaperPageDurationTime30S,
        },
        {
          ...params,
          ...(uid ? { uid } : {}),
        }
      )
    }, 30 * 1000)
  }

  const stopReport = () => {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
  }

  onMounted(() => {
    startReport()
  })

  onUnmounted(() => {
    stopReport()
  })

  return {
    startReport,
    stopReport,
  }
}
