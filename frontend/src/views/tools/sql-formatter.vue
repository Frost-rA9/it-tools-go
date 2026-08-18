<script setup lang="ts">
import { ref, watch } from 'vue'
import { NSwitch, NAlert, NCard, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import CodeOutput from '../../components/CodeOutput.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const themeVars = useThemeVars()

const sqlInput = ref("select u.id,o.amount from users u join orders o on u.id=o.user_id where u.active=true and u.age>18;")
const upperKeywords = ref(true)
const formatted = ref('')
const lineCount = ref(0)
const errorMessage = ref('')

async function run() {
  try {
    const output = await RunTool(
      'sql-formatter',
      JSON.stringify({ sql: sqlInput.value, upper_keywords: upperKeywords.value }),
    )
    const result = JSON.parse(output)
    formatted.value = result.formatted
    lineCount.value = result.line_count
    errorMessage.value = ''
  } catch (e) {
    formatted.value = ''
    lineCount.value = 0
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 200)
watch([sqlInput, upperKeywords], () => debouncedRun(), { immediate: true })
</script>

<template>
  <n-card title="SQL 格式化 — 输入" class="tool-card">
    <ToolTextarea v-model:value="sqlInput" label="输入 SQL" :rows="20" placeholder="在此输入 SQL…" monospace />

    <div class="options-row">
      <span class="option-label">关键字大写</span>
      <n-switch v-model:value="upperKeywords" />
    </div>

    <n-alert v-if="errorMessage" type="error" class="error-alert">
      {{ errorMessage }}
    </n-alert>
  </n-card>

  <n-card title="SQL 格式化 — 输出" class="tool-card">
    <CodeOutput label="格式化结果" :value="formatted" language="sql" :rows="20" />

    <div class="stats-row">
      <span class="stat">{{ lineCount }} 行</span>
    </div>

  </n-card>
</template>

<style scoped>
.options-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.option-label {
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
}

.error-alert {
  margin-bottom: 16px;
}

.stats-row {
  display: flex;
  gap: 16px;
  justify-content: center;
  margin-bottom: 16px;
}

.stat {
  font-size: 13px;
  color: v-bind('themeVars.textColor3');
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}
</style>
