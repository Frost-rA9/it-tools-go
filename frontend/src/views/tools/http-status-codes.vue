<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NText, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

interface StatusCode {
  code: number
  name: string
  description: string
  type: string
}

interface StatusCategory {
  category: string
  codes: StatusCode[]
}

const themeVars = useThemeVars()

const search = ref('')
const categories = ref<StatusCategory[]>([])

async function runSearch() {
  try {
    const output = await RunTool('http-status-codes', JSON.stringify({ query: search.value }))
    categories.value = JSON.parse(output).results
  } catch (e) {
    categories.value = []
  }
}

const debouncedSearch = useDebounceFn(runSearch, 150)
watch(search, () => debouncedSearch(), { immediate: true })
</script>

<template>
  <div class="page">
    <n-input v-model:value="search" placeholder="搜索 HTTP 状态码…" clearable />

    <div v-for="cat in categories" :key="cat.category" class="category">
      <div class="category-title">{{ cat.category }}</div>

      <n-card v-for="code in cat.codes" :key="code.code" class="code-card">
        <div class="code-header">
          <span class="code-number">{{ code.code }}</span>
          <span class="code-name">{{ code.name }}</span>
        </div>
        <div class="code-desc">
          {{ code.description }}
          <n-text v-if="code.type !== 'HTTP'" depth="3">For {{ code.type }}.</n-text>
        </div>
      </n-card>
    </div>

    <div v-if="search && !categories.length" class="empty">
      <n-text depth="3">无匹配的状态码</n-text>
    </div>
  </div>
</template>

<style scoped>
.page {
  width: 100%;
}

.category {
  margin: 24px 0;
}

.category-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
  color: v-bind('themeVars.textColor1');
}

.code-card {
  margin-bottom: 8px;
}

.code-header {
  display: flex;
  align-items: baseline;
  gap: 10px;
  margin-bottom: 4px;
}

.code-number {
  font-size: 16px;
  font-weight: 700;
  color: v-bind('themeVars.primaryColor');
  font-family: 'Cascadia Code', monospace;
}

.code-name {
  font-size: 15px;
  font-weight: 600;
  color: v-bind('themeVars.textColor1');
}

.code-desc {
  opacity: 0.7;
  font-size: 13px;
  color: v-bind('themeVars.textColor2');
}

.empty {
  padding: 24px 0;
  text-align: center;
}
</style>