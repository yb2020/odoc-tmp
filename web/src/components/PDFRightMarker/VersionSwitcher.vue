<template>
  <div
    v-if="paperVersions"
    class="pdf-versions-switcher js-pdf-versions-switcher"
  >
    <a-dropdown placement="bottomRight">
      <a @click.prevent>
        {{
          $t('viewer.totalVersions', total === 1 ? 1 : 2, {
            named: { num: total },
          })
        }}
        <DownOutlined />
      </a>
      <template #overlay>
        <a-menu class="pdf-dropdown-overlay">
          <VersionList
            :paper-versions="paperVersions"
            :paper-id="paperId"
          />
        </a-menu>
      </template>
    </a-dropdown>
  </div>
</template>
<script lang="ts" setup>
import { GetPaperVersionsResponse } from 'go-sea-proto/gen/ts/paper/Paper'
import { computed, ref } from 'vue';
import { getPaperVersions } from '~/src/api/base';
import { DownOutlined } from '@ant-design/icons-vue';
import VersionList from './VersionList.vue';

const props = defineProps<{ paperId: string }>();

const paperVersions = ref<GetPaperVersionsResponse>();

const total = computed(() => {
  if (!paperVersions.value) {
    return 0;
  }
  return (
    (paperVersions.value.publicVersions?.length || 0) +
    (paperVersions.value.privateVersions?.length || 0)
  );
});

const renderVersions = async () => {
  const res = await getPaperVersions({ paperId: props.paperId });
  // res.privateVersions = res.publicVersions
  paperVersions.value = res;
};

renderVersions();
</script>
<style lang="less">
.pdf-versions-switcher {
  font-size: 13px;
}
</style>
