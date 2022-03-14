<template>
  <div class="exchange-tag">
    <img :src="require(`@/assets/img/exchange/${exchange.toLowerCase()}.png`)" alt="Binance logo">
    <div class="name">
      <p>{{ name }}</p>
    </div>
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
    }
  }
})
</script>

<style lang="scss" scoped>
.exchange-tag {
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
}
</style>
