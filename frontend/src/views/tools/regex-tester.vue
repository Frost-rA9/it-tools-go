<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NCheckbox, NTable, NAlert, NEmpty, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

interface GroupCapture {
  name: string
  value: string
  start: number
  end: number
}

interface MatchResult {
  index: number
  value: string
  captures: GroupCapture[]
  groups: GroupCapture[]
}

const themeVars = useThemeVars()

const regex = ref('\\b\\w+@\\w+\\.\\w+\\b')
const text = ref('Contact: alice@example.com or bob@test.org for help.')
const flagG = ref(true)
const flagI = ref(false)
const flagM = ref(false)
const flagS = ref(false)
const matches = ref<MatchResult[]>([])
const errorMessage = ref('')

async function run() {
  try {
    const output = await RunTool(
      'regex-tester',
      JSON.stringify({
        regex: regex.value,
        text: text.value,
        flags: { g: flagG.value, i: flagI.value, m: flagM.value, s: flagS.value },
      }),
    )
    matches.value = JSON.parse(output).matches ?? []
    errorMessage.value = ''
  } catch (e) {
    matches.value = []
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 250)
watch([regex, text, flagG, flagI, flagM, flagS], () => debouncedRun(), { immediate: true })
</script>

<template>
  <n-card title="Regex 测试器 — 输入" class="tool-card">
    <div class="field">
      <div class="field-label">正则表达式</div>
      <n-input v-model:value="regex" placeholder="输入要测试的正则表达式…" monospace />
    </div>

    <div class="flags-row">
      <n-checkbox v-model:checked="flagG">全局 (g)</n-checkbox>
      <n-checkbox v-model:checked="flagI">忽略大小写 (i)</n-checkbox>
      <n-checkbox v-model:checked="flagM">多行 (m)</n-checkbox>
      <n-checkbox v-model:checked="flagS">单行 (s)</n-checkbox>
    </div>

    <ToolTextarea v-model:value="text" label="测试文本" :rows="12" placeholder="输入要匹配的文本…" />

    <n-alert v-if="errorMessage" type="error" class="error-alert">
      {{ errorMessage }}
    </n-alert>
  </n-card>

  <n-card title="Regex 测试器 — 匹配结果" class="tool-card">
    <n-table v-if="matches.length > 0" :bordered="false" size="small" class="result-table">
      <thead>
        <tr>
          <th style="width: 80px">索引</th>
          <th>匹配值</th>
          <th>编号捕获组</th>
          <th>命名捕获组</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(m, i) in matches" :key="i">
          <td class="mono">{{ m.index }}</td>
          <td class="mono">{{ m.value }}</td>
          <td>
            <div v-for="(c, j) in m.captures" :key="j" class="group-item">
              <span class="group-name">{{ c.name }}</span> = {{ c.value }}
              <span class="group-range">[{{ c.start }} - {{ c.end }}]</span>
            </div>
          </td>
          <td>
            <div v-for="(g, k) in m.groups" :key="k" class="group-item">
              <span class="group-name">{{ g.name }}</span> = {{ g.value }}
              <span class="group-range">[{{ g.start }} - {{ g.end }}]</span>
            </div>
          </td>
        </tr>
      </tbody>
    </n-table>

    <n-empty v-else-if="!errorMessage" description="无匹配结果" class="empty" />
  </n-card>
</template>

<style scoped>
.field {
  margin-bottom: 16px;
}

.field-label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.field :deep(input) {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}

.flags-row {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 16px;
}

.error-alert {
  margin-bottom: 16px;
}

.result-table {
  margin-bottom: 16px;
}

.mono {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
  font-size: 13px;
}

.group-item {
  font-size: 13px;
  line-height: 1.7;
  word-break: break-all;
}

.group-name {
  color: v-bind('themeVars.primaryColor');
  font-weight: 600;
}

.group-range {
  opacity: 0.6;
}

.empty {
  padding: 24px 0;
}
</style>
