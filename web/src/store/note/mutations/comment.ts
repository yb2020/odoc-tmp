import { MutationTree } from 'vuex';
import { NoteState } from '../types';
import {
  IncrementalGroupCommentRsp,
  OperationType,
  GroupComment,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/GroupNote';

export enum NoteMutationTypes {
  SET_COMMENTS = 'set_comments',
  ADD_COMMENT = 'add_comment',
  UPDATE_COMMENT = 'update_comment',
  DELETE_COMMENT = 'delete_comment',
  SET_COMMENT_MODIFIED_TIME = 'set_comment_modified_time',
}

export type Mutations<S = NoteState> = {
  [NoteMutationTypes.SET_COMMENTS](state: S, data: IncrementalGroupCommentRsp): void;
  [NoteMutationTypes.ADD_COMMENT](state: S, data: GroupComment & { markId: string }): void;
  [NoteMutationTypes.SET_COMMENT_MODIFIED_TIME](state: S, data: string): void;
  [NoteMutationTypes.UPDATE_COMMENT](
    state: S,
    data: { markId: string; commentId: string; commentContent: string }
  ): void;
  [NoteMutationTypes.DELETE_COMMENT](state: S, data: { markId: string; commentId: string }): void;
};

export const mutations: MutationTree<NoteState> & Mutations = {
  [NoteMutationTypes.SET_COMMENTS](state, data) {
    const { commentModifiedTime, incrementalGroupComments } = data;

    state.commentModifiedTime = commentModifiedTime;

    incrementalGroupComments.forEach((item) => {
      if (item.operationType === OperationType.Create) {
        const comments = state.comments[item.markId] || [];

        for (const com of item.groupComments) {
          const index = comments.findIndex((c) => c.commentId === com.commentId);
          if (index >= 0) {
            state.comments[item.markId][index] = com;
            continue;
          }

          state.comments[item.markId] = [...(state.comments[item.markId] || []), com];
        }

        return;
      }

      if (item.operationType === OperationType.Delete) {
        item.groupComments.forEach((i) => {
          const index = state.comments[item.markId].findIndex((l) => l.commentId === i.commentId);

          if (index >= 0) {
            state.comments[item.markId].splice(index, 1);
          }
        });

        return;
      }

      if (item.operationType === OperationType.Update) {
        item.groupComments.forEach((i) => {
          const index = (state.comments[item.markId] || []).findIndex(
            (l) => l.commentId === i.commentId
          );

          if (index >= 0) {
            state.comments[item.markId][index] = i;
          }
        });
        return;
      }
    });
  },

  [NoteMutationTypes.SET_COMMENT_MODIFIED_TIME](state, data) {
    state.commentModifiedTime = data;
  },

  [NoteMutationTypes.ADD_COMMENT](state, data) {
    state.comments[data.markId] = [...(state.comments[data.markId] || []), data];
  },

  [NoteMutationTypes.UPDATE_COMMENT](state, data) {
    const index = state.comments[data.markId].findIndex((com) => com.commentId === data.commentId);

    state.comments[data.markId][index] = {
      ...state.comments[data.markId][index],
      comment: data.commentContent,
    };
  },

  [NoteMutationTypes.DELETE_COMMENT](state, data) {
    const index = state.comments[data.markId].findIndex((com) => com.commentId === data.commentId);

    state.comments[data.markId].splice(index, 1);
  },
};
