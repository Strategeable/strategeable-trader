<template>
  <div class="exchange-balance">
    <div class="left">
      <img :src="`/img/${exchange}.png`" alt="">
      <p>{{ exchange }}</p>
    </div>
    <div class="right">
      <p class="btc">0.89234 BTC</p>
      <p class="estimation">~ $34212.12</p>
    </div>
    {{ balances }}
    {{ rates }}
  </div>
</template>

<script lang="ts">
import { ExchangeBalance } from '@/types/Exchange'
import { computed, defineComponent } from 'vue'
import { useStore } from 'vuex'

export default defineComponent({
  props: {
    exchange: {
      type: String,
      required: true
    }
  },
  setup (props) {
    const store = useStore()
    const balances = computed<ExchangeBalance[]>(() => store.getters.balances[props.exchange])
    const rates = computed(() => store.getters.rates)

    return {
      balances,
      rates
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
