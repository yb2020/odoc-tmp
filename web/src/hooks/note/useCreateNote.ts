import { PageSelectText, ViewerController } from '@idea/pdf-annotate-viewer';
import { WordInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/NoteModule';
import { useCommentGlobalState } from '../useNoteState';
import { styleMap } from '@/style/select';
import { ToolBarType, scaleDownRaw } from '@idea/pdf-annotate-core';
import { SelectAnnotation } from '@idea/pdf-annotate-core/render/renderSelect';
import { Nullable } from '~/src/typings/global';
import { UniTranslateResp } from '~/src/api/translate';
import { useStore } from '~/src/store';
import { PageTextRects } from '~/src/components/ToolTip/type';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { AnnotationSelect } from '~/src/stores/annotationStore/BaseAnnotationController';
import { nextTick, ref } from 'vue';
import { useWordNotes } from '@common/components/Notes/useWordNote';

const useCreateNote = ({
  pdfId,
  noteId,
  pdfViewer,
}: {
  pdfId: string;
  noteId: string;
  pdfViewer?: ViewerController;
}) => {
  const annotationStore = useAnnotationStore();
  const commentState = useCommentGlobalState();
  const store = useStore();
  const { add: addWordNote, gcolorHex } = useWordNotes(ref(pdfId), ref(noteId));

  const getRectInfo = ({
    info: pageTexts,
  }: {
    info: Nullable<PageSelectText[]>;
  }) => {
    if (!pageTexts) {
      return null;
    }

    const pageNumber = pageTexts[0].pageNum;

    const rects = pageTexts.reduce((prev: PageTextRects, current) => {
      const { rects } = current;
      const _rects = rects.map((item) => {
        return {
          ...item,
          pageNumber: current.pageNum,
        };
      });

      return [...prev, ..._rects];
    }, []);

    const rectStr = pageTexts.map((text) => text.text).join('') || '';

    return {
      pageNumber,
      rects,
      rectStr,
      rectRaw: false,
    };
  };

  const add = async (
    noteInfo:
      | {
          info: Nullable<PageSelectText[]>;
        }
      | {
          rectInfo: ReturnType<typeof getRectInfo>;
        },
    idea: string
  ) => {
    const rectInfo =
      'rectInfo' in noteInfo ? noteInfo.rectInfo : getRectInfo(noteInfo);
    if (!rectInfo) {
      return;
    }
    const { pageNumber, rects, rectStr, rectRaw } = rectInfo;

    const styleId = commentState.value.styleId;
    pdfViewer?.getDocumentViewer().clearSelection();
    const uuid = (await annotationStore.controller.addAnnotation({
      pageNumber,
      ...({
        rectStr,
      } as Partial<AnnotationSelect>),
      rectRaw,
      rectangles: rects as any[],
      type: ToolBarType.select,
      documentId: noteId,
      idea,
      styleId: +styleId,
      ...styleMap[styleId],
      commentatorInfoView: {
        avatarCdnUrl: store.state.user.userInfo?.avatarUrl ?? '',
        nickName: store.state.user.userInfo?.nickName ?? '',
        userId: store.state.user.userInfo?.id ?? '',
      },
    })) as string;

    annotationStore.showHoverNote({
      uuid,
      page: pageNumber,
    });
    await nextTick();
    annotationStore.currentAnnotationId = uuid;
  };

  const addWord = async (
    noteInfo: { info: Nullable<PageSelectText[]> },
    idea: string,
    translation: UniTranslateResp
  ) => {
    const wordInfo: WordInfo = {
      word: idea,
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      rectangle: [],
      translateInfo: translation
        ? {
            ...translation,
            targetResp: translation.targetResp ?? [],
            targetContent: [translation.targetContent],
          }
        : undefined,
    };

    const rectInfo = getRectInfo(noteInfo);
    if (!rectInfo) {
      await addWordNote({ wordInfo });
      return;
    }

    const { pageNumber, rects, rectStr } = rectInfo;
    const pdfV = pdfViewer?.getDocumentViewer().getPdfViewer();
    const pageViewport = pdfV?.getPageView(pageNumber - 1)?.viewport;
    const annotation: SelectAnnotation = {
      type: ToolBarType.Vocabulary,
      rectangles: rects.map((rect) => {
        const { x, y, width, height } = rect;
        return {
          ...rect,
          ...scaleDownRaw(pageViewport?.scale || 1, {
            x,
            y,
            width,
            height,
          }),
        };
      }),
      rectStr,
      idea,
      color: gcolorHex.value,
      pageNumber,
      tags: [],
    };

    wordInfo.rectangle =
      (annotation?.rectangles as WordInfo['rectangle']) ?? rects;
    await addWordNote({ wordInfo });
  };

  return {
    add,
    addWord,
  };
};

export default useCreateNote;
