<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NButton, useMessage } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()

const escapeInput = ref('<title>IT Tool</title>')
const escapeOutput = ref('')

const unescapeInput = ref('&lt;title&gt;IT Tool&lt;/title&gt;')
const unescapeOutput = ref('')

async function runEscape() {
  try {
    const output = await RunTool('html-entities', JSON.stringify({ text: escapeInput.value, action: 'escape' }))
    escapeOutput.value = JSON.parse(output).result
  } catch (e) {
    escapeOutput.value = ''
    message.error(String(e))
  }
}

async function runUnescape() {
  try {
    const output = await RunTool('html-entities', JSON.stringify({ text: unescapeInput.value, action: 'unescape' }))
    unescapeOutput.value = JSON.parse(output).result
  } catch (e) {
    unescapeOutput.value = ''
    message.error(String(e))
  }
}

const debouncedEscape = useDebounceFn(runEscape, 150)
const debouncedUnescape = useDebounceFn(runUnescape, 150)

watch(escapeInput, () => debouncedEscape(), { immediate: true })
watch(unescapeInput, () => debouncedUnescape(), { immediate: true })

const { copy: copyEscaped } = useClipboard({ source: escapeOutput })
const { copy: copyUnescaped } = useClipboard({ source: unescapeOutput })

async function copyEscapedResult() {
  await copyEscaped()
  message.success('转义结果已复制到剪贴板')
}

async function copyUnescapedResult() {
  await copyUnescaped()
  message.success('反转义结果已复制到剪贴板')
}
</script>

<template>
  <n-card title="转义 HTML 实体" class="tool-card">
    <ToolTextarea v-model:value="escapeInput" label="待转义字符串" :rows="3" placeholder="要转义的字符串…" />

    <ToolTextarea v-model:value="escapeOutput" label="转义后的字符串" :rows="3" readonly placeholder="转义后的字符串将显示在这里" />

    <div class="copy-row">
      <n-button type="primary" @click="copyEscapedResult">复制转义结果</n-button>
    </div>
  </n-card>

  <n-card title="反转义 HTML 实体" class="tool-card">
    <ToolTextarea v-model:value="unescapeInput" label="待反转义字符串" :rows="3" placeholder="要反转义的字符串…" />

    <ToolTextarea v-model:value="unescapeOutput" label="反转义后的字符串" :rows="3" readonly placeholder="反转义后的字符串将显示在这里" />

    <div class="copy-row">
      <n-button type="primary" @click="copyUnescapedResult">复制反转义结果</n-button>
    </div>
  </n-card>
</template>

<style scoped>

.copy-row {
  display: flex;
  justify-content: center;
}
</style>