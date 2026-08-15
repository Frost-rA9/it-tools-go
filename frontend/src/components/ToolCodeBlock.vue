<script setup lang="ts">
import { useThemeVars } from 'naive-ui'

const themeVars = useThemeVars()

withDefaults(
  defineProps<{
    label?: string
    value?: string
    align?: 'left' | 'center'
  }>(),
  {
    label: undefined,
    value: '',
    align: 'left',
  },
)
</script>

<template>
  <div class="tool-code-block">
    <div v-if="label" class="label">{{ label }}</div>
    <div class="block">
      <pre class="mono" :class="{ center: align === 'center' }">{{ value }}</pre>
    </div>
  </div>
</template>

<style scoped>
.tool-code-block {
  width: 100%;
  margin-bottom: 16px;
}

.label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.block {
  border: 1px solid v-bind('themeVars.borderColor');
  border-radius: 3px;
  background-color: v-bind('themeVars.inputColor');
  padding: 8px 12px;
  color: v-bind('themeVars.textColor1');
}

.mono {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
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
</style>
