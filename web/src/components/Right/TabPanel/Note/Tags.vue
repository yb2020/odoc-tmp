<template>
  <IdeaTag
    v-if="props.note.uuid !== NO_ANNOTATION_ID"
    class="note-tags"
    :list="annotationStore.tagList"
    :selected="props.note.tags ?? []"
    :tag-label="$t('viewer.tag')"
    @selecting="emit('selecting')"
    @selectend="emit('selectend')"
    @input-tag="inputTag($event)"
    @select-tag="selectTag($event)"
    @remove="removeTag($event)"
  />
</template>

<script setup lang="ts">
import { NO_ANNOTATION_ID } from '@/constants';
import { IdeaTag } from '@idea/aiknowledge-markdown';
//import { AnnotateTag } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { AnnotateTag } from 'go-sea-proto/gen/ts/common/AnnotateTag';
import {
  addTagToAnnotation,
  createAnnotationTag,
  deleteTagFromAnnotation,
} from '@/api/annotations';
import {
  useAnnotationStore,
  personAnnotationController,
} from '@/stores/annotationStore';
import { AnnotationAll } from '@/stores/annotationStore/BaseAnnotationController';

const props = defineProps<{
  note: AnnotationAll;
}>();

const emit = defineEmits<{
  (e: 'selecting'): void;
  (e: 'selectend'): void;
  (e: 'added'): void;
  (e: 'removed'): void;
}>();

const annotationStore = useAnnotationStore();

const localAddTag = (tag: AnnotateTag) => {
  const { pageNumber, uuid } = props.note;
  personAnnotationController.addTag(uuid, pageNumber, tag);

  emit('added');
};
const localRemoveTag = (tagId: AnnotateTag['tagId']) => {
  const { pageNumber, uuid } = props.note;
  personAnnotationController.removeTag(uuid, pageNumber, tagId);

  emit('removed');
};

const inputTag = async (tagName: AnnotateTag['tagName']) => {
  try {
    // 1. 创建新标签
    const response = await createAnnotationTag({
      tagName,
      markId: props.note.uuid,
    });
    
    // 2. 将新创建的标签关联到当前注释
    await addTagToAnnotation({ 
      tagIds: [response.tagId], 
      markId: props.note.uuid 
    });
    
    // 3. 刷新全局标签列表
    annotationStore.refreshTagList();
    
    // 4. 构造标签对象并更新本地状态
    const tag: AnnotateTag = {
      tagId: response.tagId,
      tagName,
      latestUseTime: '',
    };
    localAddTag(tag);
    
  } catch (error) {
    console.error('创建或关联标签失败:', error);
    // 可以根据需要添加用户提示，比如 message.error()
  }
};

const selectTag = async (tagId: AnnotateTag['tagId']) => {
  await addTagToAnnotation({ tagIds: [tagId], markId: props.note.uuid });
  annotationStore.refreshTagList();
  const tag = annotationStore.tagList.find((item) => item.tagId === tagId)!;
  localAddTag(tag);
};

const removeTag = async (tagId: AnnotateTag['tagId']) => {
  await deleteTagFromAnnotation({
    tagIds: [tagId],
    markId: props.note.uuid,
  });
  annotationStore.refreshTagList();
  localRemoveTag(tagId);
};
</script>

<style lang="postcss">
.note-extra {
  display: flex;
  align-items: center;

  .note-tags {
    flex: 1;
  }
}
</style>
