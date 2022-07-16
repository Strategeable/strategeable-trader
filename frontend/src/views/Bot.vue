<template>
  <div class="bot" v-if="bot">
    <div class="top">
      <router-link to="/bots">
        <fa-icon icon="arrow-left"/> Return to bots
      </router-link>
      <div class="settings">
        <fa-icon icon="sliders-h"/>
        <p>Bot settings</p>
      </div>
    </div>
    <bot-summary
      :bot="bot"
      :notClickable="false"
    />
    <div class="grid">
      <div class="stats">
        <h2>Stats</h2>
        <div class="chart">
          Chart here
        </div>
        <div class="misc">
          <div class="stat">
            <p>Win rate</p>
            <p>0.78</p>
          </div>
          <div class="stat">
            <p>Win rate</p>
            <p>0.78</p>
          </div>
          <div class="stat">
            <p>Win rate</p>
            <p>0.78</p>
          </div>
          <div class="stat">
            <p>Win rate</p>
            <p>0.78</p>
          </div>
          <div class="stat">
            <p>Win rate</p>
            <p>0.78</p>
          </div>
          <div class="stat">
            <p>Win rate</p>
            <p>0.78</p>
          </div>
        </div>
      </div>
      <div class="positions">
        <div v-if="bot && bot.status !== 'ended'">
          <h2>Running positions</h2>
          <div class="list" v-if="openPositions.length > 0">
            <position-comp
              v-for="position in openPositions"
              :key="position.id"
              :position="position"
            />
          </div>
          <div class="placeholder" v-else>
            No running positions
          </div>
        </div>
        <h2>Finished positions</h2>
        <div class="list" v-if="closedPositions.length > 0">
          <finished-positions :positions="closedPositions"/>
        </div>
        <div v-else>
          No finished positions
        </div>
      </div>
    </div>
  </div>
  <div v-else>
    Bot not found.
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from 'vue'

import { useStore } from 'vuex'
import { useRoute } from 'vue-router'

import Position from '@/types/Position'
import Bot from '@/types/Bot'
import BotSummary from '@/components/bots/BotSummary.vue'
import PositionComp from '@/components/bots/Position.vue'
import FinishedPositions from '@/components/bots/FinishedPositions.vue'

export default defineComponent({
  components: { BotSummary, PositionComp, FinishedPositions },
  setup () {
    const store = useStore()
    const route = useRoute()
    const botId: any = route.params.id
    const bot = computed<Bot>(() => store.getters.bots.find((b: Bot) => b.id === botId))
    const positions = computed<Position[]>(() => store.getters.positions[botId] || [])
    const openPositions = computed<Position[]>(() => positions.value.filter(p => ['OPEN', 'OPENING'].includes(p.state)))
    const closedPositions = computed<Position[]>(() => positions.value.filter(p => ['CLOSED', 'CLOSING'].includes(p.state)))

    return {
      bot,
      openPositions,
      closedPositions
    }
  }
})
</script>

<style lang="scss" scoped>
.bot {
  .top {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    a {
      display: block;
      color: var(--text);
      text-decoration: none;
    }
    .settings {
      display: flex;
      align-items: center;
      color: var(--primary);
      border: 1px solid transparent;
      cursor: pointer;
      padding: 0.5rem;
      transition: .3s;
      p {
        color: var(--primary);
        font-weight: 500;
      }
      svg {
        margin-right: 0.5rem;
      }
      &:hover {
        background-color: var(--background-darken);
        border: 1px solid var(--border-color);
      }
    }
  }

  .grid {
    display: grid;
    grid-template-columns: 500px 3fr;
    gap: 2.5rem;
    margin-top: 2rem;
    h2 {
      margin-bottom: 1rem;
    }
    .list, .placeholder {
      margin-bottom: 2.5rem;
    }
  }
  .chart {
    width: 100%;
    height: 350px;
    background-color: var(--background-lighten);
    border: 1px solid var(--border-color);
    padding: 1rem;
    margin-bottom: 2rem;
  }
  .misc {
    display: grid;
    grid-template-columns: 1fr 1fr;
    .stat {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 1rem;
      border-bottom: 1px solid var(--border-color);
      :first-child {
        font-weight: bold;
      }
    }
  }
}
</style>
