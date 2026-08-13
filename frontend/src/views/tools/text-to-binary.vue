<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const textInput = ref('')
const binaryOutput = ref('')

const binaryInput = ref('')
const textOutput = ref('')

async function runTextToBinary() {
  try {
    const output = await RunTool(
      'text-to-binary',
      JSON.stringify({ text: textInput.value, mode: 'text_to_binary' }),
    )
    binaryOutput.value = JSON.parse(output).result
  } catch (e) {
    binaryOutput.value = ''
    message.error(String(e))
  }
}

async function runBinaryToText() {
  try {
    const output = await RunTool(
      'text-to-binary',
      JSON.stringify({ text: binaryInput.value, mode: 'binary_to_text' }),
    )
    textOutput.value = JSON.parse(output).result
  } catch (e) {
    textOutput.value = ''
  }
}

const debouncedText = useDebounceFn(runTextToBinary, 150)
const debouncedBinary = useDebounceFn(runBinaryToText, 150)

watch(textInput, () => debouncedText(), { immediate: true })
watch(binaryInput, () => debouncedBinary(), { immediate: true })

const { copy: copyBinary } = useClipboard({ source: binaryOutput })
const { copy: copyText } = useClipboard({ source: textOutput })

async function copyBinaryResult() {
  await copyBinary()
  message.success('二进制已复制到剪贴板')
}

async function copyTextResult() {
  await copyText()
  message.success('文本已复制到剪贴板')
}
</script>

<template>
  <div class="textbinary-tool">
    <n-card title="文本转 ASCII 二进制" class="card">
      <div class="field">
        <div class="field-label">输入文本</div>
        <n-input
          v-model:value="textInput"
          type="textarea"
          :rows="4"
          placeholder="在此输入文本，如 Hello world…"
        />
      </div>

      <div class="field">
        <div class="field-label">二进制结果</div>
        <n-input
          :value="binaryOutput"
          type="textarea"
          :rows="4"
          readonly
          placeholder="文本的二进制表示将显示在这里"
        />
      </div>

      <div class="copy-row">
        <n-button type="primary" @click="copyBinaryResult">复制二进制</n-button>
      </div>
    </n-card>

    <n-card title="ASCII 二进制转文本" class="card">
      <div class="field">
        <div class="field-label">输入二进制</div>
        <n-input
          v-model:value="binaryInput"
          type="textarea"
          :rows="4"
          placeholder="在此输入二进制，如 01001000 01101001…"
        />
      </div>

      <div class="field">
        <div class="field-label">文本结果</div>
        <n-input
          :value="textOutput"
          type="textarea"
          :rows="4"
          readonly
          placeholder="二进制对应的文本将显示在这里"
        />
      </div>

      <div class="copy-row">
        <n-button type="primary" @click="copyTextResult">复制文本</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.textbinary-tool {
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
