<template>
  <div class="indicator-editor">
    <div class="indicator-name">
      <p>{{ indicatorName }}</p>
      <fa-icon
        icon="times"
        @click="$emit('delete')"
      />
    </div>
    <div
      v-if="hasTimeframe"
      class="timeframe"
    >
      <div
        class="tf"
        v-for="tf in timeframes"
        :key="tf"
        :class="{ active: finalIndicator.timeframe === tf }"
        @click="() => selectTimeframe(tf)"
      >
        {{ tf }}
      </div>
    </div>
    <div
      class="source input"
      v-if="source"
    >
      <p>Source</p>
      <v-select
        :options="options"
        label="name"
        :reduce="x => x.key"
        v-model="finalIndicator.sourceKey"
      />
    </div>
    <div
      class="input"
      v-for="field in finalIndicator.fields.filter(f => f.name !== 'Source')"
      :key="field.name"
    >
      <p>{{ field.name }}</p>
      <input
        v-if="field.type === 'number'"
        :type="field.type"
        :max="field.max"
        :min="field.min"
        :placeholder="field.default"
        v-model="field.value"
      >
      <input
        v-if="field.type === 'text'"
        :type="field.type"
        v-model="field.value"
      >
      <input
        v-if="field.type === 'checkbox'"
        :type="field.type"
        v-model="field.value"
      >
    </div>
    <div
      class="input"
    >
      <p>Offset (%)</p>
      <input type="number"
        v-model="finalIndicator.offset"
        min="0"
        max="100"
      >
    </div>
    <div
      class="input"
      v-if="finalIndicator.timeframe"
    >
      <p>Candles back</p>
      <input type="number"
        v-model="finalIndicator.candlesBack"
        min="0"
      >
    </div>
  </div>
</template>

<script lang="ts">
import indicators from '@/assets/data/indicators'
import { IndicatorSettings, TimeFrame, timeframes } from '@/types/Path'
import { defineComponent, PropType } from 'vue'

export default defineComponent({
  emits: ['delete'],
  props: {
    indicator: {
      required: true,
      type: Object as PropType<IndicatorSettings>
    }
  },
  setup (props) {
    const actualIndicator = indicators.find(indicator => indicator.key === props.indicator.indicatorKey)
    const source = props.indicator.fields.find(f => f.name === 'Source')
    const options: { key: string, name: string }[] = []

    if (source) {
      for (const src of source.options) {
        options.push({ key: src, name: indicators.find(indicator => indicator.key === src)?.name || src })
      }
      if (source.options.length === 0) {
        for (const indicator of indicators) {
          if (indicator.hasTimeframe) {
            options.push({ key: indicator.key, name: indicator.name })
          }
        }
      }
      if (source.default) {
        // eslint-disable-next-line vue/no-mutating-props
        props.indicator.sourceKey = source.default
      }
    }

    for (const field of props.indicator.fields) {
      if (field.default) {
        field.value = field.default
      }
    }

    function selectTimeframe (tf: TimeFrame) {
      // eslint-disable-next-line vue/no-mutating-props
      props.indicator.timeframe = tf
    }

    return {
      finalIndicator: props.indicator,
      indicatorName: (actualIndicator || { name: 'Name not found' }).name,
      timeframes,
      hasTimeframe: (actualIndicator || { hasTimeframe: false }).hasTimeframe,
      source,
      options,
      selectTimeframe
    }
  }
})
</script>

<style lang="scss" scoped>
.indicator-editor {
  .indicator-name {
    user-select: none;
    padding: 0.7rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border: 1px solid var(--border-color);
    border-bottom: 2px solid var(--primary);
    p {
      font-size: 14px;
      font-weight: 600;
    }
    svg {
      cursor: pointer;
    }
  }
  .timeframe {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-around;
    padding: 0.7rem 0;
    padding-bottom: 0.2rem;
    border-bottom: 2px solid var(--primary);
    margin-bottom: 1rem;
    .tf {
      user-select: none;
      font-size: 14px;
      cursor: pointer;
      width: 30px;
      margin-bottom: 0.5rem;
      &.active, &:hover {
        font-weight: bold;
        text-decoration: underline;
      }
    }
  }
  .input {
    margin-top: 1rem;
    p {
      margin-bottom: 0.5rem;
      font-size: 14px;
    }
    input {
      padding: 0.5rem;
      outline: none;
      border: 1px solid var(--border-color);
      width: 100%;
    }
  }
}
</style>
