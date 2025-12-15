<template>
  <div class="wrapper py-4 px-6 rounded-sm">
    <template v-if="finished">
      <ul
        v-if="data.length"
        class="mb-4 pt-2 px-0 flex flex-col gap-3 list-none break-all"
      >
        <li
          v-for="item in data"
          class="flex items-center gap-6"
        >
          <span
            class="flex-1 text-sm text-center py-1.5 w-56"
            :class="{
              'text-muted': !item.checked,
              'bg-muted': !item.checked,
            }"
            :style="
              item.checked
                ? {
                  backgroundColor: item.backgroundColor,
                }
                : undefined
            "
          >{{ item.type }}</span>
          <a-switch
            :checked="item.checked"
            @change="toggleChecked(item.type, $event)"
          />
        </li>
      </ul>
      <p
        v-else
        class="my-2"
      >
        {{ $t('aiHighlighter.blank') }}
      </p>
      <a-button
        class="!p-0"
        type="link"
        @click="toggleFeedback"
      >
        {{
          $t('aiHighlighter.feedback')
        }}
      </a-button>
      <Help
        :visible="feedbackVisible"
        @cancel="toggleFeedback"
      >
        uid: {{ uid }}, traceid: {{ traceid }}
      </Help>
    </template>
    <template v-else>
      <h3 class="mb-2 text-inherit">
        {{ $t('aiHighlighter.popupTt') }}
      </h3>
      <ul class="mb-6 p-0 list-none">
        <li
          v-for="x in $t('aiHighlighter.popupCt').split('\n')"
          v-html="x"
        />
      </ul>
      <p class="flex items-center m-0">
        <span class="flex-1">{{ $t('viewer.fullTextTip1', { times }) }}</span>
        <span
          v-if="status === Status.LOADING"
          class="text-progress"
        >{{
          $t('aiHighlighter.progress', { progress })
        }}</span>
        <a-button
          v-else-if="(times ?? 0) > 0"
          @click="emit('start')"
        >
          {{
            $t('aiHighlighter.btn')
          }}
        </a-button>
        <span
          v-else
          class="text-error"
        >{{
          $t('aiHighlighter.limited')
        }}</span>
      </p>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import colorString from 'color-string';
import Help from '@/components/Common/Help.vue';
import { Status } from '@/api/aiHighlighter';
import { useStore } from '~/src/store';
import { HighlightsMap } from '~/src/hooks/useAIHighlighter';

const props = defineProps<{
  status: Status;
  progress: number;
  finished?: boolean;
  times?: number;
  traceid?: string;
  highlights?: HighlightsMap;
}>();

const emit = defineEmits<{
  (e: 'start'): void;
  (e: 'toggle', type: string, checked: boolean): void;
}>();

const store = useStore();

const feedbackVisible = ref(false);
const toggleFeedback = (v = !feedbackVisible.value) => {
  feedbackVisible.value = v;
};

const uid = computed(() => store.state.user.userInfo?.id);
const data = computed(
  () =>
    Object.values(props.highlights || {})?.map((x) => {
      const c = colorString.get.rgb(x.color);
      const v = colorString.to.hex(c.slice(0, 3), 0.3);
      return {
        ...x,
        backgroundColor: v,
      };
    }) ?? []
);
const toggleChecked = (type: string, v: boolean) => {
  emit('toggle', type, v);
};
</script>

<style scoped lang="less">
.wrapper {
  max-width: 340px;
  background-color: var(--site-theme-background);
  color: var(--site-theme-text-color);
  box-shadow: var(--site-theme-shadow);

  :deep(.ant-switch):not(.ant-switch-checked) {
    background-color: var(--site-theme-switch-unchecked, rgba(0, 0, 0, 0.25));
  }
}

.text-muted {
  color: var(--site-theme-text-secondary);
}

.bg-muted {
  background-color: var(--site-theme-background-secondary);
}

.text-progress {
  color: var(--site-theme-success-color, #3da611);
}

.text-error {
  color: var(--site-theme-error-color, #e66045);
}
</style>
