<template>
  <div class="reply">
    <a-avatar
      class="avatar"
      :src="avatar"
      :size="24"
    >
      <template #icon>
        <UserOutlined />
      </template>
    </a-avatar>
    <textarea
      ref="inputRef"
      v-model.trim="input"
      class="input"
      :placeholder="placeHolder"
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
import { UserOutlined } from '@ant-design/icons-vue';

export default defineComponent({
  components: {
    UserOutlined,
  },
  props: {
    avatar: String,
    placeHolder: String,
  },
  setup(props, { emit }) {
    let input = ref('');

    const handleNote = async () => {
      emit('onBlur', input.value);

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
          input.value.slice(0, selectionStart) + '\n' + input.value.slice(selectionStart);

        setCursorPosition(target, selectionStart + 1);

        setTimeout(() => {
          autosize.update(inputRef.value);
        }, 50);
      } else if (e.keyCode == 13) {
        e.preventDefault();
        inputRef.value.blur();
      }
    };

    const handleBlur = () => {
      handleNote();
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
  margin: 8px 0;
  display: flex;
  align-items: center;

  .avatar {
    flex: 0 0 24px;
    margin-right: 8px;
  }

  .input {
    padding: 8px;
    resize: none;
    width: 100%;
    outline: none;
    height: 32px;
    background: #ffffff;
    border-radius: 2px;
    border: 1px solid #5596f2;

    font-size: 14px;
    font-weight: 400;
    color: rgba(0, 0, 0, 85%);
    line-height: 24px;
  }
}
</style>
