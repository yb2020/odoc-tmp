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
          :key="item.paperId"
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
import { computed, onMounted, onUnmounted, watch } from "vue";
import useFetchList, { QueryParam } from "~/src/hooks/useFetchList";
import {
  ReferenceSortType,
  GetCitationResponse,
  CitationInfo,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse';
import { getCitationFinal } from "~/src/api/material";
import Item from "./Item.vue";
import ErrorVue from '../Error.vue'
import EmptyVue from "../Empty.vue";
import { RefCiteTabButtonType as TabButtonType } from './type'

const props = defineProps<{ sort: ReferenceSortType, paperId: string }>();

const emit = defineEmits<{ (event: 'updateTotal', totals: Partial<Record<TabButtonType, number>>): void }>();

const pageSize = 5;

const fetchInstance = getCitationFinal()

const { fetchState, currentPageList, fetch, total, currentPage } = useFetchList<
  CitationInfo
>(
  async (queryParam: QueryParam) => {
    const res: GetCitationResponse =
      await fetchInstance.request({
        sortType: props.sort,
        paperId: props.paperId,
        pageReq: {
          pageSize: queryParam.pageSize,
          pageNum: queryParam.currentPage,
        },
      })


    if (queryParam.currentPage === 1) {
      emit('updateTotal', {
        [TabButtonType.CITATION]: res.pageResp?.total || res.citationCount || 0,
      })
    }


    return {
      list: res.citationInfoList,
      total: res.pageResp?.total || 0,
    };
  },
  pageSize,
  true
);


onUnmounted(() => {
  fetchInstance.cancle()
})

const loading = computed(() => (total.value === -1 || fetchState.pending) && !fetchState.error);

const handleChange = async (page: number) => {
  await fetch(page);
};

watch(() => props.sort, () => {
  fetch(1)
})

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
