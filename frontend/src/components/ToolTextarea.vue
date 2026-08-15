<script setup lang="ts">
import { NInput, useThemeVars } from 'naive-ui'

const themeVars = useThemeVars()
const value = defineModel<string>('value', { default: '' })

withDefaults(
  defineProps<{
    label?: string
    placeholder?: string
    rows?: number
    readonly?: boolean
    resizable?: boolean
    autosize?: boolean
    monospace?: boolean
  }>(),
  {
    label: undefined,
    placeholder: '',
    rows: 3,
    readonly: false,
    resizable: true,
    autosize: false,
    monospace: false,
  },
)
</script>

<template>
  <div class="tool-textarea">
    <div v-if="label" class="label">{{ label }}</div>
    <n-input
      v-model:value="value"
      type="textarea"
      :rows="rows"
      :placeholder="placeholder"
      :readonly="readonly"
      :resizable="resizable"
      :autosize="autosize"
      :class="{ monospace }"
    />
  </div>
</template>

<style scoped>
.tool-textarea {
  width: 100%;
  margin-bottom: 16px;
}

.label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.monospace :deep(textarea) {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}
</style>
