/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { DocFolder } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/doc/folder';
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { store } from '~/src/store';

export interface PaperFolder {
  key: string;
  title: string;
  children: PaperFolder[];
}

export const removePrefix = (key: string) => {
  return key.split('-').pop() as string;
};

export const useTextOverflow = () => {
  const lightRef = ref<HTMLDivElement | null>(null);
  const shadowRef = ref<HTMLDivElement | null>(null);
  const overflow = ref(false);
  const expand = ref(false);

  let observer: IntersectionObserver | null = null;
  onMounted(() => {
    observer = new IntersectionObserver(() => {
      overflow.value =
        lightRef.value!.offsetHeight < shadowRef.value!.offsetHeight;
    });
    observer.observe(lightRef.value!);
    observer.observe(shadowRef.value!);
  });
  onUnmounted(() => {
    observer!.disconnect();
    observer = null;
  });

  return {
    lightRef,
    shadowRef,
    overflow,
    expand,
  };
};

export type SelectedFolderMap = Record<DocFolder['id'], DocFolder>;

export const getCharLengthList = (string: string) => {
  const lengthList = Array(string.length).fill(1);
  for (let i = 0; i < string.length; i += 1) {
    const code = string.charCodeAt(i);
    if (code < 0 || code > 128) {
      lengthList[i] += 1;
    }
  }

  return lengthList;
};

export const limit30 = (string: string) => {
  const lengthList = getCharLengthList(string);
  const charList = [];
  let sum = 0;
  for (let i = 0; i < string.length; i += 1) {
    if (sum >= 30) {
      charList.push('…');
      break;
    }

    charList.push(string[i]);
    sum += lengthList[i];
  }

  return charList.join('');
};
