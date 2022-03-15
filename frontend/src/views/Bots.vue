<template>
  <div class="bots">
    <div class="section">
      <h2>Active</h2>
      <div class="active-bots">
        <bot-summary
          v-for="bot in activeBots"
          :key="bot.id"
          :bot="bot"
        />
      </div>
      <button @click="toggleLaunch">Launch new bot</button>
    </div>
    <div class="section" v-if="finishedBots.length > 0">
      <h2>Finished</h2>
      <div class="finished-bots">
        <bot-summary
          v-for="bot in finishedBots"
          :key="bot.id"
          :bot="bot"
        />
      </div>
    </div>
    <launch-new-bot
      v-if="launchBot"
      @close="toggleLaunch"
    />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref } from '@vue/runtime-core'
import { useStore } from 'vuex'

import BotSummary from '@/components/bots/BotSummary.vue'
import LaunchNewBot from '@/components/bots/LaunchNewBot.vue'
import Bot from '@/types/Bot'

export default defineComponent({
  components: { BotSummary, LaunchNewBot },
  setup () {
    const store = useStore()
    const bots = computed<Bot[]>(() => store.getters.bots)
    const activeBots = computed<Bot[]>(() => bots.value.filter((b: Bot) => b.status !== 'ended'))
    const finishedBots = computed<Bot[]>(() => bots.value.filter((b: Bot) => b.status === 'ended'))

    const launchBot = ref<boolean>(false)

    function toggleLaunch () {
      launchBot.value = !launchBot.value
    }

    return {
      activeBots,
      finishedBots,
      launchBot,
      toggleLaunch
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

.bots {
  button {
    margin-top: 1.5rem;
  }
}
</style>
