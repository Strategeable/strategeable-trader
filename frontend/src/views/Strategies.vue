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
import { computed, defineComponent } from '@vue/runtime-core'
import { useStore } from 'vuex'

export default defineComponent({
  components: { StrategiesTable },
  setup () {
    const store = useStore()
    const strategies = computed(() => store.getters.strategies)

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
