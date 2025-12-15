import { NoteLocationInfo } from '~/src/api/setting';
import { LeftSideBarType } from '~/src/components/Left/type';
import { RightSideBarType } from '~/src/components/Right/TabPanel/type';
import { LangType } from '~/src/stores/copilotType';

export const defaultCommonSettings: Required<NoteLocationInfo> = {
  currentGroupId: '',
  scale: 1.0,
  scalePresetValue: '',
  rightShow: true,
  rightTab: RightSideBarType.Note,
  rightSubTab: '',
  rightWidth: 280,
  sideBarShow: true,
  sideBarTab: LeftSideBarType.Catalog,
  sideBarWidth: 160,
  currentPage: 1,
  rightTabBars: [
    RightSideBarType.Matirial,
    RightSideBarType.Note,
    RightSideBarType.Group,
    RightSideBarType.Learn,
    RightSideBarType.Copilot,
  ].map((key) => {
    return {
      key,
      shown: true,
    };
  }),
  copilotLanguage: LangType.CHINESE,
  toolBarVisible: true,
  toolBarHeadVisible: true,
  toolBarNoteVisible: true,
};
