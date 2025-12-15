import { CpPaperInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/search/AutoCP'
import { ref, unref } from 'vue'
import _ from 'lodash'
import { fetchRecommended } from '@/common/src/api/search'

export default function useSearchRecommend(defaultKeyword: string) {
  const keyword = ref<string>(defaultKeyword || '')
  const items = ref<CpPaperInfo[]>([])
  let lastFetchId = 0 // 请求时序控制
  const fetchQuery = _.debounce(async (callback?: () => void) => {
    if (!String(unref(keyword)).trim()) {
      items.value = []
      lastFetchId = 0
      return
    }
    lastFetchId += 1
    const fetchId = lastFetchId
    const res = await fetchRecommended({
      query: unref(keyword),
      cpType: 0,
    })
    if (fetchId !== lastFetchId) {
      return
    }
    items.value = res
    if (typeof callback === 'function') {
      callback()
    }
  }, 100)

  const clear = () => {
    lastFetchId = 0
    items.value = []
  }

  return {
    keyword,
    items,
    fetchQuery,
    clear,
  }
}
