<template>
  <Tag class="literature-tag" :closable="closable" @close="handleClose"
    ><span class="text"><slot /></span
  ></Tag>
</template>
<script lang="ts" setup>
import { Tag } from 'ant-design-vue'

defineProps({
  closable: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits<{
  (event: 'close'): void
}>()

const handleClose = () => {
  emit('close')
}
</script>

<style lang="less" scoped>
.literature-tag {
  color: var(--site-theme-text-secondary);
  padding: 2px 4px;
  font-size: 12px;
  line-height: 18px;
  background: var(--site-theme-background-tertiary);
  border-radius: 2px;
  display: inline-block;
  white-space: normal;

  .text {
    word-break: break-all;
    overflow: hidden;
    text-overflow: ellipsis;
    display: inline-block;
    max-width: 155px;
    vertical-align: bottom;
  }
}
</style>
