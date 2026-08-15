<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NInputNumber, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const arabicInput = ref<number | null>(42)
const romanOutput = ref('')

const romanInput = ref('XLII')
const arabicOutput = ref('')

async function runArabicToRoman() {
  try {
    const output = await RunTool(
      'roman-numeral-converter',
      JSON.stringify({ value: arabicInput.value == null ? '' : String(arabicInput.value), mode: 'arabic_to_roman' }),
    )
    romanOutput.value = JSON.parse(output).result
  } catch (e) {
    romanOutput.value = ''
    message.error(String(e))
  }
}

async function runRomanToArabic() {
  try {
    const output = await RunTool(
      'roman-numeral-converter',
      JSON.stringify({ value: romanInput.value, mode: 'roman_to_arabic' }),
    )
    arabicOutput.value = JSON.parse(output).result
  } catch (e) {
    arabicOutput.value = ''
  }
}

const debouncedArabic = useDebounceFn(runArabicToRoman, 150)
const debouncedRoman = useDebounceFn(runRomanToArabic, 150)

watch(arabicInput, () => debouncedArabic(), { immediate: true })
watch(romanInput, () => debouncedRoman(), { immediate: true })

const { copy: copyRoman } = useClipboard({ source: romanOutput })
const { copy: copyArabic } = useClipboard({ source: arabicOutput })

async function copyRomanResult() {
  await copyRoman()
  message.success('罗马数字已复制到剪贴板')
}

async function copyArabicResult() {
  await copyArabic()
  message.success('阿拉伯数字已复制到剪贴板')
}
</script>

<template>
  
    <n-card title="阿拉伯数字转罗马数字" class="card">
      <div class="field">
        <div class="field-label">阿拉伯数字（1 - 3999）</div>
        <n-input-number v-model:value="arabicInput" :min="1" :show-button="false" style="width: 100%" />
      </div>

      <div class="field">
        <div class="field-label">罗马数字结果</div>
        <n-input :value="romanOutput" readonly placeholder="转换结果将显示在这里" />
      </div>

      <div class="copy-row">
        <n-button type="primary" @click="copyRomanResult">复制罗马数字</n-button>
      </div>
    </n-card>

    <n-card title="罗马数字转阿拉伯数字" class="card">
      <div class="field">
        <div class="field-label">罗马数字</div>
        <n-input v-model:value="romanInput" placeholder="在此输入罗马数字，如 XLII…" />
      </div>

      <div class="field">
        <div class="field-label">阿拉伯数字结果</div>
        <n-input :value="arabicOutput" readonly placeholder="转换结果将显示在这里" />
      </div>

      <div class="copy-row">
        <n-button type="primary" @click="copyArabicResult">复制阿拉伯数字</n-button>
      </div>
    </n-card>
</template>

<style scoped>
.card {
  min-width: 360px;
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
