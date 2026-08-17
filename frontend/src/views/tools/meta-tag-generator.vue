<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { NCard, NInput, NInputGroup, NInputGroupLabel, NSelect, NDynamicInput, NText, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

interface Option {
  label: string
  value?: string
  type?: string
  key?: string
  children?: Option[]
}

interface Element {
  key: string
  label: string
  placeholder: string
  type: 'input' | 'select' | 'input-multiple'
  options?: Option[]
}

interface Schema {
  name: string
  elements: Element[]
}

interface SchemasResult {
  base: Schema[]
  types: Record<string, Schema>
}

const message = useMessage()
const themeVars = useThemeVars()

const schemas = ref<SchemasResult | null>(null)
const metadata = reactive<Record<string, any>>({
  type: 'website',
  'twitter:card': 'summary_large_image',
})

// 左列的三个基础区块（General/Image/Twitter）。
const baseSections = computed<Schema[]>(() => schemas.value?.base ?? [])

// 按页面类型追加的区块（选非 website 类型时出现）。
const additionalSection = computed<Schema | null>(() => schemas.value?.types[metadata.type] ?? null)

const metaOutput = ref('')

async function fetchSchemas() {
  try {
    const output = await RunTool('meta-tag-generator', JSON.stringify({ action: 'schemas' }))
    schemas.value = JSON.parse(output).schemas
  } catch (e) {
    message.error(String(e))
  }
}

async function runGenerate() {
  try {
    const output = await RunTool('meta-tag-generator', JSON.stringify({ action: 'meta', metadata }))
    metaOutput.value = JSON.parse(output).meta
  } catch (e) {
    metaOutput.value = ''
    message.error(String(e))
  }
}

const debouncedGenerate = useDebounceFn(runGenerate, 200)

// type 变更时清空上一 schema 的键（对齐 it-tools）。
watch(
  () => metadata.type,
  (newType, oldType) => {
    const prev = schemas.value?.types[oldType]
    if (prev) {
      for (const el of prev.elements) {
        delete metadata[el.key]
      }
    }
    // 确保新类型的键存在并默认空。
    const next = schemas.value?.types[newType]
    if (next) {
      for (const el of next.elements) {
        if (!(el.key in metadata)) metadata[el.key] = ''
      }
    }
    debouncedGenerate()
  },
)

// 监听所有字段变化 → 防抖生成。
watch(metadata, () => debouncedGenerate(), { deep: true })

onMounted(async () => {
  await fetchSchemas()
  runGenerate()
})

const { copy } = useClipboard()

async function copyMeta() {
  await copy(metaOutput.value)
  message.success('meta 标签已复制到剪贴板')
}
</script>

<template>
  <div class="layout">
    <div class="left-col">
      <n-card v-for="section in baseSections" :key="section.name" :title="section.name" class="card">
        <n-input-group v-for="el in section.elements" :key="el.key" class="field">
          <n-input-group-label class="field-label">{{ el.label }}</n-input-group-label>

          <n-input v-if="el.type === 'input'" v-model:value="metadata[el.key]" :placeholder="el.placeholder" clearable class="field-input" />

          <n-dynamic-input
            v-else-if="el.type === 'input-multiple'"
            v-model:value="metadata[el.key]"
            :min="1"
            :placeholder="el.placeholder"
            :default-value="['']"
            :show-sort-button="true"
            class="field-input"
          />

          <n-select
            v-else-if="el.type === 'select'"
            v-model:value="metadata[el.key]"
            :placeholder="el.placeholder"
            :options="(el.options as any)"
            class="field-input"
          />
        </n-input-group>
      </n-card>

      <n-card v-if="additionalSection" :title="additionalSection.name" class="card">
        <n-input-group v-for="el in additionalSection.elements" :key="el.key" class="field">
          <n-input-group-label class="field-label">{{ el.label }}</n-input-group-label>

          <n-input v-if="el.type === 'input'" v-model:value="metadata[el.key]" :placeholder="el.placeholder" clearable class="field-input" />

          <n-dynamic-input
            v-else-if="el.type === 'input-multiple'"
            v-model:value="metadata[el.key]"
            :min="1"
            :placeholder="el.placeholder"
            :default-value="['']"
            :show-sort-button="true"
            class="field-input"
          />

          <n-select
            v-else-if="el.type === 'select'"
            v-model:value="metadata[el.key]"
            :placeholder="el.placeholder"
            :options="(el.options as any)"
            class="field-input"
          />
        </n-input-group>
      </n-card>
    </div>

    <n-card title="生成的 meta 标签" class="card meta-card">
      <ToolTextarea v-model:value="metaOutput" :rows="18" readonly monospace placeholder="生成的 meta 标签将显示在这里" />
      <div class="copy-row">
        <n-button type="primary" @click="copyMeta">复制 meta 标签</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.layout {
  flex: 1 1 100% !important;
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  align-items: flex-start;
  width: 100%;
}

.left-col {
  flex: 1 1 360px;
  min-width: 300px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.card {
  min-width: 0;
}

.meta-card {
  flex: 1 1 320px;
  min-width: 280px;
}

.field {
  margin-bottom: 5px;
  width: 100%;
}

.field-label {
  flex: 0 0 130px;
  text-align: right;
}

.field-input {
  flex: 1 1 auto;
  min-width: 0;
}

.copy-row {
  display: flex;
  justify-content: center;
  margin-top: 12px;
}
</style>