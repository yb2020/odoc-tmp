<template>
  <div
    class="extract-card-wrap group w-full relative px-5 pt-4 pb-6 rounded-md border border-t-4 border-solid"
    :style="{
      '--border-color-hover':
        noteStyleMap[noteMap[data?.uuid]?.styleId || '']?.color,
      'background-color': 'var(--site-theme-bg-primary)',
    }"
  >
    <close-circle-filled
      class="close-btn !text-rp-neutral-6 text-base absolute -top-2 -right-2 group-hover:block"
      @click="removeNote(data?.uuid)"
    />
    <div
      v-if="noteQuote && data?.select"
      class="note-item-select text-xs"
      :title="noteEditId ? '' : $t('common.notes.wordings.edit')"
      @click="handleEdit"
      @mousedown="mousedownReference(data)"
      @mouseup="mouseupReference(data)"
      v-html="data?.select.rectStr"
    />
    <div
      v-if="noteQuote && data?.rect"
      class="note-data-rect"
    >
      <img
        class="note-data-rect-image w-full"
        :src="data?.rect.picUrl"
        @click="modalNote = data"
      >
    </div>
    <idea-markdown
      :title="noteEditId ? '' : $t('common.notes.wordings.edit')"
      :raw="noteMap[data?.uuid]?.idea || ''"
      :uniq-id="data?.uuid"
      :editing="isEdit"
      :upload="upload"
      @clickView="handleEdit"
      @blur="submitNote(data, $event)"
      @submit="submitNote(data, $event)"
    />
    <idea-tag
      class="tag-wrap extract-card-tag-input"
      :list="userTagList || []"
      :selected="data?.tags"
      :tag-label="$t('viewer.tag')"
      @inputTag="inputTag(data, $event)"
      @selectTag="selectTag(data, $event)"
      @remove="removeTag(data, $event)"
      @click.stop
      @mouseenter="initTags"
      @focus="adjustTagDropdownPosition"
      @input="adjustTagDropdownPosition"
    />
    <div class="note-data-footer w-80">
      <!-- <div class="note-data-date">
        {{ noteMap[data.uuid]?.date }}
      </div> -->
      <div
        class="note-data-source-visit whitespace-nowrap overflow-hidden text-ellipsis cursor-pointer text-xs text-[#86919C]"
        @click="gotoPdf(data)"
      >
        <span>{{ data?.paperTitle || data?.docName }}</span>
      </div>
    </div>
    <a-dropdown
      trigger="click"
      overlay-class-name="note-item-style-overlay"
      :destroy-popup-on-hide="true"
      :get-popup-container="getPopupContainer"
    >
      <div class="note-item-style">
        <div
          :style="{
            background: noteStyleMap[noteMap[data?.uuid]?.styleId || '']?.color,
          }"
        />
      </div>
      <template #overlay>
        <a-menu class="note-item-style-overlay">
          <a-menu-item
            v-for="style in noteStyleList"
            :key="style.type"
            @click="changeNoteStyle(data, style.type)"
          >
            <em :style="{ background: style.color }" />
          </a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>
    <a-modal
      :centered="true"
      :visible="!!modalNote"
      :footer="null"
      :width="null"
      :wrap-class-name="`note-item-rect-modal ${
        noteSearchKeyword ? 'note-list-searching' : ''
      }`"
      @cancel="modalNote = null"
    >
      <div>
        <div class="note-item-rect-zoom">
          <img
            v-if="modalNote"
            :src="modalNote.rect?.picUrl || ''"
          >
        </div>
        <div class="note-item-rect-idea">
          <div
            v-if="modalNote"
            class="note-item-rect-idea-content"
            v-html="noteMap[modalNote.uuid].ideaMarkdown"
          />
          <a-tooltip
            placement="bottom"
            overlay-class-name="note-source-tips"
            :overlay-style="{
              whiteSpace: 'nowrap',
              maxWidth: '999px',
            }"
            :arrow-point-at-center="true"
            :get-popup-container="getPopupContainer"
            :title="
              $t('common.notes.wordings.jumpPaper', {
                name: modalNote && (modalNote.paperTitle || modalNote.docName),
                num: modalNote && modalNote.pageNumber,
              })
            "
          >
            <div
              class="note-item-source-link cursor-pointer"
              @click="modalNote && gotoPdf(modalNote)"
            >
              <ArrowUpOutlined />
              <span>{{ $t('common.notes.wordings.sourcePaper') }}</span>
            </div>
          </a-tooltip>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { computed, nextTick, ref } from 'vue';
import { CloseCircleFilled, ArrowUpOutlined } from '@ant-design/icons-vue';
import { message, Modal } from 'ant-design-vue';
import {
  IdeaMarkdown,
  IdeaTag,
  markdownToHtml,
} from '@idea/aiknowledge-markdown';
import { ImageStorageType, uploadImage } from '@common/api/upload';
import { noteStyleMap, noteStyleList, useNote } from '../useNote';
import { ColorKey } from '../types';
import {
  NoteAnnotation,
  deleteAnnotation,
  editAnnotation,
  createAnnotationTag,
  addTagToAnnotation,
  deleteTagFromAnnotation,
} from '~/src/api/note';
// import { createAnnotationTag } from '~/src/api/annotations';
import { useI18n } from 'vue-i18n';
import { AnnotateTag } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';

const { t } = useI18n();

const props = defineProps<{
  noteState: ReturnType<typeof useNote>;
  data: NoteAnnotation;
}>();

const {
  noteQuote,
  noteEditId,
  userTagList,
  refreshNoteExplorer,

  noteSearchKeyword,
  noteList,
  noteEditContent,
  fetchNoteTagList,
  fetchUserTagList,
} = props.noteState;

const isEdit = ref<boolean>(false);
const textareaRef = ref();

function getPopupContainer(triggerNode: HTMLElement) {
  return triggerNode.parentNode;
}

const upload = async (src: File | string) => {
  return uploadImage(src, ImageStorageType.markdown);
};

const handleEdit = () => {
  isEdit.value = true;
  nextTick(() => {
    textareaRef.value?.focus();
  });
};

const getNoteBody = (note: NoteAnnotation) => {
  if ((note.select && note.rect) || (!note.select && !note.rect)) {
    throw new Error('笔记内容有误');
  }

  return note.select || note.rect!;
};

interface NoteExtra {
  idea: string;
  ideaMarkdown: string;
  styleId: ColorKey;
}
const modalNote = ref<NoteAnnotation | null>(null);
const noteMap = computed(() => {
  const map: Record<NoteAnnotation['uuid'], NoteExtra> = {};

  noteList.value.forEach((note: NoteAnnotation) => {
    const { idea = '', styleId = ColorKey.blue } =
      note.select ?? note.rect ?? {};

    const ideaMarkdown = markdownToHtml(idea);

    map[note.uuid] = {
      idea: noteSearchKeyword.value ? removeEmTag(idea) : idea,
      ideaMarkdown,
      styleId,
    };
  });

  return map;
});
const removeEmTag = (content: string) => {
  return content.replace(/(<em>)|(<\/em>)/g, '');
};

function editNote(note: NoteAnnotation) {
  noteEditContent.value = noteMap.value[note.uuid].idea;
  noteEditId.value = note.uuid;
}

let mouseReferenceId = '';
let mouseReferenceTime = -Infinity;
function mousedownReference(note: NoteAnnotation) {
  mouseReferenceTime = Date.now();
  mouseReferenceId = note.uuid;
}

function mouseupReference(note: NoteAnnotation) {
  if (mouseReferenceId === note.uuid && Date.now() - mouseReferenceTime < 500) {
    editNote(note);
  }

  mouseReferenceId = '';
  mouseReferenceTime = -Infinity;
}

async function submitNote(note: NoteAnnotation, newIdea: string) {
  const { idea } = noteMap.value[note.uuid];
  isEdit.value = false;
  if (idea === newIdea) {
    return;
  }

  const setIdea = (value: string) => {
    const noteBody = getNoteBody(note);
    noteBody.idea = value;
    noteList.value = [...noteList.value];
  };
  setIdea(newIdea);
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { uuid, ...annotation } = note;

  try {
    await editAnnotation(annotation);
  } catch (error) {
    message.error(t('common.notes.wordings.saveNoteFail') as string);
    setIdea(idea);
    return;
  }

  message.success(t('common.notes.wordings.saveNoteSuccess') as string);
  await fetchNoteTagList();
}

async function changeNoteStyle(note: NoteAnnotation, newStyleId: ColorKey) {
  const noteBody = getNoteBody(note);
  const { styleId } = noteBody;
  noteBody.styleId = newStyleId;
  const { uuid, ...annotation } = note;
  try {
    await editAnnotation(annotation);
  } catch (error) {
    message.error(t('common.notes.wordings.saveColorFail') as string);
    noteBody.styleId = styleId;
    noteList.value = [...noteList.value];
  }
}

function removeNote(noteId: NoteAnnotation['uuid']) {
  Modal.confirm({
    title: t('common.notes.wordings.deleteNote') as string,
    content: '',
    okText: t('common.text.delete') as string,
    okType: 'danger',
    async onOk() {
      try {
        await deleteAnnotation(noteId);
      } catch (error) {
        message.error(t('common.notes.wordings.deleteNoteFail') as string);
        return;
      }

      refreshNoteExplorer();
      message.success(t('common.notes.wordings.deleteNoteSuccess') as string);
    },
  });
}

function gotoPdf(note: NoteAnnotation) {
  const { documentId = '' } = note.select ?? note.rect ?? {};
  gotoPdfPage({
    noteId: documentId,
    annotationId: note.uuid,
    pageNumber: String(note.pageNumber),
    ...{
      from_extension: String(true),
    },
  });
}

const gotoPdfPage = (params: Record<string, string>, path?: string) => {
  const url = new URL(
    //`${window.location.origin}/pdf-annotate${path || '/note'}` //旧版的跳转地址
    `${window.location.origin}/note` //新版的跳转地址
  );

  for (const key of Object.keys(params)) {
    url.searchParams.set(key, params[key]);
  }

  handleJump(url.toString());
};

const handleJump = (path: string) => {
  const a = document.createElement('a');

  a.setAttribute('href', path);

  a.setAttribute('target', '_blank');

  a.setAttribute('id', 'open-new-page');

  document.body.appendChild(a);

  a.click();

  const target = document.getElementById('open-new-page');

  if (target) {
    document.body.removeChild(target);
  }
};

function initTags() {
  if (!Array.isArray(userTagList.value)) {
    fetchUserTagList();
  }
}

// 动态调整标签下拉框位置
function adjustTagDropdownPosition() {
  nextTick(() => {
    // 使用更精准的选择器
    const tagComponent = document.querySelector('.extract-card-tag-input');
    const tagEntryContainer = tagComponent?.querySelector('.idea-tag-entry-container');
    const tagSelectContainer = document.querySelector('.idea-tag-select-container');
    
    if (tagEntryContainer && tagSelectContainer) {
      const entryRect = tagEntryContainer.getBoundingClientRect();
      const selectElement = tagSelectContainer as HTMLElement;
      
      // 精准定位输入框元素
      const inputElement = tagEntryContainer.querySelector('input');
      if (inputElement) {
        inputElement.style.width = '120px';
        inputElement.style.maxWidth = '120px';
        inputElement.style.minWidth = '120px';
        inputElement.style.boxSizing = 'border-box';
        inputElement.style.padding = '1px 3px';
        inputElement.style.border = 'none';
        inputElement.style.margin = '0';
        inputElement.style.outline = 'none';
        inputElement.style.boxShadow = 'none';
      }
      
      // 设置下拉框位置和盒模型
      selectElement.style.position = 'fixed';
      selectElement.style.top = `${entryRect.bottom + 2}px`;
      selectElement.style.left = `${entryRect.left}px`;
      selectElement.style.width = '120px';
      selectElement.style.maxWidth = '120px';
      selectElement.style.minWidth = '120px';
      selectElement.style.boxSizing = 'border-box';
      selectElement.style.padding = '0';
      selectElement.style.border = '1px solid';
      selectElement.style.margin = '0';
      selectElement.style.zIndex = '9999';
    }
  });
}

async function inputTag(note: NoteAnnotation, tagName: AnnotateTag['tagName']) {
  // 1. 创建新标签
  const response = await createAnnotationTag({
    tagName,
    markId: note.uuid,
  });
  
  // 2. 将新标签关联到当前标注（与selectTag保持一致）
  await addTagToAnnotation({
    tagIds: [response.tagId],
    markId: note.uuid,
  });
  
  // 3. 刷新用户标签列表
  fetchUserTagList();
  
  // 4. 刷新标注标签列表（与selectTag保持一致）
  fetchNoteTagList();
  
  // 5. 更新本地标签数组
  note.tags.push({
    tagId: response.tagId,
    tagName,
    latestUseTime: '0',
  });
}

async function selectTag(note: NoteAnnotation, tagId: AnnotateTag['tagId']) {
  await addTagToAnnotation({
    tagIds: [tagId],
    markId: note.uuid,
  });
  const tag = userTagList.value!.find(
    (item: AnnotateTag) => item.tagId === tagId
  )!;
  fetchNoteTagList();
  note.tags.push({ ...tag });
}

async function removeTag(note: NoteAnnotation, tagId: AnnotateTag['tagId']) {
  await deleteTagFromAnnotation({
    tagIds: [tagId],
    markId: note.uuid,
  });
  fetchNoteTagList();
  const index = note.tags.findIndex((tag: AnnotateTag) => tag.tagId === tagId);
  if (index !== -1) {
    note.tags.splice(index, 1);
  }
}
</script>

<style lang="less" scoped>
.extract-card-wrap {
  width: 368px;
  position: relative;
  z-index: 1;

  .close-btn {
    display: none;
  }
  border-color: var(--border-color-hover) transparent transparent transparent;
  &:hover {
    .close-btn {
      display: block;
    }
    border-color: var(--border-color-hover) !important;
    z-index: 10;
  }
  .note-item-select {
    color: #4e5969;
    margin-bottom: 16px;
    line-height: 20px;
    background: #f0f2f5;
    font-style: italic;
    cursor: pointer;
    word-break: break-all;
    word-wrap: break-word;
  }
  .note-data-rect {
    overflow: hidden;
    margin-bottom: 12px;
    margin-bottom: 8px;
    display: block;
  }
  .note-item-idea {
    min-height: 16px;
    cursor: pointer;
  }
  .note-item-footer {
    padding-bottom: 12px;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    .note-item-date {
      flex: 0 0 100px;
      font-size: 12px;
      color: #86919c;
      height: 20px;
      line-height: 20px;
    }
    .note-item-source-visit {
      flex: 1 1 100%;
      text-align: right;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      margin-left: 72px;
      margin-right: 20px;
      color: #86919c;
      cursor: pointer;
      > i {
        visibility: hidden;
        transform: rotate(45deg);
        margin-right: 8px;
      }
      &:hover {
        color: #4e5969;
        > i {
          visibility: visible !important;
        }
      }
    }
  }
  .note-item-delete {
    position: absolute;
    top: -14px;
    right: -14px;
    width: 30px;
    height: 30px;
    color: #86919c;
    font-size: 16px;
    z-index: 10;
    pointer-events: none;
    opacity: 0;
    transition: opacity 0.1s;
    transition-delay: 0.24s;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  &:hover {
    background: #f7f8fa;
    .note-item-delete {
      pointer-events: initial !important;
      opacity: 1;
      cursor: pointer;
    }
  }

  .note-item-style {
    position: absolute;
    right: 0;
    bottom: 45px;
    height: 30px;
    width: 40px;
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: pointer;
    > div {
      height: 10px;
      width: 10px;
      border-radius: 5px;
    }
  }

  .note-item-rect-image {
    max-height: 64px;
    max-width: 440px;
    cursor: pointer;
  }
}

:deep(.note-item-style-overlay) {
  .ant-dropdown-link {
    color: #4e5969;
  }
  .ant-dropdown-menu-item {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    font-size: 14px;
    color: #4e5969;
    height: 40px;
    width: 40px;
    transition: none;
    text-align: center;
    &:hover {
      background: #e3eef5;
    }
    .ant-dropdown-menu-title-content {
      > em:first-of-type {
        display: inline-block;
        width: 10px;
        height: 10px;
        flex: 0 0 10px;
        border-radius: 5px;
      }
    }
  }
}
</style>
<style lang="less">
.note-item-idea {
  color: var(--site-theme-text-primary);
  > p:last-of-type {
    margin-bottom: 0;
  }
}

.note-item-idea-tag {
  color: var(--site-theme-primary-color);
}
.note-list-searching {
  .note-item-idea-tag {
    color: var(--site-theme-text-tertiary) !important;
  }
}

.note-item-highlight,
.note-item-idea em,
.note-item-rect-idea-content em,
.note-item-select em {
  color: var(--site-theme-primary-color);
  background: var(--site-theme-primary-color-fade);
}

.note-item-source-link {
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  padding: 2px 8px;
  position: absolute;
  width: 84px;
  height: 24px;
  border-radius: 2px;
  font-size: 12px;
  cursor: pointer;
  color: var(--site-theme-text-tertiary);
  &:hover {
    color: var(--site-theme-primary-color);
  }
  &:active {
    background: var(--site-theme-background-hover);
  }
  > i {
    transform: rotate(45deg);
    margin-right: 8px;
  }
}

.note-item-rect-modal {
  .ant-modal-body {
    padding: 0;
  }

  .ant-modal-close {
    right: -32px;
    .ant-modal-close-x {
      width: 32px;
      height: 32px;
      background: var(--site-theme-primary-color);
      color: var(--site-theme-text-inverse);
      display: flex;
      justify-content: center;
      align-items: center;
    }
  }

  .note-item-source-link {
    right: 16px;
    bottom: 12px;
  }

  .note-item-rect-zoom {
    min-width: 400px;
    min-height: 240px;
    display: flex;
    justify-content: center;
    align-items: center;
    background: var(--site-theme-bg-soft);

    > img {
      max-width: 800px;
      max-height: 480px;
    }
  }

  .note-item-rect-idea {
    color: var(--site-theme-text-primary);
    padding-top: 18px;
    padding-bottom: 38px;
    padding-left: 24px;
    padding-right: 24px;
  }
}

.note-item-input {
  margin-bottom: 8px;
  .note_edit_container {
    margin-top: 0 !important;
    .w-e-text-container {
      padding: 3px 0 !important;
      border: 1px solid var(--site-theme-primary-color);
      border-radius: 2px;
    }
  }
}
</style>
<style lang="less" scoped>
.extract-card-wrap {
  :deep(.idea-markdown-view-container) {
    color: var(--site-theme-text-primary);
    font-family: PingFang SC;
    font-size: 14px;
    font-style: normal;
    font-weight: 600;
    line-height: 22px;
    overflow: hidden;
    /* Force break long words and URLs to prevent overflow */
    overflow-wrap: break-word !important;
    word-break: break-all !important;

    p {
      margin-bottom: 0;
    }
    img {
      max-width: 100%;
    }
  }
  .idea-markdown-edit-container {
    transition: none !important;
  }

  .idea-tag-entry-container {
    outline: 0;
  }

  .idea-tag-input {
    border: none !important;
    outline: none !important;
    box-shadow: none !important;
    width: 120px !important;
    max-width: 120px !important;
    min-width: 120px !important;
    box-sizing: border-box !important;
    padding: 1px 3px !important;
    margin: 0 !important;
  }

  .idea-tag-input:focus {
    border: none !important;
    outline: none !important;
    box-shadow: none !important;
  }

  .idea-tag-input:focus-visible {
    border: none !important;
    outline: none !important;
    box-shadow: none !important;
  }

  .idea-markdown-edit-container,
  .idea-tag-input {
    color: var(--site-theme-text-primary) !important;
    background-color: var(--site-theme-bg-primary) !important;
  }

  :deep(.idea-tag-entry-container) {
    position: relative !important;
    width: 100% !important;
    max-width: 100% !important;
    min-width: auto !important;
  }

  :deep(.idea-tag-entry-container input) {
    width: 100% !important;
    max-width: 110px !important;
    box-sizing: border-box !important;
  }

  :deep(.idea-tag-entry-container input[type="text"]) {
    width: 100% !important;
    max-width: 110px !important;
    box-sizing: border-box !important;
  }

  /* 精准定位标签输入框 - 覆盖所有可能的样式 */
  :deep(.extract-card-tag-input .idea-tag-entry-container input),
  :deep(.extract-card-tag-input .idea-tag-entry-container input[type="text"]),
  :deep(.extract-card-tag-input .idea-tag-entry-container .idea-tag-input),
  :deep(.extract-card-tag-input input),
  :deep(.extract-card-tag-input input[type="text"]),
  :deep(.extract-card-tag-input .idea-tag-input) {
    width: 120px !important;
    max-width: 120px !important;
    min-width: 120px !important;
    box-sizing: border-box !important;
    padding: 1px 3px !important;
    border: none !important;
    margin: 0 !important;
    outline: none !important;
    box-shadow: none !important;
  }

  /* 修复输入框焦点状态的蓝色边框 - 覆盖所有可能的焦点状态 */
  :deep(.extract-card-tag-input .idea-tag-entry-container input:focus),
  :deep(.extract-card-tag-input .idea-tag-entry-container input:focus-visible),
  :deep(.extract-card-tag-input .idea-tag-entry-container input:active),
  :deep(.extract-card-tag-input .idea-tag-entry-container input[type="text"]:focus),
  :deep(.extract-card-tag-input .idea-tag-entry-container input[type="text"]:focus-visible),
  :deep(.extract-card-tag-input .idea-tag-entry-container input[type="text"]:active),
  :deep(.extract-card-tag-input .idea-tag-entry-container .idea-tag-input:focus),
  :deep(.extract-card-tag-input .idea-tag-entry-container .idea-tag-input:focus-visible),
  :deep(.extract-card-tag-input .idea-tag-entry-container .idea-tag-input:active),
  :deep(.extract-card-tag-input input:focus),
  :deep(.extract-card-tag-input input:focus-visible),
  :deep(.extract-card-tag-input input:active),
  :deep(.extract-card-tag-input input[type="text"]:focus),
  :deep(.extract-card-tag-input input[type="text"]:focus-visible),
  :deep(.extract-card-tag-input input[type="text"]:active),
  :deep(.extract-card-tag-input .idea-tag-input:focus),
  :deep(.extract-card-tag-input .idea-tag-input:focus-visible),
  :deep(.extract-card-tag-input .idea-tag-input:active) {
    border: none !important;
    outline: none !important;
    box-shadow: none !important;
  }

  :deep(.idea-tag-select-container) {
    position: fixed !important;
    z-index: 9999 !important;
    background-color: var(--site-theme-bg-primary) !important;
    color: var(--site-theme-text-primary) !important;
    border: none !important;
    border-radius: 4px !important;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
    max-height: 200px !important;
    overflow-y: auto !important;
    width: 120px !important;
    min-width: 120px !important;
    max-width: 120px !important;
    box-sizing: border-box !important;
    transform: none !important;
    margin: 0 !important;
    padding: 0 !important;
    top: auto !important;
    left: auto !important;
    right: auto !important;
    bottom: auto !important;
  }

  :deep(.idea-tag-select-item) {
    height: 32px !important;
    line-height: 32px !important;
    padding-left: 12px !important;
    padding-right: 12px !important;
    overflow: hidden !important;
    text-overflow: ellipsis !important;
    white-space: nowrap !important;
    cursor: pointer !important;

    &:hover {
      background-color: var(--site-theme-bg-soft) !important;
    }
  }

  .idea-tag-select-cursor {
    background: var(--site-theme-bg-soft);
    color: var(--site-theme-text-primary);
  }

  .idea-tag-entry-container {
    position: relative !important;
    margin-top: 3px;
    // outline: solid 1px lightblue;
    display: flex;
    justify-content: flex-start;
    align-items: flex-start;
    flex-wrap: wrap;
    width: 100% !important;
    max-width: 100% !important;
    min-width: auto !important;
  }

  :deep(.idea-tag-entry-container .idea-tag-add) {
    visibility: hidden !important;
    display: inline-flex !important;
    justify-content: center !important;
    align-items: center !important;
    width: 60px !important;
    cursor: pointer !important;
    border-radius: 2px !important;
    margin: 0 4px 4px 0 !important;
    padding: 2px 8px !important;
    font-size: 14px !important;
    background: var(--site-theme-bg-soft) !important;
    color: var(--site-theme-text-tertiary) !important;
    opacity: 0 !important;
    transition: opacity 0.2s ease !important;
    box-sizing: border-box !important;
    white-space: nowrap !important;
  }

  &:hover :deep(.idea-tag-entry-container .idea-tag-add) {
    visibility: visible !important;
    opacity: 1 !important;
  }

  .idea-tag-entry-container .idea-tag-add > i {
    margin-right: 2px;
  }

  .idea-tag-entry-container .idea-tag-add > i > * {
    font-size: 10px;
  }

  .idea-tag-entry-edit {
    position: relative;
    z-index: 1;
    display: inline-block;
    height: 20px;
    top: -2px;
    margin-bottom: 4px;
  }

  .idea-tag-input {
    display: inline-block;
    box-sizing: border-box;
    padding: 1px 3px;
    width: 100% !important;
    max-width: 110px !important;
    height: 22px;
    border-radius: 2px;
    outline: 0;
  }

  :deep(.idea-tag-item-container) {
    display: inline-flex !important;
    align-items: center !important;
    justify-content: space-between !important;
    background: var(--site-theme-bg-soft);
    color: var(--site-theme-text-tertiary);
    margin: 0 4px 4px 0;
    padding: 2px 8px;
    max-width: 164px;
    border-radius: 2px;
    overflow: hidden;
    white-space: nowrap !important;
    box-sizing: border-box !important;
  }

  :deep(.idea-tag-item-name) {
    flex: 1 1 auto;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    min-width: 0 !important;
  }

  :deep(.idea-tag-item-remove) {
    flex-shrink: 0 !important;
    cursor: pointer;
    font-size: 10px;
    margin-left: 4px;
    display: inline-flex !important;
    align-items: center !important;
    justify-content: center !important;
    width: auto !important;
    height: auto !important;
  }

  .tag-wrap {
    margin-top: 24px;
    margin-right: 8px;
  }
}
</style>
