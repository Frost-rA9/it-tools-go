<script setup lang="ts">
import { computed } from 'vue'
import { NCard, NIcon, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useClipboard } from '@vueuse/core'
import hljs from 'highlight.js/lib/core'
import sql from 'highlight.js/lib/languages/sql'
import xml from 'highlight.js/lib/languages/xml'
import yaml from 'highlight.js/lib/languages/yaml'
import json from 'highlight.js/lib/languages/json'
import ini from 'highlight.js/lib/languages/ini'
import plaintext from 'highlight.js/lib/languages/plaintext'
import { Copy } from '@vicons/tabler'
import { useUiStore } from '../stores/ui'

hljs.registerLanguage('sql', sql)
hljs.registerLanguage('xml', xml)
hljs.registerLanguage('yaml', yaml)
hljs.registerLanguage('json', json)
hljs.registerLanguage('ini', ini)
hljs.registerLanguage('plaintext', plaintext)

const props = withDefaults(
  defineProps<{
    label?: string
    value: string
    language?: string
    rows?: number
  }>(),
  {
    label: undefined,
    language: 'txt',
    rows: 10,
  },
)

const message = useMessage()
const themeVars = useThemeVars()
const ui = useUiStore()
const { copy } = useClipboard()

const minHeight = computed(() => `${props.rows * 22}px`)
const themeClass = computed(() => (ui.isDarkTheme ? 'code-dark' : 'code-light'))

async function copyValue() {
  await copy(props.value)
  message.success('已复制到剪贴板')
}
</script>

<template>
  <div class="code-output" :class="themeClass">
    <div v-if="label" class="label">{{ label }}</div>
    <n-card :bordered="true" class="output-card" :content-style="{ padding: '8px 12px' }">
      <div class="code-body" :style="{ minHeight }">
        <n-scrollbar x-scrollable>
          <pre class="hljs"><code v-html="hljs.highlight(value, { language, ignoreIllegals: true }).value" /></pre>
        </n-scrollbar>
      </div>
      <div v-if="value" class="copy-btn">
        <n-button circle quaternary size="small" title="复制结果" @click="copyValue">
          <n-icon :component="Copy" :size="16" />
        </n-button>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.code-output {
  width: 100%;
}

.label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.output-card {
  position: relative;
}

.code-body {
  max-height: 420px;
}

.code-body :deep(pre) {
  margin: 0;
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  background: transparent;
  color: v-bind('themeVars.textColor1');
}

.code-body :deep(code) {
  background: transparent;
  font-family: inherit;
  font-size: inherit;
  line-height: inherit;
}

.copy-btn {
  position: absolute;
  top: 8px;
  right: 8px;
}
</style>
