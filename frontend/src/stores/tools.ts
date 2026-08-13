import { defineStore } from 'pinia'
import { ListTools } from '../../wailsjs/go/app/App'
import type { registry } from '../../wailsjs/go/models'

export type Tool = registry.Tool

interface ToolCategory {
  name: string
  tools: Tool[]
}

export const useToolsStore = defineStore('tools', {
  state: () => ({
    tools: [] as Tool[],
    loaded: false,
  }),
  getters: {
    categories(state): ToolCategory[] {
      const map = new Map<string, Tool[]>()
      for (const tool of state.tools) {
        const list = map.get(tool.category) ?? []
        list.push(tool)
        map.set(tool.category, list)
      }
      return Array.from(map.entries()).map(([name, tools]) => ({ name, tools }))
    },
    getById(state) {
      return (id: string) => state.tools.find((t) => t.id === id)
    },
  },
  actions: {
    async load() {
      if (this.loaded) return
      this.tools = await ListTools()
      this.loaded = true
    },
  },
})
