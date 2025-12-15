<template>
  <div class="config-container">
    <a-form
      ref="formRef"
      :model="googleConfig"
      :rules="rules"
      class="config-form"
      :wrapper-col="{ span: 16 }"
      :labelCol="{ span: 10 }"
    >
      <a-form-item
        label="谷歌翻译"
        name="name"
        placeholder="谷歌翻译(可再编辑)"
      >
        <a-input
          v-model:value="googleConfig.name"
          style="border: 1px solid #c9cdd4"
          :disabled="googleConfig.verified"
        />
      </a-form-item>
      <a-form-item name="googleApiKey">
        <template #label>
          <div
            style="
              color: #1d2229;
              font-family: PingFangSC-Regular, PingFang SC;
              text-align: end;
              font-size: 14px;
            "
          >
            API Key:
            <a-tooltip placement="topRight">
              <template #title>
                若不填Key，则可能存在不稳定情况
              </template>
              <QuestionCircleOutlined style="{color:'#c9cdd4'}" />
            </a-tooltip>
          </div>
        </template>



        <a-input-password
          v-model:value="googleConfig.googleApiKey"
          :disabled="googleConfig.verified"
          :visibilityToggle="!googleConfig.verified"
          style="border: 1px solid #c9cdd4"
        />
      </a-form-item>

      <a-form-item class="tencent-verify">
        <a-row
          type="flex"
          justify="end"
        >
          <a-avatar
            v-if="checkSuccess"
            :size="28"
            style="background-color: white"
          >
            <template #icon>
              <CheckCircleOutlined style="color: #52c41a" />
            </template>
          </a-avatar>
          <span
            v-if="checkSuccess"
            style="color: #52c41a"
          >验证成功</span>

          <a-button
            type="default"
            style="margin-left: 20px"
            :disabled="googleConfig.verified"
            :loading="isLoading"
            class="check-button"
            size="small"
            @click="checkAction"
          >
            {{ checkTitle }}
          </a-button>
        </a-row>
      </a-form-item>
      <a-form-item
        v-if="!submitted"
        class="google-buttons"
        :wrapper-col="{ span: 24 }"
      >
        <a-button
          :wrapper-col="{ span: 10 }"
          style="
            width: 130px;
            height: 32px;
            background: #F0F2F5;
            color: #4E5969;
            border: 0px;
          "
          @click="resetForm"
        >
          取消
        </a-button>
        <a-button
          :wrapper-col="{ span: 10 }"
          style="margin-left: 20px; width: 130px; height: 32px"
          :disabled="submitEnable"
          @click="onSubmit"
        >
          确定
        </a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script lang="ts" setup>
import {
  useTranslateStore,
  TranslateTabKey,
  GoogleTranslateStyle
} from '~/src/stores/translateStore';
import { ValidateErrorEntity } from 'ant-design-vue/es/form/interface';
import { defineComponent, reactive, ref, toRaw, onUnmounted } from 'vue';
import { CheckCircleOutlined ,QuestionCircleOutlined} from '@ant-design/icons-vue';
import { emitter, CONFIG_RESET_TYPE, CONFIG_ADD_TYPE } from '../config';
import { bridgeAdaptor } from '@/adaptor/bridge';
import { AddGoogleTranslateTpe } from '@/api/translate';
import { GoogleConfig } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/CustomTranslateInterface';

const formRef = ref();
const submitEnable = ref(true);
const checkSuccess = ref(false);
const submitted = ref(false);
const store = useTranslateStore();
const isLoading = ref(false);
const checkTitle = ref('验证');
const googleConfig = ref(store.googleConfig as GoogleConfig);
checkSuccess.value = googleConfig.value.verified;
submitted.value = googleConfig.value.verified;
/// init listener
const addListener = (e: any) => {
  if (!e) {
    return;
  }

  const deleteInfo = JSON.parse(JSON.stringify(e));

  if (deleteInfo.type === TranslateTabKey.google) {
    store.tabs = store.tabs.filter((item: any) => {
      return item.type !== TranslateTabKey.google;
    });

    reset();
  }
};
emitter.on(CONFIG_RESET_TYPE, addListener);

onUnmounted(() => {
  emitter.off(CONFIG_RESET_TYPE,addListener);

});

const rules = {
  name: [
    {
      required: true,
      message: '名称需要在3-8个字',
      trigger: 'blur',
      min: 1,
      max: 8,
    },
  ],
  googleApiKey: [{ required: false, message: 'API Key 未输入', trigger: 'blur' }],
};

const reset = () => {
  googleConfig.value.verified = false;
  checkSuccess.value = false;
  submitEnable.value = true;
  submitted.value = false;
  isLoading.value = false;
  googleConfig.value.id = '';
  googleConfig.value.name = '谷歌翻译';
  googleConfig.value.googleApiKey = '';
  googleConfig.value.createDate = '';
  store.googleConfigVersion = GoogleTranslateStyle.none;
  store.googleConfig = googleConfig.value;
  checkTitle.value = '验证';
  const idx = store.cusTabs.findIndex(
    (item: any) => item.type === TranslateTabKey.google
  );
  if (idx === -1) {
    store.cusTabs.push({
      name: googleConfig.value.name,
      type: TranslateTabKey.google,
      verified: false,
      id: '',
    });
  }
};

const checkAction = async () => {
    
  isLoading.value = true;

    
  checkTitle.value = '验证中...';

  const resp = await bridgeAdaptor.translateOnGoogle({
    text: 'hello world' as string,
    projectId: googleConfig.value.googleApiKey as string,
  });
  if (resp?.length > 0) {
    checkSuccess.value = true;
    submitEnable.value = false;
    isLoading.value = false;
    checkTitle.value = '验证';

  }else{
    isLoading.value = false;
    checkTitle.value = '验证失败';
  }
};
const onSubmit = async () => {

  googleConfig.value.createDate = new Date().getTime().toString();
  store.googleConfigVersion = googleConfig.value.googleApiKey?.length === 0 ? GoogleTranslateStyle.default : GoogleTranslateStyle.custom;
  formRef.value
    .validate()
    .then(async () => {
      const response = await AddGoogleTranslateTpe({
        name: googleConfig.value.name,
        verified: true,
        googleApiKey: googleConfig.value.googleApiKey,
        createDate: googleConfig.value.createDate,
      });
      if (response) {
        googleConfig.value.id = response;
      }
      submitted.value = true;
      googleConfig.value.verified = true;
      emitter.emit(CONFIG_ADD_TYPE, googleConfig.value);
    })
    .catch((error: ValidateErrorEntity<GoogleConfig>) => {
    
    });
};
const resetForm = () => {
  emitter.emit(CONFIG_RESET_TYPE, googleConfig.value);
  formRef.value.resetFields();
};
</script>


<style scoped  lang="less">
.config-container {
  background: white;
  padding: 20px;
  display: flex;
  height: 100%;
  width: 100%;
}

.ant-form-item {
  margin-bottom: 20px;
  color: #1d2229 !important;
}

.google-buttons {
  display: flex;
  width: 100%;
  margin-top: 50px;
  margin-left: 20px;
  justify-content: center;
}

:deep(.ant-input-password-icon) {
  color: #c9cdd4 !important;
}

:deep(.ant-form-item-label) {
  color: #1d2229 !important;
  text-align: end !important;
  font-family: PingFangSC-Regular, PingFang SC;
  font-size: 14px !important;
}

.config-form {
  width: 100%;
  color: #1d2229 !important;
  flex-direction: column;
  justify-content: space-between;

  .tencent-verify {
    display: flex;
    justify-content: flex-end;

    .check-button {
      color: #a6c6f3;
      font-size: 12px;
      border-radius: 2px;
      background: white;
    }
  }
}

:deep(.ant-input) {
  color: #1d2229 !important;
  font-family: PingFangSC-Regular, PingFang SC;
  font-size: 14px !important;
}

:deep(.ant-btn)::before {
  background: white !important;
}

:deep(.ant-form-item-required) {
  color: #1d2229 !important;
  text-align: end;
}
</style>



