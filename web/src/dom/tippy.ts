import {
  PdfFigureAndTableInfo,
  // ReferenceMarker,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse';
import { ReferenceMarker } from 'go-sea-proto/gen/ts/pdf/PdfParse';
import { App, createApp, nextTick, ref } from 'vue';
import Antd from 'ant-design-vue';
import PerfectScrollbar from 'vue3-perfect-scrollbar';
import ReferenceVue from '../components/Tippy/Reference/index.vue';
import FigureVue from '../components/Tippy/Figure.vue';
import TranslateVue from '../components/Tippy/Translate/Modal.vue';
import tippy, { Instance, Props } from 'tippy.js';
import enableTippyDraggable from './enableTippyDraggable';
import enableTippyResizable from './enableTippyResizable';
import interact from 'interactjs';
import { Nullable } from '../typings/global';
import merge from 'lodash-es/merge';
import { TranslateState, useTranslateStore } from '../stores/translateStore';
import useLocalTippy from '~/src/hooks/useLocalTippy';
import useCreateNote from '../hooks/note/useCreateNote';
import { selfNoteInfo } from '../store';
import { PageSelectText } from '~/../../packages/pdf-annotate-viewer/typing';
import { usePinaStore } from '../stores';
import i18n from '../locals/i18n';
import { UniTranslateResp } from '../api/translate';
import { useAnnotationStore } from '../stores/annotationStore';
import { usePdfStore } from '../stores/pdfStore';
import { resetLSCurrentTranslateTabKeyByDuration } from '../util/translate';
import GlossaryPopoverTableVue from '@/components/Translate/Glossary/PopoverTable.vue';
import GlossaryCreateEditVue from '@/components/Translate/Glossary/CreateEdit.vue';
import { GlossaryItem } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/GlossaryManage';
import { useTranslateLock } from '../hooks/useTranslation';
import { useRightSideTabSettings } from '../hooks/UserSettings/useSideTabSettings';
import { RightSideBarType } from '../components/Right/TabPanel/type';
import { debounce } from 'lodash-es';
import { ElementClick, PageType, reportElementClick } from '../api/report';

export interface TippyVueItem {
  tippy: Instance;
  app: App;
}

export const destroyTippyVues = (items?: TippyVueItem[]) => {
  if (!items?.length) {
    return;
  }
  items.forEach(({ tippy, app }) => {
    app.unmount();
    tippy.destroy();
  });
};

export const createReferenceTippyVue = ({
  marker,
  triggerEle,
}: {
  marker: ReferenceMarker;
  triggerEle: Element;
}): TippyVueItem => {
  const annotationStore = useAnnotationStore();
  const div = document.createElement('div');
  const app = createApp(ReferenceVue, {
    paperId: marker.paperId,
    paperTitle: marker.refRaw || marker.refContent,
    fetchFlag: true,
    updateContentHandler: () => {
      tippyInstance.setContent(instance.$el);
    },
    hideTippy: () => {
      tippyInstance.hide();
    },
    marker,
    showCiteBtn: true,
  });

  app.use(Antd);
  app.use(PerfectScrollbar);
  app.use(i18n);

  const instance = app.mount(div);
  const tippyInstance = tippy(triggerEle, {
    content: instance.$el,
    trigger: 'manual',
    arrow: false,
    theme: 'ref-paper',
    placement: 'top-end',
    maxWidth: 473,
    interactive: true,
    appendTo: document.body,
    hideOnClick: false,
    showOnCreate: true,
    zIndex: 99,
    onHide(instance) {
      instance.popper.removeAttribute('data-x');
      instance.popper.removeAttribute('data-y');
    },
    onClickOutside(instance) {
      // 引用弹窗打开时不关闭当前tippy
      if (annotationStore.showReferenceTippy) {
        return;
      }
      instance.hide();
    },
  });

  enableTippyDraggable(tippyInstance.popper);

  return {
    tippy: tippyInstance,
    app,
  };
};

export const createFigureTippyVue = ({
  info,
  triggerEle,
}: {
  info: PdfFigureAndTableInfo & { refContent: string };
  triggerEle: Element;
}): TippyVueItem => {
  const div = document.createElement('div');

  const isDingRef = ref(false); // 使用 ref 创建响应式状态

  const app = createApp(FigureVue, {
    info,
    isDing: isDingRef, // 直接传递 ref
    tippyHandler: (event: 'ding' | 'close' | 'unding') => {
      if (event === 'close') {
        tippyInstance.hide();
      } else if (event === 'ding') {
        isDingRef.value = true;
      } else {
        isDingRef.value = false;
      }
    },
  });

  app.use(Antd);
  app.use(PerfectScrollbar);
  app.use(i18n);

  const instance = app.mount(div);

  const tippyInstance = tippy(triggerEle, {
    content: instance.$el,
    trigger: 'manual',
    arrow: false,
    theme: 'ref-paper',
    placement: 'left-start',
    maxWidth: 'none',
    interactive: true,
    appendTo: document.body,
    hideOnClick: false,
    showOnCreate: true,
    popperOptions: {
      strategy: 'fixed',
      modifiers: [
        {
          name: 'flip',
          options: {
            fallbackPlacements: ['bottom', 'right'],
          },
        },
        {
          name: 'preventOverflow',
          options: {
            altAxis: true,
            tether: false,
          },
        },
        {
          name: 'eventListeners',
          enabled: false,
        },
      ],
    },
    onClickOutside(ins) {
      if (!isDingRef.value) {
        ins.hide();
      }
    },
    onHide(instance) {
      instance.popper.removeAttribute('data-x');
      instance.popper.removeAttribute('data-y');
    },
    onShown(ins) {
      ins.setContent(instance.$el);
    },
  });

  const handleFocus = (event: { currentTarget: HTMLElement }) => {
    document.body.querySelectorAll('[data-tippy-root]').forEach((item) => {
      (item as HTMLElement).style.zIndex = '999';
    });
    event.currentTarget.style.zIndex = '1000';
  };

  interact(tippyInstance.popper).on('tap', handleFocus);

  enableTippyDraggable(tippyInstance.popper, {
    listeners: {
      start: handleFocus,
    },
  });
  enableTippyResizable(
    tippyInstance.popper.querySelector('.js-figure-tippy-viewer') as HTMLElement
  );

  return {
    tippy: tippyInstance,
    app,
  };
};

export const checkMultiSegment = (pageTexts: PageSelectText[]) => {
  return (
    pageTexts.some((item) => item.multiSegment) ||
    (pageTexts.length > 1 && pageTexts.some((item) => !item.crossing))
  );
};

export interface TranslatePayload {
  pdfId: string;
  triggerEle?: Element;
  props?: Partial<Props>;
  isExistingAnnotation: boolean;
  pageTexts?: PageSelectText[];
  ocr?: {
    text: TranslateState['content'];
    addOcrNote(translation: string): Promise<void>;
  };
  translatedData?: UniTranslateResp | null;
  resetIdeaTab?: boolean;
  from?: 'selection' | 'tooltip';
}

export const createTranslateTippyVue = (() => {
  let tippyInstance: Nullable<Instance> = null;
  let app: Nullable<App> = null;
  let translateStore: ReturnType<typeof useTranslateStore> | null = null;
  let isDingRef = ref(false); // 使用 ref 创建响应式状态
  let fixed = false;
  const currentNoteInfo: { info: Nullable<PageSelectText[]> } = {
    info: null,
  };
  const show = ({
    pdfId,
    triggerEle,
    props,
    isExistingAnnotation,
    pageTexts,
    ocr,
    translatedData,
    resetIdeaTab,
    from,
  }: TranslatePayload): TippyVueItem | null => {
    if (resetIdeaTab) {
      // 每次创建翻译弹框，校正一下默认tab，超过半个小时，自动切换回idea
      // resetLSCurrentTranslateTabKeyByDuration();
    }

    console.log(
      'createTranslateTippyVue',
      pdfId,
      pageTexts,
      ocr,
      translatedData
    );

    fixed = false;

    currentNoteInfo.info = pageTexts ?? null;

    const { translateLock } = useTranslateLock();

    const { setSideTabSetting, activeTab, sideTabSettings } =
      useRightSideTabSettings();

    const updateTranslateStore = async () => {
      if (!translateStore) {
        return;
      }

      if (ocr) {
        translateStore.setContent(ocr.text);
      } else if (pageTexts) {
        const origin = pageTexts.reduce(
          (prev, current) => prev + current.text,
          ''
        );
        translateStore.setContent({
          origin,
          ocrTranslate: null,
        });
      }

      translateStore.setPdfId(pdfId);
      translateStore.isExistingAnnotation = isExistingAnnotation;
      translateStore.multiSegment = !!pageTexts && checkMultiSegment(pageTexts);
      translateStore.setExtraInfo({
        selections: pageTexts,
        ocr,
        translateData: translatedData,
      });
    };

    console.log('translateLock', translateLock.value, from, activeTab.value);

    if (translateLock.value) {
      if (
        from === 'selection' &&
        (activeTab.value !== RightSideBarType.Translate ||
          !sideTabSettings.value.shown)
      ) {
        console.warn(
          'selection translation is aborted because of the lock status.'
        );
        return null;
      }
      // 翻译锁定状态下，不创建翻译弹框
      translateStore = translateStore || useTranslateStore();
      updateTranslateStore();

      setTimeout(() => {
        setSideTabSetting({
          tab: RightSideBarType.Translate,
          shown: true,
        });
      }, 0);

      console.log('translate is locked.');
      return null;
    }

    if (tippyInstance && app) {
      updateTranslateStore();

      if (props && !isDingRef.value) {
        tippyInstance.setProps(props);
      }

      if (triggerEle) {
        tippyInstance.show();
      }

      return {
        tippy: tippyInstance,
        app,
      };
    } else if (!triggerEle) {
      console.warn('invalid none tirggerEle');
      return null;
    }

    translateStore?.setContent({
      origin: '',
      ocrTranslate: null,
    });

    const div = document.createElement('div');
    const { tippyConfig } = useLocalTippy();

    const pdfStore = usePdfStore();
    const pdfViewer = pdfStore.getViewer(selfNoteInfo.value.pdfId);
    const { add: addNote, addWord } = useCreateNote({
      pdfId,
      noteId: selfNoteInfo.value?.noteId || '',
      pdfViewer,
    });

    app = createApp(TranslateVue, {
      pdfId,
      width: Math.max(tippyConfig.value.translateWidth, 200),
      isDing: isDingRef, // 直接传递 ref
      tippyHandler: (event: 'ding' | 'close' | 'unding' | 'lock') => {
        console.log('tippyHandler', event);
        if (event === 'close') {
          tippyInstance?.hide();
        } else if (event === 'ding') {
          isDingRef.value = true;
          reportElementClick({
            element_name: ElementClick.trans_box_lock,
            page_type: PageType.note,
            type_parameter: 'none',
          });
        } else if (event === 'unding') {
          isDingRef.value = false;
        } else if (event === 'lock') {
          translateLock.value = true;
          tippyInstance?.hide();
          setTimeout(() => {
            setSideTabSetting({
              tab: RightSideBarType.Translate,
              shown: true,
            });
          }, 300);
          reportElementClick({
            element_name: ElementClick.trans_lock_right,
            page_type: PageType.note,
            type_parameter: 'none',
          });
        }
      },
      addToNoteHandler: (
        isPhrase: boolean,
        phrase: string,
        translation: string,
        translationRes: UniTranslateResp
      ) => {
        if (isPhrase) {
          addWord(currentNoteInfo, phrase, translationRes);
        } else if (translateStore?.content.ocrTranslate) {
          ocr?.addOcrNote(translation);
        } else {
          addNote(currentNoteInfo, translation);
        }
      },
      async fixPlacement() {
        if (fixed) {
          return;
        }

        fixed = true;

        if (isDingRef.value) {
          return;
        }

        await nextTick();
        tippyInstance?.setProps({
          placement: tippyInstance.props.placement,
          offset: tippyInstance.props.offset,
        });
      },
    });

    app.use(Antd);
    app.use(PerfectScrollbar);

    app.use(i18n);

    const pinaStore = usePinaStore();
    app.use(pinaStore);

    translateStore = useTranslateStore();
    updateTranslateStore();

    const instance = app.mount(div);

    const defaultProps: Partial<Props> = {
      content: instance.$el,
      trigger: 'manual',
      arrow: false,
      theme: 'ref-paper',
      placement: 'left-start',
      maxWidth: 'none',
      interactive: true,
      appendTo: document.body,
      hideOnClick: false,
      showOnCreate: true,
      popperOptions: {
        strategy: 'fixed',
        modifiers: [
          {
            name: 'flip',
            options: {
              fallbackPlacements: ['bottom', 'right'],
            },
          },
          {
            name: 'preventOverflow',
            options: {
              altAxis: true,
              tether: false,
            },
          },
          {
            name: 'eventListeners',
            enabled: false,
          },
        ],
      },
      onClickOutside(ins, event) {
        if (!isDingRef.value) {
          ins.hide();
        }
      },
      onHidden(instance) {
        setTimeout(() => {
          instance.popper.removeAttribute('data-x');
          instance.popper.removeAttribute('data-y');
          app?.unmount();
          instance.destroy();
          app = null;
          tippyInstance = null;
        }, 0);
      },
      onShown(ins) {
        ins.setContent(instance.$el);
      },
    };

    tippyInstance = tippy(triggerEle, merge(defaultProps, props));

    const handleFocus = (event: { currentTarget: HTMLElement }) => {
      document.body.querySelectorAll('[data-tippy-root]').forEach((item) => {
        (item as HTMLElement).style.zIndex = '999';
      });
      event.currentTarget.style.zIndex = '1000';
    };

    interact(tippyInstance.popper).on('tap', handleFocus);

    enableTippyDraggable(tippyInstance.popper, {
      listeners: {
        start: handleFocus,
      },
    });
    enableTippyResizable(
      tippyInstance.popper.querySelector(
        '.js-translate-tippy-viewer'
      ) as HTMLElement,
      {
        edges: { left: true, right: true, bottom: false, top: false },
        listeners: {
          end: (event: Interact.ResizeEvent) => {
            const target = event.target as HTMLElement;
            tippyConfig.value.translateWidth = target.offsetWidth;
          },
        },
      }
    );

    return {
      tippy: tippyInstance,
      app,
    };
  };
  const hide = () => {
    if (tippyInstance) {
      tippyInstance.hide();
    }
  };
  const isPinned = () => {
    return isDingRef.value;
  };
  return {
    show: debounce(show, 300),
    hide,
    isPinned,
  };
})();

export const enableTranslateTippyVueToAddNote = (enable: boolean) => {
  const translateStore = useTranslateStore();
  translateStore.enableAllowToAddNote(enable);
};
// 全局一个
export const createGlossaryTableTippyVue = (() => {
  let tippyInstance: Nullable<Instance> = null;
  return ({ triggerEle }: { triggerEle: HTMLElement }) => {
    // if (tippyInstance) {
    //   tippyInstance.show();
    //   return {
    //     tippy: tippyInstance,
    //   };
    // }

    tippyInstance?.hide();

    const div = document.createElement('div');
    const app = createApp(GlossaryPopoverTableVue, {
      close: () => {
        tippyInstance?.hide();
      },
    });
    app.use(Antd);
    app.use(PerfectScrollbar);
    app.use(i18n);

    const instance = app.mount(div);

    tippyInstance = tippy(triggerEle, {
      content: instance.$el,
      trigger: 'manual',
      arrow: false,
      theme: 'ref-paper',
      placement: 'right-start',
      maxWidth: 'none',
      interactive: true,
      appendTo: document.body,
      hideOnClick: false,
      showOnCreate: true,
      popperOptions: {
        strategy: 'fixed',
        modifiers: [
          {
            name: 'flip',
            options: {
              fallbackPlacements: ['bottom', 'right'],
            },
          },
          {
            name: 'preventOverflow',
            options: {
              altAxis: true,
              tether: false,
            },
          },
          {
            name: 'eventListeners',
            enabled: false,
          },
        ],
      },
      onHide() {
        tippyInstance = null;
      },
      onShown(ins) {
        ins.setContent(instance.$el);
      },
    });

    enableTippyDraggable(tippyInstance.popper);

    return {
      tippy: tippyInstance,
    };
  };
})();

export const createGlossaryCreateEditTippyVue = (() => {
  let tippyInstance: Nullable<Instance> = null;
  return ({
    triggerEle,
    item,
    refresh,
  }: {
    triggerEle: HTMLElement;
    item: GlossaryItem | null;
    refresh: () => void;
  }) => {
    console.log(triggerEle, item);

    // 只能开启一个编辑框
    tippyInstance?.hide();

    const div = document.createElement('div');
    const app = createApp(GlossaryCreateEditVue, {
      close: () => {
        tippyInstance?.hide();
      },
      item,
      refresh,
    });
    app.use(Antd);
    app.use(PerfectScrollbar);
    app.use(i18n);

    const instance = app.mount(div);

    tippyInstance = tippy(triggerEle, {
      content: instance.$el,
      trigger: 'manual',
      arrow: false,
      theme: 'ref-paper',
      placement: 'right-start',
      maxWidth: 'none',
      interactive: true,
      appendTo: document.body,
      hideOnClick: false,
      showOnCreate: true,
      popperOptions: {
        strategy: 'fixed',
        modifiers: [
          {
            name: 'flip',
            options: {
              fallbackPlacements: ['bottom', 'right'],
            },
          },
          {
            name: 'preventOverflow',
            options: {
              altAxis: true,
              tether: false,
            },
          },
          {
            name: 'eventListeners',
            enabled: false,
          },
        ],
      },
      onHide() {
        tippyInstance = null;
      },
      onShown(ins) {
        ins.setContent(instance.$el);
      },
    });

    enableTippyDraggable(tippyInstance.popper);

    return {
      tippy: tippyInstance,
    };
  };
})();
