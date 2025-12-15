import {
  PaperCommentLevel,
  PaperCommentDimension,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/paper/PaperCommentView';

export const DEFAULT_COMMENT = {
  commentContent: '',
  commentatorInfoView: {
    nickName: '',
    avatarCdnUrl: '',
    userId: 1,
  },
  commentApprovedCount: 0,
  approvedByMe: false,
  paperCommentLevel: -1,
  paperCommentDimensionScore: [
    {
      paperCommentDimension: 0,
      score: -1,
    },
    {
      paperCommentDimension: 1,
      score: -1,
    },
    {
      paperCommentDimension: 2,
      score: -1,
    },
    {
      paperCommentDimension: 3,
      score: -1,
    },
  ],
  anonymous: true,
};

export const COMMENT_OPTIONS = [
  {
    label: '强烈推荐',
    value: PaperCommentLevel.HIGHLY_RECOMMEND,
    i18n: 'tasks.highlyRecommend',
  },
  {
    label: '推荐',
    value: PaperCommentLevel.RECOMMEND,
    i18n: 'tasks.recommend',
  },
  {
    label: '勉强接受',
    value: PaperCommentLevel.ACCEPT,
    i18n: 'tasks.accept',
  },
  {
    label: '拒绝',
    value: PaperCommentLevel.REJECT,
    i18n: 'tasks.reject',
  },
  {
    label: '强烈拒绝',
    value: PaperCommentLevel.STRONG_REJECT,
    i18n: 'tasks.strongReject',
  },
];

export const DIMENSION_NAME_MAP = {
  [PaperCommentDimension.INNOVATION]: 'tasks.innovation',
  [PaperCommentDimension.CONSCIENTIOUSNESS]: 'tasks.strictness',
  [PaperCommentDimension.LOGIC]: 'tasks.logic',
  [PaperCommentDimension.REPRODUCIBILITY]: 'tasks.reproducibility',
};
