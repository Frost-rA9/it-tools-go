<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NButton, NAlert, NTable, NSpace, useMessage, useThemeVars } from 'naive-ui'
import { Exchange } from '@vicons/tabler'
import { useDebounceFn } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const startIp = ref('192.168.1.1')
const endIp = ref('192.168.6.255')
const result = ref<Record<string, string | number> | null>(null)
const errorMessage = ref('')
const isOrderError = ref(false)

const rows = [
  { key: 'start', label: '起始地址' },
  { key: 'end', label: '结束地址' },
  { key: 'size', label: '范围内地址数' },
  { key: 'cidr', label: '覆盖 CIDR' },
]

async function run() {
  try {
    const output = await RunTool(
      'ipv4-range-expander',
      JSON.stringify({ startIp: startIp.value.trim(), endIp: endIp.value.trim() }),
    )
    result.value = JSON.parse(output)
    errorMessage.value = ''
    isOrderError.value = false
  } catch (e) {
    result.value = null
    const msg = String(e)
    errorMessage.value = msg
    isOrderError.value = msg.includes('低于起始地址')
  }
}

const debouncedRun = useDebounceFn(run, 200)
watch([startIp, endIp], () => debouncedRun(), { immediate: true })

function switchStartEnd() {
  const tmp = startIp.value
  startIp.value = endIp.value
  endIp.value = tmp
}

function fmt(n: string | number): string {
  return typeof n === 'number' ? n.toLocaleString() : n
}
</script>

<template>
  <div class="tool-page">
    <n-card title="IPv4 范围扩展器" class="card">
      <div class="grid">
        <div class="field">
          <div class="field-label">起始地址</div>
          <n-input v-model:value="startIp" placeholder="起始 IPv4 地址…" clearable class="mono-input" />
        </div>
        <div class="field">
          <div class="field-label">结束地址</div>
          <n-input v-model:value="endIp" placeholder="结束 IPv4 地址…" clearable class="mono-input" />
        </div>
      </div>

      <n-alert v-if="errorMessage && !isOrderError" type="error" class="error-alert">
        {{ errorMessage }}
      </n-alert>

      <n-alert v-if="isOrderError" type="error" title="起始与结束地址顺序无效" class="error-alert">
        <p class="order-msg">结束地址低于起始地址，无法计算。通常交换起始与结束地址即可解决。</p>
        <n-button size="small" @click="switchStartEnd">
          <template #icon><n-icon :component="Exchange" /></template>
          交换起始与结束地址
        </n-button>
      </n-alert>

      <n-table v-if="result" :bordered="true" :single-line="false" class="result-table">
        <thead>
          <tr>
            <th />
            <th>原值</th>
            <th>覆盖后</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="{ key, label } in rows" :key="key">
            <td class="label-cell">{{ label }}</td>
            <td class="mono">
              {{
                key === 'start'
                  ? startIp.trim()
                  : key === 'end'
                    ? endIp.trim()
                    : key === 'size'
                      ? fmt(result.oldSize)
                      : '—'
              }}
            </td>
            <td class="mono">
              {{
                key === 'start'
                  ? result.newStart
                  : key === 'end'
                    ? result.newEnd
                    : key === 'size'
                      ? fmt(result.newSize)
                      : result.newCidr
              }}
            </td>
          </tr>
        </tbody>
      </n-table>
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
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
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

.order-msg {
  margin: 0 0 10px;
  opacity: 0.7;
}

.result-table {
  margin-bottom: 8px;
}

.label-cell {
  font-weight: 600;
  color: v-bind('themeVars.textColor2');
}

.mono {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}
</style>
