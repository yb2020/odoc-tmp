<template>
  <a-tooltip
    :mouse-enter-delay="1"
    placement="bottomRight"
  >
    <span
      class="translate-button"
      @click="handleClick"
    >
      <i
        v-if="!loadingOcr"
        class="aiknowledge-icon icon-translate"
        aria-hidden="true"
      />
      <i
        v-else
        class="aiknowledge-icon icon-loading ocr-spin"
        aria-hidden="true"
      />
    </span>
  </a-tooltip>
</template>
<script setup lang="ts">
import { ToolBarType } from '@idea/pdf-annotate-core';
import { TranslatePayload, createTranslateTippyVue } from '@/dom/tippy';
import { bridgeAdaptor } from '@/adaptor/bridge';
import { PageSelectText } from '~/../../packages/pdf-annotate-viewer/typing';
import { store } from '~/src/store';
import { TOOLTIP_CLASSNAME } from '~/src/dom/tooltip';
import { ref, watch } from 'vue';
import {
  DEFAULT_TRANSLATE_TAB_KEY,
  LSKEY_CURRENT_TRANSLATE_TAB,
  TranslateTabKey,
  useTranslateStore,
} from '~/src/stores/translateStore';
import { useVipStore } from '@common/stores/vip';
import { useUserStore } from '@common/stores/user';
import { postOcrZhCn } from '~/src/api/translate';
import { message } from 'ant-design-vue';
import { ResponseError } from '@common/api/type';
import { ERROR_CODE_NEED_VIP } from '@common/api/const';
import { NeedVipException } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { RectOptions } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { ElementName } from '~/src/api/report';
import { useGlossary } from '~/src/hooks/useGlossary';

const vipStore = useVipStore();
const userStore = useUserStore();
const translateStore = useTranslateStore();

const props = defineProps<{
  toolTipType: ToolBarType;
  rects: RectOptions[];
  pageTexts: PageSelectText[];
  pdfId: string;
  alive: boolean;
  getOcrImage(): Promise<string>;
  annotateId?: string;
  addOcrNote(translation: string): Promise<void>;
}>();

const loadingOcr = ref(0);

// 已移除isOcrVip计算属性，因为积分制不再需要检查VIP类型

watch(
  () => props.alive,
  (alive) => {
    if (!alive) {
      loadingOcr.value = 0;
    }
  }
);

const { glossaryChecked } = useGlossary();

const handleClick = async () => {
  if (!userStore.isLogin()) {
    bridgeAdaptor.login();
    return;
  }

  const tooltip = document.getElementsByClassName(TOOLTIP_CLASSNAME)[0]!;

  const bound = tooltip.getBoundingClientRect();

  const maxRight = Math.max(
    ...(props.rects.map((item) => item.x + item.width) as number[])
  );

  const middleRectsTop =
    (props.rects[props.rects.length - 1].y + props.rects[0].y) / 2;

  const middleBoundTop = (bound.top + bound.bottom) / 2;

  const payload: TranslatePayload = {
    pdfId: props.pdfId,
    triggerEle: tooltip,
    isExistingAnnotation: Boolean(props.annotateId),
    props: {
      offset: [
        bound.right > maxRight ? 0 : middleRectsTop - middleBoundTop,
        Math.max(bound.right, maxRight) - bound.right + 6,
      ],
      placement: 'right',
      hideOnClick: false,
    },
  };

  if (props.toolTipType !== ToolBarType.rect) {
    
    payload.pageTexts = props.pageTexts;
    createTranslateTippyVue.show({
      ...payload,
      resetIdeaTab: true,
    });
    return;
  }

  const ocrId = Math.random();
  loadingOcr.value = ocrId;

  let picBase64: string;
  try {
    picBase64 = await props.getOcrImage();
  } catch (error) {
    message.warn((error as Error)?.message ?? '');
    return;
  }
  // 已改为积分制，不再需要获取OCR剩余次数

  // 确保已获取可用渠道列表
  if (translateStore.tabs.length <= 1) {
    await translateStore.getTabs();
  }

  // 获取用户选择的渠道或默认渠道
  let channel = localStorage.getItem(LSKEY_CURRENT_TRANSLATE_TAB) || DEFAULT_TRANSLATE_TAB_KEY;

  // 检查该渠道是否在可用列表中
  if (!translateStore.tabs.includes(channel as TranslateTabKey)) {
    // 如果不在，使用第一个可用渠道或默认渠道
    channel = translateStore.tabs[0] || DEFAULT_TRANSLATE_TAB_KEY;
  }

  try {
    var result = await postOcrZhCn({
      picBase64,
      pdfId: props.pdfId,
      channel,
      useGlossary: glossaryChecked.value,
    });

    // 刷新积分
    const userStore = useUserStore();
    userStore.refreshUserCredits();
    
    // 检查是否是限制相关的响应（已经被拦截器处理过）
    if (result && result.limitHandled) {
      // 如果已经被拦截器处理，则不打开翻译组件
      if (loadingOcr.value !== ocrId) {
        return;
      }
      loadingOcr.value = 0;
      return;
    }
    
  } catch (error) {
    if (loadingOcr.value !== ocrId) {
      return;
    }

    loadingOcr.value = 0;

    const e = error as ResponseError;
    if (e?.code === ERROR_CODE_NEED_VIP) {
      vipStore.showOcrLimitDialog(e?.message, {
        exception: e?.extra as NeedVipException,
        reportParams: {
          element_name: ElementName.upperWordTranslatePopup,
        },
      });
    }

    return;
  }

  if (loadingOcr.value !== ocrId) {
    return;
  }

  loadingOcr.value = 0;
  payload.ocr = {
    text: {
      origin: result.ocrExtractText as string,
      ocrTranslate: result,
      ocrChannel: channel as TranslateTabKey,
    },
    addOcrNote: props.addOcrNote,
  };

  createTranslateTippyVue.show({
    ...payload,
    resetIdeaTab: true,
  });
};
</script>

<style scoped lang="less">
.translate-button {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  &:hover {
    background: #52565a;
  }
  .aiknowledge-icon {
    font-size: 16px;
  }
}

.ocr-spin {
  color: #1f71e0;
  animation: ocr-spin 1s linear infinite;
}

@keyframes ocr-spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
