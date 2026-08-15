<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NInputNumber, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'
import ToolCodeBlock from '../../components/ToolCodeBlock.vue'

const message = useMessage()
const themeVars = useThemeVars()

const bits = ref(2048)
const publicKeyPem = ref('')
const privateKeyPem = ref('')

async function run() {
  try {
    const output = await RunTool('rsa-key-pair-generator', JSON.stringify({ bits: bits.value }))
    const parsed = JSON.parse(output)
    publicKeyPem.value = parsed.public_key_pem
    privateKeyPem.value = parsed.private_key_pem
  } catch (e) {
    publicKeyPem.value = ''
    privateKeyPem.value = ''
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch(bits, () => debouncedRun(), { immediate: true })
</script>

<template>
  <div class="bits-row" style="flex: 0 0 100%">
    <span class="bits-label">密钥位数：</span>
    <n-input-number v-model:value="bits" :min="256" :max="16384" :step="8" style="width: 200px" />
    <n-button type="primary" @click="run">重新生成</n-button>
  </div>

  <n-card title="公钥（Public key）" class="card">
    <ToolCodeBlock :value="publicKeyPem" copyable />
  </n-card>

  <n-card title="私钥（Private key）" class="card">
    <ToolCodeBlock :value="privateKeyPem" copyable />
  </n-card>
</template>

<style scoped>
.card {
  min-width: 400px;
}

.bits-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.bits-label {
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
}
</style>
