import _ from 'lodash';
import { computed, watch } from 'vue';
import { defineStore } from 'pinia';
import {
  ToolBarType,
  PageShapeTextMap,
  PageHandwriteMap,
  PageShapeMap,
} from '@idea/pdf-annotate-core';

import {
  currentGroupId,
  currentNoteInfo,
  selfNoteInfo,
  store,
} from '~/src/store';
import { AnnotationAll } from './BaseAnnotationController';
import { GroupAnnotationController } from './GroupAnnotationController';
import { PersonAnnotationController } from './PersonAnnotationController';
import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type';
import {
  // AnnotateTag,
  ShapeType,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import {AnnotateTag} from 'go-sea-proto/gen/ts/common/AnnotateTag'
import {
  PDF_IS_SHOW_HOT_SELECT,
  PDF_IS_SHOW_MY_GROUP_SELECT,
  PDF_IS_SHOW_MY_SELECT,
  PDF_IS_SHOW_OTHER_GROUP_SELECT,
  PDF_IS_SHOW_IMAGE_GROUP_SELECT,
  PDF_SELECT_HIDE,
  PDF_SELECT_SHOW,
} from '~/src/constants';
import { ColorKey, NoteFilter } from '~/src/style/select';
import { getUserTags } from '~/src/api/base';
import { NoteSlice } from './types';
import { shapeState } from '~/src/components/Main/mouseCore';
import { usePdfStore } from '../pdfStore';

export type PageAnnotationMap = Record<string | number, AnnotationAll[]>;

export const noteBuffer = {
  handwriteBuffer: null as null | Promise<PageHandwriteMap>,
  annotationBuffer: null as null | Promise<
    [PageAnnotationMap, PageShapeTextMap, number, number]
  >,
  shapeBuffer: null as null | Promise<PageShapeMap>,
};

export const CLASS_NAME_ANNOTATION_SCREENSHOT_CAPTURING =
  'annotation-screenshot-capturing';

export interface AnnotationState {
  instantiated: boolean;
  pageMap: PageAnnotationMap;
  pageHotMap: PageAnnotationMap;
  currentAnnotationId: string;
  inputingAnnotationId: string;
  personVisible: boolean;
  hotVisible: boolean;
  hotTipVisible: boolean;
  groupSelfVisible: boolean;
  groupOtherVisible: boolean;
  groupImageVisible: boolean;
  activeOverlayAnnotateId: string;
  activeOverlayPageNumber: number;
  activeColorMap: Record<NoteFilter, boolean>;
  mapHideHoverNotes: Record<string, ReturnType<typeof setTimeout>>;
  activeHoverNotes: NoteSlice[];
  tagList: AnnotateTag[];
  showReferenceTippy: boolean;
}

export const personAnnotationController = new PersonAnnotationController();
export const groupAnnotationController = new GroupAnnotationController();

export const useAnnotationStore = defineStore('annotation', {
  state: (): AnnotationState => {
    return {
      instantiated: false,
      pageMap: {},
      pageHotMap: {},
      currentAnnotationId: '',
      inputingAnnotationId: '',
      personVisible: getLocalVisible(PDF_IS_SHOW_MY_SELECT),
      hotVisible: getLocalVisible(PDF_IS_SHOW_HOT_SELECT, false),
      hotTipVisible: false,
      groupSelfVisible: getLocalVisible(PDF_IS_SHOW_MY_GROUP_SELECT),
      groupOtherVisible: getLocalVisible(PDF_IS_SHOW_OTHER_GROUP_SELECT),
      groupImageVisible: getLocalVisible(PDF_IS_SHOW_IMAGE_GROUP_SELECT, false),
      activeOverlayAnnotateId: '',
      activeOverlayPageNumber: NaN,
      activeColorMap: {
        [ColorKey.blue]: true,
        [ColorKey.green]: true,
        [ColorKey.yellow]: true,
        [ColorKey.orange]: true,
        [ColorKey.red]: true,
        ref: true,
      },
      mapHideHoverNotes: {},
      activeHoverNotes: [],
      tagList: [],
      showReferenceTippy: false,
    };
  },
  getters: {
    controller() {
      return store.state.base.currentGroupId === SELF_NOTEINFO_GROUPID
        ? personAnnotationController
        : groupAnnotationController;
    },
    crossPageMap() {
      const cross: PageAnnotationMap = {};

      Object.values(this.pageMap).forEach((list) => {
        list.forEach((annotation) => {
          annotation.rectangles.forEach((rect) => {
            const page = rect.pageNumber || annotation.pageNumber;

            if (!(page in cross)) {
              cross[page] = [];
            }

            let index = cross[page].findIndex(
              (anno) => anno.uuid === annotation.uuid
            );
            if (index === -1) {
              cross[page].push({
                ...annotation,
                pageNumber: page,
                rectangles: [],
              });
              index += cross[page].length;
            }

            cross[page][index].rectangles.push(rect);
          });
        });
      });

      return cross;
    },
    headTailPageNumber() {
      const pageNumberList = Object.keys(this.pageMap)
        .filter((key) => this.pageMap[key])
        .map(Number);

      const head = Math.min(...pageNumberList);
      const tail = Math.max(...pageNumberList);
      return {
        head,
        tail,
      };
    },
    count() {
      const c: number = Object.values(this.pageMap).reduce(
        (sum, list) =>
          sum +
          list.reduce((num, item) => {
            num += item.isHighlight ? 0 : 1;
            return num;
          }, 0),
        0
      );
      return c;
    },
  },
  actions: {
    async refreshTagList() {
      this.tagList = await getUserTags();
    },
    async showHoverNote(y: NoteSlice, isReset = false) {
      const n = this.activeHoverNotes.find((x) => x.uuid === y.uuid);

      if (isReset) {
        this.activeHoverNotes = [y];
      } else if (!n) {
        this.activeHoverNotes.push(y);
      } else {
        clearTimeout(this.mapHideHoverNotes[y.uuid]);
      }
    },
    async mutHoverNote(uuid: string, rest: Partial<NoteSlice>) {
      const n = this.activeHoverNotes.find((x) => x.uuid === uuid);

      if (n) {
        Object.assign(n, rest);
      }
    },
    async delHoverNote(uuid: string, delay = 100) {
      const n = this.activeHoverNotes.find((x) => x.uuid === uuid);
      if (!n || n.locked) {
        return;
      }
      const action = () =>
        _.remove(this.activeHoverNotes, (x) => x.uuid === uuid);

      if (delay > 0 && n) {
        const timer = setTimeout(action, delay);
        this.mapHideHoverNotes[uuid] = timer;
      } else {
        action();
      }
    },
  },
});

function getLocalVisible(namespace: string, defaultVisible = true) {
  const localValue = localStorage.getItem(namespace);
  const visible =
    (localValue === null ? defaultVisible : false) ||
    localValue === PDF_SELECT_SHOW;
  return visible;
}

let syncAnnotationVisibleInited = false;
export const useSyncAnnotationVisible = () => {
  if (syncAnnotationVisibleInited) {
    return;
  }

  syncAnnotationVisibleInited = true;

  const annotationStore = useAnnotationStore();
  const pdfStore = usePdfStore();
  const pdfAnnotater = computed(() => {
    return pdfStore.getAnnotater(selfNoteInfo.value.noteId);
  });
  const pdfAnnotaterGroup = computed(() => {
    if (currentGroupId.value !== SELF_NOTEINFO_GROUPID) {
      return pdfStore.getAnnotater(currentNoteInfo.value?.noteId);
    }
  });

  useSyncVisibleToLocal();
  useSyncLocalToVisible();
  useSyncSvgPerson();
  useSyncSvgGroup();

  function useSyncVisibleToLocal() {
    useWatchLocalVisible(
      () => annotationStore.personVisible,
      PDF_IS_SHOW_MY_SELECT
    );
    useWatchLocalVisible(
      () => annotationStore.hotVisible,
      PDF_IS_SHOW_HOT_SELECT
    );
    useWatchLocalVisible(
      () => annotationStore.groupSelfVisible,
      PDF_IS_SHOW_MY_GROUP_SELECT
    );
    useWatchLocalVisible(
      () => annotationStore.groupOtherVisible,
      PDF_IS_SHOW_OTHER_GROUP_SELECT
    );
    useWatchLocalVisible(
      () => annotationStore.groupImageVisible,
      PDF_IS_SHOW_IMAGE_GROUP_SELECT
    );

    function useWatchLocalVisible(getValue: () => boolean, namespace: string) {
      watch(getValue, (value) => {
        console.warn({ namespace, value });

        const flag = value ? PDF_SELECT_SHOW : PDF_SELECT_HIDE;

        if (localStorage.getItem(namespace) !== flag) {
          localStorage.setItem(namespace, flag);
        }
      });
    }
  }

  function useSyncLocalToVisible() {
    setInterval(() => {
      const person = getLocalVisible(PDF_IS_SHOW_MY_SELECT);
      if (person !== annotationStore.personVisible) {
        annotationStore.personVisible = person;
      }

      const hot = getLocalVisible(PDF_IS_SHOW_HOT_SELECT, false);
      if (hot !== annotationStore.hotVisible) {
        annotationStore.hotVisible = hot;
      }

      const groupSelf = getLocalVisible(PDF_IS_SHOW_MY_GROUP_SELECT);
      if (groupSelf !== annotationStore.groupSelfVisible) {
        annotationStore.groupSelfVisible = groupSelf;
      }

      const groupOther = getLocalVisible(PDF_IS_SHOW_OTHER_GROUP_SELECT);
      if (groupOther !== annotationStore.groupOtherVisible) {
        annotationStore.groupOtherVisible = groupOther;
      }

      const groupImage = getLocalVisible(PDF_IS_SHOW_IMAGE_GROUP_SELECT);
      if (groupImage !== annotationStore.groupImageVisible) {
        annotationStore.groupImageVisible = groupImage;
      }
    }, 3000);
  }

  function useSyncSvgPerson() {
    watch(() => annotationStore.hotVisible, syncSvg);
    watch(() => annotationStore.personVisible, syncSvg);

    function syncSvg() {
      if (pdfAnnotater.value?.documentId === currentNoteInfo.value?.noteId) {
        pdfAnnotater.value?.setDisplayHandwrite(annotationStore.personVisible);
        pdfAnnotater.value?.setDisplayByType(
          String(ToolBarType.hot),
          annotationStore.hotVisible,
          annotationStore.personVisible
        );
      }
    }

    watch(currentGroupId, (id) => {
      const isSelfNoteTab = id === SELF_NOTEINFO_GROUPID;
      if (pdfAnnotater.value?.documentId === currentNoteInfo.value?.noteId) {
        pdfAnnotater.value?.setDisplayHandwrite(isSelfNoteTab);
      }

      if (!isSelfNoteTab) {
        shapeState.creating.value = ShapeType.UNRECOGNIZED;
      }
    });
  }

  function useSyncSvgGroup() {
    watch(() => annotationStore.groupSelfVisible, syncSvg);
    watch(() => annotationStore.groupOtherVisible, syncSvg);

    function syncSvg() {
      if (
        pdfAnnotaterGroup.value?.documentId === currentNoteInfo.value?.noteId
      ) {
        pdfAnnotaterGroup.value.setDisplayByUserId(
          store.state.user.userInfo?.id ?? '',
          annotationStore.groupSelfVisible,
          annotationStore.groupOtherVisible
        );
      }
    }
  }
};
