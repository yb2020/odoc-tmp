<template>
  <div
    :class="[
      'notes-container',
      tab === RightSideBarType.Group ? 'group-notes-container' : '',
    ]"
  >
    <Filter
      :activeColorMap="annotationStore.activeColorMap"
      :handleFilterChange="handleFilterChange"
      :tab="tab"
    />
    <div class="container">
      <PerfectScrollbar
        class="notes-scroller"
        @ps-scroll-y="handleScroll"
        @dblclick="handleEmptyClick(tab, pdfViewerRef)"
      >
        <EmptyComponent
          v-if="
            Object.values(notes)
              .filter(Boolean)
              .every((list) => (list as Array<any>).length === 0)
          "
          :message="
            !isOwner
              ? $t('message.noNotesTip')
              : activeTab === RightSideBarType.Group
                ? $t('teams.emptyNotesTip')
                : $t('message.doubleClickToCreateNote')
          "
          @dblclick="
            isOwner && safeInsertNoReferenceAnnotation($event, activeTab)
          "
        />
        <Notes
          v-for="(item, key) in notes"
          v-else
          :key="key"
          :notes="item"
          :pageNumber="String(key)"
          :tab="tab"
          :data-page-number="key"
        />
      </PerfectScrollbar>
      <!-- <Page
        :active-key="activeKey"
        :notes="notes"
        @setActiveKey="setActiveKey"
      /> -->
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from 'vue';
import { currentGroupId, currentNoteInfo, isOwner } from '@/store';
import Notes from './Notes.vue';
// import Page from './Page.vue';
import Filter from './Filter.vue';
import EmptyComponent from './common/Empty.vue';
import { arrow } from '~/src/dom/arrow';
import { RightSideBarType } from '../type';
import {
  safeInsertNoReferenceAnnotation,
  handleEmptyClick,
  checkNoReferenceAnnotation,
} from './annotation-state';
import { usePdfStore } from '~/src/stores/pdfStore';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { NoteFilter, ColorKey } from '~/src/style/select';
import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type';

const pdfStore = usePdfStore();
const annotationStore = useAnnotationStore();

const pdfViewerRef = computed(() => {
  return pdfStore.getViewer(currentNoteInfo.value?.pdfId);
});

const props = defineProps<{
  tab: RightSideBarType;
  activeTab: RightSideBarType;
}>();

const emit = defineEmits<{
  (e: 'counted', v: number): void;
}>();

watch(
  () => annotationStore.count,
  () => {
    emit('counted', annotationStore.count);
  }
);

const notes = computed(() => {
  if (props.activeTab !== props.tab) {
    return {};
  }

  const notes = { ...annotationStore.pageMap };

  Object.keys(notes).forEach((key) => {
    notes[key] = notes[key].filter((item) => {
      if (
        !(
          annotationStore.activeColorMap.ref ||
          item.idea ||
          checkNoReferenceAnnotation(item)
        ) ||
        (currentGroupId.value === SELF_NOTEINFO_GROUPID && item.isHighlight)
      ) {
        return false;
      }

      return annotationStore.activeColorMap[item.styleId as ColorKey];
    });
  });

  return notes;
});

const handleFilterChange = (filter: NoteFilter) => {
  annotationStore.activeColorMap[filter] =
    !annotationStore.activeColorMap[filter];
};

const activeKey = ref<number>(1);
let isClickPager = false;
// let timer: any;
// const setActiveKey = (pageNumber: number) => {
//   clearTimeout(timer);
//   isClickPager = true;
//   activeKey.value = pageNumber;
//   timer = setTimeout(() => {
//     isClickPager = false;
//   }, 1000);
// };

const changePager = (e: any) => {
  let heightTotal = 0;
  const { scrollTop, children } = e.target;
  const insideNode = [...children].find((node: HTMLElement) => {
    heightTotal += node.offsetHeight;

    return heightTotal > scrollTop;
  });
  const pageNumber = insideNode.getAttribute('data-page-number');

  activeKey.value = parseInt(pageNumber);
};

const handleScroll = (e: any) => {
  !isClickPager && changePager(e);

  if (arrow.line) {
    arrow.line.position();
  }
};
</script>

<style lang="less" scoped>
.notes-container {
  height: 100%;
  position: relative;

  .container {
    display: flex;
    height: calc(100% - 33px);
  }

  .notes-scroller {
    padding-left: 16px;
    height: 100%;
    flex: 1 1 auto;
  }
}

.group-notes-container {
  height: calc(100% - 49px);
  margin-top: 0;

  .container {
    display: flex;
    height: 100%;
  }
}
</style>
