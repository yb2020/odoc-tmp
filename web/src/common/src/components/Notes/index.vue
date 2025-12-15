<template>
  <div class="notes-container h-full">
    <main ref="mainEl" class="main flex w-full" :class="[NOTES_CLASSNAME]">
      <!-- 左侧树型 -->
      <Drawer v-model:visible="visible" class="left" :closable="collapsable" :initial-width="248" :max-width="248"
        :min-width="248" placement="left">
        <template #icon="iconProps">
          <slot name="icon" v-bind="iconProps" />
        </template>
        <a-tabs v-model:activeKey="activeKey">
          <a-tab-pane v-for="item in TAB_LIST" :key="item.key" :tab="$t(item.title)">
            <!-- 左侧树型 -->
            <NoteFolder isInClient :active-type="activeKey" :noteState="noteState" />
            <!-- 右侧展示栏 -->
            <teleport v-if="elPane" :to="elPane">
              <div class="h-full flex flex-col">
                <LoginFirst :isInClient="isInClient">
                  <template #blank>
                    <img class="w-20" src="@common/../assets/images/notes/empty-notes.svg" alt="empty" />
                    <p>
                      {{
                      $t('common.tips.empty', [
                      $t(`${activeTab?.title}`).toLocaleLowerCase(),
                      ])
                      }}
                    </p>
                  </template>
                  <component :is="TAB_COMPONENTS[activeKey]" :activeType="activeKey" :noteState="noteState" />
                </LoginFirst>
              </div>
            </teleport>
          </a-tab-pane>
        </a-tabs>
      </Drawer>
      <section id="js-right" ref="elPane" class="right" />
    </main>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useNote } from './useNote'
import NoteFolder from './NoteFolder.vue'
import Summarize from './components/Summarize.vue'
import SummarizeTabs from './components/SummarizeTabs.vue'
import Word from './components/Word.vue'
import Extract from './components/Extract.vue'
import { NoteSubType2ModuleType, NoteSubTypes } from './types'
import Drawer from '@common/components/Drawer/index.vue'
import LoginFirst from '@common/components/LoginFirst/index.vue'
import { getSummaryList, getWordList, getExtractList } from '~/src/api/note'
import { reportModuleImpression, PageType } from '@common/utils/report'
import { useElementVisibility } from '@vueuse/core'
import { NOTES_CLASSNAME } from './useNote'

const props = defineProps<{
  isInClient?: boolean
  collapsable?: boolean
}>()

const { t } = useI18n()

const TAB_LIST = [
  {
    key: NoteSubTypes.Summary,
    title: 'common.notes.summary',
  },
  {
    key: NoteSubTypes.Vocabulary,
    title: 'common.notes.vocabulary',
  },
  {
    key: NoteSubTypes.Annotation,
    title: 'common.notes.annotation',
  },
]

const TAB_FUNCS = {
  [NoteSubTypes.Summary]: getSummaryList,
  [NoteSubTypes.Vocabulary]: getWordList,
  [NoteSubTypes.Annotation]: getExtractList,
}

const TAB_COMPONENTS: Record<string, any> = {
  [NoteSubTypes.Summary]: props.isInClient ? SummarizeTabs : Summarize,
  [NoteSubTypes.Vocabulary]: Word,
  [NoteSubTypes.Annotation]: Extract,
}

const mainEl = ref()
const visible = ref(true)
const elPane = ref<HTMLElement>()
const activeKey = ref(NoteSubTypes.Summary)
const activeTab = computed(() =>
  TAB_LIST.find((x) => x.key === activeKey.value)
)
const activeFun = computed(() => TAB_FUNCS[activeKey.value])
const noteState: ReturnType<typeof useNote> = useNote(
  activeFun,
  `${t('common.text.all')}${t(activeTab.value?.title!)}`,
  activeKey,
  props.isInClient
)

watch(
  activeKey,
  () => {
    noteState.noteFolderList.value = [noteState.noteAllFolder]
    noteState.refreshNoteExplorer()
    reportModuleImpression({
      page_type: PageType.NOTE_TAB,
      module_type: NoteSubType2ModuleType[activeKey.value],
    })
  },
  { immediate: true }
)

const isMainVisible = useElementVisibility(mainEl)

const onVisibilityChange = () => {
  // 切换到笔记管理页时刷新数据
  if (document.visibilityState === 'visible' && isMainVisible.value) {
    noteState.refreshNoteExplorer()
    if (activeKey.value === NoteSubTypes.Annotation) {
      noteState.fetchUserTagList()
    }
  }
}

onMounted(() => {
  document.addEventListener('visibilitychange', onVisibilityChange, true)
})

onUnmounted(() => {
  document.removeEventListener('visibilitychange', onVisibilityChange, true)
})
</script>

<style lang="less" scoped>
.notes-container {
  position: relative;
  display: flex;
  flex-direction: column;

  main {
    overflow: hidden;
    flex: 1 1 100%;
  }
}

:deep(.main) {
  .ant-tabs {
    display: flex;
    height: 100%;
  }

  .ant-tabs-content-holder {
    flex: 1;
  }

  .ant-tabs-content {
    height: 100%;
  }

  .ant-tabs-ink-bar {
    width: 0 !important;
  }

  .ant-tabs-tab-active .ant-tabs-tab-btn {
    color: var(--site-theme-primary-color);
  }
}

:deep(.left) {
  .ant-tabs-nav {
    height: 40px;
    margin: 0;

    .ant-tabs-nav-wrap,
    .ant-tabs-nav-list {
      width: 100%;
      flex: 1 0 100%;

      .ant-tabs-tab {
        flex: 1;
        margin: 0;
        color: var(--site-theme-text-color);
        background: var(--site-theme-background-secondary);

        .ant-tabs-tab-btn {
          width: 100%;
          text-align: center;
        }
      }

      .ant-tabs-tab-active {
        background: var(--site-theme-background-tertiary);
      }
    }
  }
}

:deep(.right) {
  flex: 1;

  .ant-tabs-nav {
    background: var(--site-theme-background-secondary);
    height: 32px;
    margin: 0;

    &::before {
      border-color: var(--site-theme-background-secondary);
    }

    .ant-tabs-nav-wrap,
    .ant-tabs-nav-list {
      width: 100%;

      .ant-tabs-tab {
        padding: 6px 20px 6px 16px;
        line-height: 20px;
        color: var(--site-theme-text-color);
        background: var(--site-theme-background-tertiary);

        .ant-tabs-tab-btn {
          width: 100%;
          text-align: center;
        }

        &+.ant-tabs-tab {
          margin-left: 2px;
        }
      }

      .ant-tabs-tab-active {
        background: var(--site-theme-background-primary);
      }
    }
  }
}
</style>
