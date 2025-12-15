<template>
  <a-spin :spinning="loading">
    <ErrorVue
      v-if="fetchState.error"
      :error="fetchState.error.message"
      @redo="fetch(currentPage)"
    />
    <EmptyVue v-else-if="total === 0" />
    <div
      v-else
      class="ref-wrapper"
    >
      <div class="list">
        <Item
          v-for="(item, index) in currentPageList"
          :key="item.refIdx"
          :title="item.title"
          :index="(currentPage - 1) * pageSize + index + 1"
          :paper-id="item.paperId"
        />
      </div>
      <div
        v-if="total > 0"
        class="pagination"
      >
        <a-pagination
          simple
          size="small"
          :pageSize="5"
          :current="currentPage"
          :default-current="1"
          :total="total"
          @change="handleChange"
        />
      </div>
    </div>
  </a-spin>
</template>
<script setup lang="ts">
import { computed, watch } from 'vue';
import useFetchList, { QueryParam } from '~/src/hooks/useFetchList';
import {
  ReferenceInfo,
  ReferenceSortType,
  GetReferenceResponse,
} from 'go-sea-proto/gen/ts/pdf/PdfParse';

import { getReferenceFinal } from '~/src/api/material';
import Item from './Item.vue';
import ErrorVue from '../Error.vue';
import EmptyVue from '../Empty.vue';
import { RefCiteTabButtonType as TabButtonType } from './type';

const props =
  defineProps<{ sort: ReferenceSortType; pdfId: string; paperId: string }>();

const emit = defineEmits<{
  (event: 'updateTotal', totals: Partial<Record<TabButtonType, number>>): void;
}>();

const pageSize = 5;

const { fetchState, currentPageList, fetch, total, currentPage } =
  useFetchList<ReferenceInfo>(
    async (queryParam: QueryParam) => {
      const res: GetReferenceResponse = await getReferenceFinal({
        sortType: props.sort,
        pdfId: props.pdfId,
        paperId: props.paperId,
        pageReq: {
          pageSize: queryParam.pageSize,
          pageNum: queryParam.currentPage,
        },
      });

      if (queryParam.currentPage === 1) {
        emit('updateTotal', {
          [TabButtonType.REFERENCE]:
            res.pageResp?.total || res.referenceCount || 0,
          [TabButtonType.CITATION]:
            res.citationCount >= 0 ? res.citationCount : -1,
        });
      }

      return {
        list: res.referenceInfoList,
        total: res.pageResp?.total || 0,
      };
    },
    pageSize,
    true
  );

const loading = computed(
  () => (total.value === -1 || fetchState.pending) && !fetchState.error
);

const handleChange = async (page: number) => {
  await fetch(page);
};

watch(
  () => props.sort,
  () => {
    fetch(1);
  }
);
</script>

<style scoped lang="less">
.list {
  margin: 15px 0 10px;
  min-height: 200px;
}
.pagination {
  text-align: right;
  margin-right: 24px;
  padding-bottom: 10px;
  .ant-pagination {
    color: #4e5969;
  }
  :deep(.ant-pagination-simple) .ant-pagination-simple-pager input {
    background: rgba(255, 255, 255, 0.15);
    border: none;
  }
  :deep(.ant-pagination-item-link) {
    color: #4e5969 !important;
  }
}
</style>
