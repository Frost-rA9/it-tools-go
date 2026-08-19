<script setup lang="ts">
import { ref } from 'vue'
import { NIcon, useMessage, useThemeVars } from 'naive-ui'
import { Copy } from '@vicons/tabler'
import { useClipboard } from '@vueuse/core'

const themeVars = useThemeVars()
const message = useMessage()

const props = withDefaults(
  defineProps<{
    label?: string
    value?: string
    align?: 'left' | 'center'
    copyable?: boolean
  }>(),
  {
    label: undefined,
    value: '',
    align: 'left',
    copyable: false,
  },
)

const copySource = ref('')
const { copy } = useClipboard({ source: copySource })

async function copyValue() {
  copySource.value = props.value
  await copy()
  message.success('已复制到剪贴板')
}
</script>

<template>
  <div class="tool-code-block">
    <div v-if="label" class="label">{{ label }}</div>
    <div class="block">
      <pre class="mono" :class="{ center: align === 'center' }">{{ value }}</pre>
      <button v-if="copyable && value" class="copy-btn" title="复制" @click="copyValue">
        <n-icon :component="Copy" :size="16" />
      </button>
    </div>
  </div>
</template>

<style scoped>
.tool-code-block {
  position: relative;
  width: 100%;
  margin-bottom: 16px;
}

.label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.block {
  position: relative;
  border: 1px solid v-bind('themeVars.borderColor');
  border-radius: 3px;
  background-color: v-bind('themeVars.inputColor');
  padding: 8px 12px;
  color: v-bind('themeVars.textColor1');
}

.mono {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
  font-variant-ligatures: none;
  font-feature-settings: 'liga' 0, 'calt' 0;
  font-size: 14px;
  line-height: 1.6;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
}

.center {
  width: fit-content;
  margin-left: auto;
  margin-right: auto;
}

.copy-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 50%;
  background: transparent;
  color: v-bind('themeVars.textColor3');
  cursor: pointer;
  transition: background-color 0.2s, color 0.2s;
}

.copy-btn:hover {
  background-color: v-bind('themeVars.actionColor');
  color: v-bind('themeVars.primaryColor');
}
</style>
