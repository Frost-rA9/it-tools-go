<script setup lang="ts">
import { ref, watch } from 'vue'
import { NCard, NTable, NText, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

interface JwtClaim {
  claim: string
  value: string
  claimDescription?: string
  friendlyValue?: string
}

interface JwtResult {
  header: JwtClaim[]
  payload: JwtClaim[]
}

const message = useMessage()
const themeVars = useThemeVars()

const rawJwt = ref(
  'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c',
)
const decoded = ref<JwtResult | null>(null)
const errorText = ref('')

const sections: { key: 'header' | 'payload'; title: string }[] = [
  { key: 'header', title: 'Header' },
  { key: 'payload', title: 'Payload' },
]

async function runDecode() {
  errorText.value = ''
  try {
    const output = await RunTool('jwt-parser', JSON.stringify({ jwt: rawJwt.value }))
    decoded.value = JSON.parse(output)
  } catch (e) {
    decoded.value = null
    errorText.value = String(e).replace(/^Error:\s*/, '')
  }
}

const debouncedDecode = useDebounceFn(runDecode, 150)
watch(rawJwt, () => debouncedDecode(), { immediate: true })
</script>

<template>
  <n-card class="card">
    <ToolTextarea v-model:value="rawJwt" label="要解码的 JWT" :rows="5" placeholder="将 token 粘贴到这里…" />

    <div v-if="errorText" class="error">
      <n-text type="error">{{ errorText }}</n-text>
    </div>

    <n-table v-if="decoded" :bordered="false" class="table">
      <tbody>
        <template v-for="section in sections" :key="section.key">
          <tr>
            <td colspan="2" class="table-header">{{ section.title }}</td>
          </tr>
          <tr v-for="(c, i) in decoded[section.key]" :key="section.key + i">
            <td class="claim-cell">
              <span class="claim">{{ c.claim }}</span>
              <span v-if="c.claimDescription" class="claim-desc">({{ c.claimDescription }})</span>
            </td>
            <td class="value-cell">
              <pre class="claim-value">{{ c.value }}</pre>
              <span v-if="c.friendlyValue" class="friendly">({{ c.friendlyValue }})</span>
            </td>
          </tr>
        </template>
      </tbody>
    </n-table>
  </n-card>
</template>

<style scoped>
.card {
  min-width: 480px;
  width: 100%;
}

.error {
  margin: 8px 0 16px;
}

.table-header {
  text-align: center;
  font-weight: 600;
  color: v-bind('themeVars.textColor1');
  background-color: v-bind('themeVars.tableHeaderColor');
}

.claim-cell {
  vertical-align: top;
  width: 30%;
  padding: 8px 12px;
}

.claim {
  font-weight: 600;
  color: v-bind('themeVars.primaryColor');
}

.claim-desc {
  margin-left: 8px;
  opacity: 0.7;
  color: v-bind('themeVars.textColor3');
  font-size: 12px;
}

.value-cell {
  vertical-align: top;
  word-break: break-all;
  padding: 8px 12px;
}

.claim-value {
  margin: 0;
  font-family: 'Cascadia Code', monospace;
  font-size: 13px;
  white-space: pre-wrap;
  color: v-bind('themeVars.textColor1');
}

.friendly {
  margin-left: 8px;
  opacity: 0.7;
  color: v-bind('themeVars.textColor3');
  font-size: 12px;
}
</style>