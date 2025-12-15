import { BaseState } from './base/type';
import { DocumentsState } from './documents/type';
import { ParseState } from './parse/type';
import { UserState } from './user/type';
import { NoteState } from './note/types';
import { ShortcutsState } from './shortcuts/type';
import { CertState } from './cert';

export interface RootState {
  user: UserState;
  base: BaseState;
  documents: DocumentsState;
  parse: ParseState;
  note: NoteState;
  shortcuts: ShortcutsState;
  cert: CertState;
}
