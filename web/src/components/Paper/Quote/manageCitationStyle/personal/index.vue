<template>
  <div
    v-loading="loading"
    class="personal-style-wrap"
    data-loading-size="default"
  >
    <div v-if="fetchState && fetchState.error" class="error" @click="fetch">
      <ReloadOutlined /> {{ fetchState.error.message }}
    </div>

    <div v-else-if="styleList">
      <div class="title">
        {{ $t('home.paper.citeManage.currentStyle') }}（{{ styleList.length }}）
      </div>
      <div class="tip">{{ $t('home.paper.citeManage.tipDisplayLimit') }}</div>
      <div class="style-list-wrap">
        <draggable
          v-model="styleList"
          animation="200"
          group="personal-style-list"
          ghost-class="personal-style-list-dragging-ghost"
          @change="handleDragChange"
        >
          <transition-group>
            <Item
              v-for="(item, index) in styleList"
              :key="item.id"
              :data="item"
              :index="index"
              :length="styleList.length"
              @changeStyleSuccess="handleChangeStyleSuccess"
            ></Item>
          </transition-group>
        </draggable>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { defineComponent, ref, computed } from 'vue'
import draggable from 'vuedraggable'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { MyCslItem } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/CSL'
import Item from './item.vue'
import { sortCsl, getMyCslList } from '@/common/src/api/citation'
import useFetch from '@/common/src/hooks/useFetch'

export default defineComponent({
  components: {
    draggable,
    Item,
  },
  fetchOnServer: false,
  setup(_, { emit }) {
    const styleList = ref<MyCslItem[]>()

    const { fetchState, fetch } = useFetch(async () => {
      const res = await getMyCslList()

      styleList.value = res.list || []
    })

    const loading = computed(
      () => (!styleList.value || fetchState.pending) && !fetchState.error
    )

    const handleDragChange = async () => {
      if (!styleList.value) return

      try {
        await sortCsl({
          cslIds: styleList.value.map((item) => item.id),
        })
      } catch (error) {}
    }

    const handleChangeStyleSuccess = (id: string) => {
      fetch()
      emit('changeCslList', id)
    }

    return {
      styleList,
      fetchState,
      fetch,
      loading,
      handleDragChange,
      handleChangeStyleSuccess,
    }
  },
})
</script>
<style lang="less" scoped>
@import '../../scrollbar.less';

.personal-style-wrap {
  position: relative;
  width: 180px;
  padding: 12px 16px 0;
  border: 1px solid #e5e6eb;
  .title {
    font-family: 'Noto Sans SC';
    font-style: normal;
    font-weight: 400;
    font-size: 14px;
    line-height: 22px;
    color: #1d2229;
  }
  .tip {
    font-family: 'Noto Sans SC';
    font-style: normal;
    font-weight: 400;
    font-size: 12px;
    line-height: 18px;
    color: #86919c;
    margin: 4px 0 14px;
  }
  .style-list-wrap {
    height: 242px;
    overflow-y: auto;
    margin: 0 -16px;
  }
  .error {
    padding-top: 100px;
    cursor: pointer;
  }
}
</style>
