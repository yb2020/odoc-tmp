<template>
  <div class="other-style-item-wrap">
    <span class="title">{{ data.title }}</span>
    <div
      :class="['icon-wrap', { 'added-icon-wrap': data.addedFlag }]"
      @click="handleClick"
    >
      <i class="aiknowledge-icon icon-add" aria-hidden="true"></i>
    </div>
  </div>
</template>
<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { CslItem } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/CSL'
import { message } from 'ant-design-vue'
import { addCsl, deleteCsl } from '@common/api/citation'
import { useI18n } from 'vue-i18n'

export default defineComponent({
  props: {
    data: {
      type: Object as PropType<CslItem>,
      default: () => ({}),
    },
    addedCount: {
      type: Number,
      default: 0,
    },
  },
  setup(props, { emit }) {
    const i18n = useI18n()

    const handleClick = async () => {
      if (props.data.addedFlag && props.addedCount === 1) {
        message.warning(i18n.t('home.paper.citeManage.tipKeepLimit') as string)

        return
      }

      const request = props.data.addedFlag ? deleteCsl : addCsl

      try {
        await request({
          cslId: props.data.id,
        })

        message.success(
          props.data.addedFlag
            ? (i18n.t('home.paper.citeManage.tipRemoved') as string)
            : (i18n.t('home.paper.citeManage.tipAdded') as string)
        )

        emit('changeStyleSuccess', props.data.id)
      } catch (error) {}
    }
    return {
      handleClick,
    }
  },
})
</script>
<style lang="less" scoped>
.other-style-item-wrap {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-family: 'Lato';
  font-style: normal;
  font-weight: 400;
  font-size: 12px;
  line-height: 18px;
  color: #4e5969;
  padding: 4px 20px 4px 16px;
  &:hover {
    background: #f0f2f5;
  }
  .title {
    word-break: break-word;
    margin-right: 16px;
  }
  .icon-wrap {
    padding: 6px;
    cursor: pointer;
    line-height: 12px;
    background: #e5e6eb;
    border-radius: 2px;
    display: flex;
    .icon-add {
      font-size: 12px;
      color: #1f71e0;
      height: 12px;
      line-height: 12px;
    }
  }
  .added-icon-wrap {
    background: #f0f2f5;
    .icon-add {
      color: #a8afba;
    }
  }
}
</style>
