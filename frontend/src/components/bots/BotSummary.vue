<template>
  <div
    class="bot-summary"
    :class="bot.status"
    @click="() => notClickable ? undefined : $router.push(`/bots/${bot.id}`)"
  >
    <div class="left">
      <exchange-tag
        v-if="bot.type === 'LIVE' && exchangeConnection"
        :exchange="exchangeConnection.exchange"
        :name="exchangeConnection.exchange"
      />
      <exchange-tag
        v-else
        :exchange="bot.strategy.exchange.toUpperCase()"
        name="TEST"
      />
      <p class="strategy">{{ bot.strategy.name }}</p>
      <p class="running">{{ runningTime }}</p>
    </div>
    <div class="right">
      <p class="balance">{{ bot.startBalance }} {{ bot.quoteCurrency }} => {{ bot.currentBalance }} {{ bot.quoteCurrency }}</p>
      <div class="status" :class="bot.status">
        <p>{{ Number(changePercentage.toFixed(2)) }}%</p>
        <div class="circle"></div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onMounted, PropType, ref } from '@vue/runtime-core'
import { useStore } from 'vuex'
import humanizeDuration from 'humanize-duration'

import ExchangeTag from '@/components/bots/ExchangeTag.vue'
import Bot from '@/types/Bot'
import { ExchangeConnection } from '@/types/Exchange'

export default defineComponent({
  components: { ExchangeTag },
  props: {
    bot: {
      type: Object as PropType<Bot>,
      required: true
    },
    notClickable: {
      type: Boolean
    }
  },
  setup (props) {
    const store = useStore()
    const exchangeConnection = computed(() => {
      return store.getters.exchangeConnections
        .find((e: ExchangeConnection) => e.id === props.bot.exchangeConnectionId)
    })

    const currentTimestamp = ref<number>(Date.now())
    const changePercentage = computed(() => (props.bot.currentBalance - props.bot.startBalance) / props.bot.startBalance * 100)

    onMounted(() => {
      const interval = setInterval(() => {
        currentTimestamp.value = Date.now()
      }, 5000)

      return () => clearInterval(interval)
    })

    const runningTime = computed(() => {
      let timeDiff = currentTimestamp.value - new Date(props.bot.startDate).getTime()
      if (props.bot.status === 'ended') {
        timeDiff = new Date(props.bot.endDate || new Date()).getTime() - new Date(props.bot.startDate).getTime()
      }
      const units: any[] = ['w', 'd', 'h']

      if (timeDiff < 120 * 60 * 1000) units.push('m')

      return humanizeDuration(timeDiff, {
        units,
        maxDecimalPoints: 0
      })
    })

    return {
      exchangeConnection,
      runningTime,
      changePercentage
    }
  }
})
</script>

<style lang="scss" scoped>
.bot-summary {
  border: 1px solid var(--border-color);
  padding: 0.75rem;
  padding-bottom: 0.25rem;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  transition: .2s;
  &:hover {
    border-color: var(--primary);
  }
  .left {
    margin-right: 1rem;
  }
  .left, .right {
    display: flex;
    align-items: center;
    margin-bottom: 0.5rem;
    > * {
      margin-right: 0.75rem;
    }
  }
  .running {
    font-size: 13px;
    color: var(--text-tertiary);
    margin: 0;
  }
  .balance {
    font-size: 14px;
  }
  .status {
    display: flex;
    align-items: center;
    p {
      margin-right: 0.5rem;
      color: var(--disabled);
    }
    .circle {
      --size: 15px;
      width: var(--size);
      height: var(--size);
      background-color: var(--disabled);
      border-radius: var(--size);
      position: relative;
    }
    &.online {
      p {
        color: var(--primary);
      }
      .circle {
        background-color: var(--primary);
        &:after {
          content: '';
          position: absolute;
          left: -50%;
          top: -50%;
          width: 200%;
          height: 200%;
          border-radius: 45px;
          background-color: var(--primary);
          animation: pulse-ring 1.25s cubic-bezier(0.215, 0.61, 0.355, 1) infinite;
        }
      }
    }
    &.ended {
      p {
        color: var(--primary-darken);
      }
      .circle {
        background-color: var(--primary-darken);
      }
    }
  }
}

.ended {
  background-color: var(--background-darken);
}

@keyframes pulse-ring {
  0% {
    transform: scale(.33);
  }
  80%, 100% {
    opacity: 0;
  }
}
</style>
