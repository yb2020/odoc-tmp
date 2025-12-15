import {
  attrSelector,
  PDF_ANNOTATE_ID,
  PDF_ANNOTATIONLAYER_GROUP_NOTES,
  ToolBarType,
} from '../constants';
import { PDFJSAnnotate } from '..';
import { appendAnnotation } from '../render/appendChild';
import {
  disableUserSelect,
  enableUserSelect,
  findSVGAtPoint,
  getMetadata,
  scaleDown,
} from './utils';
import {
  repaintRect,
} from '../render/renderRect';
import { ToolBarEndEvent } from '../constants';
import setAttributes from '../utils/setAttributes';

export class RectController {
  public PDFJSAnnotateInstance: PDFJSAnnotate;
  private enabled = false;
  private path: SVGRectElement | null = null;
  private lines: any[] = [];
  private options: any = {};
  private svg: SVGSVGElement | null = null;

  public constructor(PDFJSAnnotateInstance: PDFJSAnnotate) {
    this.PDFJSAnnotateInstance = PDFJSAnnotateInstance;
  }

  private onMouseDown = (e: MouseEvent) => {
    this.path = null;
    this.svg = null;

    this.startPoint(e.clientX, e.clientY);

    document.addEventListener('pointermove', this.onMouseMove);
    document.addEventListener('pointerup', this.onMouseUp);
  };

  private onMouseUp = async () => {
    const svg = this.svg;
    const g = svg?.querySelector<SVGGElement>(
      `.${PDF_ANNOTATIONLAYER_GROUP_NOTES}`
    );

    if (
      this.lines.length &&
      svg &&
      g &&
      this.lines[0].width > 5 &&
      this.lines[0].height > 5
    ) {
      const { pageNumber } = getMetadata(svg);

      const options = {
        type: ToolBarType.rect,
        rectangles: this.lines,
        opacity: 1,
        documentId: this.PDFJSAnnotateInstance.documentId,
        pageNumber,
        tags: [],
        ...this.options,
      };

      const rect = appendAnnotation(g!, options, true) as SVGRectElement;

      if (this.path) {
        g.removeChild(this.path);
      }

      this.PDFJSAnnotateInstance.UI.emit(
        ToolBarEndEvent.EventEnd,
        options,
        svg,
        rect
      );
    }

    document.removeEventListener('pointermove', this.onMouseMove);
    document.removeEventListener('pointerup', this.onMouseUp);
  };

  private onKeyUp = (e: KeyboardEvent) => {
    if (!this.path) {
      return;
    }

    // Cancel rect if Esc is pressed
    if (e.key === 'Escape') {
      this.path.remove();
      document.removeEventListener('pointermove', this.onMouseMove);
      document.removeEventListener('pointerup', this.onMouseUp);
    }
  };

  private createPoint(x: number, y: number) {
    const svg = findSVGAtPoint(x, y);
    if (!svg) {
      return null;
    }

    this.svg = svg;

    const rect = svg.getBoundingClientRect();
    const point = scaleDown(svg, {
      x: x - rect.left,
      y: y - rect.top,
    });

    return { svg, point };
  }

  private startPoint(x: number, y: number) {
    const result = this.createPoint(x, y);
    if (!result) {
      return;
    }

    const { svg, point } = result;
    const g = svg.querySelector<SVGGElement>(
      `g.${PDF_ANNOTATIONLAYER_GROUP_NOTES}`
    )!;

    this.lines = [
      {
        x: point.x,
        y: point.y,
        width: 0,
        height: 0,
      },
    ];

    this.path = appendAnnotation(
      g,
      {
        ...this.options,
        type: ToolBarType.rect,
        rectangles: this.lines,
      },
      true
    ) as SVGRectElement;
  }

  private onMouseMove = ({ clientX, clientY }: MouseEvent) => {
    const result = this.createPoint(clientX, clientY);
    if (!result) {
      return;
    }

    const { point } = result;

    if (!this.path) {
      return;
    }

    const line = this.lines[0];

    this.lines[0] = {
      ...line,
      width: point.x - line.x,
      height: point.y - line.y,
    };

    repaintRect(this.path, this.lines[0]);
  };

  public setOptions(options: Record<string, any>) {
    this.options = {
      ...this.options,
      ...options,
    };
  }

  public setRectToDom(annotateId: string, attributes: Record<string, any>) {
    const container =
      this.PDFJSAnnotateInstance.pdfWebview.getDocumentViewer().container;

    const rect = container.querySelector(
      `rect${attrSelector(PDF_ANNOTATE_ID, annotateId)}`
    ) as SVGRectElement;

    if (!rect) {
      return {};
    }

    const svg = rect.closest('svg')!;

    const { documentId, pageNumber } = getMetadata(svg);

    setAttributes(rect, attributes);

    return { documentId, pageNumber };
  }

  public deleteRectToDom(annotateId: string) {
    const container =
      this.PDFJSAnnotateInstance.pdfWebview.getDocumentViewer().container;

    const node = container.querySelector(
      `rect${attrSelector(PDF_ANNOTATE_ID, annotateId)}`
    );

    if (!node) {
      return;
    }

    const svg = node.closest('svg');

    const { documentId } = getMetadata(svg!);

    node?.parentNode?.removeChild(node);

    return documentId;
  }

  public updateUuid(rect: SVGRectElement, uuid: string) {
    setAttributes(rect, {
      [PDF_ANNOTATE_ID]: uuid,
      uuid,
    });
  }

  public enable() {
    if (this.enabled) {
      return;
    }

    this.enabled = true;
    document.addEventListener('pointerdown', this.onMouseDown);
    document.addEventListener('keyup', this.onKeyUp);
    disableUserSelect();
  }

  public disable() {
    if (!this.enable) {
      return;
    }

    this.enabled = false;
    document.removeEventListener('pointerdown', this.onMouseDown);
    document.removeEventListener('keyup', this.onKeyUp);
    enableUserSelect();
  }
}
