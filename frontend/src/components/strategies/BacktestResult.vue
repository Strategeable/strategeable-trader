<template>
  <div class="backtest-result">
    Backtest result - from: {{ backtest.startBalance }} to: {{ backtest.endBalance }}
    <line-chart v-bind="lineChartProps"/>
  </div>
</template>

<script lang="ts">
import { BacktestResult, Position } from '@/types/Backtest'
import { computed, defineComponent, PropType } from 'vue'
import { LineChart, useLineChart } from 'vue-chart-3'
import { Chart, ChartData, registerables } from 'chart.js'
import 'chartjs-adapter-moment'

interface LineChartEntry {
  y: number
  x: Date
}

Chart.register(...registerables)
export default defineComponent({
  components: { LineChart },
  props: {
    backtest: {
      type: Object as PropType<BacktestResult>,
      required: true
    }
  },
  setup (props) {
    function calculateBalances (startBalance: number, positions: Position[]): LineChartEntry[] {
      if (!positions) return []

      let balance = startBalance
      const entries: LineChartEntry[] = []
      for (const position of positions) {
        const quoteEntrySize = position.entryValue.rate * position.entryValue.baseSize
        const quoteExitSize = position.exitValue.rate * position.exitValue.baseSize
        const difference = quoteExitSize - quoteEntrySize

        balance = balance + difference - position.entryValue.quoteFees - position.exitValue.quoteFees

        if (new Date(position.closedAt).getTime() < 37701695) continue

        entries.push({
          x: new Date(position.closedAt),
          y: balance
        })
      }

      return entries
    }

    const dataValues = computed(() => calculateBalances(props.backtest.startBalance, props.backtest.positions))

    const chartData = computed<ChartData<'line'>>(() => ({
      datasets: [
        {
          data: dataValues.value as any,
          label: 'Balance',
          borderColor: '#3781d6'
        }
      ]
    }))

    const { lineChartProps, lineChartRef } = useLineChart({
      chartData,
      options: {
        elements: {
          point: {
            radius: 0
          }
        },
        scales: {
          x: {
            type: 'time',
            ticks: {
              autoSkip: true,
              maxTicksLimit: 15
            }
          }
        }
      }
    })

    return {
      lineChartProps,
      lineChartRef
    }
  }
})
</script>
