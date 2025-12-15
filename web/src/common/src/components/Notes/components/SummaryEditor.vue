<template>
  <a-spin
    wrapperClassName="summary-editor h-full"
    :spinning="!data && loading"
  >
    <Milkdown
      v-model="summary"
      :lang="locale"
      placeholder=""
      :ignore-menu-keys="['InsertImage']"
      :default-plugins="plugins"
      :uploader="uploader"
      @update="onUpdate"
    />
    <div
      v-if="isSaving"
      class="absolute left-0 top-12 w-full flex items-center justify-center saving-indicator"
    >
      <LoadingOutlined class="text-xl" />
      <span class="text-xs ml-2">{{ $t('message.savingNote') }}</span>
    </div>
    <a-modal
      :visible="!!conflicted"
      type="confirm"
      title="检测到内容冲突"
      okText="覆盖"
      okType="danger"
      cancelText="同步远程内容"
      :cancelButtonProps="{
        type: 'primary',
      }"
      @ok="conflicted?.(true)"
      @cancel="conflicted?.(false)"
    >
      <p>
        远程内容保存时间：
        {{ dayjs(Number(data?.modifyDate)).format('YYYY-MM-DD HH:mm:ss') }}
      </p>
      <p>提示：同步远程内容后，通过Ctrl+Z可以恢复当前编辑版本</p>
    </a-modal>
  </a-spin>
</template>

<script setup lang="ts">
import _ from 'lodash';
import dayjs from 'dayjs';
import { message } from 'ant-design-vue';
import { LoadingOutlined } from '@ant-design/icons-vue';
import { useI18n } from 'vue-i18n';
import { ref, toRef } from 'vue';
import { uploadImage, ImageStorageType } from '@common/api/upload';
import Milkdown, { UpdatePayload } from '@common/components/Milkdown/index.vue';
import { DefaultPlugins } from '@common/components/Milkdown/utils';
import { GetSummaryResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/NoteModule';
import { useSummary } from '../useSummary';

const props = defineProps<{
  noteId: string;
}>();
const emit = defineEmits<{
  (e: 'loaded', data: GetSummaryResponse): void;
}>();

const noteId = toRef(props, 'noteId');
const conflicted = ref<(v: boolean) => void>();

const { t, locale } = useI18n();
const { data, summary, loading, refresh, save, isSaving } = useSummary(noteId, {
  onSuccess(res: GetSummaryResponse) {
    emit('loaded', res);
  },
  onConflict: async () => {
    return new Promise<boolean>((resolve) => {
      conflicted.value = (v: boolean) => {
        resolve(v);
        conflicted.value = undefined;
      };
    });
  },
});

const plugins = _.omit(
  {
    ...DefaultPlugins,
  },
  'block'
);

const uploader = (files: FileList) => {
  const images: File[] = [];

  for (let i = 0; i < files.length; i++) {
    const file = files.item(i);
    if (!file) {
      continue;
    }

    if (!file.type.includes('image')) {
      continue;
    }

    images.push(file);
  }

  return Promise.all(
    images.map(async (image) => {
      const src = await uploadImage(image, ImageStorageType.markdown);
      const alt = image.name;
      return {
        src,
        alt,
      };
    })
  );
};

const onUpdate = _.debounce(async ({ markdown }: UpdatePayload) => {
  if (markdown === data.value?.content) {
    return;
  }

  if (isSaving.value) {
    message.warn(`${t('message.savingNote')}...`);
    return;
  }

  await save();
}, 1500);

defineExpose({
  refresh,
});
</script>

<style lang="postcss">
.summary-editor {
  .ant-spin-container {
    height: 100%;
  }

  .milkdown-wrapper,
  .milkdown-menu-wrapper {
    height: 100%;
  }

  .milkdown-menu-wrapper {
    display: flex;
    flex-direction: column;
  }

  .milkdown-inner {
    flex: 1 1 auto;
    min-height: 0;
  }

  .milkdown {
    height: 100%;
    min-height: 100%;
    overflow: auto;

    /* 点击触发编辑 */
    .editor {
      height: 100%;
      min-height: 100%;
    }

    * {
      color: inherit;
    }
  }
  
  .saving-indicator {
    color: var(--site-theme-primary-color);
  }
  
  /* Ensure proper theming for the editor */
  .milkdown {
    background-color: var(--site-theme-bg-primary);
    color: var(--site-theme-text-primary) !important;
    
    .editor {
      background-color: var(--site-theme-bg-primary);
    }
    
    h1, h2, h3, h4, h5, h6 {
      color: var(--site-theme-text-primary);
    }
    
    a {
      color: var(--site-theme-primary-color);
    }
    
    blockquote {
      border-left-color: var(--site-theme-divider);
      background-color: var(--site-theme-bg-soft);
    }
    
    code {
      background-color: var(--site-theme-bg-soft);
    }
    
    pre {
      background-color: var(--site-theme-bg-secondary);
    }
    
    table {
      border-color: var(--site-theme-divider);
      
      th {
        background-color: var(--site-theme-bg-secondary);
        border-color: var(--site-theme-divider);
      }
      
      td {
        border-color: var(--site-theme-divider);
      }
    }
  }
}
</style>
