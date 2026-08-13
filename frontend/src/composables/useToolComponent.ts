import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useToolsStore } from '../stores/tools'

const modules = import.meta.glob('../views/tools/*.vue')

export function useToolComponent() {
  const route = useRoute()
  const store = useToolsStore()

  const tool = computed(() => store.getById(route.params.id as string))

  const component = computed(() => {
    if (!tool.value) return null
    const key = `../views/tools/${toKebabCase(tool.value.id)}.vue`
    return modules[key]
  })

  return { tool, component }
}

function toKebabCase(s: string) {
  return s.replace(/([a-z0-9])([A-Z])/g, '$1-$2').toLowerCase()
}
