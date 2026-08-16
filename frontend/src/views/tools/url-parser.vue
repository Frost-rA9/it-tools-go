<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NButton, NText, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

interface UrlPart {
  protocol: string
  username: string
  password: string
  hostname: string
  port: string
  pathname: string
  search: string
  params: { key: string; value: string }[]
}

const message = useMessage()
const themeVars = useThemeVars()

const urlInput = ref('https://user:pass@example.com:3000/path?key1=value&key2=value2#the-hash')
const parsed = ref<UrlPart | null>(null)

const props: { label: string; key: 'protocol' | 'username' | 'password' | 'hostname' | 'port' | 'pathname' | 'search' }[] = [
  { label: 'Protocol', key: 'protocol' },
  { label: 'Username', key: 'username' },
  { label: 'Password', key: 'password' },
  { label: 'Hostname', key: 'hostname' },
  { label: 'Port', key: 'port' },
  { label: 'Path', key: 'pathname' },
  { label: 'Params', key: 'search' },
]

async function runParse() {
  try {
    const output = await RunTool('url-parser', JSON.stringify({ url: urlInput.value }))
    parsed.value = JSON.parse(output)
  } catch (e) {
    parsed.value = null
    message.error(String(e))
  }
}

const debouncedParse = useDebounceFn(runParse, 150)
watch(urlInput, () => debouncedParse(), { immediate: true })

const { copy } = useClipboard()

async function copyValue(label: string, value: string) {
  await copy(value)
  message.success(`${label} 已复制到剪贴板`)
}
</script>

<template>
  <n-card class="card">
    <ToolTextarea v-model:value="urlInput" label="要解析的 URL" :rows="2" placeholder="要解析的 URL…" />

    <div class="prop-row" v-for="p in props" :key="p.key">
      <span class="prop-label">{{ p.label }}</span>
      <span class="prop-value">{{ parsed?.[p.key] || '—' }}</span>
      <n-button size="tiny" secondary @click="copyValue(p.label, parsed?.[p.key] ?? '')">复制</n-button>
    </div>

    <n-divider />

    <div class="param-row" v-for="(param, i) in parsed?.params ?? []" :key="i">
      <span class="prop-label">Param {{ i + 1 }}</span>
      <span class="param-key">{{ param.key }}</span>
      <span class="prop-value">{{ param.value }}</span>
      <n-button size="tiny" secondary @click="copyValue(`参数 ${param.key}`, param.value)">复制</n-button>
    </div>

    <div v-if="!parsed?.params?.length" class="empty">
      <n-text depth="3">无查询参数</n-text>
    </div>
  </n-card>
</template>

<style scoped>
.card {
  min-width: 480px;
  width: 100%;
}

.prop-row,
.param-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 0;
}

.prop-label {
  flex: 0 0 110px;
  font-size: 13px;
  color: v-bind('themeVars.textColor3');
}

.prop-value {
  flex: 1 1 auto;
  font-family: 'Cascadia Code', monospace;
  font-size: 13px;
  word-break: break-all;
  color: v-bind('themeVars.textColor1');
}

.param-key {
  flex: 0 0 auto;
  max-width: 180px;
  font-family: 'Cascadia Code', monospace;
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: v-bind('themeVars.primaryColor');
}

.empty {
  padding: 8px 0;
}
</style>