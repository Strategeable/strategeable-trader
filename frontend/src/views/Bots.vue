<template>
  <div class="bots">
    <div class="section">
      <h2>Active</h2>
      <div class="active-bots">
        <bot-summary
          v-for="bot in activeBots"
          :key="bot.id"
        />
      </div>
      <button>Launch new bot</button>
    </div>
    <div class="section" v-if="finishedBots.length > 0">
      <h2>Finished</h2>
      <div class="finished-bots">
        <bot-summary
          v-for="bot in finishedBots"
          :key="bot.id"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import BotSummary from '@/components/bots/BotSummary.vue'
import Bot from '@/types/Bot'
import { computed, defineComponent } from '@vue/runtime-core'
import { useStore } from 'vuex'

export default defineComponent({
  components: { BotSummary },
  setup () {
    const store = useStore()
    const bots = computed<Bot[]>(() => store.getters.bots)
    const activeBots = computed<Bot[]>(() => bots.value.filter((b: Bot) => b.status !== 'ended'))
    const finishedBots = computed<Bot[]>(() => bots.value.filter((b: Bot) => b.status === 'ended'))

    return {
      activeBots,
      finishedBots
    }
  }
})
</script>

<style lang="scss" scoped>
.active-bots, .finished-bots {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.5rem 1rem;
}
</style>
