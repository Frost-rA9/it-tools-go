<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInput, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const themeVars = useThemeVars()

interface Analysis {
  password_length: number
  entropy: number
  charset_length: number
  crack_duration: string
  seconds_to_crack: number
  score: number
}

const password = ref('')
const analysis = ref<Analysis>({
  password_length: 0,
  entropy: 0,
  charset_length: 0,
  crack_duration: '',
  seconds_to_crack: 0,
  score: 0,
})

async function run() {
  try {
    const output = await RunTool('password-strength-analyser', JSON.stringify({ password: password.value }))
    analysis.value = JSON.parse(output)
  } catch (e) {
    analysis.value = { password_length: 0, entropy: 0, charset_length: 0, crack_duration: '', seconds_to_crack: 0, score: 0 }
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch(password, () => debouncedRun(), { immediate: true })

const details = [
  { label: '密码长度：', key: 'password_length' as const },
  { label: '熵：', key: 'entropy' as const },
  { label: '字符集大小：', key: 'charset_length' as const },
  { label: '得分：', key: 'score' as const },
]

function formatValue(key: string, value: number): string {
  if (key === 'entropy') return String(Math.round(value * 100) / 100)
  if (key === 'score') return `${Math.round(value * 100)} / 100`
  return String(value)
}
</script>

<template>
  <div class="pwd-tool">
    <n-card title="密码强度分析仪" class="card">
      <div class="field">
        <div class="field-label">密码</div>
        <n-input v-model:value="password" type="password" show-password-on="click" clearable placeholder="输入要分析的密码…" />
      </div>

      <div class="crack-card">
        <div class="crack-label">暴力破解此密码所需时长</div>
        <div class="crack-duration">{{ analysis.crack_duration }}</div>
      </div>

      <div class="detail-card">
        <div v-for="d in details" :key="d.key" class="detail-row">
          <span class="detail-label">{{ d.label }}</span>
          <span class="detail-value">{{ formatValue(d.key, analysis[d.key]) }}</span>
        </div>
      </div>

      <div class="note">
        <b>注意：</b>强度基于暴力破解所需时间估算，未考虑字典攻击的可能。
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.pwd-tool {
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

.crack-card {
  border: 1px solid v-bind('themeVars.borderColor');
  border-radius: 3px;
  padding: 16px;
  text-align: center;
  margin-bottom: 16px;
}

.crack-label {
  font-size: 14px;
  opacity: 0.6;
  margin-bottom: 8px;
}

.crack-duration {
  font-size: 24px;
  font-weight: 500;
}

.detail-card {
  border: 1px solid v-bind('themeVars.borderColor');
  border-radius: 3px;
  padding: 12px 16px;
  margin-bottom: 16px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 4px 0;
}

.detail-label {
  opacity: 0.6;
}

.note {
  font-size: 13px;
  opacity: 0.7;
}
</style>
