<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NInputNumber, NButton, NAlert, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

interface FormatResult {
  label: string
  value: string
}

const value = ref('42')
const inputBase = ref<number | null>(10)
const customBase = ref<number | null>(42)
const results = ref<FormatResult[]>([])
const errorMessage = ref('')

async function run() {
  try {
    const output = await RunTool(
      'integer-base-converter',
      JSON.stringify({
        value: value.value,
        from_base: inputBase.value ?? 10,
        custom_base: customBase.value ?? 42,
      }),
    )
    results.value = JSON.parse(output).results
    errorMessage.value = ''
  } catch (e) {
    results.value = []
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch([value, inputBase, customBase], () => debouncedRun(), { immediate: true })

const copySource = ref('')
const { copy } = useClipboard({ source: copySource })

async function copyValue(v: string) {
  copySource.value = v
  await copy()
  message.success('已复制到剪贴板')
}
</script>

<template>
  <div class="baseconv-tool">
    <n-card title="整数基转换器" class="card">
      <div class="field">
        <div class="field-label">输入数字</div>
        <n-input v-model:value="value" placeholder="在此输入数字，如 42…" />
      </div>

      <div class="field">
        <div class="field-label">输入进制（2 - 64）</div>
        <n-input-number v-model:value="inputBase" :min="2" :max="64" :show-button="false" style="width: 100%" />
      </div>

      <div class="field">
        <div class="field-label">自定义输出进制（2 - 64）</div>
        <n-input-number v-model:value="customBase" :min="2" :max="64" :show-button="false" style="width: 100%" />
      </div>

      <n-alert v-if="errorMessage" type="error" class="error-alert">
        {{ errorMessage }}
      </n-alert>

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
.baseconv-tool {
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

.error-alert {
  margin-bottom: 16px;
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
  flex: 0 0 140px;
  text-align: right;
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
}

.result-value {
  flex: 1;
}
</style>
