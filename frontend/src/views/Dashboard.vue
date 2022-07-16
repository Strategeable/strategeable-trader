<template>
  <div class="dashboard">
    <div class="balance-summary">
      <div class="summary">
        <p class="header">Overall balance</p>
        <p class="balance">{{ totalValues.value }} {{ totalValues.asset }}</p>
        <p class="estimation">~ ${{ totalValues.usdValue }}</p>
      </div>
      <div class="summary">
        <p class="header">Managed by bots</p>
        <p class="balance">- BTC</p>
        <p class="estimation">~ $-</p>
      </div>
      <div class="summary">
        <p class="header">Balance + open positions</p>
        <p class="balance">- BTC</p>
        <p class="estimation">~ $-</p>
      </div>
    </div>
    <div class="content">
      <div class="wallets section">
        <h2>Wallets</h2>
        <div class="exchanges">
          <exchange-balance
            v-for="[exchange, value] in Object.entries(valuePerExchange)"
            :key="exchange"
            :exchange="(exchange as any)"
            :value="value"
          />
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
import { useStore } from '@/store'

import Bot from '@/types/Bot'
import BotSummary from '@/components/bots/BotSummary.vue'
import ExchangeBalance from '@/components/exchange/ExchangeBalance.vue'
import { ExchangeValue } from '@/types/Exchange'

export default defineComponent({
  components: {
    BotSummary, ExchangeBalance
  },
  setup () {
    const store = useStore()
    const bots = computed(() => store.getters.bots.filter((b: Bot) => b.status === 'online'))
    const rates = computed(() => store.getters.rates)
    const balances = computed(() => store.getters.balances)
    const denominateIn = computed(() => store.getters.denominateIn)
    const exchanges = computed(() => Array.from(
      new Set(store.getters.exchangeConnections.map(e => e.exchange))
    ))

    const valuePerExchange = computed(() => {
      const value: Record<string, ExchangeValue> = {}

      for (const exchange of exchanges.value) {
        if (!rates.value) continue

        let total = 0
        for (const balance of balances.value.filter(b => b.exchange === exchange)) {
          if (balance.asset === denominateIn.value) {
            total += balance.amount
            continue
          }

          const rate = rates.value.find(r => r.asset === balance.asset)
          console.log(rate, balance.asset)
          if (!rate) continue
          if (!rate.quote[denominateIn.value]) {
            const oppositeRate = rates.value.find(r => r.asset === denominateIn.value)
            if (!oppositeRate) continue
            const price = oppositeRate.quote[balance.asset]
            if (!price) continue

            total += balance.amount / price
          } else {
            total += rate.quote[denominateIn.value] * balance.amount
          }
        }

        const toUsdRate = rates.value.find(r => r.asset === denominateIn.value)
        let usdValue = 0
        if (toUsdRate) usdValue = total * (toUsdRate.quote.USDT || toUsdRate.quote.USDC || toUsdRate.quote.DAI || toUsdRate.quote.BUSD || 0)

        value[exchange] = {
          asset: denominateIn.value,
          value: Number(total.toFixed(store.getters.getAssetRounding(denominateIn.value))),
          usdValue: Number(usdValue.toFixed(store.getters.getAssetRounding('USD')))
        }
      }

      return value
    })

    const totalValues = computed(() => {
      const value: ExchangeValue = {
        asset: denominateIn.value,
        usdValue: 0,
        value: 0
      }

      for (const val of Object.values(valuePerExchange.value)) {
        value.usdValue += val.usdValue
        value.value += val.value
      }

      return value
    })

    return {
      bots,
      valuePerExchange,
      totalValues
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
