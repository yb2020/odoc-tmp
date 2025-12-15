import { message } from 'ant-design-vue';
// import {
//   // GetPaperNoteBaseInfoByIdReq,
//   // GetOwnerPaperNoteBaseInfoReq,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/NoteInfo';
import { GetPaperNoteBaseInfoByIdReq, GetOwnerPaperNoteBaseInfoReq } from 'go-sea-proto/gen/ts/note/PaperNote'
import {
  UserStatusEnum,
  // GetPdfStatusInfoResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/OpenNote';
import {
  Module,
  Store as VuexStore,
  ActionTree,
  DispatchOptions,
  MutationTree,
  ActionContext,
  CommitOptions,
} from 'vuex';

import { $GroupProceed, getGroupNote } from '~/src/api/group';
import {
  getPdfStatusInfo,
  // GetPdfStatusInfoReq,
  getOwnerPaperNoteBaseInfo,
  getPaperNoteBaseInfoById,
  NoteBaseInfo,
} from '../../api/base';
import { GetPdfStatusInfoRequest, GetPdfStatusInfoResponse } from 'go-sea-proto/gen/ts/pdf/PaperPDF'
import { RootState } from '../types';
import { BaseState, GroupNoteBaseInfo, SELF_NOTEINFO_GROUPID } from './type';
import { useCopilotStore } from '~/src/stores/copilotStore';
import { useI18n } from 'vue-i18n';

export enum BaseActionTypes {
  GET_NOTEID = 'getOwnerPaperNoteBaseInfo',
  GET_STATUSINFO = 'getPdfStatusInfo',
  GET_SELF_NOTEINFO = 'getPaperNoteBaseInfoById',
  SWITCH_TO_GROUP = 'swtichToGroup',
}

export enum BaseMutationTypes {
  SET_STATUS_INFO = 'setPdfStatusInfo',
  SET_NOTE_INFO = 'setNoteInfo',
  SET_COLLECTED = 'setCollected',
  SET_CURRENT_GROUP_ID = 'setCurrentGroupId',
  SET_GROUP_INFOS = 'setGroupInfos',
  SET_PAPER_ID = 'setPaperId',
}

type AugmentedActionContext = {
  commit<K extends keyof Mutations>(
    key: K,
    payload: Parameters<Mutations[K]>[1]
  ): ReturnType<Mutations[K]>;
} & Omit<ActionContext<BaseState, RootState>, 'commit'>;

interface Actions {
  [BaseActionTypes.GET_NOTEID](
    context: AugmentedActionContext,
    payload: GetOwnerPaperNoteBaseInfoReq
  ): Promise<string>;
  [BaseActionTypes.GET_STATUSINFO](
    context: AugmentedActionContext,
    payload: GetPdfStatusInfoRequest
  ): Promise<void>;
  [BaseActionTypes.GET_SELF_NOTEINFO](
    context: AugmentedActionContext,
    payload: GetPaperNoteBaseInfoByIdReq
  ): Promise<string>;
  [BaseActionTypes.SWITCH_TO_GROUP](
    context: AugmentedActionContext,
    payload: {
      groupId: string;
      t: ReturnType<typeof useI18n>['t'];
    }
  ): Promise<boolean>;
}

const showSwitchTip = async (
  state: BaseState,
  noteInfo: NoteBaseInfo,
  groupId: string,
  t: ReturnType<typeof useI18n>['t']
) => {
  if (
    state.currentGroupId &&
    state.noteInfoMap[state.currentGroupId]?.noteId !== noteInfo.noteId
  ) {
    const groupInfo = state.groupInfoMap[groupId];
    message.info(
      groupId === SELF_NOTEINFO_GROUPID
        ? t('message.switchPersonal')
        : `${t('message.switchTeam')} [${groupInfo.name}]`
    );
  }
};

const actions: ActionTree<BaseState, RootState> & Actions = {
  async [BaseActionTypes.GET_NOTEID](_, params) {
    const res = await getOwnerPaperNoteBaseInfo(params);
    return res.noteId;
  },
  async [BaseActionTypes.GET_STATUSINFO]({ commit }, params) {
    const res = await getPdfStatusInfo(params);
    commit(BaseMutationTypes.SET_STATUS_INFO, res);
  },
  async [BaseActionTypes.GET_SELF_NOTEINFO]({ commit }, params) {
    const res = await getPaperNoteBaseInfoById(params);
    commit(BaseMutationTypes.SET_NOTE_INFO, {
      groupId: SELF_NOTEINFO_GROUPID,
      noteInfo: res,
    });
    commit(BaseMutationTypes.SET_PAPER_ID, res.paperId);
    const copilotStore = useCopilotStore();
    copilotStore.setGptInfo(res.gptGrayTip);
    return res.pdfId;
  },

  async [BaseActionTypes.SWITCH_TO_GROUP]({ commit, state }, { groupId, t }) {
    let noteInfo = state.noteInfoMap[groupId];

    if (noteInfo) {
      showSwitchTip(state, noteInfo, groupId, t);
      commit(BaseMutationTypes.SET_CURRENT_GROUP_ID, groupId);
      return true;
    }

    if (groupId === SELF_NOTEINFO_GROUPID) {
      noteInfo = state.noteInfoMap[SELF_NOTEINFO_GROUPID];

      showSwitchTip(state, noteInfo, groupId, t);

      commit(BaseMutationTypes.SET_CURRENT_GROUP_ID, SELF_NOTEINFO_GROUPID);
      return true;
    }

    const data = await getGroupNote({
      groupId,
      paperId: state.noteInfoMap[SELF_NOTEINFO_GROUPID].paperId,
      pdfId: state.noteInfoMap[SELF_NOTEINFO_GROUPID].pdfId,
    });

    if (!data.hasPaperInGroup) {
      commit(BaseMutationTypes.SET_CURRENT_GROUP_ID, groupId);
      return false;
    }

    noteInfo = {
      ...state.noteInfoMap[SELF_NOTEINFO_GROUPID],
      noteId: data.groupNoteId,
    };

    if (data.pdfVersionSwitched) {
      noteInfo.pdfId = data.groupPdfId;
      noteInfo.sourceMark = data.sourceMark;
      noteInfo.crawlUrl = data.crawlUrl;
      noteInfo.licenceType = data.licenceType;
      noteInfo.pdfUrl = data.pdfUrl;
      noteInfo.showAnnotation = true;
    }

    commit(BaseMutationTypes.SET_NOTE_INFO, { groupId, noteInfo });

    showSwitchTip(state, noteInfo, groupId, t);

    commit(BaseMutationTypes.SET_CURRENT_GROUP_ID, groupId);

    return true;
  },
};

export type Mutations<S = BaseState> = {
  [BaseMutationTypes.SET_STATUS_INFO](
    state: S,
    payload: GetPdfStatusInfoResponse
  ): void;
  [BaseMutationTypes.SET_NOTE_INFO](state: S, payload: GroupNoteBaseInfo): void;
  [BaseMutationTypes.SET_COLLECTED](state: S, payload: boolean): void;
  [BaseMutationTypes.SET_CURRENT_GROUP_ID](state: S, groupId: string): void;
  [BaseMutationTypes.SET_GROUP_INFOS](state: S, payload: $GroupProceed[]): void;
  [BaseMutationTypes.SET_PAPER_ID](
    state: S,
    payload: NoteBaseInfo['paperId']
  ): void;
};

const mutations: MutationTree<BaseState> & Mutations = {
  [BaseMutationTypes.SET_STATUS_INFO](state, payload) {
    state.statusInfo = payload;
  },
  [BaseMutationTypes.SET_NOTE_INFO](state, payload) {
    state.noteInfoMap[payload.groupId] = payload.noteInfo;
  },
  [BaseMutationTypes.SET_COLLECTED](state, payload) {
    const selfNoteInfo = state.noteInfoMap[SELF_NOTEINFO_GROUPID];
    if (selfNoteInfo) {
      selfNoteInfo.isCollected = !!payload;
    }
  },
  [BaseMutationTypes.SET_CURRENT_GROUP_ID](state, payload) {
    state.currentGroupId = payload;
  },
  [BaseMutationTypes.SET_GROUP_INFOS](state, payload) {
    const map: Record<string, $GroupProceed> = {};
    payload.forEach((item) => {
      map[item.id] = item;
    });
    state.groupInfoMap = map;
  },
  [BaseMutationTypes.SET_PAPER_ID](state, payload) {
    state.paperId = payload;
  },
};

export type BaseStore<S = BaseState> = Omit<
  VuexStore<S>,
  'getters' | 'dispatch' | 'commit'
> & {
  commit<K extends keyof Mutations, P extends Parameters<Mutations[K]>[1]>(
    key: K,
    payload: P,
    options?: CommitOptions
  ): ReturnType<Mutations[K]>;
} & {
  dispatch<K extends keyof Actions>(
    key: K,
    payload: Parameters<Actions[K]>[1],
    options?: DispatchOptions
  ): ReturnType<Actions[K]>;
};

export const BaseModule: Module<BaseState, RootState> = {
  namespaced: true,
  mutations,
  actions,
  state: () => ({
    paperId: '',
    statusInfo: {
      pdfOnlineStatus: null,
      pdfRenderedPrivateStatus: null,
      pdfUserStatus: UserStatusEnum.UNRECOGNIZED,
      noteUserStatus: UserStatusEnum.UNRECOGNIZED,
      pdfUrl: '',
      paperId: '',
      paperTitle: '',
      isLike: false,
      pdfLikeCount: '',
      authPdfId: '',
      hasReadPermission: false,
      hasPdfAccessFlag: false,
      noteOpenAccessFlag: false,
      docName: '',
    },
    currentGroupId: '',
    noteInfoMap: {},
    groupInfoMap: {},
    tagList: [],
  }),
};
