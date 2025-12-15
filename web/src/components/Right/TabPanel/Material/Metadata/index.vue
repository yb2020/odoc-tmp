<template>
  <div class="metadata-container">
    <TitleVue
      :paper-data="paperData"
      :doc-info="docInfo"
      :doc-id="docIdForEditing"
    />
    <AuthorVue
      :with-author-link="withAuthorLink"
      :display-authors="docInfo?.displayAuthor"
      :doc-id="docIdForEditing"
    />

    <Cite
      v-if="docInfo" 
      :paper-id="docInfo.paperId" 
      :pdf-id="docInfo.pdfId"
      :page-type="PageType.library"
      @update:success="onUpdateSuccess"
    />

    <PublishDateVue
      :doc-id="docIdForEditing"
      :display-publish-date="docInfo?.displayPublishDate"
    />
    <div class="metadata-field">
      <div class="metadata-field-name">
        {{ $t('viewer.collectionLabel') }}
      </div>
      <VenueVue
        :doc-id="docIdForEditing"
        :display-venue="docInfo?.displayVenue"
      />
    </div>
    <div
      v-if="withPartition"
      class="metadata-field"
    >
      <div class="metadata-field-name">
        分区信息
      </div>
      <div
        class="metadata-field-value"
        style="pointer-events: none"
      >
        {{ docInfo?.partition ?? '-' }}
      </div>
    </div>
    <div
      v-if="isOpenPaper"
      class="metadata-field"
    >
      <div class="metadata-field-name">
        {{ $t('viewer.refsAndCitesLabel') }}
      </div>
      <RefCiteVue
        :paper-id="selfNoteInfo?.paperId ?? ''"
        :pdf-id="selfNoteInfo?.pdfId ?? ''"
      />
    </div>
    <div
      v-if="isCollected && docIdForEditing"
      class="metadata-field"
    >
      <div class="metadata-field-name">
        {{ $t('viewer.folderLabel') }}
      </div>
      <FolderVue
        :note-id="selfNoteInfo?.noteId"
        :doc-id="docIdForEditing"
      />
    </div>
    <div
      v-if="isCollected && docIdForEditing"
      class="metadata-field"
      style="height: auto"
    >
      <div class="metadata-field-name">
        {{ $t('viewer.label') }}
      </div>
      <ClassifyVue
        :doc-id="docIdForEditing"
        :classify-infos="docInfo?.classifyInfos"
      />
    </div>
    <div
      v-if="isCollected && docIdForEditing"
      class="metadata-field"
    >
      <div class="metadata-field-name">
        {{ $t('viewer.remarkLabel') }}
      </div>
      <div
        v-if="!docInfo?.remark && !remarkRef?.editFlag"
        class="metadata-field-value"
        style="cursor: pointer"
        @click="remarkRef?.startEdit()"
      />
    </div>
    <RemarkVue
      v-if="isCollected && docIdForEditing"
      ref="remarkRef"
      :doc-info="docInfo!"
    />
    <StatusVue 
      :doc-id="docInfo?.docId"
      :pdf-id="docInfo?.pdfId"
    />
    <div v-if="fetchPaperData.fetchState.error || error">
      {{ $t('message.failedToLoadDataTip') }}
    </div>
  </div>
</template>
<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { PaperDetailInfo } from 'go-sea-proto/gen/ts/paper/Paper'
import useFetch from '~/src/hooks/useFetch';
import { getPaperDetailInfo } from '~/src/api/reference';
import RefCiteVue from '../RefCite/index.vue';
import TitleVue from './Title.vue';
import AuthorVue from './Author.vue';
import FolderVue from './Folder.vue';
import PublishDateVue from './PublishDate.vue';
import VenueVue from './Venue.vue';
import RemarkVue from './Remark.vue';
import ClassifyVue from './Classify.vue';
import { DocDetailInfo } from 'go-sea-proto/gen/ts/doc/ClientDoc';
import { checkOpenPaper } from '~/src/api/helper';
import { selfNoteInfo, isOwner } from '~/src/store';
import Cite from './Cite.vue';
import StatusVue from './Status.vue';
import { PageType } from '~/src/api/report';

const props = defineProps<{
  withPartition: boolean;
  withAuthorLink: boolean;
  docInfo: null | DocDetailInfo;
  docInfoReload: () => void | Promise<any>;
  error?: null | Error;
}>();

const paperData = ref<PaperDetailInfo | null>(null);

const docIdForEditing = computed(() => {
  if (!isOwner.value) {
    return undefined;
  }

  return props.docInfo?.docId;
});
const isCollected = computed(() => selfNoteInfo.value?.isCollected ?? false)

const remarkRef = ref();
const isOpenPaper = computed(() =>
  checkOpenPaper(
    selfNoteInfo.value?.paperId ?? '',
    selfNoteInfo.value?.isPrivatePaper ?? false
  )
);

const fetchPaperData = useFetch(async () => {
  if (!isOpenPaper.value) {
    return;
  }

  const createError = () => {
    paperData.value = null;
    return new Error('Invalid paperId');
  };

  try {
    var response = await getPaperDetailInfo({
      paperId: selfNoteInfo.value?.paperId ?? '',
    });
  } catch (error) {
    throw createError();
  }

  if (!response.title) {
    throw createError();
  }

  paperData.value = response;
}, false);

watch(() => selfNoteInfo.value?.paperId ?? '', fetchPaperData.fetch, {
  immediate: true,
});

const onUpdateSuccess = async () => {
  fetchPaperData.fetch();
  await props.docInfoReload()
  if (props.docInfo) {
    selfNoteInfo.value.docName = props.docInfo.docName;
  }
};
</script>

<style lang="less">
@import './style.less';
.metadata-field {
  height: 32px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  color: var(--site-theme-text-primary);
  input {
    flex: 1 1 100%;
    height: 32px;
    border: 1px solid var(--site-theme-brand);
    border-radius: 2px;
    outline: 0;
    color: var(--site-theme-text-primary);
    text-indent: 6px;
  }
  .metadata-field-name {
    flex: 0 0 56px;
    padding-right: 4px;
    color: var(--site-theme-pdf-panel-text);
    opacity: 0.45;
  }
}
</style>
