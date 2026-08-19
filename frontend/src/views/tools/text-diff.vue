<script setup lang="ts">
import { ref, watch } from 'vue'
import { NAlert, NButton, NCard, useMessage, useThemeVars } from 'naive-ui'
import { useClipboard, useDebounceFn } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

interface Seg {
  t: string
  c: boolean
}

interface DiffRow {
  type: 'equal' | 'delete' | 'insert'
  old_no?: number
  new_no?: number
  old?: string
  new?: string
  old_segments?: Seg[]
  new_segments?: Seg[]
}

interface DiffResult {
  equal: boolean
  stats: { old_lines: number; new_lines: number; removed: number; added: number; changed: number }
  rows: DiffRow[]
}

const oldText = ref('可对比的两段文本\n第二行：保持不变\n第三行这里将被修改')
const newText = ref('可对比的两段文本\n第二行：保持不变\n第三行这里已发生变化\n末尾新增一行')
const result = ref<DiffResult | null>(null)

async function run() {
  try {
    const output = await RunTool('text-diff', JSON.stringify({ old_text: oldText.value, new_text: newText.value }))
    result.value = JSON.parse(output)
  } catch (e) {
    result.value = null
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch([oldText, newText], () => debouncedRun(), { immediate: true })

const { copy } = useClipboard()

async function copyNew() {
  await copy(newText.value)
  message.success('新文本已复制到剪贴板')
}
</script>

<template>
  <n-card title="文本比较" class="tool-card">
    <div class="inputs">
      <ToolTextarea v-model:value="oldText" label="原文本" :rows="4" placeholder="在此输入原文本…" />
      <ToolTextarea v-model:value="newText" label="新文本" :rows="4" placeholder="在此输入新文本…" />
    </div>

    <div v-if="result" class="summary">
      <n-alert v-if="result.equal" type="success" :bordered="false" class="identical">两文本相同</n-alert>
      <div v-else class="stats-row">
        <span class="stats-text">
          新增 <b class="add">{{ result.stats.added }}</b> 行 ·
          删除 <b class="del">{{ result.stats.removed }}</b> 行 ·
          修改 <b class="mod">{{ result.stats.changed }}</b> 行
        </span>
        <n-button size="small" @click="copyNew" :disabled="!newText">复制新文本</n-button>
      </div>
    </div>

    <div v-if="result && result.rows.length" class="diff-wrap">
      <div class="diff-grid">
        <div class="cell cell-head no-col">旧行号</div>
        <div class="cell cell-head">旧内容</div>
        <div class="cell cell-head no-col">新行号</div>
        <div class="cell cell-head">新内容</div>

        <template v-for="r in result.rows" :key="`${r.type}-${r.old_no ?? ''}-${r.new_no ?? ''}-${r.old_no}-${r.new_no}`">
          <div class="cell no-col">{{ r.old_no ?? '' }}</div>
          <div class="cell content" :class="r.type">
            <template v-if="r.type === 'equal'">{{ r.old }}</template>
            <template v-else-if="r.type === 'delete'">
              <span v-for="(s, i) in r.old_segments" :key="i" class="sg" :class="{ 'sg-chg': s.c }">{{ s.t }}</span>
            </template>
          </div>
          <div class="cell no-col">{{ r.new_no ?? '' }}</div>
          <div class="cell content" :class="r.type">
            <template v-if="r.type === 'equal'">{{ r.new }}</template>
            <template v-else-if="r.type === 'insert'">
              <span v-for="(s, i) in r.new_segments" :key="i" class="sg" :class="{ 'sg-chg': s.c }">{{ s.t }}</span>
            </template>
          </div>
        </template>
      </div>
    </div>
  </n-card>
</template>

<style scoped>
.inputs {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.inputs :deep(.tool-textarea) {
  flex: 1;
  min-width: 280px;
}

.summary {
  margin: 2px 0 12px;
}

.identical {
  margin-bottom: 4px;
}

.stats-row {
  display: flex;
  align-items: center;
  gap: 14px;
  flex-wrap: wrap;
}

.stats-text {
  font-size: 14px;
}

.stats-text .add {
  color: v-bind('themeVars.successColor');
}

.stats-text .del {
  color: v-bind('themeVars.errorColor');
}

.stats-text .mod {
  color: v-bind('themeVars.warningColor');
}

.diff-wrap {
  border: 1px solid v-bind('themeVars.borderColor');
  border-radius: 4px;
  overflow: hidden;
}

.diff-grid {
  display: grid;
  grid-template-columns: 48px minmax(0, 1fr) 48px minmax(0, 1fr);
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
  font-variant-ligatures: none;
  font-feature-settings: 'liga' 0, 'calt' 0;
  font-size: 13px;
}

.cell {
  padding: 2px 8px;
  line-height: 1.55;
  min-height: 24px;
  box-sizing: border-box;
  border-bottom: 1px solid v-bind('themeVars.borderColor');
  border-right: 1px solid v-bind('themeVars.borderColor');
  white-space: pre-wrap;
  word-break: break-all;
  color: v-bind('themeVars.textColor1');
}

.cell-head {
  font-weight: 700;
  font-size: 12px;
  color: v-bind('themeVars.textColor3');
  background: v-bind('themeVars.tableHeaderColor');
}

.no-col {
  color: v-bind('themeVars.textColor3');
  text-align: right;
  user-select: none;
  background: v-bind('themeVars.tableHeaderColor');
}

.cell.equal {
  background: transparent;
}

.cell.delete {
  background: color-mix(in srgb, v-bind('themeVars.errorColor') 12%, transparent);
}

.cell.insert {
  background: color-mix(in srgb, v-bind('themeVars.successColor') 12%, transparent);
}

.sg-chg {
  border-radius: 3px;
  padding: 0 2px;
}

.cell.delete .sg-chg {
  background: color-mix(in srgb, v-bind('themeVars.errorColor') 35%, transparent);
}

.cell.insert .sg-chg {
  background: color-mix(in srgb, v-bind('themeVars.successColor') 35%, transparent);
}
</style>