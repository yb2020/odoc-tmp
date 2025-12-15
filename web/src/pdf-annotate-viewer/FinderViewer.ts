import { PDFFindControllerOptions } from '@idea/pdfjs-dist/types/web/pdf_find_controller'
import FinderController, {
  DispatchPayload,
  FindState,
  InnerFinderEvent,
  SearchResult,
} from './FinderController'
import { createContainer } from './utils/ui'
import cbem from './css/style'
import { PDFDocumentProxy } from '@idea/pdfjs-dist/types/src/display/api'
import PerfectScrollbar from 'perfect-scrollbar'
import scrollIntoView from 'scroll-into-view-if-needed'
import { LoadingBem, loadingClassname, LoadingSVG } from './loading'
import debounce from 'lodash-es/debounce'
import trim from 'lodash-es/trim'
import I18nController from './I18nController'

interface FinderViewerOptions {
  pdfDocument: PDFDocumentProxy
  i18n: I18nController
}

const FinderBem = cbem('finder')

const TitleBem = FinderBem('title')

const CloseIconBem = FinderBem('close-icon')

const InputWrapBem = FinderBem('input-wrap')

const ListWrapBem = FinderBem('list-wrap')
const TotalBem = FinderBem('total')
const NumBem = FinderBem('num')
const SelectedBem = FinderBem('selected')

const classname = 'pdf-finder-container'

const styleTag = document.createElement('style')
styleTag.textContent = `
  .${classname} {
    height: 100%;
    width: 0px;
    overflow: hidden;
    background: var(--site-theme-bg-secondary, #383a3d);
    transition: transform 0.3s cubic-bezier(0.7, 0.3, 0.1, 1), height 0s ease 0.3s,
      width 0.3s ease;
  }

  .${classname} > * {
    transition: transform 0.25s cubic-bezier(0.7, 0.3, 0.1, 1),
      box-shadow 0.3s cubic-bezier(0.7, 0.3, 0.1, 1);
  }

  .${FinderBem} {
    display: flex;
    flex-direction: column;
    color: var(--site-theme-text-primary, #f0f2f5);
    &,
    * {
      box-sizing: border-box;
    }

    .${TitleBem} {
      padding: 14px 24px;
      font-size: 14px;
      line-height: 20px;
    }

    .${CloseIconBem} {
      position: absolute;
      top: 14px;
      right: 24px;
      cursor: pointer;
      color: var(--site-theme-text-primary, #f0f2f5);
    }

    .${InputWrapBem} {
      padding: 14px 24px;
      input {
        height: 42px;
        width: 100%;
        padding: 0 12px;
        border-radius: 4px;
        border: 1px solid var(--site-theme-border-color, #42454a);
        outline: none;
        background: var(--site-theme-bg-mute, #42454a);
        color: var(--site-theme-text-primary, #f0f2f5);
        &::placeholder {
          color: var(--site-theme-text-tertiary, #9e9e9e);
        }
      }
    }

    .${ListWrapBem} {
      flex: 1;
      overflow: auto;
      position: relative;
      ul {
        margin: 0;
        padding: 0;
        list-style: none;
      }
      li {
        padding: 6px 24px;
        cursor: pointer;
        &:hover {
          background-color: var(--site-theme-background-hover, #42454a);
        }
        &.${SelectedBem} {
          background-color: var(--site-theme-background-hover, #42454a);
        }
        span {
          margin-right: 8px;
        }
      }

      .${TotalBem} {
        margin: 6px 0 10px;
        padding: 0 24px;
      }
      .${NumBem} {
        margin: 10px 0;
        padding: 0 24px;
      }
    }
  }
`

if (typeof document !== 'undefined') {
  document.head.appendChild(styleTag)
}

export interface GotoMatchPayload {
  selected: {
    pageIdx: number
    matchIdx: number
  }
}

const DEFAULT_FINDER_WIDTH = 400

export default class FinderViewer {
  private pdfFinderController
  private wrapper?: HTMLDivElement
  private container?: HTMLDivElement
  private width = DEFAULT_FINDER_WIDTH
  private inputDom?: HTMLInputElement
  private searchResultWrapDom?: HTMLDivElement
  private loadingDom?: HTMLDivElement
  private lastTogglePanelTime = 0
  private i18n: I18nController

  constructor(
    { pdfDocument, i18n }: FinderViewerOptions,
    finderOptions: PDFFindControllerOptions
  ) {
    this.pdfFinderController = new FinderController(finderOptions)
    this.pdfFinderController.setDocument(pdfDocument)

    this.i18n = i18n
  }

  private loading(show: boolean) {
    if (this.loadingDom) {
      this.loadingDom.style.display = show ? 'flex' : 'none'
    }
  }

  private cleanResult() {
    if (this.searchResultWrapDom) {
      this.searchResultWrapDom.innerHTML = ''
    }
  }

  private renderResult({
    keyword,
    searchResult,
  }: {
    keyword: string
    searchResult: SearchResult
  }) {
    this.cleanResult()
    this.loading(false)
    const totalWrap = document.createElement('div')
    totalWrap.className = TotalBem
    totalWrap.innerHTML = this.i18n.t('message.totalMatches', {
      num: searchResult.total,
      keyword: `"${keyword}"`,
    })
    this.searchResultWrapDom?.appendChild(totalWrap)
    const listWrap = document.createElement('div')
    searchResult.list.forEach((page) => {
      const pageDiv = document.createElement('div')
      const pageNum = document.createElement('div')
      pageNum.className = NumBem
      pageNum.innerText = this.i18n.t('viewer.page', { num: page.page + 1 })
      pageDiv.appendChild(pageNum)
      const listUl = document.createElement('ul')
      pageDiv.appendChild(listUl)
      page.items.forEach((item, idx) => {
        const li = document.createElement('li')
        li.setAttribute('data-page', `${page.page}`)
        li.setAttribute('data-match', `${idx}`)
        li.setAttribute(
          'data-prev-offset',
          `${item.prevOffset.pageIdx}-${item.prevOffset.matchIdx}`
        )
        const strs: string[] = item.content.map((content) => {
          if (content.matched) {
            return `<span>${content.str}</span>`
          }
          return content.str
        })
        li.innerHTML = strs.join('')
        listUl.appendChild(li)
      })
      listWrap.appendChild(pageDiv)
    })

    this.searchResultWrapDom?.appendChild(listWrap)
  }

  private initWrapper() {
    if (!this.container) {
      return
    }
    this.container.innerHTML = `
      <div class="${TitleBem}">${this.i18n.t('viewer.findInDocument')}</div>
    `

    const closeIcon = document.createElement('span')

    closeIcon.classList.add(CloseIconBem)

    closeIcon.innerHTML = `
      <svg fill="currentColor" width="24px" height="24px" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <g data-name="Layer 2">
          <g data-name="close">
            <rect width="24" height="24" transform="rotate(180 12 12)" opacity="0"/><path d="M13.41 12l4.3-4.29a1 1 0 1 0-1.42-1.42L12 10.59l-4.29-4.3a1 1 0 0 0-1.42 1.42l4.3 4.29-4.3 4.29a1 1 0 0 0 0 1.42 1 1 0 0 0 1.42 0l4.29-4.3 4.29 4.3a1 1 0 0 0 1.42 0 1 1 0 0 0 0-1.42z"/>
          </g>
        </g>
      </svg>
    `

    this.container.appendChild(closeIcon)

    closeIcon.addEventListener('click', () => {
      this.toggleSearchPanel()
    })

    const inputWrap = document.createElement('div')
    inputWrap.className = InputWrapBem
    const input = document.createElement('input')
    inputWrap.appendChild(input)

    this.container.appendChild(inputWrap)
    const searchResultContainer = document.createElement('div')
    searchResultContainer.classList.add(ListWrapBem)
    const searchResultWrap = document.createElement('div')
    searchResultContainer.appendChild(searchResultWrap)
    this.container.appendChild(searchResultContainer)

    new PerfectScrollbar(searchResultWrap, {
      wheelSpeed: 1,
      wheelPropagation: true,
      minScrollbarLength: 5,
    })

    const loading = document.createElement('div')
    loading.innerHTML = LoadingSVG(25)
    loading.className = `${LoadingBem()} ${loadingClassname}`
    loading.style.display = 'none'
    searchResultContainer.appendChild(loading)

    this.searchResultWrapDom = searchResultWrap
    this.inputDom = input
    this.loadingDom = loading
  }

  private bindEvent() {
    if (!this.searchResultWrapDom) {
      return
    }
    const goToMatch = (curSelected: { pageIdx: number; matchIdx: number }) => {
      console.log('dispatch gotomatch', curSelected)
      setTimeout(() => {
        const li = this.searchResultWrapDom?.querySelector(
          `li[data-page="${curSelected?.pageIdx}"][data-match="${curSelected?.matchIdx}"]`
        ) as HTMLLIElement

        if (li) {
          this.searchResultWrapDom
            ?.querySelectorAll(SelectedBem)
            .forEach((dom) => {
              dom.classList.remove(SelectedBem)
            })
          li.classList.add(SelectedBem)
          scrollIntoView(li, {
            scrollMode: 'if-needed',
            block: 'nearest',
            inline: 'nearest',
          })
          this.pdfFinderController._eventBus.dispatch(
            InnerFinderEvent.GOTO_MATCH,
            { selected: curSelected } as GotoMatchPayload
          )
        }
      }, 0)
    }

    this.pdfFinderController.addSearchCallback((data: DispatchPayload) => {
      if (data.state === FindState.WRAPPED) {
        this.loading(false)
        // message.info('已到达文档底部，从顶部继续搜索');
        return
      }
      if (data.state === FindState.PENDING) {
        this.loading(true)
        return
      }

      if (data.state === FindState.FINISH) {
        this.renderResult({
          keyword: data.rawQuery,
          searchResult: data.source.searchResults,
        })

        this.loading(false)
        return
      }
      if (data.state === FindState.FOUND) {
        this.loading(false)
      }
    })

    const doSearch = debounce((e: KeyboardEvent) => {
      const value = trim(this.inputDom?.value)
      if (!value) {
        this.pdfFinderController.close()
        this.cleanResult()
        return
      }
      if (e.code === 'Enter') {
        return doSearchAgain()
      }
      const state = this.pdfFinderController._state
      if (state && state.query === value) {
        return
      }
      this.pdfFinderController._eventBus.dispatch('find', {
        query: value,
        type: '',
        phraseSearch: true,
        highlightAll: true,
        jumpToMatch: false,
      })
    }, 200)
    const doSearchAgain = (preSelected?: {
      pageIdx: null | number
      matchIdx: null | number
    }) => {
      if (preSelected) {
        // 通过offset指定上一处匹配
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        this.pdfFinderController._offset = {
          ...preSelected,
          wrapped: false,
        }
      }

      // 触发页面跳转及高亮渲染
      this.pdfFinderController._eventBus.dispatch('find', {
        query: this.inputDom?.value,
        type: 'again',
        phraseSearch: true,
        highlightAll: true,
        jumpToMatch: true,
      })
      // 触发高亮定位
      goToMatch(this.pdfFinderController._selected!)
    }

    this.inputDom?.addEventListener('keyup', doSearch)
    this.searchResultWrapDom.addEventListener('click', (e) => {
      const li = (e.target as HTMLElement).closest('li') as HTMLLIElement
      if (!li) {
        return
      }

      const prevOffset = li.getAttribute('data-prev-offset')?.split('-')
      const pageIdx = parseInt(prevOffset?.[0] ?? '')
      const matchIdx = parseInt(prevOffset?.[1] ?? '')

      doSearchAgain({
        pageIdx,
        // -1为了让首页的第一个顺利匹配
        matchIdx: !isNaN(matchIdx) ? matchIdx : -1,
      })
    })
  }

  toggleSearchPanel() {
    if (!this.wrapper || !this.container || !this.inputDom) {
      throw Error('search panel not init')
    }
    const now = Date.now()
    let isOpen = false
    if (this.wrapper.clientWidth) {
      // 关闭之前记录一下宽度，下次打开使用这个宽度
      if (now - this.lastTogglePanelTime > 500) {
        this.width = this.wrapper.clientWidth
      }
      this.wrapper.style.width = '0px'
      this.container.style.transform = `translateX(${this.width}px)`
      this.cleanResult()
      this.inputDom.value = ''
      this.pdfFinderController.close()
      isOpen = false
    } else {
      this.wrapper.style.width = `${this.width}px`
      this.container.style.transform = ''
      this.inputDom?.focus()
      isOpen = true
    }
    this.lastTogglePanelTime = now
    return isOpen
  }

  getPdfFinderController() {
    return this.pdfFinderController
  }

  setWrapper(wrapper: HTMLDivElement) {
    this.wrapper = wrapper
    this.container = createContainer(wrapper)
    this.wrapper.classList.add(classname)
    this.container.classList.add(FinderBem())
    this.initWrapper()
    this.bindEvent()
  }

  destroy() {
    this.pdfFinderController.setDocument(null as any)
    this.container?.remove()
  }
}
