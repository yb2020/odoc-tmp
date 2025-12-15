import { createStore } from 'vuex';
import { RootState } from './types';
import { BaseStore, BaseModule } from './base';
import { DocumentsModule, DocumentsStore } from './documents';
import { UserModule } from './user';
import { ParseModule } from './parse';
import { NoteModule } from './note';
import { ShortcutsModule } from './shortcuts';
import { CertModule } from './cert';
import { computed } from 'vue';
import { UserStatusEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/OpenNote';
import { SELF_NOTEINFO_GROUPID } from './base/type';

// 定义本地枚举，避免依赖外部模块
// export enum UserStatusEnum {
//   UNKNOWN = 0,
//   NORMAL = 1,
//   ABNORMAL = 2,
//   DELETED = 3,
//   TOURIST = 4,
//   OWNER = 5,
//   GUEST = 6,
// }

const createVueStore = () => {
  const store = createStore<RootState>({
    modules: {
      base: BaseModule,
      documents: DocumentsModule,
      user: UserModule,
      parse: ParseModule,
      note: NoteModule,
      shortcuts: ShortcutsModule,
      cert: CertModule,
    },
  });

  // just for debug
  window.store = store;
  return store;
};

export type Store = BaseStore<Pick<RootState, 'base'>> &
  DocumentsStore<Pick<RootState, 'documents'>> &
  DocumentsStore<Pick<RootState, 'user'>> &
  DocumentsStore<Pick<RootState, 'note'>> &
  DocumentsStore<Pick<RootState, 'shortcuts'>> &
  DocumentsStore<Pick<RootState, 'cert'>>;

export const store = createVueStore();

export function useStore() {
  return store;
}

export const selfNoteInfo = computed(() => {
  return store.state.base.noteInfoMap[SELF_NOTEINFO_GROUPID];
});

export const pdfStatusInfo = computed(() => {
  return store.state.base.statusInfo;
});

export const currentGroupId = computed(() => {
  return store.state.base.currentGroupId;
});

export const currentNoteInfo = computed(() => {
  return store.state.base.noteInfoMap[currentGroupId.value];
});

export const isOwner = computed(
  () => pdfStatusInfo.value.noteUserStatus === UserStatusEnum.OWNER
);

export const isSelfNoteInfo = computed(
  () => currentGroupId.value === SELF_NOTEINFO_GROUPID
);

export const ownNoteOrVisitSharedNote = computed(
  () => isOwner.value || pdfStatusInfo.value.noteOpenAccessFlag
);
