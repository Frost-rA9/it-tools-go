<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useThemeVars } from 'naive-ui'
import { useToolsStore } from '../stores/tools'

const route = useRoute()
const store = useToolsStore()
const themeVars = useThemeVars()
const tool = computed(() => store.getById(route.params.id as string))
</script>

<template>
  <div class="tool-layout">
    <div class="tool-header">
      <h1 class="tool-title">{{ tool?.name ?? '加载中…' }}</h1>
      <div class="separator" />
      <p class="tool-description">{{ tool?.description ?? '' }}</p>
    </div>
    <div class="tool-content">
      <router-view />
    </div>
  </div>
</template>

<style scoped>
.tool-layout {
  max-width: 600px;
  margin: 0 auto;
  box-sizing: border-box;
  padding: 0 16px;
}

.tool-header {
  padding: 40px 0;
  width: 100%;
}

.tool-title {
  opacity: 0.9;
  font-size: 40px;
  font-weight: 400;
  margin: 0;
  line-height: 1;
  color: v-bind('themeVars.textColor1');
}

.separator {
  width: 200px;
  height: 2px;
  background: v-bind('themeVars.textColor3');
  opacity: 0.2;
  margin: 10px 0;
}

.tool-description {
  margin: 0;
  opacity: 0.7;
  color: v-bind('themeVars.textColor2');
}

.tool-content {
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: flex-start;
  flex-wrap: wrap;
  gap: 16px;
  padding-bottom: 60px;
}
</style>
