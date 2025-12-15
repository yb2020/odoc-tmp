import * as pdfjsViewer from '@idea/pdfjs-dist/web/pdf_viewer';
import * as pdfjsLib from '@idea/pdfjs-dist';
import ThumbnailView from './ThumbnailView';
import cbem from './css/style';
import PerfectScrollbar from 'perfect-scrollbar';
import scrollIntoView from 'scroll-into-view-if-needed';
import { PDFRenderingQueue } from '@idea/pdfjs-dist/types/web/pdf_page_view';
import { createContainer } from './utils/ui';

const { getVisibleElements, watchScroll } = pdfjsViewer;

const ThumbnailsBem = cbem('thumbnails');

const classname = ThumbnailsBem.toString();

const styleTag = document.createElement('style');
styleTag.textContent = `
  .${ThumbnailsBem.toString()} {
    background-color: transparent;
  }
  
  .${ThumbnailsBem.toString()},
  .${ThumbnailsBem.toString()} * {
    box-sizing: border-box;
  }
  
  .${ThumbnailsBem.toString()} a {
    width: 130px;
  }
  
  .${ThumbnailsBem.toString()} {
    position: absolute;
    overflow: auto;
    left: 0;
    right: 0;
    bottom: 0;
    top: 0;
    padding: 10px 0;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
`;

if (typeof document !== 'undefined') {
  document.head.appendChild(styleTag);
}

interface ThumbnailViewerOptions {
  pdfDocument: pdfjsLib.PDFDocumentProxy;
  wrapper: HTMLDivElement;
  eventBus: pdfjsViewer.EventBus;
  currentPageNumber?: number;
  renderingQueue: PDFRenderingQueue;
  linkService: pdfjsViewer.PDFLinkService;
  annotationMode: number;
}

interface VisibleElementItem<T> {
  id: number;
  x: number;
  y: number;
  view: T;
  percent: number;
  widthPercent: number;
}

interface VisibleThumbs {
  first?: VisibleElementItem<ThumbnailView>;
  last: VisibleElementItem<ThumbnailView>;
  views: VisibleElementItem<ThumbnailView>[];
  ids: Set<number>;
}

export default class ThumbnailViewer {
  private pdfDocument: pdfjsLib.PDFDocumentProxy;
  private container;
  private thumbnails: ThumbnailView[] = [];
  private containerScrollbar;
  private currentPageNumber = 0;
  private renderingQueue: PDFRenderingQueue;
  private scroll;
  private eventBus;
  private annotationMode = 1;

  constructor(options: ThumbnailViewerOptions) {
    const {
      pdfDocument,
      wrapper,
      currentPageNumber,
      renderingQueue,
      linkService,
      eventBus,
      annotationMode,
    } = options;
    this.pdfDocument = pdfDocument;
    this.eventBus = eventBus;
    this.renderingQueue = renderingQueue;
    this.annotationMode = annotationMode;

    const container = createContainer(wrapper);
    container.classList.add(classname);
    container.classList.add(ThumbnailsBem());

    this.container = container;

    this.scroll = watchScroll(this.container, this.scrollUpdated.bind(this));

    this.start({
      currentPageNumber,
      linkService,
    });

    this.containerScrollbar = this.setPerfectScrollbar();
  }

  private scrollUpdated() {
    this.renderingQueue.renderHighestPriority(undefined as any);
  }

  private getScrollAhead(visible: VisibleThumbs) {
    if (visible.first?.id === 1) {
      return true;
    } else if (visible.last?.id === this.thumbnails.length) {
      return false;
    }
    return this.scroll.down;
  }

  setPerfectScrollbar() {
    const ps = new PerfectScrollbar(this.container, {
      wheelSpeed: 1,
      wheelPropagation: true,
      minScrollbarLength: 5,
    });
    window.addEventListener('resize', () => {
      ps.update();
      this.forceRendering();
    });

    return ps;
  }

  async start({
    currentPageNumber,
    linkService,
  }: {
    currentPageNumber?: number;
    linkService: pdfjsViewer.PDFLinkService;
  }) {
    const pdfDocument = this.pdfDocument;
    const pagesCount = pdfDocument.numPages;
    const firstPagePromise = pdfDocument.getPage(1);
    const firstPdfPage = await firstPagePromise;
    const viewport = firstPdfPage.getViewport({ scale: 1 });
    for (let pageNum = 1; pageNum <= pagesCount; ++pageNum) {
      const thumbnail = new ThumbnailView({
        container: this.container,
        id: pageNum,
        optionalContentConfigPromise: pdfDocument.getOptionalContentConfig(),
        viewport,
        renderingQueue: this.renderingQueue,
        linkService,
        annotationMode: this.annotationMode,
      });
      this.thumbnails.push(thumbnail);
    }
    this.thumbnails[0].setPdfPage(firstPdfPage);

    this.scrollToThumbnail(currentPageNumber || 1, {
      scrollMode: 'if-needed',
      block: 'start',
      inline: 'start',
    });

    this.forceRendering();

    this.containerScrollbar.update();

    this.eventBus.on(
      'pagechanging',
      ({ pageNumber }: { pageNumber: number }) => {
        this.scrollThumbnailIntoView(pageNumber);
      }
    );
  }

  async ensurePdfPageLoaded(thumbView: ThumbnailView) {
    if (thumbView.pdfPage) {
      return thumbView.pdfPage;
    }
    try {
      const pdfPage = await this.pdfDocument.getPage(thumbView.id);
      if (!thumbView.pdfPage) {
        thumbView.setPdfPage(pdfPage);
      }
      return pdfPage;
    } catch (reason) {
      console.error('Unable to get page for thumb view', reason);
      return null; // Page error -- there is nothing that can be done.
    }
  }

  forceRendering() {
    const visibleThumbs = this.getVisibleThumbs();
    const scrollAhead = this.getScrollAhead(visibleThumbs);
    const thumbView = this.renderingQueue.getHighestPriority(
      visibleThumbs,
      this.thumbnails,
      scrollAhead
    );

    if (thumbView) {
      this.ensurePdfPageLoaded(thumbView).then(() => {
        this.renderingQueue.renderView(thumbView);
        return;
      });
      return true;
    }
    return false;
  }

  private scrollThumbnailIntoView(pageNumber: number, scrollOptions?: any) {
    const thumbnailView = this.thumbnails[pageNumber - 1];

    if (!thumbnailView) {
      console.error('scrollThumbnailIntoView: Invalid "pageNumber" parameter.');
      return;
    }

    if (pageNumber !== this.currentPageNumber) {
      const prevThumbnailView = this.thumbnails[this.currentPageNumber - 1];
      prevThumbnailView.setUnSelected();
      thumbnailView.setSelected();
    }
    const { first, last, views } = this.getVisibleThumbs();

    if (views.length > 0 && first && last) {
      let shouldScroll = false;
      if (pageNumber <= first.id || pageNumber >= last.id) {
        shouldScroll = true;
      } else {
        for (const { id, percent } of views) {
          if (id !== pageNumber) {
            continue;
          }
          shouldScroll = percent < 100;
          break;
        }
      }
      if (shouldScroll) {
        const defaultScrollOptions =
          last.id === pageNumber
            ? {
                scrollMode: 'if-needed',
                block: 'start',
                inline: 'start',
              }
            : {
                scrollMode: 'if-needed',
                block: 'nearest',
                inline: 'nearest',
              };
        scrollIntoView(
          thumbnailView.div,
          scrollOptions || defaultScrollOptions
        );
      }
    }
    this.currentPageNumber = pageNumber;
  }

  getVisibleThumbs() {
    return getVisibleElements({
      scrollEl: this.container,
      views: this.thumbnails,
      sortByVisibility : false,
      horizontal : false,
      rtl : false,
    }) as VisibleThumbs;
  }

  scrollToThumbnail(num: number, scrollOptions?: any) {
    if (this.currentPageNumber !== num) {
      this.currentPageNumber = num;
      this.thumbnails.forEach((view) => {
        if (view.id === num) {
          view.setSelected();
        } else {
          view.setUnSelected();
        }
      });
      this.scrollThumbnailIntoView(num, scrollOptions);
    }
  }

  setAnnotateImage({
    pageIndex,
    imgUrl,
  }: {
    pageIndex: number;
    imgUrl: string;
  }) {
    const cur = this.thumbnails[pageIndex];
    if (!cur) {
      throw Error('invalid pageIndex');
    }
    cur.setAnnotateImage(imgUrl);
  }

  clear() {
    this.thumbnails.filter(e => e.getRenderingState())
                                  .forEach(e => e.reset());
  }
}
