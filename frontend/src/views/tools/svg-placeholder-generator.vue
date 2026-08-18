<script setup lang="ts">
import { ref, watch } from 'vue'
import { NButton, NCard, NColorPicker, NForm, NFormItem, NInput, NInputNumber, NSwitch, useMessage } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const width = ref(600)
const height = ref(350)
const fontSize = ref(26)
const background = ref('#cccccc')
const foreground = ref('#333333')
const exactSize = ref(true)
const text = ref('')
const svg = ref('')
const base64 = ref('')

async function generate() {
  try {
    const output = await RunTool('svg-placeholder-generator', JSON.stringify({
      width: width.value,
      height: height.value,
      font_size: fontSize.value,
      background: background.value,
      foreground: foreground.value,
      exact_size: exactSize.value,
      text: text.value,
    }))
    const result = JSON.parse(output)
    svg.value = result.svg
    base64.value = result.base64
  } catch (e) {
    svg.value = ''
    base64.value = ''
    message.error(String(e))
  }
}

const debouncedGenerate = useDebounceFn(generate, 150)
watch([width, height, fontSize, background, foreground, exactSize, text], () => debouncedGenerate(), { immediate: true })

const { copy } = useClipboard()

async function copyValue(value: string, label: string) {
  await copy(value)
  message.success(`${label} 已复制到剪贴板`)
}

function download() {
  if (!base64.value) return
  const link = document.createElement('a')
  link.href = base64.value
  link.download = 'placeholder.svg'
  link.click()
}
</script>

<template>
  <n-card title="参数与输出" class="tool-card">
    <n-form label-placement="left" label-width="110">
      <div class="form-row">
        <n-form-item label="宽度 (px)">
          <n-input-number v-model:value="width" :min="1" />
        </n-form-item>
        <n-form-item label="背景色">
          <n-color-picker v-model:value="background" :modes="['hex']" />
        </n-form-item>
      </div>
      <div class="form-row">
        <n-form-item label="高度 (px)">
          <n-input-number v-model:value="height" :min="1" />
        </n-form-item>
        <n-form-item label="文字颜色">
          <n-color-picker v-model:value="foreground" :modes="['hex']" />
        </n-form-item>
      </div>
      <div class="form-row">
        <n-form-item label="字体大小">
          <n-input-number v-model:value="fontSize" :min="1" />
        </n-form-item>
        <n-form-item label="自定义文本">
          <n-input v-model:value="text" :placeholder="`默认显示 ${width}x${height}`" />
        </n-form-item>
      </div>
      <n-form-item label="使用精确尺寸">
        <n-switch v-model:value="exactSize" />
      </n-form-item>
    </n-form>

    <ToolTextarea v-model:value="svg" label="SVG HTML 元素" :rows="8" readonly monospace />
    <ToolTextarea v-model:value="base64" label="SVG Base64" :rows="4" readonly monospace />

    <div class="actions">
      <n-button @click="copyValue(svg, 'SVG')">复制 SVG</n-button>
      <n-button @click="copyValue(base64, 'Base64')">复制 Base64</n-button>
      <n-button type="primary" @click="download">下载 SVG</n-button>
    </div>
  </n-card>

  <n-card title="预览" class="tool-card">
    <div class="preview">
      <img v-if="base64" :src="base64" alt="SVG placeholder preview" />
    </div>
  </n-card>
</template>

<style scoped>
</style>
