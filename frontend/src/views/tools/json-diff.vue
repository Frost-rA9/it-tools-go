<script setup lang="ts">
import { computed, h, ref, watch, type VNode } from 'vue'
import { NCard, NFormItem, NSwitch, NText, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

interface DiffNode {
  key: string | number
  type: 'object' | 'array' | 'value'
  status: 'added' | 'removed' | 'updated' | 'unchanged' | 'children-updated'
  oldValue?: unknown
  value?: unknown
  children?: DiffNode[]
}

interface DiffResult {
  same: boolean
  root: DiffNode | null
}

const message = useMessage()
const themeVars = useThemeVars()

const rawLeftJson = ref('')
const rawRightJson = ref('')
const onlyShowDifferences = ref(false)
const result = ref<DiffResult | null>(null)
const errorText = ref('')

async function runDiff() {
  errorText.value = ''
  if (!rawLeftJson.value.trim() || !rawRightJson.value.trim()) {
    result.value = null
    return
  }
  try {
    const output = await RunTool(
      'json-diff',
      JSON.stringify({
        left: rawLeftJson.value,
        right: rawRightJson.value,
        only_show_differences: onlyShowDifferences.value,
      }),
    )
    result.value = JSON.parse(output)
  } catch (e) {
    result.value = null
    errorText.value = String(e).replace(/^Error:\s*/, '')
  }
}

const debouncedDiff = useDebounceFn(runDiff, 150)
watch([rawLeftJson, rawRightJson, onlyShowDifferences], () => debouncedDiff(), { immediate: true })

function formatValue(v: unknown) {
  if (v === undefined) return ''
  return JSON.stringify(v)
}

// 递归渲染差异树（对齐 it-tools 的 diff-viewer）。
function renderDiff(node: DiffNode, showKeys: boolean): VNode {
  const parts: VNode[] = []
  if (showKeys && node.key !== '') {
    parts.push(h('span', { class: 'key' }, String(node.key)), h('span', null, ': '))
  }

  if (node.status === 'updated') {
    return h('li', { class: 'updated-line' }, [
      ...parts,
      h('span', { class: ['value', 'removed'] }, formatValue(node.oldValue)),
      h('span', { class: ['value', 'added'] }, formatValue(node.value)),
      ',',
    ])
  }

  if (node.type === 'array' || node.type === 'object') {
    const open = node.type === 'array' ? '[' : '{'
    const close = node.type === 'array' ? ']' : '}'
    const children = (node.children ?? []).map((c) => renderDiff(c, node.type === 'object'))
    return h('li', null, [
      h('div', { class: [node.type, node.status] }, [
        ...parts,
        open,
        children.length ? h('ul', null, children) : null,
        `${close},`,
      ]),
    ])
  }

  const valueToDisplay = node.status === 'removed' ? node.oldValue : node.value
  return h('li', null, [
    h('span', { class: [node.status, 'result'] }, [
      ...parts,
      h('span', { class: ['value', node.status] }, formatValue(valueToDisplay)),
    ]),
    ',',
  ])
}

const rootVNode = computed(() => (result.value?.root ? renderDiff(result.value.root, false) : null))
</script>

<template>
  <n-card class="card">
    <ToolTextarea v-model:value="rawLeftJson" label="第一个 JSON" :rows="10" placeholder="在此粘贴第一个 JSON…" />
    <div class="spacer" />
    <ToolTextarea v-model:value="rawRightJson" label="要比较的 JSON" :rows="10" placeholder="在此粘贴要比较的 JSON…" />

    <div v-if="errorText" class="error">
      <n-text type="error">{{ errorText }}</n-text>
    </div>

    <div v-if="result" class="diff-area">
      <div class="switch-row">
        <n-form-item label="仅显示差异" label-placement="left" class="switch-item">
          <n-switch v-model:value="onlyShowDifferences" />
        </n-form-item>
      </div>

      <div class="diff-viewer" v-if="rootVNode">
        <ul>{{ rootVNode }}</ul>
      </div>
      <div v-else class="same-text">
        <n-text depth="3">两个 JSON 相同</n-text>
      </div>
    </div>
  </n-card>
</template>

<style scoped>
.card {
  min-width: 480px;
  width: 100%;
}

.spacer {
  height: 16px;
}

.error {
  margin: 8px 0;
}

.diff-area {
  margin-top: 16px;
}

.switch-row {
  display: flex;
  justify-content: center;
  margin-bottom: 8px;
}

.switch-item {
  margin-bottom: 0;
}

.diff-viewer {
  color: v-bind('themeVars.textColor2');
  font-family: 'Cascadia Code', monospace;
  font-size: 13px;
}

.diff-viewer ul {
  list-style: none;
  padding-left: 20px;
  margin: 0;
}

.diff-viewer > ul {
  padding-left: 0;
}

.updated-line {
  padding: 3px 0;
}

.result,
.array,
.object,
.value {
  &:not(:last-child) {
    margin-right: 4px;
  }

  &.added {
    padding: 3px 5px;
    border-radius: 4px;
    background-color: color-mix(in srgb, v-bind('themeVars.successColor') 18%, transparent);
    color: v-bind('themeVars.successColor');
  }

  &.removed {
    padding: 3px 5px;
    border-radius: 4px;
    background-color: color-mix(in srgb, v-bind('themeVars.errorColor') 18%, transparent);
    color: v-bind('themeVars.errorColor');
  }
}

.added .added,
.removed .removed {
  background-color: transparent;
  color: inherit;
}

.key {
  font-weight: 600;
  color: v-bind('themeVars.textColor1');
}

.same-text {
  text-align: center;
  padding: 16px 0;
}
</style>