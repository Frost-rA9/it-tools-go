<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { NCollapseTransition, NIcon, useThemeVars } from 'naive-ui'
import { ChevronRight } from '@vicons/tabler'
import { useToolsStore } from '../stores/tools'
import { getToolIcon } from '../tools/icons'

const store = useToolsStore()
const router = useRouter()
const themeVars = useThemeVars()

const collapsedNames = ref<string[]>([])

watch(
  () => store.categories,
  (categories) => {
    const names = categories.map((c) => c.name)
    collapsedNames.value = collapsedNames.value.filter((n) => names.includes(n))
  },
  { immediate: true },
)

function isCollapsed(name: string) {
  return collapsedNames.value.includes(name)
}

function toggle(name: string) {
  const i = collapsedNames.value.indexOf(name)
  if (i >= 0) {
    collapsedNames.value.splice(i, 1)
  } else {
    collapsedNames.value.push(name)
  }
}

function go(id: string) {
  router.push({ name: 'tool', params: { id } })
}
</script>

<template>
  <div class="tool-menu">
    <div v-for="cat in store.categories" :key="cat.name" class="menu-category">
      <div class="category-header" @click="toggle(cat.name)">
        <span class="chevron" :class="{ expanded: !isCollapsed(cat.name) }">
          <n-icon :component="ChevronRight" size="16" />
        </span>
        <span class="category-name">{{ cat.name }}</span>
      </div>

      <n-collapse-transition :show="!isCollapsed(cat.name)">
        <div class="menu-wrapper">
          <div class="toggle-bar" @click="toggle(cat.name)" />
          <div class="menu-items">
            <div
              v-for="tool in cat.tools"
              :key="tool.id"
              class="tool-item"
              @click="go(tool.id)"
            >
              <n-icon :component="getToolIcon(tool.icon)" size="18" class="tool-item-icon" />
              <span>{{ tool.name }}</span>
            </div>
          </div>
        </div>
      </n-collapse-transition>
    </div>
  </div>
</template>

<style scoped>
.tool-menu {
  padding: 0 12px;
}

.menu-category {
  margin-bottom: 5px;
}

.category-header {
  display: flex;
  align-items: center;
  margin: 12px 0 0 6px;
  cursor: pointer;
  opacity: 0.6;
}

.chevron {
  display: inline-flex;
  align-items: center;
  transition: transform 0.2s;
}

.chevron.expanded {
  transform: rotate(90deg);
}

.category-name {
  margin-left: 8px;
  font-size: 13px;
}

.menu-wrapper {
  display: flex;
  flex-direction: row;
}

.menu-items {
  flex: 1;
  margin-bottom: 5px;
}

.toggle-bar {
  width: 24px;
  position: relative;
  cursor: pointer;
  opacity: 0.1;
  transition: opacity ease 0.2s;
}

.toggle-bar::before {
  content: '';
  width: 2px;
  height: 100%;
  border-radius: 2px;
  background-color: v-bind('themeVars.textColor3');
  position: absolute;
  top: 0;
  left: 14px;
}

.toggle-bar:hover {
  opacity: 0.5;
}

.tool-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  font-size: 14px;
  opacity: 0.85;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.15s;
  color: v-bind('themeVars.textColor2');
}

.tool-item-icon {
  opacity: 0.7;
}

.tool-item:hover {
  background-color: v-bind('themeVars.hoverColor');
  opacity: 1;
}
</style>
