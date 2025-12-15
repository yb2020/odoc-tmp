<template>
  <div
    v-show="noteItem"
    class="group-avatar"
    :style="{
      top: item.top + 'px',
      left: item.left + 'px',
    }"
    :data-userid="noteItem?.commentatorInfoView?.userId"
    @click="handleClick"
  >
    <a-badge
      :count="comments?.length ? `+${comments.length}` : 0"
      :number-style="{ backgroundColor: '#1F71E0', boxShadow: 'none' }"
    >
      <a-avatar
        :src="noteItem?.commentatorInfoView?.avatarCdnUrl"
        :size="32"
        class="avatar"
      >
        <template #icon>
          <UserOutlined />
        </template>
      </a-avatar>
    </a-badge>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { connect } from '~/src/dom/arrow';
import { useStore } from '~/src/store';
import { UserOutlined } from '@ant-design/icons-vue';
import { goPathPage } from '~/src/common/src/utils/url';
import { getDomainOrigin } from '~/src/util/env';
import { useAnnotationStore } from '~/src/stores/annotationStore';

const annotationStore = useAnnotationStore();

const props = defineProps<{
  item: {
    top: number;
    left: number;
    pdfAnnotateId: string;
    pageNumber: number;
  };
}>();

const store = useStore();

const noteItem = computed(() =>
  annotationStore.crossPageMap[props.item.pageNumber]?.find(
    (item) => item.uuid === props.item.pdfAnnotateId
  )
);

const comments = computed(
  () => store.state.note.comments[props.item.pdfAnnotateId]
);

const handleClick = () => {
  connect(props.item.pdfAnnotateId);

  const userId = noteItem.value?.commentatorInfoView?.userId;

  if (userId) {
    goPathPage(`${getDomainOrigin()}/user/${userId}`);
  }
};
</script>

<style lang="less" scoped>
.group-avatar {
  position: absolute;
  border-radius: 50%;
  cursor: pointer;

  .avatar {
    background: rgba(0, 0, 0, 0.3);
  }
}
</style>
