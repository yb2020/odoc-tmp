<template>
  <GlobalHeader>
    <template #title>
      <Title v-if="!IS_MOBILE && !IS_ELECTRON_MODE" />
    </template>
  </GlobalHeader>
</template>

<script lang="ts" setup>
import { IS_ELECTRON_MODE, IS_MOBILE } from '~/src/util/env';
import Title from './components/Title.vue';
import GlobalHeader from '../GlobalHeader/index.vue';
</script>
