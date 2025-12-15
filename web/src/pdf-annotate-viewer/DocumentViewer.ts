import * as pdfjsViewer from '@idea/pdfjs-dist/web/pdf_viewer'
import * as pdfjsLib from '@idea/pdfjs-dist'
import TextLayerBuilder from './TextLayerBuilder'
import PerfectScrollbar from 'perfect-scrollbar'
import { MousePoint, PageSelectText, PDFViewerScale, PointOnPdf, ViewerEvent } from './type'
import WheelController from './WheelController'
import {
  PDFDocumentProxy,
  PDFViewerOptions,
} from '@idea/pdfjs-dist/types/web/pdf_viewer'
import { createContainer } from './utils/ui'
import { InnerFinderEvent } from './FinderController'
import { GotoMatchPayload } from './FinderViewer'
import { classname, ViewerBem } from './css/DocumentViewerStyle'
import { copyTextToClipboard, getPosition } from './utils/dom'
import { normalize } from './utils/text'
import { BaseMousePoint } from './type'
// import { throttle } from 'lodash-es';
import { PDFPageViewOptions } from '@idea/pdfjs-dist/types/web/pdf_page_view'
import PDFPageView from './PDFPageView'
import TextHighlighter from './TextHighlighter'
import throttle from 'lodash-es/throttle'
import { detectCtrlPressing, JS_IGNORE_MOUSE_OUTSIDE } from './utils/mouse'

export interface DocumentViewOptions {
  scale: PDFViewerScale
  pdfDocument: PDFDocumentProxy
  wrapper: HTMLDivElement
}

type MPDFViewerOptions = Omit<PDFViewerOptions, 'container'>

type WrapSizes = { wrapWidth: number; wrapHeight: number }[]

// const DEFAULT_SCALE_DELTA = 1.1;
// const MIN_SCALE = 0.5;
// const MAX_SCALE = 4.0;
class PDFViewer extends pdfjsViewer.PDFViewer {
  createPDFPageView(options: PDFPageViewOptions) {
    return new PDFPageView(options)
  }
}

const DEFAULT_BOUND_BUFF = 6

export class DocumentViewer {
  public container: HTMLDivElement
  private pdfViewer
  private eventBus
  private boundBuff = DEFAULT_BOUND_BUFF
  private selectionEnabled = true
  private selectedTexts = ''
  public isSelecting = () => Boolean(this.deltaPoint)

  private wrapSizes: WrapSizes = []

  private isPageInit = false

  private ps?: PerfectScrollbar

  public constructor(
    { 
      wrapper,
      scale,
      pdfDocument,
    }: DocumentViewOptions,
    pdfOptions: MPDFViewerOptions
  ) {
    const container = createContainer(wrapper)
    container.classList.add(classname)

    const viewer: HTMLDivElement = document.createElement('div')
    viewer.classList.add(ViewerBem())
    viewer.classList.add('pdfViewer')

    container.append(viewer)

    const pdfViewer = new PDFViewer({
      ...pdfOptions,
      container,
      textLayerMode: 2,
    })
    pdfViewer.setDocument(pdfDocument)
    pdfViewer.setDefaultCacheSize(5)

    this.container = container

    if (localStorage.getItem('debugPerfectScrollbar')) {
      this.setPerfectScrollbar()
    } else {
      this.container.style.overflow = 'auto'
    }

    this.pdfViewer = pdfViewer
    this.eventBus = pdfOptions.eventBus

    this.start(scale);
    
    // this.bindHotkeys();
    this.initMouse();
  }

  private updateWrapSizes() {
    this.wrapSizes.length = 0
    const pageNums = this.pdfViewer.pagesCount
    for (let i = 0; i < pageNums; i += 1) {
      const pdfPageView = this.pdfViewer.getPageView(i)
      if (pdfPageView) {
        const viewport = pdfPageView.div.getBoundingClientRect()
        this.wrapSizes.push({
          wrapWidth: viewport.width,
          wrapHeight: viewport.height,
        })
      } else {
        const wrap = this.container.querySelector(
          `.page[data-page-number="${i + 1}"]`
        ) as HTMLDivElement
        const viewport = wrap.getBoundingClientRect()
        this.wrapSizes.push({
          wrapWidth: viewport.width,
          wrapHeight: viewport.height,
        })
      }
    }
  }

  private getPointOnPdfView = (p: BaseMousePoint): PointOnPdf | null => {
    // 添加 pdfViewer 和 viewer 的安全检查
    if (!this.pdfViewer || !this.pdfViewer.viewer) {
      console.warn('PDF viewer not initialized yet');
      return null;
    }
    
    if (!this.wrapSizes.length || !this.isPageInit) {
      this.updateWrapSizes()
    }
    if (!this.wrapSizes.length) {
      return null
    }
    const wx = p.x
    const wy = p.y

    const totalOffset = getPosition(this.container)

    const totalTop = this.container.scrollTop + wy - totalOffset.top

    let top = 0

    let cur = 0

    let lastTop = 0

    for (let i = 0; i < this.wrapSizes.length; i += 1) {
      top += this.wrapSizes[i].wrapHeight

      if (top >= totalTop) {
        cur = i
        break
      }

      lastTop += this.wrapSizes[i].wrapHeight
    }

    const wv = this.wrapSizes[cur]

    const curPageView: PDFPageView = this.pdfViewer.getPageView(cur)

    const pv = curPageView.pdfPage?.getViewport({
      scale:
        this.pdfViewer.currentScale *
        (pdfjsLib.PixelsPerInch as any).PDF_TO_CSS_UNITS,
    })
    if (!pv) {
      return null
    }

    const totalLeft = this.container.scrollLeft + wx - totalOffset.left

    // 再次检查 pdfViewer.viewer 是否可用（防止竞态条件）
    if (!this.pdfViewer || !this.pdfViewer.viewer) {
      console.warn('PDF viewer not available for getBoundingClientRect');
      return null;
    }

    const viewerRect = (this.pdfViewer.viewer as HTMLDivElement).getBoundingClientRect();
    const pageRect = curPageView.div.getBoundingClientRect();
    const pageOffsetLeft = pageRect.left - viewerRect.left;

    const left =
      totalLeft - pageOffsetLeft - (wv.wrapWidth - pv.width) / 2

    if (left >= 0 && left <= pv.width) {
      // 计算出当前鼠标位置在哪一页上面
      const pageOffsetTop = pageRect.top - viewerRect.top;
      const top = totalTop - pageOffsetTop - (wv.wrapHeight - pv.height) / 2

      return {
        left,
        top,
        cur,
        pv,
      }
    }
    return null
  }

  private setPerfectScrollbar() {
    // document.body.addEventListener(
    //   'mousemove',
    //   throttle(
    //     (e: MouseEvent) => {
    //       const pos = getPosition(this.container);
    //       const distanceX = pos.left + this.container.offsetWidth - e.pageX;
    //       if (Math.abs(distanceX) < 20) {
    //         this.container.classList.add('large-scrollbar-y');
    //       } else {
    //         this.container.classList.remove('large-scrollbar-y');
    //       }

    //       const distanceY = pos.top + this.container.offsetHeight - e.pageY;
    //       if (Math.abs(distanceY) < 20) {
    //         this.container.classList.add('large-scrollbar-x');
    //       } else {
    //         this.container.classList.remove('large-scrollbar-x');
    //       }
    //     },
    //     150,
    //     {
    //       leading: false,
    //     }
    //   )
    // );

    const ps = new PerfectScrollbar(this.container, {
      wheelSpeed: 0.7,
      wheelPropagation: true,
      minScrollbarLength: 5,
      suppressScrollY: false,
      suppressScrollX: false,
    })
    window.addEventListener('resize', () => {
      ps.update()
    })

    this.ps = ps

    // return ps;
  }

  private start(scale: number | string) {

    const throttleScroll = throttle(() => {
      this.ps?.update()
    }, 300)

    this.eventBus.on('pagesinit', () => {
      // We can use pdfViewer now, e.g. let's change default scale.
      this.pdfViewer.currentScaleValue = `${scale}`
      // this.containerScrollbar.update();
      /**
       * 需要延迟一下，否则div的大小是不对的
       */
      this.wrapSizes.length = 0
      setTimeout(() => {
        this.isPageInit = true
      }, 5000)
    })
    this.eventBus.on(
      'scalechanging',
      ({ scale }: { scale: number; presetValue?: string }) => {
        this.boundBuff = Math.min(DEFAULT_BOUND_BUFF * scale, 16)
        setTimeout(() => {
          // this.containerScrollbar.update();
          throttleScroll()
        }, 200)
        this.wrapSizes.length = 0
      }
    )

    // this.eventBus.on(InnerFinderEvent.HIGHLIGHT, () => {
    //   const pageNums = this.pdfViewer.pagesCount;
    //   for (let i = 0; i < pageNums; i += 1) {
    //     const pdfPageView = this.pdfViewer.getPageView(i);
    //     (pdfPageView?.textLayer as TextLayerBuilder)?.highlight();
    //   }
    // });

    this.eventBus.on('findbarclose', () => {
      const pageNums = this.pdfViewer.pagesCount
      for (let i = 0; i < pageNums; i += 1) {
        const pdfPageView = this.pdfViewer.getPageView(i)
        const textHighlighter: TextHighlighter | null =
          pdfPageView?.textLayer?.highlighter
        if (textHighlighter) {
          textHighlighter.clearMatches(true)
        }
      }
    })

    this.eventBus.on(
      InnerFinderEvent.GOTO_MATCH,
      ({ selected }: GotoMatchPayload) => {
        const pageNums = this.pdfViewer.pagesCount
        for (let i = 0; i < pageNums; i += 1) {
          const pdfPageView = this.pdfViewer.getPageView(
            i
          ) as pdfjsViewer.PDFPageView
          const textHighlighter: TextHighlighter | null =
            pdfPageView?.textLayer?.highlighter
          if (textHighlighter) {
            textHighlighter.cancelMatch()
          }
        }

        const curTextHighlighter: TextHighlighter | null =
          this.pdfViewer.getPageView(selected.pageIdx)?.textLayer?.highlighter
        if (curTextHighlighter) {
          const y = curTextHighlighter.goToMatch(selected.matchIdx)
          if (y !== null) {
            this.container.scrollTop = y
          }
        }
      }
    )
  }

  public changeScale(scale: string | number) {
    this.pdfViewer.currentScaleValue = `${scale}`;
    this.eventBus.dispatch(ViewerEvent.TRIGGER_SCALE_CHANGE, scale);
  }

  public getPoint = (event: MouseEvent) => {
    if (!this.selectionEnabled) {
      return null
    }

    return this.getPointOnPdfView({
      x: event.pageX,
      y: event.pageY,
    })
  }

  private deltaPoint: null | MousePoint = null;

  private initMouse() {
    document.body.addEventListener('pointerdown', (event) => {
      if (
        event.target &&
        !this.container.contains(event.target as HTMLElement) &&
        this.container.offsetParent &&
        !(event.target as HTMLElement).closest('.' + JS_IGNORE_MOUSE_OUTSIDE)
        // !detectCtrlPressing(e)
      ) {
        // outside
        this.clearSelection()
      }
    })

    this.container.addEventListener('pointermove', (event) => {
      const point = this.getPoint(event);
      if (!point) {
        return;
      }

      const { left, top, cur, pv } = point;

      this.setContainerCursor(cur, left, top);

      this.eventBus.dispatch(ViewerEvent.PAGE_MOUSE_EVENT, {
        event,
        point: {
          left,
          top,
          pageIndex: cur,
          viewport: pv,
        },
      });
    });

  }

  public onMouseDown = (point: PointOnPdf | null, ctrlPressing: boolean) => {    
    this.deltaPoint = null;

    //console.log('down', point);
    if (!point) {
      return;
    }

    const { left, top, cur } = point;

    const curPageView = this.pdfViewer.getPageView(cur);
    const textLayer: TextLayerBuilder = curPageView.textLayer
    const isInside = textLayer.isInTextBound(
      {
        x: left,
        y: top,
      },
      this.boundBuff
    );
    if (isInside) {
      this.deltaPoint = {
        x: left,
        y: top,
        index: cur,
      };
    }

    if (ctrlPressing) {
      textLayer.stage();
    } else {
      this.clearSelection();
    }
  }

  private setContainerCursor = (cur: number, left: number, top: number) => {
    const curTextLayer = this.pdfViewer.getPageView(cur)
      ?.textLayer as TextLayerBuilder;
    const isInside = !!curTextLayer?.isInTextBound(
      {
        x: left,
        y: top,
      },
      this.boundBuff
    );
    if (isInside) {
      this.container.style.cursor = 'text';
    } else {
      this.container.style.cursor = 'default';
    }
  }

  private textLayerSelect = (
    cur: number,
    left: number,
    top: number,
    pv: NonNullable<ReturnType<typeof this.getPointOnPdfView>>['pv']
  ) => {
    const oldPoint = this.deltaPoint as MousePoint;
    const newPoint: MousePoint = { x: left, y: top, index: cur };
    const [point1, point2] =
      newPoint.index >= oldPoint.index
        ? [oldPoint, newPoint]
        : [newPoint, oldPoint];

    const getTextlayer = (index: number): TextLayerBuilder | null => {
      return this.pdfViewer.getPageView(index)?.textLayer;
    };

    const result = [point1.index, point2.index];

    if (point1.index === point2.index) {
      const textLayer = getTextlayer(point1.index);
      if (!textLayer) {
        return;
      }

      textLayer.select(point1, point2, this.boundBuff);
      textLayer.crossing = false;
      return result;
    }

    const pageStartPoint = { x: 0, y: 0 };
    const pageEndPoint = { x: pv.width, y: pv.height };

    {
      const textLayer = getTextlayer(point1.index);
      if (!textLayer) {
        return;
      }

      textLayer.select(point1, pageEndPoint, this.boundBuff);
      textLayer.crossing = true;
    }

    for (let index = point1.index + 1; index < point2.index; index += 1) {
      const textLayer = getTextlayer(index);
      if (!textLayer) {
        return;
      }

      textLayer.select(pageStartPoint, pageEndPoint, this.boundBuff);
      textLayer.crossing = true;
    }

    {
      const textLayer = getTextlayer(point2.index);
      if (!textLayer) {
        return;
      }

      textLayer.select(pageStartPoint, point2, this.boundBuff);
      textLayer.crossing = true;
    }

    return result;
  }

  public onMouseMove = (point: PointOnPdf | null) => {
    //console.log('move')

    if (!point || !this.deltaPoint) {
      return;
    }

    const { left, top, cur, pv } = point;
    
    this.textLayerSelect(cur, left, top, pv);
  }

  public onMouseUp = (point: PointOnPdf | null) => {
    //console.log('up', point);

    if (!point || !this.deltaPoint) {
      return;
    }

    const { left, top, cur, pv } = point;

    this.textLayerSelect(cur, left, top, pv);

    this.deltaPoint = null;

    const selectedTextRects = this.getSelectedRects();
    if (selectedTextRects.length) {
      this.emitTextSelect(selectedTextRects);
    }
  }

  private getSelectedRects() {
    const selectedTextRects: PageSelectText[] = [];

    for (let i = 0; i < this.pdfViewer.pagesCount; i += 1) {
      const curTextLayer: TextLayerBuilder | undefined =
        this.pdfViewer.getPageView(i)?.textLayer

      if (!curTextLayer) {
        continue
      }

      const textRects = curTextLayer.getTextRects()
      if (!textRects) {
        continue
      }

      selectedTextRects.push({
        pageNum: i + 1,
        ...textRects,
      })
    }

    for (let i = 0; i < selectedTextRects.length - 1; i += 1) {
      const rect = selectedTextRects[i]
      if (!/\s$/.test(rect.text)) {
        rect.text += ' '
      }
    }

    return selectedTextRects
  }

  private emitTextSelect(selectedTextRects: PageSelectText[]) {
    this.eventBus.dispatch(ViewerEvent.TEXT_SELECT, selectedTextRects)

    this.selectedTexts = normalize(
      selectedTextRects.map(({ text }) => text).join('')
    )[0].toString()
  }

  public copyToClipboard() {
    const selectText = window.getSelection && window.getSelection()?.toString();
    //console.log('copyed', selectText, this.selectedTexts);
    if ((!selectText || selectText ==='\n') && this.selectedTexts) {
      copyTextToClipboard(this.selectedTexts)
    }
  }

  // bindHotkeys() {
  //   hotkeys('ctrl+c,cmd+c', () => {
  //     const selectText = window.getSelection && window.getSelection()?.toString();
  //     if (!selectText) {
  //       copyTextToClipboard(selectText || this.selectedTexts);
  //     }

  //   });
  // }

  public onDoubleClick = (e: MouseEvent) => {
    const point = this.getPointOnPdfView({
      x: e.pageX,
      y: e.pageY,
    });
    
    //console.log('dbl', point);

    if (!point) {
      return;
    }
    const { left, top, cur } = point;
    const curTextLayer = this.pdfViewer.getPageView(cur)
      ?.textLayer as TextLayerBuilder | null;
    if (!curTextLayer) {
      return;
    }

    if (detectCtrlPressing(e)) {
      curTextLayer.stage();
    } else {
      this.clearSelection();
    }

    curTextLayer?.selectWord(
      {
        x: left,
        y: top,
      },
      this.boundBuff
    );

    const selectedTextRects = this.getSelectedRects();

    if (selectedTextRects.length) {
      this.emitTextSelect(selectedTextRects);
    } else {
      this.eventBus.dispatch(ViewerEvent.EMPTY_CLICK, {
        pageNumber: cur + 1,
      });
    }
  }

  zoomIn(steps: number | null, scaleFactor?: number) {
    if (this.pdfViewer.isInPresentationMode) {
      return
    }
    this.pdfViewer.increaseScale({
      drawingDelay: 200,
      steps,
      scaleFactor,
    })
    this.eventBus.dispatch(
      ViewerEvent.TRIGGER_SCALE_CHANGE,
      this.pdfViewer.currentScale
    )
  }

  zoomOut(steps: number | null, scaleFactor?: number) {
    if (this.pdfViewer.isInPresentationMode) {
      return
    }
    this.pdfViewer.decreaseScale({
      drawingDelay: 200,
      steps,
      scaleFactor,
    })

    this.eventBus.dispatch(
      ViewerEvent.TRIGGER_SCALE_CHANGE,
      this.pdfViewer.currentScale
    )
  }

  getPdfViewer() {
    return this.pdfViewer
  }

  enableWheelToScale() {
    const wheelController = new WheelController(this)
    wheelController.bindWheelToScale()
  }

  updateSelectionColor(color?: string) {
    this.pdfViewer._pages?.forEach((page) => {
      (page.textLayer as TextLayerBuilder)?.updateRectsColor(color)
    })
  }

  enableSelection(enable: boolean) {
    this.selectionEnabled = enable
    if (!enable) {
      this.container.style.cursor = 'default'
    }
  }

  clearSelection() {
    this.selectedTexts = ''
    this.pdfViewer._pages?.forEach((page) => {
      page.textLayer?.clearTextLayer({ stage: true })
    })
  }

  setSelectedText(text: string) {
    this.selectedTexts = text
  }

  clear() {
    (this.pdfViewer._pages || [])
      .filter((e) => e.renderingState)
      .forEach((e) => e.reset())
  }

  destroy() {
    this.pdfViewer.setDocument(null as any)
  }
}
