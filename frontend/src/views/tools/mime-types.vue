<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { NCard, NSelect, NTable, NTag, NText, useMessage, useThemeVars } from 'naive-ui'
import { RunTool } from '../../../wailsjs/go/app/App'

interface MimePair {
  mime_type: string
  extensions: string[]
}

const message = useMessage()
const themeVars = useThemeVars()

const all = ref<MimePair[]>([])
const selectedMimeType = ref<string | null>(null)
const selectedExtension = ref<string | null>(null)

const mimeOptions = computed(() => all.value.map((p) => ({ label: p.mime_type, value: p.mime_type })))
const extOptions = computed(() => {
  const seen = new Set<string>()
  const opts: { label: string; value: string }[] = []
  for (const p of all.value) {
    for (const ext of p.extensions) {
      if (!seen.has(ext)) {
        seen.add(ext)
        opts.push({ label: `.${ext}`, value: ext })
      }
    }
  }
  return opts
})

const extensionsFound = computed(() => {
  const p = all.value.find((x) => x.mime_type === selectedMimeType.value)
  return p?.extensions ?? []
})

const mimeTypeFound = computed(() => {
  const p = all.value.find((x) => x.extensions.includes(selectedExtension.value ?? ''))
  return p?.mime_type ?? ''
})

onMounted(async () => {
  try {
    const output = await RunTool('mime-types', JSON.stringify({ list: true }))
    all.value = JSON.parse(output).all
  } catch (e) {
    message.error(String(e))
  }
})
</script>

<template>
  <n-card class="card">
    <div class="h2">MIME 类型转扩展名</div>
    <div class="sub">查看一个 MIME 类型关联哪些文件扩展名</div>
    <n-select v-model:value="selectedMimeType" :options="mimeOptions" filterable clearable placeholder="选择 MIME 类型… (如 application/pdf)" class="select" />

    <div v-if="extensionsFound.length" class="result">
      <span>扩展名 <n-tag round :bordered="false">{{ selectedMimeType }}</n-tag> 的 MIME 类型对应扩展名：</span>
      <div class="tags">
        <n-tag v-for="ext in extensionsFound" :key="ext" round :bordered="false" type="primary" class="tag">
          .{{ ext }}
        </n-tag>
      </div>
    </div>
  </n-card>

  <n-card class="card">
    <div class="h2">扩展名转 MIME 类型</div>
    <div class="sub">查看一个文件扩展名关联的 MIME 类型</div>
    <n-select v-model:value="selectedExtension" :options="extOptions" filterable clearable placeholder="选择扩展名… (如 pdf)" class="select" />

    <div v-if="selectedExtension" class="result">
      <span>扩展名 <n-tag round :bordered="false">{{ selectedExtension }}</n-tag> 关联的 MIME 类型：</span>
      <div class="tags">
        <n-tag round :bordered="false" type="primary">{{ mimeTypeFound }}</n-tag>
      </div>
    </div>
  </n-card>

  <n-card class="card">
    <div class="table-wrap">
      <n-table :bordered="false">
        <thead>
          <tr>
            <th>MIME 类型</th>
            <th>扩展名</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in all" :key="p.mime_type">
            <td class="mime-cell">{{ p.mime_type }}</td>
            <td>
              <n-tag v-for="ext in p.extensions" :key="ext" round :bordered="false" class="tag">
                .{{ ext }}
              </n-tag>
            </td>
          </tr>
        </tbody>
      </n-table>
    </div>
  </n-card>
</template>

<style scoped>
.card {
  min-width: 400px;
}

.h2 {
  font-size: 19px;
  font-weight: 600;
  margin-bottom: 2px;
  color: v-bind('themeVars.textColor1');
}

.sub {
  opacity: 0.8;
  font-size: 13px;
  margin-bottom: 16px;
  color: v-bind('themeVars.textColor2');
}

.select {
  margin: 16px 0;
}

.result {
  margin-top: 8px;
  font-size: 14px;
  color: v-bind('themeVars.textColor1');
}

.tags {
  margin-top: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag {
  margin-right: 4px;
}

.table-wrap {
  width: 100%;
  overflow-x: auto;
}

.mime-cell {
  font-family: 'Cascadia Code', monospace;
  font-size: 13px;
  word-break: break-all;
}
</style>