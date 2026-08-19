<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NButton, NAlert, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const expression = ref('2*sqrt(6)')
const result = ref('')
const error = ref('')

async function run() {
  const expr = expression.value.trim()
  if (!expr) {
    result.value = ''
    error.value = ''
    return
  }
  try {
    const output = await RunTool('math-evaluator', JSON.stringify({ expression: expr }))
    const parsed = JSON.parse(output)
    result.value = parsed.result
    error.value = parsed.error
  } catch (e) {
    result.value = ''
    error.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch(expression, () => debouncedRun(), { immediate: true })

const { copy } = useClipboard()

async function copyResult() {
  await copy(result.value)
  message.success('结果已复制到剪贴板')
}
</script>

<template>
  <div class="evaluator-tool">
    <n-card title="数学表达式求值器" class="tool-card">
      <div class="field">
        <div class="field-label">数学表达式</div>
        <n-input
          v-model:value="expression"
          type="textarea"
          :autosize="{ minRows: 2, maxRows: 8 }"
          placeholder="例如：2*sqrt(6)、sin(pi/2)、2^3^2、log(100)…，支持四则运算、幂、括号与内置函数"
        />
      </div>

      <n-alert v-if="error" type="error" :show-icon="true" class="error-alert">
        {{ error }}
      </n-alert>

      <div v-else class="field">
        <div class="field-label">结果</div>
        <n-input :value="result" readonly placeholder="求值结果将显示在这里" />
      </div>

      <div class="copy-row">
        <n-button type="primary" :disabled="!result" @click="copyResult">复制结果</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.tool-card {
  max-width: 640px;
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

.copy-row {
  display: flex;
  justify-content: center;
}
</style>