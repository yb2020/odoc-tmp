<template>
  <a
    target="_blank"
    href="/user/message"
  >
    <span class="rp-pointer">
      <Badge
        v-if="count > 0"
        :count="count > 99 ? '99+' : count"
      >
        <BellOutlined type="bell" />
      </Badge>

      <BellOutlined
        v-else
        type="bell"
      />
    </span>
  </a>
</template>
<script lang="ts" setup>
import { ref } from 'vue';
import { Badge } from 'ant-design-vue';
import { BellOutlined } from '@ant-design/icons-vue';

const count = ref(0);
</script>
<style lang="less" scoped>
:deep(.ant-scroll-number) {
  right: -10px;
  padding: 0 4px;
}
:deep(.anticon) {
  font-size: 20px;
  vertical-align: text-top;
}
a {
  color: inherit;
}
a:hover {
  color: #fff;
}
</style>
