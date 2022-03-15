<template>
  <default-popup @close="$emit('close')">
    <h3>Launch a new bot</h3>
    <div class="input">
      <p>Type</p>
      <v-select
        :options="['TEST', 'LIVE']"
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
      <p>Start balance</p>
      <input
        v-model="startBalance"
        type="number"
      >
    </div>
    <button>Launch bot</button>
  </default-popup>
</template>

<script lang="ts">
import { computed, defineComponent, ref } from 'vue'
import { useStore } from 'vuex'

import DefaultPopup from '@/components/popups/DefaultPopup.vue'

export default defineComponent({
  emits: ['close'],
  components: { DefaultPopup },
  setup () {
    const store = useStore()
    const strategies = computed(() => store.getters.strategies)
    const exchangeConnections = computed(() => store.getters.exchangeConnections)

    const type = ref<string>('TEST')
    const strategy = ref<string>()
    const startBalance = ref<number>()
    const exchangeConnection = ref<string>()

    return {
      type,
      strategies,
      strategy,
      startBalance,
      exchangeConnections,
      exchangeConnection
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
