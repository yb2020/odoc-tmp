<template>
  <div
    class="filter-container"
    :class="{
      'filter-container-group': tab === RightSideBarType.Group,
    }"
  >
    <template v-if="tab === RightSideBarType.Note">
      <div
        v-if="isVisitor"
        class="note-owner-info"
      >
        笔记属于
        <span
          class="note-owner-link"
          @click="goPersonPage(noteUserInfo!.id)"
        >
          <a-avatar
            v-if="noteUserInfo?.avatarUrl"
            :src="noteUserInfo!.avatarUrl"
            class="avatar"
          />
          {{ noteUserInfo?.nickName ?? '' }}
        </span>
      </div>
      <div
        v-else
        class="absolute left-4"
      >
        <MarkdownTip />
      </div>
    </template>
    <div
      class="name"
      :class="{ 'not-all': !all }"
      @click="handleClick"
    >
      <FunnelPlotOutlined />
    </div>
    <template v-if="showOptions">
      <div
        v-show="showOptions"
        ref="optionsRef"
        class="filter-color-options"
      >
        <div
          v-for="(option, key) in styleMap"
          :key="option.color"
          :value="option.color"
          :class="['option', activeColorMap[key] ? 'active' : '']"
          @click="handleFilterChange(key)"
        >
          <div class="text-container">
            <div
              class="dot"
              :style="{ backgroundColor: option.color }"
            />
            {{ $t(option.i18n) }}
          </div>

          <check-outlined class="iconcheck" />
        </div>
        <div
          value="ref"
          :class="['option', activeColorMap['ref'] ? 'active' : '']"
          @click="handleFilterChange('ref')"
        >
          {{ $t('viewer.selectedText') }}
          <check-outlined class="iconcheck" />
        </div>
      </div>
    </template>
  </div>
</template>

<script lang="ts" setup>
import { UserStatusEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/OpenNote';
import { NoteFilter, styleMap } from '@/style/select';
import { onClickOutside } from '@vueuse/core';
import { computed, ref } from 'vue';
import { RightSideBarType } from '../type';
import { FunnelPlotOutlined, CheckOutlined } from '@ant-design/icons-vue';
import { pdfStatusInfo, selfNoteInfo } from '~/src/store';
import { goPersonPage } from '~/src/common/src/utils/url';
import MarkdownTip from '@/components/Common/MarkdownTip.vue';

const isVisitor = computed(() => {
  const { noteUserStatus } = pdfStatusInfo.value;
  return (
    noteUserStatus === UserStatusEnum.TOURIST ||
    noteUserStatus === UserStatusEnum.GUEST
  );
});
const noteUserInfo = computed(() => {
  return selfNoteInfo.value?.userInfo;
});

const showOptions = ref(false);
const handleClick = () => {
  showOptions.value = !showOptions.value;
};

const props = defineProps<{
  activeColorMap: Record<NoteFilter, boolean>;
  handleFilterChange: (filter: NoteFilter) => void;
  tab: RightSideBarType;
}>();

const optionsRef = ref();

const all = computed(() => {
  return Object.values(props.activeColorMap).every(Boolean);
});

onClickOutside(optionsRef, () => {
  showOptions.value = false;
});
</script>

<style lang="postcss" scoped>
.filter-container {
  margin-right: 10px;
  padding: 0 10px 9px 0;
  position: relative;
  display: flex;
  justify-content: flex-end;

  .name {
    width: 24px;
    height: 24px;
    border-radius: 2px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--site-theme-pdf-panel-text);
    font-size: 16px;
    cursor: pointer;

    &:hover {
      background: var(--site-theme-bg-hover);
    }

    &.not-all {
      background: var(--site-theme-bg-hover);
      &:hover {
        background: var(--site-theme-bg-hover-deep);
      }
    }
  }

  .filter-color-options {
    position: absolute;
    bottom: 0;
    transform: translate(0, 100%);
    z-index: 3;
    background: var(--site-theme-pdf-panel-secondary);
    font-size: 13px;
    font-weight: 400;
    color: var(--site-theme-text-secondary);
    cursor: pointer;

    right: 0;
    margin-right: 12px;

    .option {
      width: 126px;
      height: 32px;
      display: flex;
      align-items: center;
      padding: 0 9px 0 15px;
      justify-content: space-between;
      background: var(--site-theme-bg-hover);

      &.active {
        color: var(--site-theme-text-primary);

        .iconcheck {
          display: block;
        }
      }

      .text-container {
        display: flex;
        align-items: center;

        .dot {
          width: 10px;
          height: 10px;
          border-radius: 50%;
          margin: 0 4px 0;
        }
      }
    }

    .iconcheck {
      font-size: 12px;
      color: var(--site-theme-text-primary);
      display: none;
    }
  }

  .note-owner-info {
    position: absolute;
    top: 0;
    left: 18px;
    color: var(--site-theme-text-secondary);
    font-size: 13px;
    .note-owner-link {
      margin-left: 10px;
      cursor: pointer;
      .avatar {
        margin-right: 3px;
        height: 24px;
        width: 24px;
      }
    }
  }
}

.filter-container-group {
  position: absolute;
  top: -8px;
  right: 8px;
  transform: translate(0, -100%);
  padding: 0 8px;
  z-index: 999;
  margin-right: 0;
}
</style>
