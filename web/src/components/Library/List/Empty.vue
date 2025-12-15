<template>
  <Empty
    :class="['empty', relative ? 'relative' : 'absolute']"
    :image="imgUrl"
    :image-style="{
      height: '55px',
    }"
  >
    <template #description>
      <div class="text">
        <div class="text">{{ desc }}</div>
        <Button v-if="button" class="btn" type="primary" @click="handleClick">{{
          button
        }}</Button>
      </div>
    </template>
  </Empty>
</template>
<script lang="ts" setup>
import { Empty, Button } from 'ant-design-vue'
import defaultImgUrl from '@/assets/images/empty.png'

defineProps({
  imgUrl: {
    type: String,
    default: defaultImgUrl,
  },
  desc: {
    type: String,
    default: '暂无数据',
  },
  button: {
    type: String,
    default: '',
  },
  relative: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['buttonClick'])

const handleClick = () => {
  emit('buttonClick')
}
</script>
<style lang="less" scoped>
.empty {
  left: 50%;
  transform: translateX(-50%);
  padding: 12px 0;
  .text {
    color: #a19f9d;
    font-size: 13px;
    margin-top: 16px;
    line-height: 24px;
  }
  .btn {
    margin-top: 8px;
  }
}
.relative {
  position: relative;
  margin-top: 100px;
}
.absolute {
  position: absolute;
  top: 200px;
}
</style>
