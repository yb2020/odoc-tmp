<template>
  <div class="content">
    <span class="source-title">
      {{ isUserUpload ? $t('source.personal') : $t('source.crawl') }}
    </span>
    <span
      class="jump ml-1"
      @click="handleJump"
    >{{ sourceMark }}</span>
    <!-- <span class="divide">&nbsp;&nbsp;|&nbsp;&nbsp;</span>
    <a
      href="/copyright"
      target="_blank"
    >
      <img
        class="icon"
        :src="pdfLicenceType"
        alt=""
      >
    </a> -->
    <!-- <span class="divide">&nbsp;&nbsp;|&nbsp;&nbsp;</span> -->
    <!-- <span
      class="report"
      @click.stop="handleReport"
    >{{
      $t('source.report.button')
    }}</span> -->
  </div>
  <Report v-model:visible="reportVisible" />
</template>

<script lang="ts" setup>
export interface CopyrightProps {
  licenceType: string;
  uploadUserId: string;
  crawlUrl: string;
  isUserUpload: boolean;
  sourceMark: string;
}

import { computed, ref } from 'vue';
import publicDomain from '@/assets/images/annotate-view/public_domain.png';
import cc from '@/assets/images/annotate-view/cc.png';
import declaration from '@/assets/images/annotate-view/declaration.png';
import declarationEN from '@/assets/images/annotate-view/declaration-en.png';
import Report from './Report.vue';
import { getDomainOrigin } from '~/src/util/env';
import { goPathPage } from '~/src/common/src/utils/url';
import { useI18n } from 'vue-i18n';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '~/src/hooks/useLanguage';

const { locale } = useI18n();
const { isCurrentLanguage } = useLanguage();

const licenceType = computed<Record<string, string>>(() => ({
  CC: cc,
  PUBLIC_DOMAIN: publicDomain,
  COPYRIGHT_DECLARE: isCurrentLanguage(Language.ZH_CN) ? declaration : declarationEN,
}));

const props = defineProps<CopyrightProps>();

const pdfLicenceType = computed(
  () =>
    licenceType.value[props.licenceType] || licenceType.value.COPYRIGHT_DECLARE
);

const reportVisible = ref(false);
const handleReport = () => {
  reportVisible.value = true;
};

const handleJump = () => {
  if (props.uploadUserId) {
    goPathPage(`${getDomainOrigin()}/user/${props.uploadUserId}`);
  } else {
    goPathPage(`${props.crawlUrl}`);
  }
};
</script>

<style lang="less" scoped>
.content {
  display: flex;
  align-items: center;
  font-family: PingFangSC-Regular, PingFang SC;
  font-weight: 400;
  z-index: 1;
  color: #4e5969;
  font-size: 13px;

  .divide {
    color: #e5e6eb;
  }

  .source-title {
    flex: 0 0 auto;
  }

  .jump {
    color: #245db3;
    cursor: pointer;
    flex: 0 0 auto;
  }

  .report {
    color: #4e5969;
    text-decoration: underline;
    cursor: pointer;
    flex: 0 0 auto;
    margin-right: 5px;
  }
  .icon {
    height: 14px;
  }
}

.mobile-viewport {
  .content {
    flex: 0 0 auto;
    font-size: 10rpx;
    font-weight: 400;
    color: #a8adb3;
    line-height: 18rpx;

    .report {
      display: none;
    }

    .divide {
      display: none;
    }
  }
}
</style>
