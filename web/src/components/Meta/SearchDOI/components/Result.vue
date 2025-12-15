<template>
  <div class="update-result">
    <UserInputForm
      :paper-id="paperId"
      :pdf-id="pdfId"
      :input-doc-info="result"
      :loading="false"
      :is-search="true"
      @back="handleBack"
      @submit="$emit('update:success')"
    />
  </div>
</template>
<script lang="ts">
import { computed, defineComponent } from 'vue'
import { AuthorInfo } from 'go-sea-proto/gen/ts/doc/CSL'
import { DocMetaInfoWithVenue } from '@common/api/citation'
import UserInputForm from '../../UserInput/index.vue'

export default defineComponent({
  components: {
    UserInputForm,
  },
  props: {
    result: {
      type: Object as () => DocMetaInfoWithVenue,
      default: () => null,
    },
    paperId: {
      type: String,
      default: '',
    },
    pdfId: {
      type: String,
      default: '',
    },
  },
  emits: ['update:success', 'back'],
  setup(props, { emit }) {
    const handleBack = () => {
      emit('back')
    }
    const values = computed<{ label: string; value?: string }[]>(() => {
      if (!props.result) {
        return []
      }
      const result = props.result

      const formateAuthorNames = (authorList: AuthorInfo[]) => {
        if (!authorList || !authorList.length) {
          return ''
        }
        const getFullname = (item: AuthorInfo) => {
          if (item.family || item.given) {
            return `${item.given || ''} ${item.family || ''}`
          }
          return item.literal
        }
        if (authorList.length <= 5) {
          return authorList
            .map((item) => {
              return getFullname(item)
            })
            .join(' / ')
        }
        const names = []
        for (let i = 0; i < 5; i++) {
          const item = authorList[i]
          names.push(getFullname(item))
        }
        names.push('...')
        names.push(getFullname(authorList[authorList.length - 1]))
        return names.join(' / ')
      }
      return [
        {
          label: '标题',
          value: result.title,
        },
        {
          label: '作者',
          value: formateAuthorNames(result.authorList),
        },
        {
          label: 'Doi',
          value: result.doi,
        },
        {
          label: '发布时间',
          value: result.publishDateStr,
        },
        {
          label: '收录情况',
          value: result.containerTitle?.[0],
        },
        {
          label: '分区信息',
          value: result.partition,
        },
        {
          label: '文献类型',
          value: result.docTypeName,
        },
        {
          label: '卷次',
          value: result.volume,
        },
        {
          label: '期号',
          value: result.issue,
        },
        {
          label: '页码',
          value: result.page,
        },
      ]
    })

    return {
      handleBack,
      values,
    }
  },
})
</script>
<style lang="less" scoped>
.update-result {
  display: flex;
  flex-direction: column;
  height: 100%;
  .total {
    color: #4e5969;
  }

  .list {
    height: 300px;
    overflow: auto;
    flex: 1;
    margin-bottom: 20px;
    li {
      display: flex;
      color: #1d2229;
      align-items: center;
    }
    .label {
      flex-basis: 76px;
    }

    .value {
      flex: 1;
      min-width: 0;
      display: -webkit-box;
      overflow: hidden;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }
    li + li {
      margin-top: 16px;
    }
  }
}

.btns {
  display: flex;
  justify-content: space-between;
}
.btn {
  background: #f0f2f5;
  border: none;
  color: #4e5969;
}
</style>
