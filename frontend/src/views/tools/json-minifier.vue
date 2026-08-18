<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NAlert, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import CodeOutput from '../../components/CodeOutput.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const themeVars = useThemeVars()

const jsonInput = ref('{\n  "name": "it-tools",\n  "version": 2,\n  "active": true\n}')
const minified = ref('')
const originalSize = ref(0)
const minifiedSize = ref(0)
const saved = ref(0)
const savedPercent = ref(0)
const errorMessage = ref('')

async function run() {
  try {
    const output = await RunTool('json-minifier', JSON.stringify({ json: jsonInput.value }))
    const result = JSON.parse(output)
    minified.value = result.minified
    originalSize.value = result.original_size
    minifiedSize.value = result.minified_size
    saved.value = result.saved
    savedPercent.value = result.saved_percent
    errorMessage.value = ''
  } catch (e) {
    minified.value = ''
    originalSize.value = 0
    minifiedSize.value = 0
    saved.value = 0
    savedPercent.value = 0
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 200)
watch(jsonInput, () => debouncedRun(), { immediate: true })
</script>

<template>
  <n-card title="JSON 压缩 — 输入" class="tool-card">
    <ToolTextarea v-model:value="jsonInput" label="输入 JSON" :rows="20" placeholder="在此粘贴 JSON…" monospace />

    <n-alert v-if="errorMessage" type="error" class="error-alert">
      {{ errorMessage }}
    </n-alert>
  </n-card>

  <n-card title="JSON 压缩 — 输出" class="tool-card">
    <CodeOutput label="压缩结果" :value="minified" language="json" :rows="20" />

    <div class="stats-row">
      <span class="stat">原始 {{ originalSize }} B</span>
      <span class="stat">压缩 {{ minifiedSize }} B</span>
      <span class="stat">节省 {{ saved }} B</span>
      <span class="stat">{{ savedPercent.toFixed(1) }}%</span>
    </div>
  </n-card>
</template>

<style scoped>
.error-alert {
  margin-bottom: 16px;
}

.stats-row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
  margin-bottom: 16px;
}

.stat {
  font-size: 13px;
  color: v-bind('themeVars.textColor3');
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}
</style>
