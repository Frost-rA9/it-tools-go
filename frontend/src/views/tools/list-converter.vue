<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { NCard, NInput, NSwitch, NSelect, NButton, useMessage, useThemeVars } from 'naive-ui'
import { useDebounceFn, useClipboard } from '@vueuse/core'
import ToolTextarea from '../../components/ToolTextarea.vue'
import { RunTool } from '../../../wailsjs/go/app/App'

const message = useMessage()
const themeVars = useThemeVars()

const config = reactive({
  lower_case: false,
  trim_items: true,
  remove_duplicates: true,
  keep_line_breaks: false,
  item_prefix: '',
  item_suffix: '',
  list_prefix: '',
  list_suffix: '',
  reverse_list: false,
  sort_list: '' as string,
  separator: ', ',
})

const sortOptions = [
  { label: '升序', value: 'asc' },
  { label: '降序', value: 'desc' },
]

const inputText = ref('')
const outputText = ref('')

async function run() {
  try {
    const output = await RunTool(
      'list-converter',
      JSON.stringify({ text: inputText.value, ...config }),
    )
    outputText.value = JSON.parse(output).result
  } catch (e) {
    outputText.value = ''
    message.error(String(e))
  }
}

const debouncedRun = useDebounceFn(run, 150)
watch([inputText, config], () => debouncedRun(), { immediate: true })

const { copy } = useClipboard({ source: outputText })

async function copyResult() {
  await copy()
  message.success('已复制到剪贴板')
}
</script>

<template>
  
    <n-card title="转换选项" class="tool-card">
      <div class="options-grid">
        <div class="option-row">
          <span class="option-label">去除重复项</span>
          <n-switch v-model:value="config.remove_duplicates" />
        </div>
        <div class="option-row">
          <span class="option-label">转换为小写</span>
          <n-switch v-model:value="config.lower_case" />
        </div>
        <div class="option-row">
          <span class="option-label">去除首尾空格</span>
          <n-switch v-model:value="config.trim_items" />
        </div>
        <div class="option-row">
          <span class="option-label">保留换行</span>
          <n-switch v-model:value="config.keep_line_breaks" />
        </div>
        <div class="option-row">
          <span class="option-label">排序</span>
          <n-select v-model:value="config.sort_list" :options="sortOptions" clearable placeholder="不排序" style="width: 140px" />
        </div>
        <div class="option-row">
          <span class="option-label">反转列表</span>
          <n-switch v-model:value="config.reverse_list" />
        </div>
      </div>

      <div class="field">
        <div class="field-label">分隔符</div>
        <n-input v-model:value="config.separator" placeholder=", " />
      </div>

      <div class="field">
        <div class="field-label">条目前缀 / 后缀</div>
        <div class="inline-row">
          <n-input v-model:value="config.item_prefix" placeholder="条目前缀" />
          <n-input v-model:value="config.item_suffix" placeholder="条目后缀" />
        </div>
      </div>

      <div class="field">
        <div class="field-label">列表前缀 / 后缀</div>
        <div class="inline-row">
          <n-input v-model:value="config.list_prefix" placeholder="列表前缀" />
          <n-input v-model:value="config.list_suffix" placeholder="列表后缀" />
        </div>
      </div>
    </n-card>

    <n-card title="列表转换器" class="tool-card">
      <ToolTextarea v-model:value="inputText" label="输入数据" :rows="6" placeholder="每行一个条目，在此粘贴数据…" />

      <ToolTextarea v-model:value="outputText" label="转换结果" :rows="6" readonly placeholder="转换结果将显示在这里" />

      <div class="copy-row">
        <n-button type="primary" @click="copyResult">复制结果</n-button>
      </div>
    </n-card>
</template>

<style scoped>

.options-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}

.option-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.option-label {
  font-size: 14px;
  color: v-bind('themeVars.textColor2');
}

.field {
  margin-bottom: 16px;
}

.field-label {
  font-size: 14px;
  margin-bottom: 6px;
  color: v-bind('themeVars.textColor2');
}

.inline-row {
  display: flex;
  gap: 8px;
}

.copy-row {
  display: flex;
  justify-content: center;
}
</style>
