<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NIcon, NTooltip, useThemeVars } from 'naive-ui'
import { Menu2, Home2, Search, BrandGithub, Sun, Moon } from '@vicons/tabler'
import { useUiStore } from '../stores/ui'
import CommandPalette from './CommandPalette.vue'

const router = useRouter()
const ui = useUiStore()
const themeVars = useThemeVars()
const palette = ref<InstanceType<typeof CommandPalette> | null>(null)

const GITHUB_URL = 'https://github.com/Frost-rA9/it-tools-go'

const paletteBg = computed(() => (ui.isDarkTheme ? 'rgba(255, 255, 255, 0.08)' : '#dfe4e8'))

function goHome() {
  router.push('/')
}
</script>

<template>
  <div class="toolbar">
    <n-tooltip placement="bottom">
      <template #trigger>
        <n-button quaternary circle @click="ui.isMenuCollapsed = !ui.isMenuCollapsed">
          <n-icon :size="25" :component="Menu2" />
        </n-button>
      </template>
      切换侧边栏
    </n-tooltip>

    <n-tooltip placement="bottom">
      <template #trigger>
        <n-button quaternary circle @click="goHome">
          <n-icon :size="25" :component="Home2" />
        </n-button>
      </template>
      主页
    </n-tooltip>

    <div class="palette-trigger" @click="palette?.open()">
      <span class="palette-trigger-inner">
        <n-icon :size="16" :component="Search" />
        <span class="palette-label">搜索</span>
        <span class="palette-kbd">Ctrl K</span>
      </span>
    </div>

    <n-tooltip placement="bottom">
      <template #trigger>
        <n-button quaternary circle tag="a" :href="GITHUB_URL" target="_blank">
          <n-icon :size="25" :component="BrandGithub" />
        </n-button>
      </template>
      GitHub 仓库
    </n-tooltip>

    <n-tooltip placement="bottom">
      <template #trigger>
        <n-button quaternary circle @click="ui.toggleDark()">
          <n-icon :size="25" :component="ui.isDarkTheme ? Sun : Moon" />
        </n-button>
      </template>
      {{ ui.isDarkTheme ? '切换亮色' : '切换暗色' }}
    </n-tooltip>

    <CommandPalette ref="palette" />
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 16px;
}

.palette-trigger {
  flex: 1;
  margin-left: 8px;
  margin-right: 8px;
  padding: 6px 12px;
  border-radius: 4px;
  background-color: v-bind('paletteBg');
  border: 1px solid transparent;
  cursor: pointer;
  transition: border-color 0.2s;
}

.palette-trigger:hover {
  border-color: v-bind('themeVars.primaryColor');
}

.palette-trigger-inner {
  display: flex;
  align-items: center;
  gap: 8px;
  opacity: 0.6;
}

.palette-label {
  font-size: 14px;
}

.palette-kbd {
  font-size: 12px;
  padding: 2px 6px;
  border: 1px solid currentColor;
  border-radius: 4px;
  opacity: 0.5;
}
</style>
