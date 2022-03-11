<template>
  <div class="select-source">
    <button
      @click="toggleOpen"
      class="select"
    >
      Select source
    </button>
    <div class="selector" v-if="open">
      <v-select
        :options="indicators"
        label="name"
        :reduce="x => x.key"
        v-model="sourceValue.indicatorKey"
        @option:selected="updateData"
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
            v-model="sourceValue.data[field.key]"
          >
          <input
            v-if="field.type === 'text'"
            type="text"
            :placeholder="field.default"
            v-model="sourceValue.data[field.key]"
          >
          <input
            v-if="field.type === 'checkbox'"
            type="checkbox"
            v-model="sourceValue.data[field.key]"
          >
          <select
            v-if="field.type === 'select'"
            v-model="sourceValue.data[field.key]"
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
      <button @click="$emit('select', { indicatorKey, data })">
        Select
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import indicators from '@/assets/data/indicators'
import { defineComponent, onMounted, ref, watch } from 'vue'

export default defineComponent({
  emits: ['select', 'update'],
  props: {
    source: {
      type: Object as any,
      required: true
    }
  },
  setup (props, context) {
    const open = ref<boolean>(false)
    const indicatorKeys = indicators.map(i => i.key)
    const fields = ref<any>([])
    const sourceValue = ref<any>(props.source || {
      indicatorKey: 'CANDLE_POSITION_VALUE',
      data: {
        candlePosition: 'CLOSE'
      }
    })

    watch(sourceValue, () => {
      context.emit('update', sourceValue.value)
      updateData()
    })

    onMounted(() => {
      updateData()
    })

    function toggleOpen () {
      open.value = !open.value
    }

    function updateData (e?: any) {
      let key: string
      if (!e) {
        key = sourceValue.value.indicatorKey
      } else {
        key = e.key
      }

      console.log(key)

      const indicator = indicators.find(i => i.key === key)
      if (!indicator) return

      if (e) {
        sourceValue.value.indicatorKey = e.key
      }

      for (const localKey of Object.keys(sourceValue.value.data)) {
        const field = indicator.fields.find(f => f.key === localKey)
        if (!field) {
          delete sourceValue.value.data[localKey]
        } else if (field.default) {
          sourceValue.value.data[localKey] = field.default
        }
      }

      fields.value = indicator.fields
    }

    return {
      open,
      indicatorKeys,
      indicators,
      selectedIndicatorKey: sourceValue.value.indicatorKey,
      fields,
      data: sourceValue.value.data,
      sourceValue,
      toggleOpen,
      updateData
    }
  }
})
</script>
