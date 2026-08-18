<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NButton, NDivider, NText, NIcon, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { Copy, ArrowDownRight } from '@vicons/tabler'
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

    <n-divider />

    <div class="prop-row" v-for="p in props" :key="p.key">
      <span class="prop-label">{{ p.label }}</span>
      <n-input :value="parsed?.[p.key] ?? ''" readonly size="small" class="prop-input" placeholder=" ">
        <template #suffix>
          <n-button circle quaternary size="small" @click="copyValue(p.label, parsed?.[p.key] ?? '')">
            <n-icon :component="Copy" :size="16" />
          </n-button>
        </template>
      </n-input>
    </div>

    <div class="param-row" v-for="(param, i) in parsed?.params ?? []" :key="i">
      <div class="arrow-col">
        <n-icon :component="ArrowDownRight" :size="14" class="arrow-icon" />
      </div>
      <n-input :value="param.key" readonly size="small" class="param-key-input" placeholder=" ">
        <template #suffix>
          <n-button circle quaternary size="small" @click="copyValue(`参数 ${param.key}`, param.key)">
            <n-icon :component="Copy" :size="16" />
          </n-button>
        </template>
      </n-input>
      <n-input :value="param.value" readonly size="small" class="prop-input" placeholder=" ">
        <template #suffix>
          <n-button circle quaternary size="small" @click="copyValue(`参数 ${param.key}`, param.value)">
            <n-icon :component="Copy" :size="16" />
          </n-button>
        </template>
      </n-input>
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

.prop-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.prop-label {
  flex: 0 0 110px;
  text-align: left;
  font-size: 13px;
  color: v-bind('themeVars.textColor3');
}

.prop-input {
  flex: 1;
}

.prop-input :deep(input),
.param-key-input :deep(input) {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
  font-size: 13px;
}

.param-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.arrow-col {
  flex: 0 0 110px;
  display: flex;
  justify-content: flex-start;
  padding-left: 2px;
}

.arrow-icon {
  color: v-bind('themeVars.textColor2');
}

.param-key-input {
  flex: 0 0 180px;
}

.empty {
  padding: 8px 0;
}
</style>
