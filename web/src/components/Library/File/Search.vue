<template>
  <div class="uploader-search-input">
    <a-input-search
      v-model:value="keyword"
      class="input"
      :placeholder="placeholder || $t('home.header.search.placeHolder')"
      allow-clear
      @change="change"
      @search="search(1)"
    />
    <a-spin v-if="searchResult" :spinning="loading" class="search">
      <div v-for="item in searchList" :key="item.id" class="item">
        <PaperItem
          :paper-data="item"
          size="small"
          :is-collect-paper="false"
          :is-show-venue-tagging="false"
          :collect-limit-dialog-report-params="reportParams"
          class="paper"
        />
        <slot
          name="addButton"
          :pdf-id="item.pdfId"
          :paper-id="item.id"
          :is-collected="item.isCollected"
        />
      </div>
      <a-pagination
        v-if="total > 0"
        class="pagination"
        show-quick-jumper
        size="small"
        :default-current="1"
        :total="total"
        :page-size="5"
        @change="search"
      />
      <Empty v-else-if="total === 0" :dark="false" size="small" />
    </a-spin>
    <div v-else class="query">
      <div v-for="item in queryResultList" :key="item.id" class="item">
        <PaperItem
          :paper-data="item"
          size="small"
          :is-collect-paper="false"
          :is-show-venue-tagging="false"
          :collect-limit-dialog-report-params="reportParams"
          class="paper"
        />
        <slot
          name="addButton"
          :paper-id="item.id"
          :is-collected="item.isCollected"
        />
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { computed, watch, ref } from 'vue'
import useSearchRecommend from '@/hooks/useSearchRecommend'
import { fetchSearchPaperResult } from '@/common/src/api/search'
import Empty from '@/components/Common/empty.vue'
import { $PaperDetail } from '@/common/src/api/paper'
import PaperItem from '@/components/Paper/PaperItem.vue'
import { ElementName, PageType } from '@/utils/report'

const props = defineProps({
  fromCopilot: {
    type: Boolean,
    default: false,
  },
  placeholder: {
    type: String,
    default: '',
  },
})

const reportParams = {
  page_type: PageType.library,
  element_name: ElementName.upperCollectionPopup,
}

const emit = defineEmits(['onSearch'])

const { keyword, fetchQuery, items } = useSearchRecommend('')
const queryResultList = computed(() => {
  if (!String(keyword.value).trim()) {
    return []
  }
  return items.value.map((item) => {
    return {
      id: item.id,
      title: item.title,
      isCollected: item.isCollected,
      authorList: item.authors.map((author) => {
        return { id: '', name: author }
      }),
    } as any
  })
})
watch(keyword, (value) => {
  emit('onSearch', !!String(value).trim())
})
const searchResult = ref<boolean>(false)
const total = ref(0)
const searchList = ref<$PaperDetail[]>([])
const loading = ref<boolean>(false)
const change = () => {
  if (props.fromCopilot) {
    if (!String(keyword.value).trim()) {
      searchResult.value = false
      searchList.value = []
      loading.value = false
      total.value = -1
    }
    return
  }
  searchResult.value = false
  searchList.value = []
  fetchQuery()
}
const search = async (pageNumber: number) => {
  items.value = []
  searchResult.value = true
  if (!String(keyword.value).trim() || loading.value) {
    return
  }
  try {
    loading.value = true
    if (pageNumber === 1) {
      total.value = -1
    }
    const params: Parameters<typeof fetchSearchPaperResult>['1'] = {
      keywords: keyword.value,
      page: pageNumber,
      pageSize: 5,
    }
    if (props.fromCopilot) {
      params.searchHasPublicPdf = true
      params.searchType = 1
    }
    const res = await fetchSearchPaperResult(
      params
      // context.$sentry // TODO
    )
    searchList.value = res.list
    total.value = res.total
  } catch (error) {}
  loading.value = false
}
</script>
<style lang="less" scoped>
@import '@/assets/css/antd.less';
@import '@/components/Library/Menu/Menu.less';

.uploader-search-input {
  .input {
    margin-bottom: 24px;
  }
  .item {
    display: flex;
    justify-content: space-between;
    .paper {
      margin-right: 36px;
      overflow: hidden;
    }
    .btn {
      padding: 0 16px;
    }
  }
  .item + .item {
    margin-top: 14px;
  }

  .search {
    position: relative;
    min-height: 40px;
    max-height: 507px;
    overflow-y: auto;
    width: 100%;
  }

  .pagination {
    margin-top: 16px;
  }
}
</style>
