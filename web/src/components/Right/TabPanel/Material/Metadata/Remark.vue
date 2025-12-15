<template>
  <div
    v-show="props.docInfo.remark || editFlag"
    class="metadata-remark-container"
  >
    <div
      ref="shadowRef"
      class="metadata-remark-shadow"
    >
      {{ props.docInfo.remark }}
    </div>
    <div
      ref="lightRef"
      :style="
        editFlag
          ? {
            display: 'none',
          }
          : expand
            ? {
              display: 'block',
              wordBreak: 'break-word',
            }
            : {
              display: '-webkit-box',
            }
      "
      class="metadata-remark-light"
      @click="startEdit()"
    >
      {{ props.docInfo.remark }}
      <div
        v-show="overflow || expand"
        class="metadata-remark-expand"
        @click.stop="expand = !expand"
      >
        <UpOutlined v-if="expand" />
        <DownOutlined v-else />
      </div>
    </div>
    <Textarea
      v-if="editFlag"
      ref="textareaRef"
      v-model:value="editValue"
      class="metadata-remark-textarea thin-scroll"
      :autosize="{ minRows: 1, maxRows: 18 }"
      @blur="submitEdit()"
      @press-enter="enter($event)"
    />
  </div>
</template>

<script lang="ts" setup>
import { nextTick, ref } from 'vue';
import { Textarea } from 'ant-design-vue';
import { DownOutlined, UpOutlined } from '@ant-design/icons-vue';
import { trueThrottle } from '@idea/aiknowledge-special-util/throttle'
import { updateDocRemark } from '~/src/api/material';
import { useTextOverflow } from './helper';
import { DocDetailInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/client/ClientDoc';

const props = defineProps<{
  docInfo: DocDetailInfo;
}>();

const editFlag = ref(false);
const editValue = ref('');
const textareaRef = ref<HTMLTextAreaElement | null>(null);

const { lightRef, shadowRef, overflow, expand } = useTextOverflow();

const startEdit = async () => {
  if (!props.docInfo) {
    return;
  }

  editValue.value = props.docInfo.remark;
  editFlag.value = true;
  await nextTick();
  textareaRef.value!.focus();
};

const cancelEdit = () => {
  editFlag.value = false;
  editValue.value = '';
};

const submitEdit = trueThrottle(async () => {
  try {
    await updateDocRemark({
      docId: props.docInfo.docId,
      remark: editValue.value,
    });

    props.docInfo.remark = editValue.value;
    cancelEdit();
  } catch (error) {}
}, 300, false, true);

const enter = async (event: KeyboardEvent) => {
  if (event.shiftKey) {
    return;
  }

  if (event.ctrlKey || event.altKey || event.metaKey) {
    const dom = document.querySelector<HTMLTextAreaElement>(
      '.metadata-remark-textarea'
    )!;
    const oldValue = editValue.value;
    const start = dom.selectionStart;
    const end = dom.selectionEnd;
    const newValue = oldValue.slice(0, start) + '\n' + oldValue.slice(end);
    editValue.value = newValue;
    await nextTick();
    dom.selectionStart = start + 1;
    dom.selectionEnd = dom.selectionStart;
    return;
  }

  event.preventDefault();
  submitEdit();
};

defineExpose({
  editFlag,
  startEdit,
});
</script>

<style lang="less">
.metadata-remark-container {
  position: relative;
  margin-top: -6px;
  margin-left: -8px;
}

.metadata-remark-shadow,
.metadata-remark-light {
  font-size: 13px;
  padding-top: 6px;
  padding-bottom: 6px;
  padding-left: 8px;
  border-radius: 2px;
  word-break: break-word;
  min-height: 36px;
}
.metadata-remark-shadow {
  position: absolute;
  display: block;
  opacity: 0;
  pointer-events: none;
}
.metadata-remark-light {
  position: relative;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 6;
  overflow: hidden;
  color: rgba(255, 255, 255, 0.65);
  cursor: pointer;
  .metadata-remark-expand {
    position: absolute;
    bottom: 6px;
    right: 0;
    padding-right: 2px;
    width: 40px;
    height: 20px;
    display: flex;
    align-items: center;
    background: linear-gradient(to right, transparent -14.1%, #383a3d 58.08%);
    cursor: pointer;
    display: flex;
    justify-content: flex-end;
  }
  &:hover {
    background-color: #4e5153;
    .metadata-remark-expand {
      background: linear-gradient(to right, transparent -14.1%, #4e5153 58.08%);
    }
  }
}

.metadata-remark-textarea {
  width: 100%;
  padding-left: 8px !important;
  padding-right: 8px !important;
  padding-top: 6px;
  min-height: 36px !important;
  line-height: 28px !important;
  font-size: 13px !important;
  color: #1d2229 !important;
  background-color: #fff !important;
  outline: 0;
  // border: 1px solid #1f71e0;
  // border-radius: 2px;
  resize: none;
}
</style>
