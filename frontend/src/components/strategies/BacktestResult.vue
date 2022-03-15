<template>
  <div class="backtest-result">
    <div class="banner">
      <span class="change">{{ backtestData.change }}%</span> {{ backtest.strategy.name }}
    </div>
    <div class="wrapper">
      <div class="chart">
        <line-chart v-bind="lineChartProps"/>
      </div>
      <div class="data">
        <h3>Strategy name: {{ backtest.strategy.name }}</h3>
        <p
          class="result"
        >Result <span :class="{ negative: backtestData.change < 0 }">{{ backtestData.change }}%</span></p>
        <p>{{ backtestData.winsLosses.wins }} wins / {{ backtestData.winsLosses.losses }} losses (win rate: {{ Number((backtestData.winsLosses.winRate).toFixed(2)) }})</p>
        <p>End balance: {{ backtest.endBalance }}</p>
        <button @click="$emit('restore')">Restore strategy</button>
        <p class="backtest-date">Backtested on {{ moment(backtest.startedOn).format('DD MMM HH:mm') }}</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { BacktestResult, Position } from '@/types/Backtest'
import { computed, defineComponent, PropType } from 'vue'
import { LineChart, useLineChart } from 'vue-chart-3'
import { Chart, ChartData, registerables } from 'chart.js'
import 'chartjs-adapter-moment'
import moment from 'moment'

import BacktestPosition from '@/handlers/BacktestPosition'
import { useStore } from 'vuex'

interface LineChartEntry {
  y: number
  x: Date
}

Chart.register(...registerables)
export default defineComponent({
  emits: ['restore'],
  components: { LineChart },
  props: {
    backtest: {
      type: Object as PropType<BacktestResult>,
      required: true
    }
  },
  setup (props) {
    const store = useStore()
    const theme = computed(() => store.getters.theme)
    const chartColors = computed(() => {
      return {
        xYAxis: theme.value === 'light' ? 'black' : 'white',
        lines: theme.value === 'light' ? '#eeeeff' : '#28283a'
      }
    })

    const backtestData = computed(() => {
      return {
        change: Number(((props.backtest.endBalance - props.backtest.startBalance) / props.backtest.startBalance * 100).toFixed(2)),
        fromDate: new Date(props.backtest.fromDate),
        toDate: new Date(props.backtest.toDate),
        winsLosses: calculateWinRate(props.backtest.positions)
      }
    })

    function calculateWinRate (positions: Position[]): { winRate: number, wins: number, losses: number } {
      if (!positions) return { winRate: 0, wins: 0, losses: 0 }

      let wins = 0
      let losses = 0

      for (const position of positions) {
        const pos = new BacktestPosition(position)
        if (pos.getChangePercent() >= 0) {
          wins += 1
        } else {
          losses += 1
        }
      }

      return { winRate: wins / (losses + wins), wins, losses }
    }

    function calculateBalances (startBalance: number, positions: Position[]): LineChartEntry[] {
      if (!positions) return []

      let balance = startBalance
      const entries: LineChartEntry[] = []
      entries.push({
        x: new Date(props.backtest.fromDate),
        y: balance
      })
      for (const position of positions.sort((a, b) => new Date(a.closedAt).getTime() - new Date(b.closedAt).getTime())) {
        const pos = new BacktestPosition(position)
        if (new Date(position.closedAt).getTime() <= 0) continue

        balance = balance + pos.getQuoteDifference()

        entries.push({
          x: new Date(position.closedAt),
          y: balance
        })
      }

      return entries
    }

    const dataValues = computed(() => {
      const chartValues = calculateBalances(props.backtest.startBalance, props.backtest.positions)

      return {
        chartValues
      }
    })

    const chartData = computed<ChartData<'line'>>(() => ({
      datasets: [
        {
          data: dataValues.value.chartValues as any,
          label: 'Balance',
          borderColor: '#7F79FF',
          backgroundColor: '#7F79FF',
          pointBackgroundColor: 'transparent',
          pointBorderColor: 'transparent',
          animation: false
        }
      ]
    }))

    const options = computed(() => ({
      elements: {
        point: {
          radius: 8
        }
      },
      scales: {
        x: {
          grid: {
            borderColor: chartColors.value.xYAxis,
            tickColor: chartColors.value.xYAxis,
            color: chartColors.value.lines
          },
          type: 'time',
          ticks: {
            autoSkip: true,
            maxTicksLimit: 15
          }
        },
        y: {
          grid: {
            borderColor: chartColors.value.xYAxis,
            tickColor: chartColors.value.xYAxis,
            color: chartColors.value.lines
          }
        }
      }
    }))

    const { lineChartProps, lineChartRef } = useLineChart({
      chartData,
      options
    })

    return {
      lineChartProps,
      lineChartRef,
      backtestData,
      moment
    }
  }
})
</script>

<style lang="scss" scoped>
.backtest-result {
  .banner {
    padding: 2rem;
    border: 1px solid var(--border-color);
  }
  .wrapper {
    padding: 2rem;
    display: grid;
    grid-template-columns: 700px 1fr;
    gap: 2rem;
    @media(max-width: 1150px) {
      grid-template-columns: 1fr;
    }
  }
  p {
    margin-bottom: 0.5rem;
  }
  .result {
    span {
      color: green;
      font-weight: bold;
      font-size: 20px;
      margin-left: 1rem;
      &.negative {
        color: red;
      }
    }
  }
  .backtest-date {
    margin-top: 1rem;
    color: var(--text-tertiary);
  }
}
.chart {
  width: 100%;
}
</style>
