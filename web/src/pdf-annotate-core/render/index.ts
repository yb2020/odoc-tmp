// import { WebDrawV2 } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Web';
import { WebDrawV2 } from 'go-sea-proto/gen/ts/note/web';
// import { ShapeAnnotation } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { ShapeAnnotation } from 'go-sea-proto/gen/ts/common/ShapeAnnotation';

import { PDFPageView } from '@idea/pdfjs-dist/web/pdf_viewer';
import { Annotation, appendAnnotation } from './appendChild';
import {
  PDFJSAnnotate,
  DATA_OLD_HANDWRITE_FRAME,
  setDisplay,
  PageFcanvasMap,
  clearShapeToolbar,
  numberToPx,
} from '..';
import {
  DATA_PAGE_NUMBER,
  PDF_ANNOTATE_CONTAINER,
  PDF_ANNOTATE_DOCUMENT,
  PDF_ANNOTATE_TEXTSHAPE_BOX,
  PDF_ANNOTATE_SHAPE_DIV,
  PDF_ANNOTATE_TEXTSHAPE_SCALE,
  PDF_ANNOTATE_TEXTSHAPE_SPAN,
  PDF_ANNOTATE_TEXTSHAPE_STYLE,
  PDF_ANNOTATE_HANDWRITE_CANVAS,
  PDF_ANNOTATE_PAGENUMBER,
  PDF_ANNOTATE_VIEWPORT,
  PDF_ANNOTATIONLAYER,
  PDF_ANNOTATIONLAYER_GROUP_AI,
  PDF_ANNOTATIONLAYER_GROUP_NOTES,
  PDF_ANNOTATIONLAYER_GROUP_VOCABULARY,
  READPAPER_PAGE_CONTAINER_,
  PDF_ANNOTATE_SHAPE_READONLY,
  PDF_ANNOTATE_SVG_READONLY,
} from '../constants';
import { getIpadUnit, renderSinglePageHandwrite } from './renderHandwrite';
import createStyleSheet from 'create-stylesheet';
import * as fabric from 'fabric';
import { createArrowhead, initShape, newShape } from './renderShape';
import {
  PageTextRectMap,
  TextAnnotation,
  conitnueEditTextarea,
  createTextDiv,
  editingTextarea,
  getTextBody,
  onMouseDownStartEditText,
  useStartEditText,
} from './renderTextBox';
import {
  PrepareCanvasProps,
  RenderAnnotateProps,
  RenderCanvasProps,
} from './type';

export type AnnotationListGroupByVisible = [Annotation[], boolean][];

export type PageHandwriteMap = Record<string | number, WebDrawV2[]>;
export type PageShapeTextMap = Record<string | number, TextAnnotation[]>;
export type PageShapeMap = Record<string | number, ShapeAnnotation[]>;

export const X_BEFORE_SCALE_CHANGING = 'x-beforescalechanging';
export const BEFORE_DESTROY = 'readpaper-before-destroy';
export const ANNOTATION_PAGE_RENDERED = 'annotation:pageRendered';
export const ANNOTATION_PAGESVG_RENDERED = 'annotation:pageSvgRendered';

function setFullHeightWidthStyle(
  element: SVGSVGElement | HTMLCanvasElement | HTMLDivElement
) {
  Object.assign(element.style, {
    position: 'absolute',
    top: '0',
    left: '0',
    width: '100%',
    height: '100%',
  } as CSSStyleDeclaration);
}

export const setCitationsTarget = ({
  source,
}: Pick<RenderAnnotateProps, 'source'>) => {
  const referenceLayer = source.div.querySelector(`div.${PDF_ANNOTATIONLAYER}`);

  referenceLayer?.querySelectorAll('a').forEach((item) => {
    item.setAttribute('target', '_blank');
  });
};

export const renderAnnotate = async (props: RenderAnnotateProps) => {
  const { pageNumber, source, instance } = props;
  if (instance.canvasRendering.get(pageNumber)) {
    return;
  }
  instance.canvasRendering.set(pageNumber, true);

  source.div.id = READPAPER_PAGE_CONTAINER_ + pageNumber;
  source.div.setAttribute(DATA_PAGE_NUMBER, pageNumber + '');

  const params = { ...props, pageView: source };

  await renderHandwrites(params);
  await renderTextsShapes(params);

  instance.canvasRendering.set(pageNumber, false);
  instance.UI.emit(ANNOTATION_PAGE_RENDERED, { pageNumber, source });
};

/**
 * @description 后续将数据对齐fabricjs合并canvas
 */
const renderHandwrites = async ({
  source,
  instance,
  pageNumber,
  pageView,
  handwriteBuffer,
  handwriteAndShapeVisible,
}: RenderCanvasProps) => {
  const map = await handwriteBuffer?.();
  const list = map?.[pageNumber] || [];
  if (!list?.length) {
    return;
  }

  const pageIdx = pageNumber - 1;
  let handwriteCanvas: HTMLCanvasElement;

  if (!instance.canvasElements.get(pageIdx)) {
    handwriteCanvas = document.createElement('canvas');
    instance.canvasElements.set(pageIdx, handwriteCanvas);
    handwriteCanvas.setAttribute(
      PDF_ANNOTATE_HANDWRITE_CANVAS,
      String(pageNumber)
    );
    handwriteCanvas.classList.add('handwrite-canvas');
    handwriteCanvas.style.pointerEvents = 'none';
    setDisplay(handwriteCanvas, handwriteAndShapeVisible);
  } else {
    handwriteCanvas = instance.canvasElements.get(pageIdx) as HTMLCanvasElement;
    handwriteCanvas
      .getContext('2d')
      ?.clearRect(0, 0, handwriteCanvas.width, handwriteCanvas.height);
  }

  // renderShapeDiv();
  setFullHeightWidthStyle(handwriteCanvas);

  const { height, width, scale } = pageView.viewport;

  handwriteCanvas.height = height;
  handwriteCanvas.width = width;

  renderSinglePageHandwrite(
    pageNumber,
    list,
    instance,
    false,
    scale,
    handwriteCanvas
  );

  source.div.appendChild(handwriteCanvas);
};

export const prepareTextsShapes = (props: PrepareCanvasProps) => {
  prepareTextsShapesDiv({
    ...props,
    shapeEditable: true,
  });

  prepareTexts(props);
  prepareShapes(props);
};

const prepareTextsShapesDiv = (params: PrepareCanvasProps) => {
  const { instance, pageNumber, pageView, shapeEditable } = params;
  const pageIdx = pageNumber - 1;
  // 这里不能用instance.divElements.get(pageIdx) 因为pageView可能被销毁
  let pageShape = pageView.div.querySelector<HTMLDivElement>(
    `.${PDF_ANNOTATE_SHAPE_DIV}`
  );
  if (!pageShape) {
    instance.divElements.delete(pageIdx);
    pageShape = document.createElement('div');
    pageShape.setAttribute(PDF_ANNOTATE_SHAPE_DIV, String(pageNumber));
    pageShape.classList.add(PDF_ANNOTATE_SHAPE_DIV);
    if (!shapeEditable) {
      pageShape.classList.add(PDF_ANNOTATE_SHAPE_READONLY);
    }
    pageView.div.appendChild(pageShape);
    instance.divElements.set(pageIdx, pageShape);
  }
};

const renderTextsShapes = async (params: RenderCanvasProps) => {
  const { instance } = params;
  const { container } = instance.pdfWebview.getDocumentViewer();
  clearShapeToolbar(container);
  ensureShapeClassStyle();

  prepareTextsShapesDiv(params);
  await renderTexts(params);
  await renderShapes(params);
};

const prepareTexts = ({
  instance,
  pageNumber,
  pageView,
}: PrepareCanvasProps) => {
  const pageIdx = pageNumber - 1;
  const shapeDiv = instance.divElements.get(pageIdx);
  const scale = `scale(${pageView.viewport.scale})`;
  let textsDiv = shapeDiv?.querySelector<HTMLDivElement>(
    `.${PDF_ANNOTATE_TEXTSHAPE_SCALE}`
  );
  if (textsDiv) {
    textsDiv.style.transform = scale;
    return textsDiv;
  }

  textsDiv = document.createElement('div');
  shapeDiv?.appendChild(textsDiv);
  textsDiv.classList.add(PDF_ANNOTATE_TEXTSHAPE_SCALE);
  textsDiv.style.fontSize = numberToPx(getIpadUnit(instance, pageIdx));
  textsDiv.style.transform = scale;

  return textsDiv;
};

const renderTexts = async (props: RenderCanvasProps) => {
  const { instance, pageNumber, pageView, textCallback, shapeTextBuffer } =
    props;
  const { container } = instance.pdfWebview.getDocumentViewer();
  const pageIdx = pageNumber - 1;
  const shapeDiv = instance.divElements.get(pageIdx);
  const textList = (await shapeTextBuffer?.())?.[pageNumber] ?? [];
  if (!shapeDiv || !textCallback || !textList.length) {
    return;
  }

  const textsDiv = prepareTexts(props);
  if (!textsDiv) {
    return;
  }

  textList.forEach((textItem) => {
    if (
      (!getTextBody(textItem)?.content &&
        textItem.textBox?.id !== editingTextarea.id) ||
      textsDiv.querySelector(
        `[${PDF_ANNOTATE_TEXTSHAPE_BOX}="${textItem.textBox?.id}"]`
      )
    ) {
      return;
    }

    const textBoxDiv = createTextDiv(textItem, container);
    const startEditText = useStartEditText(
      pageView.viewport.scale,
      textItem,
      textBoxDiv,
      textCallback,
      () => prepareShapes(props)
    );
    onMouseDownStartEditText(textBoxDiv, startEditText);
    textsDiv.appendChild(textBoxDiv);

    if (editingTextarea.id === textItem.textBox?.id) {
      setTimeout(() => {
        startEditText();
        conitnueEditTextarea(textBoxDiv);
      }, 1);
    }
  });
};

const prepareShapes = ({
  instance,
  pageNumber,
  pageView,
}: PrepareCanvasProps) => {
  const pageIdx = pageNumber - 1;
  const shapeDiv = instance.divElements.get(pageIdx);
  let fcanvas = PageFcanvasMap.get(pageNumber);

  if (!fcanvas) {
    const canvasElement = document.createElement('canvas');
    const { height, width } = pageView.viewport;
    canvasElement.id = DATA_SHAPE_CANVAS + pageNumber;
    canvasElement.setAttribute(DATA_SHAPE_CANVAS, String(pageNumber));
    canvasElement.height = height;
    canvasElement.width = width;

    Object.assign(canvasElement.style, {
      position: 'absolute',
      top: '0px',
      left: '0px',
    });

    shapeDiv?.appendChild(canvasElement);
    fcanvas = new fabric.Canvas(canvasElement.id, {
      selection: false,
      uniformScaling: false,
      defaultCursor: 'inherit',
    });
    PageFcanvasMap.set(pageNumber, fcanvas);
  }

  return fcanvas;
};

const renderShapes = async (props: RenderCanvasProps) => {
  const { instance, pageNumber, pageView, shapeBuffer, shapeCallback } = props;
  const pageIdx = pageNumber - 1;
  const shapeDiv = instance.divElements.get(pageIdx);
  const pageShapeMap = await shapeBuffer?.();
  const shapeList = pageShapeMap?.[pageNumber] ?? [];
  if (!shapeDiv || !shapeCallback || !shapeList.length) {
    return;
  }

  const deprecated = PageFcanvasMap.get(pageNumber);
  const deprecatedWrapper = deprecated
    ? shapeDiv.querySelector<HTMLDivElement>(`.${deprecated.containerClass}`)
    : null;
  if (deprecated && deprecatedWrapper !== deprecated.wrapperEl) {
    await deprecated.dispose();
    // deprecated.discardActiveObject();
    PageTextRectMap.delete(pageNumber);
    PageFcanvasMap.delete(pageNumber);
  }

  const fcanvas = prepareShapes(props);
  // @TODO 不要依赖text的DOM
  const textsDiv = prepareTexts(props);
  shapeList.forEach((item) => {
    const fobj = newShape(item, pageView.viewport.scale);
    initShape(item, fobj, pageView, shapeCallback);
    fcanvas.add(fobj);
    if (fobj instanceof fabric.Line) {
      const arrowhead = createArrowhead(
        item.shapeId || item.uuid,
        fobj,
        pageView.viewport.scale
      );
      textsDiv.appendChild(arrowhead);
    }
  });
};

export interface RenderAnnotationSvgProps {
  documentId: string;
  pageNumber: number;
  source: PDFPageView;
  viewport: PDFPageView['viewport'];
  instance: PDFJSAnnotate;
  annotationsAI?: Annotation[];
  annotationsVocabulary?: Annotation[];
  annotationsExtractGrouped?: AnnotationListGroupByVisible;
}

export function renderAnnotationSvg({
  documentId,
  pageNumber,
  instance,
  source,
  viewport = source.viewport,
  annotationsAI,
  annotationsVocabulary,
  annotationsExtractGrouped,
}: RenderAnnotationSvgProps) {
  let svg = source.div.querySelector<SVGSVGElement>(
    `svg.${PDF_ANNOTATIONLAYER}`
  );
  let gAI = svg?.querySelector<SVGGElement>(
    `g.${PDF_ANNOTATIONLAYER_GROUP_AI}`
  );
  let gNotes = svg?.querySelector<SVGGElement>(
    `g.${PDF_ANNOTATIONLAYER_GROUP_NOTES}`
  );
  let gVocabulary = svg?.querySelector<SVGGElement>(
    `g.${PDF_ANNOTATIONLAYER_GROUP_VOCABULARY}`
  );

  instance.UI.initMouseDown(source.div);

  const createSvg = () => {
    const el = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    el.classList.add(PDF_ANNOTATIONLAYER);
    setFullHeightWidthStyle(el);
    source.div.insertBefore(el, source.div.firstElementChild);

    el.setAttribute(PDF_ANNOTATE_CONTAINER, String(true));
    el.setAttribute(PDF_ANNOTATE_VIEWPORT, JSON.stringify(source.viewport));
    el.setAttribute(PDF_ANNOTATE_DOCUMENT, documentId);
    el.setAttribute(PDF_ANNOTATE_PAGENUMBER, String(pageNumber));

    return el;
  };
  const createG = (svg: SVGElement, klass: string) => {
    const g = document.createElementNS('http://www.w3.org/2000/svg', 'g');
    g.classList.add(klass);

    svg.appendChild(g);

    return g;
  };

  if (!svg) {
    svg = createSvg();
  }
  if (!gAI) {
    gAI = createG(svg, PDF_ANNOTATIONLAYER_GROUP_AI);
  }
  if (!gNotes) {
    gNotes = createG(svg, PDF_ANNOTATIONLAYER_GROUP_NOTES);
  }
  if (!gVocabulary) {
    gVocabulary = createG(svg, PDF_ANNOTATIONLAYER_GROUP_VOCABULARY);
  }

  if (gAI && Array.isArray(annotationsAI)) {
    gAI.innerHTML = '';
    annotationsAI.forEach((annotation) => {
      appendAnnotation(gAI!, annotation, true, viewport);
    });
  }

  if (gVocabulary && Array.isArray(annotationsVocabulary)) {
    gVocabulary.innerHTML = '';
    annotationsVocabulary.forEach((annotation) => {
      appendAnnotation(gVocabulary!, annotation, true, viewport);
    });
  }

  if (gNotes && Array.isArray(annotationsExtractGrouped)) {
    gNotes.innerHTML = '';
    annotationsExtractGrouped?.forEach(([annotationList, visible]) => {
      annotationList.forEach((annotation) => {
        appendAnnotation(gNotes!, annotation, visible, source.viewport);
      });
    });
  }

  instance.svgElements.set(pageNumber - 1, svg);

  instance.UI.emit(ANNOTATION_PAGESVG_RENDERED, { pageNumber, source });
}

function ensureShapeClassStyle() {
  if (document.getElementById(PDF_ANNOTATE_TEXTSHAPE_STYLE)) {
    return;
  }

  const style = createStyleSheet({
    [`.${PDF_ANNOTATE_SVG_READONLY} svg *`]: {
      pointerEvents: 'none',
    },
    [`.${PDF_ANNOTATE_SHAPE_DIV}.${PDF_ANNOTATE_SHAPE_READONLY} > div`]: {
      pointerEvents: 'none',
    },
    [`.${PDF_ANNOTATE_SHAPE_DIV}`]: {
      position: 'absolute',
      top: '0',
      left: '0',
      right: '0',
      bottom: '0',
      overflow: 'hidden',
    },
    ['.' + PDF_ANNOTATE_TEXTSHAPE_SCALE]: {
      position: 'absolute',
      top: '0',
      left: '0',
    },
    ['.' + PDF_ANNOTATE_TEXTSHAPE_BOX]: {
      position: 'absolute',
      overflow: 'hidden',
      wordBreak: 'break-word',
      textOverflow: 'ellipsis',
      lineHeight: '0',
    },
    [`.${PDF_ANNOTATE_TEXTSHAPE_SPAN} > br`]: {
      height: '0',
    },
    [`.${DATA_OLD_HANDWRITE_FRAME}`]: {
      position: 'absolute',
      cursor: 'pointer',
    },
  });

  style.id = PDF_ANNOTATE_TEXTSHAPE_STYLE;

  document.head.appendChild(style);
}

const DATA_SHAPE_CANVAS = 'data-shape-canvas';
