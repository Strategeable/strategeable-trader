<template>
  <div class="wrapper">
    <table cellspacing="0" cellpadding="0">
      <thead>
        <tr>
          <th>#</th>
          <th>Symbol</th>
          <th>Quote value</th>
          <th>Entry date</th>
          <th>Exit date</th>
          <th>Max drawdown</th>
          <th>PnL</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(entry, i) in actualPositions"
          :key="entry.id"
        >
          <td>{{ actualPositions.length - i }}</td>
          <td>{{ entry.symbol }}</td>
          <td>{{ Number(entry.getEntryQuoteSize().toFixed(2)) }}</td>
          <td>{{ moment(entry.openTime).format('DD-MM-YYYY HH:mm') }}</td>
          <td>{{ moment(entry.closeTime).format('DD-MM-YYYY HH:mm') }}</td>
          <td>{{ Number(entry.getMaxDrawdown().toFixed(2)) }}%</td>
          <td class="pnl" :class="{ negative: entry.getResultIncludingFees() < 0 }">{{ Number(entry.getResultIncludingFees().toFixed(2)) }}%</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType } from 'vue'
import moment from 'moment'
import Position from '@/types/Position'
import PositionHandler from '@/handlers/Position'

export default defineComponent({
  props: {
    positions: {
      type: Array as PropType<Position[]>,
      required: true
    }
  },
  setup (props) {
    const actualPositions = computed(() => props.positions
      .map(p => new PositionHandler(p))
      .sort((a, b) => new Date(b.closeTime || new Date()).getTime() - new Date(a.closeTime || new Date()).getTime())
    )
    return {
      moment,
      actualPositions
    }
  }
})
</script>

<style lang="scss" scoped>
table {
  border: 1px solid var(--border-color);
  width: 100%;
  max-width: 100%;
    min-width: 500px;
  thead {
    tr {
      background-color: var(--table-header);
      cursor: unset !important;
    }
    th {
      color: var(--text-secondary);
      background-color: var(--table-header);
    }
  }
  th, td {
    padding: 0.75rem 1.5rem;
    text-align: left;
    border: none;
    &:last-child {
      float: right;
    }
  }
  tr {
    background-color: var(--table-row);
    cursor: pointer;
  }
  tr:nth-child(even) {
    background-color: var(--table-row-alt);
  }
}

.pnl {
  color: var(--positive);
  &.negative {
    color: var(--negative);
  }
}

.wrapper {
  max-width: 100%;
  overflow-x: auto;
}
</style>
