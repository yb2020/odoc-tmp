<template>
  <div class="personal-style-item-wrap">
    <span :class="['style-name', { 'exceed-style-name': index > 4 }]">{{
      data.customDefineTitle || data.shortTitle || data.title
    }}</span>
    <DeleteOutlined
      v-if="length > 1"
      class="delete-cion"
      @click="handleDeleteStyle"
    />
  </div>
</template>
<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { MyCslItem } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/CSL'
import { message } from 'ant-design-vue'
import { DeleteOutlined } from '@ant-design/icons-vue'
import { deleteCsl } from '@common/api/citation'
import { useI18n } from 'vue-i18n'

export default defineComponent({
  components: {
    DeleteOutlined,
  },
  props: {
    data: {
      type: Object as PropType<MyCslItem>,
      default: () => ({}),
    },
    index: {
      type: Number,
      default: -1,
    },
    length: {
      type: Number,
      default: 0,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n()

    const handleDeleteStyle = async () => {
      try {
        await deleteCsl({
          cslId: props.data.id,
        })

        message.success(t('home.paper.citeManage.tipRemoved') as string)

        emit('changeStyleSuccess', props.data.id)
      } catch (error) {}
    }
    return {
      handleDeleteStyle,
    }
  },
})
</script>
<style lang="less" scoped>
.personal-style-item-wrap {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-family: 'Lato';
  font-style: normal;
  font-weight: 400;
  font-size: 13px;
  line-height: 20px;
  color: #4e5969;
  padding: 4px 16px;
  cursor: pointer;
  .delete-cion {
    display: none;
    &:hover {
      color: #1d2229;
    }
  }
  &:hover {
    background: #e8f5ff;
    .delete-cion {
      display: block;
    }
  }
  .exceed-style-name {
    color: #a8afba;
  }
}
</style>
