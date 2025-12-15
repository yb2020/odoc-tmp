import {
  LangType,
  TaskStatus,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/review/AiReviewPaper';
import { useI18n } from 'vue-i18n';

export const useAIReviewQSWordings = () => {
  const { t } = useI18n();
  const status2Txt = {
    [TaskStatus.UPLOAD]: t('common.aiReviewerQS.status.uploaded'),
    [TaskStatus.WAITING_PAY]: t('common.aiReviewerQS.status.unpaid'),
    [TaskStatus.WAITING_REVIEW]: t('common.aiReviewerQS.status.waiting'),
    [TaskStatus.REVIEWING]: t('common.aiReviewerQS.status.reviewing'),
    [TaskStatus.SUCCESS]: t('common.aiReviewerQS.status.success'),
    [TaskStatus.FAIL]: t('common.aiReviewerQS.status.failed'),
    [TaskStatus.CONSUME_BEAN_FAIL]: t('common.aiReviewerQS.status.failedBeans'),
    [TaskStatus.CANCEL]: t('common.aiReviewerQS.status.cancelled'),
    [TaskStatus.UNRECOGNIZED]: 'Unknown',
  };

  const lang2Txt = {
    [LangType.ZH]: t('common.text.chinese'),
    [LangType.EN]: t('common.text.english'),
    [LangType.UNRECOGNIZED]: 'Unknown',
  };

  return {
    status2Txt,
    lang2Txt,
  };
};
