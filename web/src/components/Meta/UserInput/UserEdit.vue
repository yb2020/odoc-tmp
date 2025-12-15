<template>
  <UserInputForm
    :paper-id="paperId"
    :pdf-id="pdfId"
    :input-doc-info="docInfo"
    :loading="loading"
    :is-search="false"
    @cancel="$emit('cancel')"
    @submit="$emit('submit')"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { DocMetaInfoWithVenue, getDocMetaInfo } from '@common/api/citation'
import UserInputForm from './index.vue'

const props = defineProps<{
  paperId: string
  pdfId: string
}>()

defineEmits<{
  (event: 'submit'): void
  (event: 'cancel'): void
}>()

const docInfo = ref<DocMetaInfoWithVenue | null>(null)
const loading = ref(false)

const fetchDocInfo = async () => {
  loading.value = true

  const responseMeta = getDocMetaInfo({
    paperId: props.paperId,
    pdfId: props.pdfId,
  })

  try {
    docInfo.value = await responseMeta
  } catch (error) {
    loading.value = false
  }

  loading.value = false
}

fetchDocInfo()
</script>
