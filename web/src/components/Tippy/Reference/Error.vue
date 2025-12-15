<template>
  <div class="error-container">
    <div class="title">
      {{ title }}
    </div>
    <div class="tip">
      <InfoCircleOutlined />{{ $t('message.inValidPaperTip') }}
    </div>
  </div>
  <div class="bottom">
    <span class="detail pr-2">{{ $t('viewer.seeMoreResults') }}</span>
    <span @click="handleClick('google')">{{ $t('viewer.googleScholar') }}</span>
    <span @click="handleClick('arxivOrg')">arxiv.org</span>
    <!-- <span
      v-if="envStore.viewerConfig.baiduScholarSearch !== false"
      @click="handleClick('baidu')"
    >{{ $t('viewer.baiduScholar') }}</span> -->
    <span @click="handleClick('semantic')">semantic scholar</span>
  </div>
</template>

<script lang="ts" setup>
import { InfoCircleOutlined } from '@ant-design/icons-vue';
import { useEnvStore } from '~/src/stores/envStore';
import { goPathPage } from '~/src/common/src/utils/url';

const props = defineProps<{ title: string }>();

const handleClick = (type: 'google' | 'baidu' | 'semantic' | 'arxivOrg') => {
  if (type === 'google') {
    goPathPage(
      `https://scholar.google.com.hk/scholar?q=${encodeURIComponent(
        props.title
      )}`
    );
  } else if (type === 'baidu') {
    goPathPage(
      `https://xueshu.baidu.com/s?wd=${encodeURIComponent(props.title)}`
    );
  } else if (type === 'semantic') {
    goPathPage(
      `https://www.semanticscholar.org/search?q=${encodeURIComponent(
        props.title
      )}&sort=relevance`
    );
  }else if(type === 'arxivOrg'){
    goPathPage(
      `http://arxiv.org/search/?query=${encodeURIComponent(
        props.title
      )}&searchtype=all&source=header`
    );
  }
};

const envStore = useEnvStore();
</script>

<style lang="less" scoped>
.error-container {
  padding: 20px 24px;

  .title {
    font-size: 14px;
    font-family: Lato-Bold, Lato;
    font-weight: bold;
    color: #262625;
    line-height: 22px;
    display: -webkit-box;
    -webkit-line-clamp: 5;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-all;
  }
  .tip {
    font-size: 14px;
    font-weight: 400;
    color: #73716f;
    line-height: 22px;
    margin-top: 12px;

    span {
      margin-right: 4px;
    }
  }
}

.bottom {
  height: 48px;
  border-top: 1px solid #e4e7ed;

  display: flex;
  align-items: center;

  span {
    margin-right: 16px;
    font-size: 12px;
    font-family: PingFangSC-Medium, PingFang SC;
    font-weight: 500;
    color: #73716f;
    cursor: pointer;

    &:hover {
      text-decoration: underline;
    }
  }

  .detail {
    margin-left: 24px;
    margin-right: 0;
    cursor: default;

    &:hover {
      text-decoration: none;
    }
  }
}
</style>
