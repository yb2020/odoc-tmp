import { escape } from 'lodash'

import { UserDocDisplayAuthor } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/UserDocManage'

export const escapeHightlight = (content: string): string => {
  return escape(content)
    .replaceAll('&lt;em&gt;', '<em>')
    .replaceAll('&lt;/em&gt;', '</em>')
}

export const removeEmTag = (content: string) => {
  return content.replace(/(<em>)|(<\/em>)/g, '')
}

export const getPopupContainer = (triggerNode: HTMLElement) =>
  (triggerNode.parentElement?.closest('.page') as HTMLElement) || document.body

export const getOneVenue = (venuInfos: string[]) => {
  return venuInfos.find(Boolean) || ''
}

export const getAuthorName = (
  author: UserDocDisplayAuthor['authorInfos'][0]
) => {
  return author.literal
}

export type { ReportParams as LimitDialogReportParams } from '@/common/src/stores/vip'

export const LIBRARY_CONTAINER_CLASSNAME = 'readpaper-library-container'
