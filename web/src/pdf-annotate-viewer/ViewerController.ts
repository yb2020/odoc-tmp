import * as pdfjsViewer from '@idea/pdfjs-dist/web/pdf_viewer'
import * as pdfjsLib from '@idea/pdfjs-dist'
import {
  PDFDocumentLoadingTask,
  DocumentInitParameters,
} from '@idea/pdfjs-dist/types/src/display/api'
import {
  LoadingBem,
  LoadingProgressBem,
  LoadingSVG,
} from './loading'
import {
  PageSelectText,
  PDFViewerScale,
  ViewerEvent,
  PageMouseEventPayload,
} from './type'
import { DocumentViewer } from './DocumentViewer'
import ThumbnailViewer from './ThumbnailViewer'
import { PDFRenderingQueue as PDFRenderingQueueType } from '@idea/pdfjs-dist/types/web/pdf_rendering_queue'
import ScrollController from './ScrollController'
import FinderViewer from './FinderViewer'
import hotkeys from 'hotkeys-js'
import DownloadController from './DownloadController'
import SANDBOX_BUNDLE_SRC from '@idea/pdfjs-dist/build/pdf.sandbox.js?url'
import { createL10nController } from './L10nController'
import { PDFViewerOptions as BasePDFViewerOptions } from '@idea/pdfjs-dist/types/web/pdf_viewer'
import I18nController from './I18nController'
import WordMatchController from './WordMatchController'
import { Language } from 'go-sea-proto/gen/ts/lang/Language';

type PDFViewerOptions = Omit<
  BasePDFViewerOptions,
  'container' | 'eventBus' | 'linkService' | 'l10n'
>

export interface ViewerControllerOptions {
  pdfDocumentParams: DocumentInitParameters & {
    retryUrl?: string
    language?: Language
  }
  containers: {
    documentWrapper: HTMLDivElement
  }
}

export class ViewerController {
  private documentWrapper: HTMLDivElement

  private pdfDocument?: pdfjsLib.PDFDocumentProxy

  private loadingTask?: PDFDocumentLoadingTask

  private eventBus = new pdfjsViewer.EventBus()

  private pdfLinkService = new ScrollController({
    eventBus: this.eventBus,
    ignoreDestinationZoom: true,
  })

  // (Optionally) enable scripting support.
  private pdfScriptingManager = new pdfjsViewer.PDFScriptingManager({
    eventBus: this.eventBus,
    sandboxBundleSrc: SANDBOX_BUNDLE_SRC,
  })

  private pdfRenderingQueue?: PDFRenderingQueueType

  private documentViewer?: DocumentViewer
  private scrollController?: ScrollController

  private finderViewer?: FinderViewer
  private thumbnailViewer?: ThumbnailViewer

  private loadingWrapper?: HTMLDivElement | null

  private downloadController: DownloadController

  private annotationMode = 1

  private pdfDocumentParams: ViewerControllerOptions['pdfDocumentParams']

  private i18nController?: I18nController

  private _handleHotkeys?: () => void

  private wordMatchController?: WordMatchController

  public constructor({ pdfDocumentParams, containers }: ViewerControllerOptions) {
    
    this.pdfDocumentParams = pdfDocumentParams;

    const { documentWrapper } = containers

    this._createLoadingTask(pdfDocumentParams)

    this.documentWrapper = documentWrapper

    this.bindHotkeys()

    this.downloadController = new DownloadController({
      url: pdfDocumentParams.url?.toString() ?? '',
      pdfScriptingManager: this.pdfScriptingManager,
    })
  }

  private _createLoadingTask(
    pdfDocumentParams: ViewerControllerOptions['pdfDocumentParams']
  ) {
    const loadingTask = pdfjsLib.getDocument({
      cMapUrl: 'https://nuxt.cdn.readpaper.com/pdfjs-dist%402.13.216/cmaps/',
      cMapPacked: true,
      enableXfa: true,
      isEvalSupported: false,
      standardFontDataUrl:
        'https://nuxt.cdn.readpaper.com/pdfjs-dist%402.13.216/standard_fonts/',
      ...pdfDocumentParams,
    })
    this.loadingTask = loadingTask
    this.loadingTask.onProgress = (progressData: {
      loaded: number
      total: number
    }) => {
      const processPercent = progressData.loaded / progressData.total || 0.99
      this.eventBus.dispatch(ViewerEvent.PROGRESS_CHANGE, processPercent)
    }
  }

  addEventListener(
    event: ViewerEvent.PROGRESS_CHANGE,
    handler: (percent: number) => void
  ): void
  addEventListener(
    event: ViewerEvent.TEXT_SELECT,
    handler: (pageTexts: PageSelectText[]) => void
  ): void
  addEventListener(
    event: ViewerEvent.EMPTY_CLICK,
    handler: (payload: { pageNumber: number }) => void
  ): void
  addEventListener(
    event: ViewerEvent.SCALE_CHANGING,
    handler: (payload: {
      scale: number
      presetValue?: string
      source: pdfjsViewer.PDFViewer
    }) => void
  ): void
  addEventListener(
    event: ViewerEvent.PAGE_RENDERED,
    handler: (payload: {
      pageNumber: number
      cssTransform: boolean
      source: pdfjsViewer.PDFPageView
    }) => void
  ): void
  addEventListener(
    event: ViewerEvent.TEXT_LAYER_RENDERED,
    handler: (payload: {
      pageNumber: number
      source: pdfjsViewer.PDFPageView
      numTextDivs: number
      error: Error
    }) => void
  ): void
  addEventListener(
    event: ViewerEvent.PAGES_INIT,
    handler: (payload: { source: pdfjsViewer.PDFViewer }) => void
  ): void
  addEventListener(
    event: ViewerEvent.PAGE_CHANGING,
    handler: (payload: {
      source: pdfjsViewer.PDFViewer
      pageNumber: number
      previous: number
    }) => void
  ): void
  addEventListener(
    event: ViewerEvent.TRIGGER_SCALE_CHANGE,
    handler: (scale: string | number) => void
  ): void
  addEventListener(
    event: ViewerEvent.PAGE_MOUSE_EVENT,
    handler: (payload: PageMouseEventPayload) => void
  ): void
  addEventListener(
    event: ViewerEvent.PDFURL_RETRY,
    handler: (payload: { pdfUrl: string }) => void
  ): void
  addEventListener(event: string, handler: unknown) {
    this.eventBus.on(event, handler as (...args: any[]) => void)
  }

  removeEventListener(event: string, handler: any) {
    this.eventBus.off(event, handler)
  }

  private hideLoading() {
    if (!this.loadingWrapper) {
      return
    }
    this.documentWrapper.parentElement?.removeChild(this.loadingWrapper)
    this.loadingWrapper = null
  }
  private loading() {
    const wrap: HTMLDivElement | null = document.createElement('div')
    wrap.className = `${LoadingBem()}`
    Object.assign(wrap.style, {
      position: 'absolute',
      top: '0',
      bottom: '0',
      left: '0',
      right: '0',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      flexDirection: 'column',
      background: 'rgba(25, 97, 205, 0.2)'
    })
    wrap.innerHTML = LoadingSVG(40)
    const progress = document.createElement('span')
    progress.className = LoadingProgressBem()
    progress.innerText = '0%'
    Object.assign(progress.style, {
      marginTop: '10px',
      color: '#ccc'
    })
    this.eventBus.on(ViewerEvent.PROGRESS_CHANGE, (percent: number) => {
      const num = Math.round(percent * 100)
      progress.innerText = `${num}%`
    })
    wrap.appendChild(progress)
    this.documentWrapper.parentElement?.appendChild(wrap)
    this.loadingWrapper = wrap
  }

  private async _retryLoadUrl(error: unknown) {
    if (!(error instanceof pdfjsLib.InvalidPDFException)) {
      this._createLoadingTask(this.pdfDocumentParams)
      try {
        const pdfDocument = await this.loadingTask!.promise
        this.eventBus.dispatch(ViewerEvent.PDFURL_RETRY, {
          pdfUrl: this.pdfDocumentParams.url,
        })
        return pdfDocument
      } catch (error) {
        if (this.pdfDocumentParams.retryUrl) {
          this.pdfDocumentParams.url = this.pdfDocumentParams.retryUrl
          this._createLoadingTask(this.pdfDocumentParams)
          const pdfDocument = await this.loadingTask!.promise
          this.eventBus.dispatch(ViewerEvent.PDFURL_RETRY, {
            pdfUrl: this.pdfDocumentParams.url,
          })
          return pdfDocument
        }
        throw error
      }
    }
    throw error
  }

  async build(scale: PDFViewerScale, pdfOptions: PDFViewerOptions) {
    if (pdfOptions.annotationMode) {
      this.annotationMode = pdfOptions.annotationMode
    }
    this.loading()

    let pdfDocument: pdfjsLib.PDFDocumentProxy | null = null

    try {
      pdfDocument = await this.loadingTask!.promise
    } catch (error) {
      pdfDocument = await this._retryLoadUrl(error)
    }

    if (!pdfDocument) {
      throw Error('Invalid pdfDocument.')
    }

    this.pdfDocument = pdfDocument
    this.pdfLinkService.setDocument(pdfDocument)

    this.i18nController = new I18nController(
      this.pdfDocumentParams.language ? Language[this.pdfDocumentParams.language] : undefined
    )

    const finderViewer = new FinderViewer(
      {
        pdfDocument: pdfDocument,
        i18n: this.i18nController,
      },
      {
        eventBus: this.eventBus,
        linkService: this.pdfLinkService,
        updateMatchesCountOnProgress: true,
      }
    )

    const documentViewer = new DocumentViewer(
      {
        pdfDocument,
        wrapper: this.documentWrapper,
        scale: scale || '1.0',
      },
      {
        eventBus: this.eventBus,
        linkService: this.pdfLinkService,
        // renderer: 'canvas',
        l10n: createL10nController(
          this.pdfDocumentParams.language ? Language[this.pdfDocumentParams.language] : undefined
        ).getL10n(),
        findController: finderViewer?.getPdfFinderController(),
        imageResourcesPath:
          'https://nuxt.cdn.readpaper.com/pdfjs-dist%402.13.216/images/',
        ...pdfOptions,
      }
    )
    this.pdfLinkService.setPdfViewerContainer(documentViewer.container)

    const pdfViewer = documentViewer.getPdfViewer()
    this.pdfLinkService.setViewer(pdfViewer)

    this.pdfRenderingQueue = pdfViewer.renderingQueue

    this.documentViewer = documentViewer
    this.scrollController = this.pdfLinkService

    this.finderViewer = finderViewer

    this.downloadController.setPDFDocument(pdfDocument)

    this.hideLoading()

    this.wordMatchController = new WordMatchController(
      this.pdfDocument,
      pdfViewer
    )

    return {
      documentViewer,
      finderViewer,
      scrollController: this.pdfLinkService,
    }
  }

  enableThumbnailViewer(thumbnailWrapper: HTMLDivElement) {
    if (!this.pdfDocument || !this.pdfRenderingQueue) {
      throw new Error('cannot enableThumbnailViewer before build is called')
    }
    const thumbnailViewer = new ThumbnailViewer({
      wrapper: thumbnailWrapper,
      eventBus: this.eventBus,
      pdfDocument: this.pdfDocument,
      renderingQueue: this.pdfRenderingQueue,
      linkService: this.pdfLinkService,
      currentPageNumber:
        this.documentViewer?.getPdfViewer()?.currentPageNumber || 1,
      annotationMode: this.annotationMode,
    })

    this.pdfRenderingQueue.isThumbnailViewEnabled = true
    this.pdfRenderingQueue.setThumbnailViewer(thumbnailViewer as any)

    this.thumbnailViewer = thumbnailViewer

    return thumbnailViewer
  }

  getDocumentViewer() {
    if (!this.documentViewer) {
      throw new Error('please call build before getDocumentViewer')
    }
    return this.documentViewer
  }

  getScrollController() {
    if (!this.scrollController) {
      throw new Error('please call build before getScrollController')
    }
    return this.scrollController
  }

  getFinderViewer() {
    if (!this.finderViewer) {
      throw new Error('please call build before getFinderView')
    }
    return this.finderViewer
  }

  getThumbnailViewer() {
    if (!this.thumbnailViewer) {
      throw new Error(
        'please call enableThumbnailViewer before getThumbnailViewer'
      )
    }
    return this.thumbnailViewer
  }

  isEnable() {
    return !!this.documentWrapper.offsetParent
  }

  bindHotkeys() {
    console.log('bindHotkeys')
    const handleHotkeys = () => {
      console.log('copy')
      if (!this.isEnable()) {
        return
      }
      setTimeout(() => {
        this.getDocumentViewer().copyToClipboard()
      }, 0);
    }
    this._handleHotkeys = handleHotkeys.bind(this)
    hotkeys('ctrl+c,cmd+c', this._handleHotkeys)

    // hotkeys('ctrl+f,cmd+f', (event) => {
    //   if (!this.isEnable()) {
    //     return
    //   }
    //   event.preventDefault()
    //   this.getFinderViewer().toggleSearchPanel()
    // })
  }

  redraw() {
    this.documentViewer?.clear()
    this.thumbnailViewer?.clear()
    this.documentViewer?.getPdfViewer()?.forceRendering(null)
  }

  getDownloadController() {
    return this.downloadController
  }

  destroy() {
    this.hideLoading()
    this.documentViewer?.destroy()
    this.thumbnailViewer?.clear()
    this.finderViewer?.destroy()
    if (this._handleHotkeys) {
      hotkeys.unbind('ctrl+c,cmd+c', this._handleHotkeys)
    }
  }

  getPdfDocument() {
    return this.pdfDocument
  }

  getWordMatchController() {
    if (!this.wordMatchController) {
      throw new Error('please call build before getWordMatchController')
    }
    return this.wordMatchController
  }
}
