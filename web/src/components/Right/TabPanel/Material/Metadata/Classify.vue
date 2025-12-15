<template>
  <div
    class="metadata-classify-cell"
    @click="showInput()"
  >
    <span
      v-if="!showSelectInput"
      class="metadata-classify-cell-view"
    >
      <TagVue
        v-for="tag in currentTags"
        :key="tag.key"
        :closable="false"
        class="metadata-classify-tag"
        :text="tag.label"
        :font-size="13"
        :dark="true"
      />
    </span>
    <Select
      v-if="showSelectInput"
      ref="selectRef"
      dropdown-class-name="metadata-classify-select"
      :value="currentTags"
      mode="multiple"
      style="width: 100%; margin: 0 !important"
      placeholder="选择标签"
      label-in-value
      :autofocus="true"
      :default-open="true"
      theme="light"
      class="select"
      option-label-prop="label"
      :default-active-first-option="false"
      :not-found-content="null"
      :filter-option="false"
      @select="handleSelect"
      @deselect="handleDeselect"
      @search="handleSearch"
      @input-key-down="handleKeydown"
    >
      <SelectOption
        v-for="item in filteredOptions"
        :key="item.key"
        :value="item.key"
        :label="item.label"
      >
        <TagVue
          style="cursor: pointer"
          class="metadata-classify-tag"
          :closable="false"
          :text="item.label"
        />
      </SelectOption>
    </Select>
  </div>
</template>
<script lang="ts">
import { computed, defineComponent, nextTick, ref } from 'vue';
import { Select } from 'ant-design-vue';

import { onClickOutside } from '@vueuse/core';
import TagVue from './Tag.vue';
import {
  addClassify,
  attachDocToClassify,
  getUserAllClassifyList,
  removeDocFromClassify,
} from '~/src/api/material';
import { DocDetailInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/client/ClientDoc';
import { delay } from '@idea/aiknowledge-special-util/delay';

interface ClassifyLabel {
  key: string;
  label: string;
}

const SelectOption = Select.Option;

function diff<T>(
  array1: T[],
  array2: T[],
  compare: (item1: T, item2: T) => boolean
) {
  const array3 = array1.filter((item1) =>
    array2.every((item2) => !compare(item1, item2))
  );
  return array3;
}

export default defineComponent({
  components: { TagVue, Select, SelectOption },
  props: {
    classifyInfos: {
      type: Array as () => DocDetailInfo['classifyInfos'],
      default: () => [],
    },
    docId: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    const allClassifyList = ref<ClassifyLabel[]>([]);
    const currentTags = computed(() => {
      return (props.classifyInfos || []).map((item) => {
        return {
          key: item.classifyId,
          label: item.classifyName,
          value: item.classifyId,
        };
      });
    });
    const showSelectInput = ref<boolean>(false);
    const filteredOptions = computed<ClassifyLabel[]>(() => {
      const list = diff(
        allClassifyList.value,
        currentTags.value,
        (item1, item2) => item1.key === item2.key
      );
      if (!inputValue.value) {
        return list;
      }
      return list.filter((item) => {
        return item.label.includes(inputValue.value);
      });
    });
    const inputValue = ref<string>('');
    const fetchAllClassifyList = async () => {
      const response = await getUserAllClassifyList({});
      allClassifyList.value =
        response?.map((item) => ({
          key: item.classifyId,
          label: item.classifyName,
        })) ?? [];
    };

    const showInput = async () => {
      if (!allClassifyList.value.length) {
        await fetchAllClassifyList();
      }
      showSelectInput.value = true;
      await nextTick();
      await delay(200);
      document
        .querySelector<HTMLInputElement>(
          '.metadata-classify-cell .ant-select-selection-overflow input'
        )
        ?.focus();
    };

    const updateOfflineClassifyInfos = (item: ClassifyLabel) => {
      if (currentTags.value.every((tag) => tag.key !== item.key)) {
        props.classifyInfos.push({
          classifyId: item.key,
          classifyName: item.label,
        });
      }
    };
    const handleEnter = async () => {
      const value = inputValue.value.trim();
      if (!value) {
        showSelectInput.value = false;
        return;
      }
      const findInCurrent = currentTags.value.find(
        (item) => item.label === value
      );
      if (findInCurrent) {
        showSelectInput.value = false;
        inputValue.value = '';
        return;
      }
      let findInAll = filteredOptions.value.find(
        (item) => item.label === value
      );
      if (!findInAll) {
        try {
          const res = await addClassify({ classifyName: value });
          findInAll = {
            key: res.classifyId,
            label: res.classifyName,
          };
          await fetchAllClassifyList();
        } catch (error) {
          return;
        }
      }
      if (findInAll) {
        const result = await attachDocToClassify({
          docId: props.docId,
          classifyId: findInAll.key,
        });
        if (result) {
          updateOfflineClassifyInfos(findInAll!);
          inputValue.value = '';
          showSelectInput.value = false;
        }
      }
    };

    const handleSelect: any = async ({ key }: { key: string }) => {
      const result = await attachDocToClassify({
        docId: props.docId,
        classifyId: key,
      });
      if (result) {
        const classify = allClassifyList.value.find(
          (item) => item.key === key
        )!;
        updateOfflineClassifyInfos(classify);
      }
    };
    const handleDeselect: any = async ({ key }: { key: string }) => {
      try {
        await removeDocFromClassify({
          classifyId: key || '',
          docId: props.docId,
        });
        fetchAllClassifyList();
        const index = props.classifyInfos.findIndex(
          (item) => item.classifyId === key
        );
        if (index !== -1) {
          props.classifyInfos.splice(index, 1);
        }
      } catch (error) {}
    };

    const selectRef = ref();

    onClickOutside(selectRef, (event: PointerEvent) => {
      if ((event.target as HTMLElement)?.closest('.metadata-classify-select')) {
        return;
      }
      inputValue.value = '';
      showSelectInput.value = false;
    });

    const handleSearch = (val: string) => {
      inputValue.value = val;
    };

    const handleKeydown = (event: KeyboardEvent) => {
      if (event.code === 'Enter' || event.key === 'Enter') {
        setTimeout(async () => {
          await handleEnter();
        }, 500);
      }
    };

    return {
      currentTags,
      showSelectInput,
      filteredOptions,
      inputValue,
      showInput,
      handleSelect,
      handleDeselect,
      selectRef,
      handleKeydown,
      handleSearch,
    };
  },
});
</script>
<style lang="less">
.metadata-classify-cell {
  width: 100%;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  height: 100%;
  min-height: 32px;
  overflow: hidden;
  .metadata-classify-cell-view {
    padding-left: 8px;
    width: 100%;
    min-height: 32px;
    &:hover {
      background: rgba(255, 255, 255, 0.08);
      border-radius: 2px;
    }
  }
  .metadata-classify-tag {
    margin-top: 4px;
    margin-bottom: 4px;
  }
  &:hover {
    cursor: pointer;
    .classify-add {
      display: flex;
      align-items: center;
    }
  }
  .select {
    margin: 4px 0;
    cursor: pointer;
  }
}

[data-theme='dark'] {
  .metadata-classify-select {
    &.ant-select-dropdown {
      background-color: #222326 !important;
    }
    .rc-virtual-list-scrollbar-thumb {
      background-color: #6c737a !important;
    }
  }
}

.metadata-classify-select.ant-select-dropdown {
  margin: 0;
  padding: 0;
  color: rgba(0, 0, 0, 0.85);
  font-variant: tabular-nums;
  line-height: 24px;
  list-style: none;
  font-feature-settings: 'tnum';
  position: absolute;
  top: -9999px;
  left: -9999px;
  z-index: 1050;
  box-sizing: border-box;
  padding: 4px 0;
  overflow: hidden;
  font-size: 14px;
  font-variant: initial;
  background-color: #fff;
  border-radius: 2px;
  outline: none;
  box-shadow:
    0 3px 6px -4px rgb(0 0 0 / 12%),
    0 6px 16px 0 rgb(0 0 0 / 8%),
    0 9px 28px 8px rgb(0 0 0 / 5%);
}
.metadata-classify-select {
  .ant-select-item {
    position: relative;
    display: block;
    min-height: 32px;
    padding: 5px 12px;
    color: rgba(0, 0, 0, 0.85);
    font-weight: normal;
    font-size: 14px;
    line-height: 22px;
    cursor: pointer;
    transition: background 0.3s ease;
  }
  .ant-select-item-option-active:not(.ant-select-item-option-disabled) {
    background-color: #f7f8fa;
  }
}
.metadata-classify-cell {
  .ant-select-focused:not(.ant-select-disabled).ant-select:not(
      .ant-select-customize-input
    )
    .ant-select-selector {
    border-color: #4792ed;
    box-shadow: 0 0 0 2px rgb(31 113 224 / 20%);
    border-right-width: 1px !important;
    outline: 0;
  }
  .ant-select-show-search.ant-select:not(.ant-select-customize-input)
    .ant-select-selector {
    cursor: text;
  }
  .ant-select:not(.ant-select-customize-input) .ant-select-selector {
    position: relative;
    background-color: #fff;
    border: 1px solid #d9d9d9;
    border-radius: 2px;
    transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
  }
  .ant-select-selection-placeholder {
    flex: 1;
    overflow: hidden;
    color: #bfbfbf;
    white-space: nowrap;
    text-overflow: ellipsis;
    pointer-events: none;
  }
  .ant-select-multiple .ant-select-selection-item {
    position: relative;
    display: flex;
    flex: none;
    box-sizing: border-box;
    max-width: 100%;
    height: 28px;
    margin-top: 2px;
    margin-bottom: 2px;
    line-height: 28px;
    // background: #f5f5f5;
    color: #000000a6;
    background-color: #f0f2f5;
    border: 1px solid #f0f0f0;
    border-radius: 2px;
    cursor: default;
    transition:
      font-size 0.3s,
      line-height 0.3s,
      height 0.3s;
    user-select: none;
    margin-inline-end: 3px;
    padding-inline-start: 8px;
    padding-inline-end: 4px;
  }
  .ant-select-multiple .ant-select-selection-item-remove {
    color: inherit;
    font-style: normal;
    line-height: 0;
    text-align: center;
    text-transform: none;
    vertical-align: -0.125em;
    text-rendering: optimizelegibility;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    display: inline-block;
    font-weight: bold;
    font-size: 10px;
    line-height: 25px;
    cursor: pointer;
    > .anticon {
      color: #000000a6;
    }
  }
  .ant-select {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
    color: #000000a6;
    font-size: 13px;
    font-variant: tabular-nums;
    line-height: 24px;
    list-style: none;
    font-feature-settings: 'tnum';
    position: relative;
    display: inline-block;
    cursor: pointer;
  }
  .ant-select-multiple .ant-select-selection-item-remove:hover {
    color: rgba(0, 0, 0, 0.75);
  }
}

.rc-virtual-list-scrollbar-thumb {
  background-color: #e5e6eb !important;
  &:hover {
    background-color: #c9cdd4 !important;
  }
}
</style>
