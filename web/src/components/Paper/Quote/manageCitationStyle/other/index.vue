<template>
  <div class="other-style-wrap">
    <div class="header">
      <span class="title">{{ $t('home.paper.citeManage.otherStyle') }}</span>
      <Input
        v-model="keyword"
        :placeholder="$t('home.paper.citeManage.tipSearch')"
        size="small"
        class="input"
        @change="handleSearchChange"
      >
        <template #prefix>
          <SearchOutlined class="search-icon" />
        </template>
      </Input>
    </div>
    <div class="other-style-content">
      <ScrollList
        :loading="loading"
        :error="error"
        :total="total"
        :hasmore="hasmore"
        scroll-element=".other-style-content"
        @scrollLoad="scrollLoad"
      >
        <template #list>
          <Item
            v-for="item in list"
            :key="item.id"
            :data="item"
            :added-count="addedCount"
            @changeStyleSuccess="handleChangeStyleSuccess"
          ></Item>
        </template>

        <template #empty>
          <div class="empty-result">
            {{ $t('home.paper.citeManage.tipBlank') }}
          </div>
        </template>
      </ScrollList>
    </div>
  </div>
</template>
<script lang="ts">
import { defineComponent, ref } from 'vue'
import { CslItem } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/CSL'
import Item from './item.vue'
import useFetchList from '@common/hooks/useFetchList'
import { getCslList } from '@common/api/citation'
import ScrollList from '@common/components/ScrollList.vue'
import { Input } from 'ant-design-vue'
import { SearchOutlined } from '@ant-design/icons-vue'

export default defineComponent({
  components: {
    ScrollList,
    Item,
    Input,
    SearchOutlined,
  },
  props: {
    addedCount: {
      type: Number,
      default: 0,
    },
  },
  fetchOnServer: false,
  setup(_, { emit }) {
    const keyword = ref<string>('')

    const { list, total, loading, fetch, error, reload, currentPage, hasmore } =
      useFetchList<CslItem>(async (queryParam) => {
        const res = await getCslList({
          pageSize: queryParam.pageSize,
          currentPage: queryParam.currentPage,
          searchContent: keyword.value,
        })

        return {
          list: res.list || [],
          total: res.total || 0,
        }
      }, 20)

    const scrollLoad = async () => {
      if (loading.value) {
        return
      }
      await fetch()
    }

    const handleSearchChange = () => {
      reload()
    }

    const handleChangeAddedFlag = (id: string) => {
      const findIndex = list.value.findIndex((item) => item.id === id)

      if (findIndex !== -1)
        list.value[findIndex].addedFlag = !list.value[findIndex].addedFlag
    }

    const handleChangeStyleSuccess = (id: string) => {
      handleChangeAddedFlag(id)

      emit('fetchMyCslList')
    }
    return {
      keyword,
      list,
      total,
      loading,
      fetch,
      error,
      reload,
      currentPage,
      hasmore,
      scrollLoad,
      handleSearchChange,
      handleChangeStyleSuccess,
      handleChangeAddedFlag,
    }
  },
})
</script>
<style lang="less" scoped>
@import '../../scrollbar.less';

.other-style-wrap {
  width: 340px;
  background: #f7f8fa;
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 24px 8px 16px;
    .title {
      font-family: 'Noto Sans SC';
      font-style: normal;
      font-weight: 400;
      font-size: 14px;
      line-height: 22px;
      color: #1d2229;
    }
    .input {
      width: 160px;
      .search-icon {
        color: #a8afba;
      }
    }
  }
  .other-style-content {
    height: 263px;
    overflow-y: auto;
    .empty-result {
      font-family: 'Noto Sans SC';
      font-style: normal;
      font-weight: 400;
      font-size: 13px;
      line-height: 20px;
      color: #a8afba;
      text-align: center;
      padding-top: 100px;
    }
  }
}
/deep/.ant-input {
  font-size: 12px;
}
/deep/.ant-input:placeholder-shown {
  font-size: 12px;
}
</style>
