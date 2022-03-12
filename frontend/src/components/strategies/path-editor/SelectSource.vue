<template>
  <div class="select-source">
    <div
      class="selector"
      v-if="ready"
    >
      <source-comp
        :value="{ indicatorKey: tree.indicatorKey, data: tree.data }"
        @update="val => updateTree(1, val)"
      />
      <div class="selector" v-if="prevIsNotFinal(tree.indicatorKey)">
        <source-comp
          :value="tree.data.source ? { indicatorKey: tree.data.source.indicatorKey, data: tree.data.source.data } : undefined"
          @update="val => updateTree(2, val)"
        />
        <div class="selector" v-if="tree.data.source && prevIsNotFinal(tree.data.source.indicatorKey)">
          <source-comp
            :value="tree.data.source.data.source ? { indicatorKey: tree.data.source.data.source.indicatorKey, data: tree.data.source.data.source.data } : undefined"
            @update="val => updateTree(3, val)"
          />
          <div class="selector" v-if="tree.data.source.data.source && prevIsNotFinal(tree.data.source.data.source.indicatorKey)">
            <source-comp
              :last="true"
              :value="{ indicatorKey: 'CANDLE_POSITION_VALUE', data: { candlePosition: tree.data.source.data.source ? tree.data.source.data.source.indicatorKey : '' } }"
              @update="val => updateTree(4, val)"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import indicators from '@/assets/data/indicators'
import { defineComponent, onMounted, ref, watch } from 'vue'

import SourceComp from '@/components/strategies/path-editor/Source.vue'

interface SourceTree {
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
    const ready = ref<boolean>()
    const tree = ref<SourceTree>({
      indicatorKey: 'CANDLE_POSITION_VALUE',
      data: {
        candlePosition: 'CLOSE'
      }
    })

    onMounted(() => {
      tree.value.indicatorKey = props.source.indicatorKey
      tree.value.data = props.source.data
      ready.value = true
    })

    watch(tree, () => {
      context.emit('update', tree.value)
    }, { deep: true })

    function updateTree (level: number, value: SourceTree) {
      if (level === 1) {
        if (!value.indicatorKey || value.indicatorKey === null) {
          tree.value.indicatorKey = ''
          tree.value.data = {}
          return
        }
        tree.value.indicatorKey = value.indicatorKey
        tree.value.data = value.data
      }
      if (level === 2) {
        if (!tree.value.data.source || !tree.value.data.source.indicatorKey) {
          tree.value.data.source = {
            indicatorKey: undefined,
            data: {}
          }
        }
        tree.value.data.source.indicatorKey = value.indicatorKey
        tree.value.data.source.data = value.data
      }
      if (level === 3) {
        if (!tree.value.data.source.data.source || !tree.value.data.source.data.source.indicatorKey) {
          tree.value.data.source.data.source = {
            indicatorKey: undefined,
            data: {}
          }
        }
        tree.value.data.source.data.source.indicatorKey = value.indicatorKey
        tree.value.data.source.data.source.data = value.data
      }
    }

    function prevIsNotFinal (val: string): boolean {
      return val !== 'CANDLE_POSITION_VALUE' && val !== ''
    }

    return {
      indicatorKeys,
      indicators,
      fields,
      tree,
      ready,
      updateTree,
      prevIsNotFinal
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
