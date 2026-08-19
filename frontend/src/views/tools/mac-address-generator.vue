<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInputNumber, NInput, NRadioGroup, NRadioButton, NButton, NAlert, NSpace, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'
import ToolCodeBlock from '../../components/ToolCodeBlock.vue'

const message = useMessage()
const themeVars = useThemeVars()

const count = ref<number | null>(1)
const prefix = ref('64:16:7F')
const caseMode = ref('upper')
const separator = ref(':')
const result = ref('')
const errorMessage = ref('')

const separators = [
  { label: ':', value: ':' },
  { label: '-', value: '-' },
  { label: '.', value: '.' },
  { label: '无', value: '' },
]

async function run() {
  try {
    const output = await RunTool(
      'mac-address-generator',
      JSON.stringify({
        count: count.value ?? 1,
        prefix: prefix.value.trim(),
        separator: separator.value,
        case: caseMode.value,
      }),
    )
    result.value = JSON.parse(output).macAddresses
    errorMessage.value = ''
  } catch (e) {
    result.value = ''
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 200)
watch([count, prefix, caseMode, separator], () => debouncedRun(), { immediate: true })

// 重新生成：直接调用后端（新随机字节），不受防抖影响
function regenerate() {
  run()
}
</script>

<template>
  <div class="tool-page">
    <n-card title="MAC 地址生成器" class="card">
      <div class="grid">
        <div class="field">
          <div class="field-label">生成数量</div>
          <n-input-number v-model:value="count" :min="1" :max="100" :show-button="false" style="width: 100%" />
        </div>
        <div class="field">
          <div class="field-label">大小写</div>
          <n-radio-group v-model:value="caseMode">
            <n-radio-button value="upper">Uppercase</n-radio-button>
            <n-radio-button value="lower">Lowercase</n-radio-button>
          </n-radio-group>
        </div>
        <div class="field">
          <div class="field-label">分隔符</div>
          <n-radio-group v-model:value="separator">
            <n-radio-button v-for="s in separators" :key="s.label" :value="s.value">
              {{ s.label }}
            </n-radio-button>
          </n-radio-group>
        </div>
      </div>

      <div class="field">
        <div class="field-label">MAC 前缀（可选）</div>
        <n-input
          v-model:value="prefix"
          placeholder="如 64:16:7F（最多 6 字节，支持 : / - / . 分隔或连续十六进制）"
          clearable
          class="mono-input"
        />
      </div>

      <n-alert v-if="errorMessage" type="error" class="error-alert">
        {{ errorMessage }}
      </n-alert>

      <ToolCodeBlock
        label="生成的 MAC 地址"
        :value="result"
        align="left"
        copyable
        class="result-block"
      />

      <n-space justify="center">
        <n-button type="primary" @click="regenerate">重新生成</n-button>
      </n-space>
    </n-card>
  </div>
</template>

<style scoped>
.tool-page {
  width: 100%;
}

.card {
  width: 100%;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}

.field {
  margin-bottom: 16px;
}

.field-label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.mono-input :deep(input) {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}

.error-alert {
  margin-bottom: 16px;
}

.result-block {
  margin-top: 4px;
}
</style>
