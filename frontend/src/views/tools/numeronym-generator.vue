<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NButton, useMessage, useThemeVars } from 'naive-ui'
import { ArrowDown } from '@vicons/tabler'
import { useClipboard, useDebounceFn } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const word = ref('internationalization')
const numeronym = ref('')

async function run() {
  try {
    const output = await RunTool('numeronym-generator', JSON.stringify({ word: word.value }))
    numeronym.value = JSON.parse(output).numeronym
  } catch (e) {
    numeronym.value = ''
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch(word, () => debouncedRun(), { immediate: true })

const { copy } = useClipboard()

async function copyResult() {
  if (!numeronym.value) return
  await copy(numeronym.value)
  message.success('已复制到剪贴板')
}
</script>

<template>
  <n-card title="数字名称生成器" class="tool-card">
    <div class="word-input">
      <n-input
        v-model:value="word"
        size="large"
        clearable
        placeholder="输入单词，如 internationalization…"
      />
    </div>

    <div class="arrow">
      <n-icon :component="ArrowDown" :size="30" />
    </div>

    <div class="result-row">
      <n-input :value="numeronym" size="large" readonly placeholder="缩写将显示在这里，如 i18n…" class="result-input" />
      <n-button size="large" @click="copyResult" :disabled="!numeronym">复制</n-button>
    </div>
  </n-card>
</template>

<style scoped>
.word-input {
  max-width: 520px;
  margin: 0 auto;
}

.arrow {
  display: flex;
  justify-content: center;
  margin: 12px 0;
  color: v-bind('themeVars.textColor3');
}

.result-row {
  display: flex;
  gap: 10px;
  max-width: 520px;
  margin: 0 auto;
}

.result-input {
  flex: 1;
}
</style>