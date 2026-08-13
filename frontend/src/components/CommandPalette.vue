<script setup lang="ts">
import { ref, computed, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { NModal, NInput, NEmpty } from 'naive-ui'
import { useMagicKeys } from '@vueuse/core'
import { useToolsStore } from '../stores/tools'

const router = useRouter()
const store = useToolsStore()

const isOpen = ref(false)
const prompt = ref('')
const inputRef = ref<InstanceType<typeof NInput> | null>(null)

const keys = useMagicKeys()
const cmdK = keys['Ctrl+K']

watch(cmdK, (v) => {
  if (v) open()
})

const filteredTools = computed(() => {
  const q = prompt.value.trim().toLowerCase()
  if (!q) return store.tools
  return store.tools.filter(
    (t) =>
      t.name.toLowerCase().includes(q) ||
      t.description.toLowerCase().includes(q) ||
      t.keywords.some((k) => k.toLowerCase().includes(q)),
  )
})

const grouped = computed(() => {
  const map = new Map<string, typeof store.tools>()
  for (const t of filteredTools.value) {
    const list = map.get(t.category) ?? []
    list.push(t)
    map.set(t.category, list)
  }
  return Array.from(map.entries()).map(([category, tools]) => ({ category, tools }))
})

function open() {
  isOpen.value = true
  nextTick(() => inputRef.value?.focus())
}

function close() {
  isOpen.value = false
  prompt.value = ''
}

function select(id: string) {
  router.push({ name: 'tool', params: { id } })
  close()
}

defineExpose({ open })
</script>

<template>
  <n-modal v-model:show="isOpen" :show-icon="false" transform-origin="center">
    <div class="palette">
      <n-input
        ref="inputRef"
        v-model:value="prompt"
        size="large"
        placeholder="输入以搜索工具…"
        clearable
        @keydown.esc="close"
      />

      <div class="palette-results">
        <template v-if="grouped.length">
          <div v-for="group in grouped" :key="group.category" class="palette-group">
            <div class="palette-category">{{ group.category }}</div>
            <div
              v-for="tool in group.tools"
              :key="tool.id"
              class="palette-item"
              @click="select(tool.id)"
            >
              {{ tool.name }}
            </div>
          </div>
        </template>
        <n-empty v-else description="无匹配工具" class="palette-empty" />
      </div>
    </div>
  </n-modal>
</template>

<style scoped>
.palette {
  max-width: 600px;
}

.palette-results {
  margin-top: 12px;
  max-height: 400px;
  overflow-y: auto;
}

.palette-group {
  margin-bottom: 8px;
}

.palette-category {
  font-size: 13px;
  font-weight: 600;
  opacity: 0.6;
  margin: 8px 4px 4px;
}

.palette-item {
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.15s;
}

.palette-item:hover {
  background-color: rgba(128, 128, 128, 0.15);
}

.palette-empty {
  margin: 16px 0;
}
</style>
