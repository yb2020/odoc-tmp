<template>
  <a-result
    class="absolute z-10 top-0 left-0 w-full h-full flex flex-col justify-center bg-white"
  >
    <template #title>
      <p class="text-rp-neutral-8 text-sm">
        请在打开的独立窗口中继续编辑<br>关闭独立窗口后会自动恢复可编辑状态
      </p>
    </template>
    <template #icon>
      <img
        src="@/assets/images/locked.svg"
        alt="Locked"
      >
    </template>
  </a-result>
</template>

<script setup lang="ts"></script>
