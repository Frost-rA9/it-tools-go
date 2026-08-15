<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInputNumber, NRadioGroup, NRadioButton, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'
import ToolCodeBlock from '../../components/ToolCodeBlock.vue'

const message = useMessage()
const themeVars = useThemeVars()

const count = ref(1)
const format = ref('raw')
const result = ref('')

async function run() {
  try {
    const output = await RunTool('ulid-generator', JSON.stringify({ count: count.value, format: format.value }))
    result.value = JSON.parse(output).result
  } catch (e) {
    result.value = ''
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch([count, format], () => debouncedRun(), { immediate: true })

const { copy } = useClipboard({ source: result })

async function copyResult() {
  await copy()
  message.success('ULID 已复制到剪贴板')
}
</script>

<template>
  <div class="ulid-tool">
    <n-card title="ULID 生成器" class="card">
      <div class="field">
        <div class="field-label">生成数量</div>
        <n-input-number v-model:value="count" :min="1" :max="100" :step="1" style="width: 200px" />
      </div>

      <div class="field">
        <div class="field-label">输出格式</div>
        <n-radio-group v-model:value="format">
          <n-radio-button value="raw">Raw</n-radio-button>
          <n-radio-button value="json">JSON</n-radio-button>
        </n-radio-group>
      </div>

      <ToolCodeBlock label="生成的 ULID" :value="result" align="center" />

      <div class="btn-row">
        <n-button type="primary" @click="copyResult">复制</n-button>
        <n-button @click="run">重新生成</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.ulid-tool {
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

.btn-row {
  display: flex;
  justify-content: center;
  gap: 12px;
}
</style>
