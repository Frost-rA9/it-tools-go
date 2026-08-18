<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NAlert, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import CodeOutput from '../../components/CodeOutput.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const themeVars = useThemeVars()

const dockerRun = ref('docker run -p 80:80 -v /var/run/docker.sock:/tmp/docker.sock:ro --restart always --log-opt max-size=1g nginx')
const composeOutput = ref('')
const warnings = ref<string[]>([])
const errorMessage = ref('')

async function run() {
  try {
    const output = await RunTool('docker-run-to-docker-compose-converter', JSON.stringify({ text: dockerRun.value }))
    const result = JSON.parse(output)
    composeOutput.value = result.compose
    warnings.value = result.warnings ?? []
    errorMessage.value = ''
  } catch (e) {
    composeOutput.value = ''
    warnings.value = []
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 300)
watch(dockerRun, () => debouncedRun(), { immediate: true })
</script>

<template>
  <n-card title="Docker Run 转 Compose — 输入" class="tool-card">
    <ToolTextarea v-model:value="dockerRun" label="docker run 命令" :rows="10" placeholder="docker run -p 80:80 nginx…" monospace />

    <n-alert v-if="errorMessage" type="error" class="error-alert">
      {{ errorMessage }}
    </n-alert>
  </n-card>

  <n-card title="Docker Run 转 Compose — 输出" class="tool-card">
    <CodeOutput label="docker-compose.yml" :value="composeOutput" language="yaml" :rows="20" />

    <n-alert v-if="warnings.length > 0" type="warning" class="warn-alert">
      <template #header>以下选项未转换到 docker-compose</template>
      <ul class="warn-list">
        <li v-for="(w, i) in warnings" :key="i">{{ w }}</li>
      </ul>
    </n-alert>
  </n-card>
</template>

<style scoped>
.error-alert {
  margin-bottom: 16px;
}

.warn-alert {
  margin-top: 16px;
}

.warn-list {
  margin: 4px 0 0;
  padding-left: 20px;
}

.warn-list li {
  font-size: 13px;
  color: v-bind('themeVars.textColor2');
}
</style>
