<template>
  <div class="metadata-card-container">
    <div
      class="metadata-card-title"
      @click="openPaperPage()"
    >
      {{ paperData.title }}
    </div>
    <div class="meta-card-author-list" />
    <div
      ref="lightRef"
      class="metadata-card-abstract-light"
      :style="{
        display: expand ? 'block' : '-webkit-box',
      }"
    >
      {{ paperData.originalAbstract }}
      <div
        v-show="overflow || expand"
        class="metadata-card-expand"
        @click="expand = !expand"
      >
        <UpOutlined v-if="expand" />
        <DownOutlined v-else />
      </div>
    </div>
    <div
      ref="shadowRef"
      class="metadata-card-abstract-shadow"
    >
      {{ paperData.originalAbstract }}
    </div>
    <div class="metadata-card-footer">
      <div class="meta-card-cite">
        被引用
      </div>
      <div class="meta-card-save">
        收藏
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { PaperDetailInfo } from 'go-sea-proto/gen/ts/paper/Paper'
import { UpOutlined, DownOutlined } from '@ant-design/icons-vue';
import { getDomainOrigin } from '~/src/util/env';
import { goPathPage } from '~/src/common/src/utils/url';
import { useTextOverflow } from './helper';

const { paperData } = defineProps<{ paperData: PaperDetailInfo }>();

const { lightRef, shadowRef, overflow, expand } = useTextOverflow();

const openPaperPage = () => {
  goPathPage(`${getDomainOrigin()}/paper/${paperData.paperId}`);
};
</script>

<style lang="less">
.metadata-card-container {
  position: relative;
  background-color: #fff;
  width: 422px;
  border-radius: 4px;
  padding-top: 20px;
  .metadata-card-title {
    margin-left: 24px;
    margin-right: 24px;
    margin-bottom: 5px;
    font-family: 'Roboto';
    font-style: normal;
    font-weight: bold;
    font-size: 14px;
    line-height: 22px;
    color: #1f71e0;
  }
  .metadata-card-abstract-light,
  .metadata-card-abstract-shadow {
    font-family: 'Roboto';
    font-style: normal;
    font-weight: 400;
    font-size: 13px;
    line-height: 20px;
    /* or 154% */
    color: #262625;
    margin-left: 24px;
    width: 377px;
    margin-bottom: 16px;
  }

  .metadata-card-abstract-light {
    position: relative;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 5;
    overflow: hidden;
  }

  .metadata-card-abstract-shadow {
    position: fixed;
    display: block;
    opacity: 0;
    pointer-events: none;
  }

  .metadata-card-expand {
    position: absolute;
    bottom: 0;
    right: 0;
    width: 40px;
    height: 20px;
    display: flex;
    align-items: center;
    background: linear-gradient(
      to right,
      rgba(255, 255, 255, 0) -14%,
      rgba(255, 255, 255, 1) 58%
    );
    cursor: pointer;
    display: flex;
    justify-content: flex-end;
  }

  .metadata-card-footer {
    height: 48px;
    border-top: 1px solid #e4e7ed;
  }
}
</style>
