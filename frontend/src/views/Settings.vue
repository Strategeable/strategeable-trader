<template>
  <div class="settings">
    <h2>Settings</h2>
    <div class="section">
      <h3>Exchange connections</h3>
      <div class="connections">
        <exchange-connection-comp
          v-for="conn in exchangeConnections"
          :key="conn.id"
          :exchangeConnection="conn"
          @delete="() => deleteConnection(conn.id)"
        />
      </div>
      <button @click="() => openCreateConnection = true">Create connection</button>
    </div>
    <default-popup @close="() => openCreateConnection = false" v-if="openCreateConnection">
      <div class="popup">
        <h3>Create exchange connection</h3>
        <div class="input">
          <p>Exchange</p>
          <v-select
            :options="exchanges"
            v-model="exchangeConnection.exchange"
          />
        </div>
        <div class="input">
          <p>Name</p>
          <input
            type="text"
            v-model="exchangeConnection.name"
          >
        </div>
        <div class="input">
          <p>API key</p>
          <input
            type="text"
            v-model="exchangeConnection.apiKey"
          >
        </div>
        <div class="input">
          <p>API secret</p>
          <input
            type="text"
            v-model="exchangeConnection.apiSecret"
          >
        </div>
        <button
          :disabled="!validCreateConnection"
          @click="createExchangeConnection"
        >Create</button>
        <p class="error" v-if="createConnectionError">{{ createConnectionError }}</p>
      </div>
    </default-popup>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref } from 'vue'
import { useStore } from 'vuex'

import ExchangeConnectionComp from '@/components/settings/ExchangeConnection.vue'
import DefaultPopup from '@/components/popups/DefaultPopup.vue'
import { ExchangeConnection, exchanges } from '@/types/Exchange'

export default defineComponent({
  components: { ExchangeConnectionComp, DefaultPopup },
  setup () {
    const store = useStore()
    const exchangeConnections = computed(() => store.getters.exchangeConnections)

    const openCreateConnection = ref<boolean>(false)
    const createConnectionError = ref<string>()
    const exchangeConnection = ref<ExchangeConnection>({
      exchange: 'binance',
      name: '',
      createdOn: '',
      apiKey: '',
      apiSecret: ''
    })
    const validCreateConnection = computed(() => {
      if (exchangeConnection.value.name === '') return false
      if (exchangeConnection.value.apiKey === '') return false
      if (exchangeConnection.value.apiSecret === '') return false
      if (!exchangeConnection.value.exchange) return false
      return true
    })

    async function createExchangeConnection () {
      if (!validCreateConnection.value) return

      const result = await store.dispatch('addExchangeConnection', exchangeConnection.value)
      if (result.error) {
        createConnectionError.value = result.error
      } else {
        openCreateConnection.value = false
        exchangeConnection.value = {
          exchange: 'binance',
          apiKey: '',
          apiSecret: '',
          name: '',
          createdOn: ''
        }
      }
    }

    function deleteConnection (id: string) {
      store.dispatch('deleteExchangeConnection', id)
    }

    return {
      exchangeConnections,
      openCreateConnection,
      exchangeConnection,
      exchanges,
      validCreateConnection,
      createConnectionError,
      createExchangeConnection,
      deleteConnection
    }
  }
})
</script>

<style lang="scss" scoped>
.settings {
  h2 {
    margin-bottom: 1.5rem;
  }
  .connections {
    margin-bottom: 1rem;
  }
  .popup {
    h3 {
      margin-bottom: 1rem;
    }
    input {
      &:not(.v-select) {
        width: 100%;
      }
    }
    .error {
      color: red;
      padding-top: 2rem;
    }
  }
}
</style>
