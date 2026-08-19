<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  NCard,
  NInput,
  NInputNumber,
  NButton,
  NTable,
  NDivider,
  NIcon,
  useMessage,
  useThemeVars,
} from 'naive-ui'
import { Plus, Trash } from '@vicons/tabler'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

interface Suite {
  title: string
  data: (number | null)[]
}

interface BenchResult {
  position: number
  title: string
  size: number
  mean: string
  variance: string
}

const suites = ref<Suite[]>([
  { title: 'Suite 1', data: [5, 10] },
  { title: 'Suite 2', data: [8, 12] },
])
const unit = ref('ms')
const results = ref<BenchResult[]>([])
const markdown = ref('')
const bulletList = ref('')

async function run() {
  try {
    const payload = {
      unit: unit.value,
      suites: suites.value.map((s) => ({
        title: s.title,
        data: s.data.filter((v): v is number => v != null),
      })),
    }
    const output = await RunTool('benchmark-builder', JSON.stringify(payload))
    const parsed = JSON.parse(output)
    results.value = parsed.results
    markdown.value = parsed.markdown
    bulletList.value = parsed.bulletList
  } catch (e) {
    results.value = []
    markdown.value = ''
    bulletList.value = ''
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch(
  [
    () => suites.value.map((s) => s.title).join('\u0000'),
    () => suites.value.map((s) => s.data.join(',')).join('\u0000'),
    unit,
  ],
  () => debouncedRun(),
  { immediate: true },
)

function addValue(index: number) {
  suites.value[index].data.push(null)
}

function removeValue(suiteIndex: number, valueIndex: number) {
  suites.value[suiteIndex].data.splice(valueIndex, 1)
}

function addSuite() {
  suites.value.push({ title: `Suite ${suites.value.length + 1}`, data: [0] })
}

function removeSuite(index: number) {
  suites.value.splice(index, 1)
}

function resetSuites() {
  suites.value = [
    { title: 'Suite 1', data: [] },
    { title: 'Suite 2', data: [] },
  ]
}

const { copy } = useClipboard()

async function copyMarkdown() {
  if (!markdown.value) return
  await copy(markdown.value)
  message.success('Markdown 表格已复制到剪贴板')
}

async function copyBulletList() {
  if (!bulletList.value) return
  await copy(bulletList.value)
  message.success('列表已复制到剪贴板')
}
</script>

<template>
  <div class="benchmark-tool">
    <div class="suites-row">
      <n-card v-for="(suite, index) in suites" :key="index" class="suite-card" :title="`套件 ${index + 1}`">
        <div class="field">
          <div class="field-label">套件名称</div>
          <n-input v-model:value="suite.title" placeholder="套件名称…" />
        </div>
        <n-divider />
        <div class="field">
          <div class="field-label">测量值</div>
          <div v-for="(v, vi) in suite.data" :key="vi" class="value-row">
            <n-input-number v-model:value="suite.data[vi]" class="value-input" :show-button="false" placeholder="测量值…" />
            <n-button size="small" quaternary circle @click="removeValue(index, vi)">
              <n-icon :component="Trash" size="18" />
            </n-button>
          </div>
          <n-button size="small" dashed block class="add-value" @click="addValue(index)">
            <n-icon :component="Plus" size="16" />
            添加测量值
          </n-button>
        </div>
        <div class="card-actions">
          <n-button v-if="suites.length > 1" size="tiny" quaternary type="error" @click="removeSuite(index)">
            删除套件
          </n-button>
          <n-button size="tiny" quaternary @click="addSuite">添加套件</n-button>
        </div>
      </n-card>
    </div>

    <n-card title="结果" class="tool-card">
      <div class="unit-row">
        <div class="unit-label">单位</div>
        <n-input v-model:value="unit" placeholder="如 ms…" class="unit-input" />
        <n-button size="small" @click="resetSuites">重置套件</n-button>
      </div>

      <n-table v-if="results.length" :bordered="true" size="small" class="result-table">
        <thead>
          <tr>
            <th>排名</th>
            <th>套件</th>
            <th>样本数</th>
            <th>均值</th>
            <th>方差</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in results" :key="r.position">
            <td>{{ r.position }}</td>
            <td>{{ r.title }}</td>
            <td>{{ r.size }}</td>
            <td>{{ r.mean }}</td>
            <td>{{ r.variance }}</td>
          </tr>
        </tbody>
      </n-table>

      <div v-else class="empty-hint">暂无数据，请添加测量值。</div>

      <div class="copy-row">
        <n-button size="small" :disabled="!markdown" @click="copyMarkdown">复制为 Markdown 表格</n-button>
        <n-button size="small" :disabled="!bulletList" @click="copyBulletList">复制为列表</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.suites-row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 16px;
}

.suite-card {
  width: 294px;
  flex-shrink: 0;
}

.field {
  margin-bottom: 12px;
}

.field-label {
  font-size: 13px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.value-row {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 6px;
}

.value-input {
  flex: 1;
}

.add-value {
  margin-top: 4px;
}

.card-actions {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-top: 8px;
}

.tool-card {
  width: 100%;
}

.unit-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.unit-label {
  color: v-bind('themeVars.textColor2');
  white-space: nowrap;
}

.unit-input {
  width: 140px;
}

.result-table {
  margin-bottom: 16px;
}

.empty-hint {
  text-align: center;
  color: v-bind('themeVars.textColor3');
  padding: 12px 0;
  font-size: 13px;
}

.copy-row {
  display: flex;
  justify-content: center;
  gap: 12px;
}
</style>