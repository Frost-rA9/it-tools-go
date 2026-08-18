<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NAlert } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import CodeOutput from '../../components/CodeOutput.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const inputText = ref('')
const outputText = ref('')
const errorMessage = ref('')

async function run() {
  try {
    const output = await RunTool('json-to-xml', JSON.stringify({ text: inputText.value }))
    outputText.value = JSON.parse(output).result
    errorMessage.value = ''
  } catch (e) {
    outputText.value = ''
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch(inputText, () => debouncedRun(), { immediate: true })

</script>

<template>
      <n-card title="JSON 转 XML — 输入" class="tool-card">
      <ToolTextarea v-model:value="inputText" label="输入 JSON" :rows="20" placeholder="在此粘贴 JSON…" />

      <n-alert v-if="errorMessage" type="error" class="error-alert">
        {{ errorMessage }}
      </n-alert>
    </n-card>

    <n-card title="JSON 转 XML — 输出" class="tool-card">
      <CodeOutput label="XML 输出" :value="outputText" language="xml" :rows="20" />
    </n-card>
</template>

<style scoped>
.error-alert {
  margin-bottom: 16px;
}

</style>
