<script setup lang="ts">
import { computed } from 'vue'
import { NConfigProvider, NGlobalStyle, NMessageProvider, darkTheme } from 'naive-ui'
import { darkThemeOverrides, lightThemeOverrides } from './theme'
import { useUiStore } from './stores/ui'

const ui = useUiStore()

const theme = computed(() => (ui.isDarkTheme ? darkTheme : null))
const themeOverrides = computed(() => (ui.isDarkTheme ? darkThemeOverrides : lightThemeOverrides))
</script>

<template>
  <n-config-provider :theme="theme" :theme-overrides="themeOverrides">
    <n-global-style />
    <n-message-provider>
      <router-view />
    </n-message-provider>
  </n-config-provider>
</template>
