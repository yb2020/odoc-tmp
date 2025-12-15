<template>
  <Modal
    :title="$t('home.library.shareFolder')"
    :visible="visible"
    :footer="null"
    destroy-on-close
    :width="384"
    @cancel="cancel"
  >
    <p class="text">
      {{ $t('home.library.copyLinkToShare1') }}
      <span style="color: #1f71e0"> “{{ folderTitle }}” </span>
      {{ $t('home.library.copyLinkToShare2') }}
    </p>
    <p class="url">{{ url }}</p>
    <div class="buttons">
      <Button @click="preview">
        <EyeOutlined class="iconeye" />
        {{ $t('home.library.preview') }}
      </Button>
      <Button type="primary" @click="copy">
        <CopyOutlined />
        {{ copied ? $t('home.library.copied') : $t('home.library.copyLink') }}
      </Button>
    </div>
  </Modal>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import { message, Modal, Button } from 'ant-design-vue'
import { EyeOutlined, CopyOutlined } from '@ant-design/icons-vue'
import { goPathPage } from '@common/utils/url'
import { useI18n } from 'vue-i18n'

const $t = useI18n().t

const props = defineProps<{
  visible: boolean
  folderKey: string
  folderTitle: string
}>()

const emit = defineEmits<{
  (event: 'cancel'): void
}>()

const url = computed(() => {
  return `${
    typeof window !== 'undefined' ? window.location.origin : ''
  }/user/collect/${props.folderKey}`
})
const copied = ref(false)
const copy = async () => {
  if (window.navigator.clipboard) {
    await window.navigator.clipboard.writeText(url.value)
  } else {
    const input = document.createElement('input')
    document.body.appendChild(input)
    input.value = url.value
    input.select()
    document.execCommand('copy')
    input.remove()
  }

  copied.value = true
  message.success($t('home.library.copySucc') as string)
}
const preview = () => {
  goPathPage(`/user/collect/${props.folderKey}`)
}
const cancel = () => {
  emit('cancel')
  copied.value = false
}
</script>

<style lang="less" scoped>
.text {
  font-size: 14px;
  color: #262625;
  line-height: 24px;
  margin-bottom: 8px;
}

.url {
  color: #73716f;
  line-height: 24px;
  margin-bottom: 40px;
  padding: 8px 12px;
  background: #f2f5f7;
}

.buttons {
  display: flex;
  justify-content: center;
  .ant-btn {
    width: 124px;
    .iconeye {
      font-size: 13px;
    }
  }
  .ant-btn + .ant-btn {
    margin-left: 16px;
  }
}
</style>
