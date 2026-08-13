import type { Component } from 'vue'
import {
  FileDigit,
  LetterX,
  LetterCaseToggle,
  Calendar,
  ArrowsLeftRight,
  Binary,
  Tools,
} from '@vicons/tabler'

// 图标名（Go 端 registry.Tool.Icon）→ 组件映射。
const toolIcons: Record<string, Component> = {
  FileDigit,
  LetterX,
  LetterCaseToggle,
  Calendar,
  ArrowsLeftRight,
  Binary,
}

// getToolIcon 返回图标名对应的组件，未找到时返回兜底图标。
export function getToolIcon(name: string): Component {
  return toolIcons[name] ?? Tools
}
