import FingerprintJS from '@fingerprintjs/fingerprintjs'

import api from '@/common/src/api/axios'
import {
  HEADER_CANCLE_AUTO_ERROR,
  REQUEST_SERVICE_NAME_SEARCH,
} from '@/common/src/api/const'

export class StorageAccessor {
  private storage: Storage
  private enable = false
  private static instance: StorageAccessor
  constructor(storage: Storage) {
    this.storage = storage
    this.canUse()
  }

  static getInstance(storage: Storage) {
    if (!StorageAccessor.instance) {
      StorageAccessor.instance = new StorageAccessor(storage)
    }
    return StorageAccessor.instance
  }

  private canUse() {
    try {
      this.storage.setItem('test', '1')
      if (this.storage.getItem('test') === '1') {
        this.enable = true
        return
      }
    } catch (error) {
      // eslint-disable-next-line no-console
      console.warn('storage is not enable')
    }
  }

  public get(key: string) {
    return this.storage.getItem(key)
  }

  public set(key: string, value: string) {
    if (!this.enable) {
      return false
    }
    try {
      this.storage.setItem(key, value)
      return true
    } catch (e) {
      return false
    }
  }

  public remove(key: string) {
    return this.storage.removeItem(key)
  }
}

export default function useBehaviorReport() {
  return (
    data: UserActionReportItem,
    cb?: ReportUserActionCallBackFn,
    delay?: number
  ) => {
    reportUserBehaviorAction(data, cb, delay)
  }
}

enum ReportPlatform {
  pc = 'pc',
  mobile = 'h5',
  unknown = 'unknown',
}

export type ReportUserActionCallBackFn = (err?: string) => void

type ReportUserActionFn = (
  data: UserActionReportItem,
  cb?: ReportUserActionCallBackFn,
  delay?: number
) => void

interface $ReportItem {
  actionTime: number
  actionType?: ReportActionType
  generalIdentifier: string
  itemId: string
  itemType?: ReportItemType
  platform: ReportPlatform
  sceneId: ReportSceneId | string
  otherInfo?: string
}

export type UserActionReportItem = Pick<
  $ReportItem,
  'actionType' | 'itemId' | 'itemType' | 'sceneId' | 'otherInfo'
>

export enum ReportActionType {
  addTag = 'add_tag',
  removeTag = 'remove_tag',
  click = 'click',
  unfoldAbstraction = 'unfold_abstraction',
  like = 'like',
  collect = 'collect',
  clickComment = 'click_comment',
  clickTopic = 'click_topic',
  show = 'show',
  clickBookshelf = 'click_bookshelf',
  clickReadpdf = 'click_readpdf',
  clickQuestonsMore = 'click_questions_more',
  clickQuestionsGowrite = 'click_questions_gowrite',
  postRelated = 'post_related',
  clickNoteGowrite = 'click_note_gowrite',
  clickRef = 'click_ref',
  clickQuote = 'click_quote',
  clickShare = 'click_share',
  search = 'search',
  refresh = 'refresh',
  relatedPaperShow = 'related_paper_show',
  relatedPaperRefresh = 'related_paper_refresh',
  relatedPaperClick = 'related_paper_click',
  topicShow = 'topic_show',
  topicClick = 'topic_click',
}

export enum ReportItemType {
  paper = 'paper',
  field = 'field',
}

export enum ReportSceneId {
  unknown = 'unknown',
  index = 'index',
  detail = 'detail',
  detailWeixin = 'detail_weixin',
  detailQQ = 'detail_qq',
  detailQzone = 'detail_qzone',
  detailWeibo = 'detail_weibo',
  indexSearch = 'index_search',
  detailSearch = 'detail_search',
  fieldOfStudy = 'feildofstudy',
}

const generateIdentifierForBrowser = async (): Promise<string> => {
  const ls = StorageAccessor.getInstance(window.localStorage)
  const key = '__paper_reader_browser_identifier_report_web__'
  const identifier = ls.get(key)
  if (identifier) {
    return identifier
  }
  try {
    const fpPromise = FingerprintJS.load()
    const fp = await fpPromise
    const result = await fp.get()
    ls.set(key, result.visitorId)
    return result.visitorId || ''
  } catch (error) {
    return ''
  }
}

export const reportUserBehaviorAction: ReportUserActionFn = (() => {
  let identifier = ''

  let platform: ReportPlatform = ReportPlatform.unknown

  // 这里用于缓存上报的数据
  const cache: $ReportItem[] = []

  let timer: number

  return async (
    item: UserActionReportItem,
    cb?: ReportUserActionCallBackFn,
    delay?: number
  ) => {
    if (!identifier) {
      identifier = await generateIdentifierForBrowser()
    }

    if (platform === ReportPlatform.unknown) {
      const isMobile =
        /(iPhone|iPod|Opera Mini|Android.*Mobile|NetFront|PSP|BlackBerry|Windows Phone)/gi.test(
          window.navigator.userAgent
        )
      platform = isMobile ? ReportPlatform.mobile : ReportPlatform.pc
    }

    cache.push({
      ...item,
      generalIdentifier: identifier,
      platform,
      actionTime: Date.now(),
    })

    if (timer) {
      window.clearTimeout(timer)
    }

    // TODO 是否考虑用sendBeacon
    // 默认是立刻上报，可以设置为延迟上报
    timer = window.setTimeout(
      () => {
        api({
          url: `${REQUEST_SERVICE_NAME_SEARCH}/interviewHistory/saveUserInterviewHistory`,
          method: 'post',
          data: cache,
          headers: {
            [HEADER_CANCLE_AUTO_ERROR]: true,
          },
        })
          .then(() => {
            cb && cb()
            // 上报成功，清除cache
            cache.length = 0
          })
          .catch((err: Error) => {
            cb && cb(err.message)
          })
      },
      typeof delay === 'number' ? delay : 0
    )
  }
})()
