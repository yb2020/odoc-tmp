<template>
  <div class="reply">
    <textarea
      ref="inputRef"
      v-model.trim="input"
      class="input"
      @click.stop="() => {}"
      @keydown="handleKeyDown"
      @blur="handleBlur"
      @input="handleInput"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref, watch } from 'vue';
import autosize from 'autosize';
import { message } from 'ant-design-vue';
import { AnnotationAll, AnnotationSelect } from '~/src/stores/annotationStore/BaseAnnotationController';
import { useAnnotationStore } from '~/src/stores/annotationStore';

const annotationStore = useAnnotationStore();

export default defineComponent({
  props: {
    item: {
      type: Object as () => AnnotationAll,
      required: true,
    },
  },
  setup(props, { emit }) {
    let input = ref('');

    const handleNote = async () => {
      if (input.value === '') {
        if (props.item.uuid === 'no-annotation') {
          annotationStore.controller.localDeleteAnnotation(
            props.item.uuid,
            props.item.pageNumber
          );
          annotationStore.currentAnnotationId = '';
        }
        return;
      }

      if (input.value === props.item.idea) {
        return;
      }

      if (
        input.value.length > 2000 ||
        ((props.item as AnnotationSelect).rectStr || '').length > 2000
      ) {
        message.error('单个笔记字数最多2000字！');
        return;
      }

      await annotationStore.controller.patchAnnotation(
        props.item.uuid, 
        {
          idea: input.value,
        }
      );

      input.value = '';
    };

    let inputRef = ref();

    onMounted(() => {
      autosize(inputRef.value);
    });

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

    const handleKeyDown = (e: KeyboardEvent) => {
      const target: any = e.target;

      if (e.keyCode == 13 && (e.ctrlKey || e.metaKey)) {
        const selectionStart = target.selectionStart;

        input.value =
          input.value.slice(0, selectionStart) +
          '\n' +
          input.value.slice(selectionStart);

        setCursorPosition(target, selectionStart + 1);

        setTimeout(() => {
          autosize.update(inputRef.value);
        }, 50);
      } else if (e.keyCode == 13) {
        e.preventDefault();
        emit('blur');
      }
    };

    const handleBlur = () => {
      handleNote();

      emit('blur');
    };

    const handleInput = () => {
      emit('onInput', input.value);
    };

    return {
      inputRef,
      input,
      handleNote,
      handleKeyDown,
      handleBlur,
      handleInput,
    };
  },
});
</script>

<style lang="less" scoped>
.reply {
  margin-top: 10px;
}

.input {
  padding: 8px;
  resize: none;
  width: 100%;
  outline: none;
  height: 40px;
  background: #ffffff;
  border-radius: 2px;
  border: 1px solid #5596f2;

  font-size: 14px;
  font-weight: 400;
  color: rgba(0, 0, 0, 85%);
  line-height: 24px;
}
</style>
