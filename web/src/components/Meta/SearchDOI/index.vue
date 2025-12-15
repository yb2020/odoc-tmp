<template>
  <div
    class="update-doi-search"
    :style="step === Step.Result ? {} : { height: '500px' }"
  >
    <Search
      v-show="step === Step.Search"
      :placeholder="placeholder"
      :from="from"
      @search="onSearch"
      @fetch="onFetch"
    />
    <List
      v-show="step === Step.List"
      :list="list"
      @select="onSelect"
      @back="onBack(Step.Search)"
    />
    <Result
      v-if="step === Step.Result"
      :result="result!"
      :paper-id="paperId"
      :pdf-id="pdfId"
      @back="onBack(Step.List)"
      @update:success="$emit('update:success')"
    />
  </div>
</template>
<script lang="ts">
import { defineComponent, PropType, ref } from 'vue'
import { DocMetaInfoSimpleVo } from 'go-sea-proto/gen/ts/doc/CSL'
import Search from './components/Search.vue'
import List from './components/List.vue'
import Result from './components/Result.vue'
import { DocMetaInfoWithVenue, withVenue, reportSearchRenewClick } from '@/api/citation'

enum Step {
  Search = 1,
  List = 2,
  Result = 3,
}

export default defineComponent({
  components: { Search, List, Result },
  props: {
    placeholder: {
      type: String,
      default: '',
    },
    from: {
      type: String as PropType<'english' | 'chinese'>,
      default: 'english',
    },
    paperId: {
      type: String,
      default: '',
    },
    pdfId: {
      type: String,
      default: '',
    },
    pageType: {
      type: String,
      required: true,
    },
  },
  emits: ['update:success'],
  setup(props) {
    const step = ref(Step.Search)
    const list = ref<DocMetaInfoSimpleVo[]>([])
    const result = ref<DocMetaInfoWithVenue | null>(null)
    const onSearch = (results: DocMetaInfoSimpleVo[]) => {
      console.log(results)
      step.value = Step.List
      list.value = results
    }
    const onBack = (val: Step) => {
      step.value = val
    }
    const onSelect = (selected: DocMetaInfoSimpleVo) => {
      reportSearch('title')
      console.log(result)
      result.value = withVenue(selected)
      step.value = Step.Result
    }

    let searchContent = ''
    const onFetch = (keyword: string) => {
      searchContent = keyword
      reportSearch('input')
    }

    const reportSearch = (
      type: Parameters<typeof reportSearchRenewClick>[0]['type']
    ) => {
      reportSearchRenewClick({
        page_type: props.pageType,
        search_content: searchContent,
        language_type: props.from,
        type,
      })
    }

    return {
      list,
      result,
      onSearch,
      step,
      Step,
      onBack,
      onSelect,
      onFetch,
    }
  },
})
</script>
<style lang="less" scoped>
.update-doi-search {
  overflow: hidden;
}
</style>
