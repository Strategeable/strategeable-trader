<template>
  <div class="strategies">
    <button @click="create">Create strategy</button>
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
import { useRouter } from 'vue-router'
import { useStore } from '@/store'

export default defineComponent({
  components: { StrategiesTable },
  setup () {
    const router = useRouter()
    const store = useStore()
    const strategies = computed(() => store.getters.strategies)

    function create () {
      router.push('/strategies/new')
    }

    return {
      strategies,
      create
    }
  }
})
</script>

<style lang="scss" scoped>
.table {
  margin-top: 2rem;
}
</style>
