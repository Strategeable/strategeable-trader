<template>
  <div class="dashboard">
    <div class="balance-summary">
      <div class="summary">
        <p class="header">Overall balance</p>
        <p class="balance">0.892342 BTC</p>
        <p class="estimation">~ $34212.12</p>
      </div>
      <div class="summary">
        <p class="header">Managed by bots</p>
        <p class="balance">0.892342 BTC</p>
        <p class="estimation">~ $34212.12</p>
      </div>
      <div class="summary">
        <p class="header">Balance + open positions</p>
        <p class="balance">0.892342 BTC</p>
        <p class="estimation">~ $34212.12</p>
      </div>
    </div>
    <div class="content">
      <div class="wallets section">
        <h2>Wallets</h2>
        <div class="exchanges">
          <exchange-balance :exchange="'binance'"/>
          <exchange-balance :exchange="'kucoin'"/>
        </div>
      </div>
      <div class="bots section">
        <h2>Running bots</h2>
        <div class="items">
          <bot-summary
            v-for="bot in bots"
            :key="bot.id"
            :bot="bot"
            :compact="true"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from 'vue'
import { useStore } from 'vuex'

import Bot from '@/types/Bot'
import BotSummary from '@/components/bots/BotSummary.vue'
import ExchangeBalance from '@/components/exchange/ExchangeBalance.vue'

export default defineComponent({
  components: {
    BotSummary, ExchangeBalance
  },
  setup () {
    const store = useStore()
    const bots = computed<Bot[]>(() => store.getters.bots.filter((b: Bot) => b.status === 'online'))

    return {
      bots
    }
  }
})
</script>

<style lang="scss" scoped>
.balance-summary {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  margin-bottom: 3rem;
  padding: 1.5rem;
  background-color: var(--background-lighten);
  border: 1px solid var(--border-color);
  .summary {
    margin-right: 5rem;
    .header {
      margin-bottom: 1rem;
      font-size: 15px;
    }
    .balance {
      font-weight: bold;
      margin-bottom: 0.25rem;
      font-size: 17px;
    }
    .estimation {
      font-size: 14px;
    }
  }
}
.content {
  display: grid;
  grid-template-columns: 3fr 4fr;
  gap: 2rem;
  .bots {
    .items {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 1rem;
    }
  }
  .exchanges {
    display: grid;
    grid-template-columns: 1fr;
    gap: 0.5rem;
  }
}
</style>
