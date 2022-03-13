<template>
  <div class="select-source">
    <div
      class="selector"
      v-if="sources.length > 0"
    >
      <div
        class="selector"
        v-for="source in sources"
        :key="source.id"
      >
        <source-comp
          :value="source"
          @update="val => updateSource(source.id, val)"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import indicators from '@/assets/data/indicators'
import { defineComponent, onMounted, ref, watch } from 'vue'

import SourceComp from '@/components/strategies/path-editor/Source.vue'
import { v4 } from 'uuid'

interface SourceTree {
  id: string
  indicatorKey: string | undefined
  data: Record<string, any>
}

interface SourceTreeWithoutId {
  indicatorKey: string | undefined
  data: Record<string, any>
}

export default defineComponent({
  components: { SourceComp },
  emits: ['update'],
  props: {
    source: {
      type: Object as any,
      required: true
    }
  },
  setup (props, context) {
    const indicatorKeys = indicators.map(i => i.key)
    const fields = ref<any[]>([])
    const sources = ref<SourceTree[]>([])

    onMounted(() => {
      let value = props.source
      while (value) {
        if (value.value.indicatorKey) {
          sources.value.push({
            id: v4(),
            indicatorKey: value.value.indicatorKey,
            data: value.value.data
          })
          if (value.value.data.source) {
            value = value.value.data.source
          } else {
            value = undefined
          }
        } else {
          value = undefined
        }
      }

      const lastSource = sources.value[sources.value.length - 1]
      const key = lastSource ? lastSource.indicatorKey || '' : ''
      if (sources.value.length > 0 && prevIsNotFinal(key)) {
        sources.value.push({
          id: v4(),
          indicatorKey: '',
          data: {}
        })
      }
    })

    function traverseTree (tree: SourceTreeWithoutId, index: number): SourceTreeWithoutId {
      let localTree: SourceTreeWithoutId = tree
      for (let i = 0; i < index; i++) {
        localTree = localTree.data.source
      }

      return localTree
    }

    function getIndicatorFieldKeys (indicatorKey: string): string[] {
      const indicator = indicators.find(i => i.key === indicatorKey)
      if (!indicator) return []

      const keys = indicator.fields.map(f => f.key)
      if (indicator.hasSource) keys.push('source')
      return keys
    }

    watch(sources, () => {
      const tree: SourceTreeWithoutId = {
        indicatorKey: '',
        data: {
          source: {
            variable: false,
            value: {}
          }
        }
      }

      for (let i = 0; i < sources.value.length; i++) {
        const source = sources.value[i]

        const tracker = traverseTree(tree, i)
        if (!tracker) continue
        tracker.indicatorKey = source.indicatorKey
        tracker.data = source.data
        if (i !== sources.value.length - 1) {
          tracker.data.source = { variable: false, value: {} }
        }

        const keys = getIndicatorFieldKeys(source.indicatorKey || '')
        for (const key of Object.keys(tracker.data)) {
          if (!keys.includes(key)) {
            delete tracker.data[key]
          }
        }
      }

      context.emit('update', tree)
    }, { deep: true })

    function prevIsNotFinal (val: string): boolean {
      const indicator = indicators.find(i => i.key === val)
      if (!indicator) return false
      return indicator.hasSource
    }

    function updateSource (id: string, data: any) {
      let tempSources: SourceTree[] = JSON.parse(JSON.stringify(sources.value))
      const idx = tempSources.findIndex(s => s.id === id)
      if (idx === -1) return
      data.id = id
      tempSources.splice(idx, 1, data)

      if (!prevIsNotFinal(data.indicatorKey) || !data.indicatorKey) {
        tempSources = tempSources.slice(0, idx + 1)
      }

      const lastSource = tempSources[tempSources.length - 1]
      const key = lastSource ? lastSource.indicatorKey || '' : ''
      if (tempSources.length > 0 && prevIsNotFinal(key)) {
        tempSources.push({
          id: v4(),
          indicatorKey: '',
          data: {}
        })
      }

      sources.value = tempSources
    }

    return {
      indicatorKeys,
      indicators,
      fields,
      sources,
      prevIsNotFinal,
      updateSource
    }
  }
})
</script>

<style lang="scss" scoped>
.selector {
  margin: 0.5rem;
  border: 1px solid var(--border-color);
  background-color: var(--background);
}
</style>
