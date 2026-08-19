<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NTable, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const TOOL_ID = 'iban-validator-and-parser'

interface IbanResult {
  valid: boolean
  errors: string[]
  qrIban: boolean
  countryCode: string
  bban: string
  friendlyFormat: string
}

interface Row {
  label: string
  value: string
  copyable: boolean
}

const rawIban = ref('')
const result = ref<IbanResult | null>(null)
const rows = ref<Row[]>([])

const ibanExamples = [
  'FR7630006000011234567890189',
  'DE89370400440532013000',
  'GB29NWBK60161331926819',
]

async function run() {
  const iban = rawIban.value
  if (!iban.trim()) {
    result.value = null
    rows.value = []
    return
  }
  try {
    const output = await RunTool(TOOL_ID, JSON.stringify({ iban }))
    const res: IbanResult = JSON.parse(output)
    result.value = res
    rows.value = [
      { label: 'IBAN 是否有效', value: res.valid ? '是' : '否', copyable: false },
      { label: 'IBAN 错误', value: res.errors.join('；'), copyable: false },
      { label: '是否 QR-IBAN', value: res.qrIban ? '是' : '否', copyable: false },
      { label: '国家代码', value: res.countryCode, copyable: true },
      { label: 'BBAN', value: res.bban, copyable: true },
      { label: '友好格式', value: res.friendlyFormat, copyable: true },
    ]
  } catch (e) {
    result.value = null
    rows.value = []
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch(rawIban, () => debouncedRun())

const { copy } = useClipboard()

async function copyValue(value: string) {
  if (!value) return
  await copy(value)
  message.success('已复制到剪贴板')
}
</script>

<template>
  <div class="iban-tool">
    <n-card title="IBAN 验证器和解析器" class="tool-card">
      <div class="field">
        <div class="field-label">IBAN</div>
        <n-input
          v-model:value="rawIban"
          placeholder="在此输入 IBAN 以验证有效性…"
          clearable
        />
      </div>

      <n-table v-if="rows.length" :bordered="true" :single-line="false" class="result-table">
        <tbody>
          <tr v-for="row in rows" :key="row.label">
            <td class="label-cell">{{ row.label }}</td>
            <td>
              <div class="value-cell">
                <span :class="{ 'error-text': row.label === 'IBAN 错误' && !!row.value }">
                  {{ row.value || '—' }}
                </span>
                <n-button v-if="row.copyable && row.value" size="tiny" quaternary @click="copyValue(row.value)">
                  复制
                </n-button>
              </div>
            </td>
          </tr>
        </tbody>
      </n-table>
    </n-card>

    <n-card title="有效 IBAN 示例" class="tool-card">
      <div v-for="iban in ibanExamples" :key="iban" class="example-row">
        <span class="example-value">{{ iban }}</span>
        <n-button size="tiny" quaternary @click="copyValue(iban)">复制</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.tool-card {
  max-width: 640px;
  margin-bottom: 16px;
}

.field {
  margin-bottom: 4px;
}

.field-label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.result-table {
  margin-top: 12px;
}

.label-cell {
  width: 140px;
  white-space: nowrap;
  font-weight: 600;
}

.value-cell {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  font-family: 'Cascadia Code', ui-monospace, monospace;
}

.error-text {
  color: v-bind('themeVars.errorColor');
}

.example-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 4px 0;
}

.example-value {
  font-family: 'Cascadia Code', ui-monospace, monospace;
}
</style>