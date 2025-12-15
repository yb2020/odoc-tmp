import { ref, nextTick } from 'vue';
import { ToolBarType } from '@idea/pdf-annotate-core';
import { delay } from '@idea/aiknowledge-special-util/delay';
import {
  AnnotationAll,
  AnnotationSelect,
} from '~/src/stores/annotationStore/BaseAnnotationController';
import { currentNoteInfo, store } from '~/src/store';
import { NO_ANNOTATION_ID } from '~/src/constants';
import { styleMap } from '~/src/style/select';
import { RightSideBarType } from '../type';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { ViewerController } from '@idea/pdf-annotate-viewer';

export const annotationDragging = ref<AnnotationAll | null>(null);
export const annotationDragHovering = ref<AnnotationAll['uuid']>('');
export const annotationDragPosition = ref(NaN);
export const annotationDragStyled = ref(false);
export const annotationDragFetching = ref(false);

export type NewAnnotation = Omit<
  AnnotationSelect,
  | 'id'
  | 'createDate'
  | 'paperId'
  | 'paperTitle'
  | 'docName'
  | 'groupId'
  | 'deleteAuthority'
  | 'position'
  | 'score'
  | 'rectangles'
  | 'markIdsOfCurrentPage'
>;

const getDefaultInsertIndex = (pageNumber: number) =>
  (useAnnotationStore().pageMap[pageNumber] || []).length;

export const getNewAnnotation = (
  pageNo?: number,
  pdfViewer?: ViewerController
): NewAnnotation => {
  const pageNumber =
    typeof pageNo === 'number'
      ? pageNo
      : (pdfViewer?.getDocumentViewer().getPdfViewer()
          .currentPageNumber as number) ?? 1;

  return {
    pdfId: currentNoteInfo.value?.pdfId ?? '',
    pageNumber,
    rectStr: '',
    type: ToolBarType.select,
    documentId: currentNoteInfo.value?.noteId ?? '',
    idea: '',
    styleId: 1,
    ...styleMap[1],
    uuid: NO_ANNOTATION_ID,
    commentatorInfoView: {
      avatarCdnUrl: store.state.user.userInfo?.avatarUrl ?? '',
      nickName: store.state.user.userInfo?.nickName ?? '',
      userId: store.state.user.userInfo?.id ?? '',
    },
    tags: [],
  };
};

export const insertNoReferenceAnnotation = async (
  pageNumber?: number,
  index?: number
) => {
  const annotation = getNewAnnotation(pageNumber);

  if (typeof index !== 'number') {
    index = getDefaultInsertIndex(annotation.pageNumber);
  }

  const annotationStore = useAnnotationStore();
  const uuid = await annotationStore.controller.onlineSaveAnnotation(
    annotation.documentId,
    annotation.pageNumber,
    annotation as any,
    index
  );
  await annotationStore.controller.loadAnnotationMap();
  await nextTick();
  await delay(100);
  console.warn('展开编辑', uuid);
  annotationStore.currentAnnotationId = uuid as string;
};

export const safeInsertNoReferenceAnnotation = (
  event: Event,
  tab: RightSideBarType,
  pageNumber?: number,
  index?: number
) => {
  if (tab !== RightSideBarType.Group) {
    event.stopPropagation();
    return insertNoReferenceAnnotation(pageNumber, index);
  }
};

export const handleEmptyClick = async (
  tab: RightSideBarType,
  pdfViewer?: ViewerController
) => {
  if (tab !== RightSideBarType.Group) {
    return;
  }

  const params = getNewAnnotation(undefined, pdfViewer);

  const annotationStore = useAnnotationStore();
  annotationStore.controller.localAddAnnotation(params);
  await nextTick();
  await delay(100);
  annotationStore.currentAnnotationId = params.uuid;
};

export const checkNoReferenceAnnotation = (note: AnnotationAll) => {
  return (
    note.type === ToolBarType.select &&
    (note as AnnotationSelect).rectStr === ''
  );
};
