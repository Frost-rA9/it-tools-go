<script setup lang="ts">
import { useThemeVars } from 'naive-ui'
import { useToolsStore } from '../stores/tools'
import ToolCard from '../components/ToolCard.vue'

const store = useToolsStore()
const themeVars = useThemeVars()
</script>

<template>
  <div class="home">
    <template v-for="cat in store.categories" :key="cat.name">
      <h3 class="category-title">{{ cat.name }}</h3>
      <div class="grid">
        <ToolCard v-for="tool in cat.tools" :key="tool.id" :tool="tool" />
      </div>
    </template>
  </div>
</template>

<style scoped>
.home {
  padding: 50px 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.category-title {
  font-size: 15px;
  font-weight: 500;
  color: v-bind('themeVars.textColor3');
  margin: 25px 0 5px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 12px;
}
</style>
