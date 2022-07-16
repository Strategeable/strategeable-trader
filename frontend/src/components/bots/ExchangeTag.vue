<template>
  <div class="exchange-tag">
    <img :src="require(`@/assets/img/exchange/${exchange.toLowerCase()}.png`)" alt="Binance logo">
    <div class="name">
      <p>{{ name }}</p>
    </div>
    <img class="quote-currency" :src="quoteCurrencyImage" alt="" />
  </div>
</template>

<script lang="ts">
import { Exchange } from '@/types/Exchange'
import { defineComponent, PropType } from 'vue'

export default defineComponent({
  props: {
    name: {
      required: true,
      type: String
    },
    exchange: {
      type: String as PropType<Exchange>,
      required: true
    },
    quoteCurrency: {
      type: String,
      default: 'USDT'
    }
  },
  setup (props) {
    const quoteCurrencyImage = `https://cdn.jsdelivr.net/npm/cryptocurrency-icons@0.16.1/svg/color/${props.quoteCurrency.toLowerCase()}.svg`

    return {
      quoteCurrencyImage
    }
  }
})
</script>

<style lang="scss" scoped>
.exchange-tag {
  position: relative;
  display: inline-flex;
  align-items: center;
  background-color: var(--tag-color);
  img {
    width: 32px;
    height: 32px;
  }
  .name {
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 0.5rem;
    p {
      color: var(--text-tertiary);
      font-size: 13px;
    }
  }
  .quote-currency {
    position: absolute;
    right: 0;
    bottom: 0;
    width: 16px;
    height: 16px;
    transform: translate(25%, 25%);
  }
}
</style>
