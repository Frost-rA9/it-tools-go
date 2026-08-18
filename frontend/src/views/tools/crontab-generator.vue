<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { NCard, NSelect, NInputNumber, NInput, NButton, NAlert, NDivider, NTable, NIcon, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { Copy } from '@vicons/tabler'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

type FieldType = 'every' | 'step' | 'range' | 'list'

interface FieldState {
  type: FieldType
  step: number
  from: number
  to: number
  valuesText: string
}

interface FieldMeta {
  key: string
  label: string
  min: number
  max: number
  unit: string
}

const fieldsMeta: FieldMeta[] = [
  { key: 'minute', label: '分钟', min: 0, max: 59, unit: '分钟' },
  { key: 'hour', label: '小时', min: 0, max: 23, unit: '小时' },
  { key: 'day_of_month', label: '日', min: 1, max: 31, unit: '天' },
  { key: 'month', label: '月', min: 1, max: 12, unit: '月' },
  { key: 'day_of_week', label: '星期', min: 0, max: 7, unit: '周' },
]

const typeOptions = [
  { label: '任意（*）', value: 'every' },
  { label: '间隔（*/N）', value: 'step' },
  { label: '范围（A-B）', value: 'range' },
  { label: '列表（A,B,C）', value: 'list' },
]

const fields = reactive<Record<string, FieldState>>({})
for (const m of fieldsMeta) {
  fields[m.key] = { type: 'every', step: 5, from: m.min, to: m.max, valuesText: '' }
}

// 表达式输入框：表单生成结果写入；用户直接编辑时进入 parse 模式
const expression = ref('* * * * *')
const summary = ref('')
const errorMessage = ref('')
const syncing = ref(false)

function buildInput() {
  const out: Record<string, unknown> = { action: 'generate' }
  for (const m of fieldsMeta) {
    const f = fields[m.key]
    const item: Record<string, unknown> = { type: f.type }
    if (f.type === 'step') item.step = f.step
    if (f.type === 'range') {
      item.from = f.from
      item.to = f.to
    }
    if (f.type === 'list') {
      item.values = f.valuesText
        .split(/[,，\s]+/)
        .map((s) => Number.parseInt(s, 10))
        .filter((n) => Number.isInteger(n))
    }
    out[m.key] = item
  }
  return out
}

async function run(input: Record<string, unknown>) {
  try {
    const output = await RunTool('crontab-generator', JSON.stringify(input))
    const result = JSON.parse(output)
    errorMessage.value = ''
    if (input.action === 'generate') {
      // 表单生成：写回表达式输入框（受控，不触发编辑联动）
      syncing.value = true
      expression.value = result.expression
      syncing.value = false
      summary.value = result.summary
    } else {
      summary.value = result.summary
    }
  } catch (e) {
    summary.value = ''
    errorMessage.value = String(e)
  }
}

// 表单变化 → generate
const debouncedGenerate = useDebounceFn(() => run(buildInput()), 150)
watch(
  () => fieldsMeta.map((m) => JSON.stringify(fields[m.key])),
  () => debouncedGenerate(),
  { immediate: true },
)

// 表达式输入框直接编辑 → parse
const debouncedParse = useDebounceFn(() => {
  if (!syncing.value && expression.value.trim()) {
    run({ action: 'parse', expression: expression.value })
  }
}, 200)
watch(expression, () => debouncedParse())

const copySource = ref('')
const { copy } = useClipboard({ source: copySource })

async function copyExpression() {
  copySource.value = expression.value
  await copy()
  message.success('Cron 表达式已复制到剪贴板')
}

// 帮助表格数据（对齐 it-tools helpers）
const helpers = [
  { symbol: '*', meaning: '任意值', example: '* * * * *', equivalent: '每分钟' },
  { symbol: '-', meaning: '范围', example: '1-10 * * * *', equivalent: '第 1 至 10 分钟' },
  { symbol: ',', meaning: '列表', example: '1,10 * * * *', equivalent: '第 1 和 10 分钟' },
  { symbol: '/', meaning: '步长', example: '*/10 * * * *', equivalent: '每 10 分钟' },
  { symbol: '@yearly', meaning: '每年 1 月 1 日午夜执行', example: '@yearly', equivalent: '0 0 1 1 *' },
  { symbol: '@monthly', meaning: '每月 1 日午夜执行', example: '@monthly', equivalent: '0 0 1 * *' },
  { symbol: '@weekly', meaning: '每周日午夜执行', example: '@weekly', equivalent: '0 0 * * 0' },
  { symbol: '@daily', meaning: '每天午夜执行', example: '@daily', equivalent: '0 0 * * *' },
  { symbol: '@hourly', meaning: '每小时整点执行', example: '@hourly', equivalent: '0 * * * *' },
  { symbol: '@reboot', meaning: '系统启动时执行', example: '@reboot', equivalent: '—' },
]
</script>

<template>
  <n-card title="Crontab 表达式生成器" class="tool-card">
      <div class="field-row" v-for="m in fieldsMeta" :key="m.key">
        <div class="field-label">{{ m.label }}</div>
        <div class="field-controls">
          <n-select v-model:value="fields[m.key].type" :options="typeOptions" class="type-select" />
          <template v-if="fields[m.key].type === 'step'">
            <span class="hint">每</span>
            <n-input-number v-model:value="fields[m.key].step" :min="1" :max="m.max - m.min + 1" :show-button="false" style="width: 110px" />
            <span class="hint">{{ m.unit }}</span>
          </template>
          <template v-else-if="fields[m.key].type === 'range'">
            <n-input-number v-model:value="fields[m.key].from" :min="m.min" :max="m.max" :show-button="false" style="width: 110px" />
            <span class="hint">到</span>
            <n-input-number v-model:value="fields[m.key].to" :min="m.min" :max="m.max" :show-button="false" style="width: 110px" />
          </template>
          <template v-else-if="fields[m.key].type === 'list'">
            <n-input
              v-model:value="fields[m.key].valuesText"
              :placeholder="`逗号分隔，范围 ${m.min}-${m.max}`"
              style="width: 260px"
            />
          </template>
        </div>
      </div>

      <n-divider />

      <div class="expr-area">
        <n-input
          v-model:value="expression"
          placeholder="* * * * *"
          class="expr-input"
          :status="errorMessage ? 'error' : undefined"
        />
        <div class="expr-suffix">
          <n-button quaternary circle size="small" :disabled="!expression" @click="copyExpression" title="复制表达式">
            <n-icon :component="Copy" />
          </n-button>
        </div>
      </div>

      <div class="cron-string" :class="{ invalid: errorMessage }">
        {{ errorMessage || summary }}
      </div>

      <n-alert v-if="errorMessage" type="error" class="error-alert">
        {{ errorMessage }}
      </n-alert>
    </n-card>

    <n-card title="Cron 语法说明" class="tool-card">
      <pre class="cron-diagram">
┌─────────── [可选] 秒 (0 - 59)
| ┌───────── 分钟 (0 - 59)
| | ┌─────── 小时 (0 - 23)
| | | ┌───── 日 (1 - 31)
| | | | ┌─── 月 (1 - 12) 或 jan,feb,mar...
| | | | | ┌─ 星期 (0 - 6，周日=0) 或 sun,mon...
| | | | | |
* * * * * * 命令</pre>

      <n-table :bordered="true" size="small" class="helper-table">
        <thead>
          <tr>
            <th>符号</th>
            <th>含义</th>
            <th>示例</th>
            <th>等价描述</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="h in helpers" :key="h.symbol">
            <td><code class="mono">{{ h.symbol }}</code></td>
            <td>{{ h.meaning }}</td>
            <td><code class="mono">{{ h.example }}</code></td>
            <td><code class="mono">{{ h.equivalent }}</code></td>
          </tr>
        </tbody>
      </n-table>
    </n-card>
</template>

<style scoped>

.field-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.field-label {
  flex: 0 0 64px;
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
}

.field-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.type-select {
  width: 150px;
}

.hint {
  font-size: 14px;
  color: v-bind('themeVars.textColor3');
}

.expr-area {
  position: relative;
}

.expr-input :deep(input) {
  font-size: 26px;
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
  text-align: center;
  padding: 6px 40px 6px 12px;
}

.expr-suffix {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 2;
}

.cron-string {
  text-align: center;
  font-size: 20px;
  opacity: 0.8;
  margin: 10px 0 0;
  min-height: 28px;
  color: v-bind('themeVars.textColor1');
}

.cron-string.invalid {
  color: v-bind('themeVars.errorColor');
  opacity: 1;
}

.error-alert {
  margin-bottom: 16px;
}

.cron-diagram {
  overflow: auto;
  padding: 10px 0;
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  color: v-bind('themeVars.textColor2');
}

.helper-table {
  margin-top: 8px;
}

.mono {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}
</style>
