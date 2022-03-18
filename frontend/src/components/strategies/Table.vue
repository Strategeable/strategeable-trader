<template>
  <div class="wrapper">
    <table cellspacing="0" cellpadding="0">
      <thead>
        <tr>
          <th>Name</th>
          <th>Created at</th>
          <th>Last edited at</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="entry in sortedData"
          :key="entry.id"
          @click="$emit('select', entry.id)"
        >
          <td>{{ entry.name }}</td>
          <td>{{ moment(entry.createdAt).format('DD-MM-YYYY HH:mm') }}</td>
          <td>{{ moment(entry.lastEdited).format('DD-MM-YYYY HH:mm') }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType } from 'vue'
import moment from 'moment'
import { Strategy } from '@/types/Strategy'

export default defineComponent({
  props: {
    data: {
      // eslint-disable-next-line
      type: Array as PropType<Strategy[]>,
      required: true
    },
    dataKey: {
      type: String,
      required: true
    }
  },
  setup (props) {
    const sortedData = computed(() => props.data.sort((a, b) => new Date(b.lastEdited).getTime() - new Date(a.lastEdited).getTime()))

    return {
      moment,
      sortedData
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

.wrapper {
  max-width: 100%;
  overflow-x: auto;
}
</style>
