<script setup lang="ts">
import { useMessage, useThemeVars } from 'naive-ui'
import { useClipboard } from '@vueuse/core'

export interface EmojiInfo {
  emoji: string
  name: string
  group: string
  keywords: string[]
  codepoints: string
  unicode: string
}

defineProps<{ info: EmojiInfo }>()

const themeVars = useThemeVars()
const message = useMessage()
const { copy } = useClipboard()

async function copyText(text: string, tip: string) {
  if (!text) return
  await copy(text)
  message.success(`${tip} 已复制到剪贴板`)
}
</script>

<template>
  <div class="emoji-card">
    <span class="emoji-char" title="点击复制 Emoji" @click="copyText(info.emoji, 'Emoji')">{{ info.emoji }}</span>
    <div class="emoji-meta">
      <div class="name">{{ info.name }}</div>
      <div class="mono-row">
        <span class="mono-item" title="点击复制码点" @click="copyText(info.codepoints, 'Codepoints')">{{ info.codepoints }}</span>
        <span class="mono-item" title="点击复制 Unicode 转义" @click="copyText(info.unicode, 'Unicode')">{{ info.unicode }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.emoji-card {
  display: flex;
  align-items: center;
  gap: 12px;
  height: 68px;
  padding: 0 12px;
  box-sizing: border-box;
  border-bottom: 1px solid v-bind('themeVars.borderColor');
  background: v-bind('themeVars.cardColor');
  cursor: default;
}

.emoji-char {
  font-size: 26px;
  line-height: 1;
  cursor: pointer;
  user-select: none;
  flex-shrink: 0;
}

.emoji-meta {
  flex: 1;
  min-width: 0;
}

.name {
  font-size: 14px;
  font-weight: 600;
  line-height: 1.4;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: v-bind('themeVars.textColor1');
}

.mono-row {
  display: flex;
  gap: 8px;
  font-size: 12px;
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
  font-variant-ligatures: none;
  font-feature-settings: 'liga' 0, 'calt' 0;
  opacity: 0.7;
  color: v-bind('themeVars.textColor2');
}

.mono-item {
  padding: 0 6px;
  border-radius: 4px;
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0.2s;
}

.mono-item:hover {
  color: v-bind('themeVars.primaryColor');
}
</style>