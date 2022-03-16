<template>
  <default-popup @close="$emit('close')">
    <h3>Launch a new bot</h3>
    <div class="input">
      <p>Type</p>
      <v-select
        :options="['TEST']"
        v-model="type"
      />
    </div>
    <div class="input">
      <p>Strategy</p>
      <v-select
        :options="strategies"
        :reduce="x => x.id"
        label="name"
        v-model="strategy"
      />
    </div>
    <div class="input" v-if="type === 'LIVE'">
      <p>Account</p>
      <v-select
        :options="exchangeConnections"
        :reduce="x => x.id"
        label="name"
        v-model="exchangeConnection"
      />
    </div>
    <div class="input">
      <p>Start balance <span v-if="strategy">({{ strategyQuoteCurrency }})</span></p>
      <input
        v-model="startBalance"
        type="number"
      >
    </div>
    <button
      @click="launch"
      :disabled="!valid"
    >Launch bot</button>
  </default-popup>
</template>

<script lang="ts">
import { computed, defineComponent, ref } from 'vue'
import { useStore } from 'vuex'

import DefaultPopup from '@/components/popups/DefaultPopup.vue'
import { Strategy } from '@/types/Strategy'

export default defineComponent({
  emits: ['close'],
  components: { DefaultPopup },
  setup (props, context) {
    const store = useStore()
    const strategies = computed(() => store.getters.strategies)
    const exchangeConnections = computed(() => store.getters.exchangeConnections)

    const errorRef = ref<string>()

    const type = ref<string>('TEST')
    const strategy = ref<string>()
    const startBalance = ref<number>()
    const exchangeConnection = ref<string>()
    const valid = computed(() => {
      if (type.value !== 'TEST' && type.value !== 'LIVE') return false
      if (!strategy.value) return false
      if (!startBalance.value) return false
      return true
    })

    const strategyQuoteCurrency = computed(() => {
      if (!strategy.value) return undefined
      const strat = strategies.value.find((s: Strategy) => s.id === strategy.value)
      if (!strat) return undefined
      return strat.quoteCurrency
    })

    async function launch () {
      const params = {
        type: type.value,
        strategyId: strategy.value,
        startBalance: startBalance.value,
        exchangeConnection: exchangeConnection.value
      }

      const result = await store.dispatch('launchBot', params)

      if (result.error) {
        errorRef.value = result.error
      } else {
        context.emit('close')
      }
    }

    return {
      type,
      strategies,
      strategy,
      startBalance,
      exchangeConnections,
      exchangeConnection,
      valid,
      strategyQuoteCurrency,
      launch
    }
  }
})
</script>

<style lang="scss" scoped>
h3 {
  margin-bottom: 1.5rem;
}

input {
  &:not(.v-select) {
    width: 100%;
  }
}

button {
  margin-top: 1rem;
}
</style>
