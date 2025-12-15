<template>
  <div class="translate-word js-interact-drag-ignore">
    <div>
      <WordPronunciation
        v-if="translatedData.britishSymbol"
        prefix="英"
        :title="translatedData.britishSymbol"
        :type="translatedData.britishFormat"
        :audio="translatedData.britishPronunciation"
      />
      <WordPronunciation
        v-if="translatedData.americaSymbol"
        prefix="美"
        :title="translatedData.americaSymbol"
        :type="translatedData.americaFormat"
        :audio="translatedData.americaPronunciation"
      />
    </div>

    <ul
      v-if="translatedData.targetResp?.length"
      class="translate-part-list"
      :style="{ fontSize: fontSize + 'px', lineHeight: '1.4' }"
    >
      <li
        v-for="item in translatedData.targetResp"
        :key="item.part"
        class="translate-part-item"
      >
        <span class="translate-part">{{ item.part }}</span>
        <div class="translate-part-txt">
          <p
            v-for="text in item.targetContent"
            :key="text"
          >
            {{ text }};
          </p>
        </div>
      </li>
    </ul>
  </div>
</template>
<script setup lang="ts">
import { UniTranslateResp } from '~/src/api/translate';
import WordPronunciation from '@common/components/Notes/components/WordPronunciation.vue';

defineProps<{
  fontSize: string;
  translatedData: UniTranslateResp;
}>();
</script>
<style scoped lang="less">
.translate-word {
  margin-top: 14px;

  .pronunciation-item + .pronunciation-item {
    margin-top: 8px;
  }

  .translate-part-list {
    padding: 0;
    margin: 12px 0 14px;
  }

  .translate-part-item {
    display: flex;
    // font-size: 16px;
    // line-height: 26px;
    color: #1d2229;

    .translate-part {
      margin-right: 4px;
    }

    .translate-part-txt {
      margin-left: 12px;
      font-weight: 400;
      color: #1d2129;
      cursor: text;

      p {
        margin: 0;
      }
    }
  }

  .translate-part-item + .translate-part-item {
    margin-top: 8px;
  }
}
</style>
