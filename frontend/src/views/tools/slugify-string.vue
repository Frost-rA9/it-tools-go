<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NButton, useMessage } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()

const input = ref('')
const slug = ref('')

async function runSlugify() {
  try {
    const output = await RunTool('slugify-string', JSON.stringify({ text: input.value }))
    slug.value = JSON.parse(output).slug
  } catch (e) {
    slug.value = ''
    message.error(String(e))
  }
}

const debouncedSlugify = useDebounceFn(runSlugify, 150)
watch(input, () => debouncedSlugify(), { immediate: true })

const { copy } = useClipboard()

async function copySlug() {
  await copy(slug.value)
  message.success('Slug 已复制到剪贴板')
}
</script>

<template>
  <n-card class="card">
    <ToolTextarea v-model:value="input" label="要 Slug 化的字符串" :rows="4" placeholder="在此输入字符串… (如 My file path)" />

    <ToolTextarea v-model:value="slug" label="Slug" :rows="4" readonly monospace placeholder="生成的 slug 将显示在这里 (如 my-file-path)" />

    <div class="copy-row">
      <n-button type="primary" :disabled="!slug" @click="copySlug">复制 Slug</n-button>
    </div>
  </n-card>
</template>

<style scoped>
.card {
  min-width: 400px;
}

.copy-row {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}
</style>