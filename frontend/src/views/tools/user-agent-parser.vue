<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { NCard, NGrid, NGi, NIcon, NTooltip, NTag, useMessage } from 'naive-ui'
import { Adjustments, Browser, Cpu, Devices, Engine } from '@vicons/tabler'
import { useDebounceFn } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

interface NameVersion {
  name: string
  version: string
}

interface UaResult {
  browser: NameVersion
  os: NameVersion
  device: { vendor: string; model: string; type: string }
}

interface Section {
  heading: string
  icon: object
  items: { label: string; value: string }[]
}

const message = useMessage()

const rawUa = ref(navigator.userAgent)
const result = ref<UaResult | null>(null)

async function runParse() {
  if (!rawUa.value.trim()) {
    result.value = null
    return
  }
  try {
    const output = await RunTool('user-agent-parser', JSON.stringify({ ua: rawUa.value }))
    result.value = JSON.parse(output)
  } catch (e) {
    result.value = null
    message.error(String(e))
  }
}

const debouncedParse = useDebounceFn(runParse, 150)
watch(rawUa, () => debouncedParse(), { immediate: true })

// ---- 前端 TS 补充（uap-go 无法提供的 engine / cpu / device.type）----

function detectEngine(ua: string): NameVersion {
  if (/Edg\//.test(ua)) return { name: 'Blink', version: '' }
  if (/Trident\//.test(ua)) return { name: 'Trident', version: (ua.match(/Trident\/([\d.]+)/) || [])[1] ?? '' }
  if (/Presto\//.test(ua)) return { name: 'Presto', version: (ua.match(/Presto\/([\d.]+)/) || [])[1] ?? '' }
  if (/OPR\//.test(ua)) return { name: 'Blink', version: '' }
  if (/Chrome\//.test(ua)) return { name: 'Blink', version: '' }
  if (/Gecko\/\d/.test(ua) && !/like Gecko/.test(ua)) return { name: 'Gecko', version: '' }
  if (/AppleWebKit\//.test(ua)) return { name: 'WebKit', version: (ua.match(/AppleWebKit\/([\d.]+)/) || [])[1] ?? '' }
  if (/Edge\//.test(ua)) return { name: 'EdgeHTML', version: (ua.match(/Edge\/([\d.]+)/) || [])[1] ?? '' }
  return { name: '', version: '' }
}

function detectCpu(ua: string): string {
  if (/aarch64|arm64/.test(ua)) return 'arm64'
  if (/x86_64|x64|amd64/.test(ua)) return 'amd64'
  if (/i386|i686|x86/.test(ua)) return 'x86'
  if (/armv7|arm\b/.test(ua)) return 'arm'
  return ''
}

function detectDeviceType(ua: string): string {
  if (/SmartTV|GoogleTV|Android TV/.test(ua)) return 'smarttv'
  if (/iPad|Tablet|Kindle|PlayBook/.test(ua)) return 'tablet'
  if (/iPhone|iPod|Android.+Mobile|Mobi/.test(ua)) return 'mobile'
  if (/PlayStation|Xbox/.test(ua)) return 'console'
  return ''
}

const sections = computed<Section[]>(() => {
  const ua = rawUa.value
  const engine = detectEngine(ua)
  const deviceType = detectDeviceType(ua)

  return [
    {
      heading: 'Browser',
      icon: Browser,
      items: [
        { label: 'Name', value: result.value?.browser.name ?? '' },
        { label: 'Version', value: result.value?.browser.version ?? '' },
      ],
    },
    {
      heading: 'Engine',
      icon: Engine,
      items: [
        { label: 'Name', value: engine.name },
        { label: 'Version', value: engine.version },
      ],
    },
    {
      heading: 'OS',
      icon: Adjustments,
      items: [
        { label: 'Name', value: result.value?.os.name ?? '' },
        { label: 'Version', value: result.value?.os.version ?? '' },
      ],
    },
    {
      heading: 'Device',
      icon: Devices,
      items: [
        { label: 'Model', value: result.value?.device.model ?? '' },
        { label: 'Type', value: result.value?.device.type || deviceType },
        { label: 'Vendor', value: result.value?.device.vendor ?? '' },
      ],
    },
    {
      heading: 'CPU',
      icon: Cpu,
      items: [{ label: 'Architecture', value: detectCpu(ua) }],
    },
  ]
})
</script>

<template>
  <n-card class="card">
    <ToolTextarea v-model:value="rawUa" label="User-Agent 字符串" :rows="2" placeholder="在此输入 User-Agent…" />

    <n-grid :x-gap="12" :y-gap="8" cols="1 s:2" responsive="screen">
      <n-gi v-for="section in sections" :key="section.heading">
        <n-card size="small" class="section-card">
          <div class="section-header">
            <n-icon size="28" :component="section.icon" :depth="3" />
            <span class="section-title">{{ section.heading }}</span>
          </div>

          <div class="section-body">
            <template v-for="item in section.items" :key="item.label">
              <n-tooltip v-if="item.value">
                <template #trigger>
                  <n-tag type="success" size="large" round :bordered="false">{{ item.value }}</n-tag>
                </template>
                {{ item.label }}
              </n-tooltip>
            </template>
          </div>

          <div class="section-fallback">
            <template v-for="item in section.items" :key="item.label">
              <span v-if="!item.value">{{ item.label }}：未知</span>
            </template>
          </div>
        </n-card>
      </n-gi>
    </n-grid>
  </n-card>
</template>

<style scoped>
.card {
  min-width: 480px;
  width: 100%;
}

.section-card {
  height: 100%;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}

.section-title {
  font-size: 17px;
  font-weight: 600;
}

.section-body {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.section-fallback {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-top: 6px;
  font-size: 13px;
  opacity: 0.7;
}
</style>