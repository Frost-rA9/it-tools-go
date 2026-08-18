<script setup lang="ts">
import { ref, watch } from 'vue'
import { NSelect, NSwitch, NCard, NAlert, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import CodeOutput from '../../components/CodeOutput.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const themeVars = useThemeVars()

const jsonInput = ref('[{"name":"Alice","age":30,"city":"Beijing"},{"name":"Bob","age":25,"city":"Shanghai"}]')
const delimiter = ref(',')
const includeHeader = ref(true)
const csvOutput = ref('')
const rows = ref(0)
const columns = ref(0)
const errorMessage = ref('')

const delimiterOptions = [
  { label: '逗号 (,)', value: ',' },
  { label: '分号 (;)', value: ';' },
  { label: 'Tab', value: '\t' },
]

async function run() {
  try {
    const output = await RunTool(
      'json-to-csv',
      JSON.stringify({ json: jsonInput.value, delimiter: delimiter.value, include_header: includeHeader.value }),
    )
    const result = JSON.parse(output)
    csvOutput.value = result.csv
    rows.value = result.rows
    columns.value = result.columns
    errorMessage.value = ''
  } catch (e) {
    csvOutput.value = ''
    rows.value = 0
    columns.value = 0
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 200)
watch([jsonInput, delimiter, includeHeader], () => debouncedRun(), { immediate: true })
</script>

<template>
  <n-card title="JSON 转 CSV — 输入" class="tool-card">
    <ToolTextarea v-model:value="jsonInput" label="输入 JSON（对象数组）" :rows="20" placeholder='[{"name":"Alice"}, ...]' monospace />

    <div class="options-row">
      <span class="option-label">分隔符</span>
      <n-select v-model:value="delimiter" :options="delimiterOptions" class="delimiter-select" />
      <span class="option-label">表头</span>
      <n-switch v-model:value="includeHeader" />
    </div>

    <n-alert v-if="errorMessage" type="error" class="error-alert">
      {{ errorMessage }}
    </n-alert>
  </n-card>

  <n-card title="JSON 转 CSV — 输出" class="tool-card">
    <CodeOutput label="CSV 结果" :value="csvOutput" language="plaintext" :rows="20" />

    <div class="stats-row">
      <span class="stat">{{ rows }} 行</span>
      <span class="stat">{{ columns }} 列</span>
    </div>
  </n-card>
</template>

<style scoped>
.options-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.option-label {
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
}

.delimiter-select {
  width: 140px;
}

.error-alert {
  margin-bottom: 16px;
}

.stats-row {
  display: flex;
  gap: 16px;
  justify-content: center;
  margin-bottom: 16px;
}

.stat {
  font-size: 13px;
  color: v-bind('themeVars.textColor3');
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}
</style>
