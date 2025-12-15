<template>
  <div
    class="page-container relative flex items-center cursor-pointer p-2 text-xs leading-3"
    @click.stop="showSelect = !showSelect"
  >
    <input
      v-model="activePage"
      :maxlength="(numPages + '').length"
      class="page-input outline-0 border-0 bg-transparent"
      :style="{ width: ((numPages + '').length || 1) + 1 + 'ch' }"
      @blur="handleBlur"
    >
    <div class="page-sum">
      / {{ numPages }}
    </div>
    <SelectOptions
      v-model:visible="showSelect"
      class="w-full"
      :options="pageOptions"
      @selectChange="changePage"
    />
  </div>
</template>

<script lang="ts" setup>
export type ToolbarPageEvent = {
  type: 'toolbar:page';
  pageNumber: number;
};

import { computed, ref, watch } from 'vue';
import SelectOptions from './Select.vue';
import { ElementClick, reportClick } from '~/src/api/report';

const emit = defineEmits<{
  (event: 'goToPage', payload: ToolbarPageEvent): void;
}>();

const props = defineProps<{ numPages: number; currentPage: number }>();

const activePage = ref(props.currentPage);

watch(
  () => props.currentPage,
  (newVal) => {
    activePage.value = newVal;
  }
);

const pageOptions = computed(() => {
  if (props.numPages <= 0) {
    return [];
  }

  const options: { title: string; value: string }[] = [];
  for (let i = 1; i <= props.numPages; i++) {
    options.push({
      title: i.toString(),
      value: i.toString(),
    });
  }
  return options;
});

const handleClick = () => {
  // goToPDFPage(+activePage.value);
  emit('goToPage', {
    type: 'toolbar:page',
    pageNumber: +activePage.value,
  });
  reportClick(ElementClick.page_adjust);
};

const handleBlur = () => {
  if (isNaN(+activePage.value)) {
    activePage.value = 1;
  }

  if (+activePage.value > props.numPages) {
    activePage.value = props.numPages;
  } else if (+activePage.value < 1) {
    activePage.value = 1;
  }

  handleClick();
};

const changePage = (page: string) => {
  activePage.value = parseInt(page);
  handleClick();
};

const showSelect = ref<boolean>(false);
</script>

<style lang="less" scoped>
.page-container {
  .page-input {
    color: #595959;
  }

  .page-sum {
    color: #9097a1;
  }
}
</style>
