import { ViewerController } from '@idea/pdf-annotate-viewer';
import {
  PDF_ANNOTATE_SHAPE_READONLY,
  PDF_ANNOTATE_TYPE,
  PDF_ANNOTATE_USERID,
  PDF_ANNOTATIONLAYER_GROUP_AI,
} from './constants';
import { UI } from './UI';
import { setDisplay } from './UI/utils';
export { setDisplay, scaleUpRaw, scaleDownRaw, SVG_CONTAINER_SELECTOR } from './UI/utils';
export { numberToPx, pxToNumber } from './utils/setAttributes';
export {
  renderAnnotate,
  renderAnnotationSvg,
  prepareTextsShapes,
  ANNOTATION_PAGE_RENDERED,
  ANNOTATION_PAGESVG_RENDERED
} from './render';
export type { 
  AnnotationListGroupByVisible,
  PageHandwriteMap,
  PageShapeTextMap,
  PageShapeMap
} from './render';
export * from './render/renderTextBox';
export * from './render/shapeCommon';
export {
  renderSinglePageHandwrite,
  getIpadUnit,
} from './render/renderHandwrite';
export * from './render/renderShape';
export * from './render/editHandwrite';
export * from './constants';
export * from './constants/color';
export { ANNOTATION_MOUSEOVER, ANNOTATION_CLICK } from './UI/event';
export { Canvas, Rect, Ellipse, Line } from 'fabric';

export class PDFJSAnnotate {
  public svgElements = new Map<number, SVGSVGElement>();
  public canvasElements = new Map<number, HTMLCanvasElement>();
  public canvasRendering = new Map<number, boolean>();
  public divElements = new Map<number, HTMLDivElement>();
  public UI: UI;

  public constructor(
    public documentId: string,
    public pdfWebview: ViewerController
  ) {
    this.UI = new UI(this);
  }

  public setDisplayByType(
    type: string,
    typeVisible: boolean,
    otherVisible: boolean
  ) {
    this.setDisplayByAttribute(
      PDF_ANNOTATE_TYPE,
      type,
      typeVisible,
      otherVisible
    );
  }

  public setDisplayByUserId(
    userId: string,
    userVisible: boolean,
    otherVisible: boolean
  ) {
    this.setDisplayByAttribute(
      PDF_ANNOTATE_USERID,
      userId,
      userVisible,
      otherVisible
    );
  }

  private setDisplayByAttribute(
    attrName: string,
    attrValue: string,
    matchVisible: boolean,
    notMatchVisible: boolean
  ) {
    this.svgElements.forEach((elm) => {
      Array.from(elm.querySelectorAll('[uuid]')).forEach((child) => {
        const value = child.getAttribute(attrName);
        setDisplay(
          child as SVGElement,
          value === attrValue ? matchVisible : notMatchVisible
        );
      });
    });
  }

  public setDisplayHandwrite(visible: boolean) {
    this.canvasElements.forEach((canvas) => {
      setDisplay(canvas, visible);
    });
    this.divElements.forEach((div) => {
      setDisplay(div, visible);
    });
  }

  public onOffAnnotation(onOff: boolean) {
    this.setTextShapeReadonly(!onOff);
    this.UI.onOffMouseMove(onOff);
  }

  private setTextShapeReadonly(readonly: boolean) {
    this.divElements.forEach(
      readonly
        ? (div) => {
            if (!div.classList.contains(PDF_ANNOTATE_SHAPE_READONLY)) {
              div.classList.add(PDF_ANNOTATE_SHAPE_READONLY);
            }
          }
        : (div) => {
            if (div.classList.contains(PDF_ANNOTATE_SHAPE_READONLY)) {
              div.classList.remove(PDF_ANNOTATE_SHAPE_READONLY);
            }
          }
    );
  }

  public clearPersonalNotes() {
    const remove = (item: Element) => item.remove();
    this.svgElements.forEach((svgEl) => {
      Array.from(svgEl.children)
        .filter((x) => !x.classList.contains(PDF_ANNOTATIONLAYER_GROUP_AI))
        .forEach(remove);
    });
    this.canvasElements.forEach(remove);
    this.canvasElements.clear();
    this.divElements.forEach(remove);
    this.divElements.clear();
  }
}
