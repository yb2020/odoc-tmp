<template>
  <div class="metadata-title-container">
    <div class="metadata-title-body" @click.stop>
      <div v-if="!editFlag" class="metadata-title-view" @click="editStart()">
        {{ docInfo?.title ?? '' }}
      </div>
      <textarea
        v-else
        v-model="editValue"
        class="metadata-title-edit thin-scroll"
        ref="editTextarea"
        @blur="editSubmit()"
        @keypress.enter.prevent="editSubmit()"
        @keyup="keyup($event)"
      ></textarea>
    </div>
    <div class="metadata-title-hover"></div>
  </div>
</template>
<script setup lang="ts">
import { nextTick, ref } from 'vue'
import { trueThrottle } from '@idea/aiknowledge-special-util/throttle'
import { DocMetaInfoSimpleVo } from 'go-sea-proto/gen/ts/doc/CSL'

const props = defineProps<{
  docInfo: DocMetaInfoSimpleVo | null
  inited: boolean
}>()

console.warn(props)

const editFlag = ref(false)
const editValue = ref('')
const editTextarea = ref<HTMLTextAreaElement | null>(null)
const editStart = async () => {
  if (!props.inited) {
    return
  }

  editValue.value = props.docInfo!.title
  editFlag.value = true
  await nextTick()
  editTextarea.value!.focus()
}
const editCancel = () => {
  editFlag.value = false
}

const keyup = (event: KeyboardEvent) => {
  if (event.code === 'Escape') {
    editCancel()
    editValue.value = ''
  }
}

const editSubmit = trueThrottle(
  async () => {
    const docName = editValue.value
    if (!docName) {
      return
    }
    props.docInfo!.title = docName
    editCancel()
  },
  300,
  false,
  true
)
</script>

<style scoped lang="less">
@import '@common/../assets/style.less';
.metadata-title-container {
  position: relative;
  border-width: 1px;
  border-style: solid;
  border-color: #c9cdd4;
  flex: 1 1 100%;
  overflow: hidden;

  .metadata-title-view {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    border-radius: 2px;
    min-height: 52px;
    max-height: 52px;
    line-height: 22px;
    padding-top: 4px;
    padding-bottom: 4px;
    padding-left: 8px;
    padding-right: 8px;
  }
  .metadata-title-edit {
    display: block;
    width: 100%;
    min-height: 32px;
    max-height: 52px;
    line-height: 22px;
    padding-top: 3px;
    padding-bottom: 3px;
    padding-left: 7px;
    padding-right: 7px;
    color: #1d2229;
    outline: 0;
    border: 1px solid #1f71e0;
    background-color: #ffffff;
    border-radius: 2px;
    resize: none;
  }
  .metadata-title-hover {
    display: none;
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    pointer-events: none;
    background-color: #fff;
    opacity: 0.08;
  }
  &:hover {
    .metadata-title-hover {
      display: block;
    }
  }
}
</style>
