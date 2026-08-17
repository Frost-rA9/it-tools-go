<script setup lang="ts">
import { computed, type ComputedRef } from 'vue'
import { NCard, NGrid, NGi, useThemeVars } from 'naive-ui'
import { useWindowSize } from '@vueuse/core'

const themeVars = useThemeVars()
const { width, height } = useWindowSize()

interface InfoItem {
  label: string
  value: ComputedRef<string>
}

interface Section {
  name: string
  information: InfoItem[]
}

const screen: InfoItem[] = [
  { label: 'Screen size', value: computed(() => `${window.screen.availWidth} x ${window.screen.availHeight}`) },
  { label: 'Orientation', value: computed(() => window.screen.orientation?.type ?? 'unknown') },
  { label: 'Orientation angle', value: computed(() => `${window.screen.orientation?.angle ?? 0}°`) },
  { label: 'Color depth', value: computed(() => `${window.screen.colorDepth} bits`) },
  { label: 'Pixel ratio', value: computed(() => `${window.devicePixelRatio} dppx`) },
  { label: 'Window size', value: computed(() => `${width.value} x ${height.value}`) },
]

const device: InfoItem[] = [
  { label: 'Browser vendor', value: computed(() => navigator.vendor || 'unknown') },
  { label: 'Languages', value: computed(() => navigator.languages?.join(', ') || navigator.language || 'unknown') },
  { label: 'Platform', value: computed(() => navigator.platform || 'unknown') },
  { label: 'User agent', value: computed(() => navigator.userAgent) },
]

const sections: Section[] = [
  { name: 'Screen', information: screen },
  { name: 'Device', information: device },
]
</script>

<template>
  <n-card v-for="section in sections" :key="section.name" :title="section.name" class="card">
    <n-grid cols="1 400:2" x-gap="12" y-gap="12">
      <n-gi v-for="item in section.information" :key="item.label" class="information">
        <div class="label">{{ item.label }}</div>
        <div class="value" v-if="item.value.value">{{ item.value.value }}</div>
        <div class="undefined-value" v-else>unknown</div>
      </n-gi>
    </n-grid>
  </n-card>
</template>

<style scoped>
.card {
  min-width: 400px;
  width: 100%;
}

.information {
  padding: 14px 16px;
  border-radius: 4px;
  background-color: color-mix(in srgb, v-bind('themeVars.textColor1') 6%, transparent);
}

.label {
  font-size: 14px;
  opacity: 0.8;
  line-height: 1;
  margin-bottom: 5px;
  color: v-bind('themeVars.textColor2');
}

.value {
  font-size: 20px;
  font-weight: 400;
  word-break: break-all;
  color: v-bind('themeVars.textColor1');
}

.undefined-value {
  opacity: 0.8;
  color: v-bind('themeVars.textColor3');
}
</style>