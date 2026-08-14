<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

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
      <div class="field">
        <div class="field-label">输入 Markdown</div>
        <n-input
          v-model:value="inputMarkdown"
          type="textarea"
          :rows="10"
          placeholder="在此粘贴 Markdown 内容…"
        />
      </div>

      <div class="field">
        <div class="field-label">输出 HTML</div>
        <n-input
          :value="outputHtml"
          type="textarea"
          :rows="10"
          readonly
          placeholder="转换得到的 HTML 将显示在这里"
        />
      </div>

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

.field {
  margin-bottom: 16px;
}

.field-label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.copy-row {
  display: flex;
  justify-content: center;
}
</style>
