/*
 * Created Date: May 26th 2022, 3:44:04 pm
 * Author: zhoupengcheng
 * -----
 * Last Modified: May 26th 2022, 3:44:04 pm
 */
import {
  actions as commentActions,
  Actions as CommentActions,
  NoteActionTypes as CommentActionTypes,
} from './comment';

export const actions = commentActions;

export type Actions = CommentActions;

export const NoteActionTypes = CommentActionTypes;

export type NoteActionTypes = typeof NoteActionTypes;
