<template>
  <div
    class="readstatus-wrapper h-fit mr-2 self-start"
    @mouseleave="onUnselecting"
  >
    <Dropdown
      arrow
      placement="bottom"
      overlay-class-name="readstatus-dropdown"
      :visible="selecting"
      :trigger="['click']"
      @visibleChange="onSelecting"
    >
      <template #overlay>
        <div @mouseenter="onSelecting" @mouseleave="onUnselecting">
          <div class="ant-popover-arrow"></div>
          <Menu @select="onSelect">
            <Item v-for="item in menuItems" :key="item">
              <div
                class="readstatus-dot"
                :style="{ backgroundColor: ReadingStatus2BadgeStatus[item] }"
              ></div>
            </Item>
          </Menu>
        </div>
      </template>
      <Tooltip placement="top" :title="badgeTxt" :arrow-point-at-center="true">
        <Progress
          v-if="stat === DocReadingStatus.READING"
          type="circle"
          size="small"
          :width="28"
          :percent="progress"
        >
          <template #format>
            <div
              class="readstatus-dot"
              :style="{ backgroundColor: ReadingStatus2BadgeStatus[stat] }"
            ></div>
          </template>
        </Progress>
        <div v-else class="w-7 h-7 flex items-center justify-center">
          <div
            class="readstatus-dot"
            :style="{ backgroundColor: ReadingStatus2BadgeStatus[stat] }"
          ></div>
        </div>
      </Tooltip>
    </Dropdown>
  </div>
</template>

<script lang="ts" setup>
import { PropType, computed, ref } from 'vue'
import { Progress, Tooltip, Menu, Dropdown } from 'ant-design-vue'
import {
  DocReadingStatus,
  UpdateReadStatusRequest,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/UserDocManage'
import type { BaseUseRequestOptions as Options } from 'ahooks-vue/dist/src/useRequest/types'
import { useRequest } from 'ahooks-vue'
import { updateDocReadStatus } from '~/src/api/document'
import { useI18n } from 'vue-i18n'

const { Item } = Menu

const ReadingStatus2BadgeStatus = {
  [DocReadingStatus.READ]: '#7BCC0A',
  [DocReadingStatus.UNREAD]: '#E66045',
  [DocReadingStatus.READING]: '#1F71E0',
  [DocReadingStatus.UNRECOGNIZED]: '#7BCC0A',
} as const

const menuItems = [
  DocReadingStatus.READ,
  DocReadingStatus.READING,
  DocReadingStatus.UNREAD,
]

const props = defineProps({
  pdfId: {
    type: String,
    default: '',
  },
  status: {
    type: Number as PropType<DocReadingStatus>,
    default: DocReadingStatus.UNRECOGNIZED,
  },
  percent: {
    type: Number,
    default: 1,
  },
})

const emit = defineEmits(['updated'])

const i18n = useI18n()

const selecting = ref(false)

const {
  status: stat,
  progress,
  loading: updating,
  run: updateReadStatus,
} = useDocReadStatus(props.status, props.percent, {
  onSuccess() {
    emit('updated', {
      status: stat.value,
      progress: progress.value,
    })
  },
})

const badgeTxt = computed(() => {
  if (stat.value === DocReadingStatus.UNREAD) {
    return i18n.t('home.library.docs.readStatus.unread')
  }
  if (stat.value === DocReadingStatus.READING) {
    return `${i18n.t('home.library.docs.readStatus.reading')}${
      progress.value || 1
    }%`
  }
  return i18n.t('home.library.docs.readStatus.read')
})

let timer: ReturnType<typeof setTimeout>
const onSelecting = () => {
  clearTimeout(timer)
  selecting.value = true
}
const onUnselecting = () => {
  timer = setTimeout(() => (selecting.value = false), 300)
}

const onSelect = ({ selectedKeys }: any) => {
  if (!updating.value) {
    updateReadStatus({
      pdfId: props.pdfId,
      status: selectedKeys[0],
    })
  }
}

function useDocReadStatus(
  s: DocReadingStatus,
  p = 1,
  opts?: Partial<Options<object, unknown[]>>
) {
  const status = ref(s)
  const progress = ref(p || 1)
  const setStatus = (x: DocReadingStatus) => (status.value = x)
  const setProgress = (x: number) => (progress.value = x)

  const rest = useRequest<object, [UpdateReadStatusRequest]>(
    async (params: UpdateReadStatusRequest) => {
      setStatus(params.status ?? s)
      if (params.status === DocReadingStatus.READING) {
        setProgress(1)
      }
      const res = await updateDocReadStatus(params)

      return res
    },
    {
      ...opts,
      manual: true,
      defaultParams: [{ pdfId: '' }],
    }
  )

  return {
    ...rest,
    status,
    progress,
  }
}
</script>

<style lang="less">
.readstatus-wrapper {
  .ant-progress .ant-progress-text {
    font-size: 0;
    width: fit-content;
  }
  .ant-badge {
    font-size: 0;
  }
  .ant-badge-status-dot {
    width: 10px;
    height: 10px;
  }
  .ant-badge-status-text {
    display: none;
  }
}
.readstatus-dropdown {
  .ant-popover-arrow {
    left: 50%;
    transform: translateX(-50%) rotate(45deg);
    border-color: transparent #fff #fff transparent;
    box-shadow: 3px 3px 7px rgba(0, 0, 0, 0.07);
  }
  &.ant-dropdown-placement-bottomCenter {
    .ant-popover-arrow {
      top: -4px;
      border-color: #fff transparent transparent #fff;
      box-shadow: -2px -2px 5px rgba(0, 0, 0, 0.06);
    }
  }
  &.ant-dropdown-placement-topCenter {
    .ant-popover-arrow {
      bottom: -4px;
      border-color: transparent #fff #fff transparent;
      box-shadow: 3px 3px 7px rgba(0, 0, 0, 0.07);
    }
  }
  .ant-menu.ant-menu-vertical {
    border: none;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);

    .ant-menu-item {
      width: 32px;
      height: 26px;
      display: flex;
      justify-content: center;
      align-items: center;
      margin: 0 !important;
      padding: 2px 6px;

      &:hover {
        background-color: #f5f5f5;
      }
    }
  }
}
.readstatus-dot {
  // display: inline-flex;
  flex: 0 0 auto;
  width: 10px;
  height: 10px;
  border-radius: 50%;
}
</style>
