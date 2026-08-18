<script setup lang="ts">
import { ref, watch } from 'vue'
import { NSelect, NSwitch, NButton, NAlert, NSpace, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const jsonInput = ref('{"name":"it-tools","version":2,"features":["json","csv"],"active":true}')
const indent = ref('2')
const sortKeys = ref(false)
const formatted = ref('')
const lineCount = ref(0)
const charCount = ref(0)
const errorMessage = ref('')

const indentOptions = [
  { label: '2 空格', value: '2' },
  { label: '4 空格', value: '4' },
  { label: 'Tab', value: '\t' },
]

async function run() {
  try {
    const output = await RunTool(
      'json-formatter',
      JSON.stringify({ json: jsonInput.value, indent: indent.value, sort_keys: sortKeys.value }),
    )
    const result = JSON.parse(output)
    formatted.value = result.formatted
    lineCount.value = result.line_count
    charCount.value = result.char_count
    errorMessage.value = ''
  } catch (e) {
    formatted.value = ''
    lineCount.value = 0
    charCount.value = 0
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 200)
watch([jsonInput, indent, sortKeys], () => debouncedRun(), { immediate: true })

const copySource = ref('')
const { copy } = useClipboard({ source: copySource })

async function copyFormatted() {
  copySource.value = formatted.value
  await copy()
  message.success('格式化结果已复制到剪贴板')
}
</script>

<template>
  <n-card title="JSON 格式化 — 输入" class="tool-card">
    <ToolTextarea v-model:value="jsonInput" label="输入 JSON" :rows="10" placeholder="在此粘贴 JSON…" monospace />

    <div class="options-row">
      <span class="option-label">缩进</span>
      <n-select v-model:value="indent" :options="indentOptions" class="indent-select" />
      <span class="option-label">按键排序</span>
      <n-switch v-model:value="sortKeys" />
    </div>

    <n-alert v-if="errorMessage" type="error" class="error-alert">
      {{ errorMessage }}
    </n-alert>
  </n-card>

  <n-card title="JSON 格式化 — 输出" class="tool-card">
    <ToolTextarea v-model:value="formatted" label="格式化结果" :rows="10" placeholder="格式化结果将显示在这里" readonly monospace />

    <div class="stats-row">
      <span class="stat">{{ lineCount }} 行</span>
      <span class="stat">{{ charCount }} 字符</span>
    </div>

    <n-space justify="center">
      <n-button type="primary" :disabled="!formatted" @click="copyFormatted">复制结果</n-button>
    </n-space>
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
