<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NButton, useMessage } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()

const inputMarkdown = ref('')
const outputHtml = ref('')

async function run() {
  try {
    const output = await RunTool('markdown-to-html', JSON.stringify({ text: inputMarkdown.value }))
    outputHtml.value = JSON.parse(output).result
  } catch (e) {
    outputHtml.value = ''
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch(inputMarkdown, () => debouncedRun(), { immediate: true })

const { copy } = useClipboard({ source: outputHtml })

async function copyResult() {
  await copy()
  message.success('HTML 已复制到剪贴板')
}
</script>

<template>
  <div class="markdownhtml-tool">
    <n-card title="Markdown 转 HTML" class="card">
      <ToolTextarea v-model:value="inputMarkdown" label="输入 Markdown" :rows="10" placeholder="在此粘贴 Markdown 内容…" />

      <ToolTextarea v-model:value="outputHtml" label="输出 HTML" :rows="10" readonly placeholder="转换得到的 HTML 将显示在这里" />

      <div class="copy-row">
        <n-button type="primary" @click="copyResult">复制 HTML</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.markdownhtml-tool {
  width: 100%;
}

.card {
  width: 100%;
}

.copy-row {
  display: flex;
  justify-content: center;
}
</style>
