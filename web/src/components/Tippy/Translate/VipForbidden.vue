<template>
  <!-- 移除会员权限检查，改为积分制，始终显示内容 -->
  <slot />
</template>
<script lang="ts" setup>
// 简化组件，只保留必要的属性
// 不再需要会员相关的导入和属性
import { TranslateTabKey } from '~/src/stores/translateStore';

defineProps<{
  currentTab: TranslateTabKey;
}>();

// 移除未使用的translateStore变量
</script>
