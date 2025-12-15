import { GroupComment } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/GroupNote';
/*
 * Created Date: March 17th 2022, 12:13:20 pm
 * Author: zhoupengcheng
 * -----
 * Last Modified: May 24th 2022, 3:21:28 pm
 */
export interface NoteState {
  commentModifiedTime: string;
  comments: Record<string, Array<GroupComment>>;
}

export { NoteSubTypes } from '@common/components/Notes/types';
