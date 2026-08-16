<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NButton, useMessage } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()

const textInput = ref('Hello world :)')
const urlEncodedOutput = ref('')

const urlInput = ref('Hello%20world%20%3A)')
const textOutput = ref('')

async function runEncode() {
  try {
    const output = await RunTool(
      'url-encoder-decoder',
      JSON.stringify({ text: textInput.value, action: 'encode' }),
    )
    urlEncodedOutput.value = JSON.parse(output).result
  } catch (e) {
    urlEncodedOutput.value = ''
    message.error(String(e))
  }
}

async function runDecode() {
  try {
    const output = await RunTool(
      'url-encoder-decoder',
      JSON.stringify({ text: urlInput.value, action: 'decode' }),
    )
    textOutput.value = JSON.parse(output).result
  } catch (e) {
    textOutput.value = ''
    message.error(String(e))
  }
}

const debouncedEncode = useDebounceFn(runEncode, 150)
const debouncedDecode = useDebounceFn(runDecode, 150)

watch(textInput, () => debouncedEncode(), { immediate: true })
watch(urlInput, () => debouncedDecode(), { immediate: true })

const { copy: copyEncoded } = useClipboard({ source: urlEncodedOutput })
const { copy: copyDecoded } = useClipboard({ source: textOutput })

async function copyEncodedResult() {
  await copyEncoded()
  message.success('编码结果已复制到剪贴板')
}

async function copyDecodedResult() {
  await copyDecoded()
  message.success('解码结果已复制到剪贴板')
}
</script>

<template>
  <n-card title="文本转 URL 编码" class="card">
    <ToolTextarea v-model:value="textInput" label="待编码文本" :rows="5" placeholder="在此输入文本…" />

    <ToolTextarea v-model:value="urlEncodedOutput" label="URL 编码结果" :rows="5" readonly placeholder="文本的 URL 编码将显示在这里" />

    <div class="copy-row">
      <n-button type="primary" @click="copyEncodedResult">复制编码结果</n-button>
    </div>
  </n-card>

  <n-card title="URL 编码转文本" class="card">
    <ToolTextarea v-model:value="urlInput" label="待解码 URL" :rows="5" placeholder="在此输入 URL 编码字符串…" />

    <ToolTextarea v-model:value="textOutput" label="解码结果" :rows="5" readonly placeholder="解码后的文本将显示在这里" />

    <div class="copy-row">
      <n-button type="primary" @click="copyDecodedResult">复制解码结果</n-button>
    </div>
  </n-card>
</template>

<style scoped>
.card {
  min-width: 400px;
}

.copy-row {
  display: flex;
  justify-content: center;
}
</style>