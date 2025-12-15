<template>
  <div class="metadata-publishDate-container">
    <Popover
      title=""
      overlay-class-name="metadata-popover"
      :visible="visible"
      trigger="click"
      placement="left"
    >
      <template #content>
        <div
          class="list-edit-publishDate"
          @click.stop
        >
          <CloseOutlined @click="cancelPublishDate()" />
          <div class="list-edit-publishDate-header">
            {{ $t('viewer.publishDate') }}
          </div>
          <div class="list-edit-publishDate-body">
            <InputNumber
              v-model:value="year"
              :max="currentYear"
              class="metadata-input"
              size="small"
              :placeholder="$t('viewer.year')"
            />
            <InputNumber
              v-model:value="month"
              :min="1"
              :max="12"
              class="metadata-input"
              size="small"
              :placeholder="$t('viewer.month')"
            />
            <InputNumber
              v-model:value="day"
              :min="1"
              :max="31"
              class="metadata-input"
              size="small"
              :placeholder="$t('viewer.date')"
            />
          </div>
          <div class="list-edit-publishDate-footer">
            <button @click="cancelPublishDate()">
              {{ $t('viewer.cancel') }}
            </button>
            <button
              :disabled="publishDateDisabled"
              @click="submitPublishDate()"
            >
              {{ $t('viewer.confirm') }}
            </button>
          </div>
        </div>
      </template>
      <div
        class="metadata-field"
        style="margin-left: -6px; padding-left: 6px"
        @click.stop
      >
        <div class="metadata-field-name">
          {{ $t('viewer.dateLabel') }}
        </div>
        <div
          class="metadata-publishDate-view"
          :style="{ cursor: props.docId ? 'pointer' : 'initial' }"
          @click="editPublishDate()"
        >
          {{ viewDate }}
          <Rollback
            v-if="props.docId && props.displayPublishDate?.rollbackEnable"
            :current-info="props.displayPublishDate?.publishDate"
            :origin-info="props.displayPublishDate?.originPublishDate"
            width="240px"
            @click="rollback()"
          />
        </div>
      </div>
    </Popover>
  </div>
</template>

<script setup lang="ts">
import { Popover, InputNumber } from 'ant-design-vue';
import { CloseOutlined } from '@ant-design/icons-vue';
import { DocDetailInfo } from 'go-sea-proto/gen/ts/doc/ClientDoc';
import { ref, computed } from 'vue';
import { updatePublishDate } from '~/src/api/material';
import Rollback from './Rollback.vue';
import { useI18n } from 'vue-i18n';
import { getCurrentLanguage } from '~/src/shared/i18n/core'

const currentYear = new Date().getFullYear();

const props = defineProps<{
  docId?: string;
  displayPublishDate?: DocDetailInfo['displayPublishDate'];
}>();

const { d } = useI18n()

// 自定义日期格式化函数，使用项目的locale转换系统
const formatDateWithCorrectLocale = (date: Date, format: string = 'short'): string => {
  try {
    // getCurrentLanguage() 现在直接返回标准格式（'en-US', 'zh-CN'）
    const standardLocale = getCurrentLanguage();
    
    console.log('PublishDate - 使用标准格式locale:', {
      标准locale: standardLocale
    });

    // 根据format类型使用不同的格式化选项
    const formatOptions: Intl.DateTimeFormatOptions = format === 'short' 
      ? { year: 'numeric', month: 'short', day: 'numeric' }
      : { year: 'numeric', month: '2-digit', day: '2-digit' };

    const formatter = new Intl.DateTimeFormat(standardLocale, formatOptions);
    return formatter.format(date);
  } catch (error) {
    console.error('PublishDate - 自定义日期格式化失败:', error);
    // 降级处理：返回简单的日期字符串
    return date.toLocaleDateString();
  }
};

const viewDate = computed(() => {
  if (!props.displayPublishDate?.publishDate) {
    return '';
  }

  const dateObject = new Date(props.displayPublishDate.publishDate.replaceAll('-', '/'));
  return formatDateWithCorrectLocale(dateObject, 'short');
});

const visible = ref(false);
const year = ref<number | undefined>();
const month = ref<number | undefined>();
const day = ref<number | undefined>();

const editPublishDate = () => {
  if (!props.docId) {
    return;
  }

  const parts = props.displayPublishDate!.publishDate.split('-');
  const nanToUndefined = (string?: string) => {
    if (!string) {
      return undefined;
    }

    return Number(string);
  };
  year.value = nanToUndefined(parts.shift());
  month.value = nanToUndefined(parts.shift());
  day.value = nanToUndefined(parts.shift());
  visible.value = true;
  document.body.addEventListener('click', cancelPublishDate, { once: true });
};
const cancelPublishDate = () => {
  document.body.removeEventListener('click', cancelPublishDate);
  visible.value = false;
  year.value = currentYear;
  month.value = undefined;
  day.value = undefined;
};
const submitPublishDate = async () => {
  const { newPublishDate } = await updatePublishDate({
    docId: props.docId,
    publishDate: [year.value, month.value, day.value]
      .filter(Boolean)
      .map((number) => String(number).padStart(2, '0'))
      .join('-'),
  });

  props.displayPublishDate!.publishDate = newPublishDate!.publishDate;
  props.displayPublishDate!.rollbackEnable = true;
  cancelPublishDate();
};

const publishDateDisabled = computed(() => {
  if (day.value) {
    if (!month.value || !year.value) {
      return true;
    }

    const date = new Date();
    date.setFullYear(year.value);
    date.setMonth(month.value - 1);
    date.setDate(day.value);

    return date.getMonth() !== month.value - 1;
  } else if (month.value) {
    return !year.value;
  } else {
    return false;
  }
});

const rollback = async () => {
  props.displayPublishDate!.publishDate =
    props.displayPublishDate!.originPublishDate;
  props.displayPublishDate!.rollbackEnable = false;

  await updatePublishDate({
    docId: props.docId!,
  });
};
</script>

<style lang="less">
@import './style.less';

.metadata-publishDate-container {
  flex: 1 1 100%;
  border-radius: 2px;
  .metadata-publishDate-view,
  input {
    height: 32px;
    line-height: 32px;
    width: 100%;
  }
  .metadata-publishDate-view {
    position: relative;
    padding-left: 8px;
    color: var(--site-theme-text-secondary);
    border-radius: 2px;
    &:hover {
      background: rgba(255, 255, 255, 0.08);
      > div {
        visibility: visible;
      }
    }
  }
  input {
    text-indent: 6px;
  }
}

.list-edit-publishDate {
  > .anticon-close {
    position: absolute;
    right: 20px;
    color: var(--site-theme-text-secondary);
  }
  .list-edit-publishDate-header {
    color: var(--site-theme-text-primary);
    font-weight: bold;
    height: 22px;
    margin-bottom: 18px;
  }
  .list-edit-publishDate-body {
    width: 220px;
    display: flex;
    justify-content: space-between;
    margin-bottom: 20px;
    > * {
      width: 66px;
    }
  }
  .list-edit-publishDate-footer {
    display: flex;
    justify-content: flex-end;
    button {
      margin-left: 12px;
      width: 64px;
      height: 24px;
      display: flex;
      align-items: center;
      justify-content: center;
      border: 0;
      outline: 0;
      border-radius: 2px;
      font-size: 12px;
      cursor: pointer;
      &:first-of-type {
        background: var(--site-theme-bg-light);
        color: var(--site-theme-text-secondary);
      }
      &:last-of-type {
        background: var(--site-theme-brand);
        color: #fff;
      }
      &:disabled {
        opacity: 0.7;
      }
    }
  }
}

[data-theme="dark"] {
  .list-edit-publishDate {
    .list-edit-publishDate-header {
      color: #ffffffd9 !important;
    }
    .list-edit-publishDate-footer {
      button {
        &:first-of-type {
          background: var(--site-theme-bg-primary) !important;
        }
      }
    }
  }
}
</style>
