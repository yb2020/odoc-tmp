<template>
  <a-modal
    :visible="visible"
    :title="isEditMode ? t('navbar.websiteformmodel.modalTitleEdit') : t('navbar.websiteformmodel.modalTitleAdd')"
    :ok-text="t('navbar.websiteformmodel.okText')"
    :cancel-text="t('navbar.websiteformmodel.cancelText')"
    @ok="handleOk"
    @cancel="handleCancel"
    :confirm-loading="confirmLoading"
    destroy-on-close
  >
    <a-form :model="formState" ref="formRef" :label-col="{ span: 4 }" :wrapper-col="{ span: 20 }">
      <a-form-item
        :label="t('navbar.websiteformmodel.nameLabel')"
        name="name"
        :rules="nameRules"
      >
        <a-input v-model:value="formState.name" :placeholder="t('navbar.websiteformmodel.namePlaceholder')" />
      </a-form-item>
      <a-form-item
        :label="t('navbar.websiteformmodel.urlLabel')"
        name="url"
        :rules="urlRules"
      >
        <a-input v-model:value="formState.url" :placeholder="t('navbar.websiteformmodel.urlPlaceholder')" />
      </a-form-item>
      <a-form-item v-if="apiError" :wrapper-col="{ offset: 4, span: 20 }">
        <a-alert :message="apiError" type="error" show-icon />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue';
import { useI18n } from 'vue-i18n';
import type { FormInstance } from 'ant-design-vue';
import { createWebsite, updateWebsite } from '@/api/nav';
import { Website } from 'go-sea-proto/gen/ts/nav/Website';
import { watch, computed } from 'vue';

const { t } = useI18n();

const props = defineProps<{
  visible: boolean;
  websiteToEdit?: Website | null;
}>();

const emit = defineEmits(['update:visible', 'success']);

const formRef = ref<FormInstance>();
const confirmLoading = ref(false);
const apiError = ref('');

const formState = reactive({
  id: null as bigint | null,
  name: '',
  url: '',
});

const isEditMode = computed(() => !!props.websiteToEdit);

const nameRules = computed(() => [
  { required: true, message: t('navbar.websiteformmodel.nameRequiredError') },
]);

const urlRules = computed(() => [
  { required: true, message: t('navbar.websiteformmodel.urlRequiredError') },
  { type: 'url', message: t('navbar.websiteformmodel.urlInvalidError') },
]);

watch(() => props.websiteToEdit, (newVal) => {
  if (newVal) {
    formState.id = newVal.id;
    formState.name = newVal.name;
    formState.url = newVal.url;
  } else {
    // Reset form for creation
    formState.id = null;
    formState.name = '';
    formState.url = '';
  }
});

const handleCancel = () => {
  formRef.value?.resetFields();
  apiError.value = '';
  emit('update:visible', false);
};

const handleOk = () => {
  formRef.value?.validate().then(async () => {
    apiError.value = ''; // 尝试提交时先清空旧错误
    confirmLoading.value = true;
    try {
      if (isEditMode.value && formState.id) {
        await updateWebsite({
          id: formState.id,
          name: formState.name,
          url: formState.url,
          iconUrl: props.websiteToEdit?.iconUrl || '',
          openType: props.websiteToEdit?.openType || 1,
        });
      } else {
        await createWebsite({
          name: formState.name,
          url: formState.url,
          iconUrl: '', // Per user request, this is hidden and should be empty
          openType: 1, // Default value
        });
      }
      emit('success');
      emit('update:visible', false);
      formRef.value?.resetFields();
    } catch (error) {
      console.error('Failed to create website:', error);
      apiError.value = isEditMode.value ? t('navbar.websiteformmodel.updateFailedError') : t('navbar.websiteformmodel.addFailedError');
    } finally {
      confirmLoading.value = false;
    }
  });
};
</script>

<style lang="less" scoped>
:deep(.ant-input) {
  color: var(--site-theme-text-primary);
  background-color: var(--site-theme-background);
  border-color: var(--site-theme-border-color);

  &::placeholder {
    color: var(--site-theme-placeholder-color);
  }

  &:hover, &:focus {
    border-color: var(--site-theme-primary-color);
  }
}

:deep(.ant-modal-content) {
  background-color: var(--site-theme-background);
}

:deep(.ant-modal-header) {
  background-color: var(--site-theme-background);
}

:deep(.ant-modal-title) {
  color: var(--site-theme-text-primary);
}

:deep(.ant-form-item-label > label) {
  color: var(--site-theme-text-primary);
}

:deep(.ant-btn) {
  &.ant-btn-primary {
    background-color: var(--site-theme-primary-color);
    border-color: var(--site-theme-primary-color);
    color: var(--site-theme-text-inverse);
    
    &:hover, &:focus {
      background-color: var(--site-theme-primary-color-hover);
      border-color: var(--site-theme-primary-color-hover);
    }
  }
  
  &.ant-btn-default {
    background-color: var(--site-theme-background-secondary);
    border-color: var(--site-theme-border-color);
    color: var(--site-theme-text-primary);
    
    &:hover, &:focus {
      border-color: var(--site-theme-primary-color);
      color: var(--site-theme-primary-color);
    }
  }
}

:deep(.ant-modal-footer) {
  border-top-color: var(--site-theme-border-color);
}
</style>
