<template>
  <table cellspacing="0" cellpadding="0">
    <thead>
      <tr>
        <th
          v-for="key in keys"
          :key="key"
        >{{ key }}</th>
      </tr>
    </thead>
    <tbody>
      <tr
        v-for="entry in data"
        :key="entry[dataKey]"
        @click="$emit('select', entry[dataKey])"
      >
        <td
          v-for="[key, value] in Object.entries(entry)"
          :key="key"
        >{{ value }}</td>
      </tr>
    </tbody>
  </table>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'

export default defineComponent({
  props: {
    data: {
      // eslint-disable-next-line
      type: Array as PropType<any[]>,
      required: true
    },
    dataKey: {
      type: String,
      required: true
    }
  },
  computed: {
    keys (): string[] {
      if (this.data.length === 0) return []
      return Object.keys(this.data[0])
    }
  }
})
</script>

<style lang="scss" scoped>
table {
  border: 1px solid var(--border-color);
  width: 100%;
  thead {
    tr {
      background-color: var(--table-header);
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
  }
  tr:nth-child(even) {
    background-color: var(--table-row-alt);
  }
}
</style>
