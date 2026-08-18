<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NInputNumber, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const input = ref('')
const cost = ref(10)
const hashed = ref('')

const compareString = ref('')
const compareHash = ref('')
const compareMatch = ref(false)

async function runHash() {
  try {
    const output = await RunTool('bcrypt', JSON.stringify({ mode: 'hash', text: input.value, cost: cost.value }))
    hashed.value = JSON.parse(output).result
  } catch (e) {
    hashed.value = ''
    message.error(String(e))
  }
}

async function runCompare() {
  try {
    const output = await RunTool('bcrypt', JSON.stringify({ mode: 'compare', text: compareString.value, hash: compareHash.value }))
    compareMatch.value = JSON.parse(output).match
  } catch (e) {
    compareMatch.value = false
  }
}

const debouncedHash = useDebounceFn(runHash, 150)
const debouncedCompare = useDebounceFn(runCompare, 150)

watch([input, cost], () => debouncedHash(), { immediate: true })
watch([compareString, compareHash], () => debouncedCompare(), { immediate: true })

const { copy } = useClipboard({ source: hashed })

async function copyHash() {
  await copy()
  message.success('BCrypt 哈希已复制到剪贴板')
}
</script>

<template>
  
    <n-card title="BCrypt 加密" class="tool-card">
      <div class="field">
        <div class="field-label">待加密字符串</div>
        <n-input v-model:value="input" placeholder="输入要加密的字符串…" />
      </div>

      <div class="field">
        <div class="field-label">Salt 轮数</div>
        <n-input-number v-model:value="cost" :min="4" :max="31" :step="1" style="width: 200px" />
      </div>

      <div class="field">
        <div class="field-label">BCrypt 哈希</div>
        <n-input :value="hashed" readonly class="hash-value" placeholder="生成的 BCrypt 哈希将显示在这里" />
      </div>

      <div class="copy-row">
        <n-button type="primary" @click="copyHash">复制哈希</n-button>
      </div>
    </n-card>

    <n-card title="对比字符串与哈希" class="tool-card">
      <div class="field">
        <div class="field-label">待对比字符串</div>
        <n-input v-model:value="compareString" placeholder="输入要对比的字符串…" />
      </div>

      <div class="field">
        <div class="field-label">BCrypt 哈希</div>
        <n-input v-model:value="compareHash" class="hash-value" placeholder="粘贴要对比的 BCrypt 哈希…" />
      </div>

      <div class="field">
        <div class="field-label">是否匹配</div>
        <span class="compare-result" :class="{ positive: compareMatch }">
          {{ compareMatch ? '匹配' : '不匹配' }}
        </span>
      </div>
    </n-card>
</template>

<style scoped>

.field {
  margin-bottom: 16px;
}

.field-label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.hash-value :deep(.n-input__input-el) {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}

.copy-row {
  display: flex;
  justify-content: center;
}

.compare-result {
  font-size: 14px;
  font-weight: 500;
  color: v-bind('themeVars.errorColor');
}

.compare-result.positive {
  color: v-bind('themeVars.successColor');
}
</style>
