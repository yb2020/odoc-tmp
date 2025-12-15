import { computed, ref } from 'vue';
import useFetch from './useFetch';

export type QueryParam = {
  currentPage: number;
  pageSize: number;
};

type ListFetchFnType<T> = (params: QueryParam) => Promise<{ list: Array<T>; total: number }>;

export default function useFetchList<T>(
  fetchFn: ListFetchFnType<T>,
  currentPageSize?: number,
  immediate?: boolean
) {
  const list = ref<T[]>([]);
  const currentPageList = ref<T[]>([]);
  const currentPage = ref<number>(0);
  const total = ref<number>(-1);
  const pageSize = currentPageSize || 10;
  const hasmore = ref<boolean>(true);
  let queryPage = 1;
  const { fetchState, fetch } = useFetch(async () => {
    const res = await fetchFn({
      currentPage: queryPage,
      pageSize,
    });

    if (res.total == -1) {
      // 说明本次请求无效，会重试
      return;
    }

    if (queryPage === 1) {
      list.value = [];
      total.value = res.total;
    }

    currentPageList.value = (res.list || []) as any;
    list.value = list.value.concat((res.list || []) as any);
    hasmore.value = total.value === -1 || queryPage * pageSize < total.value;

    currentPage.value = queryPage;

    queryPage += 1;
  }, immediate);

  const loading = computed<boolean>({
    get: () => fetchState.pending,
    set: (val: boolean) => (loading.value = val),
  });

  const error = computed<string>({
    get: () => fetchState.error?.message || '',
    set: (val: string) => (error.value = val),
  });

  const refresh = async () => {
    queryPage = 1;
    await fetch();
  };

  const fetchData = async (page?: number) => {
    if (loading.value) {
      return;
    }
    queryPage = page || currentPage.value + 1;
    await fetch();
  };

  return {
    total,
    loading,
    list,
    error,
    fetchState,
    fetch: fetchData,
    hasmore,
    refresh,
    currentPage,
    currentPageList,
  };
}
