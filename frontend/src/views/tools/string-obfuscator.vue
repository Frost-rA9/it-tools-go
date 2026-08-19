<script setup lang="ts">
import { ref, watch } from 'vue'
import { NButton, NCard, NForm, NFormItem, NInputNumber, NSwitch, useMessage } from 'naive-ui'
import { useClipboard, useDebounceFn } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()

const text = ref('Lorem ipsum dolor sit amet')
const keepFirst = ref(4)
const keepLast = ref(0)
const keepSpace = ref(true)
const result = ref('')

async function run() {
  try {
    const output = await RunTool('string-obfuscator', JSON.stringify({
      text: text.value,
      keep_first: keepFirst.value,
      keep_last: keepLast.value,
      keep_space: keepSpace.value,
    }))
    result.value = JSON.parse(output).result
  } catch (e) {
    result.value = ''
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch([text, keepFirst, keepLast, keepSpace], () => debouncedRun(), { immediate: true })

const { copy } = useClipboard()

async function copyResult() {
  if (!result.value) return
  await copy(result.value)
  message.success('已复制到剪贴板')
}
</script>

<template>
  <n-card title="字符串混淆器" class="tool-card">
    <ToolTextarea
      v-model:value="text"
      label="要混淆的字符串"
      :rows="3"
      placeholder="在此输入要混淆的字符串…"
    />

    <n-form label-placement="left" label-width="120">
      <n-form-item label="保留开头字符">
        <n-input-number v-model:value="keepFirst" :min="0" :max="100" />
      </n-form-item>
      <n-form-item label="保留结尾字符">
        <n-input-number v-model:value="keepLast" :min="0" :max="100" />
      </n-form-item>
      <n-form-item label="保留空格">
        <n-switch v-model:value="keepSpace" />
      </n-form-item>
    </n-form>

    <ToolTextarea v-model:value="result" label="混淆结果" :rows="3" readonly monospace placeholder="混淆后的字符串将显示在这里…" />

    <div class="actions">
      <n-button @click="copyResult" :disabled="!result">复制</n-button>
    </div>
  </n-card>
</template>

<style scoped>
.actions {
  display: flex;
  justify-content: center;
  gap: 12px;
}
</style>