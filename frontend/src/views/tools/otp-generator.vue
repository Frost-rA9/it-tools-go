<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { NCard, NGrid, NGi, NButton, NInput, NProgress, NImage, NText, useMessage, useThemeVars } from 'naive-ui'
import { Refresh } from '@vicons/tabler'
import { useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

interface OtpCodes {
  previous: string
  current: string
  next: string
  epoch: number
  counter: number
  secret_hex: string
  next_in: number
}

const message = useMessage()
const themeVars = useThemeVars()

const secret = ref('')
const codes = ref<OtpCodes | null>(null)
const qrDataUrl = ref('')
const uri = ref('')
const nowMs = ref(Date.now())
const secretError = ref('')

function validateSecret(s: string): boolean {
  return /^[A-Za-z2-7]+$/.test(s.trim())
}

async function generateSecret() {
  try {
    const output = await RunTool('otp-generator', JSON.stringify({ action: 'generate' }))
    secret.value = JSON.parse(output).secret
  } catch (e) {
    message.error(String(e))
  }
}

async function fetchCodes() {
  if (!validateSecret(secret.value)) {
    codes.value = null
    return
  }
  try {
    const output = await RunTool(
      'otp-generator',
      JSON.stringify({ action: 'codes', secret: secret.value.trim().toUpperCase(), now: nowMs.value }),
    )
    codes.value = JSON.parse(output)
  } catch (e) {
    codes.value = null
  }
}

async function fetchUriQr() {
  if (!validateSecret(secret.value)) {
    qrDataUrl.value = ''
    uri.value = ''
    return
  }
  const payload = { action: '', secret: secret.value.trim().toUpperCase() }
  try {
    const u = await RunTool('otp-generator', JSON.stringify({ ...payload, action: 'uri' }))
    uri.value = JSON.parse(u).uri
    const q = await RunTool('otp-generator', JSON.stringify({ ...payload, action: 'qr' }))
    qrDataUrl.value = JSON.parse(q).qr_data_url
  } catch (e) {
    message.error(String(e))
  }
}

let timer: number | undefined

onMounted(() => {
  timer = window.setInterval(() => {
    nowMs.value = Date.now()
    fetchCodes()
  }, 1000)
  if (!secret.value) generateSecret()
})

onBeforeUnmount(() => {
  if (timer !== undefined) window.clearInterval(timer)
})

watch(
  secret,
  () => {
    secretError.value =
      secret.value.trim() && !validateSecret(secret.value) ? 'Secret 应为 base32 字符串（A-Z 与 2-7）' : ''
    fetchCodes()
    fetchUriQr()
  },
  { immediate: true },
)

const epoch = computed(() => Math.floor(nowMs.value / 1000))
const interval = computed(() => (nowMs.value / 1000) % 30)
const progress = computed(() => (interval.value / 30) * 100)
const nextIn = computed(() => String(Math.floor(30 - interval.value)).padStart(2, '0'))
const counterHex = computed(() => (codes.value?.counter ?? Math.floor(epoch.value / 30)).toString(16).padStart(16, '0'))

const { copy } = useClipboard()

async function copyCode(label: string, code: string) {
  await copy(code)
  message.success(`${label} 已复制到剪贴板`)
}
</script>

<template>
  <n-card class="card">
    <n-grid :cols="1" responsive="screen" class="inner-grid">
      <n-gi>
        <n-input v-model:value="secret" placeholder="粘贴 TOTP Secret…" clearable>
          <template #suffix>
            <n-button quaternary circle size="small" :title="'生成新的随机 Secret'" @click="generateSecret">
              <n-icon :component="Refresh" />
            </n-button>
          </template>
        </n-input>
        <div v-if="secretError" class="error">
          <n-text type="error">{{ secretError }}</n-text>
        </div>

        <div class="tokens">
          <div class="token-labels">
            <span>Previous</span>
            <span class="center">Current OTP</span>
            <span class="right">Next</span>
          </div>
          <div class="token-row">
            <n-button class="token previous" @click="copyCode('Previous OTP', codes?.previous ?? '')">
              {{ codes?.previous ?? '—' }}
            </n-button>
            <n-button class="token current" @click="copyCode('Current OTP', codes?.current ?? '')">
              {{ codes?.current ?? '—' }}
            </n-button>
            <n-button class="token next" @click="copyCode('Next OTP', codes?.next ?? '')">
              {{ codes?.next ?? '—' }}
            </n-button>
          </div>
        </div>

        <n-progress type="line" :percentage="progress" :show-indicator="false" :color="themeVars.primaryColor" class="progress" />
        <div class="next-in">Next in {{ nextIn }}s</div>

        <div class="qr-area">
          <n-image v-if="qrDataUrl" :src="qrDataUrl" width="210" />
          <n-button v-if="uri" tag="a" :href="uri" target="_blank" class="uri-btn">打开 Key URI</n-button>
        </div>
      </n-gi>

      <n-gi>
        <div class="info-item">
          <span class="info-label">Secret in hexadecimal</span>
          <div class="info-value mono">{{ codes?.secret_hex ?? '—' }}</div>
        </div>
        <div class="info-item">
          <span class="info-label">Epoch</span>
          <div class="info-value mono">{{ codes?.epoch ?? epoch }}</div>
        </div>
        <div class="info-item">
          <span class="info-label">Iteration</span>
          <div class="info-value mono">Count: {{ codes?.counter ?? Math.floor(epoch / 30) }}</div>
          <div class="info-value mono">Padded hex: {{ counterHex }}</div>
        </div>
      </n-gi>
    </n-grid>
  </n-card>
</template>

<style scoped>
.card {
  min-width: 420px;
  width: 100%;
}

.inner-grid {
  max-width: 560px;
  margin: 0 auto;
}

.error {
  margin-top: 6px;
}

.tokens {
  margin: 18px 0 10px;
}

.token-labels {
  display: flex;
  width: 100%;
  margin-bottom: 6px;
  font-size: 13px;
  opacity: 0.8;
  color: v-bind('themeVars.textColor3');
}

.token-labels .center {
  flex: 1;
  text-align: center;
}

.token-labels .right {
  text-align: right;
}

.token-row {
  display: flex;
  width: 100%;
  gap: 0;
}

.token {
  flex: 1;
  height: 44px;
  font-family: 'Cascadia Code', monospace;
  font-size: 14px;
}

.token.current {
  font-size: 19px;
  border-left: 1px solid color-mix(in srgb, v-bind('themeVars.textColor1') 30%, transparent);
  border-right: 1px solid color-mix(in srgb, v-bind('themeVars.textColor1') 30%, transparent);
  border-radius: 0;
}

.token.previous {
  border-radius: 6px 0 0 6px;
}

.token.next {
  border-radius: 0 6px 6px 0;
}

.progress {
  margin-top: 10px;
}

.next-in {
  text-align: center;
  font-size: 13px;
  opacity: 0.7;
  color: v-bind('themeVars.textColor2');
}

.qr-area {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.uri-btn {
  margin-top: 4px;
}

.info-item {
  margin-bottom: 20px;
}

.info-label {
  display: block;
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.info-value {
  padding: 8px 10px;
  border-radius: 4px;
  background-color: color-mix(in srgb, v-bind('themeVars.textColor1') 6%, transparent);
  margin-bottom: 6px;
  font-size: 13px;
  word-break: break-all;
  color: v-bind('themeVars.textColor1');
}

.mono {
  font-family: 'Cascadia Code', monospace;
}
</style>