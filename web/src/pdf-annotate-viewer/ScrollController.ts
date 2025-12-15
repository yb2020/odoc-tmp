import * as pdfjsViewer from '@idea/pdfjs-dist/web/pdf_viewer';

export enum ScrollEvent {
  LAST_POSITION_UPDATED = 'x-last-position-updated',
  TRIGGER_SCROLL_CONTROLLER = 'x-trigger-scroll-controller'
}

interface LastPosition {
  left: number,
  top: number,
}

export default class ScrollController extends pdfjsViewer.PDFLinkService {

  private pdfViewerContainer?: HTMLDivElement

  private lastPdfViewerScrollPosition?: LastPosition

  private lockFlag = false

  private timer = 0

  setPdfViewerContainer(container: HTMLDivElement) {
    this.pdfViewerContainer = container
    this.pdfViewerContainer.addEventListener('scroll', () => {
      if (!this.lockFlag && this.lastPdfViewerScrollPosition) {
        this.lastPdfViewerScrollPosition = undefined
        this.eventBus.dispatch(ScrollEvent.LAST_POSITION_UPDATED, {
          lastPosition: this.lastPdfViewerScrollPosition,
          source: this,
        })
      }
    })
  }

  goToPage(val: string | number, noLastPosition?: boolean): void {
    this.lockFlag = true
    if (this.pdfViewerContainer && !noLastPosition && !this.lastPdfViewerScrollPosition) {
      this.lastPdfViewerScrollPosition = {
        left: this.pdfViewerContainer.scrollLeft,
        top: this.pdfViewerContainer.scrollTop,
      }
      this.eventBus.dispatch(ScrollEvent.LAST_POSITION_UPDATED, {
        lastPosition: this.lastPdfViewerScrollPosition,
        source: this,
      })
    }
    super.goToPage(val)
    if (this.timer) {
      clearTimeout(this.timer)
    }
    this.timer = window.setTimeout(() => {
      this.lockFlag = false
    }, 100);
    this.eventBus.dispatch(ScrollEvent.TRIGGER_SCROLL_CONTROLLER, {
      source: this,
      from: 'pageChange',
    })
  }

  scrollToPdfViewerLastPosition() {
    if (this.pdfViewerContainer && this.lastPdfViewerScrollPosition) {
      this.pdfViewerContainer.scrollLeft = this.lastPdfViewerScrollPosition.left
      this.pdfViewerContainer.scrollTop = this.lastPdfViewerScrollPosition.top
      this.lastPdfViewerScrollPosition = undefined
    }
    this.eventBus.dispatch(ScrollEvent.LAST_POSITION_UPDATED, {
      lastPosition: this.lastPdfViewerScrollPosition,
      source: this,
    })
  }

  addEventListener(event: ScrollEvent.TRIGGER_SCROLL_CONTROLLER, listener: (args: { from: 'pageChange' | 'manual', source: ScrollController }) => void): void;
  addEventListener(event: ScrollEvent.LAST_POSITION_UPDATED, listener: (args: { lastPosition: LastPosition, source: ScrollController }) => void): void;
  addEventListener(event: ScrollEvent, listener: unknown) {
    this.eventBus.on(event, listener as any)
  }

  removeEventListener(event: ScrollEvent, listener?: any) {
    this.eventBus.off(event, listener)
  }

  updatePdfViewerLastPosition(force?: boolean) {
    if (this.pdfViewerContainer && (force || !this.lastPdfViewerScrollPosition)) {
      this.lastPdfViewerScrollPosition = {
        left: this.pdfViewerContainer.scrollLeft,
        top: this.pdfViewerContainer.scrollTop,
      }

      this.eventBus.dispatch(ScrollEvent.LAST_POSITION_UPDATED, {
        lastPosition: this.lastPdfViewerScrollPosition,
        source: this,
      })
    }

    this.eventBus.dispatch(ScrollEvent.TRIGGER_SCROLL_CONTROLLER, {
      source: this,
      from: 'manual',
    })
    
  }

  getPdfViewerContainer() {
    if (!this.pdfViewerContainer) {
      throw new Error('call setPdfViewerContainer before getPdfViewerContainer')
    }
    return this.pdfViewerContainer
  }

  addLinkAttributes(link: HTMLAnchorElement, url: string): void {
    super.addLinkAttributes(link, url, true)
  }

  #goToDestinationHelper(rawDest: any, namedDest: any = null, explicitDest: any) {
    // Dest array looks like that: <page-ref> </XYZ|/FitXXX> <args..>
    const destRef = explicitDest[0];
    let pageNumber;

    if (typeof destRef === "object" && destRef !== null) {
      pageNumber = this._cachedPageNumber(destRef);

      if (!pageNumber) {
        // Fetch the page reference if it's not yet available. This could
        // only occur during loading, before all pages have been resolved.
        this.pdfDocument
          .getPageIndex(destRef)
          .then((pageIndex: number) => {
            this.cachePageRef(pageIndex + 1, destRef);
            this.#goToDestinationHelper(rawDest, namedDest, explicitDest);
            return
          })
          .catch(() => {
            console.error(
              `PDFLinkService.#goToDestinationHelper: "${destRef}" is not ` +
                `a valid page reference, for dest="${rawDest}".`
            );
          });
        return;
      }
    } else if (Number.isInteger(destRef)) {
      pageNumber = destRef + 1;
    } else {
      console.error(
        `PDFLinkService.#goToDestinationHelper: "${destRef}" is not ` +
          `a valid destination reference, for dest="${rawDest}".`
      );
      return;
    }
    if (!pageNumber || pageNumber < 1 || pageNumber > this.pagesCount) {
      console.error(
        `PDFLinkService.#goToDestinationHelper: "${pageNumber}" is not ` +
          `a valid page number, for dest="${rawDest}".`
      );
      return;
    }

    if (this.pdfHistory) {
      // Update the browser history before scrolling the new destination into
      // view, to be able to accurately capture the current document position.
      this.pdfHistory.pushCurrentPosition();
      this.pdfHistory.push({ namedDest, explicitDest, pageNumber });
    }

    this.updatePdfViewerLastPosition()

    this.lockFlag = true

    this.pdfViewer.scrollPageIntoView({
      pageNumber,
      destArray: explicitDest,
      ignoreDestinationZoom: this._ignoreDestinationZoom,
    });

    if (this.timer) {
      clearTimeout(this.timer)
    }
    this.timer = window.setTimeout(() => {
      this.lockFlag = false
    }, 100);
  }

  /**
   * This method will, when available, also update the browser history.
   *
   * @param {string|Array} dest - The named, or explicit, PDF destination.
   */
  async goToDestination(dest: string | any[]) {
    if (!this.pdfDocument) {
      return;
    }
    let namedDest, explicitDest;
    if (typeof dest === "string") {
      namedDest = dest;
      explicitDest = await this.pdfDocument.getDestination(dest);
    } else {
      namedDest = null;
      explicitDest = await dest;
    }
    if (!Array.isArray(explicitDest)) {
      console.error(
        `PDFLinkService.goToDestination: "${explicitDest}" is not ` +
          `a valid destination array, for dest="${dest}".`
      );
      return;
    }
    this.#goToDestinationHelper(dest, namedDest, explicitDest);
  }
}