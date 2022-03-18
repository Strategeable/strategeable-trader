<template>
  <div class="exchange-balance">
    <div class="left">
      <img :src="`/img/${exchange}.png`" alt="">
      <p>{{ exchange }}</p>
    </div>
    <div class="right">
      <p class="btc">{{ value.value }} {{ value.asset }}</p>
      <p class="estimation">~ ${{ Number(value.usdValue.toFixed(2)) }}</p>
    </div>
  </div>
</template>

<script lang="ts">
import { Exchange, ExchangeValue } from '@/types/Exchange'
import { defineComponent, PropType } from 'vue'

export default defineComponent({
  props: {
    exchange: {
      type: String as PropType<Exchange>,
      required: true
    },
    value: {
      type: Object as PropType<ExchangeValue>,
      required: true
    }
  }
})
</script>

<style lang="scss" scoped>
.exchange-balance {
  padding: 0.6rem 0.8rem;
  background-color: var(--background-lighten);
  border: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
  .left {
    display: flex;
    align-items: center;
    img {
      width: 32px;
      margin-right: 1rem;
    }
    p {
      text-transform: capitalize;
    }
  }
  .right {
    display: flex;
    align-items: center;
    .btc {
      font-weight: bold;
      margin-right: 1rem;
    }
    .estimation {
      font-size: 14px;
    }
  }
}
</style>
