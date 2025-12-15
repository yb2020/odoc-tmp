<template>
  <Exception
    :code="query.status"
    :message="query.message"
    :url="query.source"
  />
</template>
<script lang="ts" setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import Exception from './exception.vue'

const route = useRoute();

const query = computed(() => route.query as Record<string, string>);
</script>
