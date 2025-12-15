<template>
  <div
    v-for="item in renderedPaperVersions"
    :key="item.from"
    class="pdf-versions-block"
  >
    <div class="title">
      {{ item.title }}
    </div>
    <ul class="versions">
      <li
        v-for="version in item.versions"
        :key="version.pdfId || version.jumpUrl"
      >
        <a
          v-if="version.type === PaperVersionType.PAPER_PDF"
          class="version-inner"
          @click="handleChangeVersion(version)"
        >
          <i
            aria-hidden="true"
            class="aiknowledge-icon icon-file-pdf"
          />
          <div class="name">{{ version.datePrefix }} · {{ version.name }}</div>
          <a-tag
            v-if="version.curVersion"
            class="tag"
            color="#52C41A"
          >{{
            $t('viewer.current')
          }}</a-tag>
          <a-tag
            v-else-if="version.lastVersion"
            class="tag last-tag no-rp-theme"
            color="#e5e6eb"
          >{{ $t('viewer.last') }}</a-tag>
        </a>
        <a
          v-else
          class="version-inner link"
          :href="version.jumpUrl"
          target="_blank"
        >
          <LinkOutlined />
          <div class="name">{{ version.datePrefix }} · {{ version.name }}</div>
        </a>
      </li>
    </ul>
    <div
      v-if="item.more > 0"
      class="more"
      @click="togglePublicVersions = !togglePublicVersions"
    >
      <span v-if="togglePublicVersions">
        {{ $t('viewer.collapseHistoryVersions', { num: item.more }) }}
        <UpOutlined />
      </span>
      <span v-else>{{ $t('viewer.expandHistoryVersions', { num: item.more }) }}&nbsp;
        <DownOutlined />
      </span>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { GetPaperVersionsResponse, PaperVersionInfo, PaperVersionType } from 'go-sea-proto/gen/ts/paper/Paper'
import { computed, ref } from 'vue';
import { LinkOutlined, DownOutlined, UpOutlined } from '@ant-design/icons-vue';
import { changeCurrentPDF } from '@/api/base';
import { message } from 'ant-design-vue';
import { isElectronMode, getPdfAnnotateNoteUrl } from '~/src/util/env';
import { ElementClick, PageType, reportElementClick } from '../../api/report';
import { useI18n } from 'vue-i18n';

const props =
  defineProps<{ paperVersions: GetPaperVersionsResponse; paperId: string }>();

const { t } = useI18n();

const togglePublicVersions = ref(false);

const renderedPaperVersions = computed<
  {
    title: string;
    from: 'public' | 'private';
    versions: PaperVersionInfo[];
    more: number;
  }[]
>(() => {
  const res = props.paperVersions;

  if (!res) {
    return [];
  }
  const renderedVersions: {
    title: string;
    from: 'public' | 'private';
    versions: PaperVersionInfo[];
    more: number;
  }[] = [];
  if (res.publicVersions?.length) {
    renderedVersions.push({
      title: t(
        'viewer.publicVersions',
        res.publicVersions.length === 1 ? 1 : 2,
        { named: { num: res.publicVersions.length } }
      ),
      versions: togglePublicVersions.value
        ? [...res.publicVersions]
        : res.publicVersions.slice(0, 3),
      from: 'public',
      more: res.publicVersions.length - 3,
    });
  }
  if (res.privateVersions?.length) {
    renderedVersions.push({
      title: t('viewer.privateVersions'),
      versions: res.privateVersions,
      from: 'private',
      more: 0,
    });
  }
  return renderedVersions;
});

const handleChangeVersion = async (version: PaperVersionInfo) => {
  if (version.curVersion) {
    return;
  }
  await changeCurrentPDF({ pdfId: version.pdfId, paperId: props.paperId });
  message.success(t('message.versionChangedNoteDismatchedTip'));
  setTimeout(() => {
    if (isElectronMode()) {
      window.location.reload();
    } else {
      const noteId =
        new URL(window.location.href).searchParams.get('noteId') || '';
      const url = getPdfAnnotateNoteUrl({ noteId });
      window.location.href = url.pathname + url.search;
    }
  }, 1000);

  reportElementClick({
    page_type: PageType.note,
    type_parameter: props.paperId,
    element_name: ElementClick.version_switch,
  });
};
</script>
<style lang="less">
.pdf-versions-block {
  padding: 8px 16px;
  width: 248px;

  .title {
    color: #86919c;
    font-size: 12px;
  }

  .versions {
    font-size: 14px;

    list-style: none;
    padding: 0;

    .version-inner {
      display: flex;
      align-items: center;
      cursor: pointer;

      .tag {
        line-height: 16px;
        padding: 0 2px;
        margin-right: 0;

        &.last-tag.no-rp-theme {
          color: #4e5969;
        }
      }

      .name {
        margin-left: 4px;
        max-width: 196px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      &.link {
        color: #4e5969;
      }

      & + & {
        margin-top: 10px;
      }
    }
  }

  & + & {
    border-top: 1px solid #e5e6eb;
  }

  .more {
    color: #86919c;
    text-align: center;
    cursor: pointer;
    margin: 8px 0;
    font-size: 12px;
  }
}
</style>
