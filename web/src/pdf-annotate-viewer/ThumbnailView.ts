import { PDFPageProxy } from '@idea/pdfjs-dist/types/src/display/api'
import cbem from './css/style'
// import * as WebUiUtils from '@idea/pdfjs-dist/lib/web/ui_utils'
import { OptionalContentConfig, PDFRenderingQueue } from '@idea/pdfjs-dist/types/web/pdf_page_view'
import { PageViewport } from '@idea/pdfjs-dist/types/src/display/display_utils'
import * as pdfjsViewer from '@idea/pdfjs-dist/web/pdf_viewer';
import * as pdfjsLib from '@idea/pdfjs-dist';

const { RenderingStates, OutputScale } = pdfjsViewer
interface ThumbnailViewOptions {
  container: HTMLDivElement
  id: number
  optionalContentConfigPromise: Promise<OptionalContentConfig>
  viewport: PageViewport
  renderingQueue: PDFRenderingQueue
  linkService: pdfjsViewer.PDFLinkService
  annotationMode: number
}

const ThumbnailBem = cbem('thumbnail')

const ThumbnailNumerBem = ThumbnailBem('num')

const ThumbnainImageBem = ThumbnailBem('image')

const ThumbnailCurrentBem = ThumbnailBem('', 'current')

const ThumbnailAnnotateBem = ThumbnailBem('annotate')

const classname = ThumbnailBem.toString();

const styleTag = document.createElement('style');
styleTag.textContent = `
  .${ThumbnailBem.toString()} {
    display: block;
    margin: auto;
  }

  .${ThumbnailBem.toString()} {
    position: relative;
    padding: 8px;
  }

  .${ThumbnailBem.toString()}.${ThumbnailCurrentBem.toString()} .${ThumbnainImageBem.toString()} {
    opacity: 1;
    box-shadow: 0 0 0 4px #7cacf3;
  }
  .${ThumbnailBem.toString()}.${ThumbnailCurrentBem.toString()} .${ThumbnailNumerBem.toString()} {
    color: #7cacf3;
  }
  .${ThumbnailBem.toString()}.${ThumbnailCurrentBem.toString()} .${ThumbnailAnnotateBem.toString()} {
    opacity: 1;
  }

  .${ThumbnainImageBem.toString()} {
    height: 100%;
    width: 100%;
    opacity: 0.4;

  }

  .${ThumbnailNumerBem.toString()} {
    position: absolute;
    color: #fff;
    bottom: 8px;
    left: 0;
    right: 0;
    text-align: center;
    font-size: 12px;
    line-height: 14px;
  }

  .${ThumbnailAnnotateBem.toString()} {
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    padding: 8px;
    opacity: 0.4;
    img {
      width: 100%;
    }
  }
`;

if (typeof document !== 'undefined') {
  document.head.appendChild(styleTag);
}

const DRAW_UPSCALE_FACTOR = 2; 
const THUMBNAIL_WIDTH = 114; // px

export default class ThumbnailView {
  private container
  id
  renderingId
  div: HTMLDivElement
  private optionalContentConfigPromise
  pdfPage: null | PDFPageProxy = null
  private renderTask: pdfjsLib.RenderTask | null = null
  private canvasWidth
  private canvasHeight
  private viewport
  private ring: HTMLDivElement
  private annotateImgDiv: HTMLDivElement
  private renderingState
  private renderingQueue
  resume: null | (() => void) = null

  private annotationMode = 1;

  constructor({ 
    container, id, optionalContentConfigPromise, viewport, renderingQueue, 
    linkService,
    annotationMode,
  }:ThumbnailViewOptions) {
    this.container = container
    this.id = id
    this.renderingId = "thumbnail" + id;
    this.optionalContentConfigPromise = optionalContentConfigPromise
    this.viewport = viewport
    this.canvasWidth = THUMBNAIL_WIDTH
    this.canvasHeight = THUMBNAIL_WIDTH / viewport.width * viewport.height
    this.renderingQueue = renderingQueue
    this.annotationMode = annotationMode
    
    const anchor = document.createElement("a");
    anchor.href = `#page=${id}`;
    anchor.title = `Page ${id}`
    anchor.onclick = function () {
      linkService.goToPage(id);
      return false;
    };
    
    anchor.className = classname

    const div = document.createElement("div");
    div.className = ThumbnailBem();
    div.setAttribute("data-page-number", `${id}`);
    div.style.width = `${this.canvasWidth + 8 * 2}px`
    div.style.height = `${this.canvasHeight + 8 * 2}px`
    this.div = div;

    const imageDiv = document.createElement("div");
    imageDiv.className = ThumbnainImageBem();
    this.ring = imageDiv

    const numDiv = document.createElement('div')
    numDiv.className = ThumbnailNumerBem()
    numDiv.innerText = `${id}`

    const annotateImageDiv = document.createElement('div')
    annotateImageDiv.classList.add(ThumbnailAnnotateBem())
    this.annotateImgDiv = annotateImageDiv

    div.appendChild(imageDiv);
    div.appendChild(annotateImageDiv);
    div.appendChild(numDiv)
    anchor.appendChild(div);
    this.container.appendChild(anchor);

    this.renderingState = RenderingStates.INITIAL
  }


  private getPageDrawContext(upscaleFactor = 1) {
    // Keep the no-thumbnail outline visible, i.e. `data-loaded === false`,
    // until rendering/image conversion is complete, to avoid display issues.
    const canvas = document.createElement("canvas");

    const ctx = canvas.getContext("2d", { alpha: false });
    const outputScale = new OutputScale();

    canvas.width = (upscaleFactor * THUMBNAIL_WIDTH * outputScale.sx) | 0;
    canvas.height = (upscaleFactor * this.canvasHeight * outputScale.sy) | 0;

    const transform = outputScale.scaled
      ? [outputScale.sx, 0, 0, outputScale.sy, 0, 0]
      : undefined;

    return { ctx, canvas, transform };
  }


  private convertCanvasToImage(canvas: HTMLCanvasElement) {
    
    
    const imgUrl = canvas.toDataURL('image/png')

    const image = document.createElement("img");

    image.style.width = this.canvasWidth + "px";
    image.style.height = this.canvasHeight + "px";

    image.src = imgUrl

    this.div.setAttribute("data-loaded", 'true');
    this.ring.appendChild(image);
  }

  async draw() {
    const { pdfPage } = this
    if (!pdfPage) {
      throw new Error('pdfPage is not loaded') 
    }

    this.renderingState = RenderingStates.RUNNING

    const finishRenderTask = async (error?: Error) => {
      // The renderTask may have been replaced by a new one, so only remove
      // the reference to the renderTask if it matches the one that is
      // triggering this callback.
      if (renderTask === this.renderTask) {
        this.renderTask = null;
      }

      if (error instanceof pdfjsLib.RenderingCancelledException) {
        return;
      }
      this.renderingState = RenderingStates.FINISHED;
      this.convertCanvasToImage(canvas);

      if (error) {
        throw error;
      }
    };

    // TODO 未来需要处理rotation的情况
    const drawViewport = this.viewport.clone({
      scale: DRAW_UPSCALE_FACTOR * this.canvasHeight / this.viewport.height,
    });


    const { ctx, canvas, transform } =
      this.getPageDrawContext(DRAW_UPSCALE_FACTOR);

    
      const renderContinueCallback = (cont: () => void) => {
        
        if (!this.renderingQueue.isHighestPriority(this as any)) {
          this.renderingState = RenderingStates.PAUSED;
          this.resume = () => {

            this.renderingState = RenderingStates.RUNNING;
            cont();
          };
          return;
        }
        cont();
      };

    const renderContext = {
      canvasContext: ctx!,
      transform,
      viewport: drawViewport,
      optionalContentConfigPromise: this.optionalContentConfigPromise,
      annotationMode: this.annotationMode,
    };

    const renderTask = pdfPage.render(renderContext)

    renderTask.onContinue = renderContinueCallback;

    this.renderTask = renderTask

    const resultPromise = renderTask.promise.then(
       () => {
        return finishRenderTask();
      },
       (error) => {
        return finishRenderTask(error);
      }
    );

    // eslint-disable-next-line promise/catch-or-return
    resultPromise.finally(() => {
      // Zeroing the width and height causes Firefox to release graphics
      // resources immediately, which can greatly reduce memory consumption.
      canvas.width = 0;
      canvas.height = 0;

      // Only trigger cleanup, once rendering has finished, when the current
      // pageView is *not* cached on the `BaseViewer`-instance.
      // const pageCached = this.linkService.isPageCached(this.id);
      // if (!pageCached) {
      //   this.pdfPage?.cleanup();
      // }
    });

    return resultPromise

  }

  cancel() {
    if (this.renderTask) {
      this.renderTask.cancel()
      this.renderTask = null
      this.renderingState = RenderingStates.INITIAL
    }
    this.resume = null
  }

  setPdfPage(pdfPage: PDFPageProxy) {
    this.pdfPage = pdfPage
    this.viewport = pdfPage.getViewport({ scale: 1 });
    this.reset()
  }

  reset() {
    this.cancel()
    this.renderingState = RenderingStates.INITIAL;
    /**
       * 这里未来可能会有rotation，那么viewport会变化，高度也就遍历
       */
     const { viewport, div } = this
     this.canvasHeight = THUMBNAIL_WIDTH / viewport.width * viewport.height

     div.style.width = `${this.canvasWidth + 8 * 2}px`
     div.style.height = `${this.canvasHeight + 8 * 2}px`
   
     this.ring.removeAttribute('data-loaded')
     this.ring.innerHTML = ''

  }

  setSelected() {
    this.div.classList.add(ThumbnailCurrentBem())
  }

  setUnSelected() {
    this.div.classList.remove(ThumbnailCurrentBem())
  }

  setAnnotateImage(imgUrl: string) {
    const image = this.annotateImgDiv.querySelector('img')
    if (image) {
      image.src = imgUrl
    } else {
      const image = document.createElement('img')
      image.src = imgUrl
      this.annotateImgDiv.appendChild(image)
    }
    
  }

  getRenderingState() {
    return this.renderingState
  }
}