<template>
  <Align
    v-if="target"
    ref="aligned"
    visible
    :to="body"
    :target="target"
    :alignProps="{
      points: ['cl', 'cr'],
      // offset: [0, -10],
      // overflow: { adjustX: true, adjustY: true },
    }"
  >
    <template #align>
      <ToggleBtn
        class="relative"
        :class="isOutsideTarget && isOutsideBtn ? 'opacity-0' : ''"
        :visible="visible"
        @click="toggleVisible()"
      />
    </template>
  </Align>
  <ToggleBtn
    v-else
    class="absolute top-40 mt-2"
    :class="[visible ? 'left-full' : 'right-full']"
    :visible="visible"
    @click="toggleVisible()"
  />
</template>

<script setup lang="ts">
import Align from '@common/components/Align/index.vue'
import { useMouseInElement } from '@vueuse/core'
import { computed, ref } from 'vue'
import ToggleBtn from './ToggleBtn.vue'

const props = defineProps<{
  target?: null | HTMLElement
  visible: boolean
  toggleVisible: (v?: boolean) => void
}>()

const { body } = document
const aligned = ref()
const target = computed(() => props.target)
const button = computed(() => aligned.value?.wrapperEl as HTMLElement)

const { isOutside: isOutsideTarget } = useMouseInElement(target)
const { isOutside: isOutsideBtn } = useMouseInElement(button)
</script>
