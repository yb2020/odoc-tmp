import { PDFFindControllerOptions } from '@idea/pdfjs-dist/types/web/pdf_find_controller'
import * as pdfjsViewer from '@idea/pdfjs-dist/web/pdf_viewer'

export enum FindState {
  FOUND = 0,
  NOT_FOUND = 1,
  WRAPPED = 2,
  PENDING = 3,
  FINISH = 1000,
}
export interface SearchResult {
  list: {
    page: number
    items: SearchResultItem[]
  }[]
  total: number
}

export interface SearchResultOffset {
  pageIdx: number
  matchIdx: number
}

export interface SearchResultItem {
  offset: SearchResultOffset
  prevOffset: {
    pageIdx: number | null
    matchIdx: number | null
  }
  content: {
    str: string
    matched: boolean
  }[]
  whole: string
}

export interface DispatchPayload {
  state: FindState
  source: FinderController
  rawQuery: string
  searchResults?: SearchResult
}

export enum InnerFinderEvent {
  HIGHLIGHT = 'x-highlight-canvas',
  GOTO_MATCH = 'x-goto-match',
}

export default class FinderController extends pdfjsViewer.PDFFindController {
  highlightColor?: string

  highlightTwinklingEnabled = false

  searchResults: SearchResult = {
    list: [],
    total: 0,
  }

  private _searchCallback: ((payload: DispatchPayload) => void)[] = []

  constructor(options: PDFFindControllerOptions) {
    super(options)
    this._pendingFindMatches = new Set()
    this._eventBus.on(
      'updatetextlayermatches',
      async (payload: { source: FinderController; pageIndex: number }) => {
        if (payload.pageIndex === 0) {
          await Promise.all(
            payload.source._extractTextPromises as Promise<void>[]
          )
          this.searchResults = this._formateSearchResults()
          this._searchCallback.forEach((callback) => {
            callback({
              state: FindState.FINISH,
              source: this,
              rawQuery: this._state?.query ?? '',
              searchResults: this.searchResults,
            })
          })
        }
      }
    )
    this._eventBus._on('updatefindcontrolstate', (payload: DispatchPayload) => {
      if (!this._state) {
        return
      }
      this._searchCallback.forEach((callback) => {
        callback(payload)
      })
    })
  }

  private _formateSearchResults() {
    const list: {
      page: number
      items: SearchResultItem[]
    }[] = []
    let prevOffset: {
      pageIdx: number | null
      matchIdx: number | null
    } = {
      pageIdx: 0,
      matchIdx: null,
    }
    let total = 0
    this._pageMatches?.forEach((onePageMatches, pageIdx) => {
      if (onePageMatches.length <= 0) {
        return
      }
      total += onePageMatches.length
      const pageContent: string = this._pageContents?.[pageIdx]
      const pageDiffs: number[][] = this._pageDiffs?.[pageIdx] || []
      // console.log(11111, pageIdx, pageContent, onePageMatches, this._pageDiffs?.[pageIdx])
      let nextStart = 0
      let prevEnd = 0

      let curDiffIdx = -1
      let nextDiffIdx = -1

      const pageResultItem: {
        page: number
        items: SearchResultItem[]
      } = {
        page: pageIdx,
        items: [],
      }

      onePageMatches.forEach((match: any, idx: number) => {
        const len = this._pageMatchesLength![pageIdx][idx]

        while (curDiffIdx <= pageDiffs.length - 1) {
          if (pageDiffs[curDiffIdx + 1][0] >= match) {
            break
          }
          curDiffIdx += 1
        }

        const diff = pageDiffs[curDiffIdx] || [0, 0]

        match += -diff[1]

        const nextMatch = onePageMatches[idx + 1] || pageContent.length

        nextDiffIdx = curDiffIdx

        while (nextDiffIdx <= pageDiffs.length - 1) {
          if (pageDiffs[nextDiffIdx + 1][0] >= nextMatch) {
            break
          }
          nextDiffIdx += 1
        }

        const nextDiff = pageDiffs[nextDiffIdx] || [0, 0]

        nextStart = nextMatch - nextDiff[1]

        let start = Math.max(0, match - 100)
        start = Math.max(prevEnd, start)

        let end = Math.min(match + len + 100, pageContent.length)
        end = Math.min(nextStart, end)

        const str = pageContent.substring(start, end)
        const result: SearchResultItem = {
          whole: str,
          content: [],
          offset: {
            pageIdx,
            matchIdx: idx,
          },
          prevOffset: {
            ...prevOffset,
          },
        }
        prevOffset = {
          pageIdx,
          matchIdx: idx,
        }
        if (match > start) {
          const preStr = str.substring(0, match - start)
          const preDot = preStr.lastIndexOf('.')
          result.content.push({
            matched: false,
            str: (preDot === -1 ? '...' : '') + preStr.substring(preDot + 1),
          })
        }

        result.content.push({
          matched: true,
          str: str.substring(match - start, match - start + len),
        })

        if (match - start + len < end) {
          const lastStr = str.substring(match - start + len)
          const lastDot = lastStr.indexOf('.')
          result.content.push({
            matched: false,
            // str: lastStr
            str:
              lastStr.substring(
                0,
                lastDot === -1 ? lastStr.length : lastDot + 1
              ) + (lastDot === -1 ? '...' : ''),
          })
          prevEnd = lastDot === -1 ? end : end - (lastStr.length - lastDot)
        } else {
          prevEnd = end
        }
        pageResultItem.items.push(result)
      })
      list.push(pageResultItem)
    })

    return {
      list,
      total,
    }
  }

  addTwinkling(element: HTMLDivElement) {
    const cls = 'figure-border-animation'
    const reg = new RegExp('(\\s|^)' + cls + '(\\s|$)')
    const hasClass = element.className.match(reg)
    if (!hasClass) {
      element.className += element.className ? ' ' : '' + cls
    }
    window.setTimeout(() => {
      element.className = element.className.replace(reg, ' ')
    }, 3000)
  }

  setHighlightColor(v: string) {
    this.highlightColor = v;
  }

  toggleHighlightTwinkling(v: boolean) {
    this.highlightTwinklingEnabled = v ?? !this.highlightTwinklingEnabled;
  }

  public addSearchCallback(callback: any) {
    this._searchCallback.push(callback)
  }

  public close() {
    this._eventBus.dispatch('findbarclose', {})
    this._pageMatches = []
    this._pageMatchesLength = []
    // Currently selected match.
    this._selected = {
      pageIdx: -1,
      matchIdx: -1,
    }
    // Where the find algorithm currently is in the document.
    this._offset = {
      pageIdx: null,
      matchIdx: null,
      wrapped: false,
    }
    this._state = null
    this.searchResults = {
      total: 0,
      list: [],
    }
  }
}
