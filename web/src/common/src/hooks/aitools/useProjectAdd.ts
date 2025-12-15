import { useRequest } from 'ahooks-vue';
import { useRouter } from 'vue-router';
import { ref } from 'vue';
import { createProject, saveVersionData } from '@common/api/revise';
import { openProjectPage } from '@common/components/AIRevise/utils';
import { HEADER_CANCLE_AUTO_ERROR } from '../../api/const';

export function useProjectAdd(
  isOpenNewTab = ref(false),
  isReplace = ref(false)
) {
  const router = useRouter();
  return useRequest(
    async () => {
      const res = await createProject(undefined, {
        headers: {
          [HEADER_CANCLE_AUTO_ERROR]: 'true',
        },
      });
      await saveVersionData({
        projectId: res.projectId,
        text: '',
      });

      if (isOpenNewTab.value) {
        openProjectPage(router, { id: res.projectId }, 'resolve');
      } else {
        openProjectPage(
          router,
          { id: res.projectId },
          isReplace.value ? 'replace' : 'push'
        );
      }

      return res;
    },
    {
      manual: true,
    }
  );
}
