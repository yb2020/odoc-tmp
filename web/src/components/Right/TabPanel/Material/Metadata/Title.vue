<template>
  <div class="metadata-title-container">
    <div
      ref="triggerRef"
      class="metadata-title-body"
      @click.stop
    >
      <div
        v-if="!editFlag"
        class="metadata-title-view"
        :style="{ cursor: props.docId ? 'pointer' : 'initial' }"
        @click="editStart()"
      >
        {{ docInfo?.docName ?? '' }}
      </div>
      <textarea
        v-else
        ref="editTextarea"
        v-model="editValue"
        class="metadata-title-edit thin-scroll"
        @blur="editSubmit()"
        @keypress.enter.prevent="editSubmit()"
        @keyup="keyup($event)"
      />
    </div>
    <TippyVue
      v-if="triggerRef && props?.paperData"
      ref="tippyRef"
      :trigger-ele="triggerRef"
      placement="left-start"
      trigger="mouseenter"
      :delay="[0, 500]"
      :offset="[0, 12]"
      @onShown="onShown"
    >
      <ReferenceVue
        :paper-id="props.paperData.paperId"
        :paper-title="props.paperData.title"
        :fetch-flag="startFetch"
        :no-collect="true"
      />
    </TippyVue>
    <div class="metadata-title-hover" />
  </div>
</template>
<script setup lang="ts">
import { nextTick, ref } from 'vue';
import { PaperDetailInfo } from 'go-sea-proto/gen/ts/paper/Paper'
import { DocDetailInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/client/ClientDoc';
import TippyVue from '~/src/components/Tippy/index.vue';
import ReferenceVue from '~/src/components/Tippy/Reference/index.vue';
import { updateDoc } from '~/src/api/material';
import { trueThrottle } from '@idea/aiknowledge-special-util/throttle';
import { selfNoteInfo } from '~/src/store';

const props = defineProps<{
  paperData: PaperDetailInfo | null;
  docInfo: DocDetailInfo | null;
  docId?: string;
}>();

console.warn(props);

const editFlag = ref(false);
const editValue = ref('');
const editTextarea = ref<HTMLTextAreaElement | null>(null);
const editStart = async () => {
  if (!props.docId) {
    return;
  }

  editValue.value = props.docInfo!.docName;
  editFlag.value = true;
  await nextTick();
  editTextarea.value!.focus();
};
const editCancel = () => {
  editFlag.value = false;
};

const keyup = (event: KeyboardEvent) => {
  if (event.code === 'Escape') {
    editCancel();
    editValue.value = '';
  }
}

const editSubmit = trueThrottle(async () => {
  const docName = editValue.value;
  if (!docName) {
    return;
  }
  await updateDoc({
    docId: props.docId!,
    docName,
  });
  props.docInfo!.docName = docName;
  selfNoteInfo.value.docName = docName;
  editCancel();
}, 300, false, true);
const triggerRef = ref<HTMLDivElement | null>(null);
const tippyRef = ref<any>(null);

const handleTippyUpdate = () => {
  tippyRef.value.update();
};

const startFetch = ref(false);

const onShown = () => {
  startFetch.value = true;
};
</script>

<style scoped lang="less">
@import '~/src/assets/less/style.less';
.metadata-title-container {
  position: relative;
  margin-left: -8px;
  margin-right: -8px;
  .metadata-title-view {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    border-radius: 2px;
    max-height: 52px;
    line-height: 22px;
    padding-top: 4px;
    padding-bottom: 4px;
    padding-left: 8px;
    padding-right: 8px;
    color: var(--site-theme-text-primary);
  }
  .metadata-title-edit {
    display: block;
    width: 100%;
    height: 52px;
    line-height: 22px;
    padding-top: 3px;
    padding-bottom: 3px;
    padding-left: 7px;
    padding-right: 7px;
    color: var(--site-theme-text-primary);
    outline: 0;
    border: 1px solid var(--site-theme-brand);
    background-color: var(--site-theme-bg-light);
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
