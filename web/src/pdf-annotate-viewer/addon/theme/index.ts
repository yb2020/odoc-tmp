/* eslint-disable prefer-rest-params */
/* eslint-disable @typescript-eslint/ban-ts-comment */
import ColorConvert from './ColorConvert.js';
import * as pdfjsViewer from '@idea/pdfjs-dist/web/pdf_viewer';

import ColorsToneConfig from './colors.json'

interface WrappedCanvasRenderingContext2D extends CanvasRenderingContext2D {
  origFillRect?: CanvasRenderingContext2D['fillRect'];
}


type CanvasFuncName = Extract<
  keyof CanvasRenderingContext2D, 
  'fill' | 'fillRect' | 'fillText' | 'stroke' | 'strokeRect' | 'strokeText'
>;

interface ColorScheme {
  name: string;
  tones: ColorTone[];
}

interface ColorTone {
  name: string;
  background: string;
  foreground: string;
  colors: {
    bg: ColorConvert;
    fg: ColorConvert;
    acc: ColorConvert[];
    grad: any;
  };
  baseColors: ColorConvert[];
}

export default class PDFViewerColorToneAddon {
  // private pdfViewerController: null | ViewerController = null

  private colorSchemes: null | ColorScheme[] = null

  private styleCache = new Map()
  private readerTone: ColorTone | null = null

  private flags = { readerOn: false, isPrinting: false, shapesOn: false }

  private eventBus = new pdfjsViewer.EventBus();

  static ForceRedrawEventType = 'x-pdf-viewer-colortone-addon-forceRedraw'

  private checkFlags = () => {
    return this.flags.readerOn && !this.flags.isPrinting;
  }

  enableReaderTone(redraw?: boolean) {
    this.flags.readerOn = true;
    if (redraw) {
      this.forceRedraw();
    }
  }

  disableReaderTone(redraw?: boolean) {
    this.flags.readerOn = false;
    if (redraw) {
      this.forceRedraw();
    }
  }

  addEventListener(event: string, listener: () => void) {
    this.eventBus.on(event, listener)
  }

  removeListener(event: string, listener?: any) {
    this.eventBus.off(event, listener)
  }

  // setPDFViewer(viewer: ViewerController) {
  //   this.pdfViewerController = viewer
  // }

  forceRedraw() {
    this.eventBus.dispatch(PDFViewerColorToneAddon.ForceRedrawEventType, {})
    // if (!this.pdfViewerController) {
    //   throw Error('PDFViewerColorToneAddon: invalid pdfViewerController')
    // }
    // this.pdfViewerController.redraw()
  }


   /* Calculate a new style for given colorscheme and tone */
   calcStyle(color: ColorConvert) {
    const {grad, acc} = this.readerTone!.colors;

    if (color.chroma > 10) {
      const accents = acc.concat(this.readerTone!.baseColors);
      if (accents.length) {
        const newArr = accents.map(item => item.deltaE(color))
        const min = Math.min(...newArr)
        const index = newArr.indexOf(min)
        const style = accents[index].toHex(color.alpha)
        return style;
      }
    } else {
      const whiteL = ColorConvert.white.lightness;
      const style = grad(1 - color.lightness / whiteL).toHex(color.alpha);
      return style;
    }
  }

  /* Return fill and stroke styles */
  private getReaderStyle = (
    ctx: WrappedCanvasRenderingContext2D, 
    funcName: CanvasFuncName,
    args: any[],
    style: CanvasRenderingContext2D['fillStyle' | 'strokeStyle']
  ) => {
    const isColor = typeof style === 'string';    /* not gradient/pattern */

    if (!isColor) {
      return style;
    }

    const isText = funcName.endsWith('Text');
    const isShape = !isText && !(
      (funcName === 'fillRect') &&
      args[2] == ctx.canvas.width &&
      args[3] == ctx.canvas.height
    );
    if (isShape && !this.flags.shapesOn && style !== '#ffffff') {
      return style
    }

    if (!this.styleCache.has(style)) {
      this.styleCache.set(
        style,
        this.calcStyle(new ColorConvert(style))  
      );
    }

    return this.styleCache.get(style);
  }

  // injectCanvasProxy() {
  //   const CanvasPrototype = CanvasRenderingContext2D.prototype
  //   const fillDescriptor = Object.getOwnPropertyDescriptor(CanvasPrototype, 'fill');
  //   const fillRectDescriptor = Object.getOwnPropertyDescriptor(CanvasRenderingContext2D.prototype, 'fillRect');
  //   const fillTextDescriptor = Object.getOwnPropertyDescriptor(CanvasRenderingContext2D.prototype, 'fillText');
  //   const strokeDescriptor = Object.getOwnPropertyDescriptor(CanvasRenderingContext2D.prototype, 'stroke');
  //   const strokeRectDescriptor = Object.getOwnPropertyDescriptor(CanvasRenderingContext2D.prototype, 'strokeRect');
  //   const strokeTextDescriptor = Object.getOwnPropertyDescriptor(CanvasRenderingContext2D.prototype, 'strokeText');

  //   const proxyWrap = function() {
  //     fillDescriptor?.value.call(this)
  //   }

  //   Object.defineProperty(CanvasPrototype, 'fill', Object.assign({}, fillDescriptor, { value: proxyWrap }) )

  // }

  updateReaderColors() {
    this.styleCache.clear();
  }

  private getColorTone(theme: 'beige' | 'green' | 'dark') {
    if (!this.colorSchemes) {
      throw Error('PDFViewerColorToneAddon: call initial before changeReaderColorTone')
    }
    let idx = 0;
    if (theme === 'green') {
      idx = 1
    } else if (theme === 'dark') {
      idx = 2
    }
    return this.colorSchemes[0].tones[idx];
  }

  changeReaderColorTone(theme: 'beige' | 'green' | 'dark' | 'default') {
    // if (!this.pdfViewerController) {
    //   throw Error('PDFViewerColorToneAddon: invalid pdfViewerController')
    // }

    if (theme === 'dark') {
      this.flags.shapesOn = true
    } else {
      this.flags.shapesOn = false;
    }

    if (theme === 'default') {
      this.disableReaderTone()
      this.forceRedraw()
      return
    }
    
    if (!this.flags.readerOn) {
      this.enableReaderTone()
    }
    this.readerTone = this.getColorTone(theme)
    this.updateReaderColors()
    this.forceRedraw()
  }

  initial(theme?: 'beige' | 'green' | 'dark') {
    this.colorSchemes = ColorsToneConfig.map(scheme => {
      // const baseColors = (scheme.accents || []).map((c) => newColor(c))
      return {
        name: scheme.name,
        tones: scheme.tones.map(tone => {
          const [b, f] = [tone.background, tone.foreground].map((c) => new ColorConvert(c));
          return {
            name: tone.name,
            background: tone.background,
            foreground: tone.foreground,
            colors: {
              bg: b, fg: f, grad: b.range(f),
              acc: (tone.accents || []).map((c) => new ColorConvert(c)),
            },
            baseColors: [],
          } as ColorTone ;
        }),
      }
    });

    if (theme) {
      this.enableReaderTone()
      this.readerTone = this.getColorTone(theme)
      if (theme === 'dark') {
        this.flags.shapesOn = true;
      }
    }

    this.updateReaderColors()
    this.wrapCanvasMethod()
  }

  private wrapCanvasMethod() {
    /* Wrap canvas drawing */
    const ctxp = CanvasRenderingContext2D.prototype as WrappedCanvasRenderingContext2D;

    if (ctxp.origFillRect) {
      return
    }

    ctxp.origFillRect = ctxp.fillRect;

    const { checkFlags, getReaderStyle } = this

    const funcList: CanvasFuncName[] = [
      'fill',
      'fillRect',
      'fillText',
      'stroke',
      'strokeRect',
      'strokeText',
    ]

    funcList.forEach((func) => {
      const style: keyof CanvasRenderingContext2D = func.startsWith('fill') 
        ? 'fillStyle' 
        : 'strokeStyle'

      const originFunc = ctxp[func] as (...args: unknown[]) => unknown

      ctxp[func] = function(...args: unknown[]) {
        if (!checkFlags()) {
          return originFunc.apply(this, args)
        }

        const originStyle = this[style]

        this[style] = getReaderStyle(this, func, args, originStyle)

        const result = originFunc.apply(this, args)

        this[style] = originStyle

        return result
      }
    })
  }
}