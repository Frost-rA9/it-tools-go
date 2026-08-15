<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, NInputNumber, NSelect, NRadioGroup, NRadioButton, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'
import ToolTextarea from '../../components/ToolTextarea.vue'

const message = useMessage()
const themeVars = useThemeVars()

const version = ref('v4')
const count = ref(1)
const namespace = ref('6ba7b811-9dad-11d1-80b4-00c04fd430c8')
const name = ref('')
const result = ref('')

const versionOptions = [
  { label: 'NIL', value: 'NIL' },
  { label: 'v1', value: 'v1' },
  { label: 'v3', value: 'v3' },
  { label: 'v4', value: 'v4' },
  { label: 'v5', value: 'v5' },
]

const namespaceOptions = [
  { label: 'DNS', value: '6ba7b810-9dad-11d1-80b4-00c04fd430c8' },
  { label: 'URL', value: '6ba7b811-9dad-11d1-80b4-00c04fd430c8' },
  { label: 'OID', value: '6ba7b812-9dad-11d1-80b4-00c04fd430c8' },
  { label: 'X500', value: '6ba7b814-9dad-11d1-80b4-00c04fd430c8' },
]

const showV35 = () => version.value === 'v3' || version.value === 'v5'

async function run() {
  try {
    const output = await RunTool(
      'uuid-generator',
      JSON.stringify({ version: version.value, count: count.value, namespace: namespace.value, name: name.value }),
    )
    result.value = JSON.parse(output).result
  } catch (e) {
    result.value = ''
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch([version, count, namespace, name], () => debouncedRun(), { immediate: true })

const { copy } = useClipboard({ source: result })

async function copyResult() {
  await copy()
  message.success('UUID 已复制到剪贴板')
}
</script>

<template>
  <div class="uuid-tool">
    <n-card title="UUID 生成器" class="card">
      <div class="field">
        <div class="field-label">UUID 版本</div>
        <n-radio-group v-model:value="version">
          <n-radio-button v-for="o in versionOptions" :key="o.value" :value="o.value">{{ o.label }}</n-radio-button>
        </n-radio-group>
      </div>

      <div class="field">
        <div class="field-label">生成数量</div>
        <n-input-number v-model:value="count" :min="1" :max="50" :step="1" style="width: 200px" />
      </div>

      <template v-if="showV35()">
        <div class="field">
          <div class="field-label">命名空间</div>
          <n-select v-model:value="namespace" :options="namespaceOptions" clearable placeholder="选择预设命名空间" />
        </div>

        <div class="field">
          <div class="field-label">自定义命名空间（可选）</div>
          <n-input v-model:value="namespace" placeholder="粘贴自定义命名空间 UUID…" />
        </div>

        <div class="field">
          <div class="field-label">名称</div>
          <n-input v-model:value="name" placeholder="输入名称…" />
        </div>
      </template>

      <ToolTextarea
        v-model:value="result"
        label="生成的 UUID"
        readonly
        :rows="4"
        placeholder="生成的 UUID 将显示在这里"
        monospace
        class="uuid-display"
      />

      <div class="btn-row">
        <n-button type="primary" @click="copyResult">复制</n-button>
        <n-button @click="run">重新生成</n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.uuid-tool {
  width: 100%;
}

.card {
  width: 100%;
}

.field {
  margin-bottom: 16px;
}

.field-label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.uuid-display :deep(textarea) {
  text-align: center;
}

.btn-row {
  display: flex;
  justify-content: center;
  gap: 12px;
}
</style>
