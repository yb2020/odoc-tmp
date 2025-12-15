<template>
  <MarkerVue 
    v-for="(marker, i) in markers" 
    :key="i"
    :viewport-width="viewportWidth"
    :marker="marker"
  />
</template>
<script setup lang="ts">
import { FigureAndTableReferenceMarker } from "@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse";
import { ReferenceMarker } from 'go-sea-proto/gen/ts/pdf/PdfParse'; 
import MarkerVue from "./Marker.vue";

defineProps<{ 
    markers: ReferenceMarker[] | FigureAndTableReferenceMarker[],
    viewportWidth: number,
    uuid: 'reference' | 'figureAndTable' 
  }>();

</script>
