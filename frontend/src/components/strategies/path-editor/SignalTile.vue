<template>
  <div class="signal-tile">
    <div class="tile">
      <p>{{ tile.name }}</p>
      <div class="controls">
        <fa-icon icon="times"
          @click="deleteTile"
        />
        <fa-icon icon="pencil-alt"
          @click="toggleEdit"
        />
      </div>
    </div>
    <div
      class="editor"
      v-if="editing"
      @click.self="toggleEdit"
    >
      <div class="inner">
        <div class="top">
          <div class="input">
            <p>NAME</p>
            <input type="text" v-model="finalTile.name">
          </div>
          <div class="input">
            <p>PERSISTENCE</p>
            <input type="number" v-model="finalTile.persistence">
          </div>
        </div>
        <div class="content">
          <div class="indicator">
            <search-indicator
              v-if="!finalTile.indicatorA"
              @select="key => setIndicator('A', key)"
            />
            <edit-indicator
              v-if="finalTile.indicatorA"
              :indicator="tile.indicatorA"
              :variables="variables"
              @delete="() => tile.indicatorA = undefined"
            />
          </div>
          <div class="middle">
            <div
              class="selected operand"
            >
              <fa-icon :icon="(operands.find(o => o.key === finalTile.operand) || { icon: '' }).icon"/>
            </div>
            <div
              v-for="operand in operands.filter(o => o.key !== finalTile.operand)"
              :key="operand.key"
              class="operand"
              :class="{ selected: finalTile.operand === operand.key }"
              @click="setOperand(operand.key.toString())"
            >
              <fa-icon :icon="operand.icon"/>
            </div>
          </div>
          <div class="indicator">
            <search-indicator
              v-if="!finalTile.indicatorB"
              @select="key => setIndicator('B', key)"
            />
            <edit-indicator
              v-if="tile.indicatorB"
              :indicator="tile.indicatorB"
              :variables="variables"
              @delete="() => tile.indicatorB = undefined"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { IndicatorSettings, Operand, SignalTile, TimeFrame } from '@/types/Path'
import { defineComponent, PropType, ref } from 'vue'

import SearchIndicator from '@/components/strategies/path-editor/SearchIndicator.vue'
import EditIndicator from '@/components/strategies/path-editor/EditIndicator.vue'
import indicators from '@/assets/data/indicators'
import { Indicator } from '@/types/Indicator'
import { Variable } from '@/types/Strategy'

export default defineComponent({
  components: { SearchIndicator, EditIndicator },
  props: {
    tile: {
      required: true,
      type: Object as PropType<SignalTile>
    },
    variables: {
      required: true,
      type: Array as PropType<Variable[]>
    }
  },
  emits: ['delete'],
  setup (props, { emit }) {
    const editing = ref(false)
    const operands: { key: Operand, icon: string }[] = [
      { key: Operand.GREATER_THAN, icon: 'greater-than' },
      { key: Operand.GREATER_THAN_OR_EQUAL, icon: 'greater-than-equal' },
      { key: Operand.LOWER_THAN, icon: 'less-than' },
      { key: Operand.LOWER_THAN_OR_EQUAL, icon: 'less-than-equal' },
      { key: Operand.EQUAL, icon: 'equals' },
      { key: Operand.NOT_EQUAL, icon: 'not-equal' },
      { key: Operand.CROSS_ABOVE, icon: 'arrow-up' },
      { key: Operand.CROSS_BELOW, icon: 'arrow-down' }
    ]

    function toggleEdit () {
      editing.value = !editing.value
    }

    function deleteTile () {
      emit('delete')
    }

    function setIndicator (indicator: string, key: string) {
      const foundIndicator: Indicator = JSON.parse(JSON.stringify(indicators.find(i => i.key === key)))
      const data: Record<string, any> = {}

      for (const field of foundIndicator.fields) {
        data[field.key] = {
          variable: false,
          value: field.default
        }
      }

      const settings: IndicatorSettings = {
        timeFrame: foundIndicator.hasTimeframe ? TimeFrame.h1 : undefined,
        candlesBack: 0,
        realTime: false,
        offset: 0,
        indicatorKey: key,
        data
      }
      if (indicator === 'A') {
        props.tile.indicatorA = settings
      } else {
        props.tile.indicatorB = settings
      }
    }

    function setOperand (operand: string) {
      switch (operand) {
        case Operand.GREATER_THAN:
          props.tile.operand = Operand.GREATER_THAN
          break
        case Operand.GREATER_THAN_OR_EQUAL:
          props.tile.operand = Operand.GREATER_THAN_OR_EQUAL
          break
        case Operand.LOWER_THAN:
          props.tile.operand = Operand.LOWER_THAN
          break
        case Operand.LOWER_THAN_OR_EQUAL:
          props.tile.operand = Operand.LOWER_THAN_OR_EQUAL
          break
        case Operand.EQUAL:
          props.tile.operand = Operand.EQUAL
          break
        case Operand.NOT_EQUAL:
          props.tile.operand = Operand.NOT_EQUAL
          break
        case Operand.CROSS_ABOVE:
          props.tile.operand = Operand.CROSS_ABOVE
          break
        case Operand.CROSS_BELOW:
          props.tile.operand = Operand.CROSS_BELOW
          break
        default:
          props.tile.operand = Operand.GREATER_THAN
      }
    }

    return {
      editing,
      finalTile: props.tile,
      operands,
      toggleEdit,
      deleteTile,
      setIndicator,
      setOperand
    }
  }
})
</script>

<style lang="scss" scoped>
.signal-tile {
  .tile {
    padding: 0.7rem;
    background-color: var(--signal-tile);
    border: 1px solid var(--border-color);
    display: flex;
    justify-content: space-between;
    align-items: center;
    svg {
      cursor: pointer;
      color: var(--text-tertiary);
      &:hover {
        color: var(--text-secondary);
      }
    }
    .controls {
      display: flex;
      align-items: center;
      svg {
        margin-left: 1rem;
      }
    }
  }
  .editor {
    position: fixed;
    bottom: 0;
    top: 0;
    left: 0;
    right: 0;
    background-color: rgba(0, 0, 0, 0.3);
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 1rem;
    z-index: 10;
    .inner {
      max-width: var(--container-width);
      background-color: var(--background);
      margin: 0 auto;
      margin-top: 1rem;
      width: 100%;
      max-width: 1200px;
      min-height: 50vh;
      height: 100%;
      max-height: 70vh;
      display: flex;
      flex-direction: column;
      overflow-y: auto;
      .top {
        background-color: var(--primary);
        display: grid;
        padding: 0.7rem 1rem;
        grid-template-columns: 1fr 250px;
        gap: 1rem;
        p {
          color: white;
        }
        .input {
          display: flex;
          align-items: center;
          p {
            margin-right: 0.5rem;
            margin-bottom: 0;
          }
          input {
            width: 100%;
            border: 1px dashed var(--primary-darken);
            background-color: transparent;
          }
        }
      }
      .middle {
        padding-top: 2rem;
        .operand {
          width: 100%;
          display: flex;
          justify-content: center;
          padding: 0.8rem 0;
          font-size: 12px;
          cursor: pointer;
          &.selected, &:hover {
            background-color: var(--primary);
            color: var(--text-inverse);
          }
        }
      }
      .content {
        display: grid;
        grid-template-columns: 1fr 50px 1fr;
        gap: 1rem;
        height: 100%;
        .middle {
          background-color: var(--background-darken);
          height: 100%;
        }
        @media(max-width: 800px) {
          grid-template-columns: 1fr;
          .middle {
            padding-top: 0;
            display: flex;
            height: 40px;
          }
        }
        .indicator {
          padding: 2rem 1rem;
        }
      }
      .input {
        p {
          margin-bottom: 0.5rem;
          font-size: 14px;
          font-weight: 600;
        }
        input {
          padding: 0.5rem;
          outline: none;
          border: 1px solid var(--border-color);
          color: var(--text-inverse);
        }
      }
    }
  }
}
</style>
