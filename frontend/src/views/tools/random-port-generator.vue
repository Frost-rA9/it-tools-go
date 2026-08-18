<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInputNumber, NInput, NButton, NAlert, NSpace, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const count = ref<number | null>(5)
const min = ref<number | null>(1024)
const max = ref<number | null>(65535)
const excludeText = ref('')
const ports = ref<number[]>([])
const errorMessage = ref('')

async function run() {
  try {
    const exclude = excludeText.value
      .split(/[,，\s]+/)
      .map((s) => Number.parseInt(s, 10))
      .filter((n) => Number.isInteger(n))
    const output = await RunTool(
      'random-port-generator',
      JSON.stringify({
        count: count.value ?? 1,
        min: min.value ?? 1024,
        max: max.value ?? 65535,
        exclude,
      }),
    )
    ports.value = JSON.parse(output).ports
    errorMessage.value = ''
  } catch (e) {
    ports.value = []
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 200)
watch([count, min, max, excludeText], () => debouncedRun(), { immediate: true })

// 重新生成：直接触发后端（不受防抖影响），获得新的一组随机端口
function regenerate() {
  run()
}

const copySource = ref('')
const { copy } = useClipboard({ source: copySource })

async function copyPort(p: number) {
  copySource.value = String(p)
  await copy()
  message.success(`端口 ${p} 已复制到剪贴板`)
}
</script>

<template>
  <div class="tool-page">
    <n-card title="随机端口生成器" class="card">
      <div class="grid">
        <div class="field">
          <div class="field-label">生成数量</div>
          <n-input-number v-model:value="count" :min="1" :max="100" :show-button="false" style="width: 100%" />
        </div>
        <div class="field">
          <div class="field-label">最小值</div>
          <n-input-number v-model:value="min" :min="0" :max="65535" :show-button="false" style="width: 100%" />
        </div>
        <div class="field">
          <div class="field-label">最大值</div>
          <n-input-number v-model:value="max" :min="0" :max="65535" :show-button="false" style="width: 100%" />
        </div>
      </div>

      <div class="field">
        <div class="field-label">排除端口（逗号或空格分隔）</div>
        <n-input v-model:value="excludeText" placeholder="如 80, 443, 3306" />
      </div>

      <n-alert v-if="errorMessage" type="error" class="error-alert">
        {{ errorMessage }}
      </n-alert>

      <div class="field">
        <div class="field-label">生成的端口</div>
        <div class="result-list" v-if="ports.length">
          <div v-for="(p, idx) in ports" :key="p" class="result-row">
            <span class="result-label">端口 {{ idx + 1 }}</span>
            <n-input :value="String(p)" readonly class="result-value mono" placeholder="…" />
            <n-button size="small" @click="copyPort(p)">复制</n-button>
          </div>
        </div>
        <n-empty v-else description="暂无结果" size="small" />
      </div>

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
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
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
  flex: 0 0 80px;
  text-align: right;
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
}

.result-value {
  flex: 1;
}

.mono :deep(input) {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}
</style>
