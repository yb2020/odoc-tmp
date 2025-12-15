<template>
  <div
    v-if="!userStore.isLogin() && !isInClient"
    class="h-full flex flex-col items-center justify-center text-rp-neutral-6"
  >
    <slot name="blank" />
    <a-button
      type="primary"
      @click="userStore.openLogin()"
    >
      {{
        $t('common.tips.loginFirst')
      }}
    </a-button>
  </div>
  <slot v-else />
</template>

<script setup lang="ts">
import { useUserStore } from '@common/stores/user';

defineProps<{
  isInClient?: boolean;
}>();

const userStore = useUserStore();
</script>
