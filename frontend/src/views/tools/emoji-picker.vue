<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { NCard, NInput, NVirtualList, useMessage } from 'naive-ui'
import { Search } from '@vicons/tabler'
import { useDebounceFn } from '@vueuse/core'
import EmojiCard, { type EmojiInfo } from '../../components/EmojiCard.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()

const query = ref('')
const allEmojis = ref<EmojiInfo[]>([])
const searchResult = ref<EmojiInfo[] | null>(null)

const grouped = computed<Record<string, EmojiInfo[]>>(() => {
  const m: Record<string, EmojiInfo[]> = {}
  for (const e of allEmojis.value) {
    ;(m[e.group] ??= []).push(e)
  }
  return m
})

async function fetchAll() {
  try {
    const output = await RunTool('emoji-picker', JSON.stringify({ query: '' }))
    allEmojis.value = JSON.parse(output).emojis
  } catch (e) {
    message.error(String(e))
  }
}
onMounted(fetchAll)

async function doSearch() {
  const q = query.value.trim()
  if (!q) {
    searchResult.value = null
    return
  }
  try {
    const output = await RunTool('emoji-picker', JSON.stringify({ query: q }))
    searchResult.value = JSON.parse(output).emojis
  } catch (e) {
    searchResult.value = []
    message.error(String(e))
  }
}
const debouncedSearch = useDebounceFn(doSearch, 500)
watch(query, () => debouncedSearch())

const groupNames: Record<string, string> = {
  'Smileys & Emotion': '表情与情感',
  'People & Body': '人物与身体',
  'Animals & Nature': '动物与自然',
  'Food & Drink': '食物与饮料',
  'Travel & Places': '旅行与地点',
  Activities: '活动',
  Objects: '物品',
  Symbols: '符号',
  Flags: '旗帜',
}

function groupTitle(g: string) {
  return groupNames[g] ?? g
}
</script>

<template>
  <n-card title="Emoji 选择器" class="tool-card">
    <n-input v-model:value="query" placeholder="搜索 Emoji（如 smile、face、笑脸）…" clearable class="search-input">
      <template #prefix>
        <n-icon :component="Search" />
      </template>
    </n-input>

    <template v-if="searchResult !== null">
      <div v-if="searchResult.length === 0" class="no-result">无匹配结果</div>
      <div v-else>
        <div class="section-title">搜索结果（{{ searchResult.length }}）</div>
        <div class="plain-list">
          <EmojiCard v-for="e in searchResult" :key="e.emoji" :info="e" />
        </div>
      </div>
    </template>

    <div v-else v-for="(items, g) in grouped" :key="g" class="emoji-group">
      <div class="section-title">{{ groupTitle(g) }}（{{ items.length }}）</div>
      <n-virtual-list :items="items" :item-size="68" key-field="emoji" class="group-list">
        <template #default="{ item }">
          <EmojiCard :info="item" />
        </template>
      </n-virtual-list>
    </div>
  </n-card>
</template>

<style scoped>
.search-input {
  max-width: 520px;
}

.emoji-group {
  margin-top: 16px;
}

.section-title {
  font-size: 15px;
  font-weight: 700;
  margin-bottom: 8px;
}

.group-list {
  height: 55vh;
}

.plain-list {
  margin-top: 8px;
}

.no-result {
  margin-top: 24px;
  text-align: center;
  opacity: 0.6;
  font-size: 15px;
}
</style>