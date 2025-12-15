import {
  LSKEY_CURRENT_TRANSLATE_TAB,
  TranslateTabKey,
} from '../stores/translateStore';

const LSKEY_CURRENT_TRANSLATE_TAB_RESET_DURATION_TIME =
  'pdf-annotate/2.0/translateTabResetDurationTime';

const duration = 30 * 60 * 1000;

export const resetLSCurrentTranslateTabKeyByDuration = () => {
  try {
    const currentTab = localStorage.getItem(LSKEY_CURRENT_TRANSLATE_TAB);
    /**
     * 1、凡是切换为自定义翻译引擎的，记录原本的选项，不再调整回idea
     * 2、切换为有道/百度等的，依然默认是idea
     */

    if (
      currentTab !== TranslateTabKey.youdao &&
      currentTab !== TranslateTabKey.baidu
    ) {
      return;
    }
    const now = new Date();

    const currentTabResetTime = parseInt(
      localStorage.getItem(LSKEY_CURRENT_TRANSLATE_TAB_RESET_DURATION_TIME) ||
        now.getTime().toString(),
      10
    );
    const lastTime = new Date(currentTabResetTime);

    console.log(
      'resetLSCurrentTranslateTabKeyByDuration',
      now.getTime() - lastTime.getTime(),
      duration
    );

    if (now.getTime() - lastTime.getTime() > duration) {
      localStorage.setItem(LSKEY_CURRENT_TRANSLATE_TAB, TranslateTabKey.idea);
    }
    localStorage.setItem(
      LSKEY_CURRENT_TRANSLATE_TAB_RESET_DURATION_TIME,
      now.getTime().toString()
    );
  } catch (error) {
    console.error(error);
  }
};
