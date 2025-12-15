<!-- eslint-disable vue/no-v-html -->
<template>
  <div
    :class="[
      'remark-cell',
      {
        added: !editFlag && !remarkEscape,
        'hover-background': !editFlag,
      },
    ]"
    @click="startEdit()"
  >
    <Tooltip
      v-if="!editFlag && remarkEscape"
      placement="right"
      :title="item.remark"
      :get-popup-container="getPopupContainer"
      destroy-tooltip-on-hide
      :mouse-enter-delay="0.2"
      arrow-point-at-center
    >
      <span
        style="cursor: pointer; padding-left: 8px; padding-right: 8px"
        v-html="remarkEscape"
      ></span>
    </Tooltip>
    <Textarea
      v-if="editFlag"
      ref="textareaVm"
      v-model:value="editValue"
      :auto-size="{ minRows: 1, maxRows: 5 }"
      class="list-thin-scroll list-item-textarea"
      @blur="submitEdit()"
      @pressEnter="enter($event)"
      @keyup.capture="keyupEsc($event)"
      @click.stop
    />
  </div>
</template>
<script lang="ts" setup>
import { computed, nextTick, ref } from 'vue'
import { Textarea, Tooltip } from 'ant-design-vue'
import { UserDocInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/UserDocManage'
import { escapeHightlight, getPopupContainer } from '../helper'
import { updateDocRemark } from '@/api/material'

import { useLibraryList } from '@/stores/library/list'

const props = defineProps({
  item: {
    type: Object as () => Required<UserDocInfo>,
    required: true,
  },
  docId: {
    type: String,
    default: '',
  },
})

const storeLibraryList = useLibraryList()

const editFlag = ref(false)
const editValue = ref('')
const remarkEscape = computed(() => {
  if (props.item.searchResult.hitRemark) {
    return escapeHightlight(props.item.searchResult.hitRemark)
  }

  return props.item.remark
})
const textareaVm = ref()
let submitted = false
const startEdit = async () => {
  submitted = false
  editValue.value = props.item.remark
  editFlag.value = true
  await nextTick()
  const dom = textareaVm.value?.$el as HTMLTextAreaElement
  dom.focus()
}
const cancelEdit = () => {
  editFlag.value = false
  editValue.value = ''
}
const keyupEsc = (event: KeyboardEvent) => {
  if (event.key !== 'Escape') {
    return
  }

  event.preventDefault()
  submitted = true
  cancelEdit()
}
const submitEdit = async () => {
  if (submitted) {
    return
  }

  submitted = true

  const remark = editValue.value.trim()

  if (props.item.remark !== remark) {
    await updateDocRemark({
      docId: props.docId,
      remark,
    })
    // eslint-disable-next-line vue/no-mutating-props
    props.item.remark = remark
    storeLibraryList.paperListAll = [...storeLibraryList.paperListAll]
  }

  cancelEdit()
}
const enter = async (event: KeyboardEvent) => {
  if (event.shiftKey) {
    return
  }

  if (event.ctrlKey || event.altKey || event.metaKey) {
    const oldValue = editValue.value
    const dom = textareaVm.value?.$el as HTMLTextAreaElement
    const start = dom.selectionStart
    const end = dom.selectionEnd
    const newValue = oldValue.slice(0, start) + '\n' + oldValue.slice(end)
    editValue.value = newValue
    await nextTick()
    dom.selectionStart = start + 1
    dom.selectionEnd = dom.selectionStart
    return
  }

  event.preventDefault()
  submitEdit()
}
</script>
<style lang="less" scoped>
.remark-cell {
  display: flex;
  align-items: center;
  overflow: hidden;
  min-height: 32px;
  cursor: pointer;
  &.hover-background:hover {
    background-color: #dce0e5;
  }
  &.added:hover {
    &::after {
      content: '';
      width: 100%;
      padding: 0;
      background: #dfe6f0;
      cursor: pointer;
    }
  }
  span {
    cursor: pointer;
    max-width: 100%;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .remark-add {
    cursor: pointer;
    color: #a19f9d;
  }
  textarea {
    width: 100%;
    margin-top: 4px;
    margin-bottom: 4px;
    border: 1px solid #1f71e0;
    outline: 0;
    display: block;
  }
}
</style>
