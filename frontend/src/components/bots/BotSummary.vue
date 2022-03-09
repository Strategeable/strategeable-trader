<template>
  <div class="bot-summary">
    <div class="left">
      <exchange-tag
        :exchange="'BINANCE'"
        name="Binance main"
      />
      <p class="strategy">Trend Hero</p>
    </div>
  </div>
</template>

<script lang="ts">
import ExchangeTag from '@/components/bots/ExchangeTag.vue'
import { defineComponent } from '@vue/runtime-core'

export default defineComponent({
  components: { ExchangeTag }
})
</script>

<style lang="scss" scoped>
.bot-summary {
  border: 1px solid var(--border-color);
  padding: 0.75rem;
  display: flex;
  align-items: center;
  cursor: pointer;
  &:hover {
    border-color: var(--primary);
    transform: scale(1.01);
  }
  .left, .right {
    display: flex;
    align-items: center;
    > * {
      margin-right: 0.75rem;
    }
  }
}
</style>
