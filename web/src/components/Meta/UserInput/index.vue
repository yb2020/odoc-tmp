<template>
  <div class="user-input-container">
    <div class="user-input-field">
      <label>{{ $t('meta.detail.titleLabel') }}</label>
      <Title :inited="docInfo !== defaultDocInfo" :docInfo="docInfo" />
    </div>
    <div class="user-input-field">
      <label>{{ $t('meta.detail.authorLabel') }}</label>
      <Author
        :inited="docInfo !== defaultDocInfo"
        :docInfo="docInfo"
        :withAuthorLink="true"
        @update:authorList="handleAuthorListUpdate"
      />
    </div>
    <div class="user-input-field">
      <label>Doi</label>
      <input
        v-model="docInfo.doi"
        :placeholder="$t('meta.detail.doiPlaceholder')"
      >
    </div>
    <div class="user-input-field">
      <label>{{ $t('meta.detail.dateLabel') }}</label>
      <input
        v-model.trim="docInfo.publishDateStr"
        :placeholder="$t('meta.detail.datePlaceholder')"
      >
    </div>
    <div class="user-input-field">
      <label>{{ $t('meta.detail.venueLabel') }}</label>
      <input
        v-model.trim="docInfo.venue"
        :placeholder="$t('meta.detail.venuePlaceholder')"
      >
    </div>
    <div class="user-input-field">
      <label>{{ $t('meta.detail.partitionLabel') }}</label>
      <input 
        v-model.trim="docInfo.partition" 
        :placeholder="$t('meta.detail.partitionPlaceholder')"
      >
    </div>
    <div class="user-input-field">
      <label>{{ $t('meta.detail.typeLabel') }}</label>
      <a-select
        v-model:value="docInfo.docType"
        :allowClear="true"
        class="user-input-papertype"
        :placeholder="$t('meta.detail.typePlaceholder')"
      >
        <a-select-option v-for="item in paperTypeList" :key="item.code">
          {{ isZhCN ? item.name : item.code }}
        </a-select-option>
      </a-select>
    </div>
    <div class="user-input-field">
      <label>{{ $t('meta.detail.volumeLabel') }}</label>
      <input
        v-model.trim="docInfo.volume"
        :placeholder="$t('meta.detail.volumePlaceholder')"
      >
    </div>
    <div class="user-input-field">
      <label>{{ $t('meta.detail.issueLabel') }}</label>
      <input
        v-model.trim="docInfo.issue"
        :placeholder="$t('meta.detail.issuePlaceholder')"
      >
    </div>
    <div class="user-input-field">
      <label>{{ $t('meta.detail.pageLabel') }}</label>
      <input
        v-model.trim="docInfo.page"
        :placeholder="$t('meta.detail.pagePlaceholder')"
      >
    </div>
    <div
      class="user-input-button-group"
      :style="{ justifyContent: isSearch ? 'space-between' : 'flex-end' }"
    >
      <button v-if="isSearch" class="user-input-cancel" @click="$emit('back')">
        {{ $t('meta.back') }}
      </button>
      <button v-else class="user-input-cancel" @click="$emit('cancel')">
        {{ $t('meta.cancel') }}
      </button>
      <button
        @click="submitDocInfo"
        :class="{
          'user-input-submitting': submitting,
        }"
      >
        {{ isSearch ? $t('meta.confirm') : $t('meta.save') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from 'vue-i18n'

import {
  createDefaultDocInfo,
  postManualUpdateDocCiteInfo,
  DocMetaInfoWithVenue,
  getDocTypeList,
} from '@/api/citation'
import Title from './Title.vue'
import Author from './Author.vue'
import { delay } from '@idea/aiknowledge-special-util'
import { nextTick } from 'vue'
import { DocTypeInfo } from 'go-sea-proto/gen/ts/doc/CSL'
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '@/hooks/useLanguage';

const props = defineProps<{
  paperId: string
  pdfId: string
  inputDocInfo: DocMetaInfoWithVenue | null
  loading: boolean
  isSearch: boolean
}>()

const emit = defineEmits<{
  (event: 'submit'): void
  (event: 'cancel'): void
  (event: 'back'): void
}>()

// 语言管理
const { isCurrentLanguage } = useLanguage();
const { t } = useI18n();

const defaultDocInfo = createDefaultDocInfo()
const docInfo = computed(() => props.inputDocInfo || defaultDocInfo)
const paperTypeList = ref<DocTypeInfo[]>([])
const submitting = ref(false)

// 基于 proto 枚举的语言判断计算属性
const isZhCN = computed(() => isCurrentLanguage(Language.ZH_CN))

;(async () => {
  paperTypeList.value = await getDocTypeList({})
})()

// 处理作者列表更新
// 注意：这里直接修改props，与原始代码保持一致
// 在正常Vue实践中不推荐这样做，但这里为了保持与原始代码的兼容性
// eslint-disable-next-line vue/no-mutating-props
const handleAuthorListUpdate = (authorList) => {
  if (props.inputDocInfo) {
    // 直接修改props，与原始代码保持一致
    // eslint-disable-next-line vue/no-mutating-props
    props.inputDocInfo.authorList = authorList
  }
}

const submitDocInfo = async () => {
  await nextTick()
  await delay(300)

  if (submitting.value) {
    return
  }

  submitting.value = true

  const {
    title,
    authorList,
    doi,
    publishDateStr,
    venue,
    partition,
    docType,
    volume,
    issue,
    page,
  } = docInfo.value

  try {
    await postManualUpdateDocCiteInfo({
      paperId: props.paperId,
      pdfId: props.pdfId,
      docName: title,
      authorList,
      doi,
      publishDate: publishDateStr,
      venue,
      partition,
      docType,
      volume,
      issue,
      page,
    })
  } catch (error) {
    submitting.value = false
    return
  }

  submitting.value = false
  message.success(t('meta.updateSuccess'))
  emit('submit')
}
</script>

<style lang="less" scoped>
.user-input-container {
  display: flex;
  flex-direction: column;

  .user-input-field {
    display: flex;
    justify-content: space-between;
    color: #1d2229;
    font-size: 14px;
    margin-bottom: 8px;

    label {
      flex: 0 0 76px;
      display: flex;
      align-items: center;
      justify-content: flex-start;
    }

    > input {
      flex: 1 1 100%;
      height: 32px;
      border-radius: 2px;
      border-width: 1px;
      border-style: solid;
      border-color: #c9cdd4;
      background: #fff;
      outline: 0;
      text-indent: 12px;

      &::placeholder {
        color: #a8afba;
      }

      &:focus {
        border-color: #1f71e0;
      }
    }

    .user-input-papertype {
      flex: 1 1 100%;
      
      :deep(.ant-select-selector) {
        color: #1d2229;
      }
      
      :deep(.ant-select-dropdown) {
        .ant-select-item {
          color: #1d2229;
        }
        .ant-select-item-option-selected {
          color: #1d2229;
          font-weight: bold;
        }
      }
    }
  }

  .user-input-button-group {
    margin-top: 8px;
    display: flex;

    button {
      border: 0;
      outline: 0;
      border-radius: 2px;
      width: 88px;
      height: 32px;
      display: flex;
      justify-content: center;
      align-items: center;
      font-size: 14px;

      &.user-input-cancel {
        margin-right: 16px;
        background-color: #f0f2f5;
        color: #4e5969;
      }

      &:last-of-type {
        color: white;
        background-color: #1f71e0;

        &.user-input-submitting {
          opacity: 0.7;
          cursor: wait;
        }
      }
    }
  }
}
</style>
<style lang="less">
.user-input-container {
  .user-input-field {
    .user-input-papertype {
      > div {
        border: 1px solid #c9cdd4 !important;
      }
    }
  }
}
</style>
