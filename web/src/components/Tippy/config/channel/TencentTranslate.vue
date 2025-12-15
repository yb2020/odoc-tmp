<template>
  <div class="config-container">
    <a-form
      ref="formRef"
      :model="txConfig"
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
          v-model:value="txConfig.name"
          placeholder="腾讯翻译君(可再编辑)"
          :disabled="txConfig.verified"
          :bordered="!txConfig.verified"
          style="border: 1px solid #C9CDD4"
        />
      </a-form-item>
      <a-form-item
        label="SecretID"
        name="txSecretId"
      >
        <a-input-password
          v-model:value="txConfig.txSecretId"
          :disabled="txConfig.verified"
          :bordered="!txConfig.verified"
          :visibilityToggle="!txConfig.verified"
          style="border: 1px solid #C9CDD4;"
        />
      </a-form-item>
      <a-form-item
        label="SecretKey"
        name="txSecretKey"
      >
        <a-input-password
          v-model:value="txConfig.txSecretKey"
          :disabled="txConfig.verified"
          :bordered="!txConfig.verified"
          :visibilityToggle="!txConfig.verified"
          style="border: 1px solid #C9CDD4;"
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
              <CheckCircleOutlined style="color: #52C41A" />
            </template>
          </a-avatar>
          <span
            v-if="checkSuccess"
            style="color: #52C41A"
          >验证成功</span>
          <a-button
            style="margin-left: 20px"
            :disabled="txConfig.verified"
            class="check-button"
            size="small"
            @click="checkAction"
          >
            验证
          </a-button>
        </a-row>
      </a-form-item>
      <a-form-item
        v-if="!txConfig.verified"
        class="tencent-buttons"
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
import { useTranslateStore ,TranslateTabKey} from '@/stores/translateStore';
import { ValidateErrorEntity } from 'ant-design-vue/es/form/interface';
import { defineComponent, reactive, ref, toRaw, UnwrapRef ,onUnmounted} from 'vue';
import { CheckCircleOutlined } from '@ant-design/icons-vue';
import { emitter, CONFIG_RESET_TYPE, CONFIG_ADD_TYPE } from '../config';
import { bridgeAdaptor } from '@/adaptor/bridge';
import { AddTxTranslateType } from '@/api/translate';
import {message} from 'ant-design-vue';
import { TxConfig } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/CustomTranslateInterface';

const store = useTranslateStore();

const txConfig = ref(store.txConfig as TxConfig) ;
const formRef = ref();

const submitEnable = ref(true);
const checkSuccess = ref(false);
checkSuccess.value = txConfig.value.verified;

const  addListener =(e:any) =>{
  if (!e) {
    return;
  }
  
  const deleteInfo = JSON.parse(JSON.stringify(e));
  if(deleteInfo.type === TranslateTabKey.tencent){
     store.tabs = store.tabs.filter((item :any) => {
      return item.type !== TranslateTabKey.tencent;
    });
    reset();
  }

}
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
  txSecretId: [{ required: true, message: 'SecretId 未输入', trigger: 'blur' }],
  txSecretKey: [
    { required: true, message: 'SecretKey 未输入', trigger: 'blur' },
  ],
};

const checkForm = () => {
  if(txConfig.value.txSecretId.length == 0){
    message.error('SecretId 未输入');
    return false;
  }
  if(txConfig.value.txSecretKey.length == 0){
    message.error('SecretKey 未输入');
    return false;
  }
};

const checkAction = async () => {

  const resp = await bridgeAdaptor.translateOnTX({
    text: 'hello world',
    secretId: txConfig.value.txSecretId,
    secretKey: txConfig.value.txSecretKey,
  });
  if (resp?.length > 0) {
    checkSuccess.value = true;
    submitEnable.value = false;
  }
};
const reset = () => {
  
  txConfig.value.verified = false;
  checkSuccess.value = false;
  submitEnable.value = true;
  txConfig.value.id = '';
  txConfig.value.name = '腾讯翻译君';
  txConfig.value.txSecretId = '';
  txConfig.value.txSecretKey = '';
  txConfig.value.createDate = '';

  store.txConfig = txConfig.value;

  const item = store.cusTabs.findIndex(
        (item :any) => item.type === TranslateTabKey.tencent
  );
  if(item === -1){
    store.cusTabs.push({
      name: txConfig.value.name ,
      type: TranslateTabKey.deepl ,
      verified: false,
      id: '',
    })
  }
};

const onSubmit = async () => {
  txConfig.value.createDate = new Date().getTime().toString();
  formRef.value
    .validate()
    .then(async () => {
      const response = await AddTxTranslateType({
        name: txConfig.value.name,
        verified: true,
        txSecretId: txConfig.value.txSecretId,
        txSecretKey: txConfig.value.txSecretKey,
        createDate: txConfig.value.createDate,
      });
      if (response) {
        txConfig.value.id = response;
      }
      txConfig.value.verified = true;
      emitter.emit(CONFIG_ADD_TYPE, txConfig.value);
    })
    .catch((error: ValidateErrorEntity<TxConfig>) => {});
};

const resetForm = () => {
  emitter.emit(CONFIG_RESET_TYPE, txConfig.value);
  formRef.value.resetFields();
 
};
</script>


<style scoped lang="less">
.config-container {
  background: white;
  padding: 0px;
  display: flex;
  height: 100%;
  width: 100%;

}

.ant-form-item {
  cursor: pointer;
  margin-bottom: 20px;
}

.tencent-buttons {
  display: flex;
  width: 100%;
  margin-top: 40px;
  margin-left: 30px;
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

:deep(.ant-input-placeholder) {
  color: black !important;
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

:deep(.ant-input) {
  color: #1D2229 !important;
  font-family: PingFangSC-Regular, PingFang SC;
  font-size: 14px !important;

}



:deep(.ant-form-item-required) {
  color: #1D2229 !important;
  text-align: end !important;

}
</style>

