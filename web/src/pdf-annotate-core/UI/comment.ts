import { attrSelector, ToolBarType } from '../constants';
import { PDF_ANNOTATE_ID } from '../constants/index';
import { Annotation, appendAnnotation } from '../render/appendChild';
import {
  CacheSvg,
  getMetadata,
  scaleDown,
  splitAnnotation,
  SvgBound,
} from './utils';

import setAttributes from '../utils/setAttributes';
import { SelectAnnotation } from '../render/renderSelect';
import { ViewerController } from '@idea/pdf-annotate-viewer';

export class CommentController {
  constructor(public pdfWebview: ViewerController) {}

  addComment(options: SelectAnnotation, dry = false) {
    options.rectStr = options.rectStr || '';
    options.rectangles = options.rectangles || [];

    const cacheSvg = new CacheSvg(this.pdfWebview);
    const svg = cacheSvg.findByRect(options.rectangles[0]);

    if (!svg) {
      return;
    }

    // Initialize the annotation
    const { rectRaw, ...rest } = options;
    const annotation: SelectAnnotation = {
      ...rest,
      type: ToolBarType.select,
      tags: [],
      rectangles: options.rectangles
        .map((rect) => {
          const svg = cacheSvg.findByRect(rect);
          if (!svg) {
            return null as any;
          }

          return {
            ...rect,
            ...(rectRaw ? {} : scaleDown(svg.svg, {
              y: rect.y,
              x: rect.x,
              width: rect.width,
              height: rect.height,
            })),
          };
        })
        .filter(Boolean)
        .filter(
          (rect) =>
            rect.width > 0 && rect.height > 0 && rect.x > -1 && rect.y > -1
        ),
    };

    if (annotation.rectangles.length === 0) {
      return;
    }

    const { documentId, pageNumber } = getMetadata(svg.svg);

    let svgGroupElementList;
    if (!dry) {
      svgGroupElementList = this.addCommentToDom(annotation, true, cacheSvg);
    }

    return {
      svgGroupElementList,
      documentId,
      pageNumber,
      annotation,
    };
  }

  addCommentToDom(
    annotation: Annotation,
    visible: boolean,
    cacheSvg = new CacheSvg(this.pdfWebview)
  ) {
    const cache = cacheSvg.findByPage(annotation.pageNumber as number);

    if (!cache) {
      return;
    }

    const container = this.pdfWebview.getDocumentViewer().container;

    const hasGroup = container.querySelectorAll<SVGGElement>(
      attrSelector(PDF_ANNOTATE_ID, annotation.uuid)
    );

    if (hasGroup && hasGroup.length) {
      return Array.from(hasGroup);
    }

    const groupList = splitAnnotation(annotation as SelectAnnotation).map(
      (anno) => {
        const svg = cacheSvg.findByPage(anno.pageNumber) as SvgBound;
        const group = appendAnnotation(svg.g, anno, visible);
        if (group) {
          this.updateUuid(group, annotation.uuid as string);
        }
        return group;
      }
    );

    return groupList;
  }

  updateUuid(group: SVGGElement, uuid: string) {
    setAttributes(group, {
      [PDF_ANNOTATE_ID]: uuid,
      uuid,
    });
  }

  setCommentToDom(annotateId: string, attributes: Record<string, any>) {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const { pageNumber, ...rest } = attributes;

    const container = this.pdfWebview.getDocumentViewer().container;

    const groupList = container.querySelectorAll<SVGGElement>(
      `g${attrSelector(PDF_ANNOTATE_ID, annotateId)}`
    );

    Array.from(groupList).forEach((group) => {
      console.log({ group, attributes });
      setAttributes(group, rest);
    });
  }

  deleteCommentToDom(annotateId: string) {
    const container = this.pdfWebview.getDocumentViewer().container;

    const list = container.querySelectorAll<SVGGElement>(
      `g${attrSelector(PDF_ANNOTATE_ID, annotateId)}`
    );

    Array.from(list).forEach((g) => {
      const svg = g.parentElement;
      svg?.removeChild(g);
    });
  }
}
