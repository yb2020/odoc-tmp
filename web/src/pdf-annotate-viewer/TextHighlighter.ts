import { PageViewport } from '@idea/pdfjs-dist/types/web/interfaces'
import { TextHighlighterOptions } from '@idea/pdfjs-dist/types/web/text_highlighter'
import anime, { AnimeInstance } from 'animejs'
import { InnerFinderEvent } from './FinderController'
import { TextContentBound, TextRectCoordinate } from './type'
import { createCanvas, drawTextRectCoordinate } from './utils/canvas'
import { getHightlightRectCoordinate } from './utils/hightlight'

const clearDraw = ({
  ctx,
  coordinate,
}: {
  ctx: CanvasRenderingContext2D
  coordinate: TextRectCoordinate
}) => {
  const { x, y, width, height, offset } = coordinate
  if (offset) {
    ctx.save()
    ctx.translate(x, y)
    ctx.rotate(offset.angle)
    ctx.translate(offset.dx, offset.dy)
    ctx.clearRect(0, 0, width, height)
    ctx.restore()
  } else {
    ctx.clearRect(x || 0, y || 0, width || 0, height || 0)
  }
}

const HIGHTLIGHT_NORMAL_COLOR = 'rgba(180, 0, 170, 0.8)'
const HIGHTLIGHT_MATCH_COLOR = 'rgba(0, 100, 0, 0.85)'
const HIGHTLIGHT_MATCH_START_COLOR = 'rgba(255, 255, 255, 0.2)'
const MATCH_SCROLL_OFFSET_TOP = -50 // px

class Defer<T> {
  promise: Promise<T>;
  resolve!: (value: T | PromiseLike<T>) => void;
  reject!: (reason?: any) => void;

  constructor() {
      this.promise = new Promise<T>((resolve, reject) => {
          this.resolve = resolve;
          this.reject = reject;
      });
  }
}

export default class TextHighlighter {
  findController
  eventBus
  pageIdx
  enabled

  viewport: null | PageViewport = null
  #highlightLayerCanvas: HTMLCanvasElement | null = null
  textContentBounds: TextContentBound[] | null = null
  matches: TextRectCoordinate[] = []
  matchesAnimation: null | AnimeInstance = null
  curMatchIdx = -1

  _onUpdateTextLayerMatches: ((evt: { pageIndex: number }) => void) | null

  readyPromise: Defer<void> = new Defer()

  constructor(options: TextHighlighterOptions) {
    this.findController = options.findController
    this.eventBus = options.eventBus
    this.pageIdx = options.pageIndex
    this.enabled = false

    this._onUpdateTextLayerMatches = null
  }

  get highlightLayerCanvas() {
    return this.#highlightLayerCanvas
  }

  setTextMapping({
    viewport,
    textContentBounds,
  }: {
    viewport: PageViewport
    textContentBounds?: TextContentBound[]
  }) {
    if (!this.#highlightLayerCanvas) {
      this.#highlightLayerCanvas = createCanvas(viewport)
    } else {
      createCanvas(viewport, this.#highlightLayerCanvas)
    }
    if (textContentBounds) {
      this.textContentBounds = textContentBounds
    }
    this.viewport = viewport
    this.readyPromise.resolve();
  }

  enable() {
    if (!this.textContentBounds) {
      throw new Error('Text divs and strings have not been set.')
    }
    if (this.enabled) {
      throw new Error('TextHighlighter is already enabled.')
    }
    this.enabled = true

    if (!this._onUpdateTextLayerMatches) {
      this._onUpdateTextLayerMatches = (evt: { pageIndex: number }) => {
        if (evt.pageIndex === this.pageIdx || evt.pageIndex === -1) {
          this._updateMatches()
        }
      }
      this.eventBus.on('updatetextlayermatches', this._onUpdateTextLayerMatches)

      const selected = this.findController?.selected

      if (selected?.pageIdx === this.pageIdx) {
        setTimeout(() => {
          this.eventBus.dispatch(InnerFinderEvent.GOTO_MATCH, { selected })
        }, 0)
      }
    }

    this._updateMatches()
  }

  disable() {
    this.enabled = false
    this._updateMatches(true)
    if (this._onUpdateTextLayerMatches) {
      this.eventBus.off(
        'updatetextlayermatches',
        this._onUpdateTextLayerMatches
      )
    }
  }

  clearMatches(force?: boolean) {
    if (this.matches?.length || force) {
      this.matches = []
      this.curMatchIdx = -1
      this.#highlightLayerCanvas!.getContext('2d')?.clearRect(
        0,
        0,
        this.viewport!.width,
        this.viewport!.height
      )
    }
  }

  _renderMatches(rects: TextRectCoordinate[], color?: string, clear?: boolean) {
    const destColor =
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      color ?? this.findController.highlightColor ?? HIGHTLIGHT_NORMAL_COLOR

    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    if (this.findController.highlightTwinklingEnabled) {
      const styles = {
        background: HIGHTLIGHT_MATCH_START_COLOR,
      }
      this.matchesAnimation?.pause()
      this.matchesAnimation = anime({
        targets: styles,
        loop: 2,
        background: destColor,
        easing: 'easeInOutQuad',
        update: () => {
          this._drawMatches(rects, styles.background, true)
        },
      })
      return
    }

    this._drawMatches(rects, destColor, clear)
  }

  _drawMatches(rects: TextRectCoordinate[], color: string, clear?: boolean) {
    const highlightCtx = this.#highlightLayerCanvas!.getContext('2d')

    if (!highlightCtx) {
      return
    }
    rects.forEach((rect) => {
      if (clear) {
        clearDraw({
          ctx: highlightCtx,
          coordinate: rect,
        })
      }
      if (!rect.shouldScaleText) {
        return
      }
      drawTextRectCoordinate({
        ctx: highlightCtx,
        coordinate: rect,
        color,
      })
    })
  }

  cancelMatch() {
    if (!this.textContentBounds?.length) {
      console.warn(
        'unavaliable cancel match, textlayer is not ready. PageNum: ' +
          this.pageIdx
      )
      return null
    }
    this.drawMatch(this.curMatchIdx, HIGHTLIGHT_NORMAL_COLOR)
    this.curMatchIdx = -1
  }

  drawMatch(idx: number, color?: string) {
    if (idx < 0) {
      return []
    }
    const { findController } = this

    if (findController) {
      const pageIdx = this.pageIdx
      const matches: number[] = findController.pageMatches?.[pageIdx]
      const matchesLength = findController.pageMatchesLength?.[pageIdx]
      const result = getHightlightRectCoordinate({
        viewport: this.viewport!,
        textContentItems: this.textContentBounds!,
        matches: [matches[idx]],
        matchesLength: [matchesLength[idx]],
      })
      this._renderMatches(result, color, true)
      return result
    }
    return []
  }

  goToMatch(idx: number) {
    if (!this.textContentBounds?.length) {
      console.warn(
        'unavaliable goto match, textlayer is not ready. PageIdx: ' +
          this.pageIdx
      )
      return null
    }
    this.cancelMatch()
    this.curMatchIdx = idx
    const result = this.drawMatch(idx, HIGHTLIGHT_MATCH_COLOR)
    const divTop = (
      (this.#highlightLayerCanvas?.offsetParent as HTMLDivElement)
        ?.offsetParent as HTMLDivElement
    )?.offsetTop
    // TODO angle !== 0 的跳转
    if (result[0]?.offset) {
      return divTop
    }
    return divTop + result[0]?.y + MATCH_SCROLL_OFFSET_TOP || divTop
  }

  _convertMatches(matches: number[], matchesLength: number[]) {
    if (!matches) {
      return []
    }
    const result = getHightlightRectCoordinate({
      viewport: this.viewport!,
      textContentItems: this.textContentBounds!,
      matches,
      matchesLength,
    })
    return result
  }

  _updateMatches(reset = false) {
    if ((!this.enabled && !reset) || !this.#highlightLayerCanvas) {
      return
    }
    // Clear all current matches.
    this.clearMatches()

    if (reset) {
      return
    }

    const { findController, pageIdx } = this

    if (!findController || !findController?.highlightMatches) {
      return
    }

    // Convert the matches on the `findController` into the match format
    // used for the textLayer.
    const pageMatches = findController.pageMatches?.[pageIdx] || null
    const pageMatchesLength =
      findController.pageMatchesLength?.[pageIdx] || null
    this.matches = this._convertMatches(pageMatches, pageMatchesLength)
    this._renderMatches(this.matches)
  }
}
