<script setup lang="ts">
import { computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';

const router = useRouter();
const route = useRoute();

const message = computed(() => {
  return route.query.message as string;
});

const goBack = () => {
  router.replace({ path: '/' });
};
</script>

<template>
  <a-result
    class="mt-10"
    status="error"
    title="Error"
    :sub-title="message"
  >
    <template #extra>
      <a-button
        type="primary"
        @click="goBack"
      >
        Back Home
      </a-button>
    </template>
  </a-result>
</template>
