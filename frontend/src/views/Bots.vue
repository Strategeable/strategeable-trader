<template>
  <div class="bots">
    <div class="section">
      <h2>Active</h2>
      <div class="active-bots">
        <bot-summary/>
        <bot-summary/>
        <bot-summary/>
      </div>
    </div>
    <div class="section">
      <h2>Finished</h2>
      <div class="finished-bots">
        <bot-summary/>
        <bot-summary/>
        <bot-summary/>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import BotSummary from '@/components/bots/BotSummary.vue'
import { defineComponent } from '@vue/runtime-core'

export default defineComponent({
  components: { BotSummary }
})
</script>

<style lang="scss" scoped>
.active-bots, .finished-bots {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.5rem 1rem;
}
</style>
