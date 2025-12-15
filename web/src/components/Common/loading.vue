<template>
  <div class="loading"><a-spin :size="spinSize" /></div>
</template>
<script lang="ts" setup>
import { PropType } from 'vue'

type Size = 'small' | 'default' | 'large'

defineProps({
  spinSize: {
    type: String as PropType<Size>,
    default: 'large',
  },
})
</script>
<style lang="less" scoped>
.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  padding: 20px 0;
}
</style>
