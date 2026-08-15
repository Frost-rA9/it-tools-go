import type { Component } from 'vue'
import {
  FileDigit,
  LetterX,
  LetterCaseToggle,
  Calendar,
  ArrowsLeftRight,
  Binary,
  TextWrap,
  AlignJustified,
  List,
  Markdown,
  Braces,
  Tools,
  ArrowsShuffle,
  EyeOff,
  Lock,
} from '@vicons/tabler'

// 图标名（Go 端 registry.Tool.Icon）→ 组件映射。
const toolIcons: Record<string, Component> = {
  FileDigit,
  LetterX,
  LetterCaseToggle,
  Calendar,
  ArrowsLeftRight,
  Binary,
  TextWrap,
  AlignJustified,
  List,
  Markdown,
  Braces,
  ArrowsShuffle,
  EyeOff,
  Lock,
}

// getToolIcon 返回图标名对应的组件，未找到时返回兜底图标。
export function getToolIcon(name: string): Component {
  return toolIcons[name] ?? Tools
}
