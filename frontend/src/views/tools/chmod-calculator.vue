<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { NCard, NCheckbox, NInput, NButton, NAlert, NDescriptions, NDescriptionsItem, NTag, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

interface ChmodResult {
  octal: string
  symbolic: string
  owner: string
  group: string
  others: string
  special: number
  special_text: string
  has_suid: boolean
  has_sgid: boolean
  has_sticky: boolean
}

// 9 个权限位 + 3 个特殊位
const ownerR = ref(true)
const ownerW = ref(true)
const ownerX = ref(true)
const groupR = ref(true)
const groupW = ref(false)
const groupX = ref(true)
const othersR = ref(true)
const othersW = ref(false)
const othersX = ref(true)
const suid = ref(false)
const sgid = ref(false)
const sticky = ref(false)

const modeInput = ref('')
const result = ref<ChmodResult | null>(null)
const errorMessage = ref('')
const syncing = ref(false)

function buildOctal(): string {
  const perm = (r: boolean, w: boolean, x: boolean) => (r ? 4 : 0) + (w ? 2 : 0) + (x ? 1 : 0)
  const special = (suid.value ? 4 : 0) + (sgid.value ? 2 : 0) + (sticky.value ? 1 : 0)
  const oct = `${perm(ownerR.value, ownerW.value, ownerX.value)}${perm(groupR.value, groupW.value, groupX.value)}${perm(othersR.value, othersW.value, othersX.value)}`
  return special > 0 ? `${special}${oct}` : oct
}

function applyResult(o: ChmodResult) {
  const num = Number.parseInt(o.octal, 8)
  const bits = (v: number) => ({
    r: (v & 4) !== 0,
    w: (v & 2) !== 0,
    x: (v & 1) !== 0,
  })
  const owner = bits((num >> 6) & 7)
  const group = bits((num >> 3) & 7)
  const others = bits(num & 7)
  syncing.value = true
  ownerR.value = owner.r
  ownerW.value = owner.w
  ownerX.value = owner.x
  groupR.value = group.r
  groupW.value = group.w
  groupX.value = group.x
  othersR.value = others.r
  othersW.value = others.w
  othersX.value = others.x
  suid.value = o.has_suid
  sgid.value = o.has_sgid
  sticky.value = o.has_sticky
  syncing.value = false
}

async function run(mode: string) {
  try {
    const output = await RunTool('chmod-calculator', JSON.stringify({ mode }))
    result.value = JSON.parse(output)
    errorMessage.value = ''
    if (!syncing.value && result.value) applyResult(result.value)
  } catch (e) {
    result.value = null
    errorMessage.value = String(e)
  }
}

// 复选框变化 → 组装八进制 → 更新输入框（防循环）
watch(
  [ownerR, ownerW, ownerX, groupR, groupW, groupX, othersR, othersW, othersX, suid, sgid, sticky],
  () => {
    if (syncing.value) return
    const oct = buildOctal()
    if (modeInput.value !== oct) modeInput.value = oct
  },
)

// 输入框 → debounce 调用后端
const debouncedRun = useDebounceFn(() => run(modeInput.value), 200)
watch(modeInput, () => {
  if (!modeInput.value) return
  debouncedRun()
})

// 初始执行
run(buildOctal())

const copySource = ref('')
const { copy } = useClipboard({ source: copySource })

// it-tools 风格命令：chmod <octal> path
const chmodCommand = computed(() => (result.value ? `chmod ${result.value.octal} path` : ''))

async function copyMode() {
  copySource.value = chmodCommand.value
  await copy()
  message.success('已复制到剪贴板')
}
</script>

<template>
  <div class="tool-page">
    <n-card title="Chmod 计算器" class="card">
      <div class="field">
        <div class="field-label">输入权限（八进制或符号，如 755 / rwxr-xr-x）</div>
        <n-input v-model:value="modeInput" placeholder="755 或 rwxr-xr-x" class="mono-input" />
      </div>

      <div class="perm-table">
        <div class="perm-row">
          <div class="perm-label">属主 owner</div>
          <div class="perm-checks">
            <n-checkbox v-model:checked="ownerR">读 r</n-checkbox>
            <n-checkbox v-model:checked="ownerW">写 w</n-checkbox>
            <n-checkbox v-model:checked="ownerX">执行 x</n-checkbox>
          </div>
        </div>
        <div class="perm-row">
          <div class="perm-label">属组 group</div>
          <div class="perm-checks">
            <n-checkbox v-model:checked="groupR">读 r</n-checkbox>
            <n-checkbox v-model:checked="groupW">写 w</n-checkbox>
            <n-checkbox v-model:checked="groupX">执行 x</n-checkbox>
          </div>
        </div>
        <div class="perm-row">
          <div class="perm-label">其他 others</div>
          <div class="perm-checks">
            <n-checkbox v-model:checked="othersR">读 r</n-checkbox>
            <n-checkbox v-model:checked="othersW">写 w</n-checkbox>
            <n-checkbox v-model:checked="othersX">执行 x</n-checkbox>
          </div>
        </div>
        <div class="perm-row">
          <div class="perm-label">特殊位</div>
          <div class="perm-checks">
            <n-checkbox v-model:checked="suid">SUID 4000</n-checkbox>
            <n-checkbox v-model:checked="sgid">SGID 2000</n-checkbox>
            <n-checkbox v-model:checked="sticky">Sticky 1000</n-checkbox>
          </div>
        </div>
      </div>

      <n-alert v-if="errorMessage" type="error" class="error-alert">
        {{ errorMessage }}
      </n-alert>

      <template v-if="result">
        <div class="result-hero">
          <n-tag size="large" :bordered="false" round class="hero-tag">{{ result.octal }}</n-tag>
          <span class="hero-symbol mono">{{ result.symbolic }}</span>
        </div>

        <n-descriptions :column="3" bordered size="small" class="desc">
          <n-descriptions-item label="属主">
            <span class="mono">{{ result.owner }}</span>
          </n-descriptions-item>
          <n-descriptions-item label="属组">
            <span class="mono">{{ result.group }}</span>
          </n-descriptions-item>
          <n-descriptions-item label="其他">
            <span class="mono">{{ result.others }}</span>
          </n-descriptions-item>
        </n-descriptions>

        <div class="special-text">
          <span class="field-label">特殊位说明</span>
          <span>{{ result.special_text }}</span>
        </div>

        <div class="copy-row">
          <n-input :value="chmodCommand" readonly class="mono" placeholder="chmod 755 path" />
          <n-button type="primary" @click="copyMode">复制</n-button>
        </div>
      </template>
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

.field {
  margin-bottom: 16px;
}

.field-label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.mono-input :deep(input),
.mono {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}

.perm-table {
  margin-bottom: 16px;
  border: 1px solid v-bind('themeVars.borderColor');
  border-radius: 6px;
  overflow: hidden;
}

.perm-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  border-bottom: 1px solid v-bind('themeVars.borderColor');
}

.perm-row:last-child {
  border-bottom: none;
}

.perm-label {
  flex: 0 0 110px;
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
}

.perm-checks {
  display: flex;
  gap: 16px;
}

.error-alert {
  margin-bottom: 16px;
}

.result-hero {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.hero-tag {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
  font-size: 18px;
}

.hero-symbol {
  font-size: 18px;
  color: v-bind('themeVars.textColor1');
}

.desc {
  margin-bottom: 12px;
}

.special-text {
  margin-bottom: 16px;
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.copy-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.copy-row :deep(input) {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
}
</style>
