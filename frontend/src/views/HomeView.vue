<script setup lang="ts">
import { useThemeVars } from 'naive-ui'
import { useToolsStore } from '../stores/tools'
import ToolCard from '../components/ToolCard.vue'

const store = useToolsStore()
const themeVars = useThemeVars()
</script>

<template>
  <div class="home">
    <h3 class="section-title">全部工具</h3>
    <div class="grid">
      <ToolCard v-for="tool in store.tools" :key="tool.id" :tool="tool" />
    </div>
  </div>
</template>

<style scoped>
.home {
  padding: 50px 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.section-title {
  font-size: 15px;
  font-weight: 500;
  color: v-bind('themeVars.textColor3');
  margin: 0 0 5px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  grid-auto-rows: 180px;
  gap: 12px;
}

@media (min-width: 640px) {
  .grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 768px) {
  .grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (min-width: 1280px) {
  .grid {
    grid-template-columns: repeat(4, 1fr);
  }
}
</style>
