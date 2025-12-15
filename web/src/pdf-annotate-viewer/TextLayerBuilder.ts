import { TextContent } from '@idea/pdfjs-dist/types/src/display/api'
import { PageViewport } from '@idea/pdfjs-dist/types/web/interfaces'
import TextHighlighter from './TextHighlighter'
import { TextLayerBuilderOptions } from '@idea/pdfjs-dist/types/web/text_layer_builder'
import TextLayerRenderTask from './TextLayerRenderTask'
import { BaseMousePoint, TextContentBound, TextRectCoordinate } from './type'
import {
  getNearestDistanceBetweenPointAndRect,
  getTextRectCoordinate,
  getWordRectCoordinate,
  isInRotationBounds,
  isInRotationRect,
  optimizeRects,
} from './utils/bound'
import { normalize } from './utils/text'
import { createCanvas, drawTextRectCoordinate } from './utils/canvas'

const getRadioPoint = (point: BaseMousePoint, radio: number) => {
  return {
    x: point.x / radio,
    y: point.y / radio,
  }
}

const isBoundPoint = (p: BaseMousePoint, viewport: PageViewport) => {
  const boundPoint = {
    x: viewport.width,
    y: viewport.height,
  }
  if (p.x === boundPoint.x && p.y === boundPoint.y) {
    return true
  }
  return false
}



// const MATCH_SCROLL_OFFSET_LEFT = -400; // px

const SELECT_DEFAULT_COLOR = 'rgba(58, 87, 230, 85%)'

type Nullable<T> = {
  [P in keyof T]: T[P] | null
}

interface Options extends Omit<TextLayerBuilderOptions, 'highlighter'> {
  highlighter: TextHighlighter
}

export default class TextLayerBuilder {
  #rotation = 0

  #scale = 0

  #textContentSource: ReadableStream | TextContent | null = null

  #textLayerCanvas: HTMLCanvasElement | null = null

  div: HTMLDivElement
  renderingDone: boolean
  highlighter
  accessibilityManager
  isOffscreenCanvasSupported

  textLayerRenderTask: TextLayerRenderTask | null

  /* 选区START */
  textContentBounds: TextContentBound[] = []
  rects: TextRectCoordinate[] = []
  stagingRectsList: TextRectCoordinate[][] = []
  crossing = false
  radio = 1.0
  viewport: null | PageViewport = null
  tempTextCanvasCtx = document
    .createElement('canvas')
    .getContext('2d', { alpha: false })
  /* 选区END */

  constructor({
    highlighter = null,
    accessibilityManager = null,
    isOffscreenCanvasSupported = true,
  }: Nullable<Options>) {
    this.renderingDone = false
    this.textLayerRenderTask = null
    this.highlighter = highlighter
    this.accessibilityManager = accessibilityManager
    this.isOffscreenCanvasSupported = isOffscreenCanvasSupported

    this.div = document.createElement('div')
    this.div.className = 'textLayer';
    this.hide()
  }

  #finishRendering() {
    this.renderingDone = true

    // const endOfContent = document.createElement("div");
    // endOfContent.className = "endOfContent";
    // this.div.append(endOfContent);
  }

  #updateCanvasLayer() {
    if (this.#textLayerCanvas) {
      createCanvas(this.viewport!, this.#textLayerCanvas)
    }
  }

  #createCanvasLayer() {
    this.#textLayerCanvas = createCanvas(this.viewport!)
    this.div.appendChild(this.#textLayerCanvas)
  }

  get numTextDivs() {
    return 0
  }

  /**
   * Renders the text layer.
   * @param {PageViewport} viewport
   */
  async render(viewport: PageViewport) {
    if (!this.#textContentSource) {
      throw new Error('No "textContentSource" parameter specified.')
    }

    this.viewport = viewport

    const scale = viewport.scale * (globalThis.devicePixelRatio || 1)
    const { rotation } = viewport
    if (this.renderingDone) {
      const mustRotate = rotation !== this.#rotation
      const mustRescale = scale !== this.#scale
      if (mustRotate || mustRescale) {
        this.hide();
        
        this.#scale = scale
        this.#rotation = rotation
      }
      this.#updateCanvasLayer();
      this.textContentBounds = this.textLayerRenderTask!.updateTextContentBounds(viewport);
      this.highlighter?.setTextMapping({
        viewport,
        textContentBounds: this.textContentBounds,
      });
      this.show();
      return;
    }

    this.cancel()

    this.#createCanvasLayer();

    // this.accessibilityManager?.setTextMapping(this.textDivs);

    this.textLayerRenderTask = new TextLayerRenderTask({
      viewport,
      textContentSource: this.#textContentSource,
    })

    

    this.textContentBounds = await this.textLayerRenderTask.render(viewport)

    this.highlighter?.setTextMapping({
      viewport,
      textContentBounds: this.textContentBounds,
    });
    this.div.appendChild(this.highlighter!.highlightLayerCanvas!)

    this.#finishRendering()
    this.#scale = scale
    this.#rotation = rotation
    this.show()
    // this.accessibilityManager?.enable();
  }

  hide() {
    if (!this.div.hidden) {
      // We turn off the highlighter in order to avoid to scroll into view an
      // element of the text layer which could be hidden.
      this.highlighter?.disable();
      this.div.hidden = true
    }
  }

  show() {
    if (this.div.hidden && this.renderingDone) {
      this.div.hidden = false
      this.highlighter?.enable();
    }
  }

  /**
   * Cancel rendering of the text layer.
   */
  cancel() {
    if (this.textLayerRenderTask) {
      this.textLayerRenderTask.cancel()
      this.textLayerRenderTask = null
    }
    this.highlighter?.disable();
    this.textContentBounds.length = 0;
    this.rects.length = 0;
    // this.accessibilityManager?.disable();
    // this.textContentItemsStr.length = 0;
    // this.textDivs.length = 0;
    // this.textDivProperties = new WeakMap();
  }

  /**
   * @param {ReadableStream | TextContent} source
   */
  setTextContentSource(source: ReadableStream | TextContent) {
    this.cancel()
    this.#textContentSource = source
  }

  isInTextBound(point: BaseMousePoint, boundBuff: number) {
    if (!this.textContentBounds.length || point.x < 0 || point.y < 0) {
      return false
    }
    return isInRotationBounds(
      this.textContentBounds,
      getRadioPoint(point, this.radio),
      boundBuff
    )
  }

  stage() {
    this.stagingRectsList.push(this.rects);
    this.rects = [];
  }

  private get hasStagingRects() {
    return this.stagingRectsList.some(rects => rects.length > 0);
  }

  select(sP: BaseMousePoint, eP: BaseMousePoint, boundBuff: number) {
    if (!this.viewport || !this.#textLayerCanvas) {
      return
    }
    sP = getRadioPoint(sP, this.radio)
    eP = getRadioPoint(eP, this.radio)
    const isValidPoint = (p: BaseMousePoint) => {
      return p.x >= 0 && p.y >= 0;
    }
    const ctx = this.#textLayerCanvas.getContext('2d')
    if (
      !this.textContentBounds.length ||
      !ctx ||
      !isValidPoint(sP) ||
      !isValidPoint(eP) ||
      (Math.abs(sP.x - eP.x) <= 1 && Math.abs(sP.y - eP.y) <= 1)
    ) {
      return
    }

    const [p1, p2] = [sP, eP]

    let startIndex = -1

    let startDistance = Infinity

    const endDelta = {
      insideIndex: -1,
      insideBoundIndex: -1,
      possibleIndex: -1,
    }

    const buff = boundBuff

    if (p1.x === 0 && p1.y === 0) {
      startIndex = 0
    }

    const len = this.textContentBounds.length

    for (let i = 0; i < len; i += 1) {
      const index = i
      const bound = this.textContentBounds[index]
      if (!bound.shouldScaleText) {
        continue
      }
      // 1. 起始鼠标点必须落在加了boundBuff的bound矩形内，
      // 2. 因为加了boundBuff导致相邻的矩形会有重叠，那么鼠标点可能落在了多个bound里面
      // 3. 计算出鼠标点距离哪个去掉boundBuff的bound矩形最近，那么这个矩形就是起始bound
      // 如果点在矩形内部，那么就是这个矩形，如果点在矩形外部，那么最短的距离的点肯定是矩形的顶点
      if (isInRotationRect(bound.offset, p1, buff)) {
        const distance = getNearestDistanceBetweenPointAndRect(bound.offset, p1)

        if (distance < startDistance) {
          startIndex = index
          startDistance = distance
        }
      }

      if (
        Math.abs(p2.x - this.viewport.width) <= 1 &&
        Math.abs(p2.y - this.viewport.height) <= 1
      ) {
        endDelta.possibleIndex = this.textContentBounds.length - 1
      } else if (isInRotationRect(bound.offset, p2, buff)) {
        if (isInRotationRect(bound.offset, p2, 0)) {
          endDelta.insideIndex = index
        } else if (endDelta.insideIndex !== 0) {
          endDelta.insideBoundIndex = index
        }
      }
    }

    const endIndex =
      endDelta.insideIndex >= 0
        ? endDelta.insideIndex
        : endDelta.insideBoundIndex >= 0
        ? endDelta.insideBoundIndex
        : endDelta.possibleIndex

    if (startIndex < 0 || endIndex < 0) {
      return;
    }

    const [index1, index2, point1, point2] = startIndex < endIndex
      ? [startIndex, endIndex, p1, p2]
      : [endIndex, startIndex, p2, p1];

    this._drawRects(
      this.textContentBounds.slice(index1, index2 + 1),
      point1,
      point2
    );
  }

  clearTextLayer({
    stage = false,
    force = false,
  }) {
    if (!this.hasStagingRects && !this.rects.length && !force) {
      return;
    }

    const { width, height } = this.viewport as PageViewport;
    this.#textLayerCanvas
      ?.getContext('2d')
      ?.clearRect(0, 0, width, height);

    if (stage) {
      this.stagingRectsList = [];
    }

    this.rects = [];

    this.crossing = false;
  }

  private _drawRects(
    bounds: TextContentBound[],
    p1: BaseMousePoint,
    p2: BaseMousePoint
  ) {
    if (!this.viewport) {
      return
    }
    if (bounds.length <= 0) {
      return
    }

    // bounds = bounds.filter(b => b.shouldScaleText)

    this.clearTextLayer({ force: true })

    const tempCtx = this.tempTextCanvasCtx!

    const rects: TextRectCoordinate[] = []

    if (bounds.length === 1) {
      rects.push(
        getTextRectCoordinate({
          bound: bounds[0],
          points: [p1, p2],
          ctx: tempCtx,
          viewport: this.viewport,
        })
      )
    } else {
      const startBound = bounds.shift()
      const endBound = bounds.pop()

      rects.push(
        getTextRectCoordinate({
          bound: startBound!,
          points: [p1, null],
          ctx: tempCtx,
          viewport: this.viewport,
        })
      )

      bounds.forEach((bound) => {
        rects.push(
          getTextRectCoordinate({
            bound: bound,
            points: [null, null],
            ctx: tempCtx,
            viewport: this.viewport!,
          })
        )
      })

      rects.push(
        getTextRectCoordinate({
          bound: endBound!,
          points: [null, isBoundPoint(p2, this.viewport) ? null : p2],
          ctx: tempCtx,
          viewport: this.viewport,
        })
      )
    }

    optimizeRects(rects);

    const SPACE = ' ';
    // 空格合并到前一个字符串，如果没有前一个字符串，空格合并到后一个字符串
    for (let i = 0; i < rects.length;) {

      if (rects[i].text !== SPACE) {
        i += 1;
        continue;
      }

      const prev = rects[i - 1];
      const next = rects[i + 1];

      if (prev) {
        prev.text += SPACE;
      } else if (next && next.text !== SPACE) {
        next.text = SPACE + next.text;
      }

      rects.splice(i, 1);
    }

    this.rects = rects;

    const textCtx = this.#textLayerCanvas?.getContext('2d')

    if (!textCtx) {
      return
    }

    const redraw = (rect: TextRectCoordinate) => {
      if (!rect.shouldScaleText) {
        return
      }
      drawTextRectCoordinate({
        ctx: textCtx,
        coordinate: rect,
        color: SELECT_DEFAULT_COLOR,
      })
    }

    this.stagingRectsList.forEach(stageRects => stageRects.forEach(redraw));
    this.rects.forEach(redraw)

  }

  public getTextRects() {
    const allRects: TextRectCoordinate[] = [];
    const allTexts: string[] = [];

    const addToAll = (rect: TextRectCoordinate) => {
      allRects.push(rect);
      allTexts.push(rect.text);
    }

    this.stagingRectsList.forEach(rects => {
      rects.forEach(addToAll);

      // 多段拼接时，每段最后一个字符后面加一个空格
      const last = allTexts[allTexts.length - 1];
      if (last && !/\s$/.test(last)) {
        allTexts.push(' ');
      }
    });

    this.rects.forEach(addToAll);

    if (!allRects.length) {
      return null;
    }

    const [text] = normalize(allTexts.join('')) as string[];

    const multiSegment = this.hasStagingRects && this.rects.length > 0;
    
    return {
      multiSegment,
      crossing: this.crossing,
      rects: allRects,
      text: text.trim(),
      viewport: this.viewport as PageViewport,
    }
  }

  updateRectsColor(color?: string) {
    const textCtx = this.#textLayerCanvas?.getContext('2d')
    if (!textCtx) {
      return
    }
    textCtx.clearRect(0, 0, this.viewport!.width, this.viewport!.height)
    this.rects?.forEach((rect) => {
      drawTextRectCoordinate({
        ctx: textCtx,
        coordinate: rect,
        color: color || SELECT_DEFAULT_COLOR,
      })
    })
  }

  selectWord(point: BaseMousePoint, boundBuff: number): boolean {
    if (
      !this.textContentBounds.length ||
      point.x < 0 ||
      point.y < 0 ||
      !this.#textLayerCanvas?.getContext('2d')
    ) {
      return false
    }
    let curIndex = -1
    let curDistance = Infinity

    for (let i = 0; i < this.textContentBounds.length; i += 1) {
      const bound = this.textContentBounds[i]
      if (!bound.shouldScaleText) {
        continue
      }
      if (isInRotationRect(bound.offset, point, boundBuff)) {
        const distance = getNearestDistanceBetweenPointAndRect(
          bound.offset,
          point
        )
        if (distance < curDistance) {
          curIndex = i
          curDistance = distance
        }
      }
    }

    const curBound = this.textContentBounds[curIndex]

    if (!curBound) {
      return false
    }

    const textCtx = this.#textLayerCanvas.getContext('2d')

    const result = getWordRectCoordinate({
      bound: curBound,
      point,
      viewport: this.viewport!,
    })

    if (result) {
      this.rects.push(result)
    }

    this.rects.forEach((rect) => {
      drawTextRectCoordinate({
        ctx: textCtx!,
        coordinate: rect,
        color: SELECT_DEFAULT_COLOR,
      })
    })

    return true
  }

}
