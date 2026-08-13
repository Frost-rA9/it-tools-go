<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { NCollapse, NCollapseItem, useThemeVars } from 'naive-ui'
import { useToolsStore } from '../stores/tools'

const store = useToolsStore()
const router = useRouter()
const themeVars = useThemeVars()

const expandedNames = ref<string[]>([])

watch(
  () => store.categories,
  (categories) => {
    expandedNames.value = categories.map((c) => c.name)
  },
  { immediate: true },
)

function go(id: string) {
  router.push({ name: 'tool', params: { id } })
}
</script>

<template>
  <n-collapse v-model:expanded-names="expandedNames" class="tool-menu">
    <n-collapse-item
      v-for="cat in store.categories"
      :key="cat.name"
      :title="cat.name"
      :name="cat.name"
    >
      <div class="menu-wrapper">
        <div
          v-for="tool in cat.tools"
          :key="tool.id"
          class="tool-item"
          @click="go(tool.id)"
        >
          {{ tool.name }}
        </div>
      </div>
    </n-collapse-item>
  </n-collapse>
</template>

<style scoped>
.tool-menu {
  padding: 0 12px;
}

.menu-wrapper {
  position: relative;
  padding-left: 12px;
}

.menu-wrapper::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
  border-radius: 2px;
  background-color: v-bind('themeVars.textColor3');
  opacity: 0.2;
  transition: opacity 0.2s;
}

.menu-wrapper:hover::before {
  opacity: 0.5;
}

.tool-item {
  padding: 6px 12px;
  font-size: 14px;
  opacity: 0.85;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.15s;
  color: v-bind('themeVars.textColor2');
}

.tool-item:hover {
  background-color: v-bind('themeVars.hoverColor');
  opacity: 1;
}
</style>
