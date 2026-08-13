<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NSelect, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const formatNames = [
  'Locale string',
  'ISO 8601',
  'ISO 9075',
  'RFC 3339',
  'RFC 7231',
  'Unix timestamp',
  'Timestamp',
  'UTC format',
  'Mongo ObjectID',
  'Excel date/time',
]

const formatOptions = formatNames.map((name, i) => ({ label: name, value: i }))

interface FormatResult {
  name: string
  value: string
}

const inputDate = ref('')
const formatIndex = ref(6)
const results = ref<FormatResult[]>([])

async function run() {
  try {
    const output = await RunTool(
      'date-time-converter',
      JSON.stringify({ value: inputDate.value, format: formatIndex.value }),
    )
    const parsed = JSON.parse(output)
    results.value = parsed.results
    if (parsed.detected >= 0) {
      formatIndex.value = parsed.detected
    }
  } catch (e) {
    results.value = []
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch([inputDate, formatIndex], () => debouncedRun(), { immediate: true })

const copySource = ref('')
const { copy } = useClipboard({ source: copySource })

async function copyValue(value: string) {
  copySource.value = value
  await copy()
  message.success('已复制到剪贴板')
}
</script>

<template>
  <div class="datetime-tool">
    <n-card title="日期时间转换器" class="card">
      <div class="input-row">
        <n-input
          v-model:value="inputDate"
          class="input-field"
          placeholder="在此输入日期时间字符串…"
          clearable
        />
        <n-select
          v-model:value="formatIndex"
          :options="formatOptions"
          class="format-select"
        />
      </div>

      <div class="result-list">
        <div v-for="r in results" :key="r.name" class="result-row">
          <span class="result-label">{{ r.name }}</span>
          <n-input :value="r.value" readonly class="result-value" placeholder="无效日期…" />
          <n-button size="small" @click="copyValue(r.value)">复制</n-button>
        </div>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.datetime-tool {
  width: 100%;
}

.card {
  width: 100%;
}

.input-row {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.input-field {
  flex: 1;
}

.format-select {
  flex: 0 0 170px;
}

.result-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.result-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.result-label {
  flex: 0 0 150px;
  text-align: right;
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
}

.result-value {
  flex: 1;
}
</style>
