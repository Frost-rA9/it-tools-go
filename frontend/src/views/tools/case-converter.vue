<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

interface FormatResult {
  label: string
  value: string
}

const input = ref('lorem ipsum dolor sit amet')
const results = ref<FormatResult[]>([])

async function run() {
  try {
    const output = await RunTool('case-converter', JSON.stringify({ text: input.value }))
    results.value = JSON.parse(output).results
  } catch (e) {
    results.value = []
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch(input, () => debouncedRun(), { immediate: true })

const copySource = ref('')
const { copy } = useClipboard({ source: copySource })

async function copyValue(value: string) {
  copySource.value = value
  await copy()
  message.success('已复制到剪贴板')
}
</script>

<template>
  <div class="case-tool">
    <n-card title="大小写转换" class="card">
      <div class="field">
        <div class="field-label">输入字符串</div>
        <n-input
          v-model:value="input"
          type="textarea"
          :rows="3"
          placeholder="在此输入字符串…"
        />
      </div>

      <div class="result-list">
        <div v-for="r in results" :key="r.label" class="result-row">
          <span class="result-label">{{ r.label }}</span>
          <n-input :value="r.value" readonly class="result-value" placeholder="…" />
          <n-button size="small" @click="copyValue(r.value)">复制</n-button>
        </div>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.case-tool {
  width: 100%;
}

.card {
  width: 100%;
}

.field {
  margin-bottom: 16px;
}

.field-label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
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
  flex: 0 0 120px;
  text-align: right;
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
}

.result-value {
  flex: 1;
}
</style>
