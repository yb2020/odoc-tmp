// import { WebNoteAnnotationModel } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Web';
import { WebNoteAnnotationModel } from 'go-sea-proto/gen/ts/note/Web'
import {
  IDEAAnnotateType,
  RectOptions,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { NeedVipException } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { ToolBarType } from '@idea/pdf-annotate-core';

import { ERROR_CODE_NEED_VIP } from '@common/api/const';
import { useVipStore } from '@common/stores/vip';
import { ResponseError } from '~/src/api/type';
import { ElementName } from '~/src/api/report';
import { currentGroupId, currentNoteInfo } from '~/src/store';
import { rectStyleMap, styleMap } from '~/src/style/select';
import { NO_ANNOTATION_ID } from '~/src/constants';
import * as api from '~/src/api/annotations';
import { useAnnotationStore } from '.';
import { usePdfStore } from '../pdfStore';

export type ConvertAnnotations<OptionsType extends 'select' | 'rect'> = Omit<
  WebNoteAnnotationModel,
  'type' | 'select' | 'rect'
> &
  Omit<NonNullable<WebNoteAnnotationModel[OptionsType]>, 'rectangle'> &
  (typeof styleMap | typeof rectStyleMap)[
    | keyof typeof styleMap
    | keyof typeof rectStyleMap] & {
    type: ToolBarType;
    rectangles: RectOptions[];
    rectRaw?: boolean;
  };
export type AnnotationSelect = ConvertAnnotations<'select'>;
export type AnnotationRect = ConvertAnnotations<'rect'>;
export type AnnotationAll = AnnotationSelect | AnnotationRect;

export const convertAnnotations = (
  item: WebNoteAnnotationModel
): AnnotationAll => {
  const { type, select, rect, ...partialItem } = item;
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  const { rectangle, ...partialSelectRect } = (select || rect)!;
  let rectangles: RectOptions[] = [];
  if (rectangle) {
    rectangles = rectangles.concat(rectangle);
  }

  return {
    ...partialItem,
    ...partialSelectRect,
    ...(select
      ? styleMap[partialSelectRect.styleId as keyof typeof styleMap]
      : rectStyleMap[partialSelectRect.styleId as keyof typeof rectStyleMap]),
    type: type as unknown as ToolBarType,
    rectangles,
  };
};

export abstract class BaseAnnotationController {
  get pdfAnnotater() {
    const pdfStore = usePdfStore();

    return pdfStore.getAnnotater(currentNoteInfo.value?.noteId)!;
  }

  public abstract loadAnnotationMap(): Promise<void>;

  public async loadHotAnnotationMap() {
    useAnnotationStore().pageHotMap = {};
  }

  public abstract onlineSaveAnnotation(
    documentId: string,
    pageNumber: number,
    annotation: AnnotationAll,
    index?: number
  ): Promise<AnnotationAll['uuid'] | null>;

  public async addAnnotation(
    anno: Pick<AnnotationAll, 'documentId' | 'pageNumber'> &
      Partial<Omit<AnnotationAll, 'documentId' | 'pageNumber'>>
  ) {
    if ((anno as AnnotationSelect).rectStr === '') {
      this.localAddAnnotation(anno);
      return;
    }

    const result = this.pdfAnnotater.UI.commentController.addComment(
      anno as any
    );

    if (!result) {
      return;
    }

    const { svgGroupElementList, documentId, pageNumber, annotation } = result;

    const uuid =
      (await this.onlineSaveAnnotation(
        documentId,
        pageNumber,
        annotation as AnnotationAll
      )) ?? '';

    if (uuid) {
      svgGroupElementList?.forEach((group) => {
        this.pdfAnnotater.UI.commentController.updateUuid(
          group as SVGGElement,
          uuid
        );
      });

      this.localAddAnnotation(annotation as AnnotationAll);
    }

    return uuid;
  }

  public localAddAnnotation(
    annotation: Pick<AnnotationAll, 'documentId' | 'pageNumber'> &
      Partial<Omit<AnnotationAll, 'documentId' | 'pageNumber'>>,
    index?: number
  ) {
    const annotationStore = useAnnotationStore();
    if (!annotationStore.pageMap[annotation.pageNumber]) {
      annotationStore.pageMap[annotation.pageNumber] = [];
    }

    if (typeof index !== 'number') {
      index = annotationStore.pageMap[annotation.pageNumber].length;
    }

    annotationStore.pageMap[annotation.pageNumber].splice(
      index,
      0,
      annotation as AnnotationAll
    );
  }

  public localPatchAnnotation(
    annotationId: string,
    pageNumber: number,
    annotation: Partial<Omit<AnnotationAll, 'pageNumber'>>
  ) {
    const annotationStore = useAnnotationStore();
    const index = (annotationStore.pageMap[pageNumber] || []).findIndex(
      (item) => item.uuid === annotationId
    );

    if (index < 0) {
      return;
    }

    annotationStore.pageMap[pageNumber].splice(index, 1, {
      ...annotationStore.pageMap[pageNumber][index],
      ...annotation,
    });
  }

  public findAnnotation(annotateId: string) {
    const annotationStore = useAnnotationStore();
    let pageNumber = -1;
    let index = -1;
    let annotation: AnnotationAll | null = null;

    Object.keys(annotationStore.pageMap).some((pageNo) => {
      return annotationStore.pageMap[pageNo].some((anno, idx) => {
        if (anno.uuid !== annotateId) {
          return false;
        }

        annotation = anno;
        pageNumber = Number(pageNo);
        index = idx;
        return true;
      });
    });

    return {
      pageNumber,
      index,
      annotation: annotation as AnnotationAll | null,
    };
  }

  public async patchAnnotation(
    annotationId: string,
    annotation: Partial<Omit<AnnotationAll, 'pageNumber'>>
  ) {
    const result = this.findAnnotation(annotationId);
    if (!result.annotation) {
      return;
    }

    const page = useAnnotationStore().pageMap[result.pageNumber];
    const index = result.index;
    const annotationOrigin = result.annotation;
    const annotationChange = {
      ...annotationOrigin,
      ...annotation,
    };

    const patchDom = (anno: AnnotationAll) => {
      if ((anno as AnnotationSelect).rectStr === '') {
        return;
      }

      if (anno.type === ToolBarType.select) {
        this.pdfAnnotater.UI.commentController.setCommentToDom(
          annotationId,
          anno
        );
      } else {
        this.pdfAnnotater.UI.rectController.setRectToDom(annotationId, anno);
      }
    };

    const patchOnline = async () => {
      const { type, rectangles, isHighlight, ...rest } = annotationChange;

      const patch: any = {
        type: type as unknown as IDEAAnnotateType,
        pageNumber: rest.pageNumber,
        noteId: currentNoteInfo.value?.noteId ?? '',
        groupId: currentGroupId.value,
        isHighlight,
      };

      if (type === ToolBarType.select) {
        patch.select = {
          ...rest,
          rectangle: rectangles,
        };
      } else if (type === ToolBarType.rect) {
        const rect = {
          ...rest,
          rectangle: rectangles[0],
        };

        const { picUrl } = rect as Partial<AnnotationRect>;
        if (picUrl && !picUrl.startsWith('http')) {
          delete (rect as Partial<AnnotationRect>).picUrl;
        }
        patch.rect = rect;
      }

      await api.editAnnotation(patch);
    };

    try {
      page.splice(index, 1, annotationChange);
      await patchOnline();
      patchDom(annotationChange);
    } catch (error) {
      patchDom(annotationOrigin);
      const e = error as ResponseError;
      if (e.code === ERROR_CODE_NEED_VIP) {
        const vipStore = useVipStore();
        vipStore.showVipLimitDialog(e?.message, {
          exception: e?.extra as NeedVipException,
          reportParams: {
            element_name: ElementName.upperTeamNumNotePopup,
          },
        });
        return;
      }
      throw error;
    }
  }

  public deleteAnnotation(
    annotationId: string,
    pageNumber: number
  ): void | Promise<void> {
    const annotation = useAnnotationStore().pageMap[pageNumber].find(
      (item) => item.uuid === annotationId
    ) as AnnotationAll;

    const onlineDeleteAnnotation = async () => {
      const noteId = currentNoteInfo.value?.noteId ?? '';
      const groupId = currentGroupId.value;

      await api.deleteAnnotation({
        annotationId: annotation.uuid,
        noteId,
        groupId,
      });

      this.localDeleteAnnotation(annotation.uuid, annotation.pageNumber);
    };

    if ((annotation as AnnotationSelect).rectStr === '') {
      if (annotation.uuid === NO_ANNOTATION_ID) {
        this.localDeleteAnnotation(annotation.uuid, annotation.pageNumber);
        return;
      }

      return onlineDeleteAnnotation();
    }

    if (annotation.type === ToolBarType.rect) {
      this.pdfAnnotater.UI.rectController.deleteRectToDom(annotation.uuid);
    } else {
      this.pdfAnnotater.UI.commentController.deleteCommentToDom(
        annotation.uuid
      );
    }

    return onlineDeleteAnnotation();
  }

  public localDeleteAnnotation(annotationId: string, pageNumber: number) {
    const annotationStore = useAnnotationStore();
    const index = annotationStore.pageMap[pageNumber].findIndex(
      (item) => item.uuid === annotationId
    );

    if (index < 0) {
      return;
    }

    annotationStore.pageMap[pageNumber].splice(index, 1);

    if (annotationStore.pageMap[pageNumber].length === 0) {
      delete annotationStore.pageMap[pageNumber];
    }
  }
}
