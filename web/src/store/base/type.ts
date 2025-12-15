import { NoteBaseInfo } from '~/src/api/base';
import { $GroupProceed } from '~/src/api/group';
// import { AnnotateTag } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { AnnotateTag } from 'go-sea-proto/gen/ts/common/AnnotateTag';
import { GetPdfStatusInfoResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/OpenNote';

export interface BaseQueryParam {
  paperId: string;
}

export const SELF_NOTEINFO_GROUPID = '0';

export interface GroupNoteBaseInfo {
  groupId: string;
  noteInfo: NoteBaseInfo;
}

export interface BaseState extends BaseQueryParam {
  statusInfo: GetPdfStatusInfoResponse;
  currentGroupId: string;
  noteInfoMap: Record<string, NoteBaseInfo>; // key是groupId
  groupInfoMap: Record<string, $GroupProceed>; // key是groupId
  tagList: Array<AnnotateTag>;
}
