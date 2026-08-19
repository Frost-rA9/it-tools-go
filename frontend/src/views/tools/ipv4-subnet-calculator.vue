<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NButton, NTable, NAlert, NSpace, useMessage, useThemeVars } from 'naive-ui'
import { ArrowLeft, ArrowRight } from '@vicons/tabler'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const address = ref('192.168.0.1/24')
const result = ref<Record<string, string | number> | null>(null)
const errorMessage = ref('')

const rows = [
  { key: 'netmask', label: 'Netmask' },
  { key: 'networkAddress', label: '网络地址（Network address）' },
  { key: 'networkMask', label: '子网掩码（Network mask）' },
  { key: 'networkMaskBinary', label: '掩码二进制' },
  { key: 'cidrNotation', label: 'CIDR 表示' },
  { key: 'wildcardMask', label: '通配符掩码' },
  { key: 'networkSize', label: '网络规模（地址数）' },
  { key: 'firstAddress', label: '首个可用地址' },
  { key: 'lastAddress', label: '末个可用地址' },
  { key: 'broadcastAddress', label: '广播地址' },
  { key: 'ipClass', label: 'IP 分类' },
]

async function run() {
  try {
    const output = await RunTool('ipv4-subnet-calculator', JSON.stringify({ address: address.value.trim() }))
    result.value = JSON.parse(output)
    errorMessage.value = ''
  } catch (e) {
    result.value = null
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 200)
watch(address, () => debouncedRun(), { immediate: true })

function jumpTo(block: string) {
  if (block) address.value = block
}

const copySource = ref('')
const { copy } = useClipboard({ source: copySource })

async function copyValue(v: string | number) {
  copySource.value = String(v)
  await copy()
  message.success('已复制到剪贴板')
}
</script>

<template>
  <div class="tool-page">
    <n-card title="IPv4 子网计算器" class="card">
      <div class="field">
        <div class="field-label">IPv4 地址（可带 CIDR 前缀或子网掩码）</div>
        <n-input
          v-model:value="address"
          placeholder="如 192.168.0.1/24、192.168.0.0/255.255.255.0 或纯 IP"
          clearable
          class="mono-input"
        />
      </div>

      <n-alert v-if="errorMessage" type="error" class="error-alert">
        {{ errorMessage }}
      </n-alert>

      <n-table v-if="result" :bordered="true" :single-line="false" class="result-table">
        <tbody>
          <tr v-for="{ key, label } in rows" :key="key">
            <td class="label-cell">{{ label }}</td>
            <td class="value-cell" :class="{ mono: typeof result[key] === 'string' }">
              <span class="copyable-value" :title="'复制 ' + label" @click="copyValue(result![key])">
                {{ result[key] }}
              </span>
            </td>
          </tr>
        </tbody>
      </n-table>

      <n-space justify="center" class="nav-row">
        <n-button :disabled="!result" @click="jumpTo(String(result?.prevBlock ?? ''))">
          <template #icon><n-icon :component="ArrowLeft" /></template>
          上一块
        </n-button>
        <n-button :disabled="!result" @click="jumpTo(String(result?.nextBlock ?? ''))">
          下一块
          <template #icon><n-icon :component="ArrowRight" /></template>
        </n-button>
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

.result-table {
  margin-bottom: 16px;
}

.label-cell {
  width: 46%;
  font-weight: 600;
  color: v-bind('themeVars.textColor2');
}

.value-cell {
  color: v-bind('themeVars.textColor1');
}

.mono {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}

.copyable-value {
  cursor: pointer;
  user-select: text;
}

.copyable-value:hover {
  color: v-bind('themeVars.primaryColor');
}

.nav-row {
  margin-top: 4px;
}
</style>
