<template>
  <div class="config-container">
    <a-form
      ref="formRef"
      :model="aliConfig"
      :rules="rules"
      class="config-form"
      :wrapper-col="{ span: 16 }"
      :labelCol="{ span: 8 }"
    >
      <a-form-item
        label="名称"
        name="name"
      >
        <a-input
          v-model:value="aliConfig.name"
          placeholder="阿里翻译(可再编辑)"
          :disabled="aliConfig.verified"
          :bordered="!aliConfig.verified"
          style="border: 1px solid #C9CDD4"
        />
      </a-form-item>
      <a-form-item
        label="SecretId"
        name="aliAccessKeyId"
      >
        <a-input-password
          v-model:value="aliConfig.aliAccessKeyId"
          placeholder="input secretid"
          :disabled="aliConfig.verified"
          :bordered="!aliConfig.verified"
          :visibilityToggle="!aliConfig.verified"
          style="border: 1px solid #C9CDD4"
        />
      </a-form-item>
      <a-form-item
        label="SecretKey"
        name="aliAccessKeySecret"
      >
        <a-input-password
          v-model:value="aliConfig.aliAccessKeySecret"
          placeholder="input secretkey"
          :disabled="aliConfig.verified"
          :bordered="!aliConfig.verified"
          :visibilityToggle="!aliConfig.verified"
          style="border: 1px solid #C9CDD4"
        />
      </a-form-item>

      <a-form-item
        label="接口版本"
        name="aliInterfaceVersion"
      >
        <a-select
          v-model.value="aliConfig.aliInterfaceVersion"
          placeholder="版本选择"
          :disabled="aliConfig.verified"
        >
          <a-select-option value="general">
            通用版
          </a-select-option>
          <a-select-option value="title">
            专业版-商品标题
          </a-select-option>
          <a-select-option value="description">
            专业版-商品描述
          </a-select-option>
          <a-select-option value="communication">
            专业版-商品沟通
          </a-select-option>
          <a-select-option value="medical">
            专业版-医疗
          </a-select-option>
          <a-select-option value="social">
            专业版-社交
          </a-select-option>
          <a-select-option value="finance">
            专业版-金融
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
            :disabled="aliConfig.verified"
            class="check-button"
            size="small"
            @click="checkAction"
          >
            验证
          </a-button>
        </a-row>
      </a-form-item>

      <a-form-item
        v-if="!aliConfig.verified"
        class="ali-buttons"
        :wrapper-col="{ span: 24 }"
      >
        <a-button
          :wrapper-col="{span :10}"
          style="width:140px;height:32px;background:#F0F2F5;color:#4E5969;border:0px"
          @click="resetForm"
        >
          取消
        </a-button>
        <a-button
          :wrapper-col="{span :10}"
          style="margin-left: 20px;width:140px;height:32px"
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
import { AddTxTranslateType } from '@/api/translate';
import { AliConfig } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/CustomTranslateInterface';
import { AddAliTranslateType } from '../../../../api/translate';

/// init data
const formRef = ref();
const checkSuccess = ref(false);
const submitEnable = ref(true);
const store = useTranslateStore();
const aliConfig = ref(store.aliConfig as AliConfig);
checkSuccess.value = aliConfig.value.verified;


/// init rules
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
  aliAccessKeyId: [
    { required: true, message: 'SecretId 未输入', trigger: 'blur' },
  ],
  aliAccessKeySecret: [
    { required: true, message: 'SecretKey 未输入', trigger: 'blur' },
  ],
  aliInterfaceVersion: [
    {
      required: true,
      message: 'Please select Activity zone',
      trigger: 'change',
    },
  ],
};

const addListener =(e: any) =>{
   if (!e) {
    return;
  }
  
  const deleteInfo = JSON.parse(JSON.stringify(e));

  if(deleteInfo.type === TranslateTabKey.ali){
     
    store.tabs = store.tabs.filter((item :any) => {
      return item.type !== TranslateTabKey.ali;
    });

    reset();
  }

};
emitter.on(CONFIG_RESET_TYPE, addListener);

onUnmounted(() => {
  emitter.off(CONFIG_RESET_TYPE,addListener);

});


/// init function
const reset = () => {
  aliConfig.value.verified = false;
  checkSuccess.value = false;
  submitEnable.value = true;
  aliConfig.value.id = '';
  aliConfig.value.name = '阿里云翻译';
  aliConfig.value.aliAccessKeyId = '';
  aliConfig.value.aliAccessKeySecret = '';
  aliConfig.value.aliInterfaceVersion = 'general';

  store.aliConfig = aliConfig.value;
  const idx = store.cusTabs.findIndex(
        (item :any) => item.type === TranslateTabKey.ali
  );
  if(idx === -1){
    store.cusTabs.push({
      name: aliConfig.value.name ,
      type: TranslateTabKey.ali ,
      verified: false,
      id: '',
    })
  }
};

const onSubmit = async () => {
  aliConfig.value.createDate = new Date().getTime().toString();
  formRef.value
    .validate()
    .then(async () => {
      const response = await AddAliTranslateType({
        name: aliConfig.value.name,
        verified: true,
        aliAccessKeyId: aliConfig.value.aliAccessKeyId,
        aliAccessKeySecret: aliConfig.value.aliAccessKeySecret,
        aliInterfaceVersion: aliConfig.value.aliInterfaceVersion,
        createDate: aliConfig.value.createDate,
      });
      if (response) {
        aliConfig.value.id = response;
      }

      aliConfig.value.verified = true;
      emitter.emit(CONFIG_ADD_TYPE, aliConfig.value);
    })
    .catch((error: ValidateErrorEntity<AliConfig>) => {});
};
const resetForm = () => {
  emitter.emit(CONFIG_RESET_TYPE, aliConfig.value);
  
  formRef.value.resetFields();
};

const checkAction = async () => {
  const resp = await bridgeAdaptor.translateOnAli({
    text: 'hello world',
    scene: aliConfig.value.aliInterfaceVersion as string,
    accessKeyId: aliConfig.value.aliAccessKeyId as string,
    accessKeySecret: aliConfig.value.aliAccessKeySecret as string,
  });
  if (resp?.length > 0) {
    checkSuccess.value = true;
    submitEnable.value = false;
  }
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
  margin-bottom: 10px;
}

.ali-buttons {
  display: flex;
  width: 100%;
  margin-top: 30px;
  margin-left: 20px;
  justify-content: center;
}

.config-form {
  color: black !important;
  flex-direction: column;
  width: 100%;
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

:deep(.ant-select-selector) {
  border: 1px solid #C9CDD4 !important;
}

:deep(.ant-input-placeholder) {
  color: black !important;
}

:deep(.ant-input-password-icon) {
  color: #C9CDD4 !important;
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

:deep(.ant-select-selection-placeholder) {
  color: black !important;
}

:deep(.ant-select-arrow) {
  color: #C9CDD4 !important;
}

:deep(.ant-form-item-required) {
  color: black !important;
  text-align: end;
}

:deep(.ant-form-item-label) {
  color: #1D2229 !important;
  text-align: end !important;
  font-family: PingFangSC-Regular, PingFang SC;
  font-size: 14px !important;
}


:deep(.ant-select-dropdown) {
  z-index: 10003 !important;

}
</style>
