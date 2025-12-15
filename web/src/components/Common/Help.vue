<template>
  <a-modal
    maskClosable
    wrapClassName="modal-help"
    :visible="visible"
    :width="300"
    :centered="true"
    :footer="null"
    @cancel="emit('cancel')"
  >
    <div class="p-4 flex flex-col items-center text-center text-rp-neutral-10">
      <!-- <template v-if="envStore.viewerConfig.feedbackEmail || isWebEN">
        <p class="tip">Please contact us at</p>
        <a class="" :href="`mailto:${envStore.viewerConfig.feedbackEmail}`">
          {{ envStore.viewerConfig.feedbackEmail }}
        </a>
        <a-button
          size="large"
          type="primary"
          @click.stop="copyToPaste(envStore.viewerConfig.feedbackEmail!)"
          ><copy-outlined />Copy email</a-button
        >
      </template> -->
      <!-- <img
        class="mt-4 w-[200px]"
        src="https://readpaper.com/doc/assets/img/Qrcode.png"
        alt="QRCode"
      >
      <p class="mt-3 mb-0">
        {{ $t('wordings.contactWechat') }}
      </p>
      <p class="mb-0 text-xs text-rp-neutral-6">
        WeChat: readpaper888
      </p>
      <div class="mt-3 text-xs text-rp-neutral-6">
        <slot v-if="uid">
          uid: {{ uid }}
        </slot>
      </div> -->
      <a :href="userGuideLink" target="_blank" class="user-guide-link">{{ $t('home.upload.errorDetail.guide') }}</a>
    </div>
  </a-modal>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { Popover } from 'ant-design-vue';
import { useStore } from '~/src/store';
import { useI18n } from 'vue-i18n';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '@/hooks/useLanguage';

// import { CopyOutlined } from '@ant-design/icons-vue';
// import { copyToPaste } from '@/util/copy';
// import { useEnvStore } from '~/src/stores/envStore';


const visible = defineModel('visible', { default: false });

const emit = defineEmits<{
  (event: 'cancel'): void;
}>();

const { locale } = useI18n();

// 语言管理
const { isCurrentLanguage } = useLanguage();

// 创建用户指南链接的计算属性
const userGuideLink = computed(() => {
  if (isCurrentLanguage(Language.ZH_CN)) {
    return '/docs/zh/guide';
  }
  return '/docs/guide';
});

const store = useStore();
// const envStore = useEnvStore();


const uid = computed(() => store.state.user.userInfo?.id);
</script>

<style lang="less">
.modal-help {
  .ant-modal-close {
    color: theme('colors.rp-neutral-6');
    &:hover {
      color: theme('colors.rp-neutral-10');
    }
  }
  .ant-modal-content {
    background-color: #fff;
  }
}
.user-guide-link {
  display: block;
  color: #1f71e0;
  text-decoration: none;
  padding: 8px 0;
  
  &:hover {
    color: #0d5cbf;
    text-decoration: underline;
  }
}
</style>
