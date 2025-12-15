<template>
  <div
    v-if="notes && notes.length"
    class="notes"
    :style="{
      paddingTop:
        annotationStore.headTailPageNumber.head === Number(pageNumber)
          ? '9px'
          : '1px',
    }"
  >
    <div
      :class="['page', `page-${pageNumber}`]"
      :title="tab === RightSideBarType.Note ? '双击此处可以添加笔记' : ''"
      @dblclick="safeInsertNoReferenceAnnotation($event, tab, Number(pageNumber), 0)"
    >
      {{ $t('viewer.page', { num: $n(parseInt(pageNumber), 'integer') }) }}
    </div>
    <div
      class="note-prev"
      :title="tab === RightSideBarType.Note ? '双击此处可以添加笔记' : ''"
      @dblclick="safeInsertNoReferenceAnnotation($event, tab, Number(pageNumber), 0)"
    >
      <div
        :style="{ background: isInsertToFirst ? 'var(--site-theme-brand)' : 'transparent' }"
      />
    </div>
    <template v-if="tab === RightSideBarType.Note">
      <Item
        v-for="(item, index) in notes"
        :key="item.uuid"
        :note="item"
        :index="index"
        :note-list="notes"
        :isOwner="isOwner"
      />
    </template>
    <template v-if="tab === RightSideBarType.Group">
      <GroupItem
        v-for="item in notes"
        :key="item.uuid"
        :note="item"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { AnnotationAll } from '~/src/stores/annotationStore/BaseAnnotationController';
import { RightSideBarType } from '../type';
import Item from './Item.vue';
import GroupItem from './group/Item.vue';
import {
  safeInsertNoReferenceAnnotation,
  annotationDragHovering,
  annotationDragPosition,
} from './annotation-state';
import { computed } from 'vue';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { isOwner } from '~/src/store';

const annotationStore = useAnnotationStore();
const props = defineProps<{
  pageNumber: string;
  notes: AnnotationAll[];
  tab: RightSideBarType;
}>();

const isInsertToFirst = computed(() => {
  return (
    props.notes.length > 0 &&
    annotationDragHovering.value === props.notes[0].uuid &&
    annotationDragPosition.value === 0
  );
});
</script>

<style scoped lang="less">
.notes {
  width: 100%;
  padding-right: 16px;

  .page {
    font-family: 'Noto Sans SC';
    font-style: normal;
    font-weight: 400;
    font-size: 13px;
    line-height: 20px;
    color: var(--site-theme-pdf-panel-text);
    display: flex;
    align-items: flex-end;
    padding-bottom: 2px;
    opacity: 0.8;
  }

  .note-prev {
    height: 10px;
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    div {
      flex: 0 0 96%;
      height: 3px;
      pointer-events: none;
    }
  }
}
</style>
