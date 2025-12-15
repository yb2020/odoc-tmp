import { ref, onMounted } from 'vue'
import { defineStore } from 'pinia'

import { useVueUpload } from '@/utils/pdf-upload/index.js'
import api from '@/common/src/api/axios'

export const useUpload = defineStore('upload', () => {
  const upload = useVueUpload({
    ref,
    onMounted: onMounted,
  })

  return upload
})
