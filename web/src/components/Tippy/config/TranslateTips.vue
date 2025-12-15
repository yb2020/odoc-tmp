<template>
  <div class="tips-container">
    <div v-if="selectedType === 'other'">
      <div class="tips-title">
        <span>
          <p>
            ReadPaper固定有道、百度、IDEA三个内置翻译渠道,除此之外,更推出用户添加自定义翻译接口,满足个性化的翻译需求、从而得到更高质量的翻译结果。
          </p>
        </span>
      </div>

      <div class="tips-button-container">
        <a
          class="tips-button"
          @click="toConfigPage"
        > 如何添加? </a>
      </div>
    </div>
    <TencentTranslate v-else-if="selectedType === 'tencent'" />
    <AlibabaTranslate v-else-if="selectedType === 'ali'" />
    <GoogleTranslate v-else-if="selectedType === 'google'" />
    <DeeplTranslate v-else-if="selectedType === 'deepl'" />
  </div>
</template>

<script lang="ts">
import { ref, defineComponent, inject, onUnmounted } from 'vue';
import TencentTranslate from './channel/TencentTranslate.vue';
import AlibabaTranslate from './channel/AlibabaTranslate.vue';
import GoogleTranslate from './channel/GoogleTranslate.vue';
import DeeplTranslate from './channel/DeeplTranslate.vue';
import { useTranslateStore, TranslateTabKey } from '@/stores/translateStore';

// import OpenAITranslate from './channel/OpenAITranslate.vue';
import { emitter, CONFIG_TYPE, CONFIG_RESET_TYPE } from './config';
export default defineComponent({
  components: {
    TencentTranslate,
    AlibabaTranslate,
    GoogleTranslate,
    DeeplTranslate,
    // OpenAITranslate,
  },

  

  

  setup() {
    const toConfigPage = () => {
      window.open('https://docs.qq.com/doc/DWWhuQnh1TEV0UFlH');
    };

    return {
      toConfigPage,
    };
  },

  data() {
    return {
      selectedType: TranslateTabKey.other,
    };
  },

  computed: {},

  created() {
    const configTypeCallback = (e:any) =>{
       let config = JSON.parse(JSON.stringify(e));
       this.selectedType = config.type;
    }
    emitter.on(CONFIG_TYPE, configTypeCallback);

    const resetCallback = (e:any) =>{
      this.selectedType = TranslateTabKey.other;
    }
    emitter.on(CONFIG_RESET_TYPE, resetCallback);

    onUnmounted(() => {
      emitter.off(CONFIG_TYPE, configTypeCallback);
      emitter.off(CONFIG_RESET_TYPE, resetCallback);
    });
   
  },
});
</script>

<style scoped lang="less">
.tips-container {
  position: relative;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;

  .tips-title {
    padding: 10px 5px;
    font-size: 14px;
    color: #333;
    font-family: PingFangSC-Regular, PingFang SC;
    line-height: 20px;
  }

  .tips-button-container {
    position: absolute;
    top: 100%;
    /* Push the button below the span */
    right: 0;

    .tips-button {
      color: #1890ff;
      font-size: 14px;
      padding: 10px 25px;
      font-family: PingFangSC-Regular, PingFang SC;

      text-decoration: underline;
    }
  }
}
</style>
