<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NAlert, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'
import ToolCodeBlock from '../../components/ToolCodeBlock.vue'

const message = useMessage()
const themeVars = useThemeVars()

const ip = ref('192.168.1.1')
const result = ref<Record<string, string | number> | null>(null)
const errorMessage = ref('')

const sections = [
  { key: 'decimal', label: '十进制（Decimal）' },
  { key: 'hexadecimal', label: '十六进制（Hexadecimal）' },
  { key: 'binary', label: '二进制（Binary）' },
  { key: 'ipv6', label: 'IPv6 映射（完整）' },
  { key: 'ipv6Short', label: 'IPv6 映射（压缩）' },
]

async function run() {
  try {
    const output = await RunTool('ipv4-address-converter', JSON.stringify({ ip: ip.value.trim() }))
    result.value = JSON.parse(output)
    errorMessage.value = ''
  } catch (e) {
    result.value = null
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 200)
watch(ip, () => debouncedRun(), { immediate: true })
</script>

<template>
  <div class="tool-page">
    <n-card title="IPv4 地址转换器" class="card">
      <div class="field">
        <div class="field-label">IPv4 地址</div>
        <n-input
          v-model:value="ip"
          placeholder="输入 IPv4 地址，如 192.168.1.1"
          clearable
          class="mono-input"
        />
      </div>

      <n-alert v-if="errorMessage" type="error" class="error-alert">
        {{ errorMessage }}
      </n-alert>

      <template v-if="result">
        <ToolCodeBlock
          v-for="{ key, label } in sections"
          :key="key"
          :label="label"
          :value="String(result[key])"
          align="center"
          copyable
        />
      </template>
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
</style>
