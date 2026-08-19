<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  NCard,
  NForm,
  NFormItem,
  NInputNumber,
  NDatePicker,
  NSelect,
  NStatistic,
  useMessage,
  useThemeVars,
} from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

// 参考示例：3 分钟洗 5 个盘子，共 500 个盘子。
const unitCount = ref(500)
const unitPerTimeSpan = ref(5)
const timeSpan = ref(3)
const timeSpanUnit = ref(60_000) // 分钟
const startedAt = ref(Date.now())

const timeUnitOptions = [
  { label: '毫秒', value: 1 },
  { label: '秒', value: 1000 },
  { label: '分钟', value: 60_000 },
  { label: '小时', value: 3_600_000 },
  { label: '天', value: 86_400_000 },
]

const durationText = ref('')
const endAtText = ref('')

async function run() {
  if (unitCount.value <= 0 || unitPerTimeSpan.value <= 0 || timeSpan.value <= 0) {
    durationText.value = ''
    endAtText.value = ''
    return
  }
  try {
    const output = await RunTool(
      'eta-calculator',
      JSON.stringify({
        unitCount: unitCount.value,
        unitPerTimeSpan: unitPerTimeSpan.value,
        timeSpan: timeSpan.value,
        timeSpanUnitMultiplier: timeSpanUnit.value,
        startedAtMs: startedAt.value,
      }),
    )
    const parsed = JSON.parse(output)
    durationText.value = parsed.durationText
    endAtText.value = parsed.endAtText
  } catch (e) {
    durationText.value = ''
    endAtText.value = ''
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch([unitCount, unitPerTimeSpan, timeSpan, timeSpanUnit, startedAt], () => debouncedRun(), {
  immediate: true,
})
</script>

<template>
  <div class="eta-tool">
    <n-card title="ETA 计算器" class="tool-card">
      <n-form label-placement="top" :show-feedback="false">
        <div class="grid-2">
          <n-form-item label="待消耗的单位总数">
            <n-input-number v-model:value="unitCount" :min="1" class="full" />
          </n-form-item>
          <n-form-item label="消耗开始时间">
            <n-date-picker v-model:value="startedAt" type="datetime" class="full" />
          </n-form-item>
        </div>

        <n-form-item label="每个时间段消耗的单位数" class="consume-form-item">
          <div class="consume-row">
            <n-input-number v-model:value="unitPerTimeSpan" :min="1" class="consume-count" />
            <div class="consume-meta">
              <span class="hint">单位 /</span>
              <n-input-number v-model:value="timeSpan" :min="1" class="consume-span" />
              <n-select v-model:value="timeSpanUnit" :options="timeUnitOptions" class="consume-unit" />
            </div>
          </div>
        </n-form-item>
      </n-form>

      <div class="result-grid">
        <n-card title="总时长" size="small" class="stat-card">
          <n-statistic :value="durationText" />
        </n-card>
        <n-card title="预计结束时间" size="small" class="stat-card">
          <n-statistic :value="endAtText" />
        </n-card>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.full {
  width: 100%;
}

.grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.consume-form-item {
  margin-top: 12px;
}

.consume-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  width: 100%;
}

.consume-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.consume-count {
  width: 100%;
}

.consume-span {
  flex: 1;
  min-width: 0;
}

.consume-unit {
  width: 110px;
  flex-shrink: 0;
}

.hint {
  color: v-bind('themeVars.textColor3');
  white-space: nowrap;
  flex-shrink: 0;
}

.result-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 8px;
}
</style>