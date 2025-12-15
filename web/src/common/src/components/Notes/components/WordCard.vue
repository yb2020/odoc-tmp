<template>
  <div
    ref="rootEl"
    class="word-card w-full relative"
    style="background-color: var(--site-theme-bg-secondary);"
    :style="cardStyle"
  >
    <component
      :is="$slots.close || DeleteOutlined"
      class="btn-delete !absolute top-2.5 right-4 py-1.5 cursor-pointer hidden group-hover:block text-base"
      :style="{ color: 'var(--site-theme-text-inverse)' }"
      @click="handleDetele"
    />
    <div
      class="word-tt px-4 py-2.5 text-lg text-rp-neutral-10 bg-[#3A57E640]"
      :style="titleStyle"
    >
      {{ data?.word }}
    </div>
    <div
      class="word-ct relative px-4 pt-3.5 pb-4"
      :class="{ group: isHiddenContent }"
    >
      <div
        v-if="isHiddenContent"
        class="masked h-full w-full text-center group-hover:hidden flex justify-center items-center"
        :style="contentStyle"
      >
        <img :src="HiddenWord">
      </div>
      <div
        class="overflow-hidden group-hover:flex flex flex-col"
        :class="{ hidden: isHiddenContent }"
        :style="contentStyle"
      >
        <!-- 音标 -->
        <div
          v-if="!disabledPhonetic && prons.length"
          class="phonetic flex flex-wrap items-center gap-2 mb-3"
        >
          <WordPronunciation
            v-for="item in prons"
            class="px-2 py-0.5 border border-solid border-rp-dark-8 rounded-s"
            :prefix="item.prefix"
            :title="item.symbol"
            :type="item.format"
            :icon="HornImg"
            :audio="item.pronunciation"
          />
        </div>
        <!-- 释义 -->
        <div
          class="word-note overflow-auto text-rp-neutral-10 flex-1"
          :class="isLimitedHeight ? 'h-16' : ''"
        >
          <IdeaMarkdown
            ref="textareaRef"
            :raw="textareaV"
            :uniq-id="data.id"
            :editing="isEdit"
            :min-rows="1"
            :max-rows="5"
            :upload="upload"
            @blur="handleSubmit($event)"
            @submit="handleSubmit($event)"
            @click-view="handleEdit"
            @click.stop="() => {}"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { WordInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/NoteModule';
import { computed, ref, StyleValue, watch } from 'vue';
import { DeleteOutlined } from '@ant-design/icons-vue';
import { IdeaMarkdown } from '@idea/aiknowledge-markdown';
import WordPronunciation from '@common/components/Notes/components/WordPronunciation.vue';
import { useWordNoteOp } from '@common/components/Notes/useWordNote';
import { ImageStorageType, uploadImage } from '@common/api/upload';
import HornImg from '~common/assets/images/notes/horn.svg';
import HiddenWord from '~common/assets/images/notes/hidden_word.png';

const props = defineProps<{
  data: WordInfo;
  disabled?: boolean;
  disabledPhonetic?: boolean;
  isLimitedHeight?: boolean;
  isHiddenContent?: boolean;
  cardStyle?: StyleValue;
  titleStyle?: StyleValue;
  contentStyle?: StyleValue;
}>();
const emits = defineEmits<{
  (e: 'deleted'): void;
  (e: 'editing'): void;
  (e: 'edited'): void;
  (e: 'submited', x: string): void;
  (e: 'submitTargetContent', x: string): void;
}>();

const upload = async (src: File | string) => {
  return uploadImage(src, ImageStorageType.markdown);
};

const wid = computed(() => props.data.id || '');
const prons = computed(() => {
  const {
    britishFormat,
    britishPronunciation,
    britishSymbol,
    americaFormat,
    americaPronunciation,
    americaSymbol,
  } = props.data.translateInfo || {};
  return [
    britishSymbol
      ? {
          prefix: '英',
          symbol: britishSymbol,
          format: britishFormat,
          pronunciation: britishPronunciation,
        }
      : (false as never),
    americaSymbol
      ? {
          prefix: '美',
          symbol: americaSymbol,
          format: americaFormat,
          pronunciation: americaPronunciation,
        }
      : (false as never),
  ].filter(Boolean);
});
const content = computed(() => {
  return (
    props.data.translateInfo?.targetResp
      ?.map((x) => `${x.part || ''} ${x.targetContent.join(';')}`)
      .join('\n') ||
    props.data.translateInfo?.targetContent[0] ||
    ''
  );
});

const rootEl = ref();
const isEdit = ref<boolean>(false);
const textareaV = ref(content.value);
const textareaRef = ref();

watch(content, () => (textareaV.value = content.value));

const handleEdit = () => {
  if (props.disabled) {
    return;
  }
  isEdit.value = true;
  emits('editing');
  // nextTick(() => {
  //   textareaRef.value?.focus();
  // });
};

const handleSubmit = async (v: string) => {
  if (isUpdating.value) {
    return;
  }
  // 这里不需要等待
  updateWord({
    wordInfo: {
      ...props.data.translateInfo,
      targetContent: [v],
      targetResp: [],
      glossaryList: [],
    },
  });
  isEdit.value = false;
  textareaV.value = v;
  emits('edited');
  emits('submited', v);
  emits('submitTargetContent', v);
};

const handleDetele = async () => {
  if (isDeleting.value) {
    return;
  }
  // Modal.confirm({
  //   title: t('message.confirmToDeleteNoteTip'),
  //   okText: t('common.text.delete'),
  //   onOk: async () => {
  await delWord();
  emits('deleted');
  //   },
  //   okButtonProps: {
  //     danger: true,
  //   },
  //   cancelButtonProps: { type: 'primary' },
  // });
};

const { updateWord, isUpdating, delWord, isDeleting } = useWordNoteOp(wid);

defineExpose({
  isEdit,
  rootEl,
  handleEdit,
});
</script>

<style lang="less" scoped>
@import '~common/assets/functions.less';
.word-card {
  .btn-delete {
    top: 10px;
    right: 16px;
    color: #86919c;
    display: none;
  }

  &:hover {
    .btn-delete {
      display: block;
    }
  }

  :deep(.pronunciation-item) {
    padding: 2px 8px;

    border: 1px solid theme('colors.rp-dark-8');
    border-radius: 4px;

    .content {
      color: theme('colors.rp-neutral-8');
      font-size: 13px;
      line-height: 20px;
      padding: 0;
      background: transparent;
      .icon {
        width: 16px;
      }
    }

    &.active {
      border-color: theme('colors.rp-darkblue-7');
      .content {
        color: theme('colors.rp-darkblue-7');
      }
    }
  }

  :deep(.idea-markdown-view-container) {
    .psm(rgba(201, 205, 212, 1));
    padding: 0;
    border: 1px solid transparent;
  }
  :deep(.idea-markdown-edit-container) {
    // 与滚动条占用宽度保持一致
    padding: 0 6px 0 0;
  }
}
</style>
