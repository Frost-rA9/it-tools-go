<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { NCard, NInput, NButton, NAlert, NSpace, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'
import ToolCodeBlock from '../../components/ToolCodeBlock.vue'

const message = useMessage()
const themeVars = useThemeVars()

const macAddress = ref('d2:5f:61:07:3d:63')

const ula = ref('')
const firstBlock = ref('')
const lastBlock = ref('')
const errorMessage = ref('')

// 对齐 it-tools macAddressValidation：2-5 组“XX:/-”加末组 XX（3-6 字节）。
const MAC_RE = /^([0-9A-Fa-f]{2}[:-]){2,5}([0-9A-Fa-f]{2})$/
const macValid = computed(() => MAC_RE.test(macAddress.value.trim()))

async function run() {
  if (!macValid.value) {
    ula.value = firstBlock.value = lastBlock.value = ''
    errorMessage.value = `MAC 地址格式无效：${macAddress.value}（期望如 d2:5f:61:07:3d:63）`
    return
  }
  try {
    const output = await RunTool(
      'ipv6-ula-generator',
      JSON.stringify({ macAddress: macAddress.value.trim() }),
    )
    const parsed = JSON.parse(output)
    ula.value = parsed.ula
    firstBlock.value = parsed.firstRoutableBlock
    lastBlock.value = parsed.lastRoutableBlock
    errorMessage.value = ''
  } catch (e) {
    ula.value = firstBlock.value = lastBlock.value = ''
    errorMessage.value = String(e)
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch(macAddress, () => debouncedRun(), { immediate: true })

// 重新生成：直接调用后端（新时间戳 → 新 ULA），不受防抖影响
function regenerate() {
  run()
}
</script>

<template>
  <div class="tool-page">
    <n-card title="IPv6 ULA 生成器" class="card">
      <n-alert type="info" class="info-alert">
        本工具采用 IETF 建议的 RFC 4193 方法 1：将当前时间戳与 MAC 地址拼接后做 SHA1 哈希，取其低
        40 bits 生成随机的唯一本地地址（ULA）前缀。地址前缀固定为 <span class="mono-inline">fc00::/7</span>（本地分配位为 1，即
        <span class="mono-inline">fd</span>）。
      </n-alert>

      <div class="field">
        <div class="field-label">MAC 地址</div>
        <n-input
          v-model:value="macAddress"
          placeholder="输入 MAC 地址，如 d2:5f:61:07:3d:63"
          clearable
          class="mono-input"
        />
      </div>

      <n-alert v-if="errorMessage" type="error" class="error-alert">
        {{ errorMessage }}
      </n-alert>

      <template v-if="ula">
        <ToolCodeBlock
          label="IPv6 ULA 前缀"
          :value="ula"
          align="center"
          copyable
        />
        <ToolCodeBlock
          label="首个可路由块（/64）"
          :value="firstBlock"
          align="center"
          copyable
        />
        <ToolCodeBlock
          label="末个可路由块（/64）"
          :value="lastBlock"
          align="center"
          copyable
        />
      </template>

      <n-space justify="center">
        <n-button type="primary" :disabled="!macValid" @click="regenerate">重新生成</n-button>
      </n-space>
    </n-card>
  </div>
</template>

<style scoped>
.tool-page {
  width: 100%;
}

.card {
  width: 100%;
}

.info-alert {
  margin-bottom: 16px;
}

.field {
  margin-bottom: 16px;
}

.field-label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.mono-inline {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}

.mono-input :deep(input) {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}

.error-alert {
  margin-bottom: 16px;
}
</style>