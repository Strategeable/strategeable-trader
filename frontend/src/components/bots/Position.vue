<template>
  <div class="position">
    <div class="left">
      <div class="symbol">
        <div class="images">
          <img :src="baseCurrencyImage" :alt="baseCurrency">
          <img :src="quoteCurrencyImage" :alt="quoteCurrency">
        </div>
        <p>{{ position.symbol }}</p>
      </div>
      <p class="hold-time">{{ holdTime }}</p>
    </div>
    <div class="right">
      <p class="value">{{ Number(realPosition.getEntryQuoteSize().toFixed(2)) }} => {{ Number(realPosition.getEntryQuoteSize().toFixed(2)) }}</p>
      <p class="change">10.26%</p>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onMounted, PropType, ref } from 'vue'
import humanizeDuration from 'humanize-duration'

import Position from '@/types/Position'
import PositionHandler from '@/handlers/Position'

export default defineComponent({
  props: {
    position: {
      type: Object as PropType<Position>,
      required: true
    }
  },
  setup (props) {
    const realPosition = new PositionHandler(props.position)
    const [baseCurrency, quoteCurrency] = props.position.symbol.split('/')
    const baseCurrencyImage = `https://cdn.jsdelivr.net/npm/cryptocurrency-icons@0.16.1/svg/color/${baseCurrency.toLowerCase()}.svg`
    const quoteCurrencyImage = `https://cdn.jsdelivr.net/npm/cryptocurrency-icons@0.16.1/svg/color/${quoteCurrency.toLowerCase()}.svg`

    const currentTimestamp = ref<number>(Date.now())

    const holdTime = computed(() => {
      const timeDiff = currentTimestamp.value - new Date(props.position.openTime).getTime()
      const units: any[] = ['w', 'd', 'h']

      if (timeDiff < 120 * 60 * 1000) units.push('m')

      return humanizeDuration(timeDiff, {
        units,
        maxDecimalPoints: 0
      })
    })

    onMounted(() => {
      const interval = setInterval(() => {
        currentTimestamp.value = Date.now()
      }, 5000)

      return () => clearInterval(interval)
    })

    return {
      quoteCurrency,
      baseCurrency,
      baseCurrencyImage,
      quoteCurrencyImage,
      holdTime,
      realPosition
    }
  }
})
</script>

<style lang="scss" scoped>
.position {
  padding: 1rem;
  border: 1px solid var(--border-color);
  background-color: var(--background-lighten);
  display: flex;
  align-items: center;
  justify-content: space-between;
  .left, .right {
    display: flex;
    align-items: center;
  }
  .symbol {
    display: flex;
    align-items: center;
    margin-right: 1rem;
    p {
      font-weight: bold;
    }
    .images {
      margin-right: 1rem;
      img {
        z-index: 10;
        width: 30px;
      }
      :nth-child(2) {
        margin-left: -10px;
        z-index: 0;
      }
    }
  }
  .hold-time {
    color: var(--text-tertiary);
    font-size: 14px;
  }
  .value {
    margin-right: 1rem;
    font-size: 14px;
  }
  .change {
    color: var(--positive);
    font-weight: bold;
    &.negative {
      color: var(--negative);
    }
  }
}
</style>
