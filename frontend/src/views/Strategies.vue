<template>
  <div class="strategies">
    <button>Create strategy</button>
    <div class="table">
      <strategies-table
        :data="strategies"
        dataKey="id"
        @select="id => $router.push(`/strategies/${id}`)"
      />
    </div>
  </div>
</template>

<script lang="ts">
import StrategiesTable from '@/components/strategies/Table.vue'
import { defineComponent } from '@vue/runtime-core'

export default defineComponent({
  components: { StrategiesTable },
  setup () {
    const strategies = [
      { id: '1', name: 'Trend Hero', markets: 'BTC, USDT', running: 1, lastEdited: new Date() },
      { id: '2', name: 'Trend Hero', markets: 'BTC, USDT', running: 1, lastEdited: new Date() },
      { id: '3', name: 'Trend Hero', markets: 'BTC, USDT', running: 1, lastEdited: new Date() }
    ]

    return {
      strategies
    }
  }
})
</script>

<style lang="scss" scoped>
.table {
  margin-top: 2rem;
}
</style>
