<script setup lang="ts">
import { ref, watch } from 'vue'
import { NButton, NCard, NColorPicker, NForm, NFormItem, NGrid, NGi, NImage, NInput, NSelect, useMessage } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const text = ref('https://example.com')
const foreground = ref('#000000ff')
const background = ref('#ffffffff')
const level = ref('medium')
const dataURL = ref('')

const levels = ['low', 'medium', 'quartile', 'high'].map((value) => ({ label: value, value }))

async function generate() {
  if (!text.value) {
    dataURL.value = ''
    return
  }
  try {
    const output = await RunTool('qr-code-generator', JSON.stringify({
      text: text.value,
      foreground: foreground.value,
      background: background.value,
      level: level.value,
    }))
    dataURL.value = JSON.parse(output).data_url
  } catch (e) {
    dataURL.value = ''
    message.error(String(e))
  }
}

const debouncedGenerate = useDebounceFn(generate, 150)
watch([text, foreground, background, level], () => debouncedGenerate(), { immediate: true })

function download() {
  if (!dataURL.value) return
  const link = document.createElement('a')
  link.href = dataURL.value
  link.download = 'qr-code.png'
  link.click()
}
</script>

<template>
  <n-card class="card">
    <n-grid x-gap="16" y-gap="16" cols="1 600:3">
      <n-gi span="2">
        <n-form label-placement="left" label-width="130" label-align="right" class="form">
          <n-form-item label="文本">
            <n-input v-model:value="text" placeholder="输入链接或文本…" />
          </n-form-item>
          <n-form-item label="前景色">
            <n-color-picker v-model:value="foreground" :modes="['hex']" />
          </n-form-item>
          <n-form-item label="背景色">
            <n-color-picker v-model:value="background" :modes="['hex']" />
          </n-form-item>
          <n-form-item label="纠错级别">
            <n-select v-model:value="level" :options="levels" />
          </n-form-item>
        </n-form>
      </n-gi>
      <n-gi>
        <div class="preview">
          <n-image v-if="dataURL" :src="dataURL" width="200" />
          <n-button :disabled="!dataURL" @click="download">下载二维码</n-button>
        </div>
      </n-gi>
    </n-grid>
  </n-card>
</template>

<style scoped>
.card { min-width: 400px; }
.preview { display: flex; flex-direction: column; align-items: center; gap: 12px; }
</style>
