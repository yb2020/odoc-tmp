import { IDEAAnnotateType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
// import { WebNoteAnnotationModel } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Web';
import { WebNoteAnnotationModel } from 'go-sea-proto/gen/ts/note/Web'
import { OperationType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/GroupNote';
import { ToolBarType } from '@idea/pdf-annotate-core';

import { currentGroupId, currentNoteInfo, store } from '~/src/store';

import {
  AnnotationAll,
  AnnotationRect,
  AnnotationSelect,
  BaseAnnotationController,
  convertAnnotations,
} from '~/src/stores/annotationStore/BaseAnnotationController';
import { ElementName, reportAddNote } from '~/src/api/report';
import * as api from '~/src/api/annotations';

import { useAnnotationStore } from '.';
import { getGroupNotes } from '~/src/api/groupNote';
import { NO_ANNOTATION_ID } from '~/src/constants';
import { arrow } from '~/src/dom/arrow';
import { styleMap } from '~/src/style/select';
import { ResponseError } from '~/src/api/type';
import { useVipStore } from '@common/stores/vip';
import { ERROR_CODE_NEED_VIP } from '@common/api/const';
import { NeedVipException } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';

const NOTE_MODIFIED_TIME_0 = '0';

export class GroupAnnotationController extends BaseAnnotationController {
  private noteModifiedTime = NOTE_MODIFIED_TIME_0;

  public async loadAnnotationMap() {
    await this.syncGroupAnnotation(NOTE_MODIFIED_TIME_0);
  }

  public clearGroupAnnotation() {
    this.noteModifiedTime = NOTE_MODIFIED_TIME_0;
    useAnnotationStore().pageMap = {};
  }

  public async syncGroupAnnotation(noteModifiedTime?: string) {
    const annotationStore = useAnnotationStore();

    if (!noteModifiedTime) {
      ({ noteModifiedTime } = this);
    }

    const { groupNoteModifiedTime, incrementalGroupNote } =
      await fetchGroupNotes(currentNoteInfo.value?.noteId, noteModifiedTime);

    if (noteModifiedTime === NOTE_MODIFIED_TIME_0) {
      const pageMap: typeof annotationStore.pageMap = {};
      incrementalGroupNote.forEach((item) => {
        item.webNoteAnnotationModels.forEach((anno) => {
          if (!pageMap[anno.pageNumber]) {
            pageMap[anno.pageNumber] = [];
          }

          pageMap[anno.pageNumber].push(anno);
        });
      });

      Object.keys(annotationStore.pageMap).forEach((pageNumber) => {
        annotationStore.pageMap[pageNumber]
          .filter((anno) => anno.uuid === NO_ANNOTATION_ID)
          .forEach((anno) => {
            pageMap[pageNumber].push(anno);
          });
      });

      annotationStore.pageMap = pageMap;
      this.noteModifiedTime = groupNoteModifiedTime;
      return;
    }

    incrementalGroupNote.forEach((item) => {
      if (item.operationType === OperationType.Create) {
        this.deltaCreate(item);
      }

      if (item.operationType === OperationType.Delete) {
        this.deltaDelete(item);
      }

      if (item.operationType === OperationType.Update) {
        this.deltaUpdate(item);
      }
    });

    this.noteModifiedTime = groupNoteModifiedTime;
  }

  private deltaCreate(delta: GroupNoteDelta) {
    const annotationStore = useAnnotationStore();
    delta.webNoteAnnotationModels.forEach((item: AnnotationAll) => {
      const index = (annotationStore.pageMap[item.pageNumber] || []).findIndex(
        (i) => i.uuid === item.uuid
      );

      console.log('note create', index, item);

      if ((item as AnnotationSelect).rectStr !== '') {
        const show =
          item.commentatorInfoView?.userId === store.state.user.userInfo?.id
            ? annotationStore.groupSelfVisible
            : annotationStore.groupOtherVisible;

        this.pdfAnnotater.UI.commentController.addCommentToDom(
          item as any,
          show
        );
      }

      if (index >= 0) {
        annotationStore.pageMap[item.pageNumber].splice(index, 1, item);
        return;
      }

      annotationStore.pageMap[item.pageNumber] = [
        ...(annotationStore.pageMap[item.pageNumber] || []),
        item,
      ];
    });
  }

  private deltaDelete(delta: GroupNoteDelta) {
    const annotationStore = useAnnotationStore();
    delta.webNoteAnnotationModels.forEach((item) => {
      annotationStore.pageMap[item.pageNumber] =
        annotationStore.pageMap[item.pageNumber] || [];

      const index = annotationStore.pageMap[item.pageNumber].findIndex(
        (a) => a.uuid === item.uuid
      );

      console.log('note delete', index, item);

      if (index >= 0) {
        if ((item as AnnotationSelect).rectStr !== '') {
          if (item.type === ToolBarType.select) {
            this.pdfAnnotater.UI.commentController.deleteCommentToDom(
              item.uuid
            );
          } else {
            this.pdfAnnotater.UI.rectController.deleteRectToDom(item.uuid);
          }
        }

        annotationStore.pageMap[item.pageNumber].splice(index, 1);
      }
    });
  }

  private deltaUpdate(delta: GroupNoteDelta) {
    const annotationStore = useAnnotationStore();
    delta.webNoteAnnotationModels.forEach((item) => {
      annotationStore.pageMap[item.pageNumber] =
        annotationStore.pageMap[item.pageNumber] || [];

      const index = annotationStore.pageMap[item.pageNumber].findIndex(
        (a) => a.uuid === item.uuid
      );

      console.log('note update', index, item);

      if (index === -1) {
        return;
      }

      if ((item as AnnotationSelect).rectStr !== '') {
        if (item.type === ToolBarType.select) {
          this.pdfAnnotater.UI.commentController.setCommentToDom(
            item.uuid,
            item
          );
        } else {
          this.pdfAnnotater.UI.rectController.setRectToDom(item.uuid, item);
        }

        if (arrow.annotationId === item.uuid) {
          arrow.line!.color =
            styleMap[item.styleId as keyof typeof styleMap].color;
        }
      }

      annotationStore.pageMap[item.pageNumber].splice(index, 1, item);
    });
  }

  public async onlineSaveAnnotation(
    documentId: string,
    pageNumber: number,
    annotation: AnnotationAll
  ) {
    const { type, rectangles, ...rest } = annotation;

    reportAddNote();

    const params = {
      type: type as unknown as IDEAAnnotateType,
      pageNumber,
      noteId: currentNoteInfo.value?.noteId,
      groupId: currentGroupId.value,
      pdfId: currentNoteInfo.value?.pdfId,
    } as Partial<WebNoteAnnotationModel> as WebNoteAnnotationModel;

    if (type === ToolBarType.select) {
      params.select = {
        ...(rest as AnnotationSelect),
        rectangle: rectangles,
      };
    } else if (type === ToolBarType.rect) {
      params.rect = {
        ...(rest as AnnotationRect),
        rectangle: rectangles[0],
      };
    }

    let id = '';
    try {
      id = await api.addGroupAnnotation(params);
    } catch (err) {
      const e = err as ResponseError;
      if (e.code === ERROR_CODE_NEED_VIP) {
        const vipStore = useVipStore();
        vipStore.showVipLimitDialog(e?.message, {
          exception: e?.extra as NeedVipException,
          reportParams: {
            element_name: ElementName.upperTeamNumNotePopup,
          },
        });
      }
      throw err;
    }

    annotation.uuid = id;
    annotation.pageNumber = pageNumber;
    annotation.documentId = documentId;

    return id;
  }
}

export interface GroupNoteDelta {
  operationType: OperationType;
  webNoteAnnotationModels: Array<ReturnType<typeof convertAnnotations>>;
}

export interface GroupNoteResponse {
  incrementalGroupNote: GroupNoteDelta[];
  groupNoteModifiedTime: string;
}

export const fetchGroupNotes = async (
  documentId: string,
  noteModifiedTime = NOTE_MODIFIED_TIME_0
) => {
  const res = await getGroupNotes({
    noteId: documentId,
    noteModifiedTime,
    operationTypes:
      noteModifiedTime === NOTE_MODIFIED_TIME_0
        ? [OperationType[2]]
        : [OperationType[0], OperationType[1], OperationType[2]],
  } as any);

  res.incrementalGroupNote.forEach((item) => {
    item.webNoteAnnotationModels = item.webNoteAnnotationModels.map((item) =>
      convertAnnotations(item)
    ) as any;
  });

  return res as unknown as GroupNoteResponse;
};
