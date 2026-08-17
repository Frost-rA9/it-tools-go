<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NButton, NInput, NStatistic, NScrollbar, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const username = ref('')
const password = ref('')
const header = ref('')

async function runGenerate() {
  try {
    const output = await RunTool(
      'basic-auth-generator',
      JSON.stringify({ username: username.value, password: password.value }),
    )
    header.value = JSON.parse(output).header
  } catch (e) {
    header.value = ''
    message.error(String(e))
  }
}

const debouncedGenerate = useDebounceFn(runGenerate, 150)
watch([username, password], () => debouncedGenerate(), { immediate: true })

const { copy } = useClipboard()

async function copyHeader() {
  await copy(header.value)
  message.success('Authorization 头已复制到剪贴板')
}
</script>

<template>
  <n-card class="card">
    <div class="field">
      <span class="field-label">Username</span>
      <n-input v-model:value="username" placeholder="你的用户名…" clearable />
    </div>

    <div class="field">
      <span class="field-label">Password</span>
      <n-input v-model:value="password" placeholder="你的密码…" clearable type="password" />
    </div>

    <div class="header-card">
      <n-statistic label="Authorization header:" class="header-statistic">
        <n-scrollbar x-scrollable class="header-scroll">
          {{ header }}
        </n-scrollbar>
      </n-statistic>
    </div>

    <div class="copy-row">
      <n-button type="primary" @click="copyHeader">复制 Authorization 头</n-button>
    </div>
  </n-card>
</template>

<style scoped>
.card {
  min-width: 420px;
  width: 100%;
}

.field {
  margin-bottom: 20px;
}

.field-label {
  display: block;
  font-size: 14px;
  margin-bottom: 8px;
  color: v-bind('themeVars.textColor2');
}

.header-card {
  margin-top: 4px;
  padding: 12px 16px;
  border-radius: 6px;
  background-color: color-mix(in srgb, v-bind('themeVars.primaryColor') 8%, transparent);
}

.header-statistic :deep(.n-statistic-value__content) {
  font-family: 'Cascadia Code', monospace;
  font-size: 17px;
  white-space: nowrap;
  color: v-bind('themeVars.textColor1');
}

.header-scroll {
  max-width: 100%;
}

.copy-row {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>