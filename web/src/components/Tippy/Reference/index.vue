<template>
  <div :class="['reference-box', { loading: !fetchState.error && !paperData }]">
    <ErrorVue
      v-if="fetchState.error"
      :title="paperTitle"
    />
    <InfoVue
      v-else-if="paperData"
      :showCiteBtn="showCiteBtn"
      :paper-data="paperData"
      :no-collect="noCollect"
      :marker="marker"
      :paper-data-reload="fetch"
      @hide="onHide"
    />
    <a-spin v-else />
  </div>
</template>
<script setup lang="ts">
import { PaperDetailInfo } from 'go-sea-proto/gen/ts/paper/Paper'
import { nextTick, ref, watch } from "vue";
import { getPaperDetailInfo } from "~/src/api/reference";
import useFetch from "~/src/hooks/useFetch";
import ErrorVue from './Error.vue'
import InfoVue from "./Info.vue";
import {
  ReferenceMarker,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse';

const props = defineProps<{ 
  showCiteBtn?: boolean;
  paperId: string;
  paperTitle: string; 
  fetchFlag: boolean;
  noCollect?: boolean;
  marker?: ReferenceMarker;
  updateContentHandler?: () => void;
  hideTippy?: () => void
}>();

const emit = defineEmits<{
  (event: 'updateContent'): void
  (event: 'hide'): void
}>()

const paperData = ref<PaperDetailInfo | null>(null)

const { fetch, fetchState } = useFetch(async () => {
  if (!props.paperId) {
    emit('updateContent')
    throw new Error('Invalid paperId');
  }
  try {
    const res = await getPaperDetailInfo({ paperId: props.paperId });
    if (!res.title) {
      emit('updateContent')
      props.updateContentHandler?.();
      throw new Error('Invalid paperId');
    }
    paperData.value = res;
    nextTick(() => {
      props.updateContentHandler?.();
      emit('updateContent')
    })
  } catch (err) {
    props.updateContentHandler?.();
    emit('updateContent')
    throw new Error('Invalid paperId');
  }
}, false);

watch(() => props.fetchFlag, (newVal) => {
  if (newVal && !paperData.value) {
    fetch()
  }
}, {
  immediate: true,
})

const onHide = () => {
  props.hideTippy?.()
  emit('hide')
}

</script>


<style scoped lang="less">
.reference-box {
  width: 473px;
  background-color: #fff;
  &.loading {
    text-align: center;
    line-height: 120px;
    min-height: 120px;
  }
}
</style>
