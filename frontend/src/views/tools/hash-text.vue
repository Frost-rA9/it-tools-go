<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NSelect, NButton, NInput, NDivider, NIcon, useMessage, useThemeVars } from 'naive-ui'
import { Copy } from '@vicons/tabler'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'
import ToolTextarea from '../../components/ToolTextarea.vue'

interface DigestResult {
  algo: string
  digest: string
}

const message = useMessage()
const themeVars = useThemeVars()

const textInput = ref('')
const encoding = ref('Hex')

const encodingOptions = [
  { label: 'Binary (base 2)', value: 'Bin' },
  { label: 'Hexadecimal (base 16)', value: 'Hex' },
  { label: 'Base64 (base 64)', value: 'Base64' },
  { label: 'Base64url (base 64 with url safe chars)', value: 'Base64url' },
]

const results = ref<DigestResult[]>([])

async function run() {
  try {
    const output = await RunTool('hash-text', JSON.stringify({ text: textInput.value, encoding: encoding.value }))
    results.value = JSON.parse(output).results
  } catch (e) {
    results.value = []
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch([textInput, encoding], () => debouncedRun(), { immediate: true })

const copySource = ref('')
const { copy } = useClipboard({ source: copySource })

async function copyDigest(value: string) {
  copySource.value = value
  await copy()
  message.success('摘要已复制到剪贴板')
}
</script>

<template>
  <div class="hashtext-tool">
    <n-card title="哈希文本" class="card">
      <ToolTextarea
        v-model:value="textInput"
        label="待哈希文本"
        :rows="4"
        placeholder="在此输入要哈希的文本…"
      />

      <n-divider style="margin: 20px 0" />

      <div class="field">
        <div class="field-label">摘要编码</div>
        <n-select v-model:value="encoding" :options="encodingOptions" />
      </div>

      <div class="result-list">
        <div v-for="r in results" :key="r.algo" class="algo-group">
          <span class="algo-label">{{ r.algo }}</span>
          <n-input :value="r.digest" readonly class="algo-input" placeholder="…">
            <template #suffix>
              <n-button quaternary circle size="small" @click="copyDigest(r.digest)">
                <n-icon :component="Copy" />
              </n-button>
            </template>
          </n-input>
        </div>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.hashtext-tool {
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

.result-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.algo-group {
  display: flex;
  width: 100%;
}

.algo-label {
  flex: 0 0 120px;
  display: flex;
  align-items: center;
  align-self: stretch;
  padding: 0 12px;
  background-color: v-bind('themeVars.actionColor');
  border: 1px solid v-bind('themeVars.borderColor');
  border-right: none;
  border-radius: 3px 0 0 3px;
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
  box-sizing: border-box;
}

.algo-input {
  flex: 1;
  min-width: 0;
  border-radius: 0 3px 3px 0;
}

.algo-input :deep(.n-input__border),
.algo-input :deep(.n-input__state-border) {
  border-left: none;
}
</style>
