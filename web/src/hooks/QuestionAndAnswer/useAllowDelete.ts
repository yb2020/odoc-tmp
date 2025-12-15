import { computed, ref } from "vue";
import { useStore } from "~/src/store";
import { getAllowedDeleteQuestionAndAnswerUserIdList } from '~/src/api/question'

const useAllowDelete = () => {
  const store = useStore();

  const user = computed(() => store.state.user);

  const AllowedDeleteUserIdList = ref<string[]>([]);

  const isAllowedDelete = computed(() =>
    AllowedDeleteUserIdList.value.includes(user.value?.userInfo?.id || '')
  );

  const getAllowedDeleteUserIdList = async () => {
    const res = await getAllowedDeleteQuestionAndAnswerUserIdList();
    AllowedDeleteUserIdList.value = res;
  };

  getAllowedDeleteUserIdList();

  return { isAllowedDelete }
}

export default useAllowDelete