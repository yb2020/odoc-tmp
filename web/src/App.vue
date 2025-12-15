<script setup lang="ts">
// This starter template is using Vue 3 <script setup> SFCs
// Check out https://v3.vuejs.org/api/sfc-script-setup.html#sfc-script-setup
import LimitDialog from '@/common/src/components/Pay/LimitDialog.vue';
</script>

<template>
  <router-view />
  <LimitDialog />
</template>

<style lang="less">
html,
body,
#app {
  height: 100%;
  margin: 0;
  .wang_plugin_tag {
    margin: 0 3px;
    color: #1f71e0;
  }
}
</style>
