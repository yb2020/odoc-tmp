import { RootState } from '../types';
import {
  Module,
  Store as VuexStore,
  CommitOptions,
  DispatchOptions,
} from 'vuex';
import { Mutations, mutations } from './mutations';
import { Actions, actions } from './actions';
import { NoteState } from './types';

export { NoteMutationTypes } from './mutations';
export { NoteActionTypes } from './actions';

export type NoteStore<S = NoteState> = Omit<
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

export const NoteModule: Module<NoteState, RootState> = {
  namespaced: true,
  state: () => ({
    comments: {},
    commentModifiedTime: '0',
  }),
  actions,
  mutations,
};
