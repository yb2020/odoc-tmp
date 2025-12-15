<template>
  <Layout
    :title="item ? $t('glossary.editTitle') : $t('glossary.addTitle')"
    noDing
    :tippyHandler="onTippyHandler"
    group="glossary"
    :style="{ width: '400px' }"
  >
    <div class="p-3 glossary-create-edit">
      <a-form
        ref="formRef"
        :model="form"
        layout="vertical"
      >
        <a-form-item
          :label="$t('glossary.table.original')"
          required
          name="originalText"
        >
          <a-textarea
            v-model:value="form.originalText"
            class="!pb-6 js-interact-drag-ignore"
            :auto-size="{ minRows: 3, maxRows: 5 }"
          />
          <a-checkbox
            v-model:checked="form.matchCase"
            class="absolute bottom-1 right-1 rp-light-theme"
          >
            {{ $t('glossary.table.caseSensitive') }}
          </a-checkbox>
        </a-form-item>
        <a-form-item
          :label="$t('translate.translation')"
          required
          name="translationText"
        >
          <a-textarea
            v-model:value="form.translationText"
            class="!pb-6 js-interact-drag-ignore"
            :disabled="translationTextDisabled"
            :auto-size="{ minRows: 3, maxRows: 5 }"
          />
          <a-checkbox
            v-model:checked="form.ignored"
            class="absolute bottom-1 right-1 rp-light-theme"
            @change="onIgnoreChange"
          >
            {{ $t('glossary.table.noTranslation') }}
          </a-checkbox>
        </a-form-item>
        <a-form-item>
          <div class="flex justify-end">
            <a-button
              size="small"
              @click="onCancel"
            >
              {{
                $t('viewer.cancel')
              }}
            </a-button>
            <a-button
              size="small"
              class="ml-2"
              type="primary"
              :loading="submitLoading"
              @click="onSubmit"
            >
              {{ $t('glossary.save') }}
            </a-button>
          </div>
        </a-form-item>
      </a-form>
    </div>
  </Layout>
</template>
<script setup lang="ts">
import { GlossaryItem } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/GlossaryManage';
import Layout from '@/components/Tippy/Layout/index.vue';
import { computed, ref } from 'vue';
import { FormInstance } from 'ant-design-vue/es/form/Form';
import { addGlossaryItem, updateGlossaryItem } from '~/src/api/glossary';
import { message } from 'ant-design-vue';
import { useI18n } from 'vue-i18n';
const props = defineProps<{
  item: GlossaryItem | null;
  close: () => void;
  refresh: () => void;
}>();

const onTippyHandler = (event: 'ding' | 'close' | 'unding' | 'lock') => {
  if (event === 'close') {
    props.close();
  }
};

const form = ref({
  originalText: props.item?.originalText || '',
  translationText: props.item?.translationText || '',
  matchCase: props.item?.matchCase || false,
  ignored: props.item?.ignored || false,
});

const formRef = ref<FormInstance | null>(null);

const translationTextDisabled = computed(() => {
  return form.value.ignored;
});

const { t } = useI18n();

const submitLoading = ref(false);
const onSubmit = async () => {
  if (submitLoading.value) return;
  submitLoading.value = true;
  try {
    if (form.value.ignored) {
      await formRef.value?.validateFields(['originalText']);
    } else {
      await formRef.value?.validate();
    }
    if (props.item) {
      // 编辑
      await updateGlossaryItem({
        ...form.value,
        id: props.item.id,
      });
      message.success(t('glossary.editSuccessTip'));
    } else {
      // 添加
      await addGlossaryItem(form.value);
      message.success(t('glossary.addSuccessTip'));
    }
    await props.refresh();
    props.close();
  } catch (error) {}
  submitLoading.value = false;
};

const onCancel = () => {
  props.close();
};

const onIgnoreChange = (checked: boolean) => {
  if (checked && !form.value.translationText.trim()) {
    formRef.value?.clearValidate('translationText');
  }
};
</script>
<style less scoped>
.glossary-create-edit {
  :deep(.ant-col.ant-form-item-label label) {
    color: #000;
  }
  :deep(
      .ant-form-item-control-input
        .ant-form-item-control-input-content
        .ant-input
    ) {
    border-color: rgb(217, 217, 217);
    color: rgba(0, 0, 0, 0.65) !important;
    &::placeholder {
      color: #aaa;
    }
  }
}
</style>
