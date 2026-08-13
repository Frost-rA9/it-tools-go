import { defineStore } from 'pinia'
import { useDark, useMediaQuery, useStorage, useToggle } from '@vueuse/core'

export const useUiStore = defineStore('ui', () => {
  const isDarkTheme = useDark({ initialValue: 'light' })
  const toggleDark = useToggle(isDarkTheme)
  const isSmallScreen = useMediaQuery('(max-width: 700px)')
  const isMenuCollapsed = useStorage('isMenuCollapsed', false)

  return { isDarkTheme, toggleDark, isSmallScreen, isMenuCollapsed }
})
