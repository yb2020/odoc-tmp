<template>
  <div class="flex items-center justify-center h-full text-rp-neutral-10">
    <div class="text-sm mt-4 text-center">
      <div class="text-base mb-1">
        {{ $t('translate.emptyTip') }}
      </div>
      <p>💡{{ $t('translate.translateTip', { cmd: CtrlOrCommand }) }}</p>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { getPlatformKey } from '~/src/store/shortcuts';

const CtrlOrCommand = getPlatformKey() === 'win32' ? 'Ctrl' : 'Command';
</script>


<style lang="less" scoped>
.flex {
  display: flex;
  color: var(--site-theme-pdf-panel-text);
}
</style>
