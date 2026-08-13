<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

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
  <div class="textunicode-tool">
    <n-card title="文本转 Unicode" class="card">
      <div class="field">
        <div class="field-label">输入文本</div>
        <n-input
          v-model:value="textInput"
          type="textarea"
          :rows="4"
          placeholder="在此输入文本，如 Hello…"
        />
      </div>

      <div class="field">
        <div class="field-label">Unicode 结果</div>
        <n-input
          :value="unicodeOutput"
          type="textarea"
          :rows="4"
          readonly
          placeholder="文本的 Unicode 表示将显示在这里"
        />
      </div>

      <div class="copy-row">
        <n-button type="primary" @click="copyUnicodeResult">复制 Unicode</n-button>
      </div>
    </n-card>

    <n-card title="Unicode 转文本" class="card">
      <div class="field">
        <div class="field-label">输入 Unicode</div>
        <n-input
          v-model:value="unicodeInput"
          type="textarea"
          :rows="4"
          placeholder="在此输入 Unicode，如 &#72;&#105;…"
        />
      </div>

      <div class="field">
        <div class="field-label">文本结果</div>
        <n-input
          :value="textOutput"
          type="textarea"
          :rows="4"
          readonly
          placeholder="Unicode 对应的文本将显示在这里"
        />
      </div>

      <div class="copy-row">
        <n-button type="primary" @click="copyTextResult">复制文本</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.textunicode-tool {
  display: flex;
  flex-direction: column;
  gap: 16px;
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
