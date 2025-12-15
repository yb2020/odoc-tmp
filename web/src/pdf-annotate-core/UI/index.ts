import { PDFJSAnnotate, PDF_ANNOTATE_SVG_READONLY } from '..';
import EventEmitter from 'events';
import isMobile from 'is-mobile';
import { useAnnotationMouseMove, useAnnotationMouseDown } from './event';
// import { disableEdit, enableEdit } from './edit';
// import { disablePen, enablePen, setPen } from './pen';
// import { disablePoint, enablePoint } from './point';
// import { disableText, enableText, setText } from './text';

// import { enableSelect, disableSelect } from './select';
// import { enableLine, setLine, disableLine } from './line';
import { RectController } from './rect';

// import { addUnderline, setUnderline, deleteUnderline } from './underline';
// import { addHighlight, setHighlight, deleteHighlight } from './highlight';
import { CommentController } from './comment';

// import { enableCircle, disableCircle, setCircle } from './circle';

// import undoManager from './utils/undo';

export class UI extends EventEmitter {
  protected disabled = false;
  protected mousedownFlags = new Map<Element, boolean>();

  public rectController: RectController;
  public commentController: CommentController;

  public initAnnotationMouseDown: ReturnType<typeof useAnnotationMouseDown>;
  private initAnnotationMouseMove: ReturnType<typeof useAnnotationMouseMove>;

  public constructor(public PDFJSAnnotateInstance: PDFJSAnnotate) {
    super();

    this.rectController = new RectController(PDFJSAnnotateInstance);

    this.commentController = new CommentController(
      PDFJSAnnotateInstance.pdfWebview
    );

    this.initAnnotationMouseDown = useAnnotationMouseDown(this);
    this.initAnnotationMouseMove = useAnnotationMouseMove(this);
    this.onOffMouseMove(true);
  }

  public onOffMouseMove(onOff: boolean) {
    const MOUSEMOVE = isMobile({
      tablet: true,
    })
      ? 'pointerdown'
      : 'pointermove';
    const hasClass = document.body.classList.contains(
      PDF_ANNOTATE_SVG_READONLY
    );

    document.removeEventListener(MOUSEMOVE, this.initAnnotationMouseMove);

    if (onOff) {
      document.addEventListener(MOUSEMOVE, this.initAnnotationMouseMove);

      if (hasClass) {
        document.body.classList.remove(PDF_ANNOTATE_SVG_READONLY);
      }
    } else {
      if (!hasClass) {
        document.body.classList.add(PDF_ANNOTATE_SVG_READONLY);
      }
    }
  }

  initMouseDown(el: SVGSVGElement | HTMLDivElement) {
    if (!this.mousedownFlags.get(el)) {
      this.initAnnotationMouseDown(el);
      this.mousedownFlags.set(el, true);
    }
  }

  enable() {
    this.disabled = false;
  }

  disable() {
    this.disabled = true;
  }

  emit(eventName: string | symbol, ...args: any[]) {
    if (this.disabled) {
      return false;
    }

    return super.emit(eventName, ...args);
  }
}
