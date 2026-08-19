<script setup lang="ts">
import { ref } from 'vue'
import { NCard, NInputNumber, NInputGroup, NInputGroupLabel, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

interface Scale {
  scale: string
  title: string
  unit: string
}

const scales: Scale[] = [
  { scale: 'kelvin', title: '开尔文', unit: 'K' },
  { scale: 'celsius', title: '摄氏', unit: '°C' },
  { scale: 'fahrenheit', title: '华氏', unit: '°F' },
  { scale: 'rankine', title: '兰金', unit: '°R' },
  { scale: 'delisle', title: '德利尔', unit: '°De' },
  { scale: 'newton', title: '牛顿', unit: '°N' },
  { scale: 'reaumur', title: '列氏', unit: '°Ré' },
  { scale: 'romer', title: '罗氏', unit: '°Rø' },
]

const values = ref<Record<string, number | null>>({
  kelvin: 273.15,
  celsius: 0,
  fahrenheit: 32,
  rankine: 491.67,
  delisle: 559.73,
  newton: 0,
  reaumur: 0,
  romer: 7.5,
})

async function run(from: string) {
  const value = values.value[from]
  if (value == null) return
  try {
    const output = await RunTool(
      'temperature-converter',
      JSON.stringify({ value, from }),
    )
    const parsed = JSON.parse(output)
    for (const r of parsed.results) {
      values.value[r.scale] = r.value
    }
  } catch (e) {
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)

// 编辑任一温标输入即触发布局换算（含被编辑行回填两位小数）。
function onScaleUpdate(scale: string) {
  debouncedRun(scale)
}
</script>

<template>
  <div class="temperature-tool">
    <n-card title="温度转换器" class="tool-card">
      <div class="unit-list">
        <n-input-group v-for="s in scales" :key="s.scale" class="unit-row">
          <n-input-group-label class="unit-label">{{ s.title }}</n-input-group-label>
          <n-input-number
            v-model:value="values[s.scale]"
            class="unit-input"
            :show-button="false"
            :placeholder="s.title"
            @update:value="() => onScaleUpdate(s.scale)"
          />
          <n-input-group-label class="unit-suffix">{{ s.unit }}</n-input-group-label>
        </n-input-group>
      </div>
      <div class="hint">编辑任一温标，其余温标将同步换算（保留两位小数）。</div>
    </n-card>
  </div>
</template>

<style scoped>
.tool-card {
  max-width: 640px;
}

.unit-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.unit-label {
  width: 100px;
  justify-content: flex-start;
}

.unit-suffix {
  width: 60px;
}

.unit-input {
  flex: 1;
}

.hint {
  margin-top: 16px;
  font-size: 13px;
  color: v-bind('themeVars.textColor3');
}
</style>