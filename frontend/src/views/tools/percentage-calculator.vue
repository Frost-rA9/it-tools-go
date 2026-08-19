<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInputNumber, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

// 第一组：X% of Y
const percentX = ref<number | null>(50)
const percentY = ref<number | null>(200)
const percentResult = ref('')

// 第二组：X 是 Y 的百分之几
const ratioX = ref<number | null>(25)
const ratioY = ref<number | null>(200)
const ratioResult = ref('')

// 第三组：从 X 到 Y 的增减百分比
const changeFrom = ref<number | null>(50)
const changeTo = ref<number | null>(75)
const changeResult = ref('')

async function runPercentOf() {
  if (percentX.value == null || percentY.value == null) return
  try {
    const output = await RunTool(
      'percentage-calculator',
      JSON.stringify({ mode: 'percent_of', x: percentX.value, y: percentY.value }),
    )
    percentResult.value = JSON.parse(output).result
  } catch (e) {
    percentResult.value = ''
    message.error(String(e))
  }
}

async function runWhatPercent() {
  if (ratioX.value == null || ratioY.value == null) return
  try {
    const output = await RunTool(
      'percentage-calculator',
      JSON.stringify({ mode: 'what_percent', x: ratioX.value, y: ratioY.value }),
    )
    ratioResult.value = JSON.parse(output).result
  } catch (e) {
    ratioResult.value = ''
  }
}

async function runChange() {
  if (changeFrom.value == null || changeTo.value == null) return
  try {
    const output = await RunTool(
      'percentage-calculator',
      JSON.stringify({ mode: 'change', x: changeFrom.value, y: changeTo.value }),
    )
    changeResult.value = JSON.parse(output).result
  } catch (e) {
    changeResult.value = ''
  }
}

const debouncedPercent = useDebounceFn(runPercentOf, 150)
const debouncedRatio = useDebounceFn(runWhatPercent, 150)
const debouncedChange = useDebounceFn(runChange, 150)

watch([percentX, percentY], () => debouncedPercent(), { immediate: true })
watch([ratioX, ratioY], () => debouncedRatio(), { immediate: true })
watch([changeFrom, changeTo], () => debouncedChange(), { immediate: true })

const { copy } = useClipboard()

async function copyResult(value: string) {
  if (!value) return
  await copy(value)
  message.success('已复制到剪贴板')
}
</script>

<template>
  <div class="percentage-tool">
    <n-card title="X% 的 Y 是多少" class="tool-card">
      <div class="calc-row">
        <n-input-number v-model:value="percentX" :min="0" class="num-input" placeholder="X" />
        <span class="op-label">% 的</span>
        <n-input-number v-model:value="percentY" :min="0" class="num-input" placeholder="Y" />
        <span class="op-eq">=</span>
        <n-input-number :value="Number(percentResult) || null" readonly class="num-input result" placeholder="结果" />
        <n-button size="small" :disabled="!percentResult" @click="copyResult(percentResult)">复制</n-button>
      </div>
    </n-card>

    <n-card title="X 是 Y 的百分之几" class="tool-card">
      <div class="calc-row">
        <n-input-number v-model:value="ratioX" :min="0" class="num-input" placeholder="X" />
        <span class="op-label">是</span>
        <n-input-number v-model:value="ratioY" :min="0" class="num-input" placeholder="Y" />
        <span class="op-label">的百分之</span>
        <n-input-number :value="Number(ratioResult) || null" readonly class="num-input result" placeholder="结果" />
        <n-button size="small" :disabled="!ratioResult" @click="copyResult(ratioResult)">复制</n-button>
      </div>
    </n-card>

    <n-card title="从 X 到 Y 的增减百分比" class="tool-card">
      <div class="calc-row">
        <n-input-number v-model:value="changeFrom" :min="0" class="num-input" placeholder="From" />
        <span class="op-label">到</span>
        <n-input-number v-model:value="changeTo" :min="0" class="num-input" placeholder="To" />
        <span class="op-eq">=</span>
        <n-input-number :value="Number(changeResult) || null" readonly class="num-input result" placeholder="结果" />
        <n-button size="small" :disabled="!changeResult" @click="copyResult(changeResult)">复制</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.tool-card {
  margin-bottom: 16px;
}

.calc-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.num-input {
  width: 130px;
}

.num-input.result {
  font-weight: 600;
}

.op-label,
.op-eq {
  color: v-bind('themeVars.textColor2');
  white-space: nowrap;
}

.op-eq {
  color: v-bind('themeVars.textColor1');
}
</style>