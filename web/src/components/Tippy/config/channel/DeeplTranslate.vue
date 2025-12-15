<template>
  <div class="config-container">
    <a-form
      ref="formRef"
      :model="deeplConfig"
      :rules="rules"
      class="config-form"
      :wrapper-col="{ span: 16 }"
      :labelCol="{ span: 8}"
    >
      <a-form-item
        label="名称"
        name="name"
      >
        <a-input
          v-model:value="deeplConfig.name"
          placeholder="deepl翻译(可再编辑)"
          :disabled="deeplConfig.verified"
          style="border: 1px solid #C9CDD4"
        />
      </a-form-item>
      <a-form-item
        label="Key"
        name="deepLKey"
      >
        <a-input-password
          v-model:value="deeplConfig.deepLKey"
          :disabled="deeplConfig.verified"
          :visibilityToggle="!deeplConfig.verified"
          style="border: 1px solid #C9CDD4;"
        />
      </a-form-item>
      <a-form-item
        label="API"
        name="deepLApi"
      >
        <a-select
          v-model.value="deeplConfig.deepLApi"
          placeholder="Free"
          :disabled="deeplConfig.verified"
        >
          <a-select-option value="Free">
            Free
          </a-select-option>
          <a-select-option value="Pro">
            Pro
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item
        label="Formaility"
        name="deepLFormality"
      >
        <a-select
          v-model.value="deeplConfig.deepLFormality"
          placeholder="默认"
          :disabled="deeplConfig.verified"
        >
          <a-select-option value="default">
            默认
          </a-select-option>
          <a-select-option value="more">
            倾向于正式语言
          </a-select-option>
          <a-select-option value="less">
            倾向于非正式语言
          </a-select-option>
        </a-select>
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
              <CheckCircleOutlined style="color: #52C41A" />
            </template>
          </a-avatar>
          <span
            v-if="checkSuccess"
            style="color: #52C41A"
          >验证成功</span>
          <a-button
            style="margin-left: 20px"
            :disabled="deeplConfig.verified"
            class="check-button"
            size="small"
            @click="checkAction"
          >
            验证
          </a-button>
        </a-row>
      </a-form-item>

      <a-form-item
        v-if="!deeplConfig.verified"
        class="deepl-buttons"
        :wrapper-col="{ span: 24 }"
      >
        <a-button
          :wrapper-col="{span :10}"
          style="width:150px;height:32px;background:#F0F2F5;color:#4E5969;border:0px"
          @click="resetForm"
        >
          取消
        </a-button>
        <a-button
          :wrapper-col="{span :10}"
          style="margin-left: 20px;width:150px;height:32px"
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
import { useTranslateStore ,TranslateTabKey} from '~/src/stores/translateStore';
import { ValidateErrorEntity } from 'ant-design-vue/es/form/interface';
import { defineComponent, reactive, ref, toRaw, onUnmounted } from 'vue';
import { CheckCircleOutlined } from '@ant-design/icons-vue';
import { emitter, CONFIG_RESET_TYPE, CONFIG_ADD_TYPE } from '../config';
import { bridgeAdaptor } from '@/adaptor/bridge';
import { DeepLConfig } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/CustomTranslateInterface';
import { AddDeeplTranslateType, AddTxTranslateType } from '@/api/translate';

/// init data
const formRef = ref();
const submitEnable = ref(true);
const checkSuccess = ref(false);
const store = useTranslateStore();
const deeplConfig = ref(store.deeplConfig as DeepLConfig)  ;
checkSuccess.value = deeplConfig.value.verified;
/// init form rules
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
  deepLKey: [{ required: true, message: 'Key 未输入' ,trigger: 'blur'}],
  deepLApi: [
    { required: true, message: 'Please select API', trigger: 'change' },
  ],
  deepLFormality: [
    { required: true, message: 'Please select 环境', trigger: 'change' },
  ],
};

/// init listener
const addListener =(e: any) =>{
   if (!e) {
    return;
  }
  const deleteInfo = JSON.parse(JSON.stringify(e));
  if(deleteInfo.type === TranslateTabKey.deepl){

    store.tabs = store.tabs.filter((item :any) => {
      return item.type !== TranslateTabKey.deepl;
    });
    reset();
  }

};
emitter.on(CONFIG_RESET_TYPE, addListener);

onUnmounted(() => {
  emitter.off(CONFIG_RESET_TYPE,addListener);

});


/// init function
const checkAction = async () => {

  const resp = await bridgeAdaptor.translateOnDeepl({
    authKey: deeplConfig.value.deepLKey as string,
    text: 'hello world',
    api: deeplConfig.value.deepLApi as string,
  });
  if (resp?.length > 0) {
    checkSuccess.value = true;
    submitEnable.value = false;
  }
};

/// reset form model data
const reset = () => {
  
  deeplConfig.value.verified = false;
  checkSuccess.value = false;
  submitEnable.value = true;
  deeplConfig.value.id = '';
  deeplConfig.value.deepLApi = 'Free';
  deeplConfig.value.name = 'DeepL翻译';
  deeplConfig.value.deepLKey = '';
  deeplConfig.value.deepLFormality = 'default';
  deeplConfig.value.createDate = '';
  store.deeplConfig = deeplConfig.value;

  const item = store.cusTabs.findIndex(
        (item :any) => item.type === TranslateTabKey.deepl
  );
  if(item === -1){
    store.cusTabs.push({
      name: deeplConfig.value.name ,
      type: TranslateTabKey.deepl ,
      verified: false,
      id: '',
    })
  }
  
};
/// submit form data
const onSubmit = async () => {
  deeplConfig.value.createDate = new Date().getTime().toString();
  formRef.value
    .validate()
    .then(async () => {
      const response = await AddDeeplTranslateType({
        name: deeplConfig.value.name,
        verified: true,
        deepLKey: deeplConfig.value.deepLKey,
        deepLApi: deeplConfig.value.deepLApi,
        deepLFormality: deeplConfig.value.deepLFormality,
        createDate: deeplConfig.value.createDate,
      });
      if (response) {
        deeplConfig.value.id = response;
      }
      deeplConfig.value.verified = true;
      emitter.emit(CONFIG_ADD_TYPE, deeplConfig.value);
    })
    .catch((error: ValidateErrorEntity<DeepLConfig>) => {
      console.log('error', error);
    });
};

/// reset form data
const resetForm = () => {
  emitter.emit(CONFIG_RESET_TYPE, deeplConfig.value);
  
  formRef.value.resetFields();
};
</script>


<style scoped  lang="less">
.config-container {
  background: white;
  padding: 0px;
  display: flex;
  height: 100%;
  width: 100%;
}

.ant-form-item {
  cursor: pointer;
  margin-bottom: 15px;
}

.deepl-buttons {
  display: flex;
  width: 100%;
  margin-top: 30px;
  margin-left: 20px;
  justify-content: center;
}

.config-form {
  width: 100%;
  color: #C9CDD4 !important;
  flex-direction: column;
  justify-content: space-between;

  .tencent-verify {
    display: flex;
    justify-content: flex-end;

    .check-button {
      font-size: 12px;
      border-radius: 2px;
      height: 24px;
      width: 48px;

    }
  }
}

:deep(.ant-input) {
  color: #1D2229 !important;
  font-family: PingFangSC-Regular, PingFang SC;
  font-size: 14px !important;
}

:deep(.ant-form-item-required) {
  color: #1D2229 !important;
}

:deep(.ant-form-item) {
  color: #C9CDD4 !important;
}

:deep(.ant-form-item-label) {
  color: #1D2229 !important;
  text-align: end !important;
  font-family: PingFangSC-Regular, PingFang SC;
  font-size: 14px !important;
}

:deep(.ant-input-password-icon) {
  color: #C9CDD4 !important;
}

:deep(.ant-select-selection-placeholder) {
  color: black !important;
}

:deep(.ant-select-option) {
  color: black !important;
}


:deep(.ant-select-item-option-content) {
  color: black !important;
}


:deep(.ant-select-selection-item) {
  color: #1D2229 !important;
}



:deep(.ant-select-selector) {
  border: 1px solid #C9CDD4 !important;
}


:deep(.ant-form-item-label) {
  color: #1D2229 !important;
  text-align: end !important;
}

:deep(.ant-select-arrow) {
  color: #C9CDD4 !important;
  text-align: end;
}
</style>
