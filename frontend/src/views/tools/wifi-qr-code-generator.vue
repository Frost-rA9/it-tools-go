<script setup lang="ts">
import { ref, watch } from 'vue'
import { NButton, NCard, NCheckbox, NColorPicker, NForm, NFormItem, NImage, NInput, NSelect, NText, useMessage } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const encryption = ref('WPA')
const ssid = ref('')
const password = ref('')
const hidden = ref(false)
const eapMethod = ref('')
const anonymous = ref(false)
const identity = ref('')
const phase2 = ref('')
const foreground = ref('#000000ff')
const background = ref('#ffffffff')
const dataURL = ref('')

const encryptionOptions = [
  { label: '无密码', value: 'nopass' },
  { label: 'WPA/WPA2', value: 'WPA' },
  { label: 'WEP', value: 'WEP' },
  { label: 'WPA2-EAP', value: 'WPA2-EAP' },
]
const eapOptions = ['MD5', 'POTP', 'GTC', 'TLS', 'IKEv2', 'SIM', 'AKA', "AKA'", 'TTLS', 'PWD', 'LEAP', 'PSK', 'FAST', 'TEAP', 'EKE', 'NOOB', 'PEAP']
  .map((value) => ({ label: value, value }))
const phase2Options = ['None', 'MSCHAPV2'].map((value) => ({ label: value, value }))

async function generate() {
  if (!ssid.value) {
    dataURL.value = ''
    return
  }
  try {
    const output = await RunTool('wifi-qr-code-generator', JSON.stringify({
      ssid: ssid.value,
      password: password.value,
      encryption: encryption.value,
      eap_method: eapMethod.value,
      hidden: hidden.value,
      anonymous: anonymous.value,
      identity: identity.value,
      phase2: phase2.value,
      foreground: foreground.value,
      background: background.value,
    }))
    dataURL.value = JSON.parse(output).data_url
  } catch (e) {
    dataURL.value = ''
    message.error(String(e))
  }
}

const debouncedGenerate = useDebounceFn(generate, 150)
watch(
  [encryption, ssid, password, hidden, eapMethod, anonymous, identity, phase2, foreground, background],
  () => debouncedGenerate(),
  { immediate: true },
)

function download() {
  if (!dataURL.value) return
  const link = document.createElement('a')
  link.href = dataURL.value
  link.download = 'wifi-qr-code.png'
  link.click()
}
</script>

<template>
  <n-card class="card">
    <div class="grid">
      <div>
        <n-form label-placement="left" label-width="130">
          <n-form-item label="加密方式">
            <n-select v-model:value="encryption" :options="encryptionOptions" />
          </n-form-item>
          <n-form-item label="SSID">
            <div class="inline-field">
              <n-input v-model:value="ssid" placeholder="WiFi 名称…" />
              <n-checkbox v-model:checked="hidden">隐藏 SSID</n-checkbox>
            </div>
          </n-form-item>
          <n-form-item v-if="encryption !== 'nopass'" label="密码">
            <n-input v-model:value="password" type="password" placeholder="WiFi 密码…" />
          </n-form-item>
          <n-form-item v-if="encryption === 'WPA2-EAP'" label="EAP 方法">
            <n-select v-model:value="eapMethod" :options="eapOptions" filterable />
          </n-form-item>
          <n-form-item v-if="encryption === 'WPA2-EAP'" label="Identity">
            <div class="inline-field">
              <n-input v-model:value="identity" placeholder="EAP Identity…" />
              <n-checkbox v-model:checked="anonymous">Anonymous</n-checkbox>
            </div>
          </n-form-item>
          <n-form-item v-if="encryption === 'WPA2-EAP'" label="EAP Phase 2">
            <n-select v-model:value="phase2" :options="phase2Options" filterable />
          </n-form-item>
          <n-form-item label="前景色">
            <n-color-picker v-model:value="foreground" :modes="['hex']" />
          </n-form-item>
          <n-form-item label="背景色">
            <n-color-picker v-model:value="background" :modes="['hex']" />
          </n-form-item>
        </n-form>
      </div>
      <div class="preview" v-if="dataURL">
        <n-image :src="dataURL" width="200" alt="wifi-qrcode" />
        <n-button @click="download">下载二维码</n-button>
      </div>
      <div v-else class="empty">
        <n-text depth="3">填写 SSID 后生成二维码</n-text>
      </div>
    </div>
  </n-card>
</template>

<style scoped>
.card { min-width: 500px; }
.grid { display: grid; grid-template-columns: minmax(0, 1fr) 230px; gap: 24px; align-items: center; }
.inline-field { display: flex; align-items: center; gap: 10px; width: 100%; }
.inline-field > .n-input { flex: 1; min-width: 0; }
.preview { display: flex; flex-direction: column; align-items: center; gap: 12px; }
.empty { text-align: center; padding: 24px 0; }
@media (max-width: 700px) {
  .card { min-width: 400px; }
  .grid { grid-template-columns: 1fr; }
}
</style>
