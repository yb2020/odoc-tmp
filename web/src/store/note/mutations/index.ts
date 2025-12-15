/*
 * Created Date: May 26th 2022, 3:56:58 pm
 * Author: zhoupengcheng
 * -----
 * Last Modified: May 26th 2022, 3:56:58 pm
 */
import {
  mutations as commentMutations,
  Mutations as NoteMutation,
  NoteMutationTypes as _NoteMutationTypes,
} from './note';
import {
  mutations as noteMutations,
  Mutations as CommentMutations,
  NoteMutationTypes as CommentMutationTypes,
} from './comment';

export const mutations = {
  ...commentMutations,
  ...noteMutations,
};

export type Mutations = CommentMutations & NoteMutation;

export const NoteMutationTypes = {
  ..._NoteMutationTypes,
  ...CommentMutationTypes,
};

export type NoteMutationTypes = typeof NoteMutationTypes;
