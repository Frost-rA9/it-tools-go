<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useThemeVars } from 'naive-ui'
import type { Tool } from '../stores/tools'

const props = defineProps<{ tool: Tool }>()
const router = useRouter()
const themeVars = useThemeVars()

function go() {
  router.push({ name: 'tool', params: { id: props.tool.id } })
}
</script>

<template>
  <div class="tool-card" @click="go">
    <div class="tool-card-name">{{ tool.name }}</div>
    <div class="tool-card-desc">{{ tool.description }}</div>
  </div>
</template>

<style scoped>
.tool-card {
  background: v-bind('themeVars.cardColor');
  border: 2px solid v-bind('themeVars.borderColor');
  border-radius: 4px;
  padding: 16px;
  cursor: pointer;
  height: 100%;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.tool-card:hover {
  border-color: v-bind('themeVars.primaryColor');
}

.tool-card-name {
  font-size: 18px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor1');
}

.tool-card-desc {
  font-size: 14px;
  color: v-bind('themeVars.textColor3');
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
