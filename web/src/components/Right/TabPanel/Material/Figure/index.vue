<template>
  <a-spin :spinning="total === -1 && !error">
    <div class="figure-material">
      <div class="tab-title">
        {{ $t('viewer.figuresAndTables') }}
      </div>
      <ErrorVue
        v-if="error"
        :error="error.message"
        @redo="handleRedo"
      />
      <Empty v-else-if="total === 0" />
      <div
        v-else
        class="list"
      >
        <Item
          v-for="(item, index) in list"
          :key="index"
          :item="item"
          :pdfId="pdfId"
        />
      </div>
    </div>
  </a-spin>
</template>
<script lang="ts" setup>
import { computed } from 'vue';
import { useStore } from '~/src/store';
import { ParseActionTypes } from '~/src/store/parse';
import Empty from '../Empty.vue';
import ErrorVue from '../Error.vue';
import Item from './Item.vue';

const props = defineProps<{ pdfId: string }>();

const store = useStore();

const list = computed(
  () => store.state.parse.pdfFigureAndTableInfos[props.pdfId]?.list
);

const error = computed(
  () => store.state.parse.pdfFigureAndTableInfos[props.pdfId]?.error
);

const total = computed(() => (list.value ? list.value.length : -1));

const handleRedo = () => {
  store.dispatch(
    `parse/${ParseActionTypes.GET_PDF_FIGURE_AND_TABLE_INFOS}`,
    props.pdfId
  );
};
</script>
<style scoped lang="less">
.figure-material {
  font-size: 13px;
  line-height: 18px;
  margin: 10px 0;
  // background: #282b2e;
  position: relative;
  padding-bottom: 10px;
  // margin-bottom: 70px;
  color: var(--site-theme-pdf-panel-text);
  min-height: 120px;
  .loading {
    height: 50px;
    position: relative;
  }
  .tab-title {
    padding: 10px 0;
    // border: 1px solid #e4e7ed;
    font-size: 14px;
    line-height: 26px;
    font-weight: 600;
  }

  .pagination {
    text-align: right;
    padding-bottom: 20px;
  }
}
</style>
