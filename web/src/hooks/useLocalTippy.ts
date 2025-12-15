import { useLocalStorage } from '@vueuse/core'

export interface LocalTippyConfig {
  translateWidth: number,
}

export default function useLocalTippy() {
  const tippyConfig = useLocalStorage<LocalTippyConfig>('pdf-annotate/2.0/tippy', {
    translateWidth: 400,
  })

  const savaLocalTippyConfig = (values: Partial<LocalTippyConfig>) => {
    tippyConfig.value = {
      ...tippyConfig.value,
      ...values,
    }
  }
  return { tippyConfig, savaLocalTippyConfig }
}