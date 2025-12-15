<template>
  <div
    v-if="!editFlag"
    class="metadata-venue-view"
    :style="{ cursor: props.docId ? 'pointer' : 'initial' }"
    @click="startEdit()"
  >
    {{ props.displayVenue?.venue ?? '' }}
    <Rollback
      v-if="props.docId && props.displayVenue?.rollbackEnable"
      :current-info="props.displayVenue?.venue"
      :origin-info="props.displayVenue?.originVenue"
      width="360px"
      @click="rollback()"
    />
  </div>
  <input
    v-else
    ref="inputRef"
    v-model="editValue"
    class="metadata-venue-input"
    @blur="submitEdit()"
    @keypress.enter.prevent="submitEdit()"
    @keyup="keyup($event)"
  >
</template>

<script setup lang="ts">
import { nextTick, ref } from 'vue';
import { DocDetailInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/client/ClientDoc';
import Rollback from './Rollback.vue';
import { updateVenue } from '~/src/api/material';
import { trueThrottle } from '@idea/aiknowledge-special-util/throttle';

const props = defineProps<{
  docId?: string;
  displayVenue?: DocDetailInfo['displayVenue'];
}>();

const editFlag = ref(false);
const editValue = ref('');
const inputRef = ref<HTMLInputElement | null>(null);
const startEdit = async () => {
  if (!props.docId || !props.displayVenue) {
    return;
  }

  editValue.value = props.displayVenue.venue;
  editFlag.value = true;
  await nextTick();
  inputRef.value?.focus();
};

const cancelEdit = () => {
  editFlag.value = false;
  editValue.value = '';
};

const keyup = (event: KeyboardEvent) => {
  if (event.code === 'Escape') {
    cancelEdit();
  }
}

const submitEdit = trueThrottle(
  async () => {
    if (props.displayVenue!.venue === editValue.value) {
      cancelEdit();
      return;
    }

    const venue = editValue.value;
    if (!venue) {
      return;
    }

    await updateVenue({
      docId: props.docId!,
      venue,
    });

    props.displayVenue!.venue = venue;
    props.displayVenue!.rollbackEnable =
      venue !== props.displayVenue!.originVenue;
    cancelEdit();
  },
  300,
  false,
  true
);

const rollback = async () => {
  props.displayVenue!.venue = props.displayVenue!.originVenue;
  props.displayVenue!.rollbackEnable = false;

  try {
    await updateVenue({
      docId: props.docId!,
    });
  } catch (error) {
    console.warn('todo');
    return;
  }
};
</script>

<style lang="less">
.metadata-venue-view {
  flex: 1 1 100%;
  position: relative;
  padding-left: 8px;
  height: 32px;
  line-height: 32px;
  border-radius: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--site-theme-text-secondary);
  &:hover {
    background: rgba(255, 255, 255, 0.08);
    > div {
      visibility: visible;
    }
  }
}

.metadata-venue-input {
  background-color: var(--site-theme-bg-light);
}
</style>
