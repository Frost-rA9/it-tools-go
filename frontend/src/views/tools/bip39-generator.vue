<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NGrid, NGi, NInput, NSelect, NButton, NIcon, useMessage, useThemeVars } from 'naive-ui'
import { Refresh, Copy } from '@vicons/tabler'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const language = ref('English')
const entropy = ref('')
const mnemonic = ref('')
const entropyError = ref('')
const mnemonicError = ref('')

const languageOptions = [
  { label: 'English', value: 'English' },
  { label: 'Chinese simplified', value: 'Chinese simplified' },
  { label: 'Chinese traditional', value: 'Chinese traditional' },
  { label: 'Czech', value: 'Czech' },
  { label: 'French', value: 'French' },
  { label: 'Italian', value: 'Italian' },
  { label: 'Japanese', value: 'Japanese' },
  { label: 'Korean', value: 'Korean' },
  { label: 'Portuguese', value: 'Portuguese' },
  { label: 'Spanish', value: 'Spanish' },
]

let updating = false

async function run(mode: 'entropy-to-mnemonic' | 'mnemonic-to-entropy', text: string, target: 'mnemonic' | 'entropy') {
  try {
    const output = await RunTool('bip39-generator', JSON.stringify({ mode, language: language.value, text }))
    const value = JSON.parse(output).result
    updating = true
    if (target === 'mnemonic') {
      mnemonic.value = value
      mnemonicError.value = ''
    } else {
      entropy.value = value
      entropyError.value = ''
    }
    setTimeout(() => { updating = false }, 0)
  } catch (e) {
    if (target === 'mnemonic') {
      mnemonic.value = ''
      mnemonicError.value = String(e)
    } else {
      entropy.value = ''
      entropyError.value = String(e)
    }
  }
}

const debouncedToMnemonic = useDebounceFn((t: string) => run('entropy-to-mnemonic', t, 'mnemonic'), 150)
const debouncedToEntropy = useDebounceFn((t: string) => run('mnemonic-to-entropy', t, 'entropy'), 150)

watch(entropy, (v) => {
  if (updating) return
  entropyError.value = ''
  if (!v.trim()) { mnemonic.value = ''; return }
  debouncedToMnemonic(v)
}, { immediate: true })

watch(mnemonic, (v) => {
  if (updating) return
  mnemonicError.value = ''
  if (!v.trim()) { entropy.value = ''; return }
  debouncedToEntropy(v)
})

watch(language, () => {
  if (entropy.value.trim()) debouncedToMnemonic(entropy.value)
})

async function refreshEntropy() {
  try {
    const output = await RunTool('bip39-generator', JSON.stringify({ mode: 'generate', language: language.value }))
    entropy.value = JSON.parse(output).result
    entropyError.value = ''
  } catch (e) {
    message.error(String(e))
  }
}

const entropyCopySource = ref('')
const mnemonicCopySource = ref('')
const { copy: copyEntropy } = useClipboard({ source: entropyCopySource })
const { copy: copyMnemonic } = useClipboard({ source: mnemonicCopySource })

async function copyEntropyResult() {
  entropyCopySource.value = entropy.value
  await copyEntropy()
  message.success('熵已复制到剪贴板')
}

async function copyMnemonicResult() {
  mnemonicCopySource.value = mnemonic.value
  await copyMnemonic()
  message.success('助记词已复制到剪贴板')
}
</script>

<template>
  <div class="bip39-tool">
    <n-card title="BIP39 密码生成器" class="card">
      <n-grid cols="3" :x-gap="16" responsive="screen">
        <n-gi span="1">
          <div class="field-label">语言：</div>
          <n-select v-model:value="language" :options="languageOptions" filterable />
        </n-gi>
        <n-gi span="2">
          <div class="field-label">熵（seed）：</div>
          <div class="input-row">
            <n-input v-model:value="entropy" class="mono-input" placeholder="十六进制熵，如 32 位…" />
            <n-button quaternary circle @click="refreshEntropy" title="随机生成">
              <n-icon :component="Refresh" :size="18" />
            </n-button>
            <n-button quaternary circle @click="copyEntropyResult" title="复制">
              <n-icon :component="Copy" :size="18" />
            </n-button>
          </div>
          <div v-if="entropyError" class="feedback">{{ entropyError }}</div>
        </n-gi>
      </n-grid>

      <div class="mnemonic-row">
        <div class="field-label">助记词：</div>
        <div class="input-row">
          <n-input v-model:value="mnemonic" class="mono-input" placeholder="生成的助记词将显示在这里" />
          <n-button quaternary circle @click="copyMnemonicResult" title="复制">
            <n-icon :component="Copy" :size="18" />
          </n-button>
        </div>
        <div v-if="mnemonicError" class="feedback">{{ mnemonicError }}</div>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.bip39-tool {
  width: 100%;
}

.card {
  width: 100%;
}

.field-label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.mnemonic-row {
  margin-top: 16px;
}

.input-row {
  display: flex;
  align-items: center;
  gap: 4px;
}

.mono-input {
  flex: 1;
}

.mono-input :deep(.n-input__input-el) {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}

.feedback {
  font-size: 12px;
  color: v-bind('themeVars.errorColor');
  margin-top: 4px;
}
</style>
