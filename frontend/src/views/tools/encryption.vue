<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NSelect, NAlert, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'
import ToolTextarea from '../../components/ToolTextarea.vue'

const message = useMessage()
const themeVars = useThemeVars()

const algoOptions = [
  { label: 'AES', value: 'AES' },
  { label: 'TripleDES', value: 'TripleDES' },
  { label: 'Rabbit', value: 'Rabbit' },
  { label: 'RC4', value: 'RC4' },
]

const encryptText = ref('Lorem ipsum dolor sit amet')
const encryptSecret = ref('my secret key')
const encryptAlgo = ref('AES')
const encryptOutput = ref('')

const decryptText = ref('U2FsdGVkX1/EC3+6P5dbbkZ3e1kQ5o2yzuU0NHTjmrKnLBEwreV489Kr0DIB+uBs')
const decryptSecret = ref('my secret key')
const decryptAlgo = ref('AES')
const decryptOutput = ref('')
const decryptError = ref('')

async function runEncrypt() {
  decryptError.value = ''
  try {
    const output = await RunTool(
      'encryption',
      JSON.stringify({ mode: 'encrypt', text: encryptText.value, secret: encryptSecret.value, algo: encryptAlgo.value }),
    )
    encryptOutput.value = JSON.parse(output).result
  } catch (e) {
    encryptOutput.value = ''
    message.error(String(e))
  }
}

async function runDecrypt() {
  decryptOutput.value = ''
  decryptError.value = ''
  try {
    const output = await RunTool(
      'encryption',
      JSON.stringify({ mode: 'decrypt', text: decryptText.value, secret: decryptSecret.value, algo: decryptAlgo.value }),
    )
    decryptOutput.value = JSON.parse(output).result
  } catch (e) {
    decryptError.value = '解密失败：密钥错误或密文损坏'
  }
}

const debouncedEncrypt = useDebounceFn(runEncrypt, 150)
const debouncedDecrypt = useDebounceFn(runDecrypt, 150)

watch([encryptText, encryptSecret, encryptAlgo], () => debouncedEncrypt(), { immediate: true })
watch([decryptText, decryptSecret, decryptAlgo], () => debouncedDecrypt(), { immediate: true })
</script>

<template>
  
    <n-card title="加密" class="tool-card">
      <div class="side-grid">
        <ToolTextarea
          v-model:value="encryptText"
          label="待加密文本"
          :rows="4"
          placeholder="要加密的字符串…"
          monospace
        />
        <div class="side">
          <div class="field">
            <div class="field-label">密钥口令</div>
            <n-input v-model:value="encryptSecret" placeholder="输入密钥口令…" />
          </div>
          <div class="field">
            <div class="field-label">加密算法</div>
            <n-select v-model:value="encryptAlgo" :options="algoOptions" />
          </div>
        </div>
      </div>

      <ToolTextarea
        v-model:value="encryptOutput"
        label="加密结果"
        :rows="3"
        readonly
        placeholder="加密后的密文将显示在这里"
        monospace
      />
    </n-card>

    <n-card title="解密" class="tool-card">
      <div class="side-grid">
        <ToolTextarea
          v-model:value="decryptText"
          label="待解密密文"
          :rows="4"
          placeholder="粘贴要解密的密文…"
          monospace
        />
        <div class="side">
          <div class="field">
            <div class="field-label">密钥口令</div>
            <n-input v-model:value="decryptSecret" placeholder="输入密钥口令…" />
          </div>
          <div class="field">
            <div class="field-label">解密算法</div>
            <n-select v-model:value="decryptAlgo" :options="algoOptions" />
          </div>
        </div>
      </div>

      <n-alert v-if="decryptError" type="error" class="error-alert">{{ decryptError }}</n-alert>
      <ToolTextarea
        v-else
        v-model:value="decryptOutput"
        label="解密结果"
        :rows="3"
        readonly
        placeholder="解密后的明文将显示在这里"
        monospace
      />
    </n-card>
</template>

<style scoped>

.side-grid {
  display: flex;
  gap: 16px;
  align-items: flex-start;
}

.side-grid :deep(.tool-textarea) {
  flex: 1;
}

.side {
  flex: 0 0 260px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.field {
  margin-bottom: 12px;
}

.field-label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.error-alert {
  margin-bottom: 16px;
}
</style>
