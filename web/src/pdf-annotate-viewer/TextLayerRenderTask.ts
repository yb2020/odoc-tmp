import { TextContent, TextItem, TextMarkedContent } from "@idea/pdfjs-dist/types/src/display/api";
import { PageViewport } from "@idea/pdfjs-dist/types/src/display/display_utils";
import * as SharedUtil from "@idea/pdfjs-dist/lib/shared/util";
import { TextContentBound, TextContentChar } from "./type";
import { getTextProperties } from "./utils/text";
import { TextLayerRenderParameters } from '@idea/pdfjs-dist/types/src/display/text_layer'

function isTextItem(item: TextItem | TextMarkedContent): item is TextItem {
  return (item as TextItem).str !== undefined
}

const DEFAULT_FONT_SIZE = 30;


export default class TextLayerRenderTask {

  #tempTextCanvasCtx;

  _textContentSource;
  _isReadableStream;
  _reader: ReadableStreamDefaultReader<TextContent> | null;
  _canceled;

  _textContent: TextContent | null = null;

  constructor({
    textContentSource,
  }: Pick<TextLayerRenderParameters, 'textContentSource' | 'viewport'>) {
    this._textContentSource = textContentSource;
    this._isReadableStream = textContentSource instanceof ReadableStream;

    this._reader = null;

    this._canceled = false;

    // The temporary canvas is used to measure text length in the DOM.
    const canvas = document.createElement("canvas");
    canvas.height = canvas.width = DEFAULT_FONT_SIZE;
    this.#tempTextCanvasCtx = canvas.getContext("2d", { alpha: false })!;

  }

  /**
   * Cancel rendering of the textLayer.
   */
  cancel() {
    this._canceled = true;
    if (this._reader) {
      this._reader
        .cancel(new SharedUtil.AbortException("TextLayer task cancelled."))
        .catch(() => {
          // Avoid "Uncaught promise" messages in the console.
        });
      this._reader = null;
    }
  }

  updateTextContentBounds(viewport: PageViewport) {
    return this.buildTextContentBounds(this._textContent!, viewport)
  }

  buildTextContentBounds(textContent: TextContent, viewport: PageViewport) {
    const textContentBounds: TextContentBound[] = []
    textContent.items.forEach((item) => {
      if (isTextItem(item)) {
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        const textProperties = getTextProperties({
          viewport,
          geom: item,
          ctx: this.#tempTextCanvasCtx,
          styles: textContent.styles,
        })
        if (textProperties) {
          textContentBounds.push({
            geom: item as TextItem & { charsArray: TextContentChar[] },
            str: item.str,
            ...textProperties,
          })
        }

      }

    })
    return textContentBounds
  }


  /**
   * @private
   */
  async render(viewport: PageViewport) {
    if (!this._isReadableStream) {
      return this.buildTextContentBounds(this._textContentSource as TextContent, viewport)
    }
    this._reader = (this._textContentSource as ReadableStream<TextContent>).getReader();

    // eslint-disable-next-line promise/catch-or-return
    this._reader.closed.finally(() => {
      this._reader?.releaseLock();
    })

    const tmpTextContent: TextContent = {
      items: [],
      styles: {},
    }

    const pump = async (): Promise<void> => {
      const res = await (this._reader as ReadableStreamDefaultReader<TextContent>).read()
      if (res.done === true) {
        return
      }

      res.value?.items.forEach((item) => {
        if (isTextItem(item)) {
          tmpTextContent.items.push(item);
        }
      });

      Object.assign(tmpTextContent.styles, res.value?.styles);
      
      return pump();
    }

    try {
      await pump();
    } catch (error) {}

    this._textContent = tmpTextContent;

    return this.buildTextContentBounds(tmpTextContent, viewport)
  }
}