<template>
  <div class="config-container">
    <a-form
      ref="formRef"
      v-model="formState"
      :rules="rules"
      class="config-form"
      :label-col="labelCol"
      :wrapper-col="wrapperCol"
    >
      <a-form-item
        label="名称"
        name="name"
        placeholder="OpenAI翻译(可再编辑)"
      >
        <a-input v-model="formState.name" />
      </a-form-item>
      <a-form-item
        label="API Key"
        name="apikey"
      >
        <a-input-password
          v-model="formState.apikey"
          placeholder="input password"
        />
      </a-form-item>
      <a-form-item
        label="API 域名"
        name="zone"
      >
        <a-input-password
          v-model="formState.zone"
          placeholder="input password"
        />
      </a-form-item>

      <a-form-item class="openai-buttons">
        <a-button
          type="primary"
          @click="onSubmit"
        >
          Create
        </a-button>
        <a-button
          style="margin-left: 20px"
          @click="resetForm"
        >
          Reset
        </a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script lang="ts">
import { ValidateErrorEntity } from 'ant-design-vue/es/form/interface';
import { defineComponent, reactive, ref, toRaw, UnwrapRef } from 'vue';
import { bridgeAdaptor } from '@/adaptor/bridge';

interface FormState {
  name: string;
  apikey: string;
  zone: string;
}
export default defineComponent({
  setup() {
     const formRef = ref();
    const formState: UnwrapRef<FormState> = reactive({
      name: '',
      apikey: '',
      zone: '',
    });
    const rules = {
      name: [
        { required: true, message: '名称需要在3-8个字', trigger: 'blur',min: 1, max: 8},
      ],
      apikey: [
        { required: true, message: 'apikey 未输入', trigger: 'blur',min:1 },
      ],
      zone: [
        { required: true, message: 'API域名 未输入', trigger: 'blur',min:1 },
      ],
    };
    const onSubmit = () => {
      formRef.value
        .validate()
        .then(() => {
          console.log('values', formState, toRaw(formState));
        })
        .catch((error: ValidateErrorEntity<FormState>) => {
          console.log('error', error);
        });
    };
    const resetForm = () => {
      formRef.value.resetFields();
    };
    return {

       formRef,
        labelCol: { span: 10 },
      wrapperCol: { span: 14 },
                    formState,

      rules,
      onSubmit,
      resetForm,
    };
  },
});
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
  cursor: pointer;
  margin-bottom: 20px;

}

.openai-buttons {
  display: flex;
  padding: 30px 0px;
  justify-content: center;
  align-items: center;
}



:deep(.ant-input-password-icon) {
  color: #C9CDD4 !important;
}


.config-form {
  color: black !important;
  flex-direction: column;
  justify-content: space-between;
}


:deep(.ant-input) {
  color: #1D2229 !important;
  border: 1px solid #C9CDD4 !important;
  font-family: PingFangSC-Regular, PingFang SC;
  font-size: 14px !important;
}

:deep(.ant-form-item-required) {
  color: black !important;
}

:deep(.ant-select-selection-placeholder) {
  color: black !important;
}

:deep(.ant-select-arrow) {
  color: black !important;
}
</style>
