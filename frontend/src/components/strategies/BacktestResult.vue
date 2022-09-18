<template>
  <div class="backtest-result" :class="{ open }">
    <div class="banner">
      <div class="left">
        <span v-if="!isNaN(backtestData.change)" class="change" :class="{ negative: backtestData.change < 0 }">{{ backtestData.change }}%</span>
        <span v-else class="change">-</span>
        {{ backtest.strategy.name }}
      </div>
      <div class="right">
        <p>Backtested on {{ moment(backtest.startedOn).format('DD MMM HH:mm') }}</p>
        <fa-icon :icon="open ? 'caret-up' : 'caret-down'" />
      </div>
    </div>
    <div class="wrapper">
      <div class="inner">
        <div class="chart">
          <line-chart v-bind="lineChartProps"/>
        </div>
        <div class="data">
          <div>
            <p>{{ moment(backtest.fromDate).format('DD-MM-YYYY HH:mm') }} until {{ moment(backtest.toDate).format('DD-MM-YYYY HH:mm') }}</p>
            <p
              class="result"
            >
              Result
              <span v-if="!isNaN(backtestData.change)" :class="{ negative: backtestData.change < 0 }">{{ backtestData.change }}%</span>
              <span v-else>-</span>
            </p>
            <p>{{ backtestData.winsLosses.wins }} wins / {{ backtestData.winsLosses.losses }} losses (win rate: {{ Number((backtestData.winsLosses.winRate || 0).toFixed(2)) }})</p>
            <p>Max drawdown {{ calculateMaxDrawdown(backtestData.balances.map(x => x.y), true) }}%</p>
            <p>Max change {{ calculateMaxChange(backtestData.balances.map(x => x.y), backtest.startBalance, true) }}%</p>
          </div>
          <div class="bottom">
            <button v-if="!backtest.finished" @click="$emit('stop')">Stop backtest</button>
            <button class="outline" @click="$emit('restore')">Restore strategy</button>
            <button class="outline" @click="$emit('export')">Export strategy</button>
          </div>
        </div>
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
import { useStore } from '@/store'

interface LineChartEntry {
  y: number
  z: number
  x: Date
}

Chart.register(...registerables)
export default defineComponent({
  emits: ['restore', 'export', 'stop'],
  components: { LineChart },
  props: {
    backtest: {
      type: Object as PropType<BacktestResult>,
      required: true
    },
    open: {
      type: Boolean,
      required: true
    }
  },
  setup (props) {
    const store = useStore()
    const theme = computed(() => store.getters.theme)
    const chartColors = computed(() => {
      return {
        xYAxis: theme.value === 'light' ? 'black' : 'white',
        lines: theme.value === 'light' ? '#eeeeff' : '#1e1e2a'
      }
    })

    const quoteSymbol = props.backtest.strategy.symbols[0].split('/')[1]

    const backtestData = computed(() => ({
      change: Number(((props.backtest.endBalance - props.backtest.startBalance) / props.backtest.startBalance * 100).toFixed(2)),
      fromDate: new Date(props.backtest.fromDate),
      toDate: new Date(props.backtest.toDate),
      winsLosses: calculateWinRate(props.backtest.positions),
      balances: calculateBalances(props.backtest.startBalance, props.backtest.positions)
    }))

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
        y: balance,
        z: balance
      })

      // Cumulatively build up the historical balance of this backtest
      // by adding results of the positions saved for it in the database
      for (const position of positions.sort((a, b) => new Date(a.closedAt).getTime() - new Date(b.closedAt).getTime())) {
        const pos = new BacktestPosition(position)
        if (new Date(position.closedAt).getTime() <= 0) continue

        balance = balance + pos.getQuoteDifference()

        entries.push({
          x: new Date(position.closedAt),
          y: balance,
          z: 0
        })
      }

      return entries
    }

    function calculateMaxDrawdown (balances: number[], round?: boolean): number {
      if (balances.length <= 1) return 0

      let maxDrawdown = 0
      let currentMax = balances[0]
      let currentPrice: number
      let currentDrawdown: number

      for (let i = 1; i < balances.length; i++) {
        currentPrice = balances[i]
        currentDrawdown = (currentPrice - currentMax) / currentMax
        maxDrawdown = currentDrawdown < maxDrawdown ? currentDrawdown : maxDrawdown
        currentMax = currentPrice > currentMax ? currentPrice : currentMax
      }

      const result = maxDrawdown * 100
      if (round) return Number(result.toFixed(2))
      return result
    }

    function calculateMaxChange (balances: number[], startBalance: number, round?: boolean): number {
      if (balances.length <= 1) return 0
      const highestBalance = balances.reduce((acc, curr) => curr > acc ? curr : acc, 0)

      const change = (highestBalance - startBalance) / startBalance * 100
      if (round) return Number(change.toFixed(2))
      return change
    }

    const chartData = computed<ChartData<'line'>>(() => ({
      datasets: [
        {
          data: props.backtest.positions.map(p => ({ x: p.openedAt, y: p.entryValue.baseSize * p.entryValue.rate })) as any,
          label: 'Buy',
          borderColor: 'transparent',
          backgroundColor: 'transparent',
          pointBackgroundColor: 'green',
          pointBorderColor: 'green',
          pointRadius: 2,
          animation: false
        },
        {
          data: props.backtest.positions.map(p => ({ x: p.closedAt, y: p.exitValue.baseSize * p.exitValue.rate })) as any,
          label: 'Sell',
          borderColor: 'transparent',
          backgroundColor: 'transparent',
          pointBackgroundColor: 'red',
          pointBorderColor: 'red',
          pointRadius: 2,
          animation: false
        },
        {
          data: backtestData.value.balances as any,
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
      moment,
      quoteSymbol,
      calculateMaxDrawdown,
      calculateMaxChange
    }
  }
})
</script>

<style lang="scss" scoped>
.backtest-result {
  margin-bottom: 0.5rem;
  .banner {
    padding: 1rem;
    border: 1px solid var(--border-color);
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    transition: .2s ease-in-out;
    .change {
      color: var(--positive);
      &.negative {
        color: var(--negative);
      }
    }
    .left, .right {
      display: flex;
      align-items: center;
      p {
        font-size: 12px;
        color: var(--text-secondary);
        margin: 0;
        margin-right: 1rem;
      }
    }
    span {
      width: 90px;
      display: block;
      font-weight: bold;
    }
    &:hover {
      background-color: var(--background-lighten);
    }
  }
  .wrapper {
    max-height: 0;
    overflow-y: hidden;
    transition: .4s ease-in-out;
    .inner {
      padding: 2rem 0;
      display: grid;
      grid-template-columns: 700px 1fr;
      gap: 2rem;
    }
    @media(max-width: 1150px) {
      grid-template-columns: 1fr;
    }
  }
  p {
    margin-bottom: 0.5rem;
  }
  .result {
    span {
      color: var(--positive);
      font-weight: bold;
      font-size: 20px;
      margin-left: 1rem;
      &.negative {
        color: var(--negative);
      }
    }
  }
  &.open {
    .wrapper {
      max-height: 800px;
    }
    .banner {
      background-color: var(--background-lighten);
    }
  }
}
.chart {
  padding: 1rem;
  width: 100%;
  background-color: var(--background-lighten);
  border: 1px solid var(--border-color);
}
.data {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  .bottom {
    button {
      margin-right: 1rem;
    }
  }
}
</style>
