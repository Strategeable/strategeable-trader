<template>
  <div class="bot-summary">
    <div class="left">
      <exchange-tag
        v-if="bot.type === 'LIVE'"
        :exchange="bot.exchange"
        :name="'name (todo)'"
      />
      <exchange-tag
        v-else
        :exchange="bot.strategy.exchange.toUpperCase()"
        name="TEST"
      />
      <p class="strategy">{{ bot.strategy.name }}</p>
    </div>
  </div>
</template>

<script lang="ts">
import ExchangeTag from '@/components/bots/ExchangeTag.vue'
import Bot from '@/types/Bot'
import { defineComponent, PropType } from '@vue/runtime-core'

export default defineComponent({
  components: { ExchangeTag },
  props: {
    bot: {
      type: Object as PropType<Bot>,
      required: true
    }
  }
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
