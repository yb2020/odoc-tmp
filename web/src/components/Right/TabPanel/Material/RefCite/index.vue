<template>
  <div
    ref="refCiteRef"
    class="metadata-field-value"
    style="color: rgba(255, 255, 255, 0.85)"
  >
    <span
      class="ref-cite-open"
      @click="handleTab(TabButtonType.REFERENCE)"
    >
      ref
      {{
        totalValue[TabButtonType.REFERENCE] >= 0
          ? ` （${totalValue[TabButtonType.REFERENCE]}）`
          : ''
      }}
    </span>
    <!-- /&nbsp;
    <span
      class="ref-cite-open"
      @click="handleTab(TabButtonType.CITATION)"
    >
      cite
      {{
        totalValue[TabButtonType.CITATION] >= 0
          ? ` （${totalValue[TabButtonType.CITATION]}）`
          : ''
      }}
    </span> -->
  </div>
  <Tippy
    v-if="refCiteRef"
    ref="tippyRef"
    :trigger-ele="refCiteRef"
    :max-width="336"
    placement="left-start"
    trigger="click"
    :offset="[0, 76]"
  >
    <div class="ref-cite-container">
      <div
        class="ref-cite-close"
        @click="tippyRef?.hide()"
      >
        <CloseOutlined />
      </div>
      <div class="tab-buttons">
        <div
          :class="['btn', { current: tabValue === TabButtonType.REFERENCE }]"
          @click="handleTab(TabButtonType.REFERENCE)"
        >
          {{ $t('viewer.references') }}
          <span v-if="totalValue[TabButtonType.REFERENCE] >= 0">· {{ totalValue[TabButtonType.REFERENCE] }}</span>
        </div>
        <!-- <div
          :class="['btn', { current: tabValue === TabButtonType.CITATION }]"
          @click="handleTab(TabButtonType.CITATION)"
        >
          {{ $t('viewer.citations') }}
          <span v-if="totalValue[TabButtonType.CITATION] >= 0">· {{ totalValue[TabButtonType.CITATION] }}</span>
        </div> -->
      </div>
      <div class="sort">
        <span class="sort-label">{{ $t('viewer.sort') }}: &nbsp;</span>
        <select
          :value="selectValue"
          class="select"
          @change="handleSort"
        >
          <option :value="ReferenceSortType.DEFAULT">
            {{ $t('viewer.default') }}
          </option>
          <!-- <option :value="ReferenceSortType.CITATION">
            {{ $t('viewer.citationCount') }}
          </option> -->
          <option :value="ReferenceSortType.PUBLISH_DATE">
            {{ $t('viewer.publishDate') }}
          </option>
        </select>
      </div>
      <keep-alive>
        <RefVue
          v-if="tabValue === TabButtonType.REFERENCE"
          :sort="selectValue"
          :currentTab="tabValue"
          :pdf-id="pdfId"
          :paper-id="paperId"
          @update-total="onUpdateTotal"
        />
        <CiteVue
          v-else
          :sort="selectValue"
          :currentTab="tabValue"
          :pdf-id="pdfId"
          :paper-id="paperId"
          @update-total="onUpdateTotal"
        />
      </keep-alive>
    </div>
  </Tippy>
</template>
<script setup lang="ts">
import { ref } from 'vue';
import { CloseOutlined } from '@ant-design/icons-vue';
import { ReferenceSortType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse';
import Tippy from '~/src/components/Tippy/index.vue';
import { RefCiteTabButtonType as TabButtonType } from './type';
import RefVue from './Ref.vue';
import CiteVue from './Cite.vue';
import { isNumber } from '@vueuse/core';

defineProps<{ paperId: string; pdfId: string }>();

const emit = defineEmits(['close']);

const refCiteRef = ref<HTMLDivElement | null>(null);
const tippyRef = ref();
const selectValue = ref<ReferenceSortType>(ReferenceSortType.DEFAULT);
const handleSort = (e: any) => {
  selectValue.value = parseInt(e.target?.value) || 0;
};

const tabValue = ref<TabButtonType>(TabButtonType.REFERENCE);
const handleTab = (value: TabButtonType) => {
  tabValue.value = value;
};

const totalValue = ref<{ [_key in TabButtonType]: number }>(
  {} as { [_key in TabButtonType]: number }
);

const onUpdateTotal = (totals: Partial<Record<TabButtonType, number>>) => {
  if (isNumber(totals[TabButtonType.REFERENCE])) {
    totalValue.value[TabButtonType.REFERENCE] = totals[
      TabButtonType.REFERENCE
    ] as number;
  }
  if (isNumber(totals[TabButtonType.CITATION])) {
    totalValue.value[TabButtonType.CITATION] = totals[
      TabButtonType.CITATION
    ] as number;
  }
};
</script>

<style scoped lang="less">
@import '../Metadata/style.less';

[data-theme="dark"] {
  .ref-cite-container {
    background-color: #222326 !important;
  }
}

.ref-cite-container {
  min-width: 336px;
  background-color: #fff;
  padding-top: 16px;
  padding-bottom: 14px;
}
.tab-buttons {
  display: flex;
  justify-content: flex-start;
  margin-bottom: 18px;
  margin-left: 14px;
  .btn {
    background: rgba(255, 255, 255, 0.05);
    color: #86919c;
    padding: 3px 10px;
    cursor: pointer;
    &:first-child {
      border-radius: 4px 0 0 4px;
    }
    &:last-child {
      border-radius: 0 4px 4px 0;
    }
    &.current {
      color: #1d2229;
    }
  }
}

.sort {
  margin-left: 24px;
  color: #86919c;
  .select {
    color: #4e5969;
    background-color: #ffffff;
    border: none;
    outline: none;
    cursor: pointer;
  }
  .arrow {
    opacity: 0.3;
    font-size: 10px;
    display: inline-block;
    cursor: pointer;
  }
}

.ref-cite-open {
  color: #1f65c4;
  cursor: pointer;
  &:hover {
    color: #4387d9;
  }
}
.ref-cite-close {
  position: absolute;
  right: 12px;
  top: 12px;
  height: 30px;
  width: 30px;
  display: flex;
  color: #86919c;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}
</style>
