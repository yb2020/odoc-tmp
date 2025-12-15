import { createSharedComposable, useLocalStorage } from '@vueuse/core';
import { PDF_READER } from '@/common/src/constants/storage-keys';

export const LSKeyForGlossary = PDF_READER.GLOSSARY;

export const useGlossary = createSharedComposable(() => {
  const glossaryChecked = useLocalStorage(LSKeyForGlossary, true);

  return {
    glossaryChecked,
  };
});
