import { message } from 'ant-design-vue';
import {
  IDEAAnnotateType,
//   // AnnotateTag,
//   // ShapeAnnotation,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { ShapeAnnotation } from 'go-sea-proto/gen/ts/common/ShapeAnnotation'
import { AnnotateTag } from 'go-sea-proto/gen/ts/common/AnnotateTag'
// import type { IDEAAnnotateType } from 'go-sea-proto/gen/ts/common/IDEAAnnotateType'

import {
  OperationTypeOfSort,
  // WebNoteAnnotationModel,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Web';
import { WebNoteAnnotationModel } from 'go-sea-proto/gen/ts/note/Web'
import {
  ToolBarType,
  PageHandwriteMap,
  PageShapeTextMap,
  PageShapeMap,
  TextAnnotation,
} from '@idea/pdf-annotate-core';
import { selfNoteInfo } from '@/store';
import { reportAddNote } from '~/src/api/report';

import * as api from '~/src/api/annotations';
import {
  useAnnotationStore,
  PageAnnotationMap,
  noteBuffer,
} from '~/src/stores/annotationStore';
import {
  AnnotationAll,
  AnnotationRect,
  AnnotationSelect,
  BaseAnnotationController,
  convertAnnotations,
} from '~/src/stores/annotationStore/BaseAnnotationController';
import { createElementVNode } from 'vue';
import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type';
import i18n from '~/src/locals/i18n';
import { ResponseError } from '~/src/api/type';
import { ERROR_CODE_UNCERT } from '~/src/store/cert';

export class PersonAnnotationController extends BaseAnnotationController {
  public async loadAnnotationMap() {
    if (!noteBuffer.annotationBuffer) {
      this.refreshAnnotationBuffer();
    }

    if (!noteBuffer.handwriteBuffer) {
      this.refreshHandwriteBuffer();
    }

    if (!noteBuffer.shapeBuffer) {
      this.refreshShapeBuffer();
    }

    [useAnnotationStore().pageMap] = (await noteBuffer.annotationBuffer) as [
      PageAnnotationMap,
      PageShapeTextMap,
      number,
      number,
    ];
  }

  private get annotationParams() {
    return {
      noteId: selfNoteInfo.value?.noteId ?? '',
      groupId: SELF_NOTEINFO_GROUPID,
    };
  }

  private refreshAnnotationBuffer() {
    noteBuffer.annotationBuffer = (async () => {
      try {
        // eslint-disable-next-line no-var
        var { annotationList, handTextList } = await api.getByNote(
          this.annotationParams
        );
      } catch (error) {
        noteBuffer.annotationBuffer = null;
        message.error('获取笔记数据失败！');
        throw error;
      }

      const pageMap = this.createPageMap(annotationList);
      const PageShapeTextMap = createMap(handTextList);

      return [
        pageMap,
        PageShapeTextMap,
        annotationList.length,
        handTextList.length,
      ];
    })();
  }

  private refreshHandwriteBuffer() {
    noteBuffer.handwriteBuffer = (async (): Promise<PageHandwriteMap> => {
      try {
        // eslint-disable-next-line no-var
        var handwriteList = await api.getHandwrite(this.annotationParams);
      } catch (error) {
        noteBuffer.handwriteBuffer = null;
        message.error('获取手写数据失败');
        return {};
      }

      const pageHandwriteMap = createMap(handwriteList);
      return pageHandwriteMap;
    })();
  }

  private refreshShapeBuffer() {
    noteBuffer.shapeBuffer = (async () => {
      try {
        // eslint-disable-next-line no-var
        var shapeList = await api.getShapeList({
          noteId: this.annotationParams.noteId,
        });
      } catch (error) {
        noteBuffer.shapeBuffer = null;
        message.error('获取图形数据失败');
        return {};
      }

      const pageShapeMap: PageShapeMap = {};
      shapeList.forEach((shape) => {
        if (!pageShapeMap[shape.pageNumber]) {
          pageShapeMap[shape.pageNumber] = [];
        }

        pageShapeMap[shape.pageNumber].push(shape);
      });

      return pageShapeMap;
    })();
  }

  private createPageMap(response: WebNoteAnnotationModel[]) {
    const pageMap = createMap(response.filter(Boolean).map(convertAnnotations));

    return pageMap;
  }

  public async loadHotAnnotationMap() {
    useAnnotationStore().pageHotMap = await this.getHotAnnotationMap();
  }

  private latestHotAnnotationPdfId?: string;
  private latestHotAnnotationMap: Promise<PageAnnotationMap> | null = null;
  private async getHotAnnotationMap() {
    const pdfId = selfNoteInfo.value?.pdfId ?? '';
    if (
      pdfId === this.latestHotAnnotationPdfId &&
      this.latestHotAnnotationMap
    ) {
      return this.latestHotAnnotationMap;
    }

    this.latestHotAnnotationPdfId = pdfId;
    this.latestHotAnnotationMap = (async () => {
      try {
        // eslint-disable-next-line no-var
        var response = await api.getHotAnnotations({
          pdfId,
        });
      } catch (error) {
        message.error('获取数据失败！');
        this.latestHotAnnotationPdfId = undefined;
        this.latestHotAnnotationMap = null;
        throw error;
      }

      return this.createPageMap(response);
    })();

    return this.latestHotAnnotationMap;
  }

  public async onlineSaveAnnotation(
    documentId: string,
    pageNumber: number,
    annotation: AnnotationAll,
    index?: number
  ) {
    const { type, rectangles, ...rest } = annotation;

    reportAddNote();

    const params = {
      type: type as unknown as IDEAAnnotateType,
      pageNumber,
      position: typeof index !== 'number' || index < 0 ? -1 : index,
      pdfId: selfNoteInfo.value.pdfId,
      noteId: selfNoteInfo.value.noteId,
    } as Partial<WebNoteAnnotationModel> as WebNoteAnnotationModel;

    params.markIdsOfCurrentPage =
      params.position === -1
        ? []
        : useAnnotationStore().pageMap[pageNumber]?.map((item) => item.uuid) ??
          [];

    if (type === ToolBarType.select) {
      params.isHighlight = annotation.isHighlight;
      params.select = {
        ...(rest as AnnotationSelect),
        rectangle: rectangles,
      };
    } else if (type === ToolBarType.rect) {
      params.rect = {
        ...(rest as AnnotationRect),
        rectangle: rectangles[0],
      };
      params.groupId = SELF_NOTEINFO_GROUPID;
      (params as any).noteId = selfNoteInfo.value?.noteId;
    } else if (params.type === IDEAAnnotateType.IDEAAnnotateTypeTextBox) {
      params.textBox = annotation.textBox;
    }

    try {
      // eslint-disable-next-line no-var
      var id = await api.postSaveAnnotation(params);
    } catch (error) {
      const responseError = error as ResponseError;
      if (responseError.code === ERROR_CODE_UNCERT) {
        message.error(responseError.message);
      } else if (
        responseError instanceof Object &&
        responseError.code === api.ERROR_CODE_SORT
      ) {
        console.error('调整顺序失败', error);
        this.needRefresh();
      } else {
        message.error(
          error instanceof Object ? (error as any).message || '' : error
        );
      }

      return null;
    }

    annotation.uuid = id;
    annotation.pageNumber = pageNumber;
    annotation.documentId = documentId;

    this.refreshAnnotationBuffer();

    return id;
  }

  public async patchAnnotation(
    annotationId: string,
    annotation: Partial<Omit<AnnotationAll, 'pageNumber'>>
  ) {
    await super.patchAnnotation(annotationId, annotation);
    this.refreshAnnotationBuffer();
  }

  public deleteAnnotation(
    annotationId: string,
    pageNumber: number
  ): void | Promise<void> {
    const result = super.deleteAnnotation(annotationId, pageNumber);
    if (!result) {
      return;
    }

    return (async () => {
      await result;
      this.refreshAnnotationBuffer();
    })();
  }

  public async sortAnnotation(
    annotationId: string,
    toPageNumber: number,
    notesSort: Record<number, string[]>
  ) {
    const annotationStore = useAnnotationStore();

    console.warn('sort_notes', notesSort);

    const restore: Record<number, AnnotationAll[]> = {};
    const update: Record<number, AnnotationAll[]> = {};
    const annotationList: AnnotationAll[] = [];
    const pageNumberList = Object.keys(notesSort);
    pageNumberList.forEach((pageNumber) => {
      restore[Number(pageNumber)] = annotationStore.pageMap[pageNumber];
      annotationList.push(...annotationStore.pageMap[pageNumber]);
    });
    pageNumberList.forEach((pageNumber) => {
      update[Number(pageNumber)] = notesSort[Number(pageNumber)].map((id) => {
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        const annotation = annotationList.find((item) => item.uuid === id)!;
        if (annotation.pageNumber === Number(pageNumber)) {
          return annotation;
        }

        return {
          ...annotation,
          pageNumber: Number(pageNumber),
        };
      });
    });

    const request = api.putSortAnnotations({
      noteId: selfNoteInfo.value?.noteId ?? '',
      pageNumber: toPageNumber,
      markIdsWithSort: notesSort[toPageNumber].map((markId) => ({
        markId,
        operationTypeOfSort:
          annotationId === markId
            ? OperationTypeOfSort.addToPage
            : OperationTypeOfSort.exchangeInPage,
      })),
    });

    localSortAnnotation(update);
    try {
      await request;
    } catch (error) {
      const responseError = error as ResponseError;
      if (
        responseError instanceof Object &&
        responseError.code === ERROR_CODE_UNCERT
      ) {
        message.error(responseError.message);
      } else if (
        responseError instanceof Object &&
        responseError.code === api.ERROR_CODE_SORT
      ) {
        console.error('调整顺序失败', error);
        this.needRefresh();
      } else {
        message.error(responseError.message);
      }

      localSortAnnotation(restore);
    }

    function localSortAnnotation(
      partialPageMap: Record<number, AnnotationAll[]>
    ) {
      Object.keys(partialPageMap).forEach((pageNumber) => {
        annotationStore.pageMap[pageNumber] =
          partialPageMap[Number(pageNumber)];
      });
    }
  }

  private needRefresh = () => {
    const a = createElementVNode(
      'a',
      {
        attrs: {
          href: 'javascript:void(0)',
        },
      },
      `[${i18n.global.t('viewer.refresh')}]`
    );
    const text = i18n.global
      .t('message.inValidDataNeedToRefreshTip')
      .split('%%');
    text.splice(1, 0, a as any);
    const span = createElementVNode('span', {}, text);
    message.error({
      content: span,
      onClick: () => {
        this.refreshAnnotationBuffer();
        this.loadAnnotationMap();
      },
    });
  };

  public addTag(annotationId: string, pageNumber: number, tag: AnnotateTag) {
    const annotation = useAnnotationStore().pageMap[pageNumber].find(
      (item) => item.uuid === annotationId
    );
    annotation?.tags.push({ ...tag });
  }

  public removeTag(
    annotationId: string,
    pageNumber: number,
    tagId: AnnotateTag['tagId']
  ) {
    const annotation = useAnnotationStore().pageMap[pageNumber].find(
      (item) => item.uuid === annotationId
    );
    if (!annotation) {
      return;
    }

    const index = annotation.tags.findIndex((item) => item.tagId === tagId);
    if (index !== -1) {
      annotation.tags.splice(index, 1);
    }
  }

  public async createShape(shape: ShapeAnnotation) {
    let uuid: string;
    try {
      uuid = await api.saveShape({
        // pdfId: selfNoteInfo.value.pdfId ?? '',
        noteId: this.annotationParams.noteId,
        shapeAnnotation: {
          ...shape,
          uuid: '',
        },
      });
    } catch (error) {
      return;
    }

    this.refreshShapeBuffer();
    return uuid;
  }

  public async updateShape(shape: ShapeAnnotation) {
    try {
      await api.updateShape({
        annotations: [shape],
      });
    } catch (error) {
      return;
    }

    this.refreshShapeBuffer();
  }

  public async deleteShape(shapeId: ShapeAnnotation['uuid']) {
    try {
      await api.deleteShape({
        shapeIds: [shapeId],
      });
    } catch (error) {
      return;
    }

    this.refreshShapeBuffer();
  }

  public async deleteTextBox(annotationId: string) {
    await api.deleteAnnotation({
      annotationId,
      noteId: this.annotationParams.noteId,
      groupId: this.annotationParams.groupId,
    });

    const [, textMap] = await noteBuffer.annotationBuffer!;

    Object.values(textMap).some((list) => {
      const index = list.findIndex((item) => item.textBox?.id === annotationId);

      if (index !== -1) {
        list.splice(index, 1);
        return true;
      }
    });
  }

  public async editTextBox(textItem: TextAnnotation) {
    await api.editAnnotation({
      ...textItem,
      pdfId: selfNoteInfo.value.pdfId,
      noteId: this.annotationParams.noteId,
      groupId: this.annotationParams.groupId,
    } as any);

    this.refreshAnnotationBuffer();
  }
}

function createMap<Item extends { pageNumber: number }>(list: Item[]) {
  const map: Record<string | number, Item[]> = {};
  list.forEach((item) => {
    if (!map[item.pageNumber]) {
      map[item.pageNumber] = [];
    }

    map[item.pageNumber].push(item);
  });

  return map;
}
