<template>
  <Notes collapsable>
    <template #icon="{ visible, toggleVisible, ...rest }">
      <Toggle
        v-bind="rest"
        :target="!visible ? target : undefined"
        :visible="visible"
        :toggle-visible="toggleVisible"
      />
    </template>
  </Notes>
</template>
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import Notes from '@common/components/Notes/index.vue'
import Toggle from '@/components/Layout/Toggle.vue'
import { PAGE_ROUTE_NAME } from '@/routes/type'
import { PageType, useReportVisitDuration } from '@common/utils/report'
import { NOTES_CLASSNAME } from '@common/components/Notes/useNote'

const target = ref<HTMLElement>()

onMounted(() => {
  target.value =
    document.querySelector<HTMLElement>(
      `.js-sidebar li[name="${PAGE_ROUTE_NAME.NOTES}"]`
    ) || undefined
})

useReportVisitDuration(
  () => '',
  () => ({
    page_type: PageType.NOTE_TAB,
    type_parameter: 'none',
  }),
  () => {
    const container = document.getElementsByClassName(NOTES_CLASSNAME)[0]
    return !container || !container.parentElement
  }
)
</script>
