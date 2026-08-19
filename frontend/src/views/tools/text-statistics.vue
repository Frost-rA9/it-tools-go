<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NStatistic, useMessage } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()

const text = ref('')
const characters = ref(0)
const words = ref(0)
const lines = ref(0)
const sizeText = ref('0 Bytes')

async function run() {
  try {
    const output = await RunTool('text-statistics', JSON.stringify({ text: text.value }))
    const result = JSON.parse(output)
    characters.value = result.characters
    words.value = result.words
    lines.value = result.lines
    sizeText.value = result.size_text
  } catch (e) {
    characters.value = 0
    words.value = 0
    lines.value = 0
    sizeText.value = ''
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch(text, () => debouncedRun(), { immediate: true })
</script>

<template>
  <n-card title="文本统计" class="tool-card">
    <ToolTextarea
      v-model:value="text"
      label="输入文本"
      :rows="8"
      placeholder="在此输入文本，实时统计字符、单词、行数与字节大小…"
    />

    <div class="stats-row">
      <n-statistic label="字符数" :value="characters" />
      <n-statistic label="单词数" :value="words" />
      <n-statistic label="行数" :value="lines" />
      <n-statistic label="字节大小" :value="sizeText" />
    </div>
  </n-card>
</template>

<style scoped>
.stats-row {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.stats-row :deep(.n-statistic) {
  flex: 1;
  min-width: 120px;
}
</style>