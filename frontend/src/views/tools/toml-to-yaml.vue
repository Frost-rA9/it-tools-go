<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NButton, NAlert, useMessage } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()

const inputText = ref('')
const outputText = ref('')
const errorMessage = ref('')

async function run() {
  try {
    const output = await RunTool('toml-to-yaml', JSON.stringify({ text: inputText.value }))
    outputText.value = JSON.parse(output).result
    errorMessage.value = ''
  } catch (e) {
    outputText.value = ''
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch(inputText, () => debouncedRun(), { immediate: true })

const { copy } = useClipboard({ source: outputText })

async function copyResult() {
  await copy()
  message.success('YAML 已复制到剪贴板')
}
</script>

<template>
  <div class="tomlyaml-tool">
    <n-card title="TOML 转 YAML" class="card">
      <ToolTextarea v-model:value="inputText" label="输入 TOML" :rows="10" placeholder="在此粘贴 TOML…" />

      <n-alert v-if="errorMessage" type="error" class="error-alert">
        {{ errorMessage }}
      </n-alert>

      <ToolTextarea v-model:value="outputText" label="YAML 输出" :rows="10" readonly placeholder="转换得到的 YAML 将显示在这里" />

      <div class="copy-row">
        <n-button type="primary" @click="copyResult">复制 YAML</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.tomlyaml-tool {
  width: 100%;
}

.card {
  width: 100%;
}

.error-alert {
  margin-bottom: 16px;
}

.copy-row {
  display: flex;
  justify-content: center;
}
</style>
