<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NButton, useMessage } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()

const textInput = ref('')
const unicodeOutput = ref('')

const unicodeInput = ref('')
const textOutput = ref('')

async function runTextToUnicode() {
  try {
    const output = await RunTool(
      'text-to-unicode',
      JSON.stringify({ text: textInput.value, mode: 'text_to_unicode' }),
    )
    unicodeOutput.value = JSON.parse(output).result
  } catch (e) {
    unicodeOutput.value = ''
    message.error(String(e))
  }
}

async function runUnicodeToText() {
  try {
    const output = await RunTool(
      'text-to-unicode',
      JSON.stringify({ text: unicodeInput.value, mode: 'unicode_to_text' }),
    )
    textOutput.value = JSON.parse(output).result
  } catch (e) {
    textOutput.value = ''
  }
}

const debouncedText = useDebounceFn(runTextToUnicode, 150)
const debouncedUnicode = useDebounceFn(runUnicodeToText, 150)

watch(textInput, () => debouncedText(), { immediate: true })
watch(unicodeInput, () => debouncedUnicode(), { immediate: true })

const { copy: copyUnicode } = useClipboard({ source: unicodeOutput })
const { copy: copyText } = useClipboard({ source: textOutput })

async function copyUnicodeResult() {
  await copyUnicode()
  message.success('Unicode 已复制到剪贴板')
}

async function copyTextResult() {
  await copyText()
  message.success('文本已复制到剪贴板')
}
</script>

<template>
  
    <n-card title="文本转 Unicode" class="tool-card">
      <ToolTextarea v-model:value="textInput" label="输入文本" :rows="4" placeholder="在此输入文本，如 Hello…" />

      <ToolTextarea v-model:value="unicodeOutput" label="Unicode 结果" :rows="4" readonly placeholder="文本的 Unicode 表示将显示在这里" />

      <div class="copy-row">
        <n-button type="primary" @click="copyUnicodeResult">复制 Unicode</n-button>
      </div>
    </n-card>

    <n-card title="Unicode 转文本" class="tool-card">
      <ToolTextarea v-model:value="unicodeInput" label="输入 Unicode" :rows="4" placeholder="在此输入 Unicode，如 &amp;#72;&amp;#105;…" />

      <ToolTextarea v-model:value="textOutput" label="文本结果" :rows="4" readonly placeholder="Unicode 对应的文本将显示在这里" />

      <div class="copy-row">
        <n-button type="primary" @click="copyTextResult">复制文本</n-button>
      </div>
    </n-card>
</template>

<style scoped>

.copy-row {
  display: flex;
  justify-content: center;
}
</style>
