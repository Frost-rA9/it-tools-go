<script setup lang="ts">
import { computed, defineAsyncComponent } from 'vue'
import { NEmpty } from 'naive-ui'
import { useToolComponent } from '../composables/useToolComponent'

const { tool, component } = useToolComponent()

const toolComponent = computed(() => {
  if (!component.value) return null
  return defineAsyncComponent(component.value as any)
})
</script>

<template>
  <component v-if="tool && toolComponent" :is="toolComponent" />
  <n-empty v-else-if="tool" description="该工具界面尚未实现" class="tool-empty" />
  <n-empty v-else description="未找到工具" class="tool-empty" />
</template>

<style scoped>
.tool-empty {
  margin-top: 40px;
}
</style>
