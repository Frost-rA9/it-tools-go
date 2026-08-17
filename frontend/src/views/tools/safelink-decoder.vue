<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NButton, useMessage } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()

const inputUrl = ref('')
const decodedUrl = ref('')

async function runDecode() {
  try {
    const output = await RunTool('safelink-decoder', JSON.stringify({ url: inputUrl.value }))
    decodedUrl.value = JSON.parse(output).decoded_url
  } catch (e) {
    decodedUrl.value = ''
    message.error(String(e))
  }
}

const debouncedDecode = useDebounceFn(runDecode, 150)
watch(inputUrl, () => debouncedDecode(), { immediate: true })

const { copy } = useClipboard()

async function copyDecoded() {
  await copy(decodedUrl.value)
  message.success('解码 URL 已复制到剪贴板')
}
</script>

<template>
  <n-card class="card">
    <ToolTextarea v-model:value="inputUrl" label="输入 Outlook SafeLink URL" :rows="4" placeholder="在此粘贴 Outlook SafeLink 链接…" />

    <ToolTextarea v-model:value="decodedUrl" label="解码后的 URL" :rows="4" readonly monospace placeholder="解码后的 URL 将显示在这里" />

    <div class="copy-row">
      <n-button type="primary" :disabled="!decodedUrl" @click="copyDecoded">复制解码 URL</n-button>
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
  margin-top: 16px;
}
</style>