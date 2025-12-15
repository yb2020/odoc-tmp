<template>
  <a-spin :spinning="loading">
    <ul
      v-if="projects?.length"
      class="flex flex-wrap gap-4"
    >
      <li
        v-for="item in projects"
        class="group relative w-56 h-[280px] p-4 bg-white rounded-2xl border border-white hover:border-rp-blue-6"
      >
        <div
          class="h-full flex flex-col cursor-pointer"
          @click="handleOpen(item.projectId)"
        >
          <h3
            class="text-base text-rp-neutral-10 overflow-hidden whitespace-nowrap text-ellipsis"
          >
            {{ item.simpleTitle || $t('common.aitools.untitled') }}
          </h3>
          <p
            class="flex-1 min-h-0 mt-3 text-sm text-rp-neutral-8 line-clamp-[9]"
          >
            {{ item.simpleText }}
          </p>
          <time class="text-rp-neutral-8">{{
            formatRecentDate(parseInt(item.version!, 10))
          }}</time>
          <aside class="hidden group-hover:block absolute right-4 bottom-4">
            <a-button
              :loading="deletingMap[item.projectId]"
              type="text"
              size="small"
              class="!p-0"
              @click.stop.prevent="onDelete(item)"
            >
              <DeleteOutlined class="hover:text-rp-blue-6" />
            </a-button>
          </aside>
        </div>
      </li>
    </ul>
    <a-result
      v-else
      :sub-title="$t('common.aitools.recentEmpty')"
      class="!pt-28 text-rp-neutral-8"
    >
      <template #icon>
        <img
          :src="SvgEmpty"
          alt="empty"
          class="inline"
        >
      </template>
    </a-result>
  </a-spin>
</template>

<script setup lang="ts">
import { useRequest } from 'ahooks-vue';
import { DeleteOutlined } from '@ant-design/icons-vue';
import { Modal } from 'ant-design-vue';
import SvgEmpty from '~common/assets/images/aitools/empty.svg';
import { deleteProject } from '@common/api/revise';
import { formatRecentDate } from '@common/utils/format';
import { useI18n } from 'vue-i18n';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { SimpleProjectInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText';
import { openProjectPage } from '../utils';

const props = defineProps<{
  loading: boolean;
  projects: SimpleProjectInfo[];
  openInNewTab?: boolean;
}>();

const emit = defineEmits<{
  (e: 'deleted', id: string): void;
}>();

const { t } = useI18n();

const { run: doDeleteProject } = useRequest(
  async ({ projectId }: { projectId: string }) => {
    try {
      deletingMap.value[projectId] = true;
      await deleteProject({
        projectId,
      });
      emit('deleted', projectId);
    } finally {
      deletingMap.value[projectId] = false;
    }
  },
  {
    manual: true,
  }
);

const deletingMap = ref(
  {} as {
    [k: string]: boolean;
  }
);
const onDelete = ({ projectId }: { projectId: string }) => {
  Modal.confirm({
    title: t('common.text.tips'),
    content: t('common.aitools.recentDeleteTip'),
    onOk() {
      doDeleteProject({
        projectId,
      });
    },
  });
};

const router = useRouter();

const handleOpen = (projectId: string) => {
  openProjectPage(
    router,
    { id: projectId },
    props.openInNewTab ? 'resolve' : 'push'
  );
};
</script>
