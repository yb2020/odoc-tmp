<template>
  <div class="metadata-author-container">
    <div
      ref="triggerRef"
      class="metadata-author-light"
      :style="{
        visibility: ready ? 'visible' : 'hidden',
        cursor: props.inited ? 'pointer' : 'initial',
      }"
    >
      {{ authorsView }}
    </div>
    <div ref="shadowRef" class="metadata-author-shadow">
      {{ authorsView }}
    </div>
    <TippyVue
      v-if="triggerRef && props.inited"
      ref="tippyRef"
      :trigger-ele="(triggerRef as any)"
      placement="right-start"
      :offset="[0, 20]"
      trigger="click"
      :z-index="9999"
      :appendTo="$getContainer"
      @onShown="docId = props.inited"
      @onHide="docId = false"
    >
      <AuthorEditVue
        :author-list-backup="docInfo.authorList"
        :with-author-link="withAuthorLink"
        @cancel="cancel()"
        @update="update($event)"
      />
    </TippyVue>
  </div>
</template>

<script setup lang="ts">
import { AuthorInfo, DocMetaInfoSimpleVo } from 'go-sea-proto/gen/ts/doc/CSL'
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import TippyVue from '@common/components/Tippy/index.vue'
import AuthorEditVue from './AuthorEdit.vue'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  inited: boolean
  docInfo: DocMetaInfoSimpleVo
  withAuthorLink?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:authorList', authorList: AuthorInfo[]): void
}>()

const safeAuthors = computed(() => {
  return props.docInfo.authorList ?? []
})
const initMiddleNumber = () => {
  return Math.max(0, safeAuthors.value.length - 2)
}
const middleNumber = ref(initMiddleNumber())
const ready = ref(false)
const authors = computed(() => {
  const [head, ...rest] = safeAuthors.value
  const tail = rest.pop()
  const after: string[] = []
  if (head) {
    after.push(head.literal)
  }

  after.push(...rest.slice(0, middleNumber.value).map((item) => item.literal))

  if (middleNumber.value < rest.length) {
    after.push('…')
  }

  if (tail) {
    after.push(tail.literal)
  }

  return after
})

console.warn({ authors, middleNumber })

const { t } = useI18n()
const authorsView = computed(() => {
  return authors.value.join(' / ') || t('meta.detail.authorPlaceholder')
})

const sync = async () => {
  await nextTick()
  if (
    shadowRef.value!.offsetHeight > triggerRef.value!.offsetHeight &&
    middleNumber.value > 0
  ) {
    middleNumber.value -= 1
  } else {
    ready.value = true
  }
}

watch(safeAuthors, () => {
  ready.value = false
  middleNumber.value = initMiddleNumber()
})
watch(authors, sync)
onMounted(sync)

const triggerRef = ref<HTMLDivElement | null>(null)
const shadowRef = ref<HTMLDivElement | null>(null)
const tippyRef = ref()
const docId = ref(false)
const cancel = () => {
  tippyRef.value.hide()
  docId.value = false
}

const update = (authorList: AuthorInfo[]) => {
  emit('update:authorList', authorList)
  tippyRef.value.hide()
}
</script>

<style lang="less" scoped>
.metadata-author-container {
  flex: 1 1 100%;
  position: relative;
  padding-left: 8px;
  padding-right: 8px;
  padding-top: 4px;
  padding-bottom: 4px;
  border-radius: 2px;

  border-width: 1px;
  border-style: solid;
  border-color: #c9cdd4;

  .metadata-author-light,
  .metadata-author-shadow {
    font-size: 14px;
    word-break: break-all;
  }

  .metadata-author-light {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .metadata-author-shadow {
    pointer-events: none;
    opacity: 0;
    position: absolute;
  }
}
</style>
<style lang="less" scoped>
.metadata-author-container {
  &:hover {
    background: rgba(255, 255, 255, 0.08);

    .metadata-rollback {
      top: 50% !important;
      transform: translateY(-50%);
      visibility: visible !important;
    }
  }
}
</style>
