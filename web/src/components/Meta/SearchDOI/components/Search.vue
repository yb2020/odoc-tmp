<template>
  <div class="update-search">
    <a-input
      v-model:value="value"
      :placeholder="placeholder"
      @pressEnter="handleSearch"
    />
    <div class="btn">
      <a-button type="primary" block :loading="loading" @click="handleSearch">
        {{ $t('meta.search') }}
      </a-button>
    </div>
  </div>
</template>
<script lang="ts">
import { defineComponent, ref, PropType } from 'vue'
import { enDocCiteSearch, zhDocCiteSearch } from '@/api/citation'
import { message } from 'ant-design-vue'

export default defineComponent({
  layout: 'none',
  props: {
    placeholder: {
      type: String,
      default: '',
    },
    from: {
      type: String as PropType<'english' | 'chinese'>,
      default: 'english',
    },
  },
  emits: ['search', 'fetch'],
  setup(props, { emit }) {
    const value = ref('')

    const loading = ref(false)

    const handleSearch = async () => {
      if (loading.value) {
        return
      }
      const val = value.value.trim()
      if (!val) {
        message.error('请输入搜索内容')
        return
      }

      emit('fetch', val)

      loading.value = true
      try {
        const result =
          props.from === 'english'
            ? await enDocCiteSearch({
                searchContent: val,
              })
            : await zhDocCiteSearch({
                searchContent: val,
              })

        value.value = ''
        emit('search', result)
      } catch (error) {
        message.error((error as Error).message)
        console.log('searchError', typeof error)
      } finally {
        loading.value = false
      }
    }
    return {
      handleSearch,
      value,
      loading,
    }
  },
})
</script>
<style lang="less" scoped>
.update-search {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  .btn {
    padding: 40px 86px 0;
    width: 100%;
  }
}
.ant-input {
    color: var(--site-theme-text-primary);
}
.ant-input::placeholder {
    color: var(--site-theme-placeholder-color);
}
</style>
