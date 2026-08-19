<script setup lang="ts">
import { ref, watch } from 'vue'
import { NButton, NCard, NForm, NFormItem, NSlider, NSwitch, useMessage } from 'naive-ui'
import { useClipboard, useDebounceFn } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()

const paragraphs = ref(1)
const sentences = ref<[number, number]>([3, 8])
const words = ref<[number, number]>([8, 15])
const startWithLoremIpsum = ref(true)
const asHTML = ref(false)
const outputText = ref('')

async function run() {
  try {
    const output = await RunTool('lorem-ipsum-generator', JSON.stringify({
      paragraphs: paragraphs.value,
      sentence_min: sentences.value[0],
      sentence_max: sentences.value[1],
      word_min: words.value[0],
      word_max: words.value[1],
      start_with_lorem_ipsum: startWithLoremIpsum.value,
      as_html: asHTML.value,
    }))
    outputText.value = JSON.parse(output).text
  } catch (e) {
    outputText.value = ''
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch(
  [paragraphs, sentences, words, startWithLoremIpsum, asHTML],
  () => debouncedRun(),
  { immediate: true, deep: true },
)

const { copy } = useClipboard()

async function copyResult() {
  if (!outputText.value) return
  await copy(outputText.value)
  message.success('已复制到剪贴板')
}
</script>

<template>
  <n-card title="Lorem ipsum 生成器" class="tool-card">
    <n-form label-placement="left" label-width="120">
      <n-form-item label="段落数">
        <n-slider v-model:value="paragraphs" :min="1" :max="20" :step="1" />
      </n-form-item>
      <n-form-item label="每段句数">
        <n-slider v-model:value="sentences" :min="1" :max="50" :step="1" range />
      </n-form-item>
      <n-form-item label="每句词数">
        <n-slider v-model:value="words" :min="1" :max="50" :step="1" range />
      </n-form-item>
      <n-form-item label="以 Lorem ipsum 开头">
        <n-switch v-model:value="startWithLoremIpsum" />
      </n-form-item>
      <n-form-item label="输出 HTML">
        <n-switch v-model:value="asHTML" />
      </n-form-item>
    </n-form>

    <ToolTextarea
      v-model:value="outputText"
      label="生成结果"
      :rows="10"
      readonly
      monospace
      placeholder="调整参数或点击刷新生成占位文本…"
    />

    <div class="actions">
      <n-button @click="copyResult">复制</n-button>
      <n-button type="primary" @click="run">刷新</n-button>
    </div>
  </n-card>
</template>

<style scoped>
.actions {
  display: flex;
  justify-content: center;
  gap: 12px;
}
</style>