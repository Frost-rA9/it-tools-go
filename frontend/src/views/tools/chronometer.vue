<script setup lang="ts">
import { onUnmounted, ref } from 'vue'
import { NCard, NButton } from 'naive-ui'
import { useThemeVars } from 'naive-ui'

const themeVars = useThemeVars()

const isRunning = ref(false)
const accumulated = ref(0)
const startedAt = ref(0)
const display = ref(0)

let rafId: number | null = null

function tick() {
  display.value = accumulated.value + (isRunning.value ? Date.now() - startedAt.value : 0)
  rafId = requestAnimationFrame(tick)
}

function start() {
  startedAt.value = Date.now()
  isRunning.value = true
  rafId = requestAnimationFrame(tick)
}

function stop() {
  accumulated.value += Date.now() - startedAt.value
  isRunning.value = false
  if (rafId !== null) cancelAnimationFrame(rafId)
  rafId = null
}

function reset() {
  if (isRunning.value) {
    // 运行中重置：保持计时继续，仅清零累计
    accumulated.value = 0
    startedAt.value = Date.now()
  } else {
    accumulated.value = 0
    display.value = 0
  }
}

onUnmounted(() => {
  if (rafId !== null) cancelAnimationFrame(rafId)
})

// formatMs 与参考项目 chronometer.service.ts 一致：[H:]MM:SS.mmm。
function formatMs(msTotal: number) {
  const ms = msTotal % 1000
  const secs = ((msTotal - ms) / 1000) % 60
  const mins = (((msTotal - ms) / 1000 - secs) / 60) % 60
  const hrs = (((msTotal - ms) / 1000 - secs) / 60 - mins) / 60
  const hrsString = hrs > 0 ? `${hrs.toString().padStart(2, '0')}:` : ''
  return `${hrsString}${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}.${ms
    .toString()
    .padStart(3, '0')}`
}
</script>

<template>
  <div class="chronometer-tool">
    <n-card title="秒表" class="tool-card">
      <div class="duration">{{ formatMs(display) }}</div>
      <div class="controls">
        <n-button v-if="!isRunning" type="primary" size="large" @click="start">开始</n-button>
        <n-button v-else type="warning" size="large" @click="stop">停止</n-button>
        <n-button size="large" @click="reset">重置</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.tool-card {
  max-width: 640px;
}

.duration {
  text-align: center;
  font-size: 40px;
  font-family: 'Cascadia Code', ui-monospace, monospace;
  color: v-bind('themeVars.textColor1');
  margin: 20px 0;
}

.controls {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-bottom: 8px;
}
</style>