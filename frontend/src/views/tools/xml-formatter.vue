<script setup lang="ts">
import { ref, watch } from 'vue'
import { NSelect, NAlert, NCard, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import CodeOutput from '../../components/CodeOutput.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const themeVars = useThemeVars()

const xmlInput = ref('<root><user id="1"><name>Alice</name><tags><tag>a</tag><tag>b</tag></tags></user><empty/></root>')
const indent = ref('2')
const formatted = ref('')
const lineCount = ref(0)
const errorMessage = ref('')

const indentOptions = [
  { label: '2 空格', value: '2' },
  { label: '4 空格', value: '4' },
  { label: 'Tab', value: '\t' },
]

async function run() {
  try {
    const output = await RunTool('xml-formatter', JSON.stringify({ xml: xmlInput.value, indent: indent.value }))
    const result = JSON.parse(output)
    formatted.value = result.formatted
    lineCount.value = result.line_count
    errorMessage.value = ''
  } catch (e) {
    formatted.value = ''
    lineCount.value = 0
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 200)
watch([xmlInput, indent], () => debouncedRun(), { immediate: true })
</script>

<template>
  <n-card title="XML 格式化 — 输入" class="tool-card">
    <ToolTextarea v-model:value="xmlInput" label="输入 XML" :rows="20" placeholder="在此输入 XML…" monospace />

    <div class="options-row">
      <span class="option-label">缩进</span>
      <n-select v-model:value="indent" :options="indentOptions" class="indent-select" />
    </div>

    <n-alert v-if="errorMessage" type="error" class="error-alert">
      {{ errorMessage }}
    </n-alert>
  </n-card>

  <n-card title="XML 格式化 — 输出" class="tool-card">
    <CodeOutput label="格式化结果" :value="formatted" language="xml" :rows="20" />

    <div class="stats-row">
      <span class="stat">{{ lineCount }} 行</span>
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

.indent-select {
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
