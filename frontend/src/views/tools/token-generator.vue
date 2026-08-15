<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NSwitch, NSlider, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'
import ToolTextarea from '../../components/ToolTextarea.vue'

const message = useMessage()
const themeVars = useThemeVars()

const length = ref(64)
const withUppercase = ref(true)
const withLowercase = ref(true)
const withNumbers = ref(true)
const withSymbols = ref(false)
const token = ref('')

async function regenerate() {
  try {
    const output = await RunTool(
      'token-generator',
      JSON.stringify({
        length: length.value,
        with_uppercase: withUppercase.value,
        with_lowercase: withLowercase.value,
        with_numbers: withNumbers.value,
        with_symbols: withSymbols.value,
      }),
    )
    token.value = JSON.parse(output).result
  } catch (e) {
    token.value = ''
    message.error(String(e))
  }
}

const debouncedRegenerate = useDebounceFn(regenerate, 150)

watch([length, withUppercase, withLowercase, withNumbers, withSymbols], () => debouncedRegenerate(), { immediate: true })

const { copy } = useClipboard({ source: token })

async function copyToken() {
  await copy()
  message.success('Token 已复制到剪贴板')
}
</script>

<template>
  <div class="token-tool">
    <n-card title="Token 生成器" class="card">
      <div class="switch-grid">
        <div class="switch-col">
          <div class="switch-row">
            <span class="switch-label">大写字母</span>
            <n-switch v-model:value="withUppercase" />
          </div>
          <div class="switch-row">
            <span class="switch-label">小写字母</span>
            <n-switch v-model:value="withLowercase" />
          </div>
        </div>
        <div class="switch-col">
          <div class="switch-row">
            <span class="switch-label">数字</span>
            <n-switch v-model:value="withNumbers" />
          </div>
          <div class="switch-row">
            <span class="switch-label">符号</span>
            <n-switch v-model:value="withSymbols" />
          </div>
        </div>
      </div>

      <div class="slider-row">
        <div class="field-label">长度（{{ length }}）</div>
        <n-slider v-model:value="length" :min="1" :max="512" :step="1" />
      </div>

      <ToolTextarea
        v-model:value="token"
        readonly
        placeholder="生成的 Token 将显示在这里"
        class="token-display"
      />

      <div class="btn-row">
        <n-button type="primary" @click="copyToken">复制 Token</n-button>
        <n-button @click="regenerate">重新生成</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.token-tool {
  display: flex;
  width: 100%;
}

.card {
  width: 100%;
}

.switch-grid {
  display: flex;
  gap: 48px;
  justify-content: center;
  margin-bottom: 20px;
}

.switch-col {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.switch-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  width: 140px;
}

.switch-label {
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
}

.slider-row {
  margin-bottom: 16px;
}

.btn-row {
  display: flex;
  justify-content: center;
  gap: 12px;
}

.token-display :deep(textarea) {
  text-align: center;
}
</style>
