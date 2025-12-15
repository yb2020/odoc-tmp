<script setup lang="ts">
import { onMounted } from 'vue';

import {
  ApplyQuotaStatus,
  GetApplyQuotaResp,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText';
import Icon from '@ant-design/icons-vue';
import ApplyModal from '@common/components/AITools/Apply/index.vue';
import { useUserStore } from '@common/stores/user';
import useApplyPermission from '@common/hooks/aitools/useApplyPermission';
import { initWelcomeAnimate } from './animate';

const emit = defineEmits<{
  (e: 'applied', x: GetApplyQuotaResp): void;
}>();

const userStore = useUserStore();
const { applyStatus, applyLoading, applyPermission } = useApplyPermission();

const handleLogin = () => {
  userStore.openLogin();
};

const handleApply = async () => {
  if (applyLoading.value) {
    return;
  }
  try {
    const res = await applyPermission();

    emit('applied', res!);
  } catch (e) {}
};

onMounted(() => {
  initWelcomeAnimate('js-particles-js');
});
</script>
<template>
  <div class="welcome-bg h-full">
    <div class="flex flex-col items-center justify-center h-full">
      <h1 class="mb-8">
        <img
          src="@common/../assets/images/aitools/polish-logo.svg"
          alt="AI Polish"
          class="h-22"
        >
      </h1>
      <p
        class="text-rp-neutral-10 font-medium text-2xl leading-[36px] text-center mb-12"
      >
        {{ $t('common.aitools.welcome.title') }}<br>
        <span v-if="!userStore.isLogin()">
          {{ $t('common.aitools.welcome.tip1') }}
        </span>
      </p>
      <div class="mb-24">
        <a-button
          v-if="!userStore.isLogin()"
          class="welcome-btn rp-bg-linear"
          @click="handleLogin"
        >
          {{ $t('common.account.login') }} /
          {{ $t('common.account.register') }}
        </a-button>
        <a-button
          v-else-if="applyStatus === ApplyQuotaStatus.APPLYING"
          class="welcome-btn"
          disabled
        >
          {{ $t('common.aitools.welcome.applied') }}
        </a-button>
        <a-button
          v-else
          class="welcome-btn"
          :loading="applyLoading"
          @click="handleApply"
        >
          <icon :style="{ color: 'white' }">
            <template #component>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="20"
                height="20"
                viewBox="0 0 20 20"
                fill="none"
              >
                <path
                  fill-rule="evenodd"
                  clip-rule="evenodd"
                  d="M13.9882 4.19643V5.26786C13.9882 5.36607 13.9078 5.44643 13.8096 5.44643H5.23818C5.13996 5.44643 5.05961 5.36607 5.05961 5.26786V4.19643C5.05961 4.09821 5.13996 4.01786 5.23818 4.01786H13.8096C13.9078 4.01786 13.9882 4.09821 13.9882 4.19643ZM5.05961 7.41072C5.05961 7.3125 5.13996 7.23214 5.23818 7.23214H9.34532C9.44353 7.23214 9.52389 7.3125 9.52389 7.41072V8.48214C9.52389 8.58036 9.44353 8.66072 9.34532 8.66072H5.23818C5.13996 8.66072 5.05961 8.58036 5.05961 8.48214V7.41072ZM3.27389 17.3214H8.45246C8.55068 17.3214 8.63103 17.4018 8.63103 17.5V18.75C8.63103 18.8482 8.55068 18.9286 8.45246 18.9286H2.38103C1.98594 18.9286 1.66675 18.6094 1.66675 18.2143V0.714286C1.66675 0.319196 1.98594 0 2.38103 0H16.6667C17.0618 0 17.381 0.319196 17.381 0.714286V8.03572C17.381 8.13393 17.3007 8.21429 17.2025 8.21429H15.9525C15.8542 8.21429 15.7739 8.13393 15.7739 8.03572V1.60714H3.27389V17.3214ZM17.6479 10.2093L19.7908 12.3522C19.8572 12.4185 19.9098 12.4972 19.9457 12.5838C19.9816 12.6705 20.0001 12.7634 20.0001 12.8572C20.0001 12.951 19.9816 13.0438 19.9457 13.1305C19.9098 13.2171 19.8572 13.2959 19.7908 13.3622L13.153 20H10.0001V16.8471L16.6379 10.2093C16.7042 10.143 16.783 10.0903 16.8696 10.0544C16.9563 10.0185 17.0491 10 17.1429 10C17.2367 10 17.3296 10.0185 17.4163 10.0544C17.5029 10.0903 17.5816 10.143 17.6479 10.2093ZM11.4287 18.5714H12.5615L16.1329 15L15.0001 13.8672L11.4287 17.4387V18.5714ZM16.0102 12.8572L17.1429 13.9899L18.2757 12.8572L17.1429 11.7244L16.0102 12.8572Z"
                  fill="white"
                />
              </svg>
            </template>
          </icon>
          {{ $t('common.aitools.welcome.apply') }}
        </a-button>
      </div>
    </div>
    <div
      id="js-particles-js"
      class="absolute top-0 left-0 bottom-0 right-0 -z-[1]"
    />
  </div>
  <ApplyModal />
</template>
<style scoped lang="less">
.welcome {
  &-btn.ant-btn {
    background: linear-gradient(90deg, #7bb3ff 0%, #e78ef7 100%) !important;
    color: #fff;
    font-weight: 500;
    font-size: 18px;
    line-height: 25px;
    height: 56px;
    width: 320px;
    border: none;
    border-radius: 30px;
  }
}
</style>
