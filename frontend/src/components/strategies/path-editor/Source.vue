<template>
  <div class="source" v-if="ready">
    <v-select
      :options="indicators.filter(i => i.hasTimeframe)"
      label="name"
      :reduce="x => x.key"
      :disabled="last"
      v-model="selectedIndicatorKey"
    />
    <div class="fields">
      <div
        class="input"
        v-for="field in fields"
        :key="field.key"
      >
        <p>{{ field.name }}</p>
        <input
          v-if="field.type === 'number'"
          type="number"
          :max="field.max"
          :min="field.min"
          :placeholder="field.default"
          v-model="sourceValue.data[field.key].value"
        >
        <input
          v-if="field.type === 'text'"
          type="text"
          :placeholder="field.default"
          v-model="sourceValue.data[field.key].value"
        >
        <input
          v-if="field.type === 'checkbox'"
          type="checkbox"
          v-model="sourceValue.data[field.key].value"
        >
        <select
          v-if="field.type === 'select'"
          v-model="sourceValue.data[field.key].value"
        >
          <option
            v-for="option in field.options"
            :key="option"
            :value="option"
          >
            {{ option }}
          </option>
        </select>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import indicators from '@/assets/data/indicators'
import { Data } from '@/types/Path'
import { computed, defineComponent, onMounted, PropType, ref, watch } from 'vue'

interface SourceTree {
  indicatorKey: string | undefined
  data: Record<string, Data>
}

export default defineComponent({
  emits: ['update'],
  props: {
    value: {
      type: Object as PropType<SourceTree>
    },
    last: {
      type: Boolean
    }
  },
  setup (props, context) {
    const selectedIndicatorKey = ref<string>()
    const ready = ref<boolean>(false)
    const fields = computed(() => {
      const indicator = indicators.find(i => i.key === selectedIndicatorKey.value)
      if (!indicator) return []
      return indicator.fields
    })

    const sourceValue = ref<any>({
      indicatorKey: undefined,
      data: {}
    })

    onMounted(() => {
      if (props.value) {
        const data = JSON.parse(JSON.stringify(props.value))

        sourceValue.value = {
          indicatorKey: data.indicatorKey,
          data: data.indicatorKey ? data.data : {}
        }
        selectedIndicatorKey.value = data.indicatorKey
      }

      ready.value = true
    })

    // TODO: this causes a lot of recursive updates,
    // which is not intended. Needs work.
    watch(sourceValue, () => {
      context.emit('update', sourceValue.value)
    }, { deep: true })

    watch(selectedIndicatorKey, () => {
      sourceValue.value.indicatorKey = selectedIndicatorKey.value
      if (!selectedIndicatorKey.value) {
        sourceValue.value.data = {}
      } else {
        const indicatorData: Record<string, Data> = {}
        const indicator = indicators.find(i => i.key === selectedIndicatorKey.value)

        // Set the fields that the indicator has on it
        // so that they can be modified/configured by the inputs
        // in the template
        for (const field of (indicator || { fields: [] }).fields) {
          indicatorData[field.key] = {
            variable: false,
            value: field.default
          }
        }

        sourceValue.value.data = indicatorData
      }
    })

    return {
      indicators,
      selectedIndicatorKey,
      fields,
      sourceValue,
      ready
    }
  }
})
</script>

<style lang="scss" scoped>
.source {
  padding: 0.5rem;
  padding-bottom: 0;
  .fields {
    margin-top: 0.5rem;
    display: flex;
    flex-wrap: wrap;
    .input {
      margin-right: 0.5rem;
      margin-bottom: 0.5rem;
      p {
        margin-bottom: 0.3rem;
        font-size: 13px;
        color: var(--text-secondary);
      }
      input, select {
        padding: 0.3rem;
        border: 1px solid var(--border-color);
      }
    }
  }
}
</style>
