import { useRequest } from 'ahooks-vue';
import { Ref } from 'vue';
import { getRecentProjects } from '@common/api/revise';

export function useProjects(ready: Ref<boolean>) {
  return useRequest(
    async () => {
      const res = await getRecentProjects();

      return res.simpleProjectInfos;
    },
    {
      ready,
    }
  );
}
