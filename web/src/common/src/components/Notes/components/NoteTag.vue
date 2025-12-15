<template>
  <div class="note-tag-container px-4">
    <div class="note-tag-title">
      {{ $t('common.notes.tag') }}
    </div>
    <div class="note-tag-list">
      <div class="note-tag-item">
        <span
          class="note-tag-name"
          :class="{
            'note-tag-selected': noteTagSelected.length === 0,
          }"
          @click="clickAllTag()"
        >{{ $t('common.text.all').trim().toUpperCase() }}</span>
      </div>
      <div
        v-for="item in noteTagList"
        :key="item.tagId"
        class="note-tag-item"
      >
        <span
          v-if="noteTagEditId === item.tagId"
          class="note-tag-edit"
        >
          <input
            v-model.trim="noteTagEditName"
            :data-note-tag-id="item.tagId"
            @keypress.enter="submitNoteTag(item)"
            @blur="submitNoteTag(item)"
          >
        </span>
        <a-dropdown
          v-else
          :trigger="['contextmenu']"
          overlay-class-name="note-tag-dropdown"
        >
          <span
            class="note-tag-name"
            :class="{
              'note-tag-selected': noteTagSelected.includes(item.tagId),
            }"
            @click="clickNoteTag(item.tagId)"
          >
            {{ item.tagName }}
          </span>
          <template #overlay>
            <a-menu>
              <a-menu-item @click="editNoteTag(item)">
                {{
                  $t('common.text.rename')
                }}
              </a-menu-item>
              <a-menu-item
                class="delete-menu-item"
                @click="removeNoteTag(item)"
              >
                {{ $t('common.text.delete') }}
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { nextTick } from 'vue';
import { message, Modal } from 'ant-design-vue';
import lodash from 'lodash';
import { useNote } from '../useNote';
import { NoteTag } from '../types';
import { deleteAnnotateTag, renameAnnotateTag } from '@common/api/note';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const props = defineProps<{
  noteState: ReturnType<typeof useNote>;
}>();

const {
  noteTagList,
  noteTagSelected,
  noteTagEditId,
  noteTagEditName,
  fetchNoteList,
  fetchNoteTagList,
  fetchUserTagList,
} = props.noteState;

function clickNoteTag(tagId: NoteTag['tagId']) {
  if (noteTagSelected.value.includes(tagId)) {
    noteTagSelected.value = noteTagSelected.value.filter(
      (id: string) => id !== tagId
    );
  } else {
    noteTagSelected.value = [tagId];
  }

  fetchNoteList(1);
}

function clickAllTag() {
  noteTagSelected.value = [];
  fetchNoteList(1);
}

async function editNoteTag(tag: NoteTag) {
  noteTagEditId.value = '';
  noteTagEditName.value = tag.tagName;
  noteTagEditId.value = tag.tagId;
  await nextTick();
  const input: HTMLInputElement | null = document.querySelector(
    `[data-note-tag-id="${tag.tagId}"]`
  );
  input?.focus();
}

function stopEditNoteTag() {
  noteTagEditId.value = '';
  noteTagEditName.value = '';
}

const submitNoteTag = lodash.throttle(
  async (tag: NoteTag) => {
    if (tag.tagName === noteTagEditName.value) {
      stopEditNoteTag();
      return;
    }

    const notCnEnNoDash = /[^\u4E00-\u9FA5_\-a-zA-Z0-9]/;
    const match = notCnEnNoDash.exec(noteTagEditName.value);
    if (match) {
      message.warn(
        t('common.notes.wordings.noSpecial', {
          word: match[0],
        }) as string
      );
      return;
    }

    const restore = tag.tagName;
    tag.tagName = noteTagEditName.value;
    stopEditNoteTag();
    noteTagList.value = [...noteTagList.value];
    try {
      await renameAnnotateTag({
        tagId: tag.tagId,
        tagName: tag.tagName,
      });
    } catch (error) {
      tag.tagName = restore;
      noteTagList.value = [...noteTagList.value];
      return;
    }

    fetchUserTagList();
    fetchNoteList();
    message.success(t('common.notes.wordings.renameSuccess') as string);
  },
  500,
  { leading: true, trailing: false }
);

function removeNoteTag(tag: NoteTag) {
  Modal.confirm({
    title: t('common.notes.wordings.deleteTag', {
      tag: tag.tagName,
    }) as string,
    content: '',
    okText: t('common.text.delete') as string,
    okType: 'danger',
    async onOk() {
      try {
        await deleteAnnotateTag({ tagId: tag.tagId });
        message.success(t('common.notes.wordings.deleteTagSuccess') as string);
      } catch (error) {
        fetchNoteTagList();
        message.error(t('common.notes.wordings.deleteTagFail') as string);
        return;
      }

      await fetchNoteTagList();
      fetchNoteList();
      fetchUserTagList();
    },
  });
}

function cancelEdit() {
  noteTagEditId.value = '';
}
</script>

<style lang="less" scoped>
.note-tag-container {
  padding-bottom: 4px;
  padding-right: 10px;
  display: flex;
  .note-tag-title {
    padding-top: 12px;
    flex: 0 0 64px;
    font-size: 14px;
    color: var(--site-theme-text-secondary);
    > * {
      margin-left: 2px;
      &:hover {
        color: var(--site-theme-primary-color);
      }
    }
  }
  .note-tag-list {
    flex: 1 1 100%;
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: flex-start;
    .note-tag-item {
      margin-top: 12px;
      margin-right: 12px;
      .note-tag-name {
        height: 22px;
        padding: 2px 8px;
        display: inline-flex;
        justify-content: center;
        align-items: center;
        border-radius: 2px;
        font-size: 12px;
        cursor: pointer;
        background: var(--site-theme-bg-soft);
        color: var(--site-theme-primary-color);
        transition: all 0.2s ease;
        
        &:hover {
          background: var(--site-theme-primary-color-fade);
        }
        
        &.note-tag-selected {
          background: var(--site-theme-primary-color);
          color: var(--site-theme-text-inverse);
        }
      }
      .note-tag-edit {
        position: relative;
        height: 24px;
        display: inline-block;
        span {
          position: absolute;
          height: 100%;
          z-index: 10;
          display: flex;
          align-items: center;
          font-size: 12px;
          color: var(--site-theme-text-tertiary);
          top: 1px;
          left: 7px;
        }
        input {
          text-indent: 12px;
          border: 1px solid var(--site-theme-primary-color);
          border-radius: 2px;
          color: var(--site-theme-text-primary);
          background-color: var(--site-theme-bg-primary);
          height: 100%;
          font-size: 12px;
          outline: 0;
          
          &:focus {
            box-shadow: 0 0 0 2px var(--site-theme-primary-color-fade);
          }
        }
      }
    }
  }
}
</style>

<style lang="less">
.note-tag-dropdown {
  .ant-dropdown-menu {
    background-color: var(--site-theme-bg-primary);
    box-shadow: var(--site-theme-shadow-2);
    
    .ant-dropdown-menu-item {
      color: var(--site-theme-text-secondary);
      
      &:hover {
        background-color: var(--site-theme-background-hover);
        color: var(--site-theme-primary-color);
      }
      
      &.delete-menu-item {
        color: var(--site-theme-error);
      }
    }
  }
}
</style>
