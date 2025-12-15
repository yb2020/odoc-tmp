<template>
  <div
    v-show="!hidden"
    ref="rootEl"
    class="aihighlight-card relative pt-1.5"
    :style="styles"
  >
    <u
      class="absolute h-3.5 w-full left-0 top-0 rounded-[6px]"
      :style="{
        backgroundColor: data.color,
      }"
    />
    <div class="inner p-3 bg-white rounded-[6px] text-rp-neutral-10">
      <h3
        class="mb-2 flex items-center text-[15px] text-inherit font-medium leading-6"
      >
        <span class="flex-1">{{ data.scimType
        }}{{
          $t('common.symbol.parenthesis', [$t('aiHighlighter.translation')])
        }}</span>
        <a-tooltip placement="left">
          <template #title>
            {{ $t('translate.addToAnnotation') }}
          </template>
          <i
            class="aiknowledge-icon flex items-center h-6 text-2xl icon-add-to-note cursor-pointer"
            @click="addToNote"
          />
        </a-tooltip>
      </h3>
      <LoadingOutlined
        v-if="!data.translation"
        class="text-sm"
      />
      <p
        v-else
        class="mb-0 text-sm leading-[22px]"
      >
        {{ data.translation }}
      </p>
    </div>
    <aside class="h-2.5" />
  </div>
</template>

<script setup lang="ts">
import _ from 'lodash';
import { LoadingOutlined } from '@ant-design/icons-vue';
import { SCIMItem } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/scim/AiSCIMInfo';
import { ViewerController } from '@idea/pdf-annotate-viewer';
import { Rectangle } from '@idea/pdf-annotate-core/render/renderSelect';
import { useLocalStorage } from '@vueuse/core';
import { computed, onMounted, ref } from 'vue';
import { useRequest } from 'ahooks-vue';
import { selfNoteInfo } from '@/store';
import {
  DEFAULT_TRANSLATE_TAB_KEY,
  LSKEY_CURRENT_TRANSLATE_TAB,
  TranslateTabKey,
} from '@/stores/translateStore';
import useCreateNote from '@/hooks/note/useCreateNote';
import { useTranslateApi } from '@/hooks/useTranslation';
import { getTextRectByAspectRadio } from '@/util/popup';
import { ElementClick, reportClick } from '@/api/report';
import { useUserStore } from '~/src/common/src/stores/user';

const props = defineProps<{
  hidden?: boolean;
  to: HTMLElement;
  target: SVGGElement;
  data: SCIMItem & {
    translation?: string;
  };
  pdfViewInstance: ViewerController;
}>();

const emit = defineEmits(['noted', 'hide']);

const rootEl = ref();
const styles = computed(() => {
  const { target, to } = props;
  const { translation, content } = props.data;
  if (!target || !to) {
    return {};
  }

  const { width: minWidth } = target.getBoundingClientRect();
  const { width: maxWidth } = to.getBoundingClientRect();
  const width = getTextRectByAspectRadio(
    translation || content,
    4 / 3,
    {
      'font-size': '14px',
      'line-height': '22px',
    },
    minWidth,
    maxWidth
  );

  return {
    width: `${width}px`,
    minWidth: `${minWidth}px`,
    maxWidth: `${maxWidth}px`,
  };
});

const transChannel = useLocalStorage<TranslateTabKey>(
  LSKEY_CURRENT_TRANSLATE_TAB,
  DEFAULT_TRANSLATE_TAB_KEY
);
const { run: fetchTranslate, loading: translating } = useTranslateApi();
const doTranslate = async () => {
  // 用户悬浮一段时间后才开始翻译
  setTimeout(async () => {
    if (props.data.translation || !rootEl.value) {
      return;
    }

    const res = await fetchTranslate(transChannel.value, props.data.content);

    if (res?.targetContent) {
      props.data.translation = res.targetContent;
    }

    if (props.hidden) {
      emit('hide');
    }
    // 刷新积分
    const userStore = useUserStore();
    userStore.refreshUserCredits();
  }, 500);
};

onMounted(() => {
  doTranslate();
});

const { add: addNote } = useCreateNote({
  pdfId: selfNoteInfo.value?.pdfId || '',
  noteId: selfNoteInfo.value?.noteId || '',
  pdfViewer: props.pdfViewInstance,
});

const {
  data: added,
  run: addToNote,
  loading: adding,
} = useRequest(
  async () => {
    if (added.value && adding.value) {
      return;
    }
    if (!props.data.translation) {
      return;
    }

    reportClick(ElementClick.scim_reading);

    const { target } = props;
    const pageNumber = parseInt(
      target.getAttribute('page-number') as string,
      10
    );
    const rects = _.attempt<Rectangle[]>(() =>
      JSON.parse(target.getAttribute('rectangles') as string)
    );
    const rectStr = props.data.content;
    if (Number.isNaN(pageNumber) || rects instanceof Error || !rectStr) {
      return;
    }

    await addNote(
      {
        rectInfo: {
          pageNumber,
          rectStr,
          rectRaw: true,
          rects: rects.map((x) => ({
            ...x,
            text: '',
            shouldScaleText: false,
          })),
        },
      },
      props.data.translation
    );
    emit('noted');

    return true;
  },
  {
    manual: true,
  }
);

defineExpose({
  rootEl,
  translating,
});
</script>

<style lang="less" scoped>
.inner {
  filter: drop-shadow(2px 6px 16px rgba(0, 0, 0, 0.12))
    drop-shadow(4px 10px 28px rgba(0, 0, 0, 0.08));
}
</style>
