<template>
  <div
    ref="wrapper"
    class="relative h-full flex flex-col px-4"
  >
    <div class="w-full min-w-[960px] max-w-[1192px] mx-auto">
      <section class="banner mt-6">
        <a
          class="banner-link h-24 w-full flex items-center justify-center"
          @click.prevent="openBeginnerGuide(isCurrentLanguage(Language.EN_US))"
        >
          <span class="text-white text-2xl tracking-widest">{{ $t('common.aitools.guide') }}<ArrowRightOutlined /></span>
        </a>
      </section>
      <section class="operations mt-6 flex items-center gap-6">
        <BtnAdd
          :isOpenNewTab="openInNewTab"
          @preclick="precheck"
          @added="onAdded"
          @errvip="onVipErr"
        />
        <BtnUpload
          always-new
          @preclick="precheck"
          @import:success="onImported"
          @errvip="onVipErr"
        />
      </section>
    </div>
    <section class="projects flex flex-col flex-1 min-h-0">
      <h2 class="w-full min-w-[960px] max-w-[1192px] mx-auto text-lg mt-6 mb-6">
        {{ $t('common.aitools.rencetTt') }}
        <RedoOutlined
          class="cursor-pointer"
          @click="handleRefresh"
        />
      </h2>
      <div class="flex-1 min-h-0 psm overflow-auto">
        <div class="min-w-[960px] max-w-[1192px] mx-auto">
          <Projects
            :loading="loading"
            :projects="rencentProjects"
            :openInNewTab="openInNewTab"
            @deleted="onDeleted"
          />
        </div>
      </div>
    </section>
    <VipGuard
      v-if="wrapper"
      v-model:visible="guardVisible"
      :closable="true"
      :container="wrapper"
      wrapClassName="!absolute"
      :maskStyle="{
        position: 'absolute',
      }"
    />
  </div>
</template>

<script setup lang="ts">
import { isNil } from 'lodash';
import { SimpleProjectInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText';
import { useRouter } from 'vue-router';
import { ArrowRightOutlined } from '@ant-design/icons-vue';
import { useProjects } from '@common/hooks/aitools/useProjects';
import BtnAdd from './BtnAdd/index.vue';
import BtnUpload from './BtnUpload/index.vue';
import Projects from './Projects/index.vue';
import VipGuard from './VipGuard/index.vue';
import { openProjectPage, openBeginnerGuide } from './utils';
import { computed, ref } from 'vue';
import { RedoOutlined } from '@ant-design/icons-vue';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '@/hooks/useLanguage';

const props = defineProps<{
  loading: boolean;
  projects: SimpleProjectInfo[];
  openInNewTab?: boolean;
  precheck: (e: MouseEvent) => void;
}>();

const emit = defineEmits<{
  (e: 'added', id?: string): void;
  (e: 'deleted', id: string): void;
  (e: 'refresh'): void;
}>();

const wrapper = ref<HTMLElement>();
const guardVisible = ref(false);

// 语言管理
const { isCurrentLanguage } = useLanguage();

const router = useRouter();
const isFetchSelf = computed(() => isNil(props.projects));
const rencentProjects = computed(() =>
  isFetchSelf.value ? selfProjects.value || [] : props.projects
);

const { data: selfProjects, refresh } = useProjects(isFetchSelf);

const onVipErr = () => {
  guardVisible.value = true;
};

const onAdded = (id: string) => {
  if (isFetchSelf.value) {
    // 后台数据同步没那么快
    setTimeout(refresh, 500);
  }

  emit('added', id);
};

const onImported = (id: string) => {
  openProjectPage(
    router,
    {
      id,
    },
    props.openInNewTab ? 'resolve' : 'push'
  );

  onAdded(id);
};

const onDeleted = (id: string) => {
  if (isFetchSelf.value) {
    refresh();
  }

  emit('deleted', id);
};

const handleRefresh = () => {
  if (isFetchSelf.value) {
    refresh();
  } else {
    emit('refresh');
  }
};
</script>

<style lang="less">
.banner-link {
  border-radius: 16px;
  background: url('~common/assets/images/aitools/bg-guide.png') no-repeat center;
  background-size: cover;
}
</style>
