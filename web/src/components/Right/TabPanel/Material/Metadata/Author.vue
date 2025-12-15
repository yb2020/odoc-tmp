<template>
  <div class="metadata-author-container">
    <div
      ref="triggerRef"
      class="metadata-author-light"
      :style="{
        visibility: ready ? 'visible' : 'hidden',
        cursor: props.docId ? 'pointer' : 'initial',
      }"
    >
      {{ authorsView }}
    </div>
    <div
      ref="shadowRef"
      class="metadata-author-shadow"
    >
      {{ authorsView }}
    </div>
    <TippyVue
      v-if="triggerRef && props.docId"
      ref="tippyRef"
      :trigger-ele="triggerRef"
      placement="left-start"
      :offset="[0, 20]"
      trigger="click"
      @onShown="docId = props.docId!"
      @onHide="docId = ''"
    >
      <AuthorEditVue
        :doc-id="docId"
        :with-author-link="withAuthorLink"
        @cancel="cancel()"
        @update="update($event)"
      />
    </TippyVue>
    <Rollback
      v-if="props.displayAuthors?.rollbackEnable"
      :current-info="props.displayAuthors?.authors.join(' / ')"
      :origin-info="props.displayAuthors?.originAuthors.join(' / ')"
      width="500px"
      @click="rollback()"
    />
  </div>
</template>

<script setup lang="ts">
import { DisplaySimpleAuthors } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/doc/UserDoc';
import { DocDetailInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/client/ClientDoc';
import { computed, nextTick, onMounted, ref, watch } from 'vue';
import { updateAuthors } from '~/src/api/material';
import TippyVue from '../../../../Tippy/index.vue';
import AuthorEditVue from './AuthorEdit.vue';
import Rollback from './Rollback.vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  docId?: string;
  displayAuthors?: DocDetailInfo['displayAuthor'];
  withAuthorLink: boolean;
}>();

const safeAuthors = computed(() => {
  return props.displayAuthors?.authors ?? [];
});
const initMiddleNumber = computed(() => {
  return Math.max(0, safeAuthors.value.length - 2);
});
const middleNumber = ref(initMiddleNumber.value);
const ready = ref(false);
const authors = computed(() => {
  const [head, ...rest] = safeAuthors.value;
  const tail = rest.pop();
  const after: string[] = [];
  if (head) {
    after.push(head);
  }

  after.push(...rest.slice(0, middleNumber.value));

  if (middleNumber.value < rest.length) {
    after.push('…');
  }

  if (tail) {
    after.push(tail);
  }

  return after;
});

const { t } = useI18n();

const authorsView = computed(() => {
  return authors.value.join(' / ') || t('message.noAuthorDataTip');
});

const sync = async () => {
  await nextTick();
  if (
    shadowRef.value!.offsetHeight > triggerRef.value!.offsetHeight &&
    middleNumber.value > 0
  ) {
    middleNumber.value -= 1;
  } else {
    ready.value = true;
  }
};

watch(safeAuthors, () => {
  ready.value = false;
  middleNumber.value = initMiddleNumber.value;
});
watch(authors, sync);
onMounted(sync);

const triggerRef = ref<HTMLDivElement | null>(null);
const shadowRef = ref<HTMLDivElement | null>(null);
const tippyRef = ref();
const docId = ref('');
const cancel = () => {
  tippyRef.value.hide();
  docId.value = '';
};

const update = (newAuthor: DisplaySimpleAuthors) => {
  props.displayAuthors!.authors = newAuthor!.authors.map((item) => item.name);
  props.displayAuthors!.rollbackEnable = newAuthor!.rollbackEnable;
};

const rollback = async () => {
  const newAuthor = await updateAuthors({
    docId: props.docId!,
  });
  update(newAuthor!);
};
</script>

<style lang="less" scoped>
.metadata-author-container {
  position: relative;
  margin-bottom: 12px;
  margin-left: -8px;
  margin-right: -8px;
  padding-left: 8px;
  padding-right: 8px;
  padding-top: 4px;
  padding-bottom: 4px;
  border-radius: 2px;
  .metadata-author-light,
  .metadata-author-shadow {
    font-size: 13px;
    color: var(--site-theme-text-secondary);
    word-break: break-all;
  }

  .metadata-author-light {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .metadata-author-shadow {
    pointer-events: none;
    opacity: 0;
    position: absolute;
  }
}
</style>
<style lang="less">
.metadata-author-container {
  &:hover {
    background: rgba(255, 255, 255, 0.08);
    .metadata-rollback {
      top: 50% !important;
      transform: translateY(-50%);
      visibility: visible !important;
    }
  }
}
</style>
