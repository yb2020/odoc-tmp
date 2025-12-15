<template>
  <slot
    :loading="loading"
    :onClick="onAddProj"
  >
    <a-button
      type="primary"
      shape="round"
      size="large"
      class="!px-6 !flex items-center justify-center text-base font-bold"
      @click="onAddProj"
    >
      <template #icon>
        <LoadingOutlined v-if="loading" />
        <PlusOutlined
          v-else
          class="text-xl"
        />
      </template>{{ $t('common.aitools.addProj') }}
    </a-button>
  </slot>
</template>

<script setup lang="ts">
import { toRef, watch } from 'vue';
import { useProjectAdd } from '@common/hooks/aitools/useProjectAdd';
import { PlusOutlined, LoadingOutlined } from '@ant-design/icons-vue';
import { ResponseError } from '@common/api/type';
import { ERROR_CODE_NEED_VIP } from '@common/api/const';
import { message } from 'ant-design-vue';

const props = defineProps<{
  isOpenNewTab?: boolean;
}>();
const emit = defineEmits<{
  (e: 'preclick', evt: MouseEvent): void;
  (e: 'added', id: string): void;
  (e: 'errvip', err: ResponseError): void;
}>();

const isOpenNewTab = toRef(props, 'isOpenNewTab');

const { loading, run, error } = useProjectAdd(isOpenNewTab);
const onAddProj = async (e: MouseEvent) => {
  emit('preclick', e);
  if (!e.defaultPrevented) {
    const res = await run();

    emit('added', res?.projectId!);
  }
};

watch(error, () => {
  if (
    error.value instanceof ResponseError &&
    error.value.code === ERROR_CODE_NEED_VIP
  ) {
    emit('errvip', error.value);
  } else {
    message.error(error.value?.message);
  }
});
</script>
