<template>
  <div class="private-comment">
    <i
      class="aiknowledge-icon icon-lock-fill"
      aria-hidden="true"
    />
    {{ $t('message.privateNoteTip') }}
  </div>
</template>
<script lang="ts" setup></script>
<style lang="less" scoped>
.private-comment {
  width: 100%;
  height: 100%;
  z-index: 3;

  background: rgba(39, 43, 46, 0.65);
  backdrop-filter: blur(4px);

  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;

  position: absolute;
  left: 0;
  bottom: 0;

  overflow: hidden;

  font-size: 14px;
  font-weight: 600;
  color: #a19f9d;
  line-height: 24px;

  .aiknowledge-icon {
    font-size: 24px;
    color: #a19f9d;
    margin-bottom: 16px;
  }
}
</style>
