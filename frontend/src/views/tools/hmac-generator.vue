<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NSelect, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'
import ToolTextarea from '../../components/ToolTextarea.vue'

const message = useMessage()
const themeVars = useThemeVars()

const text = ref('')
const secret = ref('')
const algo = ref('SHA256')
const encoding = ref('Hex')
const result = ref('')

const algoOptions = [
  { label: 'MD5', value: 'MD5' },
  { label: 'RIPEMD160', value: 'RIPEMD160' },
  { label: 'SHA1', value: 'SHA1' },
  { label: 'SHA3', value: 'SHA3' },
  { label: 'SHA224', value: 'SHA224' },
  { label: 'SHA256', value: 'SHA256' },
  { label: 'SHA384', value: 'SHA384' },
  { label: 'SHA512', value: 'SHA512' },
]

const encodingOptions = [
  { label: 'Binary (base 2)', value: 'Bin' },
  { label: 'Hexadecimal (base 16)', value: 'Hex' },
  { label: 'Base64 (base 64)', value: 'Base64' },
  { label: 'Base64url (base 64 with url safe chars)', value: 'Base64url' },
]

async function run() {
  try {
    const output = await RunTool(
      'hmac-generator',
      JSON.stringify({ text: text.value, secret: secret.value, algo: algo.value, encoding: encoding.value }),
    )
    result.value = JSON.parse(output).result
  } catch (e) {
    result.value = ''
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch([text, secret, algo, encoding], () => debouncedRun(), { immediate: true })

const { copy } = useClipboard({ source: result })

async function copyResult() {
  await copy()
  message.success('HMAC 已复制到剪贴板')
}
</script>

<template>
  <div class="hmac-tool">
    <n-card title="HMAC 生成器" class="card">
      <ToolTextarea v-model:value="text" label="待计算文本" :rows="4" placeholder="输入要计算 HMAC 的文本…" />

      <div class="field">
        <div class="field-label">密钥</div>
        <n-input v-model:value="secret" clearable placeholder="输入密钥…" />
      </div>

      <div class="select-row">
        <div class="select-col">
          <div class="field-label">哈希函数</div>
          <n-select v-model:value="algo" :options="algoOptions" />
        </div>
        <div class="select-col">
          <div class="field-label">输出编码</div>
          <n-select v-model:value="encoding" :options="encodingOptions" />
        </div>
      </div>

      <ToolTextarea v-model:value="result" label="HMAC 结果" readonly :rows="3" placeholder="HMAC 将显示在这里" monospace />

      <div class="copy-row">
        <n-button type="primary" @click="copyResult">复制 HMAC</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.hmac-tool {
  width: 100%;
}

.card {
  width: 100%;
}

.field {
  margin-bottom: 16px;
}

.select-row {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.select-col {
  flex: 1;
  min-width: 0;
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
