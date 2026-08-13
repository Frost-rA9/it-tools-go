<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NLayout, NLayoutSider, useThemeVars } from 'naive-ui'
import { useToolsStore } from '../stores/tools'
import { useUiStore } from '../stores/ui'
import ToolMenu from '../components/ToolMenu.vue'
import Toolbar from '../components/Toolbar.vue'
import heroGradient from '../assets/hero-gradient.svg'

const router = useRouter()
const store = useToolsStore()
const ui = useUiStore()
const themeVars = useThemeVars()
onMounted(() => store.load())
</script>

<template>
  <n-layout has-sider class="app-layout">
    <n-layout-sider
      bordered
      collapse-mode="width"
      :width="240"
      :collapsed-width="0"
      :collapsed="ui.isMenuCollapsed"
      :show-trigger="false"
      :native-scrollbar="false"
      class="app-sider"
    >
      <div class="hero-wrapper" @click="router.push('/')">
        <img :src="heroGradient" class="hero-gradient" alt="" />
        <div class="hero-text">
          <div class="hero-title">IT - TOOLS</div>
          <div class="hero-divider" />
          <div class="hero-subtitle">跨平台开发者工具集</div>
        </div>
      </div>
      <div class="sider-content">
        <ToolMenu />
      </div>
    </n-layout-sider>

    <n-layout class="app-content">
      <Toolbar />
      <n-scrollbar class="content-scroll">
        <router-view />
      </n-scrollbar>
    </n-layout>
  </n-layout>
</template>

<style scoped>
.app-layout {
  height: 100vh;
}

.hero-wrapper {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 10;
  overflow: hidden;
  cursor: pointer;
}

.hero-gradient {
  display: block;
  width: 100%;
  height: 220px;
  margin-top: -65px;
}

.hero-text {
  position: absolute;
  left: 0;
  width: 100%;
  text-align: center;
  top: 16px;
  color: #fff;
}

.hero-title {
  font-size: 25px;
  font-weight: 600;
  letter-spacing: 1px;
}

.hero-divider {
  width: 50px;
  height: 2px;
  border-radius: 4px;
  background-color: v-bind('themeVars.primaryColor');
  margin: 8px auto;
}

.hero-subtitle {
  font-size: 14px;
  opacity: 0.85;
}

.sider-content {
  padding-top: 160px;
  padding-bottom: 200px;
}

.content-scroll {
  height: calc(100vh - 48px);
}
</style>
