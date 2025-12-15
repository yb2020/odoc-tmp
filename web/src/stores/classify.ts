import { UserDocClassifyInfo } from 'go-sea-proto/gen/ts/doc/UserDocManage'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getUserAllClassifyList } from '@/api/material'

export const useClassify = defineStore('classify', () => {
  let inited: Promise<void> | null = null

  const classifyList = ref<UserDocClassifyInfo[]>([])

  const fetchClassifyList = async () => {
    const res = await getUserAllClassifyList({})
    console.log('jinzhi:debug:getUserAllClassifyList', res)
    if (res) {
      classifyList.value = res.map(item => ({
        classifyId: item.classify_id,
        classifyName: item.classify_name,
        docId: item.doc_id,
      }))
    }
  }

  const initClassifyList = async () => {
    if (!inited) {
      inited = fetchClassifyList()
    }

    return inited
  }

  const refreshClassifyList = async () => {
    if (!inited) {
      return initClassifyList()
    }

    return fetchClassifyList()
  }

  return {
    classifyList,
    initClassifyList,
    refreshClassifyList,
  }
})
