import { computed, ref, watch } from 'vue';
import { isOwner } from '~/src/store';

const NOTE_INFO_IGNORED = 'NOTE_INFO_IGNORED';
// TODO 后续要实现跨tab的sessionStorage
export const ignored = ref(NOTE_INFO_IGNORED in sessionStorage);

watch(ignored, (value) => {
  if (value) {
    sessionStorage.setItem(NOTE_INFO_IGNORED, '');
  } else {
    sessionStorage.removeItem(NOTE_INFO_IGNORED);
  }
});

export const FOOTER_NOTE_INFO_HEIGHT = 132;
export const footerNoteInfoHeight = computed(() => {
  if (isOwner.value) {
    return NaN;
  }

  if (ignored.value) {
    return NaN;
  }

  return FOOTER_NOTE_INFO_HEIGHT + 8;
});
