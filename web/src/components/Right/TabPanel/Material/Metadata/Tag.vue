<template>
  <Tag
    class="literature-tag"
    :class="[dark ? 'literature-tag-dark' : 'literature-tag-light']"
    :closable="props.closable"
    :style="{ fontSize: `${props.fontSize}px` }"
    @close="handleClose"
  >
    <span class="text">{{ props.text }}</span>
  </Tag>
</template>
<script lang="ts" setup>
import { Tag } from 'ant-design-vue';

const props = defineProps({
  closable: {
    type: Boolean,
    default: true,
  },
  text: {
    type: String,
    default: '',
  },
  dark: {
    default: false,
  },
  fontSize: {
    type: Number,
    default: 12,
  },
});

const emit = defineEmits(['close']);

const handleClose = () => {
  emit('close');
};
</script>

<style lang="less">
.literature-tag.ant-tag {
  padding: 2px 4px;
  line-height: 18px;
  border-radius: 2px;
  display: inline-block;
  white-space: normal;
  border: 0;
  &.literature-tag-light {
    color: #000000a6;
    background: #f0f2f5;
  }
  &.literature-tag-dark {
    color: #ffffffd9;
    background-color: #ffffff26;
  }
  .text {
    word-break: break-all;
    overflow: hidden;
    text-overflow: ellipsis;
    display: inline-block;
    // max-width: 155px;
    vertical-align: bottom;
  }
}
</style>
