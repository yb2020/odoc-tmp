import {
  addComment,
  getGroupComments,
  updateComment,
  deleteComment,
} from './../../../api/groupNote';
import { ActionContext, ActionTree } from 'vuex';
import { RootState } from '../../types';
import { NoteMutationTypes, Mutations } from '../mutations';
import { NoteState } from '../types';
import {
  GroupNoteCommentCreateReq,
  GroupNoteCommentUpdateReq,
  IncrementalGroupCommentReq,
  GroupNoteCommentDeleteReq,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/GroupNote';
import { NeedVipException } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { ResponseError } from '~/src/api/type';
import { ElementName } from '~/src/api/report';
import { ERROR_CODE_NEED_VIP } from '@common/api/const';
import { useVipStore } from '@common/stores/vip';

export enum NoteActionTypes {
  GET_GROUP_COMMENTS = 'get_group_comments',
  ADD_COMMENT = 'add_comment',
  UPDATE_COMMENT = 'update_comment',
  DELETE_COMMENT = 'delete_comment',
}

type AugmentedActionContext = {
  commit<K extends keyof Mutations>(
    key: K,
    payload: Parameters<Mutations[K]>[1]
  ): ReturnType<Mutations[K]>;
} & Omit<ActionContext<NoteState, RootState>, 'commit'>;

export interface Actions {
  [NoteActionTypes.GET_GROUP_COMMENTS](
    context: AugmentedActionContext,
    payload: IncrementalGroupCommentReq
  ): Promise<void>;

  [NoteActionTypes.ADD_COMMENT](
    context: AugmentedActionContext,
    payload: GroupNoteCommentCreateReq
  ): void;

  [NoteActionTypes.UPDATE_COMMENT](
    context: AugmentedActionContext,
    payload: GroupNoteCommentUpdateReq
  ): void;

  [NoteActionTypes.DELETE_COMMENT](
    context: AugmentedActionContext,
    payload: GroupNoteCommentDeleteReq
  ): void;
}

export const actions: ActionTree<NoteState, RootState> & Actions = {
  async [NoteActionTypes.GET_GROUP_COMMENTS](
    { commit, rootState, state },
    params
  ) {
    const comments = await getGroupComments(params);

    if (params.commentModifiedTime === '0') {
      state.comments = {};
    }

    commit(NoteMutationTypes.SET_COMMENTS, comments);
  },

  async [NoteActionTypes.ADD_COMMENT]({ commit, state, rootState }, params) {
    const { commentedUserName, markId, ..._params } = params as any;

    const res = await addComment(_params).catch((err) => {
      const e = err as ResponseError;
      if (e.code === ERROR_CODE_NEED_VIP) {
        useVipStore().showVipLimitDialog(e.message, {
          exception: e.extra as NeedVipException,
          reportParams: {
            element_name: ElementName.upperTeamNumNotePopup,
          },
        });
      }
      throw e;
    });

    commit(NoteMutationTypes.ADD_COMMENT, {
      commentId: res.data.commentId,
      comment: params.commentContent,
      commentatorInfoView: {
        avatarCdnUrl: rootState.user.userInfo?.avatarUrl || '',
        nickName: rootState.user.userInfo?.nickName || '',
        userId: rootState.user.userInfo?.id || '',
      },
      commentedUserName: commentedUserName || '',
      deleteAuthority: true,
      markId: markId || '',
    } as any);
  },

  async [NoteActionTypes.UPDATE_COMMENT]({ commit }, params) {
    const { markId, ..._params } = params as any;

    await updateComment(_params);

    commit(NoteMutationTypes.UPDATE_COMMENT, {
      ..._params,
      markId,
    });
  },

  async [NoteActionTypes.DELETE_COMMENT]({ commit }, params) {
    const { markId, ..._params } = params as any;

    await deleteComment(params);

    commit(NoteMutationTypes.DELETE_COMMENT, {
      ..._params,
      markId,
    });
  },
};
