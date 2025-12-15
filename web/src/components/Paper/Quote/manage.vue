<template>
  <div class="manage-citation-style-wrap">
    <ManagePersonal
      ref="managePersonalRef"
      @changeCslList="handleChangeCslList"
    ></ManagePersonal>
    <ManageOther
      ref="manageOtherRef"
      :added-count="addedCount"
      @fetchMyCslList="handleFetchMyCslList"
    ></ManageOther>
  </div>
</template>
<script lang="ts">
import { defineComponent, ref, onMounted, computed } from 'vue'
import ManagePersonal from './manageCitationStyle/personal/index.vue'
import ManageOther from './manageCitationStyle/other/index.vue'

export default defineComponent({
  components: {
    ManagePersonal,
    ManageOther,
  },
  layout: 'none',
  setup() {
    const managePersonalRef = ref()

    const manageOtherRef = ref()

    const handleChangeCslList = (id: string) => {
      manageOtherRef.value?.handleChangeAddedFlag(id)
    }

    const handleFetchMyCslList = () => {
      managePersonalRef.value?.fetch()
    }

    onMounted(() => {
      document.body.style.minWidth = 'unset'
      document.body.style.width = '100%'
      document.body.style.backgroundColor = '#fff'
    })

    const addedCount = computed(
      () => managePersonalRef.value?.styleList?.length
    )
    return {
      manageOtherRef,
      managePersonalRef,
      handleChangeCslList,
      handleFetchMyCslList,
      addedCount,
    }
  },
})
</script>
<style lang="less" scoped>
.manage-citation-style-wrap {
  display: flex;
  height: 312px;
}
</style>
