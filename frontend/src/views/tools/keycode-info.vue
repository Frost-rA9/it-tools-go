<script setup lang="ts">
import { computed, ref } from 'vue'
import { NCard, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useEventListener, useClipboard } from '@vueuse/core'

const message = useMessage()
const themeVars = useThemeVars()

const event = ref<KeyboardEvent>()

useEventListener(document, 'keydown', (e) => {
  event.value = e
})

const fields = computed(() => {
  if (!event.value) return []
  const e = event.value
  return [
    { label: 'Key', value: e.key, placeholder: '按键名称…' },
    { label: 'Keycode', value: String(e.keyCode), placeholder: '按键码…' },
    { label: 'Code', value: e.code, placeholder: '物理键代码…' },
    { label: 'Location', value: String(e.location), placeholder: '位置…' },
    {
      label: 'Modifiers',
      value: [e.metaKey && 'Meta', e.shiftKey && 'Shift', e.ctrlKey && 'Ctrl', e.altKey && 'Alt']
        .filter(Boolean)
        .join(' + '),
      placeholder: 'None',
    },
  ]
})

const { copy } = useClipboard()

async function copyValue(label: string, value: string) {
  await copy(value)
  message.success(`${label} 已复制到剪贴板`)
}
</script>

<template>
  <n-card class="card">
    <div class="press-area">
      <div class="pressed-key">{{ event?.key ?? '—' }}</div>
      <div class="hint">按下键盘上的按键以查看其信息</div>
    </div>

    <div class="prop-row" v-for="f in fields" :key="f.label">
      <span class="prop-label">{{ f.label }}</span>
      <span class="prop-value">{{ f.value || '—' }}</span>
      <n-button size="tiny" secondary :disabled="!f.value" @click="copyValue(f.label, f.value)">
        复制
      </n-button>
    </div>
  </n-card>
</template>

<style scoped>
.card {
  min-width: 400px;
}

.press-area {
  text-align: center;
  padding: 28px 0;
  margin-bottom: 16px;
  border-radius: 6px;
  background-color: color-mix(in srgb, v-bind('themeVars.primaryColor') 8%, transparent);
}

.pressed-key {
  font-size: 44px;
  font-weight: 600;
  line-height: 1;
  margin-bottom: 10px;
  color: v-bind('themeVars.textColor1');
}

.hint {
  opacity: 0.7;
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
}

.prop-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
}

.prop-label {
  flex: 0 0 90px;
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
}

.prop-value {
  flex: 1 1 auto;
  font-family: 'Cascadia Code', monospace;
  font-size: 14px;
  word-break: break-all;
  color: v-bind('themeVars.textColor1');
}
</style>