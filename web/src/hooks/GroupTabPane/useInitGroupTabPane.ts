import { ref } from "vue"
import { $GroupProceed, getMyCreatedGroupList, getMyJoinGroupList } from "~/src/api/group"
import { useStore } from "~/src/store"
import { BaseMutationTypes } from "~/src/store/base"

export const useInitGroupTabPane = () => {
  const store = useStore()
  const createdGroupList = ref<$GroupProceed[]>()
  const joinedGroupList = ref<$GroupProceed[]>()
  const getGroupList = () => {
    
  }
  // const getGroupList = async () => {
  //   const [res0, res1] = await Promise.all([
  //     getMyCreatedGroupList({ currentPage: 1, pageSize: 50 }),
  //     getMyJoinGroupList({ currentPage: 1, pageSize: 50 })
  //   ])

  //   createdGroupList.value = res0.list || []
  //   joinedGroupList.value = res1.list || []

  //   const groupList = [...createdGroupList.value, ...joinedGroupList.value]

  //   if (groupList.length) {
  //     store.commit(`base/${BaseMutationTypes.SET_GROUP_INFOS}`,  groupList)
  //     return groupList
  //   }
  //   return null;
  // }
  return {
    createdGroupList,
    joinedGroupList,
    getGroupList,
  }
}

