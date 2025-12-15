<template>
  <div
    ref="scrollElement"
    :class="{ more: total === 0 }"
  >
    <slot
      v-if="total === 0"
      name="empty"
    >
      <Empty />
    </slot>
    <template v-else>
      <slot />
      <LoadingOutlined v-show="loading && !error" />
      <div
        v-if="error"
        class="center rp-pointer"
        @click="loadData(true)"
      >
        <ReloadOutlined /> {{ error }}
      </div>
      <div
        v-else-if="!hasmore"
        class="center"
      >
        <slot
          v-if="$slots.nomore"
          name="nomore"
        />
        <span v-else>{{ $t('message.allData') }}</span>
      </div>
    </template>
  </div>
</template>
<script lang="ts" setup>
import { onMounted, onUnmounted, ref } from 'vue';
import { Empty } from 'ant-design-vue';
import lodash from 'lodash';
import { LoadingOutlined, ReloadOutlined } from '@ant-design/icons-vue';

const props = defineProps<{
  loading?: boolean;
  total?: number;
  error?: string;
  hasmore?: boolean;
}>();

const emit = defineEmits<{
  (e: 'scrollLoad', isRetry?: boolean): void;
}>();

const loadData = (fromErrorClick?: boolean) => {
  emit('scrollLoad', fromErrorClick);
};

const scrollLoad = lodash.debounce((e) => {
  if (!props.hasmore || props.error) {
    return;
  }
  // 总高度
  const scrollHeight = e.target.scrollHeight;
  // 滚动距离
  const scrolleTop = e.target.scrollTop;
  // 窗口高度
  const offsetHeight = e.target.offsetHeight;
  if (scrollHeight - scrolleTop - offsetHeight < 200) {
    loadData(false);
  }
}, 300);

const scrollElement = ref<HTMLElement>();

onMounted(() => {
  scrollElement.value?.addEventListener('scroll', scrollLoad, false);
});

onUnmounted(() => {
  scrollElement.value?.removeEventListener('scroll', scrollLoad, false);
});
</script>
<style lang="less" scoped>
.more {
  height: 85%;
}
.center {
  color: #a6a4a1;
  text-align: center;
  font-size: 12px;
  line-height: 20px;
  margin: 40px 0;
}
</style>
