<script setup lang="ts">
import { useRouter } from 'vue-router'
import { NIcon, useThemeVars } from 'naive-ui'
import type { Tool } from '../stores/tools'
import { getToolIcon } from '../tools/icons'

const props = defineProps<{ tool: Tool }>()
const router = useRouter()
const themeVars = useThemeVars()

function go() {
  router.push({ name: 'tool', params: { id: props.tool.id } })
}
</script>

<template>
  <div class="tool-card-item" @click="go">
    <n-icon :component="getToolIcon(tool.icon)" size="40" class="tool-card-icon" />
    <div class="tool-card-name">{{ tool.name }}</div>
    <div class="tool-card-desc">{{ tool.description }}</div>
  </div>
</template>

<style scoped>
.tool-card-item {
  display: flex;
  flex-direction: column;
  background: v-bind('themeVars.cardColor');
  border: 2px solid v-bind('themeVars.borderColor');
  border-radius: 4px;
  padding: 16px;
  cursor: pointer;
  height: 100%;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.tool-card-item:hover {
  border-color: v-bind('themeVars.primaryColor');
}

.tool-card-item-icon {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  color: v-bind('themeVars.textColor3');
}

.tool-card-item-name {
  font-size: 18px;
  line-height: 24px; /* 显式行高（18px 字体），消除中英文字体默认行高差异 */
  margin: 5px 0 6px;
  color: v-bind('themeVars.textColor1');
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tool-card-item-desc {
  font-size: 14px;
  line-height: 1.5;
  color: v-bind('themeVars.textColor3');
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  height: 42px; /* 固定两行高度（14px × 1.5 × 2），保证所有卡片等高 */
}
</style>
