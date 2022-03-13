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
        :class="{ active: finalIndicator.timeFrame === tf }"
        @click="() => selectTimeframe(tf)"
      >
        {{ tf }}
      </div>
    </div>
    <div
      class="source input"
      v-if="hasSource"
    >
      <p>Source</p>
      <select-source
        @update="updateSource"
        :source="getIndicatorSource()"
      />
    </div>
    <div
      class="input"
      v-for="key in Object.keys(indicator.data || {}).filter(f => f !== 'source')"
      :key="key"
    >
      <p>{{ key }}</p>
      <input
        v-if="getFieldValues(key, 'type') === 'number'"
        type="number"
        :max="getFieldValues(key, 'max')"
        :min="getFieldValues(key, 'min')"
        :placeholder="getFieldValues(key, 'default')"
        v-model="indicator.data[key].value"
      >
      <input
        v-if="getFieldValues(key, 'type') === 'text'"
        :type="getFieldValues(key, 'type')"
        v-model="indicator.data[key].value"
      >
      <input
        v-if="getFieldValues(key, 'type') === 'checkbox'"
        :type="getFieldValues(key, 'type')"
        v-model="indicator.data[key].value"
      >
      <select
        v-if="getFieldValues(key, 'type') === 'select'"
        v-model="indicator.data[key].value"
      >
        <option
          v-for="option in getFieldValues(key, 'options')"
          :key="option"
          :value="option"
        >
          {{ option }}
        </option>
      </select>
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
      v-if="finalIndicator.timeFrame"
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
import { defineComponent, PropType, ref } from 'vue'

import SelectSource from '@/components/strategies/path-editor/SelectSource.vue'

export default defineComponent({
  components: {
    SelectSource
  },
  emits: ['delete'],
  props: {
    indicator: {
      required: true,
      type: Object as PropType<IndicatorSettings>
    }
  },
  setup (props) {
    const actualIndicator = indicators.find(indicator => indicator.key === props.indicator.indicatorKey)
    const options = ref<any[]>([])

    function selectTimeframe (tf: TimeFrame) {
      // eslint-disable-next-line vue/no-mutating-props
      props.indicator.timeFrame = tf
    }

    function getFieldValues (key: string, item: string): any {
      if (!actualIndicator) return ''

      const field: any = actualIndicator.fields.find(f => f.key === key)
      return field[item]
    }

    function hasTimeframe (): boolean {
      return (actualIndicator || { hasTimeframe: false }).hasTimeframe
    }

    function hasSource (): boolean {
      return (actualIndicator || { hasSource: false }).hasSource
    }

    function getIndicatorSource (): any {
      if (!props.indicator.data.source && hasSource()) {
        props.indicator.data.source = {
          variable: false,
          value: {
            indicatorKey: 'CANDLE_POSITION_VALUE',
            data: {
              candlePosition: {
                variable: false,
                value: 'CLOSE'
              }
            }
          }
        }
      }
      return props.indicator.data.source
    }

    function updateSource (source: any) {
      props.indicator.data.source = {
        value: source,
        variable: false
      }
    }

    return {
      finalIndicator: props.indicator,
      indicatorName: (actualIndicator || { name: 'Name not found' }).name,
      timeframes,
      hasTimeframe: hasTimeframe(),
      hasSource: hasSource(),
      options,
      actualIndicator,
      selectTimeframe,
      getFieldValues,
      getIndicatorSource,
      updateSource
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

.source {
  padding: 0.6rem;
  background-color: var(--background-darken);
  border: 1px solid var(--border-color);
}
</style>
