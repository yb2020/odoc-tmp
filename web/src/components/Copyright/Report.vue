<template>
  <a-modal
    v-model:visible="visible"
    :title="$t('source.report.title')"
    width="478px"
    class="report"
    :closable="false"
  >
    <template #footer>
      <a-button
        key="back"
        @click="handleCancel"
      >
        {{
          $t('source.report.close')
        }}
      </a-button>
    </template>
    <div>
      <p>
        {{
          $t('source.report.content', { email: 'report_readpaper@idea.edu.cn' })
        }}
      </p>
      <a
        href="http://docs.qq.com/doc/DVXNmc0ZwaHhxdmlS"
        target="_blank"
        rel="noopener noreferrer"
      >{{ $t('source.report.tpl1') }}></a>
      <br>
      <a
        href="http://docs.qq.com/doc/DVUpqdFZlQUJpTGlr"
        target="_blank"
        rel="noopener noreferrer"
      >{{ $t('source.report.tpl2') }}></a>
    </div>
  </a-modal>
</template>
<script lang="ts" setup>
const visible = defineModel('visible', { default: false });

const handleCancel = () => {
  visible.value = false;
};
</script>

<style lang="less" scoped>
.report {
  p {
    height: 50px;
    font-size: 15px;
    font-weight: 400;
    line-height: 25px;
    margin-bottom: 20px;
  }
  a {
    font-size: 15px;
    color: #3b9aeb;
    line-height: 25px;
  }
}
</style>
