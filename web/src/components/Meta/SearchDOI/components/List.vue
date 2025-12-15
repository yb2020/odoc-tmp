<template>
  <div class="update-list">
    <div class="total">
      {{
        $t('meta.detail.totalSearchResults', list.length, {
          named: { total: list.length },
        })
      }}
    </div>
    <ul v-if="list.length" class="list">
      <li v-for="(item, idx) in list" :key="idx" @click="handleSelect(item)">
        {{ item.title }}
      </li>
    </ul>
    <div class="btns">
      <a-button class="btn" @click="handleBack">{{ $t('meta.back') }}</a-button>
    </div>
  </div>
</template>
<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { DocMetaInfoSimpleVo } from 'go-sea-proto/gen/ts/doc/CSL'

export default defineComponent({
  props: {
    list: {
      type: Array as PropType<DocMetaInfoSimpleVo[]>,
      default: () => [],
    },
  },
  setup(_, { emit }) {
    const handleBack = () => {
      emit('back')
    }
    const handleSelect = (result: DocMetaInfoSimpleVo) => {
      emit('select', result)
    }
    return {
      handleBack,
      handleSelect,
    }
  },
})
</script>
<style lang="less" scoped>
.update-list {
  display: flex;
  flex-direction: column;
  height: 100%;
  justify-content: space-between;
  .total {
    color: #4e5969;
  }

  .list {
    margin: 24px 0;
    overflow: auto;
    flex: 1;
    li {
      color: #1d2229;
      line-height: 22px;
      display: -webkit-box;
      overflow: hidden;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      cursor: pointer;
    }
    li + li {
      margin-top: 22px;
    }
  }
}

.btns {
  display: flex;
  justify-content: space-between;
}
.btn {
  background: #f0f2f5;
  border: none;
  color: #4e5969;
}
</style>
