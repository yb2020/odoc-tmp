import * as pdfjsViewer from '@idea/pdfjs-dist/web/pdf_viewer';
import type { DocumentViewer } from './DocumentViewer'

const WHEEL_ZOOM_DISABLED_TIMEOUT = 1000 // ms

export default class WheelController {
  private zoomDisabledTimeout = 0
  private pdfDocumentViewer
  private _props = {
    _wheelUnusedFactor: 1,
    _wheelUnusedTicks: 0,
  }
  constructor(pdfDocumentViewer: DocumentViewer) {
    this.pdfDocumentViewer = pdfDocumentViewer
  }

  

  private setZoomDisabledTimeout() {
    if (this.zoomDisabledTimeout) {
      window.clearTimeout(this.zoomDisabledTimeout)
    }
    this.zoomDisabledTimeout = window.setTimeout(() => {
      this.zoomDisabledTimeout = 0
    }, WHEEL_ZOOM_DISABLED_TIMEOUT)
  }

  private _accumulateFactor(previousScale: number, factor: number, prop: '_wheelUnusedFactor') {
    if (factor === 1) {
      return 1;
    }
    // If the direction changed, reset the accumulated factor.
    if ((this._props[prop] > 1 && factor < 1) || (this._props[prop] < 1 && factor > 1)) {
      this._props[prop] = 1;
    }

    const newFactor =
      Math.floor(previousScale * factor * this._props[prop] * 100) /
      (100 * previousScale);
    this._props[prop] = factor / newFactor;

    return newFactor;
  }

  private _accumulateTicks(ticks: number, prop: '_wheelUnusedTicks') {
    if ((this._props[prop] > 0 && ticks < 0) || (this._props[prop] < 0 && ticks > 0)) {
      this._props[prop] = 0;
    }
    this._props[prop] += ticks;
    const wholeTicks = Math.trunc(this._props[prop]);
    this._props[prop] -= wholeTicks;
    return wholeTicks;
  }

  private _centerAtPos(previousScale: number, x: number, y: number) {
    const pdfViewer = this.pdfDocumentViewer.getPdfViewer()
    const scaleDiff = pdfViewer.currentScale / previousScale - 1;
    if (scaleDiff !== 0) {
      const [top, left] = pdfViewer.containerTopLeft;
      pdfViewer.container.scrollLeft += (x - left) * scaleDiff;
      pdfViewer.container.scrollTop += (y - top) * scaleDiff;
    }
  }

  bindWheelToScale() {
    const handler = (evt: WheelEvent) => {
      const pdfViewer = this.pdfDocumentViewer.getPdfViewer()
      const container = this.pdfDocumentViewer.container
      if (!container.offsetParent) {
        // 说明不可见，隐藏状态
        return
      }

      if (evt.target && !container.contains(evt.target as HTMLElement)) {
        return
      }

      if (pdfViewer.isInPresentationMode) {
        return
      }

      const deltaMode = evt.deltaMode
      let { deltaY } = evt;
      if (
        (navigator as any).userAgentData?.platform === 'Windows'
        || navigator.userAgent.includes('Windows')
        || navigator.platform === 'Win32'
      ) {
        deltaY = Math.sign(deltaY) * Math.min(64, Math.abs(deltaY));       
      }

      let scaleFactor = Math.exp(-deltaY / 100)
      const isPinchToZoom =
        (evt.ctrlKey ) &&
        deltaMode === WheelEvent.DOM_DELTA_PIXEL &&
        evt.deltaX === 0 &&
        Math.abs(scaleFactor - 1) < 0.05 &&
        evt.deltaZ === 0
      if (
        isPinchToZoom ||
        evt.ctrlKey ||
        evt.metaKey
      ) {
        evt.preventDefault()
        if (this.zoomDisabledTimeout || document.visibilityState === 'hidden') {
          return
        }
        const previousScale = pdfViewer.currentScale;

        if (isPinchToZoom) {
          scaleFactor = this._accumulateFactor(
            previousScale,
            scaleFactor,
            '_wheelUnusedFactor'
          )
          if (scaleFactor < 1) {
            this.pdfDocumentViewer.zoomOut(null, scaleFactor)
          } else if (scaleFactor > 1) {
            this.pdfDocumentViewer.zoomIn(null, scaleFactor)
          } else {
            return
          }
        } else {
          const delta = pdfjsViewer.normalizeWheelEventDirection({
            deltaX: evt.deltaX,
            deltaY,
          })
          let ticks = 0
          if (
            deltaMode === WheelEvent.DOM_DELTA_LINE ||
            deltaMode === WheelEvent.DOM_DELTA_PAGE
          ) {
            if (Math.abs(delta) >= 1) {
              ticks = Math.sign(delta)
            } else {
              ticks = this._accumulateTicks(
                delta,
                '_wheelUnusedTicks'
              )
            }
          } else {
            const PIXELS_PER_LINE_SCALE = 30
            ticks = this._accumulateTicks(
              delta / PIXELS_PER_LINE_SCALE,
              '_wheelUnusedTicks'
            )
          }
          if (ticks < 0) {
            this.pdfDocumentViewer.zoomOut(-ticks)
          } else if (ticks > 0) {
            this.pdfDocumentViewer.zoomIn(ticks)
          } else {
            return
          }
        }
        this._centerAtPos(
          previousScale,
          evt.clientX,
          evt.clientY
        )
      } else {
        this.setZoomDisabledTimeout()
      }
    }

    /**
     * 这里使用perfectScroll之后有个问题，在缩放的时候无法禁止掉ps的滚动，所以就改了ps的源码，做了简单的处理
     * 不知道这里有没有更好的解法方案
     */
    document.addEventListener('wheel', handler, {
      passive: false,
    })
  }
}
