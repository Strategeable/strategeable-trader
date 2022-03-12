<template>
  <div class="backtest-result">
    Backtest result - from: {{ backtest.startBalance }} to: {{ backtest.endBalance }}
  </div>
</template>

<script lang="ts">
import { BacktestResult } from '@/types/Backtest'
import { defineComponent, PropType } from 'vue'

export default defineComponent({
  props: {
    backtest: {
      type: Object as PropType<BacktestResult>,
      required: true
    }
  },
  setup () {
    return {
    }
  }
})
</script>
