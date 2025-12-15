<script setup lang="ts"></script>
<template>
  <div class="border rounded-xl border-rp-neutral-3 flex flex-col bg-white">
    <div class="border-b border-rp-neutral-3">
      <slot name="title" />
    </div>
    <slot />
  </div>
</template>
