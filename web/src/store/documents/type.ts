import { ToolBarType } from '@idea/pdf-annotate-core';
import { NoteLocationInfo } from '~/src/api/setting';
import { Nullable } from '~/src/typings/global';
import { SideTabCommonSettings } from '~/src/hooks/UserSettings/useSideTabSettings';
import { RightSideBarType } from '~/src/components/Right/TabPanel/type';
export interface DocumentsState {
  isFullPage: boolean;
  toolBarType: ToolBarType;
  userSettingInfo: Required<NoteLocationInfo>;
  rightSideTabSettings: Nullable<SideTabCommonSettings>;
  disabledRightTabs: RightSideBarType[];
}
