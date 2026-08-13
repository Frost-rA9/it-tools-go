<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NSwitch, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const encodeUrlSafe = ref(false)
const textInput = ref('')
const base64Output = ref('')

const decodeUrlSafe = ref(false)
const base64Input = ref('')
const textOutput = ref('')

async function runEncode() {
  try {
    const output = await RunTool(
      'base64-string-converter',
      JSON.stringify({ text: textInput.value, action: 'encode', url_safe: encodeUrlSafe.value }),
    )
    base64Output.value = JSON.parse(output).result
  } catch (e) {
    base64Output.value = ''
    message.error(String(e))
  }
}

async function runDecode() {
  try {
    const output = await RunTool(
      'base64-string-converter',
      JSON.stringify({ text: base64Input.value, action: 'decode', url_safe: decodeUrlSafe.value }),
    )
    textOutput.value = JSON.parse(output).result
  } catch (e) {
    textOutput.value = ''
  }
}

const debouncedEncode = useDebounceFn(runEncode, 150)
const debouncedDecode = useDebounceFn(runDecode, 150)

watch([textInput, encodeUrlSafe], () => debouncedEncode(), { immediate: true })
watch([base64Input, decodeUrlSafe], () => debouncedDecode(), { immediate: true })

const { copy: copyBase64 } = useClipboard({ source: base64Output })
const { copy: copyText } = useClipboard({ source: textOutput })

async function copyBase64Result() {
  await copyBase64()
  message.success('Base64 字符串已复制到剪贴板')
}

async function copyTextResult() {
  await copyText()
  message.success('文本已复制到剪贴板')
}
</script>

<template>
  <div class="base64-tool">
    <n-card title="文本转 Base64" class="card">
      <div class="switch-row">
        <span class="switch-label">URL 安全编码</span>
        <n-switch v-model:value="encodeUrlSafe" />
      </div>

      <div class="field">
        <div class="field-label">待编码文本</div>
        <n-input
          v-model:value="textInput"
          type="textarea"
          :rows="5"
          placeholder="在此输入文本…"
        />
      </div>

      <div class="field">
        <div class="field-label">Base64 编码结果</div>
        <n-input
          :value="base64Output"
          type="textarea"
          :rows="5"
          readonly
          placeholder="文本的 Base64 编码将显示在这里"
        />
      </div>

      <div class="copy-row">
        <n-button type="primary" @click="copyBase64Result">复制 Base64</n-button>
      </div>
    </n-card>

    <n-card title="Base64 转文本" class="card">
      <div class="switch-row">
        <span class="switch-label">URL 安全解码</span>
        <n-switch v-model:value="decodeUrlSafe" />
      </div>

      <div class="field">
        <div class="field-label">待解码 Base64</div>
        <n-input
          v-model:value="base64Input"
          type="textarea"
          :rows="5"
          placeholder="在此输入 Base64 字符串…"
        />
      </div>

      <div class="field">
        <div class="field-label">解码结果</div>
        <n-input
          :value="textOutput"
          type="textarea"
          :rows="5"
          readonly
          placeholder="解码后的文本将显示在这里"
        />
      </div>

      <div class="copy-row">
        <n-button type="primary" @click="copyTextResult">复制文本</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.base64-tool {
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: 100%;
}

.card {
  width: 100%;
}

.switch-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.switch-label {
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
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
