<template>
  <div ref="textareaRef">
    <div
      v-if="!editing"
      class="text js-interact-drag-ignore"
      @click="showEdit"
    >
      <div
        class="max-h-40 overflow-y-auto scrollbar"
        :style="{ fontSize: fontSize + 'px', lineHeight: '1.4' }"
      >
        <template v-for="(item, index) in shownInput">
          <span
            v-if="!item.target"
            :key="`text-${index}`"
          >{{
            item.origin
          }}</span>
          <a-popover
            v-else
            :key="`glossary-${index}`"
            :title="null"
            placement="top"
            :getPopupContainer="getPopupContainer"
          >
            <template #content>
              <div class="flex gap-4 justify-between items-center mb-2">
                <div class="text-base">
                  {{ $t('glossary.terms') }}
                </div>
                <ExtraIcons
                  only-add-to-parse
                  :input="item.origin"
                  :translated-content="item.target"
                  :translated-data="{
                    targetResp: [{ targetContent: [item.target], part: '' }],
                    glossaryList: [],
                    targetContent: item.target,
                  }"
                  :add-to-note-handler="addToNoteHandler"
                />
              </div>
              <div>{{ item.target }}</div>
            </template>
            <span class="text-rp-blue-6 cursor-pointer">{{ item.origin }}</span>
          </a-popover>
        </template>
      </div>
    </div>
    <a-textarea
      v-show="editing"
      v-model:value="input"
      class="js-interact-drag-ignore"
      :placeholder="$t('translate.placeholder')"
      :auto-size="{ minRows: 1, maxRows: 6 }"
      :style="{
        color: 'rgba(0,0,0,.85)',
        border: '1px solid #f5f7fa',
        background: '#fefefe',
        boxShadow: '0 4px 16px 0 rgba(12, 53, 115, 0.2)',
        fontSize: '14px',
      }"
      @blur="onTranslate"
      @keydown="handleKeyDown"
    />
  </div>
</template>
<script lang="ts" setup>
import { computed, nextTick, ref } from 'vue';
import { useTranslateStore } from '~/src/stores/translateStore';
import { GlossaryInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/TranslateProto';
import ExtraIcons from './ExtraIcons.vue';
import { UniTranslateResp } from '~/src/api/translate';

const props = defineProps<{
  fontSize: string;
  glossaryList?: GlossaryInfo[];
  addToNoteHandler: (
    isPhrase: boolean,
    phrase: string,
    translation: string,
    translationRes: UniTranslateResp
  ) => void;
}>();

const store = useTranslateStore();
const input = ref(store.content.origin);
const preInput = ref(store.content.origin);

const emit = defineEmits<{
  (event: 'translate', input: string): void;
}>();

const editing = ref(false);
const textareaRef = ref<HTMLDivElement>();
const showEdit = () => {
  editing.value = true;
  nextTick(() => {
    const div = textareaRef.value as HTMLDivElement;
    div.querySelector('textarea')?.focus();
  });
};

/*
 * 设置输入域(input/textarea)光标的位置
 * @param {HTMLInputElement/HTMLTextAreaElement} elem
 * @param {Number} index
 */
function setCursorPosition(elem: HTMLTextAreaElement, index: number) {
  const val = elem.value;
  const len = val.length; // 超过文本长度直接返回

  if (len < index) return;

  setTimeout(function () {
    elem.focus();

    if (elem.setSelectionRange) {
      // 标准浏览器
      elem.setSelectionRange(index, index);
    }
  }, 10);
}

const onTranslate = () => {
  editing.value = false;
  if (input.value === preInput.value) {
    return;
  }
  nextTick(() => {
    emit('translate', input.value);
    preInput.value = input.value;
  });
};

const handleKeyDown = (e: KeyboardEvent) => {
  return;
  const target: any = e.target;

  if (e.keyCode == 13 && (e.ctrlKey || e.metaKey)) {
    const selectionStart = target.selectionStart;

    input.value =
      input.value.slice(0, selectionStart) +
      '\n' +
      input.value.slice(selectionStart);

    setCursorPosition(target, selectionStart + 1);
  } else if (e.keyCode == 13) {
    e.preventDefault();

    onTranslate();
  }
};

const shownInput = computed<{ origin: string; target: string }[]>(() => {
  if (!props.glossaryList?.length) {
    return [{ origin: input.value, target: '' }];
  }
  const glossaryList = props.glossaryList.sort((a, b) => a.start - b.start);
  const arr: { origin: string; target: string }[] = [];
  let start = 0;
  console.log('input.value', glossaryList);
  glossaryList.forEach((glossary) => {
    if (glossary.start > start) {
      arr.push({
        origin: input.value.slice(start, glossary.start),
        target: '',
      });
    }
    arr.push({
      origin: input.value.slice(glossary.start, glossary.end + 1),
      target: glossary.translationText,
    });
    start = glossary.end + 1;
  });
  if (start < input.value.length) {
    arr.push({ origin: input.value.slice(start), target: '' });
  }
  return arr;
});

const getPopupContainer = () =>
  (textareaRef.value as HTMLElement) || document.body;

defineExpose({
  startTranslate: onTranslate,
  getOriginInput: () => input.value,
  updateOriginInput: (value: string) => {
    input.value = value;
    preInput.value = value;
  },
});
</script>
<style scoped lang="less">
.text {
  background: #f5f7fa;
  border-radius: 4px;
  font-size: 14px;
  font-family: Lato-Regular, Lato;
  font-weight: 400;
  color: #1d2129;
  line-height: 22px;
  padding: 8px;
  cursor: text;

  .ellipsis {
    text-overflow: ellipsis;
    overflow: hidden;
    -webkit-line-clamp: 6;
    -webkit-box-orient: vertical;
    display: -webkit-box;
  }
}

.scrollbar {
  scrollbar-width: thin;
  scrollbar-color: #d1d5da #f5f7fa;

  &::-webkit-scrollbar {
    width: 6px; //修改垂直滚动条宽度
    height: 6px; //修改水平滚动条宽度
  }

  &::-webkit-scrollbar-thumb {
    border-radius: 10px;
    background: transparent;
  }
}
</style>
