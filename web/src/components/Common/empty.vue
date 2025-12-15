<template>
  <Empty
    :class="['empty', { 'dark-style': dark }]"
    :image="imgUrl"
    :image-style="{
      height: size === 'small' ? '40px' : '60px',
    }"
  >
    <template #description>
      <span class="text">
        {{ content || $t('home.global.noData') }}
      </span>
    </template>
  </Empty>
</template>
<script lang="ts">
import { Empty } from 'ant-design-vue'
import { defineComponent, PropType } from 'vue'
import imgUrl from '@/assets/images/empty.png'

export default defineComponent({
  props: {
    dark: {
      type: Boolean,
      default: false,
    },
    size: {
      type: String as PropType<'normal' | 'small'>,
      default: 'normal',
    },
    content: {
      type: String,
      default: '',
    },
  },
  setup() {
    return {
      imgUrl,
    }
  },
  components: { Empty },
})
</script>
<style lang="less" scoped>
.empty {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 12px 0;
  margin: 0;
  height: 100%;
  flex-direction: column;
  .text {
    color: #a19f9d;
    font-size: 13px;
  }
}
.dark-style {
  background: #f9f9f9;
  .text {
    color: rgba(0, 0, 0, 0.34);
  }
}
</style>
